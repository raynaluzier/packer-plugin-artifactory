
packer {
    required_plugins {
        artifactory = {
            version = ">= 0.0.30"
            source  = "github.com/raynaluzier/artifactory"
        }
    }
}


data "artifactory" "basic-example" {
    # Provide via environment variables
    artifactory_token     = "eyJ2ZXIiOiIyIiwidHlwIjoiSldUIiwiYWxnIjoiUlMyNTYiLCJraWQiOiJ2V3V1Tm9uZE0wVVBtRVdtR3cwMEpXTnYzUW9fc2N4WW92WGcwVm15cHdzIn0.eyJzdWIiOiJqZmFjQDAxamY5eTN4cnFhbmFuMWo1anIwZ3Awa2MxL3VzZXJzL3BhY2tlciIsInNjcCI6ImFwcGxpZWQtcGVybWlzc2lvbnMvYWRtaW4iLCJhdWQiOiIqQCoiLCJpc3MiOiJqZmZlQDAxamY5eTN4cnFhbmFuMWo1anIwZ3Awa2MxIiwiaWF0IjoxNzM0NDYyOTUyLCJqdGkiOiIzNzliOThlNy0yNWZiLTQ2ZTgtOGFjYy04N2NhM2Y2MmQ4OTIiLCJ0aWQiOiJhMGZ0emZqNHgwOTRpIn0.jOff5zBZ70cb_8uJZWWaPnUO4Ub0JqFYrY57JHmqRHy-gqsZRxtBjCNpNATFkjzvuEtoN_qUbRaEXC2FSHKgjiA7boCOjCI_Hu1_Nsls2RyRsCzbYlLdHvfGMiQoh2DjFYZ49YWNQ8QfGuQ6TpUWR5FxT0Dl-R2hkemzVVxOG-lRFC_QIEJSSKR53zvhRBELMcYbJqSJ4NIRdJoOsON-NZCqLP4pOcV3U3u1-JFpndgewlq6jZk5z9eppSbtItmin8LYlb56QCtgNG9DH41wrHJUmuwlQXRDPcZF6gr0y1yi4M5QPtj5BB5IPeYKUchC-T-W6xKnUvXhMDCS5SbpRw"
    artifactory_server    = "https://riverpointtechnology.jfrog.io/artifactory/api"
    artifactory_logging   = "INFO"

    artifact_name = "test-artifact"
    file_type     = "txt"
    
    filter = {
        release = "latest-stable"
    }
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
}