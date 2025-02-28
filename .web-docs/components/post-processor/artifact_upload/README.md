# JFrog Artifactory Post-Processor

Type:  `artifactory-upload`

The Artifactory post-provisioner `artifactory-upload` is used to upload a newly created artifact image (of type OVA, OVF, or VMTX), and it's associated image files, into JFrog Artifactory.


## Advisements
* When uploading, if the files already exist in the target location, they will be overwritten. 

* When uploading, the files are placed into a directory named after the image withing Artifactory. 
Ex: If the target path is /win-local-libs/, the image file 'win2022.ova' will be placed in /win-local-libs/win2022/.


## Housekeeping
* Artifactory property key/values, artifact URIs, download URIs, Artifactory paths (/repo/folder/...), and file names are **CASE SENSITIVE**. There are a few exceptions, however, it's best to assume case sensitivity for successful outcomes. This is a behavior of the Artifactory API and not something we can control.


## Configuration Reference

- `artifactory_server` (string) - Required; The API address of the Artifactory server (ex: https://server.domain.com:8081/artifactory/api). The URL will differ slightly between cloud-hosted and self-hosted instanced.
    * Environment variable: `ARTIFACTORY_SERVER`
- `artifactory_token` (string) - Required; The Artifactory account Identity Token used to authenticate with the Artifactory server and perform operations. Results are limited to whatever the account has access to. If the account can only "see" a single repository, then the results will only include content from that single repository.
    * Environment variable: `ARTIFACTORY_TOKEN`
- `source_path` (string) - Required; The directory path where the source image file(s) are located (for ex. "C:\\lab" or "/lab")
- `target_path` (string) - *Optional; The target path (/repo/folder/) within Artifactory where the artifact should be uploaded to. If NOT populated, you MUST use `existing_uri_target` instead. The image files will automatically be placed within a subfolder in this path named after the image. For example: /repo/folder --> /repo/folder/image1111/image1111.ova
- `image_type` (string) - Required; The type of image that will be uploaded; supported types are 'ova', 'ovf', and 'vmtx'.
- `image_name` (string) - Required; The base image name
- `existing_uri_target` (string) - *Optional; The URI address of an existing artifact. The plugin will parse this address to determine the /repo/folder/path and set this as the `target_path` for the new artifact.


## Output Data

None


## Basic Example Usage

**Upload Image**
```hcl
	post-processor "artifactory-upload" {
		artifactory_token     = var.artif_token  
        artifactory_server    = var.artif_server
			
		source_path = "c:\\lab"
		target_path = "/test-packer-plugin/win"
		image_type  = "ova"
		image_name  = "test-artifact"
	}
```

**Upload Image Using Existing Artifact Path**
```hcl
	post-processor "artifactory-upload" {
		artifactory_token     = var.artif_token  
        artifactory_server    = var.artif_server 
			
		source_path         = "c:\\lab"
		existing_uri_target = "https://server.domain.com/artifactory/api/storage/test-packer-plugin/existing-artifact.txt"
		image_type  = "ova"
		image_name  = "test-artifact"
	}
```

## FAQ
* Can I use this to upload image files other than OVA, OVF, VMTX?
  - No, this component validates only the expected files that exist with each of these image package types and then uploads them to Artifactory.

* How do I separate these image files from other image files in Artifactory?
  - When files are uploaded, they are placed into their own image name-based folder. For example:  win2022-25-01-23.ova would be uploaded to `/repo/folder/win2022-25-01-23/win2022-25-01-23.ova`.

* What if I have existing folder/files with the same name in my Artifactory?
  - The uploaded files will overwrite the existing files. If this is not desired, rename your image files to include a unique suffix, such as a version or date value. For example, instead of `win2022`, use something like `win2022-1.0.1` or `win2022-250125`, etc.

* I get an error when trying to use a Windows path like this: 'E:\lab-servs\' as my source path. This is the right path. Why is this happening?
  - Windows-based paths must be properly escaped. Instead, use a path like this: 'E:\\lab-servs\\'. This is not an issue with Linux-based paths.

* Is the target path case sensitive?
  - Yes, Artifactory is very particular about paths, artifact names, and property key/values. If the target path is incorrectly cased, Artifactory will think this is a different path and throw an error.

* Is the image name case sensitive?
  - Yes and no. The component will be able to validate the existance of the files in the source path regardless of case, however, the image name provided will also be used to create image name-based directory within Artifactory. Artifactory IS case sensitive about pathing, therefore it will view this as a separate directory, even if one already exists with the same name but different case, and it will create a new directory with the case as it was input.

* Is the image type case sensitive.
  - No.
