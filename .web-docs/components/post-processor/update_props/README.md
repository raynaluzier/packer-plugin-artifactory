# JFrog Artifactory Post-Processor

Type:  `artifactory-update-props`

The Artifactory post-provisioner `artifactory-update-props` is used to assign one or more properties to an artifact within Artifactory.


## Advisements
* If the property key already exists on the artifact, the new value will simply be updated.

* If an incorrectly case property key is passed to the artifact, the artifact will treat that property as a completely separate property.


## Housekeeping
* Artifactory property key/values, artifact URIs, download URIs, Artifactory paths (/repo/folder/...), and file names are **CASE SENSITIVE**. There are a few exceptions, however, it's best to assume case sensitivity for successful outcomes. This is a behavior of the Artifactory API and not something we can control.


## Configuration Reference

- `artifactory_server` (string) - Required; The API address of the Artifactory server (ex: https://server.domain.com:8081/artifactory/api). The URL will differ slightly between cloud-hosted and self-hosted instanced.
    * Environment variable: `ARTIFACTORY_SERVER`
- `artifactory_token` (string) - Required; The Artifactory account Identity Token used to authenticate with the Artifactory server and perform operations. Results are limited to whatever the account has access to. If the account can only "see" a single repository, then the results will only include content from that single repository.
    * Environment variable: `ARTIFACTORY_TOKEN`
- `artifact_uri` (string) - Required; The URI of the image artifact. The file type should be OVA, OVF, or VMTX. All standard files for the given image type will be included (ex: OVF images also include .MDF and .VMDK files; these will be included automatically).
- `properties` (map[string]string) - Required; The key/value pairs of one or more properties to apply to the artifact. Even if the property key already exists, the value will simply be updated.
** NOTE: Property key/values are CASE SENSITIVE. Therefore, passing incorrectly cased property keys will create a NEW property in that case. 

For example, if `testartifact.txt` has the key/value property of 'release=latest-stable' and the key/value property 'RELEASE=stable' is passed to this same artifact, the artifact will then have BOTH entries rather than updating the original. It will have:
- release = latest-stable
- RELEASE = stable


## Output Data

None


## Basic Example Usage

**Update Artifact with Single Property**
```hcl
	post-processor "artifactory-update-props" {
		artifactory_token     = var.artif_token  
    	artifactory_server    = var.artif_server
			
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

## FAQ
* Is the artifact URI case sensitive?
  - Yes. Artifactory is very particular about casing with regards to paths, artifacts, and properties. If the case does not match, Artifactory will think this is a different artifact and throw an error that it can't find it.
  - You can use the `artifactory` datasource to locate the artifact within Artifactory, which will then output the artifact_uri as one of it's available outputs. This value can then be used in the `artifactory-upload` component by referencing the schema in dot-notation (for example: `data.artifactory.basic-example.artifact_uri`).

* What happens if I update a property that already exists with the same value?
  - Nothing. The update will say it completed successfully.

* How do I update a property value for a property key that already exists?
  - Simply pass the same property key (with correct case) and the new property value.

* Are property keys/values case sensitive?
  - Yes. Artifactory is very particular about casing with regards to paths, artifacts, and properties. If the property key case does not match, Artifactory will think this is a different property and assign the "new" key/value pair to the artifact. If the property value does not match, Artifactory will update the property value to match what was input.

* Does the property have to exist in Artifactory first before I can assign it?
  - No. Artifactory will take whatever property key/value you provide and assign it to the artifact.

* What if I want to update properties on multiple artifacts?
  - Use separate `artifactory-update` blocks for each artifact. Each block can have one or more property key/values provided.