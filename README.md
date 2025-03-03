# Packer Plugin JFrog Artifactory
The `Artifactory` plugin can be used with HashiCorp [Packer](https://www.packer.io) to locate, retrieve, import into vCenter as a VM Template, and upload and update custom images into Artifactory. For the list of available functionality for this plugin, see [docs](docs). 

NOTE: The `artifactory-import` component makes use of the VMware OVFTool. This must be installed on the machine where the plugin will be executed from to function properly.

## Installation

### Using pre-built releases

#### Using the `packer init` command

Starting from version 1.7, Packer supports a new `packer init` command allowing
automatic installation of Packer plugins. Read the
[Packer documentation](https://www.packer.io/docs/commands/init) for more information.

To install this plugin, copy and paste this code into your Packer configuration .
Then, run [`packer init`](https://www.packer.io/docs/commands/init).

```hcl
packer {
  required_plugins {
    artifactory = {
      version = ">= 1.0.12"
      source  = "github.com/raynaluzier/artifactory"
    }
  }
}
```

#### Manual installation

You can find pre-built binary releases of the plugin [here](https://github.com/raynaluzier/packer-plugin-artifactory/releases).
Once you have downloaded the latest archive corresponding to your target OS,
uncompress it to retrieve the plugin binary file corresponding to your platform.
To install the plugin, please follow the Packer documentation on
[installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).


### From Sources

If you prefer to build the plugin from sources, clone the GitHub repository
locally and run the command `go build` from the root
directory. Upon successful compilation, a `packer-plugin-artifactory` plugin
binary file can be found in the root directory.
To install the compiled plugin, please follow the official Packer documentation
on [installing a plugin](https://www.packer.io/docs/extending/plugins/#installing-plugins).


### Configuration

For more information on how to configure the plugin, please read the
documentation located in the [`docs/`](docs) directory.
