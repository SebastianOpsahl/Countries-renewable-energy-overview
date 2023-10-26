package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"groupXX/firebase"
	"groupXX/functions"
	"groupXX/structures"
)

func CurrentHandler(w http.ResponseWriter, r *http.Request) {
	//takes method of the request, if GET then forward to function, else write to the user that only GET is allowed
	switch r.Method {
	case http.MethodGet:
		CurrentGetHandler(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+
			"' are supported.", http.StatusNotImplemented)
		return
	}
}

func CurrentGetRequest(w http.ResponseWriter, r *http.Request) (country string, neighbours bool, err error) {
	basePath := structures.RENEWABLECURRENT_PATH

	//returns path other than basePath
	country = r.URL.Path[len(basePath):]
	err = firebase.UpdateCalls(country)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	//returns parsed query parameters in a map
	queryParams := r.URL.Query()
	//gets the neighbours query from the map
	neighboursStr := queryParams.Get("neighbours")

	//temporary value
	neighbours = false

	//if the user inputted neighbour specifications
	if neighboursStr != "" {
		//if so turn into bool
		neighbours, err = strconv.ParseBool(neighboursStr)
		if err != nil {
			http.Error(w, "Error parsing neighbours value from string to bool: "+err.Error(), http.StatusBadRequest)
			return
		}
	}
	return country, neighbours, nil
}

func CurrentGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	countryName, neighbours, err := CurrentGetRequest(w,r)
	if err != nil{
		log.Printf("Error parsing URL")
	}
	
	//calls file to return the countries as a struct with the specification of country name, true because it only returns
	//the current value 2021 not history, nil and nil for the optional value of begin and end year of search
	data, err := functions.ReadCountryInfo(w, countryName, true, nil, nil)
	if err != nil {
		log.Printf("Error reading CSV file: %v", err)
	}
	if data == nil {
		fmt.Fprint(w, "No return for the given search found")
	} else {

		//if the user wants to see neighbours aswell
		if neighbours {
			for _, entry := range data {
				//calls function to retrieve neighbours
				neighbours, err := functions.RetrieveNeighbours(w, entry.Country)

				if err != nil {
					log.Printf("Error retrieving neighbours: %v", err)
					http.Error(w, "Error retrieving neighbours: "+err.Error(), http.StatusInternalServerError)
					return
				}

				//if the country has neighbours loop trough neighbours
				if len(neighbours) > 0 {
					for _, currentNeighbour := range neighbours {
						//update number of calls
						err = firebase.UpdateCalls(currentNeighbour)
						if err != nil {
							http.Error(w, err.Error(), http.StatusInternalServerError)
						}
						neigh, err := functions.ReadCountryInfo(w, currentNeighbour, true, nil, nil)
						if err != nil {
							log.Printf("Error reading CSV file: %v", err)
						} 

						//in cases where the neighbour counrty doesn't exist in the csv file
						if len(neigh) > 0 {
							data = append(data, neigh[0])
						}
					}
				}
			}
		}
		functions.PrintData(w, data)
	}
}
