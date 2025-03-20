//go:generate packer-sdc mapstructure-to-hcl2 -type Config
package artifactUploadOther

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
	ArtifactoryToken       string `mapstructure:"artifactory_token" required:"true"`
	ArtifactoryServer      string `mapstructure:"artifactory_server" required:"true"`
	SourcePath			   string `mapstructure:"source_path" required:"true"`
	ArtifactoryPath		   string `mapstructure:"artifactory_path" required:"true"`
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

	if p.config.ArtifactoryToken == "" {
		token := os.Getenv("ARTIFACTORY_TOKEN")
		if token == "" {
			log.Fatal("---> Missing Artifactory identity token. The token is required to complete tasks against Artifactory.")
		}	
	}

	if p.config.ArtifactoryServer == "" {
		serverApi := os.Getenv("ARTIFACTORY_SERVER")
		if serverApi == "" {
			log.Fatal("---> Missing Artifactory server API address. The server API address is required to communicate with Artifactory.")	
		}
	}

	if p.config.SourcePath == "" {
		log.Fatal("---> Please provide the source path to the artifact(s) to upload.")
		log.Fatal("Source path should be in the form of either 'h:\\lab\\' or '/lab/'")
	}

	if p.config.ArtifactoryPath == "" {
		log.Fatal("---> Please provide the Artifactory /repo/folder/path where the artifact(s) should be uploaded to.")
	}

	if len(p.config.FileList) <= 0 {
		log.Fatal("---> Please add one or more files to the file_list for upload.")
	}

	return nil
}

func (p *PostProcessor) PostProcess(ctx context.Context, ui packersdk.Ui, source packersdk.Artifact) (packersdk.Artifact, bool, bool, error) {
	var token, serverApi, sourcePath, artifPath, result string
	var err error
	var fileList, failList []string

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

	if p.config.SourcePath != "" {
		sourcePath = p.config.SourcePath
	}

	if p.config.ArtifactoryPath != "" {
		artifPath = p.config.ArtifactoryPath
	}

	if len(p.config.FileList) > 0 {
		fileList = p.config.FileList
	}

	for _, file := range fileList {
		result, err = tasks.UploadGeneralArtifact(serverApi, token, sourcePath, artifPath, file)
		if result == "Failed" && err == nil {
			log.Println("File not found: " + file)
			failList = append(failList, file)
		} else if err != nil {
			log.Println("Error uploading: " + file)
			failList = append(failList, file)
		} else {
			log.Println("Successfully uploaded file: " + file)
		}
	}

	if len(failList) > 0 {
		log.Println("[WARN] Unable to upload one or more artifacts.")
		log.Println("[WARN] Please check the file name(s) and Artifactory path provided; both are CASE-SENSITIVE.")
		log.Println("---> The following files were not uploaded: ")
		for _, f := range failList {
			log.Println(f)
		}
		err := errors.New("Unable to upload one or more artifacts.")
		return source, false, false, err

	} else {
		ui.Say("Artifact upload(s) complete.")
		return source, true, true, nil
	}

}