Type:  `artifactory-import`

# JFrog Artifactory Data Source

The Artifactory data source is used to download an artifact image of type OVA, OVF, or VMTX and it's standard associated image files that are stored in JFrog Artifactory to an accessible datastore path, convert the image to a VMX, import into vCenter, and then mark it as a template which is then ready for the vsphere-clone builder to consume. 

This component is meant to be used in conjunction with the [vSphere-Clone](https://developer.hashicorp.com/packer/integrations/hashicorp/vsphere/latest/components/builder/vsphere-clone) Builder.

## Configuration Reference

- `artifactory_server` (string) - Required; The API address of the Artifactory server (ex: https://server.domain.com:8081/artifactory/api). The URL will differ slightly between cloud-hosted and self-hosted instanced.
    * Environment variable: `ARTIFACTORY_SERVER`
- `artifactory_token` (string) - Required; The Artifactory account Identity Token used to authenticate with the Artifactory server and perform operations. Results are limited to whatever the account has access to. If the account can only "see" a single repository, then the results will only include content from that single repository.
    * Environment variable: `ARTIFACTORY_TOKEN`
- `vcenter_server` (string) - Required; The FQDN or IP address of the target vCenter where the image should be imported.
    * Environment variable: `VCENTER_SERVER`
- `vcenter_user` (string) - Required; The vCenter service or user account (ex: jsmith@domain.com) that will be querying for the corresponding resource pool and folder vSphere IDs and importing the VMX into vCenter. This account needs access to the target datacenter, cluster, datastore, resource pool, and folder within vCenter for the tasks to complete successfully.
    * Environment variable: `VCENTER_USER`
- `vcenter_password` (string) - Required; The vCenter password for the associated service or user account provided above.
    * Environment variable: `VCENTER_PASSWORD`
- `datacenter_name` (string) - Optional, but recommended; Target datacenter name the template will be imported into; Will try to use the default datacenter if this is left blank.
    * Environment variable: `VCENTER_DATACENTER`
- `datastore_name` (string) - Required; The name of the datastore the image files were downloaded to.
    * Environment variable: `VCENTER_DATASTORE`
- `cluster_name` (string) - Optional, but recommended; Target cluster name the template will be imported into. This is used to get the resource pool ID. If left blank, will try to use the default datacenter and resource pool. 
    * Environment variable: `VCENTER_CLUSTER`
- `folder_name` (string) - Optional, but recommended; The name of the vCenter folder where the template should be placed; Will try to use the default root folder of the datacenter if left blank.
    * Environment variable: `VCENTER_FOLDER`
- `respool_name` (string) - Optional, but recommended; The name of the target resource pool where the template should reside; Will try to use the default resource pool if left blank.
    * Environment variable: `VCENTER_RESOURCE_POOL`
- `output_dir` (string) - Required; The path to an accessible datastore where the downloaded image files should be placed; Process will automatically place image files into their own folder based on the image name, so this doesn't need to be included. (ex: 'H:\\lab-servers' or '/lab-servers/')
    * Environment variable: `OUTPUTDIR`
- `download_uri` (string) - Required; The Artifactory download URI of the image artifact that should be downloaded. The artifact should be either an OVA, OVF, or VMTX file. Any standard vSphere files associated with one of those image types will also be downloaded automatically. 


## Output Data

None


## Basic Example Usage

```hcl
data "artifactory-import" "basic-example" {
    artifactory_token   = "artifactory_token"
    artifactory_server  = "https://server.domain.com:8081/artifactory/api"

    vcenter_server      = "vc01.domain.com"
    vcenter_user        = "jsmith@domain.com"
    vcenter_password    = "P@ssW0rd!"
    datacenter_name     = "Lab"
    datastore_name      = "lab-servs"
    cluster_name        = "lab-cluster01"
    folder_name         = "Templates"
    respool_name        = "cluster01-pool"

    output_dir          = "H:\\lab-servs
    download_uri        = "https://server.domain.com:8081/artifactory/api/storage/lab-repo/win/win2022.ova"
}
```

## Basic Example Usage With Variable References
In this example, we use the artifactory data source to get the image's download URI, which is an output of that data source, and feed it into the artifact import data source.

```hcl
data "artifactory" "basic-example" {
	# Provide via environment variables
	artifactory_token     = var.artif_token  
	artifactory_server    = var.artif_server

	artifact_name = "test-artifact"
	file_type     = "ova"
		
	filter = {
		release = "latest-stable"
	}
}

data "artifactory-import" "basic-example" {
	artifactory_token   = var.artif_token
	artifactory_server  = var.artif_server

	vcenter_server      = var.vc_server
	vcenter_user        = var.vc_user
	vcenter_password    = var.vc_password
	datacenter_name     = var.vc_datacenter
	datastore_name      = var.vc_datastore
	cluster_name        = var.vc_cluster
	folder_name         = var.vc_folder
	respool_name        = var.vc_respool

	output_dir          = var.output_directory
	download_uri        = data.artifactory.basic-example.download_uri
}
```