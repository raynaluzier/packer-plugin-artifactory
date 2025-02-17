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
	vsCommon "github.com/raynaluzier/vsphere-go-sdk/common"
	vsGov "github.com/raynaluzier/vsphere-go-sdk/govmomi"
	vsVm "github.com/raynaluzier/vsphere-go-sdk/vm"
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
	OutputDir				string `mapstructure:"output_dir" required:"false"`   // required if using downloading
	DownloadUri				string `mapstructure:"download_uri" required:"false"` // required if using downloading
	// Convert and import to vCenter without first downloading image (i.e. image was already downloaded previously)
	// Defaults to false
	ImportNoDownload		bool   `mapstructure:"import_no_download" required:"false"`
	SourceImagePath			string `mapstructure:"source_path" required:"false"` // required if bool is true
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
			log.Fatal("Missing the target vCenter datastore name.")
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
	var downloadUri, sourcePath string

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

	var importNoDownload = false  // default; whether we will convert and import the image into vCenter without downloading first (i.e. image already downloaded previously)
	if d.config.ImportNoDownload == true {
		importNoDownload = true
	}

	if d.config.SourceImagePath != "" {
		sourcePath = d.config.SourceImagePath
	}

	//------------------------------------------------------------------------------------------------------
	// Gets required info to facilitate the vCenter import process
	vcToken := vsCommon.VcenterAuth(vcUser, vcPass, vcServer)

	folderId, err := vsGov.GetFolderId(vcUser, vcPass, vcServer, folderName, dcName)
	log.Println("Folder ID: " + folderId)
	if err != nil {
		log.Printf("Error getting folder ID: %s\n", err)
	}

	resPoolId, err := vsGov.GetResPoolId(vcUser, vcPass, vcServer, resPoolName, dcName, clusterName)
	log.Println("Resource Pool ID: " + resPoolId)
	if err != nil {
		log.Printf("Error getting resource pool ID: %s\n", err)
	}
	
	// If we are downloading first, parse for needed details, then proceed with download, conversion, import, and templating 
	if importNoDownload == false && outputDir != "" && downloadUri != "" {
		imageFileName := artifCommon.ParseUriForFilename(downloadUri)
		imageName     := artifCommon.ParseFilenameForImageName(imageFileName)

		// Download Artifacts
		log.Println("Downloading artifacts from Artifactory....")
		downloadResult := artifTasks.DownloadArtifacts(serverApi, token, downloadUri, outputDir)
		log.Println("Download Result: " + downloadResult)
		
		var missingInputsMsg  = "Missing required inputs"
		var downloadFailedMsg = "File download failed"

		// If the download result doesn't contain one of these msgs, proceed with import
		if !strings.Contains(downloadResult, missingInputsMsg) || !strings.Contains(downloadResult, downloadFailedMsg) {
			log.Println("Image download completed successfully.")
			log.Println("Checking image type and converting if necessary. This may time some time...")

			fileType, sourcePath, targetPath := vsVm.SetPathsFromDownloadUri(outputDir, downloadUri)
			convertResult := vsVm.ConvertImageByType(fileType, sourcePath, targetPath)
			log.Println("Conversion Result: " + convertResult)

			if convertResult != "Failed" {
				log.Println("Setting vmPathName....")
				vmPathName := vsVm.SetVmPathName(sourcePath, dsName)

				log.Println("Beginning import into vCenter....")
				statusCode := vsVm.RegisterVm(vcToken, vcServer, dcName, vmPathName, imageName, folderId, resPoolId)
				log.Println("Status Code of Register VM task: ", statusCode)
	
				if statusCode == "200" {
					log.Println("Import successful. Marking image as a VM Template...")
					tempResult := vsGov.MarkAsTemplate(vcUser, vcPass, vcServer, imageName, dcName)
					log.Println(tempResult)
	
					if strings.Contains(tempResult, "Success") {
						log.Println("The image import and template conversion completed successfully.")
					} else {
						log.Fatal("Error: Unable to import and/or convert the image into a VM Template.")
					}
				} else {
					log.Fatal("Error registering VMX file with vCenter.")
				}
			} else {
				log.Fatal("Error during image type check and file conversion process.")
			}
		} else {
			log.Fatal("Error: Failures occurred during image download.")
		}
	} else {   // no download flag is true
		// If we're skipping the download and going straight to checking the image, converting if needed, and then importing and templating...
		// 'sourcePath' is the full path to the source image file including filename
		var imageFileName, sourceFolderPath, vmPathName, fileType, convertResult string
		var postConvTargetPath, postConvTargetFilePath, imageName string
		log.Println("'import_no_download' flag is set to TRUE. Skipping artifact download....")

		if sourcePath != "" {
			targetPath := vsVm.SetPathNoDownload(sourcePath)					                // Ex: ova/ovf = E:\Lab, vmtx = E:\Lab\win22\win22.vmx
			isWinPath := vsCommon.CheckPathType(sourcePath)
			if isWinPath == true {
				imageFileName, sourceFolderPath = vsCommon.FileNamePathFromWin(sourcePath)		// Ex: E:\Lab\win22\win22.ova, returns: win22.ova, E:\Lab\win22\
				imageName = vsCommon.ParseFilenameForImageName(imageFileName)		            // Ex: rhel9.ova, returns rhel9
				if fileType != "vmtx" {		// vmtx files have a target path that includes full path to VMX file, the other types just have a folder target
					postConvTargetPath = targetPath
				} else {
					postConvTargetPath = sourceFolderPath
					// since we're grabbing the sourceFolderPath regardless of type, we can use this VMTX postConvert value as it will be the same
				}
				} else {
				imageFileName, sourceFolderPath = vsCommon.FileNamePathFromLnx(sourcePath)		// Ex: /lab/rhel9/rhel9.ova, returns: rhel9.ova, /lab/rhel9/
				imageName = vsCommon.ParseFilenameForImageName(imageFileName)		            // Ex: rhel9.ova, returns rhel9
				if fileType != "vmtx" {
					postConvTargetPath = targetPath + imageName
					postConvTargetPath = vsCommon.CheckAddSlashToPath(postConvTargetPath)
				} else {
					postConvTargetPath = sourceFolderPath
				}
			}
			
			fileType = vsCommon.GetFileType(imageFileName)						               // rhel9.ova, returns ova

			log.Println("Image Filename: " + imageFileName)
			log.Println("Image Name: " + imageName)
			log.Println("File Type: " + fileType)
			log.Println("Source Path: " + sourcePath)
			log.Println("Target Path: " + targetPath)
			log.Println("Source Folder Path: " + sourceFolderPath)
			log.Println("Post Conversion Target Path: " + postConvTargetPath)

			// If this is an OVF image, we need to first move the image files into a sub dir called "ovf_files" and update the conversion source path to here
			// If not, we'll get a file conflict with the disk file(s)
			if fileType == "ovf" {
				log.Println("OVF file detected...")
				log.Println("Moving OVF files into subdirectory of source path called 'ovf_files'...")
				destDir := sourceFolderPath + "ovf_files"		                // ex: 'E:\\path\\to\\win2022\\ovf_files'
				destDir = vsCommon.CheckAddSlashToPath(destDir)                 // add ending slash by os type; 'E:\\path\\to\\win2022\\ovf_files\\'
				moveList, err := vsVm.SetOvfFileList(sourcePath)                // Get list of OVF files to move
				err = vsCommon.MoveFiles(moveList, sourceFolderPath, destDir)   // [file list], 'E:\\path\\to\\win2022\\', 'E:\\path\\to\\win2022\\ovf_files\\'
				if err != nil {
					log.Printf("Error moving files: %v\n", err)
				} else {
					log.Println("Files moved successfully!")
				}

				// Setting the new conversion source path to the ovf_file dir for the conversion process only
				newSourcePath := destDir + imageFileName				        // ex: E:\\path\\to\\win2022\\ovf_files\\win2022.ovf  
				log.Println("Checking image type and converting if necessary. This may time some time...")
				convertResult = vsVm.ConvertImageByType(fileType, newSourcePath, targetPath)
				
			} else {  // ova and vmtx don't need to be moved first
				log.Println("Checking image type and converting if necessary. This may time some time...")
				convertResult = vsVm.ConvertImageByType(fileType, sourcePath, targetPath)
			}

			log.Println("Conversion Result: " + convertResult)

			if convertResult != "Failed" {
				log.Println("Setting vmPathName....")
				postConvTargetFilePath = postConvTargetPath + imageName + ".vmx"

				vmPathName = vsVm.SetVmPathName(postConvTargetFilePath, dsName)
				log.Println("vmPathName: " + vmPathName)

				log.Println("Beginning import into vCenter....")
				statusCode := vsVm.RegisterVm(vcToken, vcServer, dcName, vmPathName, imageName, folderId, resPoolId)
				log.Println("Status Code of Register VM task: ", statusCode)
	
				if statusCode == "200" {
					log.Println("Import successful. Marking image as a VM Template...")
					tempResult := vsGov.MarkAsTemplate(vcUser, vcPass, vcServer, imageName, dcName)
					log.Println(tempResult)
	
					if strings.Contains(tempResult, "Success") {
						log.Println("The image import and template conversion completed successfully.")
					} else {
						log.Fatal("Error: Unable to import and/or convert the image into a VM Template.")
					}
				} else {
					log.Fatal("Error registering VMX file with vCenter.")
				}
			} else {
				log.Fatal("Error during image type check and file conversion process.")
			}
		}
	}
	
	output := DatasourceOutput{}

	return hcl2helper.HCL2ValueFromConfig(output, d.OutputSpec()), nil
}