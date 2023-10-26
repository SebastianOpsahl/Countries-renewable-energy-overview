package handlers

import (
	"net/http"
	"testing"
	"net/http/httptest"

	"github.com/stretchr/testify/assert"
)

func TestCurrentGetRequest(t *testing.T) {
	testCases := []struct {
		url       string
		country   string
		neighbours bool
	}{
		{
			url:       "/energy/v1/renewables/current/germany?neighbours=true",
			country:   "germany",
			neighbours: true,
		},
		{
			url:       "/energy/v1/renewables/current/germany",
			country:   "germany",
			neighbours: false,
		},
		{
			url:       "/energy/v1/renewables/current/spain?neighbours=false",
			country:   "spain",
			neighbours: false,
		},
	}

	for _, tc := range testCases {
		t.Run("", func(t *testing.T) {
			req, err := http.NewRequest("GET", tc.url, nil)
			assert.NoError(t, err)

			// Create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
			rr := httptest.NewRecorder()

			country, neighbours, err := CurrentGetRequest(rr, req)

			assert.NoError(t, err)
			assert.Equal(t, tc.country, country)
			assert.Equal(t, tc.neighbours, neighbours)
		})
	}
}