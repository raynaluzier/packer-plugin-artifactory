# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

data "artifactory" "basic-example" {
    # --> Provide creds via environment variables
    artifactory_token     = "1234567890abcdefghijkl1234567890mnopqrstuv"
    artifactory_server    = "https://myserver.com:8081/artifactory/api"
    
    artifactory_outputdir = "C:\\lab\\output-directory"
    artifactory_logging   = "INFO"

    artifact_name = "test-artifact"
    file_type     = "txt"
    channel       = "windows-iis-prod"
    
    filter = {
        release = "latest-stable"
        testing = "passed"
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