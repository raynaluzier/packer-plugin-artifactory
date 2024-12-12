
data "artifactory" "basic-example" {
    # Provide via environment variables
    //artifactory_token     = ""
    //artifactory_server    = ""      // https://myserver.com:8081/artifactory/api
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