//go:generate packer-sdc mapstructure-to-hcl2 -type Config
package artifactUpload

import (
	"context"
	"log"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/common"
	packersdk "github.com/hashicorp/packer-plugin-sdk/packer"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	artifactorysdk "github.com/raynaluzier/artifactory-go-sdk/common"
	"github.com/raynaluzier/artifactory-go-sdk/tasks"
)

type Config struct {
	SourcePath			   string `mapstructure:"source_path" required:"true"`
	// If not provided, then can reference an existing artifact URI to parse for the target
	TargetPath			   string `mapstructure:"target_path" required:"false"`  // either this or existing uri target
	// Optional for potential distinguishing values such as version, date, etc where the image name is always the same
	// Will use '-' as a separator; if blank, will be ignored
	FileSuffix			   string `mapstructure:"file_suffix" required:"false"`
	ExistingUriTarget	   string `mapstructure:"existing_uri_target" required:"false"`

	common.PackerConfig	  `mapstructure:",squash"`

	AritfactoryToken       string `mapstructure:"artifactory_token" required:"true"`
	ArtifactoryServer      string `mapstructure:"artifactory_server" required:"true"`
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

	if p.config.SourcePath == "" {
		log.Fatal("Please provide the source path to the artifact to upload.")
		log.Fatal("Source path should be in the form of either 'h:\\lab\\artifact.ext' or '/lab/artifact.ext'")
	}

	if p.config.TargetPath == "" && p.config.ExistingUriTarget == "" {
		log.Fatal("Please provide either a target path or an existing artifact URI (from data source) to reference as a target location.")
		log.Fatal("If using an existing artifact URI, the artifact's path will be parsed and used as the target for the new artifact.")
		log.Fatal("Otherwise, the target path should be in the form of '/repo/folder/path")
	}

	if p.config.TargetPath != "" && p.config.ExistingUriTarget != "" {
		log.Println("Values have been provided for both the target path AND an existing URI target.")
		log.Println("The value provided in the target path will be used.")
	}

	return nil
}

func (p *PostProcessor) PostProcess(ctx context.Context, ui packersdk.Ui, source packersdk.Artifact) (packersdk.Artifact, bool, bool, error) {
	var sourcePath, targetPath, fileSuffix string
	//token     := p.config.PackerConfig.PackerUserVars["artifactory_token"]
	//serverApi := p.config.PackerConfig.PackerUserVars["artifactory_server"]

	token := p.config.AritfactoryToken			// testing
	serverApi := p.config.ArtifactoryServer		// testing

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

	if p.config.FileSuffix == "" {
		fileSuffix = ""
	} else {
		fileSuffix = p.config.FileSuffix
	}

	downloadUri, artifactUri, err := tasks.UploadArtifact(serverApi, token, sourcePath, targetPath, fileSuffix)
	
	if err != nil {
		log.Fatal("Unable to upload the artifact - ", err)
	} else {
		ui.Say("Artifact upload completed.")
		log.Println("Download URI for new artifact: " + downloadUri)
		log.Println("Artifact URI for new artifact: " + artifactUri)
	}

	return source, true, true, nil
}