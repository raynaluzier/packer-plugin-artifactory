package artifactImage

import (
	"fmt"
	"os"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	"github.com/zclconf/go-cty/cty"

	"github.com/raynaluzier/go-artifactory/common"
	"github.com/raynaluzier/go-artifactory/operations"
	"github.com/raynaluzier/go-artifactory/search"
)

type Config struct {
	AritfactoryToken       string `mapstructure:"artifactory_token" required:"true"`
	ArtifactoryServer      string `mapstructure:"artifactory_server" required:"true"`
	// Defaults to user's home dir if blank
	ArtifactoryOutputDir   string `mapstructure:"artifactory_output_dir" required:"false"`
	// Defaults to 'INFO'
	ArtifactoryLogging     string `mapstructure:"artifactory_logging" required:"false"`

	// Full or partial name of the artifact
	ArtifactName           string `mapstructure:"artifact_name" required:"true"`
	// File extension; defaults to '.vmxt' if left blank
	ArtifactFileType       string `mapstructure:"file_type" required:"true"`
	// Channel is technically a property; if it exists, will be appended to the kvProperties []string
	ArtifactChannel        string `mapstructure:"channel" required:"false"`
	// Key/value pairs of properties to filter on
	ArtifactFilter         map[string]string `mapstructure:"filter" required:"false"`
}

type Datasource struct {
	config Config
}

type DatasourceOutput struct {
	Name        string `mapstructure:"name"`
	Created     string `mapstructure:"creation_date"`
	ArtifactUri	string `mapstructure:"artifcat_uri"`
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
			fmt.Errorf("Please provide an Artifactory Identity Token.")
		}
	}
	
	if d.config.ArtifactoryServer == "" {
		server := os.Getenv("ARTIFACTORY_SERVER")
		if server == "" {
			fmt.Errorf("Please provide the URL to the Artifactory server (ex: https://server.com:8081/artifactory/api).")
		}
	}

	if d.config.ArtifactName == "" {
		fmt.Errorf("Please provide the full or partial artifact name.")
	}

	if d.config.ArtifactFileType == "" {
		fmt.Errorf("Please provide the source image's extension type; for example .vmxt.")
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
	var strErr string

	// Environment related
	if d.config.AritfactoryToken == "" {
		token := os.Getenv("ARTIFACTORY_TOKEN")
		if token != "" {
			d.config.AritfactoryToken = token
		}
	}
	
	if d.config.ArtifactoryServer == "" {
		server := os.Getenv("ARTIFACTORY_SERVER")
		if server != "" {
			d.config.ArtifactoryServer = server
		}
	}

	if d.config.ArtifactoryOutputDir == "" {
		outputDir := os.Getenv("ARTIFACTORY_OUTPUTDIR")
		if outputDir != "" {
			d.config.ArtifactoryOutputDir = outputDir
		}
	} // If outputDir is still "", the user's home dir will be used

	if d.config.ArtifactoryLogging == "" {
		logLevel := os.Getenv("ARTIFACTORY_LOGGING")
		if logLevel == "" {
			d.config.ArtifactoryLogging = "INFO"
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

	// Search for artifact by name
	listArtifacts, err := search.GetArtifactsByName(artifName)
	if err != nil {
		strErr = fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Error getting list of matching artifacts - " + strErr)
	}
	
	listByFileType, err := search.FilterListByFileType(ext, listArtifacts)
	if err != nil {
		strErr = fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Error filtering artifacts by file type - " + strErr)
	}

	if len(listByFileType) == 1 {
		artifactUri = listByFileType[0]
	} else if len(listByFileType) > 1 && len(kvProperties) != 0 {
		artifactUri, err = operations.FilterListByProps(listByFileType, kvProperties)
		if err != nil {
			strErr = fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error("Error filtering artifacts by properties - " + strErr)
		}
	} else {
		// If no props but more than one artif in list, return latest
		artifactUri, err = operations.GetLatestArtifactFromList(listByFileType)
		if err != nil {
			strErr = fmt.Sprintf("%v\n", err)
			common.LogTxtHandler().Error("Error getting latest artifact from list - " + strErr)
		}
	}

	artifactName := operations.GetArtifactNameFromUri(artifactUri)

	// Get other useful artifact details
	createDate, err := operations.GetCreateDate(artifactUri)
	if err != nil {
		strErr = fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Unable to get create date of artifact - " + strErr)
	}
	
	downloadUri, err := operations.GetDownloadUri(artifactUri)
	if err != nil {
		strErr = fmt.Sprintf("%v\n", err)
		common.LogTxtHandler().Error("Unable to get download URI - " + strErr)
	}
	
	output := DatasourceOutput{
		Name: 			artifactName,
		Created: 		createDate,
		ArtifactUri: 	artifactUri,
		DownloadUri: 	downloadUri,
	}

	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}