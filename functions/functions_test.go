package functions

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
	"encoding/json"
	"io"
	"os"
	"sort"

	"groupXX/structures"
)

func TestIsCountryCode(t *testing.T) {
	rand.Seed(time.Now().UnixNano())

	//generates random strings of different lengths
	for i := 0; i < 100; i++ {
		length := rand.Intn(6) // Generate a random length between 0 and 5
		randomStr := randomString(length)

		//checks if the IsCountryCode function returns the correct output
		result := IsCountryCode(randomStr)
		if length == 3 && !result {
			t.Errorf("Expected true for %s, but got false", randomStr)
		} else if length != 3 && result {
			t.Errorf("Expected false for %s, but got true", randomStr)
		}
	}
}

// function to create random string
func randomString(n int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	str := make([]byte, n)
	for i := range str {
		str[i] = letters[rand.Intn(len(letters))]
	}
	return string(str)
}

func TestRetrieveNeighbours(t *testing.T) {
	testCases := []struct {
		countryName        string
		expectedNeighbours []string
	}{
		{
			// Based on country name expects the following neighbours
			countryName:        "Netherlands",
			expectedNeighbours: []string{"Belgium", "Germany"},
		},
		{
			countryName:        "Sri Lanka",
			expectedNeighbours: []string{"India"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.countryName, func(t *testing.T) {
			rr := httptest.NewRecorder()

			neighbours, err := RetrieveNeighbours(rr, tc.countryName)
			if err != nil {
				t.Errorf("RetrieveNeighbours failed: %v", err)
				return
			}


			fmt.Printf("%s neighbours: %v\n", tc.countryName, neighbours) // Add this line to print the actual neighbors

			if len(neighbours) != len(tc.expectedNeighbours) {
				t.Errorf("Expected %d neighbours, got %d", len(tc.expectedNeighbours), len(neighbours))
				return
			}

			sort.Strings(neighbours)
			sort.Strings(tc.expectedNeighbours)

			for i, expectedNeighbour := range tc.expectedNeighbours {
				if expectedNeighbour != neighbours[i] {
					t.Errorf("Expected neighbour %s, got %s", expectedNeighbour, neighbours[i])
				}
			}
		})
	}
}


func TestGetCountryDataFromFile(t *testing.T) {
	// List of test cases with random country names and sample data
	testCases := []struct {
		countryName string
		sampleData  structures.Country
	}{
		{
			// Based on country expects these values
			countryName: "Slovenia",
			sampleData: structures.Country{
				Borders:     []string{"AUT", "HRV", "ITA", "HUN"},
				CountryCode: "SI",
				Name: struct {
					Common string `json:"common"`
				}{
					Common: "Slovenia",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.countryName, func(t *testing.T) {
			// Create a ResponseRecorder to capture the output of the GetCountryData function
			rr := httptest.NewRecorder()

			// Call GetCountryData function with the test country file and the test case's country name
			result, err := GetCountryData(rr, tc.countryName, structures.TESTCOUNTRYFILE)
			if err != nil {
				t.Errorf("GetCountryData failed: %v", err)
				return
			}

			// Check if the result contains the expected country data
			if len(result) == 0 {
				t.Errorf("No country data returned for %s", tc.countryName)
				return
			}

			countryFound := false
			for _, country := range result {
				if strings.EqualFold(country.Name.Common, tc.countryName) {
					countryFound = true

					if country.CountryCode != tc.sampleData.CountryCode {
						t.Errorf("Expected CountryCode %s, got %s", tc.sampleData.CountryCode, country.CountryCode)
					}

					for i, border := range country.Borders {
						if border != tc.sampleData.Borders[i] {
							t.Errorf("Expected border %s, got %s", tc.sampleData.Borders[i], border)
						}
					}

					break
				}
			}

			if !countryFound {
				t.Errorf("Expected country data for %s not found", tc.countryName)
			}
		})
	}
}


//checks based on the struct if the printing is accurate
func TestPrintData(t *testing.T) {
	testCases := []struct {
		name     string
		data     interface{}
		hasError bool
		expected string
	}{
		{
			name: "TestPrintData_Map",
			data: map[string]interface{}{
				"countryName": "United States",
				"countryCode": "USA",
				"year":        2022,
				"percentage":  65.5,
			},
			hasError: false,
			expected: `{"countryCode":"USA","countryName":"United States","percentage":65.5,"year":2022}`,
		},
		{
			name:     "TestPrintData_InvalidData",
			data:     func() {},
			hasError: true,
			expected: "Error encoding API response: json: unsupported type: func()",
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			w := httptest.NewRecorder()

			PrintData(w, tc.data)

			response := w.Result()
			body, err := io.ReadAll(response.Body)
			if err != nil {
				t.Errorf("Error reading response body: %v", err)
			}

			actual := strings.TrimSpace(string(body))

			if tc.hasError && response.StatusCode != http.StatusInternalServerError {
				t.Errorf("Expected status code: %d, got: %d", http.StatusInternalServerError, response.StatusCode)
			}

			if actual != tc.expected {
				t.Errorf("Expected output: %s, got: %s", tc.expected, actual)
			}
		})
	}
}

// function used to stub the web service which provides country information in JSON format
// this will be used to test function on it without relying on the access to the web service
func StubbingCountryData() {
	//requires the URL of the webservice
	resp, err := http.Get("http://129.241.150.113:8080/v3.1/all")
	if err != nil {
		fmt.Printf("Error fetching data: %v\n", err)
		return
	}
	defer resp.Body.Close()

	//reads all it's content
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("Error reading response body: %v\n", err)
		return
	}

	//parses the JSON data into a list of Country structs
	var countries []structures.Country
	err = json.Unmarshal(body, &countries)
	if err != nil {
		fmt.Printf("Error unmarshaling JSON data: %v\n", err)
		return
	}

	//saves the list of Country structs to a JSON file which will contain all the country information
	file, err := os.Create("./countriesData.json")
	if err != nil {
		fmt.Printf("Error creating file: %v\n", err)
		return
	}
	defer file.Close()

	//sets the file to be JSON
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(countries); err != nil {
		fmt.Printf("Error writing JSON data to file: %v\n", err)
		return
	}

	fmt.Println("Successfully saved countriesData.json")
}
