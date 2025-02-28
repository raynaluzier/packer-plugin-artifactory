Type: `artifactory-download-other`

# JFrog Artifactory Data Source

The Artifactory data source is used to download any supporting artifacts that go with a given image file (type OVA, OVF, or VMTX). This could be scripts, metadata, documents, other image file types, or image files that shouldn't be converted. **No files files are converted or imported into vCenter as part of this.**

This data source differs from the `artifactory-import` data source in that the `artifactory-import` process is used to download OVA, OVF, or VMTX files which are comprised of one or more specific files, depending on the image type. Depending on the image type, the associated files are determined, validated, and optionally downloaded before being converted to VMX, imported into vCenter, and marked as a VM Template.


## Advisements
* When downloading, if the files already exist in the target location, they will be overwritten. 


## Housekeeping
* Artifactory property key/values, artifact URIs, download URIs, Artifactory paths (/repo/folder/...), and file names are **CASE SENSITIVE**. There are a few exceptions, however, it's best to assume case sensitivity for successful outcomes. This is a behavior of the Artifactory API and not something we can control.


## Configuration Reference

- `artifactory_server` (string) - Required; The API address of the Artifactory server (ex: https://server.domain.com:8081/artifactory/api). The URL will differ slightly between cloud-hosted and self-hosted instanced.
    * Environment variable: `ARTIFACTORY_SERVER`
- `artifactory_token` (string) - Required; The Artifactory account Identity Token used to authenticate with the Artifactory server and perform operations. Results are limited to whatever the account has access to. If the account can only "see" a single repository, then the results will only include content from that single repository.
    * Environment variable: `ARTIFACTORY_TOKEN`
- `output_dir` (string) - Required; The directory where the artifacts should be downloaded to; ensure this is properly escaped as necessary.
    * Environment variable: `OUTPUTDIR`
- `artifactory_path` (string) - Required; The repo path within Artifactory where the artifact(s) to be downloaded reside(s) (ex: /repo/folder).
- `file_list` ([]string) - Required; The list of file names with extensions to be downloaded; each file should be in quotes.


## Output Data

No outputs.


## Basic Example Usage

```hcl
data "artifactory-download-other" "basic-example" {
		artifactory_token     = var.artif_token  
		artifactory_server    = var.artif_server

		output_dir       = "c:\\lab\\output-test\\"
		artifactory_path = "/test-repo/testing/"
		file_list        = ["testfile3.txt", "testfile4.txt"]
}
```

## FAQ
* What if I want to store these files with the image files?
  - Specify the output directory of the image files.
  
* What if I have existing files with the same name in my output directory?
  - The downloaded files will overwrite the existing files.

* Can this do single file downloads?
  - Yes.

* Are the files in the file list case sensitive?
  - Yes. Artifactory is very particular about casing with regards to paths, artifacts, and properties. If the case does not match, Artifactory will think this is a different artifact and throw an error that it can't find it.

* Is the Artifactory path case sensitive?
  - Yes. See previous note.

* Is the output directory case sensitive?
  - No.

* Can I use this component to import my file into vCenter?
  - No. This component does not do any importing. It simple downloads a defined list of one or more files.

* What if I need to download files from different locations?
  - Use a separate instance of this component for each location.
