# JFrog Artifactory Data Source

Type:  `artifactory`

The Artifactory data source is used to filter and identify an artifact image stored in JFrog Artifactory, and then output the artifact's name, URI, created date, and download URI. 

The use of property key(s)/value(s) as filter parameters can further assist in identifying the correct image. If more than one artifact matches the input parameters, the latest artifact will be returned.


## Configuration Reference

- `artifactory_server` (string) - Required; The API address of the Artifactory server (ex: https://server.domain.com:8081/artifactory/api). The URL will differ slightly between cloud-hosted and self-hosted instanced.
    * Environment variable: `ARTIFACTORY_SERVER`
- `artifactory_token` (string) - Required; The Artifactory account Identity Token used to authenticate with the Artifactory server and perform operations. Results are limited to whatever the account has access to. If the account can only "see" a single repository, then the results will only include content from that single repository.
    * Environment variable: `ARTIFACTORY_TOKEN`

- `artifact_name` (string) - Required; The full or partial name of the artifact/image to search for (ex: win-22).
- `file_type` (string) - Required; The file extension of the desired artifact (ex: vmtx). If left blank, this will default to 'vmtx'.
- `filter` (map[string]string) - Optional; The key/value pairs of artifact properties to filter the artifact by.
- `channel` (string) - Optional; Similar concept to HCP Packer; the channel name assigned to a given artifact. This is simply a property VALUE to the key 'channel'. To be valid, an artifact must have a property named 'channel' assigned with the desired value (ex: 'windows-iis-prod').


## Output Data

- `artifactName` (string) - The name of the artifact.
- `createdDate` (string) - The date the artifact was created.
- `artifactUri` (string) - The URI of the artifact.
- `downloadUri` (string) - The download URI of the artifact.


## Basic Example Usage

**Search for Image by Name Only**
```hcl
data "artifactory" "basic-example" {
    artifactory_token     = "artifactory_token"
    artifactory_server    = "https://server.domain.com:8081/artifactory/api"

    artifact_name = "test-artifact"
    file_type     = "txt"
}
```

**Search for Image Name with Property Filters**
```hcl
data "artifactory" "basic-example" {
    artifactory_token     = "artifactory_token"
    artifactory_server    = "https://server.domain.com:8081/artifactory/api"

    artifact_name = "test-artifact"
    file_type     = "txt"
    
    filter = {
        release = "latest-stable"
        testing = "passed"
    }
}
```

**Search for Image by Channel Property**
```hcl
data "artifactory" "basic-example" {
    artifactory_token     = "artifactory_token"
    artifactory_server    = "https://server.domain.com:8081/artifactory/api"

    artifact_name = "test-artifact"
    file_type     = "txt"
    channel       = "windows-iis-lab"
}
```