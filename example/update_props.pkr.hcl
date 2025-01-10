
packer {
    required_plugins {
        artifactory = {
            version = ">= 1.0.8"
            source  = "github.com/raynaluzier/artifactory"
        }
    }
}

variable "artif_token" {
    type        = string
    description = "Identity token of the Artifactory account with access to execute commands"
    default     = env("ARTIFACTORY_TOKEN")
}

variable "artif_server" {
    type        = string
    description = "The Artifactory API server address"
    default     = env("ARTIFACTORY_SERVER")
}

data "artifactory" "basic-example" {
    # Provide via environment variables
    artifactory_token     = var.artif_token  
    artifactory_server    = var.artif_server 
    artifactory_logging   = "INFO" 

    artifact_name = "test-artifact"
    file_type     = "txt"
}

locals {
    name         = data.artifactory.basic-example.name
    create_date  = data.artifactory.basic-example.creation_date
    artifact_uri = data.artifactory.basic-example.artifact_uri
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
            "echo artifact name: ${local.name}",
            "echo artifact creation date: ${local.create_date}",
            "echo artifact URI: ${local.artifact_uri}",
            "echo artifact download URI: ${local.download_uri}"
        ]
    }

    post-processor "artifactory-update-props" {
        artifactory_token     = var.artif_token  
        artifactory_server    = var.artif_server 

		artifact_uri = local.artifact_uri
		properties   = {
		    release = "latest-stable"
		}
	}
}