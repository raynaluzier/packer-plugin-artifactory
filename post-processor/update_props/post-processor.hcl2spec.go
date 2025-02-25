// Code generated by "packer-sdc mapstructure-to-hcl2"; DO NOT EDIT.

package artifactUpdateProps

import (
	"github.com/hashicorp/hcl/v2/hcldec"
	"github.com/zclconf/go-cty/cty"
)

// FlatConfig is an auto-generated flat version of Config.
// Where the contents of a field with a `mapstructure:,squash` tag are bubbled up.
type FlatConfig struct {
	ArtifactoryToken   *string           `mapstructure:"artifactory_token" required:"true" cty:"artifactory_token" hcl:"artifactory_token"`
	ArtifactoryServer  *string           `mapstructure:"artifactory_server" required:"true" cty:"artifactory_server" hcl:"artifactory_server"`
	ArtifactUri        *string           `mapstructure:"artifact_uri" required:"true" cty:"artifact_uri" hcl:"artifact_uri"`
	ArtifactProperties map[string]string `mapstructure:"properties" required:"true" cty:"properties" hcl:"properties"`
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
		"artifact_uri":       &hcldec.AttrSpec{Name: "artifact_uri", Type: cty.String, Required: false},
		"properties":         &hcldec.AttrSpec{Name: "properties", Type: cty.Map(cty.String), Required: false},
	}
	return s
}
