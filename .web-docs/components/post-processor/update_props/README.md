# JFrog Artifactory Post-Processor

Type:  `artifactory-update-props`

The Artifactory post-provisioner `artifactory-update-props` is used to assign one or more properties to an artifact within Artifactory.


## Configuration Reference

- `artifactory_server` (string) - Required; The API address of the Artifactory server (ex: https://server.domain.com:8081/artifactory/api). The URL will differ slightly between cloud-hosted and self-hosted instanced.
    * Environment variable: `ARTIFACTORY_SERVER`
- `artifactory_token` (string) - Required; The Artifactory account Identity Token used to authenticate with the Artifactory server and perform operations. Results are limited to whatever the account has access to. If the account can only "see" a single repository, then the results will only include content from that single repository.
    * Environment variable: `ARTIFACTORY_TOKEN`
- `logging` (string) - Optional; The logging level to use (INFO, WARN, ERROR, DEBUG). This defaults to 'INFO' if left blank.
    * Environment variable: `LOGGING`
- `artifact_uri` (string) - Required; The URI of the image artifact. The file type should be OVA, OVF, or VMTX. All standard files for the given image type will be included (ex: OVF images also include .MDF and .VMDK files; these will be included automatically).
- `properties` (map[string]string) - Required; The key/value pairs of one or more properties to apply to the artifact.

## Output Data

None


## Basic Example Usage

**Update Artifact with Single Property**
```hcl
	post-processor "artifactory-update-props" {
		artifactory_token     = var.artif_token  
    	artifactory_server    = var.artif_server
		logging               = "DEBUG" 
			
		artifact_uri = "${var.artif_server}/storage/test-packer-plugin/test-artifact.ova"
		properties   = {
			release = "latest-stable"
		}
	}
```

**Update Artifact with Multiple Properties**
```hcl
	post-processor "artifactory-update-props" {
		artifactory_token     = var.artif_token  
    	artifactory_server    = var.artif_server 
			
		artifact_uri = "${var.artif_server}/storage/test-packer-plugin/test-artifact.ova"
		properties   = {
			release = "latest-stable"
			testing = "passed"
		}
	}
```