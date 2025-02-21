//go:generate packer-sdc mapstructure-to-hcl2 -type Config
package artifactUpdateProps

import (
	"context"
	"log"
	"os"

	"github.com/hashicorp/hcl/v2/hcldec"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/raynaluzier/artifactory-go-sdk/tasks"
)

type Config struct {
	AritfactoryToken       string `mapstructure:"artifactory_token" required:"true"`
	ArtifactoryServer      string `mapstructure:"artifactory_server" required:"true"`
	ArtifactUri			   string `mapstructure:"artifact_uri" required:"true"`
	ArtifactProperties	   map[string]string `mapstructure:"properties" required:"true"`
	// Defaults to INFO
	Logging                string `mapstructure:"logging" required:"false"`
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

	if p.config.AritfactoryToken == "" {
		log.Fatal("---> Missing Artifactory identity token. The token is required to complete tasks against Artifactory.")
	}

	if p.config.ArtifactoryServer == "" {
		log.Fatal("---> Missing Artifactory server API address. The server API address is required to communicate with Artifactory.")
	}

	if p.config.ArtifactUri == "" {
		log.Fatal("---> Missing Artifact URI. The new Artifact URI is required to update the artifact's properties.")
	}

	if len(p.config.ArtifactProperties) == 0 {
		log.Fatal("---> Missing Artifact properties. At least one key/value pair is required to update the artifact's properties.")
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
	var kvProperties []string
	var token, serverApi, artifactUri, logLevel string

	if p.config.AritfactoryToken != "" {
		token = p.config.AritfactoryToken
	}
	
	if p.config.ArtifactoryServer != "" {
		serverApi = p.config.ArtifactoryServer
	}

	if p.config.ArtifactUri != "" {
		artifactUri = p.config.ArtifactUri
	}

	if len(p.config.ArtifactProperties) != 0 {
		kvProperties = BuildProps(p.config.ArtifactProperties)
	}

	if p.config.Logging == "" {
		logLevel := os.Getenv("LOGGING")
		if logLevel != "" {
			p.config.Logging = logLevel
		}
	}

	statusCode, err := tasks.SetProps(serverApi, token, logLevel, artifactUri, kvProperties)
	if statusCode == "204" {
		ui.Say("Property assignment to artifact was successful.")
	} else {
		log.Fatal("Unable to update the artifact properties - ", err)
	}

	return source, true, true, nil
}