// Code generated by "packer-sdc mapstructure-to-hcl2"; DO NOT EDIT.

package artifactUpload

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

// FlatConfig is an auto-generated flat version of Config.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatConfig struct {
	AritfactoryToken  *string  `mapstructure:"artifactory_token" required:"true" cty:"artifactory_token" hcl:"artifactory_token"`
	ArtifactoryServer *string  `mapstructure:"artifactory_server" required:"true" cty:"artifactory_server" hcl:"artifactory_server"`
	SourcePath        *string  `mapstructure:"source_path" required:"true" cty:"source_path" hcl:"source_path"`
	ArtifactoryPath   *string  `mapstructure:"artifactory_path" required:"true" cty:"artifactory_path" hcl:"artifactory_path"`
	FolderName        *string  `mapstructure:"folder_name" required:"false" cty:"folder_name" hcl:"folder_name"`
	FileList          []string `mapstructure:"file_list" required:"true" cty:"file_list" hcl:"file_list"`
}

// FlatMapstructure returns a new FlatConfig.
// FlatConfig is an auto-generated flat version of Config.
// Where the contents a fields with a `mapstructure:,squash` tag are bubbled up.
func (*Config) FlatMapstructure() interface{ HCL2Spec() map[string]hcldec.Spec } {
	return new(FlatConfig)
}

// HCL2Spec returns the hcl spec of a Config.
// This spec is used by HCL to read the fields of Config.
// The decoded values from this spec will then be applied to a FlatConfig.
func (*FlatConfig) HCL2Spec() map[string]hcldec.Spec {
	s := map[string]hcldec.Spec{
		"artifactory_token":  &hcldec.AttrSpec{Name: "artifactory_token", Type: cty.String, Required: false},
		"artifactory_server": &hcldec.AttrSpec{Name: "artifactory_server", Type: cty.String, Required: false},
		"source_path":        &hcldec.AttrSpec{Name: "source_path", Type: cty.String, Required: false},
		"artifactory_path":   &hcldec.AttrSpec{Name: "artifactory_path", Type: cty.String, Required: false},
		"folder_name":        &hcldec.AttrSpec{Name: "folder_name", Type: cty.String, Required: false},
		"file_list":          &hcldec.AttrSpec{Name: "file_list", Type: cty.List(cty.String), Required: false},
	}
	return s
}
