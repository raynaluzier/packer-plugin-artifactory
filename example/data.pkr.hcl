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