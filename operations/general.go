package operations

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"packer-plugin-artifactory/common"
	"path"
	"strings"
)

type Contents struct {
	Child	 		string
	IsFolder		bool
}

type artifJson struct {
	Repo			string 	`json:"repo"`
	Path			string	`json:"path"`
	Created			string	`json:"created"`
	CreatedBy		string	`json:"createdBy"`
	LastModified	string	`json:"lastModified"`
	ModifiedBy		string	`json:"modifiedBy"`
	LastUpdated		string	`json:"lastUpdated"`
	DownloadUri 	string 	`json:"downloadUri"`
	MimeType 		string	`json:"mimeType"`
	Size			string	`json:"size"`
	Checksums	struct {
		Sha1		string	`json:"sha1"`
		Md5			string	`json:"md5"`
		Sha256		string	`json:"sha256"`
	}	`json:"checksums"`
	OriginalChecksums	struct {
		Sha1		string	`json:"sha1"`
		Md5			string	`json:"md5"`
		Sha256		string	`json:"sha256"`				
	}   `json:"originalChecksums"`
	Uri 			string	`json:"uri"`
}

var request *http.Request
var err error
var foundPaths []string


func ListRepos() ([]string, error) {
	// Gets list of Repos in Artifactory instance
	var listRepos []string
	artifBase, bearer := common.AuthCreds()
	requestPath := artifBase + "/repositories"

	request, err = http.NewRequest("GET", requestPath, nil)
	request.Header.Add("Authorization", bearer)
	
	client := &http.Client{}
	response, err := client.Do(request)
	if err != nil {
		log.Println("Error on response.\n[ERROR] - ", err)
	}
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	//fmt.Println(string(body))

	// JSON return is an array of strings '[{"key":"repo_name1, "type":"LOCAL"...}, {"key":"repo_name2"}...]'
	type reposJson struct {
		Key 		string	`json:"key"`
		Description	string	`json:"description"`
		Type		string	`json:"type"`
		Url 		string	`json:"url"`
		PackageType	string	`json:"packageType"`
	}

	// Unmarshal the JSON return
	var jsonData []reposJson
	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		fmt.Printf("Could not unmarshal %s\n", err)
	}

	// As long as the results are not empty, parse the results append the repo names to 
	// the list of strings
	if len(jsonData) != 0 {
		for _, k := range jsonData {
			listRepos = append(listRepos, k.Key)
		}
		return listRepos, nil
	} else {
		err := errors.New("No repos found")
		return nil, err
	}
	return listRepos, nil
}

func GetItemChildren(item string) ([]Contents, error) {
	// Returns the children of the given item and whether that child is a folder or not (bool)
	// Item can represent a repo name or a combo of repo/child_folder/subchild_folder/etc

	// If item is the full path and filename to the artifact itself, no results will be returned as artifacts
	// do not have children. However, artifacts can be children themselves.
	artifBase, bearer := common.AuthCreds()
	requestPath := artifBase + "/storage/" + item

	type itemResults struct {
		Repo			string		`json:"repo"`
		Path			string		`json:"path"`
		Created			string		`json:"created"`
		LastModified 	string		`json:"lastModified"`
		LastUpdated		string		`json:"lastUpdated"`
		Children []struct {
			Uri		string		`json:"uri"`
			Folder	bool		`json:"folder"`
		}
		Uri				string		`json:"uri"`
		CreatedBy		string		`json:"createdBy"`	// Not exposed at repo level
		ModifiedBy		string		`json:"modifiedBy"`	// Not exposed at repo level
	}

	var childDetails []Contents

	// If the item (folder or artifact) is not empty, get the details of the item
	if (item != "") {
		request, err = http.NewRequest("GET", requestPath, nil)
		request.Header.Add("Authorization", bearer)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Println("Error on response.\n[ERROR] - ", err)
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		//fmt.Println(string(body))

		// Unmarshal the JSON return
		var jsonData *itemResults
		err = json.Unmarshal(body, &jsonData)
		if err != nil {
			fmt.Printf("Could not unmarshal %s\n", err)
		}

		// If the item has children, parse the data and return the abbreviated
		// URI ('/folder', '/folder/artifact.ext', etc) and whether the child item is a folder or not (bool)
		if len(jsonData.Children) != 0 {
			for idx, c := range jsonData.Children {
				c = jsonData.Children[idx]
				childDetails = append(childDetails, Contents{Child: c.Uri, IsFolder: c.Folder})
			}
			return childDetails, nil
		} else {
			// If no children found, we return empty contents; this isn't an error condition
			fmt.Println("No child objects found for " + item)
			return childDetails, nil
		}

		/*
		for idx := 0; idx < len(childDetails); idx++ {
			fmt.Println(childDetails[idx].Child, childDetails[idx].IsFolder)
		}*/
		// ex: [{/test-artifact-1.1.txt false} {/test-artifact-1.2.txt false} {/test-artifact-1.3.txt false}] <nil>
	} else {
		err := errors.New("No item or path provided. Unable to get child items without parent item/path.")
		return nil, err
	}

	if err != nil {
		fmt.Println("Unable to parse URL")
		return nil, err
	}

	return childDetails, nil
}


func GetArtifactPath(artifName string) ([]string, error) {
	// Takes in an artifact's name and searches Artifactory, returning the path to the artifact
	// Searches are CASE SENSITIVE
	// A path will be returned for every artifact FILE who's name includes the search string (e.g. paths for 
	// 'win2022' and 'win2022-iis' would both be returned)
	// Multiple version files for a given artifact will result in the same path being added to the list multiple times
	// So we will search for and remove duplicates before returning the results
	var childList []Contents
	var listOfPaths []string
	foundPaths = nil

	if artifName != "" {
		// Get a list of the available repos
		listRepos, err := ListRepos()
		if err != nil || listRepos[0] == "" || len(listRepos) == 0 {
			err := errors.New("No repos found")
			return nil, err
		}

		if len(listRepos) != 0 {
			for idx := 0; idx < len(listRepos); idx++ {
				childList, err = GetItemChildren(listRepos[idx])
				if len(childList) != 0 {
					listOfPaths = RecursiveSearch(childList, artifName, listRepos[idx], foundPaths)
				}
			}

			if len(listOfPaths) > 1 {
				//We'll search the list for duplicates and remove them
				listOfPaths = common.RemoveDuplicateStrings(listOfPaths)
				if len(listOfPaths) > 1 {
					fmt.Println("More than one possible artifact path found")
				}
				return listOfPaths, nil
			} else if len(listOfPaths) == 1 && listOfPaths[0] != "" {
				return listOfPaths, nil
			} else if len(listOfPaths) == 0 || listOfPaths[0] == "" {
				err := errors.New("Unable to find path to artifact")
				return nil, err
			}
		} else {
			err := errors.New("List of repos to check is empty. Either there are no repos or you do not have sufficient permissions to the repo(s).")
			return nil, err
		}
	} else {
		err := errors.New("Unable to determine path to artifact without the artifact name")
		return nil, err
	}

	// Now we have a list of repos to check through... can do another search for props...
	return listOfPaths, err
}

func RecursiveSearch(list []Contents, artifName, searchPath string, foundPaths []string) ([]string) {
	// Recursively searches a list of child items for the specificied artifact name 
	// For each child item in the list, if item isn't a folder, checks if the child item contains
	// the desired artifact name. If so, the matching item's path will be added to the foundPath list. 
	// If not, the search path will be updated to check the next layer down, and the search will run again
	var nextList []Contents
	var currentPath string
	currentPath = searchPath

	if len(list) != 0 {
		for item := 0; item < len(list); item++ {					// For each item in list...
			if list[item].IsFolder == false {						// If not a folder, does artifact match?
				// 'Contains' search is case sensitive; so we'll convert the input artifact name and convert to both cases and recheck
				lowStr := common.ConvertToLowercase(artifName)
				upStr := common.ConvertToUppercase(artifName)

				if strings.Contains(list[item].Child, artifName) {      // If we don't find it initially, we'll check with cases converted
					foundPaths = append(foundPaths, searchPath)         // If found, item's path appended to found list
				} else if strings.Contains(list[item].Child, lowStr) {
					foundPaths = append(foundPaths, searchPath)
				} else if strings.Contains(list[item].Child, upStr) {
					foundPaths = append(foundPaths, searchPath)
				}
			} else {  // IsFolder == true; so we get its children and repeat the search
				searchPath = currentPath + list[item].Child		   // 1st "/repo" + "/folder", 2nd "/repo/folder" + "/folder", etc
				nextList, err = GetItemChildren(searchPath)
				if len(nextList) != 0 {
					foundPaths = RecursiveSearch(nextList, artifName, searchPath, foundPaths)
				}
			}
		}
	}
	return foundPaths
}

func GetDownloadUri(artifPath, artifNameExt string) (string, error) {
	// Requires full path to the artifact, include full artifact name with extention
	// Gets the artifact details and will return the download URI used to retrieve the artifact
	artifBase, bearer := common.AuthCreds()
	requestPath := artifBase + "/storage" + artifPath + "/" + artifNameExt
	var downloadUri string

	// Ensures the required path components are not empty before doing a GET request
	if (artifPath != "") && (artifNameExt != "") {
		request, err = http.NewRequest("GET", requestPath, nil)
		request.Header.Add("Authorization", bearer)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Println("Error on response.\n[ERROR] - ", err)
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)

		// Unmarshal the JSON return
		var jsonData *artifJson
		err = json.Unmarshal(body, &jsonData)
		if err != nil {
			fmt.Printf("Could not unmarshal %s\n", err)
		}

		// As long as the Download URI field is not empty, parse and return the Download URI
		if jsonData.DownloadUri != "" {
			downloadUri = jsonData.DownloadUri
			return downloadUri, nil
		} else {
			err = errors.New("There is no download URI for the artifact")
			return "", err
		}
	} else {
		// If the required path components are not supplied, we'll throw an error
		message := ("Supplied artifact path: " + artifPath + " and full artifact name: " + artifNameExt)
		fmt.Println(message)
		err := errors.New("Unable to get artifact details without full path to the artifact")
		return "", err
	}

	if err != nil {
		fmt.Println("Unable to parse URL")
		return "", err
	}
	return downloadUri, nil
}

func GetCreateDate(artifactUri string) (string, error) {
	// Requires full path to the artifact, include full artifact name with extention
	// Gets the artifact details and will return the string date created
	_, bearer := common.AuthCreds()
	var createdDate string

	// Ensures the required path components are not empty before doing a GET request
	if (artifactUri != "") {
		request, err = http.NewRequest("GET", artifactUri, nil)
		request.Header.Add("Authorization", bearer)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Println("Error on response.\n[ERROR] - ", err)
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)

		// Unmarshal the JSON return
		var jsonData *artifJson
		err = json.Unmarshal(body, &jsonData)
		if err != nil {
			fmt.Printf("Could not unmarshal %s\n", err)
		}

		// As long as the Created field is not empty, parse and return the Create Date
		if jsonData.Created != "" {
			createdDate = jsonData.Created
			return createdDate, nil
		} else {
			err = errors.New("There is no create date for the artifact")
			return "", err
		}
	} else {
		// If the required path components are not supplied, we'll throw an error
		message := ("Supplied artifact path: " + artifactUri)
		fmt.Println(message)
		err := errors.New("Unable to get artifact details without full path to the artifact")
		return "", err
	}

	if err != nil {
		fmt.Println("Unable to parse URL")
		return "", err
	}
	return createdDate, nil
}

func RetrieveArtifact(downloadUri string) (string, error) {
	// Gets the artifact via provided Download URI and copies it to the output directory specified in
	// the environment variables file
	var outputDir string
	_, bearer := common.AuthCreds()

	// If no output directory path was provided, the artifact file will be downloaded to the top-level
	// directory of this code
	if len(os.Getenv("OUTPUTDIR")) != 0 {
		OUTPUTDIR := os.Getenv("OUTPUTDIR")
		outputDir = common.EscapeSpecialChars(OUTPUTDIR)  // Ensure special characters are escaped
		outputDir = common.CheckAddSlashToPath(outputDir) // Ensure path ends with appropriate slash type
	} else {  // There's no OUTPUTDIR env var...
		fmt.Println("No output directory provided; output will be at top-level directory")
		outputDir = ""
	}

	// If we have a download URI, get the artifact and download it
	if downloadUri != "" {
		request, err = http.NewRequest("GET", downloadUri, nil)
		request.Header.Add("Authorization", bearer)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Println("Error on response.\n[ERROR] - ", err)
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		//fmt.Println(string(body))   // prints the contents of the file

		if response.StatusCode == 404 {
			err := errors.New("File not found.")
			return "File download failed.", err
		} else {  // File was found
			// Create file name from download URI path of artifact
			fileUrl, err := url.Parse(downloadUri)
			if err != nil {
				log.Fatal(err)
			}

			// Get the file name from the path
			path := fileUrl.Path
			segments := strings.Split(path, "/")
			fileName := segments[len(segments)-1]

			// Creates the file at the defined path
			// Will overwrite the file if it already exists
			newFile, err := os.Create(outputDir + fileName)   
			if err != nil {
				log.Fatal(err)
				return "Error creating file at target location.", err
			}
			err = os.WriteFile(outputDir + fileName, body, 0777)  //set this to something else...
			if err != nil {
				log.Fatal(err)
				return "Error downloading file to target location.", err
			}
			defer newFile.Close()
		}
	} else {
		err := errors.New("No download URI was provided. Unable to download the artifact without the download URI.")
		return "File download failed.", err
	}

	return "Completed file download", nil
}

func UploadFile(sourcePath, targetPath, fileSuffix string) (string, error) {
	// sourcePath includes full file path; needs proper escape chars (ex: h:\\lab\\artifact.txt OR /lab/artifact.txt)
	// targetPath will be /repo-key/folder/path/
	// The target FILENAME will match the source filename as it exists in the source directory
	// fileSuffix is a placeholder for potential distinguishing values such as dates, versions, etc. 
		// If "", it will be ignored. Otherwise, it will be appended to target FILENAME.

	var downloadUri string
	var filePath string
	var fileName string
	var found bool
	artifBase, bearer := common.AuthCreds()
	separater := "-"										  // If adding a file suffix (like date, version, etc), use this separater between filename and suffix
	trimmedBase := artifBase[:len(artifBase)-4]               // Removing '/api' from base URI

	if len(sourcePath) != 0 && targetPath != "" { 
		// We need to ensure the provided source path/file are valid and exist
		if len(path.Ext(sourcePath)) != 0 {		                  // Ensures file with extension exists in source path
			sourcePath = common.EscapeSpecialChars(sourcePath)
			targetPath = common.EscapeSpecialChars(targetPath)
			targetPath = common.CheckAddSlashToPath(targetPath)
			
			// Determine source filename and source file path by platform type
			winPath := common.CheckPathType(sourcePath)
			if winPath == true {
				segments := strings.Split(sourcePath, "\\")	  	  // Split source path into segments
				fileName = segments[len(segments)-1]			  // Determine filename from path
				filePath = sourcePath[:len(sourcePath)-len(fileName)]  // Determine path without filename
			} else {   // Unix path
				segments := strings.Split(sourcePath, "/")	 	  // Split source path into segments
				fileName = segments[len(segments)-1]              // Determine filename from path
				filePath = sourcePath[:len(sourcePath)-len(fileName)]  // Determine path without filename				
			}
			
			// Get all files in the provided source directory
			filesInDirectory, err := os.ReadDir(filePath)
			if err != nil {
				return "", err
			}
			
			// For each file in the source directory, do a case insensitive file name comparison for a match
			// As Artifactory cares about case here, we want to make sure the filename supplied matches the case of the filename that actually exists in the source path
			for _, file := range filesInDirectory {
				isSameStr := common.StringCompare(fileName, file.Name())            // Filename from provided source path vs. filename pulled directly from source path
				if isSameStr == true {												// If true, we know files are the same
					found = true													// Mark that we found a matching file
					isExactStr, err := common.SearchForExactString(file.Name(), fileName)  // Now, checks if cases matches
					if err != nil {
						fmt.Println("Error searching for exact string: ", err)
					}
					
					if isExactStr == false {										// Files are the same, but provided and actual cases are different
						fileName = file.Name() 										// Set the provided filename to match the actual filename so we'll use to the correct case
					}
				}
			}

			// If we couldn't find a matching file at all, then we throw an error
			if found == false {
				err := errors.New("Unable to validate existance of source file. Source file doesn't exist.")
				return "", err
			}
			
			// We now have a validated source path and filename
			// Set target filename = filename + fileSuffix (if not blank)
			if len(fileSuffix) != 0 || fileSuffix != "" {							// If a file suffix (like version, date, etc) was provided...
				fileExt := path.Ext(fileName)										// Returns .[ext]
				justName := strings.Trim(fileName, fileExt)							// Trim off extension
				fileName = justName + separater + fileSuffix + fileExt
			}   // If blank, then the original filename will be used
			
			// Now we're ready to form our request inputs and make the API call
			newArtifactPath := trimmedBase + targetPath + fileName                  // Forms: http://artifactory_base_api_url/repo-key/folder/artifact.txt
			data := strings.NewReader("@/" + sourcePath)                            // Formats the payload appropriately
			
			// Makes the API call to upload the specified file
			request, err = http.NewRequest("PUT", newArtifactPath, data)
			request.Header.Add("Authorization", bearer)
	
			client := &http.Client{}
			response, err := client.Do(request)
			if err != nil {
				log.Println("Error on response.\n[ERROR] - ", err)
			}
			defer response.Body.Close()
			body, err := io.ReadAll(response.Body)
			// Prints the response which includes the details about the new artifact
			fmt.Println(string(body))
	
			// Unmarshals the JSON data
			var jsonData *artifJson
			err = json.Unmarshal(body, &jsonData)
			if err != nil {
				fmt.Printf("Could not unmarshal %s\n", err)
			}
	
			// As long as the Download URI is populated, we'll parse and return the Download URI
			if jsonData.DownloadUri != "" {
				downloadUri = jsonData.DownloadUri
				return downloadUri, nil
			} else {  // Otherwise, we will throw an error
				err = errors.New("There is no download URI for the artifact")
				return "", err
			}
		} else {  // No file extension found
			err = errors.New("No file extension found in source path. Ensure source includes path and source file with extension.")
			return "", err
		}
	} else {
		// If the required components are not supplied, we will throw an error
		message := ("Supplied source path: " + sourcePath + ", target path: " + targetPath)
		err := errors.New("Cannot upload file without source path/file, target path, and artifact file name")
		fmt.Println(message)
		return "", err
	}
	if err != nil {
		fmt.Println("Unable to parse URL")
		return "", err
	}
	return downloadUri, nil
}

func DeleteArtifact(artifUri string) (string, error) {
    // Takes in artifact's URI and executes a DELETE call against it
	_, bearer := common.AuthCreds()

	if artifUri != "" { 
		request, err = http.NewRequest("DELETE", artifUri, nil)
		request.Header.Add("Authorization", bearer)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Println("Error on response.\n[ERROR] - ", err)
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		fmt.Println(string(body))
		
		// If the request is successful, it will simply return a status code of 204
		if response.StatusCode == 204 {
			fmt.Println("Request completed successfully")
			statusCode = "204"
		} else {
			// If the request fails, it will return a status code of 400
			fmt.Println("Unable to complete request")
			statusCode = "404"
		}
	} else {
		// If the required component is not supplied, we will throw an error
		message := ("Supplied artifact path is: " + artifUri)
		fmt.Println(message)
		err := errors.New("Unable to DELETE item without artifact URI.")
		return "", err
	}

	if err != nil {
		fmt.Println("Unable to parse URL")
		return "", err
	}

	return statusCode, nil
}
