package functions

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"context"
	"strings"

	"groupXX/firebase"
	"groupXX/structures"
)

//used to check if we should use country name or country code
func IsCountryCode(countryInput string) bool {
	if len(countryInput) == 3 {
		return true
	} else {
		return false
	}
}

//functions to retrieve the specified country info
func ReadCountryInfo(w http.ResponseWriter, searchInput string, current bool, begin *int, end *int) ([]structures.DataEntry, error) {
	//if not specified search input just read the file completly because for no specified country 
	//this becomes more effecient, if current only write for current year, else all
	if searchInput == ""{
		if current == true{
			//a map presaved with only current countries
			return OnlyCurrent, nil
		} else{
			allCountries, err := RetrieveAll(structures.FILEPATH, false)
			if err != nil{
				log.Printf("Error retrieving countries from file: %v", err)
				http.Error(w, "Error retrieving countries from file: "+err.Error(), http.StatusInternalServerError)
			}
			return allCountries, nil
		}
	}

	//generate cache key
	cacheKey := fmt.Sprintf("%s_%v_%v_%v", searchInput, current, begin, end)
	
	ctx := context.Background()

	//Try getting data from cache
	cachedData, err := firebase.GetCachedData(ctx, cacheKey)
	if err != nil {
		log.Printf("Error getting cached data: %v", err)
	} else if cachedData != nil {
		return cachedData, nil
	}

	//call ExtractByMap to get matching countries
	matchingCountries, err := ExtractByMap(searchInput)
	if err != nil {
		return nil, err
	} else if matchingCountries == nil{
		return nil, nil
	}

	//filter the matching countries based on the additional criteria
	var data []structures.DataEntry
	//goes trough all matching countries
	for _, entry := range matchingCountries {
		//if the user want current year which as of our data is in 2021
		if current {
			if entry.Year != 2021 {
				continue
			}
		}

		//specifications for begin and end, where it can have both or either
		if begin != nil && end != nil {
			if entry.Year < *begin || entry.Year > *end {
				continue
			}
		} else if begin != nil {
			if entry.Year < *begin {
				continue
			}
		} else if end != nil {
			if entry.Year > *end {
				continue
			}
		}

		//appends the data that have "passed" all checks and haven't been "continue'd"
		data = append(data, entry)
	}

	//if the data wasn't found in cache, cache the data
	if cachedData == nil{
		err := firebase.SetCachedData(ctx, cacheKey, data)
		if err != nil {
			log.Printf("Error setting cached data: %v", err)
		}
	}
	return data, nil
}

//function to retrieve a countries neighbour
func RetrieveNeighbours(w http.ResponseWriter, searchCountry string) ([]string, error) {
	//sets the contet type to JSON format
	w.Header().Add("content-type", "application/json")
	//string list which stroes country names
	var borderNames []string

	//get's main countries
	countryBorders, err := GetCountryData(w, searchCountry,  structures.COUNTRYSEARCH)
	if err != nil {
		//if the country wasn't found don't report error because a country isn't obligated to have neighbours
		//so just return empty list
		if err.Error() == "Country not found" {
			return []string{}, nil
		}
		log.Printf("Error getting country API: %v", err)
		http.Error(w, "Error getting country API: "+err.Error(), http.StatusInternalServerError)
		return nil, err
	}

	//loops trough each country retrieved from the restcountries api
	for _, country := range countryBorders {
		//loops trough that counties border countries
		for _, borderCountry := range country.Borders {
			//gets the data for that border country
			borderCountries, err := GetCountryData(w, borderCountry, structures.COUNTRYSEARCH)
			if err != nil {
				if err.Error() == "Country not found" {
					continue
				}
				log.Printf("Error getting country API: %v", err)
				http.Error(w, "Error getting country API: "+err.Error(), http.StatusInternalServerError)
				return nil, err
			}
			//appends the found names of the border countries to the list of border names
			borderNames = append(borderNames, borderCountries[0].Name.Common)
		}
	}

	return borderNames, nil
}

//TEST DENNE TA URL SOM PARAMETER
//function to get country data, returns list of the country struct
func GetCountryData(w http.ResponseWriter, country string, path string) ([]structures.Country, error) {
    var responseCountryBody []byte
    var err error

    if path == structures.TESTCOUNTRYFILE {
        // Read the test country file
        responseCountryBody, err = ioutil.ReadFile(path)
        if err != nil {
            log.Printf("Error reading test country file: %v", err)
            http.Error(w, "Error reading test country file: "+err.Error(), http.StatusInternalServerError)
            return nil, err
        }
    } else {
        var responseCountry *http.Response
		//
		if IsCountryCode(path){
        	responseCountry, err = http.Get("http://129.241.150.113:8080/v3.1/alpha/" + url.PathEscape(country))
		} else {
			responseCountry, err = http.Get("http://129.241.150.113:8080/v3.1/name/" + url.PathEscape(country))
		}

        if err != nil {
            log.Printf("Error getting country API: %v for %s", err, country)
            http.Error(w, "Error getting country API: "+err.Error(), http.StatusInternalServerError)
            return nil, err
        }

        defer responseCountry.Body.Close()

        if responseCountry.StatusCode != http.StatusOK {
            if responseCountry.StatusCode == http.StatusNotFound {
                return []structures.Country{}, fmt.Errorf("Country not found")
            }
            err := fmt.Errorf("Country API returned a non-200 status code: %v", responseCountry.StatusCode)
            log.Printf("Error getting country API: %v", err)
            http.Error(w, err.Error(), http.StatusInternalServerError)
            return nil, err
        }

        responseCountryBody, err = ioutil.ReadAll(responseCountry.Body)
        if err != nil {
            log.Printf("Error reading country API response: %v", err)
            http.Error(w, "Error reading country API response: "+err.Error(), http.StatusInternalServerError)
            return nil, err
        }
    }

    var countries []structures.Country
    err = json.Unmarshal(responseCountryBody, &countries)
    if err != nil {
        log.Printf("Error decoding country API response: %v", err)
        http.Error(w, "Error decoding country API response: "+err.Error(), http.StatusInternalServerError)
        return nil, err
    }

	if path == structures.TESTCOUNTRYFILE {
        matchingCountries := make([]structures.Country, 0)
        for _, c := range countries {
            if strings.EqualFold(c.Name.Common, country) {
                matchingCountries = append(matchingCountries, c)
            }
        }
        return matchingCountries, nil
    }

    return countries, nil
}

// function which outputs data in a desired format
func PrintData(w http.ResponseWriter, data interface{}) {
	encoded, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error encoding API response: %v", err)
		http.Error(w, "Error encoding API response: "+err.Error(), http.StatusInternalServerError)
	}
	fmt.Fprintln(w, string(encoded))
}

