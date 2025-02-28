
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

source "null" "basic-example" {
    communicator = "none"
}

build {
    sources = [
        "source.null.basic-example"
    ]

    post-processor "artifactory-upload-other" {
        artifactory_token     = var.artif_token  
        artifactory_server    = var.artif_server 
			
		source_path      = "c:\\lab\\test-dir"
		file_list        = ["testfile1.txt", "testfile2.txt"]
		artifactory_path = "/test-libs-local/"
	}
}