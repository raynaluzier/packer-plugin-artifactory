# JFrog Artifactory Post-Processor

Type:  `artifactory-upload-other`

The Artifactory post-processor is used to upload any supporting artifacts that go with a given image file (type OVA, OVF, or VMTX). This could be scripts, metadata, documents, etc.

This post-processor differs from the `artifactory-upload` post-processor in that the `artifactory-upload` process is used to upload OVA, OVF, or VMTX files which are comprised of one or more specific files, depending on the image type.

## Configuration Reference

- `artifactory_server` (string) - Required; The API address of the Artifactory server (ex: https://server.domain.com:8081/artifactory/api). The URL will differ slightly between cloud-hosted and self-hosted instanced.
    * Environment variable: `ARTIFACTORY_SERVER`
- `artifactory_token` (string) - Required; The Artifactory account Identity Token used to authenticate with the Artifactory server and perform operations. Results are limited to whatever the account has access to. If the account can only "see" a single repository, then the results will only include content from that single repository.
    * Environment variable: `ARTIFACTORY_TOKEN`	
- `source_path` (string) - Required; The directory path where the files to be uploaded reside.
- `artifactory_path` (string) - Required; The repo path within Artifactory where the artifact(s) are to uploaded to (ex: /repo/folder).
- `folder_name` (string) - Optional; If files are to be uploaded to the same directory as a given image, set this to be the image name; otherwise leave blank.
  ** You may optionally include the folder name in the Artifactory path above. This was only separated out to specifically draw attention to whether or not to place the file(s) with an image, but handling it with this input is not required.
- `file_list` ([]string) - Required; The list of file names with extensions to be uploaded; each file should be in quotes.

## Output Data

None


## Basic Example Usage

**Upload Artifact(s) to a Specific Folder in Artifactory (like a folder for an associated image)
```hcl
	post-processor "artifactory-upload-other" {
        artifactory_token     = var.artif_token  
        artifactory_server    = var.artif_server
			
		source_path      = "c:\\lab\\test-dir"
		file_list        = ["testfile1.txt", "testfile2.txt"]
		artifactory_path = "/rpt-libs-local/"
		folder_name      = "win2022"
	}
```

**Upload Artifact(s) to a Different Folder in Artifactory from the Image**
```hcl
	post-processor "artifactory-upload-other" {
        artifactory_token     = var.artif_token  
        artifactory_server    = var.artif_server 
			
		source_path      = "c:\\lab\\test-dir"
		file_list        = ["scriptA.ps1", "scriptB.ps1"]
		artifactory_path = "/rpt-libs-local/scripts/"
	}
```