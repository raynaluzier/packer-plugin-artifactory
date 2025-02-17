package artifactImport

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
)

const testDatasourceImportHCL2Basic = `
	packer {
		required_plugins {
			artifactory = {
				version = ">= 1.0.25"
				source  = "github.com/raynaluzier/artifactory"
			}
		}
	}

	variable "artif_token" {
		type        = string
		description = "Identity token of the Artifactory account with access to execute commands"
		sensitive   = true
		default     = env("ARTIFACTORY_TOKEN")
	}

	variable "artif_server" {
		type        = string
		description = "The Artifactory API server address"
		default     = env("ARTIFACTORY_SERVER")
	}

	variable "vc_server" {
		type        = string
		description = "The vCenter Server FQDN"
		default     = env("VCENTER_SERVER")
	}

	variable "vc_user" {
		type        = string
		description = "The vCenter username"
		default     = env("VCENTER_USER")
	}

	variable "vc_password" {
		type        = string
		description = "The vCenter username"
		sensitive   = true
		default     = env("VCENTER_PASSWORD")
	}

	variable "vc_datacenter" {
		type        = string
		description = "vCenter datacenter name where the image will be imported to"
		default     = env("VCENTER_DATACENTER")
	}

	variable "vc_datastore" {
		type        = string
		description = "Datastore name where the image will be imported from"
		default     = env("VCENTER_DATASTORE")
	}

	variable "vc_cluster" {
		type        = string
		description = "vCenter cluster name where the image will be imported to; used to find the resource pool"
		default     = env("VCENTER_CLUSTER")
	}

	variable "vc_folder" {
		type        = string
		description = "vCenter folder name where the image will be imported to"
		default     = env("VCENTER_FOLDER")
	}

	variable "vc_respool" {
		type        = string
		description = "vCenter resource pool name where the image will be imported to"
		default     = env("VCENTER_RESOURCE_POOL")
	}

	variable "output_directory" {
		type        = string
		description = "Directory path where the image will be downloaded to"
		default     = env("OUTPUTDIR")
	}

	data "artifactory" "basic-example" {
		# Provide via environment variables
		artifactory_token     = var.artif_token  
		artifactory_server    = var.artif_server

		artifact_name = "ub20pkrt-10031746"
		file_type     = "ovf"
		
		filter = {
			release = "stable"
		}
	}

	data "artifactory-import" "basic-example" {
		artifactory_token   = var.artif_token
		artifactory_server  = var.artif_server

		vcenter_server      = var.vc_server
		vcenter_user        = var.vc_user
		vcenter_password    = var.vc_password
		datacenter_name     = var.vc_datacenter
		datastore_name      = var.vc_datastore
		cluster_name        = var.vc_cluster
		folder_name         = var.vc_folder
		respool_name        = var.vc_respool

		//output_dir          = var.output_directory
		//download_uri        = data.artifactory.basic-example.download_uri

		import_no_download  = true
		source_path         = "E:\\path\\testing\\img2\\img2.ovf"
		target_path         = "E:\\path\\converted\\img2\\img2.vmx"
	}

	locals {
		download_uri = data.artifactory.basic-example.download_uri
	}

	source "null" "basic-example" {
		communicator = "none"
	}

	build {
		sources = [
			"source.null.basic-example"
		]

		provisioner "shell-local" {
			inline = [
				"echo artifact download URI: ${local.download_uri}"
			]
		}
	}
`
// Run with: PACKER_ACC=1 go test -count 1 -v ./... -timeout 180m
func TestAccDatasourceImport_Artifactory(t *testing.T) {
	testCase := &acctest.PluginTestCase{
		Name: "artifactory_datasource_import basic_test",
		Setup: func() error {
			return nil
		},
		Teardown: func() error {
			return nil
		},
		Template: testDatasourceImportHCL2Basic,
		Type: "basic-example",
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}

			logs, err := os.Open(logfile)
			if err != nil {
				return fmt.Errorf("Unable find %s", logfile)
			}
			defer logs.Close()

			logsBytes, err := io.ReadAll(logs)
			if err != nil {
				return fmt.Errorf("Unable to read %s", logfile)
			}
			logsString := string(logsBytes)
			log.Println(logsString)

			outputLog := "The image import and template conversion completed successfully."
			if matched, _ := regexp.MatchString(outputLog+".*", logsString); !matched {
				t.Fatalf("logs doesn't contain expected output %q", logsString)
			}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}