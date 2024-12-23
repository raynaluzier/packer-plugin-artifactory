The JFrog Artifactory plugin can be used with HashiCorp Packer to identify, retrieve, and work with artifacts. This plugin currently comes with a data source to locate target artifacts and retrieve information necessary to work with them depending on the strategy you want to use. 

Future components will include functionality to download artifacts, and create and manipulate artifacts and their properties after build time.

### Installation

To install this plugin, copy and paste this code into your Packer configuration, then run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    name = {
      # source represents the GitHub URI to the plugin repository without the `packer-plugin-` prefix.
      source  = "github.com/raynaluzier/artifactory"
      version = ">=0.0.1"
    }
  }
}
```

Alternatively, you can use `packer plugins install` to manage installation of this plugin.

```sh
$ packer plugins install github.com/raynaluzier/artifactory
```

### Components



#### Data Sources

- [data source](/packer/integrations/hashicorp/artifactory/latest/components/data-source/datasource) - Filter and locate target artifacts and retrieve information necessary to work with them.

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
  artifactory_server = "https://server.com:8081/artifactory/api"
  artifactory_token  = "1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t"
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
  default     = env("ARTIFACTORY_TOKEN")
}

variable "artif_server" {
  type        = string
  description = "The Artifactory API server address"
  default     = env("ARTIFACTORY_SERVER")
}
```
Set the Variables in a .auto.pkrvars.hcl File:
```
$ artif_server = "https://server.domain.com:8081/artifactory/api"
$ artif_token  = "1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t"
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
  default     = env("ARTIFACTORY_TOKEN")
}

variable "artif_server" {
  type        = string
  description = "The Artifactory API server address"
  default     = env("ARTIFACTORY_SERVER")
}

```
Set the Environment Variables:
```
$ export ARTIFACTORY_SERVER=https://server.domain.com:8081/artifactory/api
$ export ARTIFACTORY_TOKEN=1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t
```

**Usage:  .env**
If using a `.env` file, include it at the root of your configuration. Ensure an entry is added for it to your `.gitignore` file to prevent storing sensitive data in your repository.

Ensure a corresponding input variable definition is present for each as shown below, and use the `var.` reference in the data source.

Variable Definitions:
```
variable "artif_token" {
  type        = string
  description = "Identity token of the Artifactory account with access to execute commands"
  default     = env("ARTIFACTORY_TOKEN")
}

variable "artif_server" {
  type        = string
  description = "The Artifactory API server address"
  default     = env("ARTIFACTORY_SERVER")
}
```
Set the Environment Variables in the .env File:
```
ARTIFACTORY_SERVER=https://server.domain.com:8081/artifactory/api
ARTIFACTORY_TOKEN=1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t
```

### Troubleshooting
To gather additional information about the processes happening, change the logging level to 'DEBUG'. This can be done by setting Logging in the data source configuration or via an environment variable.

Data Source Example:
```
data "artifactory" "basic-example" {
  artifactory_server  = "https://server.domain.com:8081/artifactory/api"
  artifactory_token   = "1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t"
  artifactory_logging = "DEBUG"
}
```

Exported Environment Variable:
```
$ export ARTIFACTORY_SERVER=https://server.domain.com:8081/artifactory/api
$ export ARTIFACTORY_TOKEN=1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t
$ export ARTIFACTORY_LOGGING=DEBUG
```

Environment Variable in .env:
```
ARTIFACTORY_SERVER=https://server.domain.com:8081/artifactory/api
ARTIFACTORY_TOKEN=1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t
ARTIFACTORY_LOGGING=DEBUG
```