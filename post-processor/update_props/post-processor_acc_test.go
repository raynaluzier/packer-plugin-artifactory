package artifactUpdateProps

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

var artifactUri string
var statusCode  string

const testDirName      = "test-directory"
const testArtifactName = "test-artifact.txt"
const artifactSuffix   = ""
const artifactContents = "Just some test content."
var kvProps []string
var downloadUri string

var token = ""    // Update for testing
var server = ""   // Update for testing

// Run with: PACKER_ACC=1 go test -count 1 -v ./... -timeout=120
func TestAccPostProcessorUpdate_Artifactory(t *testing.T) {

	// Prep test artifact
	testDirPath        := common.CreateTestDirectory(testDirName)
	testArtifactPath   := common.CreateTestFile(testDirPath, testArtifactName, artifactContents)
	uploadTestArtifact := true
	testArtifactUpdate := SetTemplate(testArtifactPath)
	kvProps = append(kvProps,"release=stable")

	log.Println("Test Directory Created: " + testDirPath)
	log.Println("Test Artifact Created: " + testArtifactPath)
 
	testCase := &acctest.PluginTestCase{
		Name: "artifactory_postprocessor_update_test",
		Setup: func() error {
			status, err := tasks.SetupTest(server, token, testArtifactPath, artifactSuffix, kvProps, uploadTestArtifact)
			fmt.Println("Status of setup: " + status)

			if err != nil {
				log.Fatal(err)
			}
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
		Template: testArtifactUpdate,
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

			propLog := "Property assignment to artifact was successful."
			if matched, _ := regexp.MatchString(propLog+".*", logsString); !matched {
				t.Fatalf("logs do not contain expected output %q", logsString)
			}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}

func SetTemplate(testArtifactPath string) string {
	template := `
	packer {
		required_plugins {
			artifactory = {
				version = ">= 1.0.8"
				source  = "github.com/raynaluzier/artifactory"
			}
		}
	}

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

	source "null" "basic-example" {
		communicator = "none"
	}

	build {
		sources = [
			"source.null.basic-example"
		]

		post-processor "artifactory-update-props" {
		    artifactory_token     = var.artif_token  
        	artifactory_server    = var.artif_server 

			artifact_uri = "${var.artif_server}/storage/test-packer-plugin/test-artifact.txt"
			properties   = {
				release = "latest-stable"
			}
		}
	}
	`
	return template
}