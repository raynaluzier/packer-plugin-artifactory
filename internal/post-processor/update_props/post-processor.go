//go:generate packer-sdc mapstructure-to-hcl2 -type Config
package artifactUpdateProps

import (
	"context"
	"errors"
	"os"

	"github.com/hashicorp/hcl/v2/hcldec"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/raynaluzier/artifactory-go-sdk/tasks"
)

type Config struct {
	ArtifactoryToken   string            `mapstructure:"artifactory_token" required:"true"`
	ArtifactoryServer  string            `mapstructure:"artifactory_server" required:"true"`
	ArtifactUri        string            `mapstructure:"artifact_uri" required:"true"`
	ArtifactProperties map[string]string `mapstructure:"properties" required:"true"`
}

type PostProcessor struct {
	config Config
}

func (p *PostProcessor) ConfigSpec() hcldec.ObjectSpec { return p.config.FlatMapstructure().HCL2Spec() }

func (p *PostProcessor) Configure(raws ...interface{}) error {
	err := config.Decode(&p.config, nil, raws...)
	if err != nil {
		return err
	}

	if p.config.ArtifactoryToken == "" {
		return errors.New("missing Artifactory identity token")
	}

	if p.config.ArtifactoryServer == "" {
		return errors.New("missing Artifactory server API address")
	}

	if p.config.ArtifactUri == "" {
		return errors.New("artifact_uri is required")
	}

	if len(p.config.ArtifactProperties) == 0 {
		return errors.New("artifact properties must contain at least one key/value pair")
	}

	return nil
}

func BuildProps(kvInput map[string]string) []string {
	var props []string
	if len(kvInput) > 0 {
		for key, value := range kvInput {
			kvPair := key + "=" + value
			props = append(props, kvPair)
		}
	}
	return props
}

func (p *PostProcessor) PostProcess(ctx context.Context, ui packersdk.Ui, source packersdk.Artifact) (packersdk.Artifact, bool, bool, error) {
	var kvProperties []string
	var token, serverApi, artifactUri string

	if p.config.ArtifactoryToken == "" {
		token = os.Getenv("ARTIFACTORY_TOKEN")
	} else {
		token = p.config.ArtifactoryToken
	}

	if p.config.ArtifactoryServer == "" {
		serverApi = os.Getenv("ARTIFACTORY_SERVER")
	} else {
		serverApi = p.config.ArtifactoryServer
	}

	if p.config.ArtifactUri != "" {
		artifactUri = p.config.ArtifactUri
	}

	if len(p.config.ArtifactProperties) != 0 {
		kvProperties = BuildProps(p.config.ArtifactProperties)
	}

	statusCode, err := tasks.SetProps(serverApi, token, artifactUri, kvProperties)
	if err != nil {
		return source, false, false, err
	}
	if statusCode != "204" {
		return source, false, false, errors.New("unable to update artifact properties")
	}

	ui.Say("Property assignment to artifact was successful.")
	return source, true, true, nil
}
