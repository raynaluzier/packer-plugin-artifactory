The JFrog Artifactory plugin can be used with HashiCorp Packer to identify, retrieve, and upload artifacts. It provides three data sources and three post-processors.

The data sources locate artifacts, download supporting files, and optionally download and import OVA, OVF, or VMTX images into vCenter as templates. Image files are checked and converted when necessary before import.

The post-processors upload image files, upload supporting files, and update artifact properties.

NOTE: The `artifactory-import` component makes use of the VMware OVFTool. This must be installed on the machine where the plugin will be executed from to function properly.


### Installation

To install this plugin, copy and paste this code into your Packer configuration, then run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    artifactory = {
      source  = "github.com/raynaluzier/artifactory"
      version = ">=1.0.27"
    }
  }
}
```

Alternatively, you can use `packer plugins install` to manage installation of this plugin.

```sh
$ packer plugins install github.com/raynaluzier/artifactory
```

### Components
The following components are available with the Artifactory plugin.

#### Data Sources
- [artifactory](https://github.com/raynaluzier/packer-plugin-artifactory/blob/main/docs/datasources/datasource.mdx) - Filter and locate target artifacts and retrieve information necessary to work with them.

- [artifactory-import](https://github.com/raynaluzier/packer-plugin-artifactory/blob/main/docs/datasources/artifact_import.mdx) - Optionally download image artifacts (OVA, OVF, or VMTX) and/or import into vCenter as a template (for use with vsphere-clone builder plugin).

- [artifactory-download-other](https://github.com/raynaluzier/packer-plugin-artifactory/blob/main/docs/datasources/download_other.mdx) - Download non-image artifacts that may support the image itself, such as script files, metadata, etc.

#### Post-Processors
- [artifactory-upload](https://github.com/raynaluzier/packer-plugin-artifactory/blob/main/docs/post-processors/artifact_upload.mdx) - Upload newly built image artifacts (OVA, OVF, or VMTX) to JFrog Artifactory.

- [artifactory-update-props](https://github.com/raynaluzier/packer-plugin-artifactory/blob/main/docs/post-processors/update_props.mdx) - Update the properties of an existing or newly created image artifact stored in Jfrog Artifactory.

- [artifactory-upload-other](https://github.com/raynaluzier/packer-plugin-artifactory/blob/main/docs/post-processors/upload_other.mdx) - Upload non-image artifacts that may support the image itself, such as script files, metadata, etc.

### Authentication
There are several ways to provide credentials for JFrog Artifactory authentication, which uses a bearer token when making each underlying request. The following authentication methods are supported:

- Static credentials
- Variables file
- Environment variables
    - Exported
    - .env File

#### Static Credentials
Static credentials, though not recommended, can be provided using the server API address and identity token as follows:

```hcl
data "artifactory" "basic-example" {
  artifactory_server = "https://server.domain.com:8081/artifactory/api"
  artifactory_token  = "1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t"
}

data "artifactory-import" "basic-example" {
  artifactory_server = "https://server.domain.com:8081/artifactory/api"
  artifactory_token  = "1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t"

  vcenter_server     = "vc01.domain.com"
  vcenter_user       = "auser@domain.com"
  vcenter_password   = "MyP@$$w0rd!"
}
```
#### Variables File
Input variables values can be supplied through the use of either a `.pkrvars.hcl` or `.auto.pkrvars.hcl` file.

Ensure a corresponding input variable definition is present for each as shown below, and use the `var.` reference in the data source.

**Usage: Variables File**
Variable Definitions:
```
variable "artif_token" {
  type        = string
  description = "Identity token of the Artifactory account with access to execute commands"
  sensitive   = true
  default     = env("ARTIFACTORY_TOKEN")
}

variable "artif_server" {
  type        = string
  description = "The Artifactory API server address"
  default     = env("ARTIFACTORY_SERVER")
}

variable "vc_server" {
  type        = string
  description = "vCenter Server FQDN/IP address"
  default     = env("VCENTER_SERVER")
}

variable "vc_user" {
  type        = string
  description = "vCenter User account"
  default     = env("VCENTER_USER")
}

variable "vc_password" {
  type        = string
  description = "vCenter User account password"
  sensitive   = true
  default     = env("VCENTER_PASSWORD")
}
```
Set the Variables in a .auto.pkrvars.hcl File:
```
artif_server = "https://server.domain.com:8081/artifactory/api"
artif_token  = "1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t"
vc_server    = "vc01.domain.com"
vc_user      = "auser@domain.com"
vc_password  = "MyP@$$w0rd!"

```

#### Environment Variables
Credentials can be provided via the `ARTIFACTORY_SERVER` and `ARTIFACTORY_TOKEN` environment variables, which represent the target Artifactory server instance's API address and Artifactory identity token of the account with access to query artifacts and execute operations.

Ensure a corresponding input variable definition is present for each as shown below, and use the `var.` reference in the data source.

**Usage:  Exported**
Variable Definitions:
```
variable "artif_token" {
  type        = string
  description = "Identity token of the Artifactory account with access to execute commands"
  sensitive   = true
  default     = env("ARTIFACTORY_TOKEN")
}

variable "artif_server" {
  type        = string
  description = "The Artifactory API server address"
  default     = env("ARTIFACTORY_SERVER")
}

variable "vc_server" {
  type        = string
  description = "vCenter Server FQDN/IP address"
  default     = env("VCENTER_SERVER")
}

variable "vc_user" {
  type        = string
  description = "vCenter User account"
  default     = env("VCENTER_USER")
}

variable "vc_password" {
  type        = string
  description = "vCenter User account password"
  sensitive   = true
  default     = env("VCENTER_PASSWORD")
}
```
Set the Environment Variables:
```
export ARTIFACTORY_SERVER=https://server.domain.com:8081/artifactory/api
export ARTIFACTORY_TOKEN=1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t
export VCENTER_SERVER=vc01.domain.com
export VCENTER_USER=auser@domain.com
export VCENTER_PASSWORD='MyP@$$w0rd!'
```

By default, vCenter TLS certificates are verified. For legacy environments with self-signed certificates, set `VSPHERE_INSECURE_SKIP_TLS_VERIFY=true`; use this only when certificate verification cannot be configured.

**Usage:  .env**
If using a `.env` file, include it at the root of your configuration. Ensure an entry is added for it to your `.gitignore` file to prevent storing sensitive data in your repository.

Ensure a corresponding input variable definition is present for each as shown below, and use the `var.` reference in the data source.

Variable Definitions:
```
variable "artif_token" {
  type        = string
  description = "Identity token of the Artifactory account with access to execute commands"
  sensitive   = true
  default     = env("ARTIFACTORY_TOKEN")
}

variable "artif_server" {
  type        = string
  description = "The Artifactory API server address"
  default     = env("ARTIFACTORY_SERVER")
}

variable "vc_server" {
  type        = string
  description = "vCenter Server FQDN/IP address"
  default     = env("VCENTER_SERVER")
}

variable "vc_user" {
  type        = string
  description = "vCenter User account"
  default     = env("VCENTER_USER")
}

variable "vc_password" {
  type        = string
  description = "vCenter User account password"
  sensitive   = true
  default     = env("VCENTER_PASSWORD")
}
```
Set the Environment Variables in the .env File:
```
ARTIFACTORY_SERVER=https://server.domain.com:8081/artifactory/api
ARTIFACTORY_TOKEN=1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t
VCENTER_SERVER=vc01.domain.com
VCENTER_USER=auser@domain.com
VCENTER_PASSWORD=MyP@$$w0rd!
```
