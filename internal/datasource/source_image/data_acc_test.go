package artifactImage

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"testing"

	"github.com/hashicorp/packer-plugin-sdk/acctest"
	"github.com/raynaluzier/artifactory-go-sdk/common"
	"github.com/raynaluzier/artifactory-go-sdk/tasks"
)

const testDatasourceHCL2Basic = `
	packer {
		required_plugins {
			artifactory = {
				version = ">= 1.0.44"
				source  = "github.com/raynaluzier/artifactory"
			}
		}
	}

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

	data "artifactory" "basic-example" {
		# Provide via environment variables
		artifactory_token     = var.artif_token  
		artifactory_server    = var.artif_server

		artifact_name = "test-artifact"
		file_type     = "txt"
		
		filter = {
			release = "latest-stable"
		}
	}

	locals {
		name         = data.artifactory.basic-example.name
		create_date  = data.artifactory.basic-example.creation_date
		artifact_uri = data.artifactory.basic-example.artifact_uri
		download_uri = data.artifactory.basic-example.download_uri
	}

	source "null" "basic-example" {
		communicator = "none"
	}

	build {
		sources = [
			"source.null.basic-example"
		]

		provisioner "shell-local" {
			inline = [
				"echo artifact name: ${local.name}",
				"echo artifact creation date: ${local.create_date}",
				"echo artifact URI: ${local.artifact_uri}",
				"echo artifact download URI: ${local.download_uri}"
			]
		}
	}
	`

var artifactUri string
var statusCode  string

const testDirName      = "test-directory"
const testArtifactName = "test-artifact.txt"
const artifactContents = "Just some test content."
var kvProps []string
const uploadTestArtifact = true
const setAsOva           = false

var token = ""    // Update for testing
var server = ""   // Update for testing

// Run with: PACKER_ACC=1 go test -count 1 -v ./... -timeout 120m
func TestAccDatasource_Artifactory(t *testing.T) {

	// Prep test artifact
	testDirPath  := common.CreateTestDirectory(testDirName)
	testArtifactPath := common.CreateTestFile(testDirPath, testArtifactName, artifactContents, setAsOva)
	kvProps = append(kvProps,"release=latest-stable")
	uploadTestArtifact := true

	log.Println("Test Directory Created: " + testDirPath)
	log.Println("Test Artifact Created: " + testArtifactPath)
 
	testCase := &acctest.PluginTestCase{
		Name: "artifactory_datasource_basic_test",
		Setup: func() error {
			artifactUri, err := tasks.SetupTest(server, token, testArtifactPath, kvProps, uploadTestArtifact)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Test Artifact Created: " + artifactUri)
			return nil
		},
		Teardown: func() error {
			// Deletes locally created test artifact and test directory
			common.DeleteTestFile(testArtifactPath)
			common.DeleteTestDirectory(testDirPath)
	
			// Deletes test artifact and repo from Artifactory
			statusCode := tasks.TeardownTest(server, token)
			if statusCode == "200" {
				log.Println("Test environment successfully torn down.")
			} else {
				log.Println("Unable to teardown test environment.")
			}
			log.Println("Status code of teardown operation: " + statusCode)

			return nil
		},
		Template: testDatasourceHCL2Basic,
		Type: "basic-example",
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}

			logs, err := os.Open(logfile)
			if err != nil {
				return fmt.Errorf("Unable find %s", logfile)
			}
			defer logs.Close()

			logsBytes, err := io.ReadAll(logs)
			if err != nil {
				return fmt.Errorf("Unable to read %s", logfile)
			}
			logsString := string(logsBytes)
			log.Println(logsString)

			artifactUriLog := fmt.Sprintf("null.basic-example: artifact URI: %s", artifactUri)
			if matched, _ := regexp.MatchString(artifactUriLog+".*", logsString); !matched {
				t.Fatalf("logs doesn't contain expected output %q", logsString)
			}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}