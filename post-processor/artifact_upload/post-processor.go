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
	artifactorysdk "github.com/raynaluzier/artifactory-go-sdk/common"
	"github.com/raynaluzier/artifactory-go-sdk/tasks"
)

type Config struct {
	ArtifactoryToken       string `mapstructure:"artifactory_token" required:"true"`
	ArtifactoryServer      string `mapstructure:"artifactory_server" required:"true"`
	SourcePath			   string `mapstructure:"source_path" required:"true"`
	// If not provided, then can reference an existing artifact URI to parse for the target
	TargetPath			   string `mapstructure:"target_path" required:"false"`  // either this or existing uri target
	// Valid values are "ova", "ovf", and "vmtx"
	ImageType              string `mapstructure:"image_type" required:"true"`
	// Base image name (ex: win2022 or rhel9)
	ImageName              string `mapstructure:"image_name" required:"true"`
	ExistingUriTarget	   string `mapstructure:"existing_uri_target" required:"false"`
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
		log.Println("---> Please provide the source path to the artifact to upload.")
		log.Fatal("Source path should be in the form of either 'h:\\lab\\artifact.ext' or '/lab/artifact.ext'")
	}

	if p.config.TargetPath == "" && p.config.ExistingUriTarget == "" {
		log.Println("---> Please provide either a target path OR use the 'existing_artifact_uri' input (which can be populated manually or from data source) to reference as a target location.")
		log.Println("If using an existing artifact URI, the artifact's path will be parsed and used as the target for the new artifact.")
		log.Fatal("Otherwise, the target path should be in the form of '/repo/folder/path.")
	}

	if p.config.TargetPath != "" && p.config.ExistingUriTarget != "" {
		log.Println("Values have been provided for both the target path AND an existing URI target.")
		log.Println("The value provided in the target path will be used.")
	}

	if p.config.ImageType == "" {
		log.Fatal("---> Please provide the image file type that will be uploaded: 'ova', 'ovf', or 'vmtx'.")
	}

	if p.config.ImageName == "" {
		log.Fatal("---> Please provide the name of the image; examples: win2022, rhel9, win22_25_01_25...")
	}

	return nil
}

func (p *PostProcessor) PostProcess(ctx context.Context, ui packersdk.Ui, source packersdk.Artifact) (packersdk.Artifact, bool, bool, error) {
	var token, serverApi, sourcePath, targetPath, imageType, imageName string

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

	if p.config.TargetPath != "" && p.config.ExistingUriTarget == "" {
		targetPath = p.config.TargetPath
	} else if p.config.TargetPath != "" && p.config.ExistingUriTarget != "" {
		targetPath = p.config.TargetPath
	} else if p.config.TargetPath == "" && p.config.ExistingUriTarget != "" {
		targetPath = artifactorysdk.ParseArtifUriForPath(serverApi, p.config.ExistingUriTarget)
	}

	if p.config.ImageType != "" {
		imageType = p.config.ImageType
	}

	if p.config.ImageName != "" {
		imageName = p.config.ImageName
	}

	log.Println("Server Address: " + serverApi)
	log.Println("Image Name: " + imageName)
	log.Println("Image Type: " + imageType)
	log.Println("Source Path: " + sourcePath)
	log.Println("Target Path: " + targetPath)

	log.Println("Preparing to check and upload image artifact(s)...")
	result := tasks.UploadArtifacts(serverApi, token, imageType, imageName, sourcePath, targetPath)

	if result != "End of upload process" {
		log.Println("Unable to upload artifacts - " + result)
		log.Println("The Artifactory path is CASE SENSITIVE. Please verify the Artifactory path is correct, image type is either OVA/OVF/VMTX, and the image file(s) exist in the source path.")
		err := errors.New("Unable to upload image artifacts.")
		return source, false, false, err

	} else {
		log.Println(result + " - Uploaded image artifacts for: " + imageName + "." + imageType)
		log.Println("---> Upload of image artifact(s) completed.")
		return source, true, true, nil
	}

}