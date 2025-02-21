//go:generate packer-sdc mapstructure-to-hcl2 -type Config,DatasourceOutput
package artifactImage

import (
	"log"
	"os"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/raynaluzier/artifactory-go-sdk/tasks"
	"github.com/zclconf/go-cty/cty"
)

// --> If making changes to this section, make sure the hcl2spec gets updated as well!
type Config struct {
	AritfactoryToken       string `mapstructure:"artifactory_token" required:"true"`
	ArtifactoryServer      string `mapstructure:"artifactory_server" required:"true"`
	// Full or partial name of the artifact
	ArtifactName           string `mapstructure:"artifact_name" required:"true"`
	// File extension; defaults to '.vmtx' if left blank
	ArtifactFileType       string `mapstructure:"file_type" required:"true"`
	// Channel is technically a property; if it exists, will be appended to the kvProperties []string
	ArtifactChannel        string `mapstructure:"channel" required:"false"`
	// Key/value pairs of properties to filter on
	ArtifactFilter         map[string]string `mapstructure:"filter" required:"false"`
}

type Datasource struct {
	config Config
}

// --> If making changes to this section, make sure the hcl2spec gets updated as well!
type DatasourceOutput struct {
	Name        string `mapstructure:"name"`
	Created     string `mapstructure:"creation_date"`
	ArtifactUri	string `mapstructure:"artifact_uri"`
	DownloadUri string `mapstructure:"download_uri"`
}

func (d *Datasource) ConfigSpec() hcldec.ObjectSpec { 
	return d.config.FlatMapstructure().HCL2Spec() 
}

func (d *Datasource) Configure(raws ...interface{}) error {
	err := config.Decode(&d.config, nil, raws...)
	if err != nil {
		return err
	}

	if d.config.AritfactoryToken == "" {
		token := os.Getenv("ARTIFACTORY_TOKEN")
		if token == "" {
			log.Fatal("---> Please provide an Artifactory Identity Token.")
		}
	}
	
	if d.config.ArtifactoryServer == "" {
		server := os.Getenv("ARTIFACTORY_SERVER")
		if server == "" {
			log.Fatal("---> Please provide the URL to the Artifactory server (ex: https://server.com:8081/artifactory/api).")
		}
	}

	if d.config.ArtifactName == "" {
		log.Fatal("---> Please provide the full or partial artifact name.")
	}

	if d.config.ArtifactFileType == "" {
		log.Fatal("---> Please provide the source image's extension type; for example '.vmtx' or 'vmtx'.")
	}
	return nil
}

func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func BuildPropFilters(kvInput map[string]string) ([]string) {
	var filterOptions []string
	if len(kvInput) > 0 {
		for key, value := range kvInput {
			filter := key + "=" + value
			filterOptions = append(filterOptions, filter)
		}
	}
	return filterOptions
}

func (d *Datasource) Execute() (cty.Value, error) {
	var artifName string
	var ext string
	var kvProperties []string
	var artifactUri string

	// Environment related
	if d.config.AritfactoryToken == "" {
		token := os.Getenv("ARTIFACTORY_TOKEN")
		if token != "" {
			d.config.AritfactoryToken = token
		}
	}
	
	if d.config.ArtifactoryServer == "" {
		serverApi := os.Getenv("ARTIFACTORY_SERVER")
		if serverApi != "" {
			d.config.ArtifactoryServer = serverApi
		}
	}

	// Artifact Related
	if d.config.ArtifactName != "" {
		artifName = d.config.ArtifactName
	}

	if d.config.ArtifactFileType != "" {
		ext = d.config.ArtifactFileType
	}

	if len(d.config.ArtifactFilter) != 0 {
		kvProperties = BuildPropFilters(d.config.ArtifactFilter)
	}

	// Channel is technically a property
	if d.config.ArtifactChannel != "" {
		channelProp := "channel=" + d.config.ArtifactChannel
		kvProperties = append(kvProperties, channelProp)
	}

	// Search for artifact and return details
	artifactUri, artifactName, createDate, downloadUri, err := tasks.GetImageDetails(d.config.ArtifactoryServer, d.config.AritfactoryToken, artifName, ext, kvProperties)
	if err != nil {
		log.Println(err)
		if artifactUri == "" {
			log.Println("[ERROR] ----> No Artifact URI was returned; this is likely because no matching artifact was found.")
		} else if artifactName == "" || createDate == "" || downloadUri == "" {
			log.Println("[ERROR] ----> An Artifact URI was returned, however an error was encountered retrieving one or more artifact details.")
		}
	}
	
	output := DatasourceOutput{
		Name: 	artifactName,
		Created: 	createDate,
		ArtifactUri: 	artifactUri,
		DownloadUri: 	downloadUri,
	}

	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}