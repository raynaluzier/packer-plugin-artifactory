// Copyright (c) HashiCorp, Inc.
// SPDX-License-Identifier: MPL-2.0

package main

import (
	"fmt"
	"os"

	artifactImport "packer-plugin-artifactory/internal/datasource/artifact_import"
	artifactDownloadOther "packer-plugin-artifactory/internal/datasource/download_other"
	artifactImage "packer-plugin-artifactory/internal/datasource/source_image"
	artifactUpload "packer-plugin-artifactory/internal/post-processor/artifact_upload"
	artifactUpdateProps "packer-plugin-artifactory/internal/post-processor/update_props"
	artifactUploadOther "packer-plugin-artifactory/internal/post-processor/upload_other"

	"github.com/hashicorp/packer-plugin-sdk/plugin"
	"github.com/hashicorp/packer-plugin-sdk/version"
)

var (
	// Version is the main version number that is being run at the moment.
	Version = "1.1.1"

	// VersionPrerelease is A pre-release marker for the Version. If this is ""
	// (empty string) then it means that it is a final release. Otherwise, this
	// is a pre-release such as "dev" (in development), "beta", "rc1", etc.
	VersionPrerelease = ""

	// PluginVersion is used by the plugin set to allow Packer to recognize
	// what version this plugin is.

	//PluginVersion = version.NewRawVersion(Version + "-" + VersionPrerelease)
	PluginVersion = version.NewRawVersion(Version)
)

func main() {
	pps := plugin.NewSet()
	pps.RegisterDatasource(plugin.DEFAULT_NAME, new(artifactImage.Datasource))
	pps.RegisterDatasource("import", new(artifactImport.Datasource))
	pps.RegisterDatasource("download-other", new(artifactDownloadOther.Datasource))
	pps.RegisterPostProcessor("upload", new(artifactUpload.PostProcessor))
	pps.RegisterPostProcessor("upload-other", new(artifactUploadOther.PostProcessor))
	pps.RegisterPostProcessor("update-props", new(artifactUpdateProps.PostProcessor))
	pps.SetVersion(PluginVersion)
	err := pps.Run()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
