# JFrog Artifactory Post-Processor

Type:  `artifactory-upload-other`

The Artifactory post-processor is used to upload any supporting artifacts that go with a given image file (type OVA, OVF, or VMTX). This could be scripts, metadata, documents, etc.

This post-processor differs from the `artifactory-upload` post-processor in that the `artifactory-upload` process is used to upload OVA, OVF, or VMTX files which are comprised of one or more specific files, depending on the image type.


## Housekeeping
* Artifactory property key/values, artifact URIs, download URIs, Artifactory paths (/repo/folder/...), and file names are **CASE SENSITIVE**. There are a few exceptions, however, it's best to assume case sensitivity for successful outcomes. This is a behavior of the Artifactory API and not something we can control.


## Configuration Reference

- `artifactory_server` (string) - Required; The API address of the Artifactory server (ex: https://server.domain.com:8081/artifactory/api). The URL will differ slightly between cloud-hosted and self-hosted instanced.
    * Environment variable: `ARTIFACTORY_SERVER`
- `artifactory_token` (string) - Required; The Artifactory account Identity Token used to authenticate with the Artifactory server and perform operations. Results are limited to whatever the account has access to. If the account can only "see" a single repository, then the results will only include content from that single repository.
    * Environment variable: `ARTIFACTORY_TOKEN`	
- `source_path` (string) - Required; The directory path where the files to be uploaded reside.
- `artifactory_path` (string) - Required; The repo path within Artifactory where the artifact(s) are to uploaded to (ex: /repo/folder).
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

## FAQ
* What if I want to store these additional files with the image?
  - Include the image name-based folder in the artifactory path. So if the image was uploaded using the `artifactory-upload` post-processor, we know that it will be stored in it's own image name-based folder (Ex: /myrepo/win/**win2022**/win2022.ova). In this case, simply specify `/myrepo/win/win2022/` as the artifactory path for these other files.

* What if I want to store some files in one place and some files in another?
  - Use multiple `artifactory-upload` blocks.

* Is the Artifactory path case sensitive?
  - Yes. Artifactory is very particular about casing with regards to paths, artifacts, and properties. If the case does not match, Artifactory will think this is a different artifact and throw an error that it can't find it.

* Are the files in the file list case sensitive?
  - Yes and no. If you specify the file name in lowercase but the file exists in uppercase in the source directory, the file will still be successfully verified and the file will be uploaded in its original case.
  - However, if the file already exists in the same path in Artifactory, but in a different case, Artifactory will view it as a different file and upload a separate copy in the differing case.

* Is the source path case sensitive?
  - No.