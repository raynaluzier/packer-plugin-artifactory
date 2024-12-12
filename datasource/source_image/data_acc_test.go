package artifactImage

import (
	_ "embed"
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

// go:embed test-fixtures/template.pkr.hcl
var testDatasourceHCL2Basic string

const testDirName      = "test-directory"
const testArtifactName = "test-artifact.txt"
const artifactSuffix   = ""
const artifactContents = "Just some test content."

const testRepoName       = "test-packer-plugin"
const testRepoConfigName = "repository-config.json"
const repoConfigContents = "{ \"key\": \"" + testRepoName + "\",\"rclass\": \"local\", \"description\": \"temporary; test repo for packer plugin acceptance testing\"}"

var artifactUri string
var statusCode  string

// Run with: PACKER_ACC=1 go test -count 1 -v ./datasource/source_image/data_acc_test.go -timeout=120
func TestAccDatasource_Artifactory(t *testing.T) {
	// Prep test artifact
	testDirPath  := common.CreateTestDirectory(testDirName)
	testArtifact := common.CreateTestFile(testDirPath, testArtifactName, artifactContents)
		
	// Prep Repo Config File
	configFilePath := common.CreateTestFile(testDirPath, testRepoConfigName, repoConfigContents)

	testCase := &acctest.PluginTestCase{
		Name: "artifactory_datasource_basic_test",
		Setup: func() error {
			artifactUri = TestAccDatasource_TestSetup(t, configFilePath, testArtifact)
			log.Println("Test Artifact Created: " + artifactUri)
			return nil
		},
		Teardown: func() error {
			statusCode = TestAccDatasource_TestTeardown(t, configFilePath, testArtifact, testDirPath, testRepoName, artifactUri)
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
			
			artifactUriLog := fmt.Sprintf("null.basic-example: artifact URI: %s", artifactUri)
			if matched, _ := regexp.MatchString(artifactUriLog+".*", logsString); !matched {
				t.Fatalf("logs doesn't contain expected output %q", logsString)
			}
			return nil
		},
	}
	acctest.TestPlugin(t, testCase)
}

// Must have Artifactory instance licensed at Pro or higher, access to create/remove repos and artifacts
func TestAccDatasource_TestSetup(t *testing.T, configFilePath, testArtifact string) string {
	datasource := Datasource{
		config: Config{},
	}
	var kvProps []string
	kvProps = append(kvProps,"release=latest-stable")

	artifactUri, err := tasks.SetupTest(datasource.config.ArtifactoryServer, datasource.config.AritfactoryToken, testRepoName, configFilePath, testArtifact, artifactSuffix, kvProps)
	if err != nil {
		log.Fatal(err)
	}
	return artifactUri
}

func TestAccDatasource_TestTeardown(t *testing.T, configFilePath, testArtifact, testDirPath, testRepoName, artifactUri string) string {
	datasource := Datasource{
		config: Config{},
	}

	// Delete local test files
	common.DeleteTestFile(configFilePath)
	common.DeleteTestFile(testArtifact)

	// Delete locat test directory
	common.DeleteTestDirectory(testDirPath)

	// Delete test artifact from test repo, then delete test repo
	statusCode := tasks.TeardownTest(datasource.config.ArtifactoryServer, datasource.config.AritfactoryToken, testRepoName, artifactUri)
	if statusCode == "204" {
		fmt.Println("Test environment successfully torn down.")
	} else {
		fmt.Println("Unable to delete test artifact and/or test repo.")
	}
	return statusCode
}