package common

import (
	"strings"
	"encoding/json"
	"log"
	"regexp"
	"fmt"
	"os"
)

func AuthCreds() (string, string) {
	artifServer := os.Getenv("ARTIF")
	token := os.Getenv("TOKEN")
	bearer := "Bearer " + token
	return artifServer, bearer
}

func ConvertToLowercase(inputStr string) string {
	// Converts string to lowercase
	lowerStr := strings.ToLower(inputStr)
	return lowerStr
}

func ConvertToUppercase(inputStr string) string {
	// Converts string to uppercase
	upperStr := strings.ToUpper(inputStr)
	return upperStr
}

func RemoveDuplicateStrings(listOfStrings []string) ([]string) {
	// Searches list of strings and removes duplicates
	allStrings := make(map[string]bool)

	list := []string{}
	for _, item := range listOfStrings {
		if _, value := allStrings[item]; !value {
			allStrings[item] = true
			list = append(list, item)
		}
	}
	return list
}

func ReturnWithDupCounts(listOfStrings []string) (map[string]int) {
	// Count occurances of each string and returns map of strings and their duplicate counts
	countMap := make(map[string]int)
	
	for _, str := range listOfStrings {
		countMap[str]++
	}
	return countMap
}

func ReturnDuplicates(countMap map[string]int) []string {
	// Takes in a count map of strings and their number of duplicate occurances (map[str1:1, str2:5, str3:1])
	// For any strings with more than one occurance, the string is added to the duplicates lists and returned
	duplicates := []string{}

	for str, count := range countMap {
		if count > 1 {
			duplicates = append(duplicates, str)
		}
	}
	return duplicates
}

func SetArtifUriFromDownloadUri(downloadUri string) string {
	downloadUri = strings.Replace(downloadUri, "8082", "8081", 1)  // Modify the server port from 8082 to 8081
	artifServer := os.Getenv("ARTIF")                              // http://server.com:8081/artifactory/api
	trimmedServer := strings.TrimSuffix(artifServer, "/api")	   // http://server.com:8081/artifactory
	artifSuffix := strings.TrimPrefix(downloadUri, trimmedServer)  // /repo-key/folder/path/artifact.ext
	artifUri := artifServer + "/storage" + artifSuffix      // http://server.com:8081/artifactory/storage/repo-key/folder/path/artifact.ext
	
	return artifUri
}

func SearchForExactString(searchTerm, inputStr string) (bool, error) {
	// Searches an input string for an exact search term
	// For example: "win2022" will return true if input string is "win2022", false if "win2022-iis"
	result, err := regexp.MatchString("(?sm)^" + searchTerm + "$", inputStr)
	if err != nil {
		fmt.Println("Error searching for : " + searchTerm)
		return result, err
	}
	return result, err
}

func EscapeSpecialChars(input string) (string) {
	// Takes the output directory provided from the environment variable and adds escape characters
	// For Ex: F:\mypath\ becomes F:\\mypath\\
	var js json.RawMessage
	// Replace newlines with space rather than escaping them
	input = strings.ReplaceAll(input, "\n", " ")
	// Done to take the help of the json.Unmarshal function
	jsonString := createJsonString(input)
	byteValue := []byte(jsonString)
	err := json.Unmarshal(byteValue, &js)

	// Escape spechail characters only if JSON unmarshal results in an error
	if err != nil {
		out, err := json.Marshal(input)
		if err != nil {
			log.Printf("json marshalling failed with error : %v", err)
			return input
		} else {
			// JSON marshal quotes the entire string which results in double quotes at beginning/end of string
			return string(out[1 : len(out)-1])
		}
	}
	return input
}

func createJsonString(input string) string {
	// Used with SpecialChars function to properly format output directories that may include "\" in path
	jsonString := "{\"key\":\""
	endJson := "\"}"
	jsonString = jsonString + input + endJson
	return jsonString
}

func CheckPathType(path string) bool {
	// Checks path to see if path is Unix-based (has '/') or Windows-based (has '\')
	isWinPath := strings.Contains(path, "\\")
	return isWinPath
}

func StringCompare(inputStr, actualStr string) bool {
	// Performs case INSENSITIVE comparision of strings (like file names); returns true if they match
	// Does NOT do partial string comparisons; "win" and "win-2022" will be false
	if strings.EqualFold(inputStr, actualStr) {
		return true
	} else {   // Different strings
		return false
	}
}

func CheckAddSlashToPath(path string) string {
	// Based on path type (Win vs Unix), checks path to see if it ends with appropriate back or forward slash, if not, will add as appropriate
	// This is to ensure the output directory path provided is formatted as required
	lastChar := path[len(path)-1:]
	winPath := CheckPathType(path)

	if winPath == true {
		if lastChar == "\\" {
			fmt.Println("Path: '" + path + "' is formatted properly")
			return path
		} else {
			// Add backslash to path
			path = path + "\\"
			return path
		}
	} else {  // Unix Path
		if lastChar == "/" {
			fmt.Println("Path: '" + path + "' is formatted properly")
			return path
		} else {
			// Add forwardslash to path
			path = path + "/"
			return path
		}
	}
	return path
}

func ContainsSpecialChars(strings []string) bool {
    // Checks for the special characters disallowed by Artifactory in Properties
	// Returns true if ANY of the chars are found; false if not
	pattern := regexp.MustCompile("[(){}\\[\\]*+^$\\/~`!@#%&<>;, ]")  // add '=' back later
	for idx := 0; idx< len(strings); idx++ {
		if pattern.MatchString(strings[idx]) {
			return true
		}
	}
	return false
}
