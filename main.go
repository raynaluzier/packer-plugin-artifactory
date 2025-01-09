// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	artifactImage "packer-plugin-artifactory/datasource/source_image"
	artifactUpload "packer-plugin-artifactory/post-processor/artifact_upload"
	artifactUpdateProps "packer-plugin-artifactory/post-processor/update_props"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
	"github.com/hashicorp/packer-plugin-sdk/version"
)

var (
	// Version is the main version number that is being run at the moment.
	Version = "1.0.3"

	// VersionPrerelease is A pre-release marker for the Version. If this is ""
	// (empty string) then it means that it is a final release. Otherwise, this
	// is a pre-release such as "dev" (in development), "beta", "rc1", etc.
	VersionPrerelease = "dev"

	// PluginVersion is used by the plugin set to allow Packer to recognize
	// what version this plugin is.
	PluginVersion = version.NewRawVersion(Version + "-" + VersionPrerelease)
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterDatasource(plugin.DEFAULT_NAME, new(artifactImage.Datasource))
	pps.RegisterPostProcessor("upload", new(artifactUpload.PostProcessor))
	pps.RegisterPostProcessor("update-props", new(artifactUpdateProps.PostProcessor))
	pps.SetVersion(PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
