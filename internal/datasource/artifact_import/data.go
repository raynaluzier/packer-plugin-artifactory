//go:generate packer-sdc mapstructure-to-hcl2 -type Config,DatasourceOutput
package artifactImport

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	artifCommon "github.com/raynaluzier/artifactory-go-sdk/common"
	artifTasks "github.com/raynaluzier/artifactory-go-sdk/tasks"
	vsTasks "github.com/raynaluzier/vsphere-go-sdk/tasks"
	"github.com/zclconf/go-cty/cty"
)

// --> If making changes to this section, make sure the hcl2spec gets updated as well!
type Config struct {
	ArtifactoryToken  string `mapstructure:"artifactory_token" required:"true"`
	ArtifactoryServer string `mapstructure:"artifactory_server" required:"true"`

	VcenterServer     string `mapstructure:"vcenter_server" required:"true"`
	VcenterUser       string `mapstructure:"vcenter_user" required:"true"`
	VcenterPassword   string `mapstructure:"vcenter_password" required:"true"`
	VcenterDatacenter string `mapstructure:"datacenter_name" required:"true"`
	VcenterDatastore  string `mapstructure:"datastore_name" required:"true"`
	VcenterCluster    string `mapstructure:"cluster_name" required:"true"`
	// Will use default root folder of datacenter if left blank
	VcenterFolder string `mapstructure:"folder_name" required:"false"`
	// Will use default pool if left blank
	VcenterResourcePool string `mapstructure:"respool_name" required:"false"`

	// Accessible datastore path for downloaded files
	OutputDir   string `mapstructure:"output_dir" required:"false"`   // required if using downloading
	DownloadUri string `mapstructure:"download_uri" required:"false"` // required if using downloading
	DsImagePath string `mapstructure:"ds_image_path" required:"false"`
	// Convert and import to vCenter without first downloading image (i.e. image was already downloaded previously)
	// Defaults to false
	ImportNoDownload bool   `mapstructure:"import_no_download" required:"false"`
	SourceImagePath  string `mapstructure:"source_path" required:"false"` // required if bool is true
}

type Datasource struct {
	config Config
}

// --> If making changes to this section, make sure the hcl2spec gets updated as well!
type DatasourceOutput struct{}

func (d *Datasource) ConfigSpec() hcldec.ObjectSpec {
	return d.config.FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Configure(raws ...interface{}) error {
	err := config.Decode(&d.config, nil, raws...)
	if err != nil {
		return err
	}

	if d.config.ArtifactoryToken == "" {
		token := os.Getenv("ARTIFACTORY_TOKEN")
		if token == "" {
			return errors.New("please provide an Artifactory identity token")
		}
	}

	if d.config.ArtifactoryServer == "" {
		server := os.Getenv("ARTIFACTORY_SERVER")
		if server == "" {
			return errors.New("please provide the Artifactory server API URL (ex: https://server.com:8081/artifactory/api)")
		}
	}

	if d.config.VcenterServer == "" {
		vcServer := os.Getenv("VCENTER_SERVER")
		if vcServer == "" {
			return errors.New("missing target vCenter server")
		}
	}

	if d.config.VcenterUser == "" {
		vcUser := os.Getenv("VCENTER_USER")
		if vcUser == "" {
			return errors.New("missing vCenter username")
		}
	}

	if d.config.VcenterPassword == "" {
		vcPass := os.Getenv("VCENTER_PASSWORD")
		if vcPass == "" {
			return errors.New("missing vCenter password")
		}
	}

	if d.config.VcenterDatacenter == "" {
		dcName := os.Getenv("VCENTER_DATACENTER")
		if dcName == "" {
			return errors.New("missing target vCenter datacenter")
		}
	}

	if d.config.VcenterDatastore == "" {
		dsName := os.Getenv("VCENTER_DATASTORE")
		if dsName == "" {
			return errors.New("missing target vCenter datastore")
		}
	}

	if d.config.VcenterCluster == "" {
		clusterName := os.Getenv("VCENTER_CLUSTER")
		if clusterName == "" {
			return errors.New("missing target vCenter cluster")
		}
	}

	if d.config.VcenterFolder == "" {
		folderName := os.Getenv("VCENTER_FOLDER")
		if folderName == "" {
			log.Println("---> No target vCenter folder was provided.")
			log.Println("**** The default root folder will be used instead.")
			log.Println("**** If this is not desired, then please provide a specific vCenter folder name.")
		}
	}

	if d.config.VcenterResourcePool == "" {
		resPoolName := os.Getenv("VCENTER_RESOURCE_POOL")
		if resPoolName == "" {
			log.Println("---> No target vCenter resource pool was provided.")
			log.Println("**** The default resource pool will be used instead.")
			log.Println("**** If this is not desired, then please provide a specific vCenter resource pool name.")
		}
	}

	if d.config.OutputDir == "" && d.config.ImportNoDownload == false {
		outputDir := os.Getenv("OUTPUTDIR")
		if outputDir == "" {
			log.Println("---> No output directory was provided.")
			return errors.New("please provide an accessible datastore output directory")
		}
	}

	if d.config.DownloadUri == "" && d.config.ImportNoDownload == false {
		log.Println("---> No download URI for the artifact was provided. This is required if the artifact should be downloaded before importing into vCenter.")
		return errors.New("provide download_uri or set import_no_download to TRUE and provide full path to the source_path where the image file is located")
	}

	if d.config.ImportNoDownload == true {
		if d.config.SourceImagePath == "" {
			log.Println("---> The 'import_no_download' flag is set to TRUE.")
			return errors.New("source_path is required when import_no_download is true; must be full path to the image file (e.g., '/lab/win22/win22.ova'); if using Windows path, ensure it's properly escaped with double-backslashes")
		}
	}
	return nil
}

func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Execute() (cty.Value, error) {
	var downloadUri, sourcePath, importResult, token, serverApi, outputDir string
	var vcServer, vcUser, vcPass, dcName, dsName, clusterName, resPoolName, folderName, dsImagePath string

	// Artifactory related
	if d.config.ArtifactoryToken == "" {
		token = os.Getenv("ARTIFACTORY_TOKEN")
	} else {
		token = d.config.ArtifactoryToken
	}

	if d.config.ArtifactoryServer == "" {
		serverApi = os.Getenv("ARTIFACTORY_SERVER")
	} else {
		serverApi = d.config.ArtifactoryServer
	}

	// vCenter Related
	if d.config.VcenterServer == "" {
		vcServer = os.Getenv("VCENTER_SERVER")
	} else {
		vcServer = d.config.VcenterServer
	}

	if d.config.VcenterUser == "" {
		vcUser = os.Getenv("VCENTER_USER")
	} else {
		vcUser = d.config.VcenterUser
	}

	if d.config.VcenterPassword == "" {
		vcPass = os.Getenv("VCENTER_PASSWORD")
	} else {
		vcPass = d.config.VcenterPassword
	}

	if d.config.VcenterDatacenter == "" {
		dcName = os.Getenv("VCENTER_DATACENTER")
	} else {
		dcName = d.config.VcenterDatacenter
	}

	if d.config.VcenterDatastore == "" {
		dsName = os.Getenv("VCENTER_DATASTORE")
	} else {
		dsName = d.config.VcenterDatastore
	}

	if d.config.VcenterCluster == "" {
		clusterName = os.Getenv("VCENTER_CLUSTER")
	} else {
		clusterName = d.config.VcenterCluster
	}

	if d.config.VcenterResourcePool == "" {
		resPoolName = os.Getenv("VCENTER_RESOURCE_POOL")
	} else {
		resPoolName = d.config.VcenterResourcePool
	}

	if d.config.VcenterFolder == "" {
		folderName = os.Getenv("VCENTER_FOLDER")
	} else {
		folderName = d.config.VcenterFolder
	}

	// Other

	if d.config.OutputDir == "" {
		outputDir = os.Getenv("OUTPUTDIR")
	} else {
		outputDir = d.config.OutputDir
	}

	if d.config.DownloadUri != "" {
		downloadUri = d.config.DownloadUri
	}

	if d.config.DsImagePath != "" {
		dsImagePath = d.config.DsImagePath
	} else {
		dsImagePath = ""
	}

	var importNoDownload = false // default; whether we will convert and import the image into vCenter without downloading first (i.e. image already downloaded previously)
	if d.config.ImportNoDownload == true {
		importNoDownload = true
	}

	if d.config.SourceImagePath != "" {
		sourcePath = d.config.SourceImagePath
	}

	//------------------------------------------------------------------------------------------------------
	// Gets required info to facilitate the vCenter import process

	folderId, resPoolId, err := vsTasks.GetResourceIds(vcUser, vcPass, vcServer, dcName, folderName, resPoolName, clusterName)
	if err != nil {
		return cty.NilVal, fmt.Errorf("get vCenter folder and resource pool IDs: %w", err)
	}

	// If we are downloading first, parse for needed details, then proceed with download, conversion, import, and templating
	if importNoDownload == false && outputDir != "" && downloadUri != "" {
		imageFileName := artifCommon.ParseUriForFilename(downloadUri)
		imageName := artifCommon.ParseFilenameForImageName(imageFileName)

		// Download Artifacts
		downloadResult := artifTasks.DownloadArtifacts(serverApi, token, downloadUri, outputDir)
		log.Println("Download Result: " + downloadResult)

		if downloadResult == "End of download process" {
			log.Println("Image download completed successfully.")
			log.Println("Checking image type and converting if necessary. This may time some time...")

			importResult = vsTasks.ConvertImportFromDownload(vcUser, vcPass, vcServer, outputDir, downloadUri, dcName, dsName, dsImagePath, imageName, folderId, resPoolId)
		} else {
			return cty.NilVal, fmt.Errorf("download image artifacts: %s", downloadResult)
		}
	} else { // no download flag is true
		log.Println("Checking image type and converting if necessary. This may time some time...")
		importResult = vsTasks.ConvertImportNoDownload(vcUser, vcPass, vcServer, dcName, dsName, sourcePath, dsImagePath, folderId, resPoolId)
	}
	log.Println("Import Result: " + importResult)

	if importResult == "Success" {
		log.Println("Process completed successfully.")
	} else {
		return cty.NilVal, fmt.Errorf("vCenter import did not complete successfully: %s", importResult)
	}

	output := DatasourceOutput{}

	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}
