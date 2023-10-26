package functions

import (
	"testing"
	"unicode"

	"groupXX/structures"
)

//VALID
//to test if it can retrieve all countries
func TestRetrieveAll(t *testing.T) {
	//call the RetrieveAll with the test data set
	data, err := RetrieveAll("./testData.csv", false)
	if err != nil {
		t.Errorf("Expected no error, but got: %v", err)
	}

	expectedData := []structures.DataEntry{
		{
			Country:     "United States",
			CountryCode: "USA",
			Year:        2020,
			Percentage:  10.5,
		},
		{
			Country:     "United States",
			CountryCode: "USA",
			Year:        2021,
			Percentage:  10.6,
		},
		{
			Country:     "Canada",
			CountryCode: "CAN",
			Year:        2021,
			Percentage:  29.8,
		},
		{
			Country:     "Brazil",
			CountryCode: "BRA",
			Year:        2021,
			Percentage:  46.2,
		},
	}

	for i, entry := range data {
		if entry != expectedData[i] {
			t.Errorf("Expected entry %d to be %+v, but got %+v", i, expectedData[i], entry)
		}
	}
}

func TestFindMatchingCountries(t *testing.T) {
	testCases := []struct {
		data             []structures.DataEntry
		country          string
		expectedMatching []structures.DataEntry
	}{
		{
			data: []structures.DataEntry{
				{Country: "United States", CountryCode: "US"},
				{Country: "Canada", CountryCode: "CA"},
				{Country: "United Kingdom", CountryCode: "GB"},
			},
			country: "United States",
			expectedMatching: []structures.DataEntry{
				{Country: "United States", CountryCode: "US"},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.country, func(t *testing.T) {
			matchingCountries := FindMatchingCountries(tc.data, tc.country)

			if len(matchingCountries) != len(tc.expectedMatching) {
				t.Errorf("Expected %d matching countries, got %d", len(tc.expectedMatching), len(matchingCountries))
				return
			}

			for i, expected := range tc.expectedMatching {
				actual := matchingCountries[i]
				if expected.Country != actual.Country || expected.CountryCode != actual.CountryCode {
					t.Errorf("Expected matching country %+v, got %+v", expected, actual)
				}
			}
		})
	}
}

var sampleData = []structures.DataEntry{
	{Country: "United States", CountryCode: "US"},
	{Country: "United Kingdom", CountryCode: "GB"},
	{Country: "Spain", CountryCode: "ES"},
	{Country: "France", CountryCode: "FR"},
	{Country: "Germany", CountryCode: "DE"},
}

func insertBSTNode(node *structures.BSTNode, data structures.DataEntry) *structures.BSTNode {
	if node == nil {
		newNode := &structures.BSTNode{
			Letter: unicode.ToUpper(rune(data.Country[0])),
			Data:   []structures.DataEntry{data},
		}
		return newNode
	}

	firstLetter := unicode.ToUpper(rune(data.Country[0]))

	if firstLetter < node.Letter {
		node.Left = insertBSTNode(node.Left, data)
	} else if firstLetter > node.Letter {
		node.Right = insertBSTNode(node.Right, data)
	} else {
		node.Data = append(node.Data, data)
	}

	return node
}

func setUpBST(data []structures.DataEntry) *structures.BSTNode {
	var root *structures.BSTNode

	for _, entry := range data {
		root = insertBSTNode(root, entry)
	}

	return root
}

func setUpPartitionedData(data []structures.DataEntry) map[rune][]structures.DataEntry {
	partitionedData := make(map[rune][]structures.DataEntry)

	for _, entry := range data {
		firstLetter := unicode.ToUpper(rune(entry.Country[0]))
		partitionedData[firstLetter] = append(partitionedData[firstLetter], entry)
	}

	return partitionedData
}

func TestExtractByMap(t *testing.T) {
	SpecifiedData = setUpBST(sampleData)

	testCases := []struct {
		country          string
		expectedMatching []structures.DataEntry
	}{
		{
			country: "United States",
			expectedMatching: []structures.DataEntry{
				{Country: "United States", CountryCode: "US"},
			},
		},
		// Add more test cases here
	}

	for _, tc := range testCases {
		t.Run(tc.country, func(t *testing.T) {
			matchingCountries, err := ExtractByMap(tc.country)

			if err != nil {
				t.Errorf("ExtractByMap failed: %v", err)
				return
			}

			if len(matchingCountries) != len(tc.expectedMatching) {
				t.Errorf("Expected %d matching countries, got %d", len(tc.expectedMatching), len(matchingCountries))
				return
			}

			for i, expected := range tc.expectedMatching {
				actual := matchingCountries[i]
				if expected.Country != actual.Country || expected.CountryCode != actual.CountryCode {
					t.Errorf("Expected matching country %+v, got %+v", expected, actual)
				}
			}
		})
	}
}

func TestSearchBST(t *testing.T) {
	SpecifiedData = setUpBST(sampleData)

	testCases := []struct {
		letter           rune
		expectedMatching []structures.DataEntry
	}{
		{
			letter: 'U',
			expectedMatching: []structures.DataEntry{
				{Country: "United States", CountryCode: "US"},
				{Country: "United Kingdom", CountryCode: "GB"},
			},
		},
		// Add more test cases here
	}

	for _, tc := range testCases {
		t.Run(string(tc.letter), func(t *testing.T) {
			matchingEntries := SearchBST(SpecifiedData, tc.letter)

			if len(matchingEntries) != len(tc.expectedMatching) {
				t.Errorf("Expected %d matching entries, got %d", len(tc.expectedMatching), len(matchingEntries))
				return
			}

			for i, expected := range tc.expectedMatching {
				actual := matchingEntries[i]
				if expected.Country != actual.Country || expected.CountryCode != actual.CountryCode {
					t.Errorf("Expected matchingentry %v, got %v", expected, actual)
				}
		}
	})
	}
}

func TestPartitionDataByFirstLetter(t *testing.T) {
	sampleData := []structures.DataEntry{
		{Country: "United States", CountryCode: "US"},
		{Country: "United Kingdom", CountryCode: "GB"},
		{Country: "France", CountryCode: "FR"},
		{Country: "Germany", CountryCode: "DE"},
		{Country: "Spain", CountryCode: "ES"},
	}

	expectedPartitions := map[rune][]structures.DataEntry{
		'U': {
			{Country: "United States", CountryCode: "US"},
			{Country: "United Kingdom", CountryCode: "GB"},
		},
		'F': {
			{Country: "France", CountryCode: "FR"},
		},
		'G': {
			{Country: "Germany", CountryCode: "DE"},
		},
		'S': {
			{Country: "Spain", CountryCode: "ES"},
		},
	}

	root := PartitionDataByFirstLetter(sampleData)

	for letter, expectedEntries := range expectedPartitions {
		matchingEntries := SearchBST(root, letter)

		if len(matchingEntries) != len(expectedEntries) {
			t.Errorf("Expected %d matching entries for letter %c, got %d", len(expectedEntries), letter, len(matchingEntries))
			continue
		}

		for i, expected := range expectedEntries {
			actual := matchingEntries[i]
			if expected.Country != actual.Country || expected.CountryCode != actual.CountryCode {
				t.Errorf("Expected matching entry %v, got %v", expected, actual)
			}
		}
	}
}