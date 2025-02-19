//go:generate packer-sdc mapstructure-to-hcl2 -type Config
package artifactDownloadOther

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
	OutputDir			   string `mapstructure:"output_dir" required:"true"`
	ArtifactoryPath        string `mapstructure:"artifactory_path" required:"true"`
	FileList               []string `mapstructure:"file_list" required:"true"`
}

type Datasource struct {
	config Config
}

// --> If making changes to this section, make sure the hcl2spec gets updated as well!
type DatasourceOutput struct {}

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
			log.Fatal("Please provide an Artifactory Identity Token.")
		}
	}
	
	if d.config.ArtifactoryServer == "" {
		server := os.Getenv("ARTIFACTORY_SERVER")
		if server == "" {
			log.Fatal("Please provide the URL to the Artifactory server (ex: https://server.com:8081/artifactory/api).")
		}
	}

	if d.config.OutputDir == "" {
		log.Println("Please provide the path to the desired output directory for the files being downloaded.")
		log.Fatal("Path should include proper escape characters where necessary.")
	}

	if d.config.ArtifactoryPath == "" {
		log.Fatal("Please provide the repo path in Artifactory where the artifacts reside (ex: /repo/folder/).")
	}

	if len(d.config.FileList) <= 0 {
		log.Println("Please provide a list of one or more filenames to be downloaded.")
		log.Fatal("Ex:  file_list = [\"file1.txt\", \"file2.txt\"]")
	}

	return nil
}

func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Execute() (cty.Value, error) {
	var serverApi, token, artifPath, outputDir, result string
	var fileList []string
	var err error

	// Environment related
	if d.config.AritfactoryToken == "" {
		token = os.Getenv("ARTIFACTORY_TOKEN")
		if token != "" {
			d.config.AritfactoryToken = token
		}
	}
	
	if d.config.ArtifactoryServer == "" {
		serverApi = os.Getenv("ARTIFACTORY_SERVER")
		if serverApi != "" {
			d.config.ArtifactoryServer = serverApi
		}
	}

	// Artifact Related
	if d.config.OutputDir == "" {
		outputDir = os.Getenv("OUTPUTDIR")
		if outputDir != "" {
			d.config.OutputDir = outputDir
		}
	}

	if d.config.ArtifactoryPath != "" {
		artifPath = d.config.ArtifactoryPath
	}

	if len(d.config.FileList) != 0 {
		fileList = d.config.FileList
	}

	dirCheck, err := os.Stat(outputDir)
	log.Println("Checking if output directory exists... ", dirCheck)
	if os.IsNotExist(err) {
		err = os.MkdirAll(outputDir, 0755)
		if err != nil {
			log.Println("Error creating directory: " + outputDir + " - ", err)
		} else {
			log.Println("Successfully created directory: " + outputDir)
		}
	}

	for _, file := range fileList {
		task := "Downloading: " + file
		result, err = tasks.DownloadGeneralArtifact(serverApi, token, outputDir, artifPath, file, task)
		log.Println("Result of download: " + file + " - " + result)
	}

	if err != nil {
		log.Println("There were errors downloading one or more files")
		log.Println(err)
	}

	output := DatasourceOutput{}

	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}