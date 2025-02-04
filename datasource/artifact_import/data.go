//go:generate packer-sdc mapstructure-to-hcl2 -type Config,DatasourceOutput
package artifactImport

import (
	"log"
	"os"
	"strings"

	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/hashicorp/packer-plugin-sdk/hcl2helper"
	"github.com/hashicorp/packer-plugin-sdk/template/config"
	artifCommon "github.com/raynaluzier/artifactory-go-sdk/common"
	artifTasks "github.com/raynaluzier/artifactory-go-sdk/tasks"
	vsphereTasks "github.com/raynaluzier/vsphere-go-sdk/tasks"
	"github.com/zclconf/go-cty/cty"
)

// --> If making changes to this section, make sure the hcl2spec gets updated as well!
type Config struct {
	AritfactoryToken       string `mapstructure:"artifactory_token" required:"true"`
	ArtifactoryServer      string `mapstructure:"artifactory_server" required:"true"`
	
	VcenterServer			string `mapstructure:"vcenter_server" required:"true"`
	VcenterUser				string `mapstructure:"vcenter_user" required:"true"`
	VcenterPassword			string `mapstructure:"vcenter_password" required:"true"`
	// Will use default DC if left blank
	VcenterDatacenter		string `mapstructure:"datacenter_name" required:"false"`
	VcenterDatastore		string `mapstructure:"datastore_name" required:"true"`
	// Used to get Res Pool ID; will use default DC and pool if left blank
	VcenterCluster			string `mapstructure:"cluster_name" required:"false"`
	// Will use default root folder of datacenter if left blank
	VcenterFolder			string `mapstructure:"folder_name" required:"false"`
	// Will use default pool if left blank
	VcenterResourcePool		string `mapstructure:"respool_name" required:"false"`
	// Accessible datastore path for downloaded files
	OutputDir				string `mapstructure:"output_dir" required:"true"`
	
	DownloadUri				string `mapstructure:"download_uri" required:"true"`

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

	if d.config.VcenterServer == "" {
		vcServer := os.Getenv("VCENTER_SERVER")
		if vcServer == "" {
			log.Fatal("Missing the target vCenter Server name.")
		}
	}

	if d.config.VcenterUser == "" {
		vcUser := os.Getenv("VCENTER_USER")
		if vcUser == "" {
			log.Fatal("Missing the vCenter Server username. This is required for authentication.")
		}
	}

	if d.config.VcenterPassword == "" {
		vcPass := os.Getenv("VCENTER_PASSWORD")
		if vcPass == "" {
			log.Fatal("Missing the vCenter Server password. This is required for authentication.")
		}
	}

	if d.config.VcenterDatacenter == "" {
		dcName := os.Getenv("VCENTER_DATACENTER")
		if dcName == "" {
			log.Println("The target vCenter datacenter was not provided. The default datacenter will be used.")
			log.Println("**** If this is not desired, then please provide the datacenter name.")
		}
	}

	if d.config.VcenterDatastore == "" {
		dsName := os.Getenv("VCENTER_DATASTORE")
		if dsName == "" {
			log.Println("Missing the target vCenter datastore name.")
		}
	}

	if d.config.VcenterCluster == "" {
		clusterName := os.Getenv("VCENTER_CLUSTER")
		if clusterName == "" {
			log.Println("No target vCenter was provided, which is used to locate the target resource pool.")
			log.Println("**** The default datacenter and resource pool will be used instead.")
			log.Println("**** If this is not desired, then please provide a specific datacenter, cluster, and resource pool information.")
		}
	}

	if d.config.VcenterFolder == "" {
		folderName := os.Getenv("VCENTER_FOLDER")
		if folderName == "" {
			log.Println("No target vCenter folder was provided.")
			log.Println("**** The default root folder will be used instead.")
			log.Println("**** If this is not desired, then please provide a specific vCenter folder name.")
		}
	}

	if d.config.VcenterResourcePool == "" {
		resPoolName := os.Getenv("VCENTER_RESOURCE_POOL")
		if resPoolName == "" {
			log.Println("No target vCenter resource pool was provided.")
			log.Println("**** The default resource pool will be used instead.")
			log.Println("**** If this is not desired, then please provide a specific vCenter resource pool name.")
		}
	}

	if d.config.OutputDir == "" {
		outputDir := os.Getenv("OUTPUTDIR")
		if outputDir == "" {
			log.Println("No output directory was provided.")
			log.Println("**** Output will be user's home directory.")
			log.Println("**** If importing to vCenter, please provide the directory path to an accessible datastore.")
		}
	}

	if d.config.DownloadUri == "" {
		log.Println("Please provide download URI for image artifact.")
	}

	return nil
}

func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Execute() (cty.Value, error) {
	var downloadUri string

	// Artifactory related
	if d.config.AritfactoryToken == "" {
		d.config.AritfactoryToken = os.Getenv("ARTIFACTORY_TOKEN")
	}
	token := d.config.AritfactoryToken
	
	if d.config.ArtifactoryServer == "" {
		d.config.ArtifactoryServer = os.Getenv("ARTIFACTORY_SERVER")
	}
	serverApi := d.config.ArtifactoryServer

	// vCenter Related
	if d.config.VcenterServer == "" {
		d.config.VcenterServer = os.Getenv("VCENTER_SERVER")
	}
	vcServer := d.config.VcenterServer

	if d.config.VcenterUser == "" {
		d.config.VcenterUser = os.Getenv("VCENTER_USER")
	}
	vcUser := d.config.VcenterUser

	if d.config.VcenterPassword == "" {
		d.config.VcenterPassword = os.Getenv("VCENTER_PASSWORD")
	}
	vcPass := d.config.VcenterPassword

	if d.config.VcenterDatacenter == "" {
		d.config.VcenterDatacenter = os.Getenv("VCENTER_DATACENTER")
	}
	dcName := d.config.VcenterDatacenter

	if d.config.VcenterDatastore == "" {
		d.config.VcenterDatastore = os.Getenv("VCENTER_DATASTORE")
	}
	dsName := d.config.VcenterDatastore

	if d.config.VcenterCluster == "" {
		d.config.VcenterCluster = os.Getenv("VCENTER_CLUSTER")
	}
	clusterName := d.config.VcenterCluster

	if d.config.VcenterResourcePool == "" {
		d.config.VcenterResourcePool = os.Getenv("VCENTER_RESOURCE_POOL")
	}
	resPoolName := d.config.VcenterResourcePool

	if d.config.VcenterFolder == "" {
		d.config.VcenterFolder = os.Getenv("VCENTER_FOLDER")
	}
	folderName := d.config.VcenterFolder

	// Other Info
	if d.config.OutputDir == "" {
		d.config.OutputDir = os.Getenv("OUTPUTDIR")
	}
	outputDir := d.config.OutputDir

	if d.config.DownloadUri != "" {
		downloadUri = d.config.DownloadUri
	}

	imageFileName := artifCommon.ParseUriForFilename(downloadUri)
	imageName     := artifCommon.ParseFilenameForImageName(imageFileName)

	// Download Artifacts
	downloadResult := artifTasks.DownloadArtifacts(serverApi, token, downloadUri, outputDir)
	
	var missingInputsMsg  = "Missing required inputs"
	var downloadFailedMsg = "File download failed"

	// If the download result doesn't contain one of these msgs, proceed with import
	if !strings.Contains(downloadResult, missingInputsMsg) || !strings.Contains(downloadResult, downloadFailedMsg) {
		log.Println("Image download completed successfully. Beginning import into vCenter...")
		importResult := vsphereTasks.ImportVm(vcUser, vcPass, vcServer, dcName, dsName, imageName, folderName, resPoolName, clusterName)

		if strings.Contains(importResult, "Success") {
			log.Println("The image import and template conversion completed successfully.")
		} else {
			log.Fatal("Error: Unable to import and/or convert the image into a VM Template.")
		}
	} else {
		log.Fatal("Error: Failures occurred during image download.")
	}
	
	output := DatasourceOutput{}

	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}