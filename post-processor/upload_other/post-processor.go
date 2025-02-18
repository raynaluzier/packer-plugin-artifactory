//go:generate packer-sdc mapstructure-to-hcl2 -type Config
package artifactUpload

import (
	"context"
	"errors"
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
	SourcePath			   string `mapstructure:"source_path" required:"true"`
	ArtifactoryPath		   string `mapstructure:"artifactory_path" required:"true"`
	FolderName             string `mapstructure:"folder_name" required:"false"`
	FileList         	   []string `mapstructure:"file_list" required:"true"`
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
		token := os.Getenv("ARTIFACTORY_TOKEN")
		if token == "" {
			log.Fatal("Missing Artifactory identity token. The token is required to complete tasks against Artifactory.")
		}	
	}

	if p.config.ArtifactoryServer == "" {
		serverApi := os.Getenv("ARTIFACTORY_SERVER")
		if serverApi == "" {
			log.Fatal("Missing Artifactory server API address. The server API address is required to communicate with Artifactory.")	
		}
	}

	if p.config.SourcePath == "" {
		log.Fatal("Please provide the source path to the artifact(s) to upload.")
		log.Fatal("Source path should be in the form of either 'h:\\lab\\' or '/lab/'")
	}

	if p.config.ArtifactoryPath == "" {
		log.Fatal("Please provide the Artifactory /repo/folder/path where the artifact(s) should be uploaded to.")
	}

	if p.config.FolderName == "" {
		log.Println("No folder name was provided. Therefore, artifact(s) will be placed in the root of the Artifactory path provided.")
		log.Println("If this is not desired, please provide an folder name. To place these with a specific image file, use the image name as the folder name.")
	}

	if len(p.config.FileList) <= 0 {
		log.Fatal("Please add one or more files to the file_list for upload.")
	}

	return nil
}

func (p *PostProcessor) PostProcess(ctx context.Context, ui packersdk.Ui, source packersdk.Artifact) (packersdk.Artifact, bool, bool, error) {
	var token, serverApi, sourcePath, artifPath, folderName, result string
	var err error
	var fileList []string

	if p.config.AritfactoryToken != "" {
		token = p.config.AritfactoryToken
	} else {
		token = os.Getenv("ARTIFACTORY_TOKEN")
	}
	
	if p.config.ArtifactoryServer != "" {
		serverApi = p.config.ArtifactoryServer
	} else {
		serverApi = os.Getenv("ARTIFACTORY_SERVER")
	}

	if p.config.SourcePath != "" {
		sourcePath = p.config.SourcePath
	}

	if p.config.ArtifactoryPath != "" {
		artifPath = p.config.ArtifactoryPath
	}

	if p.config.FolderName == "" {
		folderName = ""
	} else {
		folderName = p.config.FolderName
	}

	if len(p.config.FileList) > 0 {
		fileList = p.config.FileList
	}

	for _, file := range fileList {
		result, err = tasks.UploadGeneralArtifact(serverApi, token, sourcePath, artifPath, file, folderName)
		if err != nil {
			log.Println(result)
			log.Fatal("Error uploading: " + file)
		} else {
			log.Println("Successfully uploaded file: " + file)
		}
	}

	if err != nil {
		log.Fatal("Unable to upload one or more artifacts.")
		err := errors.New("Unable to upload artifacts.")
		return source, false, false, err

	} else {
		ui.Say("Artifact upload(s) completed.")
		return source, true, true, nil
	}

}