package artifactUpload

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
const testArtifactName = "test-artifact.txt"  // this gets renamed to an ova file later
const artifactContents = "Just some test content."
var kvProps []string
var downloadUri string
const setAsOva = true

var token = ""    // Update for testing
var server = ""   // Update for testing

// Run with: PACKER_ACC=1 go test -count 1 -v ./... -timeout 120m
func TestAccPostProcessorUpload_Artifactory(t *testing.T) {

	// Prep test artifact
	testDirPath        := common.CreateTestDirectory(testDirName)
	// Returns $HOME_DIR/test-directory/test-artifact.ova
	testArtifactPath   := common.CreateTestFile(testDirPath, testArtifactName, artifactContents, setAsOva)  // the artifact will be renamed to test-artifact.ova
	uploadTestArtifact := false  // Don't need to upload artifact as part of test setup; test itself will do this.
	testArtifactUpload := SetTemplate(testArtifactPath)
	kvProps = append(kvProps,"release=latest-stable")

	log.Println("Test Directory Created: " + testDirPath)
	log.Println("Test Artifact Created: " + testArtifactPath)
 
	testCase := &acctest.PluginTestCase{
		Name: "artifactory_postprocessor_upload_test",
		Setup: func() error {
			status, err := tasks.SetupTest(server, token, testArtifactPath, kvProps, uploadTestArtifact)
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
		Template: testArtifactUpload,
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

			downloadUriLog := "Download URI for new artifact: " + downloadUri
			if matched, _ := regexp.MatchString(downloadUriLog+".*", logsString); !matched {
				t.Fatalf("logs doesn't contain expected output %q", logsString)
			}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}

func SetTemplate(testDirPath string) string {
	newPath := common.EscapeSpecialChars(testDirPath)  // $HOME_DIR/test-directory/

	template := `
	packer {
		required_plugins {
			artifactory = {
				version = ">= 1.0.27"
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

	source "null" "basic-example" {
		communicator = "none"
	}

	build {
		sources = [
			"source.null.basic-example"
		]
		
		post-processor "artifactory-upload" {
		    artifactory_token     = var.artif_token  
        	artifactory_server    = var.artif_server 
			
			source_path = "` + newPath + `"
			target_path = "/test-packer-plugin"
			image_name  = "test-artifact"
			image_type  = "ova"
		}
	}
	`
	return template
}