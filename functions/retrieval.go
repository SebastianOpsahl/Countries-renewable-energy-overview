package functions

import (
	"encoding/csv"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
	"unicode"

	"groupXX/structures"
)

// BST to store DataEntry structs which are equal to a line of the csv file based on first letter of
// countryname
var SpecifiedData *structures.BSTNode

// function initialized at start without having to be called
func init() {
	//retrievals all
	arraysWithData, err := RetrieveAll(structures.FILEPATH, false)
	if err != nil {
		fmt.Printf("Error reading CSV file")
		return
	}

	//turns into structure for faster retrieval
	SpecifiedData = PartitionDataByFirstLetter(arraysWithData)
}

// list of DataEntry struct which all are of current year
var OnlyCurrent []structures.DataEntry

func init() {
	var err error
	//calls RetrieveAll function with current set to true
	OnlyCurrent, err = RetrieveAll(structures.FILEPATH, true)
	if err != nil {
		log.Printf("Error retrieving countries from file: %v", err)
	}
}

//Time complexity in O notation: O((log c)+d)
//where c is the amount of different country first letter and d is the maximum number of countries sharing the same first letter
//in country name because our algorithm stores presaves the countries with their information in maps based on the first letter
//of the country name instead of in the CSV file. The maps are then placed placed in a BST based on the size of the map key.
//It then goes into that map and goes trough all of the countries in that map and checking which matches with the specifications.

//Example worst scenario:
//20 countries with different first letter, where of each there are 10
//just reading trough file: 200
//our algorithm: (log₂(21)) - 1 ≈ 4 -> 4 + 10 = 14

// extracts array of structs based on country name
func ExtractByMap(country string) ([]structures.DataEntry, error) {
	if SpecifiedData == nil {
		return nil, fmt.Errorf("Error: no data loaded")
	}
	countryLetter := unicode.ToUpper(rune(country[0]))
	//returns a map with the countries with that first letter
	entriesWithKey := SearchBST(SpecifiedData, countryLetter)
	if len(entriesWithKey) == 0 {
		//no point in writing error for the server, since there isn't a problem with the server, just wrong input
		//which will be dealt with in the function it is called from
		return nil, nil
	}
	//based on the countries with that first letter find the matching one
	matchingCountries := FindMatchingCountries(entriesWithKey, country)
	return matchingCountries, nil
}

// searches the binary tree based on letter input
func SearchBST(node *structures.BSTNode, letter rune) []structures.DataEntry {
	//if node doesn't exist
	if node == nil {
		return nil
	}
	//go left is letter is smaller, right for bigger and do recusive call
	if letter < node.Letter {
		return SearchBST(node.Left, letter)
	} else if letter > node.Letter {
		return SearchBST(node.Right, letter)
	}
	return node.Data
}

// extracts the entries with given key
func ExtractEntriesWithKey(partitionedData map[rune][]structures.DataEntry, key rune) []structures.DataEntry {
	//takes to uppercode so regardless of the case of user input matches database
	upperKey := unicode.ToUpper(key)

	//creates storage for entries
	entries := make([]structures.DataEntry, 0)
	if upperData, ok := partitionedData[upperKey]; ok {
		//appends with that key
		entries = append(entries, upperData...)
	}

	return entries
}

func FindMatchingCountries(data []structures.DataEntry, country string) []structures.DataEntry {
	//storage for matching countries
	matchingCountries := make([]structures.DataEntry, 0)
	countryLower := strings.ToLower(country)

	//compares on a lowercase level so the cases is no problem
	for _, entry := range data {
		var entryCountryLower string
		if IsCountryCode(countryLower) {
			entryCountryLower = strings.ToLower(entry.CountryCode)
		} else {
			entryCountryLower = strings.ToLower(entry.Country)
		}
		if countryLower == entryCountryLower {
			matchingCountries = append(matchingCountries, entry)
		}
	}
	if len(matchingCountries) == 0 {
		return nil
	}

	return matchingCountries
}

// based on lists of DataEntry structs returns them mapped into structured form
func PartitionDataByFirstLetter(data []structures.DataEntry) *structures.BSTNode {
	//sets root
	var root *structures.BSTNode

	for _, entry := range data {
		//retrieves first letter
		firstLetter := unicode.ToUpper(rune(entry.Country[0]))
		//calls inserIntoBST for each map
		root = InsertIntoBST(root, firstLetter, entry)
	}
	return root
}

func RetrieveAll(filePath string, current bool) ([]structures.DataEntry, error) {
	//list of structs
	var data []structures.DataEntry

	//opens given file
	file, err := os.Open(filePath)
	if err != nil {
		log.Printf("Error opening CSV file: %v", err)
		return nil, err
	}
	//closes it at the end of the function
	defer file.Close()

	//creates a new csv reader
	reader := csv.NewReader(file)

	//reads first line and don't do anything about it because it is headers
	//skips the first lines so don't care about the value
	_, err = reader.Read()
	if err != nil {
		log.Printf("Error, not being able to read CSV file: %v", err)
		return nil, err
	}

	//infinite loop
	for {
		//reads line
		record, err := reader.Read()

		//if end of file break the loop
		if err == io.EOF {
			break
		}
		//if error occured
		if err != nil {
			log.Printf("Error while reading the CSV file: %v", err)
			return nil, err
		}

		//new instance of struct
		entry := structures.DataEntry{}

		//fills in the variables in the struct
		entry.Country = record[0]
		entry.CountryCode = record[1]

		//reads year as string and converts to int
		year, err := strconv.Atoi(record[2])
		if err != nil {
			log.Printf("error parsing year: %v", err)
			continue
		}

		entry.Year = year

		//checks if current is true and if the entry year is not 2021, continue to the next iteration
		if current && entry.Year != 2021 {
			continue
		}

		//percentage
		percentage, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			log.Printf("error parsing percentage: %v", err)
			continue
		}
		entry.Percentage = percentage

		//appends newly created struct to list of structs
		data = append(data, entry)
	}

	//if successful return the data and nil error occurred
	return data, nil
}

// take node and letter
func InsertIntoBST(node *structures.BSTNode, letter rune, entry structures.DataEntry) *structures.BSTNode {
	//if node is nil return the BST node with it's data
	if node == nil {
		return &structures.BSTNode{Data: []structures.DataEntry{entry}, Letter: letter}
	}
	//left for less right for greater else become root
	if letter < node.Letter {
		node.Left = InsertIntoBST(node.Left, letter, entry)
	} else if letter > node.Letter {
		node.Right = InsertIntoBST(node.Right, letter, entry)
	} else {
		node.Data = append(node.Data, entry)
	}
	return node
}
