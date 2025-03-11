Type:  `artifactory-import`

# JFrog Artifactory Data Source

The Artifactory data source is used to download an artifact image of type OVA, OVF, or VMTX and it's standard associated image files that are stored in JFrog Artifactory to an accessible datastore path, convert the image to a VMX, import into vCenter, and then mark it as a template which is then ready for the vsphere-clone builder to consume.

If the download portion is unnecessary (say the download occurred in a previous run, but the run had to be stopped or didn't complete successfully for some reason), there is an option to pick up the process starting with the conversion step. The conversion process checks the image file first, so if the file is already in VMX format, it will move on to importing to vCenter. If the file is still in OVA, OVF, or VMTX format, then the conversion will be done.

It is assumed that the datastore path where the image files were downloaded is the same where the converted image files should be stored. This is the default behavior, but deviations from this are checked and handled behind the scenes.

This component is meant to be used in conjunction with the [vSphere-Clone](https://developer.hashicorp.com/packer/integrations/hashicorp/vsphere/latest/components/builder/vsphere-clone) Builder.


## Advisements
* When converting an image that was thin provisioned or thick provisioned, eager zeroed to start, the OVFTOOL automatically converts the image to be thick provisioned, lazy zeroed. Converting the resulting template back to a virtual machine and attempting to power it on may result in the following error: **"Unsupported or invalid disk type 2 for 'scsi0:0'. Ensure that the disk has been imported."** 
    ![Unsupported or Invalid Disk Error](https://github.com/raynaluzier/packer-plugin-artifactory/tree/main/docs/datasources/unsupported_disk_error.jpg)

To resolve this, edit the settings of the virtual machine. Expand the hard disk settings and change the **Virtual Device Node** to **IDE 0**. Repeat this for each disk. The machine will then power on successfully.
    ![Edit Settings](https://github.com/raynaluzier/packer-plugin-artifactory/tree/main/docs/datasources/edit_settings_disk.jpg)

* When downloading, if the files already exist in the target location, they will be overwritten. 

* When downloading and/or converting image files, the files are placed into a directory named after the image. 
Ex: If the output directory is H:\\lab-servs, the image file 'win2022.ova' will be placed in H:\\lab-servs\\win2022\\win2022.ova, and when the OVA is unpackaged, the resulting files will be in H:\\lab-servs\\win2022\\.


## Housekeeping
* Artifactory property key/values, artifact URIs, download URIs, Artifactory paths (/repo/folder/...), and file names are **CASE SENSITIVE**. There are a few exceptions, however, it's best to assume case sensitivity for successful outcomes. This is a behavior of the Artifactory API and not something we can control. 

* This process does NOT cleanup any image files that remain after the image conversion to VMX.

* If opting to use only the convert/import piece without first downloading AND the image files are OVA/OVF format, if the parent directory of the image file is named differently than the image file itself (for example: "E:\lab\win2022.ovf" or "/lab/testing/rhel9.ova"), the OVFTOOL will automatically place the converted image files into a sub-directory named after the image file which is a behavior that can't be changed. So using our examples, this results in "E:\lab\win2022\win2022.vmx" or "/lab/testing/rhel9/rhel9.vmx". This means that the original OVA/OVF files will still reside in (our example) "E:\lab" or "/lab/testing" while the VMX files are in a sub-directory. Especially as this resulting path will be used to import into vCenter, this may not be ideal behavior. So to avoid this situation for convert/import-only scenarios, ensure the source files reside in a directory that's the same as the image name (ex: "E:\lab\win2022\win2022.ovf" or "/lab/testing/rhel9/rhel9.ova").

* An OVF image is typically comprised of an **.ovf**, **.mf**, and one or more **-disk#.vmdk** (ex: win2022-disk1.vmdk) files. When the OVF is unpacked, it creates a disk file(s) that's the same name as the disk file(s) included with the OVF. Because we are placing the converted image files into the same directory as the source image (which is the same directory that is using during the import into vCenter), this creates a file conflict and the process fails. Therefore, when the image type is OVF, the OVF files will first be moved into a subdirectory called `ovf_files` and then the conversion will kick off. The converted files can then go to the intended source directory without conflict. OVA and VMTX files do not have this issue.

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
- `datacenter_name` (string) - Required; Target datacenter name the template will be imported into; Will try to use the default datacenter if this is left blank.
    * Environment variable: `VCENTER_DATACENTER`
- `datastore_name` (string) - Required; The name of the datastore the image files were downloaded to.
    * Environment variable: `VCENTER_DATASTORE`
- `cluster_name` (string) - Required; Target cluster name the template will be imported into. This is used to get the resource pool ID. If left blank, will try to use the default datacenter and resource pool. 
    * Environment variable: `VCENTER_CLUSTER`
- `folder_name` (string) - Optional, but recommended; The name of the vCenter folder where the template should be placed; Will try to use the default root folder of the datacenter if left blank.
    * Environment variable: `VCENTER_FOLDER`
- `respool_name` (string) - Optional, but recommended; The name of the target resource pool where the template should reside; Will try to use the default resource pool if left blank.
    * Environment variable: `VCENTER_RESOURCE_POOL`
- `import_no_download` (bool) - Optional; Whether we should skip the initial download process from Artifactory in the event the image file(s) have already been downloaded, maybe from a previous run/process. This signals the workflow to check the image type and convert if necessary, then import the image into vCenter and mark as a template. Defaults to FALSE.
    **If set to TRUE, then a value for `source_path` is required.**
- `output_dir` (string) - Optional; The path to an accessible datastore where the downloaded image files should be placed; Process will automatically place image files into their own folder based on the image name, so this doesn't need to be included. (ex: 'H:\\lab-servers' or '/lab-servers/').
    **If `import_no_download` is set to FALSE (default), then values for `output_dir` and `download_uri` are required**
    * Environment variable: `OUTPUTDIR`
- `download_uri` (string) - Optional; The Artifactory download URI of the image artifact that should be downloaded. The artifact should be either an OVA, OVF, or VMTX file. Any standard vSphere files associated with one of those image types will also be downloaded automatically. 
    **If `import_no_download` is set to FALSE (default), then values for `output_dir` and `download_uri` are required**
- `ds_image_path` (string) - Optional; Usually for Linux paths; the datastore path of the image without the mount point portion (see **USAGE** below for specifics). This is used to build the `vmPathName` that vCenter uses to import the VMX into vCenter before marking it as a VM Template.

  Usually in cases where the Packer machine is a Linux box, it will have a mount point/share to the datastore in the form of (for example) '/mnt/share/dev-servers/' where '/mnt/share' is the mount point seen from the Packer machine and '/dev-servers/' would be a folder on the datastore where different machines/templates live. In our example, the path to the converted image from Packer's perspective is '/mnt/share/dev-servers/rhel9/rhel9.vmx'.  

  Packer needs to be able to access '/mnt/share/dev-servers' (as in our example) to validate that the image files exist before we try to import them. However, when we build the required `vmPathName` that points to the image's VMX file stored on the datastore mounted in vCenter (after conversion), vCenter has no idea where '/mnt/share/dev-servers/rhel9/rhel9.vmx' is because all it sees from the datastore is '/dev-servers/rhel9/rhel9.vmx'. 

  Since mount point/share names and structures can vary wildly depending on individual needs/preferences, there's no way to automate the detection of what this portion of the path would be and edit it out. We don't typically have this issue in Windows because we're mapping a drive letter to the datastore (for example:  E:\\Lab\\dev-servers), drive letters will always be a single character that we can trim off, and the remaining path is a path that should exist on the datastore accessible to vCenter. Therefore, if you are running Packer from a machine that's using a mount point/share where part of this path doesn't exist on the datastore, then please use the `ds_image_path` input to provide the datastore folder path. 
  
  **USAGE:**
  * If `import_no_download` is set to TRUE, `ds_image_path` should be the same as the `source_path` EXCEPT without the mount/share information.
      Ex:  import_no_download = true
           source_path = "/mnt/share/lab-servers/rhel9/rhel9.ova"
           ds_image_path = "/lab-servers/rhel9/rhel9.ova"
           
  * If `import_no_download` is set to FALSE, `ds_image_path` should be the same as the `output_dir` EXCEPT without the mount/share information.
      Ex:  import_no_download = false // or this input is excluded entirely
           source_path = "/mnt/share/lab-servers/"
           ds_image_path = "/lab-servers/"

- `source_path` (string) - Optional; Full file path (ex: `/mnt/share/folder/image1234/image1234.ova`) to the source image file (should be OVA, OVF, VMTX, or VMX if it's already in that format). As the image files must be in VMX format (essentially a VM) first for this plugin component to do the import, the image files will be examined to determined whether the image needs to be converted to VMX format. If not, the conversion step is skipped and the import will proceed.
    **If `import_no_download` is set to TRUE, then a value for `source_path` is required.**


## Output Data

None


## Basic Example Usage, Downloading Artifacts and Importing to vCenter

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

    output_dir          = "/mnt/lab/lab-servs/"
    download_uri        = "https://server.domain.com:8081/artifactory/api/storage/lab-repo/win/win2022.ova"
    ds_image_path       = "/lab-servs/"                  // <-- Note this is the same as the output_dir EXCEPT without the mount/share info
}
```

## Basic Example Usage, Skipping the Artifact Download Step

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

    import_no_download  = true
    source_path         = "/mnt/share/rhel9/rhel9.ova"
    ds_image_path       = "/rhel9/rhel9.ova"             // <-- Note this is the same as the source_path EXCEPT without the mount/share info
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

## FAQ
* How do I separate these image files from other image files in my output directory?
  - When the image files are downloaded, they are placed into their own image name-based folder. For example:  win2022-25-01-23.ova would be downloaded to `.../win2022-25-01-23/win2022-25-01-23.ova`.

* What if I have existing files with the same name in my output directory?
  - The downloaded files will overwrite the existing files. If this is not desired, set the 'import_no_download' flag and use the existing image file in your output directory as the source.
  - Alternatively, specify a different directory.

* Does the output directory I specify have to exist first?
  - No, the process will create the directory path if it doesn't already exist.
  - **As this process imports the template into vCenter, you should specify an accessible datastore as the output directory.**

* Is the download URI case sensitive?
  - It depends on the behavior of the particular Artifactory API. Typically, paths and artifact names are case sensitive. However, it is NOT case sensitive in this instance.

* What if I don't want my OVA converted into a VM Template?
  - Use the 'artifactory-download-other' component to download the OVA into a desired location. But then you will have to manually add the OVA to vCenter as that functionality is not currently available.

* What if I want to pull down an ISO for use with the 'vsphere-iso' plugin?
  - Use the 'artifactory-download-other' component to download the OVA into a desired location. Use this location as the ISO source with the 'vsphere-iso' plugin.

* What if I have other files, like scripts or YAML files, that I need for the configuration of the image?
  - Use the 'artifactory-download-other' component to download the supporting files into a desired location.

* When I try to boot the template machine that was created after converting my OVA/OVF, I get an "Unsupported or invalid disk type 2 for 'scsi0:0'" error. Can you fix this?
  - Unfortunately, no. This is a known potential behavior when the OVFTool converts from a differing disk format type (like `thin provisioned`) to it's resulting `thick provisioned lazy zeroed` type, and we do not maintain this tool. However, steps to remediate can be found in the **'Advisements'** section above.

* I get an error when trying to use a Windows path like this: 'E:\lab-servs\'. This is the right path. Why is this happening?
  - Windows-based paths must be properly escaped. Instead, use a path like this: 'E:\\lab-servs\\'. This is not an issue with Linux-based paths.

* Can this component add image files into a Content Library?
  - No, this component does not work with a Content Library. The Content Library is an entirely different process model with it's own separate set of APIs.

* My Packer system runs Linux and I've got a share to the datastore mounted. Everything runs file until it gets to the import process and then it fails. Why is this happening?
  - The share on Packer is using a mount point which exists soley on the Packer system, but doesn't exist on the datastore itself. Therefore, it's unknown to vCenter. Use the `ds_image_path` input to specify the folder path on the datastore without the mount point portion. See above for examples.