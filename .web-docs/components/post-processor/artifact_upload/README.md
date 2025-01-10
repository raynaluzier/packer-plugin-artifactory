# JFrog Artifactory Post-Processor

Type:  `artifactory-upload`

The Artifactory post-provisioner `artifactory-upload` is used to upload a newly created artifact into JFrog Artifactory and displays the artifact's download URI and the artifact URI in the UI.


## Configuration Reference

- `artifactory_server` (string) - Required; The API address of the Artifactory server (ex: https://server.domain.com:8081/artifactory/api). The URL will differ slightly between cloud-hosted and self-hosted instanced.
- `artifactory_token` (string) - Required; The Artifactory account Identity Token used to authenticate with the Artifactory server and perform operations. Results are limited to whatever the account has access to. If the account can only "see" a single repository, then the results will only include content from that single repository.

- `source_path` (string) - Required; The full directory path with filename where the source artifact is located (for ex. "C:\\lab\\artifact.txt" or "/lab/artifact.txt")
- `target_path` (string) - *Optional; The target path (/repo/folder/path) within Artifactory where the artifact should be uploaded to. If NOT populated, you MUST use `existing_uri_target` instead.
- `file_suffix` (string) - Optional; Distinguishing file name suffix to use such as version, date, ect. where the base image name is always the same.
- `existing_uri_target` (string) - *Optional; The URI address of an existing artifact. The plugin will parse this address to determine the /repo/folder/path and set this as the `target_path` for the new artifact.

## Output Data

None


## Basic Example Usage

**Upload Image with No File Suffix**
```hcl
	post-processor "artifactory-upload" {
		artifactory_token     = var.artif_token  
        artifactory_server    = var.artif_server 
			
		source_path = "c:\\lab\\artifact.txt"
		target_path = "/test-packer-plugin/win"
	}
```

**Upload Image with File Suffix**
```hcl
	post-processor "artifactory-upload" {
		artifactory_token     = var.artif_token  
        artifactory_server    = var.artif_server 
			
		source_path = "c:\\lab\\artifact.txt"
		target_path = "/test-packer-plugin/win"
		file_suffix = "acc-test1"
	}
```

**Upload Image Using Existing Artifact Path**
```hcl
	post-processor "artifactory-upload" {
		artifactory_token     = var.artif_token  
        artifactory_server    = var.artif_server 
			
		source_path         = "c:\\lab\\artifact.txt"
		existing_uri_target = "https://server.domain.com/artifactory/api/storage/test-packer-plugin/existing-artifact.txt"
	}
```