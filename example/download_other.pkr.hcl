
packer {
    required_plugins {
        artifactory = {
            version = ">= 1.0.27"
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

data "artifactory-download-other" "basic-example" {
		artifactory_token     = var.artif_token  
		artifactory_server    = var.artif_server
        logging               = "INFO"

		output_dir       = "c:\\lab\\output-test\\"
		artifactory_path = "/test-libs-local/testing/"
		file_list        = ["testfile3.txt", "testfile4.txt"]
}

source "null" "basic-example" {
    communicator = "none"
}

build {
    sources = [
        "source.null.basic-example"
    ]
}