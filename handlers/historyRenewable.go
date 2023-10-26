package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sort"
	"strconv"
	"strings"

	"groupXX/firebase"
	"groupXX/functions"
	"groupXX/structures"
)

func HistoryHandler(w http.ResponseWriter, r *http.Request) {
	//takes method of the request, if GET then forward to function, else write to the user that only GET is allowed
	switch r.Method {
	case http.MethodGet:
		HistoryGetHandler(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' are supported.", http.StatusNotImplemented)
		return
	}
}

func HistoryGetRequest(w http.ResponseWriter, r *http.Request) (country string, begin int, end int, sorting bool, err error) {
	w.Header().Add("content-type", "application/json")

	//sets the basepath so we can work on top of that
	basePath := "/energy/v1/renewables/history/"
	if !strings.HasPrefix(r.URL.Path, basePath) {
		http.Error(w, "Invalid URL path", http.StatusBadRequest)
		return
	}

	//extract the country value from the path
	country = strings.TrimPrefix(r.URL.Path, basePath)
	err = firebase.UpdateCalls(country)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//extract the query parameters
	queryParams := r.URL.Query()
	beginStr := queryParams.Get("begin")
	endStr := queryParams.Get("end")
	sortingByValueStr := queryParams.Get("sortByValue")

	begin = 0
	end = 0

	//convert begin and end years from strings to integers if they are specified if not let them be 0
	//the 0 value will be used as an identifier when it is called 
	if beginStr != "" {
		begin, err = strconv.Atoi(beginStr)
		if err != nil {
			http.Error(w, "Error parsing begin year string to integer", http.StatusBadRequest)
			return
		}
	}

	if endStr != "" {
		end, err = strconv.Atoi(endStr)
		if err != nil {
			http.Error(w, "Error parsing end year string to integer", http.StatusBadRequest)
			return
		}
	}

	//parse the sortingByValue query parameter as a boolean, if it's provided
	sorting = false
	if sortingByValueStr != "" {
		sorting, err = strconv.ParseBool(sortingByValueStr)
		if err != nil {
			http.Error(w, "Error parsing sorting value string to bool", http.StatusBadRequest)
			return
		}
	}

	return country, begin, end, sorting, nil
}

func HistoryGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	countryName, begin, end, sorting, err := HistoryGetRequest(w,r)
	if err != nil{
		log.Printf("Error parsing URL")
	}

	var data []structures.DataEntry

	//handle the begin and end specifications, turned them into pointers to deal with their absence
	if begin == 0 && end == 0 {
		data, err = functions.ReadCountryInfo(w, countryName, false, nil, nil)
	} else if begin == 0 && end != 0 {
		data, err = functions.ReadCountryInfo(w, countryName, false, nil, &end)
	} else if begin != 0 && end == 0 {
		data, err = functions.ReadCountryInfo(w, countryName, false, &begin, nil)
	} else if begin != 0 && end != 0 {
		data, err = functions.ReadCountryInfo(w, countryName, false, &begin, &end)
	}

	//based on the potential calls of the ReadCountryInfo, checks if data is returned (found)
	if data == nil {
		fmt.Fprint(w, "No return for the given search found")
	} else {
		if err != nil {
			log.Printf("Error reading CSV file: %v", err)
		}

		//sort the data based on Percentage if sorting is true
		if sorting {
			sort.Sort(ByPercentage(data))
		}

		functions.PrintData(w, data)
	}
}

//list of DataEntry structs to store in a structured way
type ByPercentage []structures.DataEntry

//functions to return length, swap elements and return the less big
func (a ByPercentage) Len() int           { return len(a) }
func (a ByPercentage) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByPercentage) Less(i, j int) bool { return a[i].Percentage < a[j].Percentage }