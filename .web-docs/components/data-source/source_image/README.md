# JFrog Artifactory Data Source

Type:  `artifactory`

The Artifactory data source is used to filter and identify an artifact image stored in JFrog Artifactory, and then output the artifact's name, URI, created date, and download URI. 

The use of property key(s)/value(s) as filter parameters can further assist in identifying the correct image. If more than one artifact matches the input parameters, the latest artifact will be returned.


## Housekeeping
* Artifactory property key/values, artifact URIs, download URIs, Artifactory paths (/repo/folder/...), and file names are **CASE SENSITIVE**. There are a few exceptions, however, it's best to assume case sensitivity for successful outcomes. This is a behavior of the Artifactory API and not something we can control.


## Configuration Reference

- `artifactory_server` (string) - Required; The API address of the Artifactory server (ex: https://server.domain.com:8081/artifactory/api). The URL will differ slightly between cloud-hosted and self-hosted instanced.
    * Environment variable: `ARTIFACTORY_SERVER`
- `artifactory_token` (string) - Required; The Artifactory account Identity Token used to authenticate with the Artifactory server and perform operations. Results are limited to whatever the account has access to. If the account can only "see" a single repository, then the results will only include content from that single repository.
    * Environment variable: `ARTIFACTORY_TOKEN`

- `artifact_name` (string) - Required; The full or partial name of the artifact/image to search for (ex: win-22).
- `file_type` (string) - Required; The file extension of the desired artifact (ex: vmtx). If left blank, this will default to 'vmtx'.
- `filter` (map[string]string) - Optional; The key/value pairs of artifact properties to filter the artifact by.
- `channel` (string) - Optional; Similar concept to HCP Packer; the channel name assigned to a given artifact. This is simply a property VALUE to the key 'channel'. To be valid, an artifact must have a property named 'channel' assigned with the desired value that designates an environment/tier/system type/etc that it is meant for (ex: 'windows-iis-prod').


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

```hcl
data "artifactory" "basic-example" {
    artifactory_token     = "artifactory_token"
    artifactory_server    = "https://server.domain.com:8081/artifactory/api"

    artifact_name = "test-artifact"
    file_type     = "txt"
    
    filter = {
        mytag = "some_value"
        color = "blue"
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

## FAQ
* I'm not sure what to use for the 'channel' option? Where do I find that?
  - This is meant to mimic the Channel option found in HCP Packer. In this case, it's nothing more than a property key assigned to your artifact within Artifactory with a corresponding value that should match the type of environment/build that it's intended for. 
  
  To be used as intended, it should only be assigned to a single image artifact.

  If you do not have a 'channel' property assigned to an artifact, then it won't be of use. 

* Should the file type be in the format of '**.**ova' or 'ova'?
  - The file type supports either format.

* Does this component only support OVA, OVF, or VMTX file types?
  - No, while the 'artifactory-import' process DOES (because it's meant for a specific purpose), this component can locate whatever file type you have in your Artifactory instance.

* Are the property keys and values case sensitive?
  - Yes. Artifactory is very particular about casing for paths, artifacts, and properties and views different casing as a different item. If the case is not correct for either the KEY or VALUE, Artifactory will not be able to find the property.

* Is the 'channel' option case sensitive?
  - Yes. This is technically a property key/value and treated exactly the same as any other Artifactory property.

* Is the artifact name case sensitive?
  - No. While Artifactory is particular about case typically, in this case, the search is case insensitive.

* What if I have multiple artifacts with the same name?
  - The component will search for all artifacts that contain the artifact name provided. It will then filter those artifacts by file type. Next it will filter based on matching all of the property key/values, if provided. If the results return more than one option, the artifact with the most recent creation date is returned.

  While the search is pretty accurate, if you give extremely vague parameters, it's possible you won't get the result you expect. If this is the case, try providing a bit more detail/more complete information in the parameters.

* Can I provide a partial artifact name?
  - Yes. Please see note above about how searches are conducted. If you aren't getting the result you expect, try providing a bit more detail/more complete information in the parameters.