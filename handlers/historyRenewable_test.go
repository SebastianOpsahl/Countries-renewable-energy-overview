package handlers

import (
	"testing"
	"reflect"
	"sort"
	"net/http/httptest"
	"net/http"

	"github.com/stretchr/testify/assert"

	"groupXX/structures"

)

func TestByPercentage(t *testing.T) {
	testCases := []struct {
		input    ByPercentage
		expected ByPercentage
	}{
		{
			input: ByPercentage{
				structures.DataEntry{Percentage: 5.0},
				structures.DataEntry{Percentage: 1.0},
				structures.DataEntry{Percentage: 3.0},
			},
			expected: ByPercentage{
				structures.DataEntry{Percentage: 1.0},
				structures.DataEntry{Percentage: 3.0},
				structures.DataEntry{Percentage: 5.0},
			},
		},
		// Add more test cases as needed.
	}

	for _, tc := range testCases {
		t.Run("Sorting ByPercentage", func(t *testing.T) {
			// Sort the input ByPercentage slice using sort.Sort().
			sort.Sort(tc.input)

			// Compare the sorted input slice with the expected sorted slice.
			if !reflect.DeepEqual(tc.input, tc.expected) {
				t.Errorf("Incorrect sorting: got %v, expected %v", tc.input, tc.expected)
			}
		})
	}
}

func TestHistoryGetRequest(t *testing.T) {
	testCases := []struct {
		url       string
		country   string
		begin     int
		end       int
		sorting   bool
	}{
		{
			url:       "/energy/v1/renewables/history/germany?begin=1990&end=2020&sortByValue=true",
			country:   "germany",
			begin:     1990,
			end:       2020,
			sorting:   true,
		},
		{
			url:       "/energy/v1/renewables/history/germany",
			country:   "germany",
			begin:     0,
			end:       0,
			sorting:   false,
		},
		{
			url:       "/energy/v1/renewables/history/spain?begin=1990",
			country:   "spain",
			begin:     1990,
			end:       0,
			sorting:   false,
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			req, err := http.NewRequest("GET", tc.url, nil)
			assert.NoError(t, err)

			// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()

			country, begin, end, sorting, err := HistoryGetRequest(rr, req)

			assert.NoError(t, err)
			assert.Equal(t, tc.country, country)
			assert.Equal(t, tc.begin, begin)
			assert.Equal(t, tc.end, end)
			assert.Equal(t, tc.sorting, sorting)
		})
	}
}
