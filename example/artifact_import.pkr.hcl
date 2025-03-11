# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

packer {
    required_plugins {
        artifactory = {
            version = ">= 1.1.0"
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

variable "download_addr" {
	type        = string
	description = "Download URI of the image artifact to be downloaded"
}

data "artifactory" "basic-example" {
	# Provide via environment variables
	artifactory_token     = var.artif_token  
	artifactory_server    = var.artif_server

	artifact_name = "test-artifact"
	file_type     = "ova"
		
	filter = {
		release = "latest-stable"
	}
}

data "artifactory-import" "basic-example" {
	artifactory_token   = "artifactory_token"
	artifactory_server  = "https://server.domain.com:8081/artifactory/api"

	vcenter_server      = var.vc_server
	vcenter_user        = var.vc_user
	vcenter_password    = var.vc_password
	datacenter_name     = var.vc_datacenter
	datastore_name      = var.vc_datastore
	cluster_name        = var.vc_cluster
	folder_name         = var.vc_folder
	respool_name        = var.vc_respool

	ouput_dir           = var.output_directory							// Ex: "/mnt/share/lab-servs/"  // remove if 'import_no_download = true'
	download_uri        = data.artifactory.basic-example.download_uri     								// remove if 'import_no_download = true'
	//ds_image_path     = "/lab-servs/"

	//import_no_download = true
	//source_path        = "/mnt/share/lab-servs/img22/img22.ova"
	//ds_image_path      = "/lab-servs/img22.ova"
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