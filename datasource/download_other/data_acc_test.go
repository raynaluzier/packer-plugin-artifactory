package artifactDownloadOther

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
const testArtifact     = "testfile1.txt"
const artifactSuffix   = ""
const artifactContents = "Just some test content."
var kvProps []string
var downloadUri string

var token  = ""    // Update for testing
var server = ""   // Update for testing

// Run with: PACKER_ACC=1 go test -count 1 -v ./... -timeout 120m
func TestAccDatasourceDownload_Artifactory(t *testing.T) {
	// Prep test artifact
	testDirPath        := common.CreateTestDirectory(testDirName)
	testArtifactPath   := common.CreateTestFile(testDirPath, testArtifact, artifactContents)
	uploadTestArtifact := true
	testArtifactDownloadOther := SetTemplate(testDirPath)
	kvProps = append(kvProps,"release=latest-stable")

	log.Println("Test Directory Created: " + testDirPath)
	log.Println("Test Artifact Created: " + testArtifactPath)
 
	testCase := &acctest.PluginTestCase{
		Name: "artifactory_datasource_download_other_test",
		Setup: func() error {
			artifactUri, err := tasks.SetupTest(server, token, testArtifactPath, artifactSuffix, kvProps, uploadTestArtifact)
			if err != nil {
				log.Fatal(err)
			}
			log.Println("Test Artifact Created: " + artifactUri)
			common.DeleteTestFile(testArtifactPath) // Delete test file locally so we can test downloading here
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
		Template: testArtifactDownloadOther,
		Type: "basic-example",
		Check: func(buildCommand *exec.Cmd, logfile string) error {
			if buildCommand.ProcessState != nil {
				if buildCommand.ProcessState.ExitCode() != 0 {
					return fmt.Errorf("Bad exit code. Logfile: %s", logfile)
				}
			}

			items, _ := os.ReadDir(testDirPath)
			for _, item := range items {
				if item.Name() == testArtifact {
					log.Println("Test artifact downloaded successfully.")
				} else {
					log.Println("Test artifact download failed.")
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
			
			outputLog := "Test artifact downloaded successfully."
			if matched, _ := regexp.MatchString(outputLog+".*", logsString); !matched {
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
				version = ">= 1.0.26"
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
	
	data "artifactory-download-other" "basic-example" {
		artifactory_token     = var.artif_token  
		artifactory_server    = var.artif_server

		output_dir       = "` + newPath +`"
		artifactory_path = "/test-packer-plugin"
		file_list        = ["` + testArtifact + `"]
	}

	source "null" "basic-example" {
		communicator = "none"
	}

	build {
		sources = [
			"source.null.basic-example"
		]
	}
	`
	return template
}