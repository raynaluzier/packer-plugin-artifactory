package operations

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"packer-plugin-artifactory/common"
	"sort"
	"strings"
)

//var request *http.Request - declared in general.go
//var err error - declared in general.go
var statusCode string

// Used to store and structure the returned key/values that can be anything
type prop struct {
	Name string
	Value string
}

func GetArtifactPropVals(artifUri string, listPropKeys []string) (interface{}, error){
	// Takes in the URI of the artifact, plus one or more property keys
	// Returns the values for only the properties included in the URI for the given artifact
	// Search is CASE SENSTIVE
	_, bearer := common.AuthCreds()
	var properties []prop
	
	// Ensures the components of the required full artifact path along with at least one or more properties is not empty
	if artifUri != "" {
		// Determines whether we will format a list of properties first, or pass a single property
		// before making the API call
		if len(listPropKeys) > 1 {
			// If there's more than one property name supplied, adds the required ',' separater between them
			strProps := strings.Join(listPropKeys, ",")
			request, err = http.NewRequest("GET", artifUri + "?properties=" + strProps, nil)
		} else if len(listPropKeys) == 1 && listPropKeys[0] != "" {
			request, err = http.NewRequest("GET", artifUri + "?properties=" + listPropKeys[0], nil)
		} else {
			err := errors.New("Unable to search for Artifact properties without one or more property names")
			return nil, err
		}
		request.Header.Add("Authorization", bearer)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Println("Error on response.\n[ERROR] - ", err)
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)
		//fmt.Println(string(body))

		if err != nil || response.StatusCode == 404 {
			err := errors.New("No matching property(ies) could be found")
			return nil, err
		} else {  // There was no error
			// Declares a map whose key type is a string with any value type
			// This is used because the returned JSON data is unstructured; 'properties' contains one or more key/values that
			// correspond to a property name and property value that can be anything
			var result map[string]any

			// Unmarshal the JSON return
			err = json.Unmarshal(body, &result)
			if err != nil {
				fmt.Printf("Could not unmarshal %s\n", err)
			}

			// As long as the property results are not empty, parse thru the property keys and values
			// The property keys are returned as a string, but the values must be converted to string first
			// and the surrounding [ ] brackets are trimmed off
			// Each key/value pair are stored in a struct of type 'prop' and returned, allowing for easier parsing later
			var strValue string

			parseProps := result["properties"].(map[string]any)
			if len(parseProps) != 0 {
				for k, v := range parseProps {
					strValue = fmt.Sprintf("%v", v)
					strValue = strings.Trim(strValue, "]")
					strValue = strings.Trim(strValue, "[")
					//fmt.Println(k, strValue)
					properties = append(properties, prop{Name: k, Value: strValue})
				}
				return properties, nil
				/*for idx := 0; idx < len(properties); idx++ {
					fmt.Println(properties[idx].Name, properties[idx].Value)
				}*/
			} else {
				err := errors.New("No results returned")
				return nil, err
			}
		}

	} else {
		// If no properties were supplied, we'll throw an error
		message := ("Artifact URI is: " + artifUri)
		if len(listPropKeys) != 0 && listPropKeys[0] != "" {
			// If the artifact URI wasn't provided, we'll throw an error
			fmt.Println(message)
			err := errors.New("Unable to search for Artifact properties without the artifact's URI")
			return nil, err
		} else {
			// If at least one property wasn't provided, we'll throw an error
			fmt.Println(message)
			err := errors.New("Unable to search for Artifact properties without the artifact's URI and one or more property names")
			return nil, err
		}
	}

	if err != nil {
		fmt.Println("Unable to parse URL")
		return nil, err
	}

	return properties, nil
	
}

func GetAllPropsForArtifact(artifUri string) (interface{}, error) {
	// Takes in the URI of a given artifact and pulls all the properties and their values assigned to the artifact
	_, bearer := common.AuthCreds()
	var properties [] prop

	if artifUri != "" {
		request, err = http.NewRequest("GET", artifUri + "?properties", nil)
		request.Header.Add("Authorization", bearer)

		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Println("Error on response.\n[ERROR] - ", err)
		}
		defer response.Body.Close()
		body, err := io.ReadAll(response.Body)

		if err != nil || response.StatusCode == 404 {
			err := errors.New("No matching property(ies) could be found")
			return nil, err
		} else {
			// Declares a map whose key type is a string with any value type
			// This is used because the returned JSON data is unstructured; 'properties' contains one or more key/values that
			// correspond to a property name and property value that can be anything
			var result map[string]any

			// Unmarshal the JSON return
			err = json.Unmarshal(body, &result)
			if err != nil {
				fmt.Printf("Could not unmarshal %s\n", err)
			}

			// As long as the property results are not empty, parse thru the property keys and values
			// The property keys are returned as a string, but the values must be converted to string first
			// and the surrounding [ ] brackets are trimmed off
			// Each key/value pair are stored in a struct of type 'prop' and returned, allowing for easier parsing later
			var strValue string

			parseProps := result["properties"].(map[string]any)
			if len(parseProps) != 0 {
				for k, v := range parseProps {
					strValue = fmt.Sprintf("%v", v)
					strValue = strings.Trim(strValue, "]")
					strValue = strings.Trim(strValue, "[")
					//fmt.Println(k, strValue)
					properties = append(properties, prop{Name: k, Value: strValue})
				}
				return properties, nil
				/*for idx := 0; idx < len(properties); idx++ {
					fmt.Println(properties[idx].Name, properties[idx].Value)
				}*/
			} else {
				err := errors.New("No results returned")
				return nil, err
			}
		}
	} else {
		message := ("Artifact URI is: " + artifUri)
		fmt.Println(message)
		err := errors.New("Unable to retrieve properties of the artifact without the Artifact's URI.")
		return nil, err
	}
}

func FilterListByProps(listArtifUris, listKvProps []string) (string, error) {
	// Takes in a list of artifact URIs; for each URI, pulls its list of properties
	// Then compares the list of one or more key/value pairs ('key=value') provided to the key/values assigned to the artifact
	// If there's a match, the artifact URI will be added to the filteredList
	// The same URI may be added to the list for multiple property matches, 
	
	var foundList []string
	var filteredList []string
	var structData []map[string]interface{}
	numProps := len(listKvProps)
	var dateMap []map[string]string
	var foundItem string

	if len(listArtifUris) != 0 && len(listKvProps) != 0 {
		for a := 0; a < len(listArtifUris); a++ {
			// For each artifact URI in list, get it's properties/values; there can be one or more properties/values assigned
			artifProps, err := GetAllPropsForArtifact(listArtifUris[a])  // ex return: [{release stable} {testing passed}]
			if err != nil {
				//log.Println("No properties returned for artifact: " + listArtifUris[a])
			} else {
				// Convert custom data type 'prop' object passed out as interface{} into JSON format
				jsonBytes, err := json.Marshal(artifProps)
				if err != nil {
					fmt.Println("Unable to marshal data - ", err)
				}
				// Convert the JSON data into a map of arbitrary values to support any type (in this case, our custom 'prop' type)
				err = json.Unmarshal([]byte(jsonBytes), &structData)
				if err != nil {
					fmt.Println("Unable to unmarshal data - ", err)
				}

				// For each returned key/value property assigned to the artifact...
				for idx := 0; idx < len(structData); idx++ {
					// Convert each pair to a string and format to match the listKvProps input ('key=value')
					propName := fmt.Sprintf("%v", structData[idx]["Name"]) 
					propVal := fmt.Sprintf("%v", structData[idx]["Value"])
					propCompare := propName + "=" + propVal

					// For each key/value property provided as input, compare it to what was returned from the artifact
					// If the property key/value pair matches, add the artifact to the foundList
					// If more than one property key/value pair was input as a filter, an instance of the matching artifact 
					// will be added added to the foundList; we'll handle this situation below
					// So inputting 3 property key/value pairs will result in 3 instances of the artifact in the foundList
					for k := 0; k < len(listKvProps); k++ {
						if propCompare == listKvProps[k] {
							fmt.Println("Property found: " + listArtifUris[a])
							foundList = append(foundList, listArtifUris[a])
						}
					}
				}
				
			}
		}

		// For ex: If 3 property key/value pairs were input as filters, we'd expect that any artifact that has ALL of those matching
		// properties is more likely the artifact we're looking for. However, it's possible that multiple artifacts have those same property
		// key/value pairs ('release=stable', 'testing=passed') for a given artifact (new 'win-22' image built over multiple days).
		
		// If only one artifact is present in the foundList, we'll return this as the artifact. 
		// If multiple artifacts (ex. 'win-22') with the same exact properties are returned, then we'll return the latest based on 'created' date

		// If more than one artifact is returned in the foundList...
		if len(foundList) > 1 {
			// Count the occurance of duplicate artifacts and return a map of the artifact and duplicate count
			countMap := common.ReturnWithDupCounts(foundList)
			fmt.Println(countMap)
			for str, count := range countMap {
				// If the number of duplicate artifacts found matches the number of input property key/value pairs, add them to a filter list
				if count == numProps {
					filteredList = append(filteredList, str)
				}
			}

			// For each artifact in the filter list, we grab it's 'created' date and add that artifact and date to an array of maps
			for i := 0; i < len(filteredList); i++ {
				addMap := make(map[string]string)
				created, err := GetCreateDate(filteredList[i])
				if err != nil {
					fmt.Println("Error getting created date")
				}
				//fmt.Println(created)

				addMap["artifact"] = filteredList[i]
				addMap["created"] = created
				dateMap = append(dateMap, addMap)

				//fmt.Println(dateMap[i]["artifact"], dateMap[i]["created"])
			}
			// Sort by 'created' date and return the latest instance of the artifact
			sort.Slice(dateMap, func(i, j int) bool { 
				return dateMap[i]["created"] < dateMap[j]["created"]
			})
			latest := len(dateMap) - 1
			foundItem = dateMap[latest]["artifact"]
			return foundItem, nil

		  // If only one matching artifact was found, return that artifact
		} else if len(foundList) == 1 {
			foundItem = foundList[0]
			return foundItem, nil
		} else {   // If no artifacts were founds...
			err := errors.New("No matching artifacts were found.")
			return "", err
		}
	}
	return foundItem, nil
}

func SetArtifactProps(artifUri string, listKvProps []string) (string, error) {
	// Takes in list of property key/value pairs and assigns them to the provided artifact
	// Inputs are CASE SENSITIVE
	// Requires artifact's URI, plus one or more property keys/values
	// Disallows special characters: )( }{ ][ *+^$\/~`!@#%&<>;, and the SPACE character
	_, bearer := common.AuthCreds()
	requestPath := artifUri + "?properties="

	//Checks for special characters in prop key/values; throws error if found
	if common.ContainsSpecialChars(listKvProps) == true {
		err := errors.New("Properties cannot contain special characters --> )( }{ ][ *+^$\\/~`!@#%&<>;, and SPACE")
		fmt.Println("Special character found")
		return "", err
	} else {
		// Ensures the artifact URI along with at least one or more property keys/values is not empty
		if artifUri != "" && len(listKvProps) != 0 {
			// Determines whether we will format a list of property keys/values first, or pass a single property key/value pair
			// before making the API call
			if len(listKvProps) > 1 {
				// If there's more than one property keys/values supplied, adds the required ';' separater between them
				strProps := strings.Join(listKvProps, ";")
				//fmt.Println(strProps)
				request, err = http.NewRequest("PUT", requestPath + strProps, nil)
			} else if len(listKvProps) == 1 && listKvProps[0] != "" {
				request, err = http.NewRequest("PUT", requestPath + listKvProps[0], nil)
			} else {
				err := errors.New("Unable to set Artifact properties without one or more property names and values")
				return "", err
			}
			request.Header.Add("Authorization", bearer)
			
			client := &http.Client{}
			response, err := client.Do(request)
			if err != nil {
				log.Println("Error on response.\n[ERROR] - ", err)
			}
			defer response.Body.Close()

			// If the request is successful, it will simply return a status code of 204
			if response.StatusCode == 204 {
				fmt.Println("Request completed successfully")
				statusCode = "204"
			} else {
				// If the request fails, it will return a status code of 400
				fmt.Println("Unable to complete request")
				statusCode = "400"
			}
		} else {
			// If no property keys/values were supplied, we'll throw an error
			numProps := len(listKvProps)
			if numProps != 0 {
				// If the artifact URI wasn't provided, we'll throw an error
				err := errors.New("Unable to set Artifact properties without artifact's URI")
				return "", err
			} else {
				// If at lease one property key/value weren't provided, we'll throw an error
				err := errors.New("Unable to set Artifact properties without artifact's URI and one or more property names")
				return "", err
			}
		}
	}

	if err != nil {
		fmt.Println("Unable to parse URL")
		return "", err
	}

	return statusCode, nil
}

func DeleteArtifactProps(artifUri string, listProps []string) (string, error) {
	// Takes in the artifact's URI and one or more property keys and removes the corresponding property from the provided artifact
	// If a property is provided that doesn't exist (which includes incorrectly cased properties), the API ignores this and will return a successful response
	_, bearer := common.AuthCreds()
	requestPath := artifUri + "?properties="

	// Ensures the artifact URI along with at least one or more property keys is not empty
	if artifUri != "" && len(listProps) != 0 {
		// Determines whether we will format a list of property keys first, or pass a single property key
		// before making the API call
		if len(listProps) > 1 {
			// If there's more than one property keys supplied, adds the required ',' separater between them
			strProps := strings.Join(listProps, ",")
			//fmt.Println(strProps)
			request, err = http.NewRequest("DELETE", requestPath + strProps, nil)
		} else if len(listProps) == 1 && listProps[0] != "" {
			request, err = http.NewRequest("DELETE", requestPath + listProps[0], nil)
		} else {
			err := errors.New("Unable to delete Artifact properties without one or more property names")
			return "", err
		}
		request.Header.Add("Authorization", bearer)
		
		client := &http.Client{}
		response, err := client.Do(request)
		if err != nil {
			log.Println("Error on response.\n[ERROR] - ", err)
		}
		defer response.Body.Close()

		// If the request is successful, it will simply return a status code of 204
		if response.StatusCode == 204 {
			fmt.Println("Request completed successfully")
			statusCode = "204"
		} else {
			// If the request fails, it will return a status code of 400
			fmt.Println("Unable to complete request")
			statusCode = "400"
		}
	} else {
		// If no property keys were supplied, we'll throw an error
		numProps := len(listProps)
		message := ("Artifact Path is: " + artifUri)
		if numProps != 0 {
			// If the artifact URI was not provided, we'll throw an error
			fmt.Println(message)
			err := errors.New("Unable to delete Artifact properties without artifact URI")
			return "", err
		} else {
			// If at lease one property key wasn't provided, we'll throw an error
			fmt.Println(message)
			err := errors.New("Unable to delete Artifact properties without artifact URI and one or more property names")
			return "", err
		}
	}

	if err != nil {
		fmt.Println("Unable to parse URL")
		return "", err
	}

	return statusCode, nil
}

