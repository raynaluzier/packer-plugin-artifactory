The JFrog Artifactory plugin can be used with HashiCorp Packer to identify, retrieve, and work with artifacts. This plugin comes with a data source to locate target artifacts and retrieve information necessary to work with them depending on the strategy you want to use, as well as two post-processors. One post-processor enables the upload of a newly created artifact into Artifactory and the other enables the assignment of one or more properties to an artifact.

Future components will include functionality to download artifacts.

### Installation

To install this plugin, copy and paste this code into your Packer configuration, then run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    artifactory = {
      source  = "github.com/raynaluzier/artifactory"
      version = ">=1.0.8"
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

- [artifactory](/packer/integrations/jfrog/artifactory/latest/components/data-source/source_image/README.md) - Filter and locate target artifacts and retrieve information necessary to work with them.

#### Post-Processors
- [artifactory-upload](/packer/integrations/jfrog/artifactory/latest/components/post-processors/artifact_upload/README.md) - Upload artifact to Artifactory.
- [artifactory-update-props](/packer/integrations/jfrog/artifactory/latest/components/post-processors/update_props/README.md) - Assign one or more properties to an artifact in Artifactory.

### Authentication
There are several ways to provide credentials for JFrog Artifactory authentication, which uses a bearer token when making each underlying request. The following authentication methods are supported:

- Static credentials
- Environment variables
    - Exported
    - .env File

#### Static Credentials
Static credentials can be provided using the server API address and identity token as follows:

```hcl
data "artifactory" "basic-example" {
  artifactory_server = "https://server.com:8081/artifactory/api"
  artifactory_token  = "1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t"
}
```
#### Environment Variables
Credentials can be provided via the `ARTIFACTORY_SERVER` and `ARTIFACTORY_TOKEN` environment variables, which represent the target Artifactory server instance's API address and Artifactory identity token of the account with access to query artifacts and execute operations.

**Usage:  Exported**
```
$ export ARTIFACTORY_SERVER="https://server.com:8081/artifactory/api"
$ export ARTIFACTORY_TOKEN="1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t"
packer build template.pkr.hcl
```

**Usage:  .env**
If using a `.env` file, include it at the root of your configuration. Ensure an entry is added for it to your `.gitignore` file to prevent storing sensitive data in your repository.
```
ARTIFACTORY_SERVER=https://server.com:8081/artifactory/api
ARTIFACTORY_TOKEN=1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t
```

### Troubleshooting
To gather additional information about the processes happening, change the logging level to 'DEBUG'. This can be done by setting Logging in the data source configuration or via an environment variable.

Data Source Example:
```hcl
data "artifactory" "basic-example" {
  artifactory_server  = "https://server.com:8081/artifactory/api"
  artifactory_token   = "1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t"
  artifactory_logging = "DEBUG"
}
```

Exported Environment Variable:
```
$ export ARTIFACTORY_SERVER="https://server.com:8081/artifactory/api"
$ export ARTIFACTORY_TOKEN="1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t"
$ export ARTIFACTORY_LOGGING="DEBUG"
packer build template.pkr.hcl
```

Environment Variable in .env:
```
ARTIFACTORY_SERVER=https://server.com:8081/artifactory/api
ARTIFACTORY_TOKEN=1a2b3c4d5e6f7g8h9i0j1k2l3m4n5o6p7q8r9s0t
ARTIFACTORY_LOGGING=DEBUG
```