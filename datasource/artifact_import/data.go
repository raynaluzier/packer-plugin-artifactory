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
	vsTasks "github.com/raynaluzier/vsphere-go-sdk/tasks"
	"github.com/zclconf/go-cty/cty"
)

// --> If making changes to this section, make sure the hcl2spec gets updated as well!
type Config struct {
	AritfactoryToken       string `mapstructure:"artifactory_token" required:"true"`
	ArtifactoryServer      string `mapstructure:"artifactory_server" required:"true"`
	
	VcenterServer			string `mapstructure:"vcenter_server" required:"true"`
	VcenterUser				string `mapstructure:"vcenter_user" required:"true"`
	VcenterPassword			string `mapstructure:"vcenter_password" required:"true"`
	VcenterDatacenter		string `mapstructure:"datacenter_name" required:"true"`
	VcenterDatastore		string `mapstructure:"datastore_name" required:"true"`
	VcenterCluster			string `mapstructure:"cluster_name" required:"true"`
	// Will use default root folder of datacenter if left blank
	VcenterFolder			string `mapstructure:"folder_name" required:"false"`
	// Will use default pool if left blank
	VcenterResourcePool		string `mapstructure:"respool_name" required:"false"`

	// Accessible datastore path for downloaded files
	OutputDir				string `mapstructure:"output_dir" required:"false"`   // required if using downloading
	DownloadUri				string `mapstructure:"download_uri" required:"false"` // required if using downloading
	// Convert and import to vCenter without first downloading image (i.e. image was already downloaded previously)
	// Defaults to false
	ImportNoDownload		bool   `mapstructure:"import_no_download" required:"false"`
	SourceImagePath			string `mapstructure:"source_path" required:"false"` // required if bool is true
	// Defaults to INFO
	Logging                string `mapstructure:"logging" required:"false"`
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
			log.Fatal("Missing the target vCenter datacenter name.")
		}
	}

	if d.config.VcenterDatastore == "" {
		dsName := os.Getenv("VCENTER_DATASTORE")
		if dsName == "" {
			log.Fatal("Missing the target vCenter datastore name.")
		}
	}

	if d.config.VcenterCluster == "" {
		clusterName := os.Getenv("VCENTER_CLUSTER")
		if clusterName == "" {
			log.Fatal("Missing the target vCenter cluster.")
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

	if d.config.OutputDir == "" && d.config.ImportNoDownload == false {
		outputDir := os.Getenv("OUTPUTDIR")
		if outputDir == "" {
			log.Println("No output directory was provided.")
			log.Fatal("Please provide the directory path to an accessible datastore.")
		}
	}

	if d.config.DownloadUri == "" && d.config.ImportNoDownload == false {
		log.Println("No download URI for the artifact was provided. This is required if the artifact should be downloaded before importing into vCenter.")
		log.Fatal("If the image does not need to be downloaded first, please set 'import_no_download' to TRUE, and provide full file path to the source directory where the image file (OVA, OVF, or VMTX) is located.")
	}

	if d.config.ImportNoDownload == true {
		if d.config.SourceImagePath == "" {
			log.Println("The 'import_no_download' flag is set to TRUE.")
			log.Fatal("The 'source_path' to the full path for the image file (OVA, OVF, or VMTX) is required. Ex: '/lab/win22/win22.ova' If using a Windows path, ensure it is properly escaped with double-backslashes.")
		}
	}
	return nil
}

func (d *Datasource) OutputSpec() hcldec.ObjectSpec {
	return (&DatasourceOutput{}).FlatMapstructure().HCL2Spec()
}

func (d *Datasource) Execute() (cty.Value, error) {
	var downloadUri, sourcePath, importResult string

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
	if d.config.Logging == "" {
		d.config.Logging = os.Getenv("LOGGING")
	}
	logLevel := d.config.Logging

	if d.config.OutputDir == "" {
		d.config.OutputDir = os.Getenv("OUTPUTDIR")
	}
	outputDir := d.config.OutputDir

	if d.config.DownloadUri != "" {
		downloadUri = d.config.DownloadUri
	}

	var importNoDownload = false  // default; whether we will convert and import the image into vCenter without downloading first (i.e. image already downloaded previously)
	if d.config.ImportNoDownload == true {
		importNoDownload = true
	}

	if d.config.SourceImagePath != "" {
		sourcePath = d.config.SourceImagePath
	}

	//------------------------------------------------------------------------------------------------------
	// Gets required info to facilitate the vCenter import process

	folderId, resPoolId, err := vsTasks.GetResourceIds(vcUser, vcPass, vcServer, logLevel, dcName, folderName, resPoolName, clusterName)
	if err != nil {
		log.Fatal("Error getting folder and resource pool IDs")
	}
	
	// If we are downloading first, parse for needed details, then proceed with download, conversion, import, and templating 
	if importNoDownload == false && outputDir != "" && downloadUri != "" {
		imageFileName := artifCommon.ParseUriForFilename(downloadUri)
		imageName     := artifCommon.ParseFilenameForImageName(imageFileName)

		// Download Artifacts
		downloadResult := artifTasks.DownloadArtifacts(serverApi, token, logLevel, downloadUri, outputDir)
		log.Println("Download Result: " + downloadResult)
		
		var missingInputsMsg  = "Missing required inputs"
		var downloadFailedMsg = "File download failed"

		// If the download result doesn't contain one of these msgs, proceed with import
		if !strings.Contains(downloadResult, missingInputsMsg) || !strings.Contains(downloadResult, downloadFailedMsg) {
			log.Println("Image download completed successfully.")
			log.Println("Checking image type and converting if necessary. This may time some time...")

			importResult = vsTasks.ConvertImportFromDownload(vcUser, vcPass, vcServer, logLevel, outputDir, downloadUri, dcName, dsName, imageName, folderId, resPoolId)
		} else {
			log.Fatal("Error: Failures occurred during image download.")
		}
	} else {   // no download flag is true
		log.Println("Checking image type and converting if necessary. This may time some time...")
		importResult = vsTasks.ConvertImportNoDownload(vcUser, vcPass, vcServer, logLevel, dcName, dsName, sourcePath, folderId, resPoolId)
	}
	log.Println("Import Result: " + importResult)

	if importResult == "Success" {
		log.Println("Process completed successfully.")
	} else if importResult == "Failed" || importResult == "" {
		log.Fatal("Process did not complete successfully.")
	}
	
	output := DatasourceOutput{}

	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}