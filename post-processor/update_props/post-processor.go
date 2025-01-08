//go:generate packer-sdc mapstructure-to-hcl2 -type Config
package artifactUpdateProps

import (
	"context"
	"log"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/raynaluzier/artifactory-go-sdk/tasks"
)

type Config struct {
	common.PackerConfig	  `mapstructure:",squash"`
	ArtifactUri			   string `mapstructure:"artifact_uri" required:"true"`
	ArtifactProperties	   map[string]string `mapstructure:"properties" required:"true"`
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

	if p.config.ArtifactUri == "" {
		log.Fatal("Missing Artifact URI. The new Artifact URI is required to update the artifact's properties.")
	}

	if len(p.config.ArtifactProperties) == 0 {
		log.Fatal("Missing Artifact properties. At least one key/value pair is required to update the artifact's properties.")
	}

	return nil
}

func BuildProps(kvInput map[string]string) ([]string) {
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
	token     := p.config.PackerConfig.PackerUserVars["artifactory_token"]
	serverApi := p.config.PackerConfig.PackerUserVars["artifactory_server"]
	var kvProperties []string
	var artifactUri string

	if p.config.ArtifactUri != "" {
		artifactUri = p.config.ArtifactUri
	}

	if len(p.config.ArtifactProperties) != 0 {
		kvProperties = BuildProps(p.config.ArtifactProperties)
	}

	statusCode, err := tasks.SetProps(serverApi, token, artifactUri, kvProperties)
	if statusCode == "204" {
		ui.Say("Property assignment to artifact was successful.")
	} else {
		log.Fatal("Unable to update the artifact properties - ", err)
	}

	return source, true, true, nil
}