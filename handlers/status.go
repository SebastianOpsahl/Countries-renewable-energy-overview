package handlers

//missing notification_db and webhooks!!!!!!!!!!!!!!!!!!!

import (
	"context"
	"log"
	"net/http"
	"strings"
	"time"

	"groupXX/firebase"
	"groupXX/functions"
	"groupXX/structures"
)

// variable which equals to the time it is started (when the application starts)
var startTime = time.Now()

func StatusHandler(w http.ResponseWriter, r *http.Request) {
	//takes method of the request, if GET then forward to function, else write info to user
	switch r.Method {
	case http.MethodGet:
		StatusGetHandler(w, r)
	default:
		http.Error(w, "This service only offers the <a href=\"/diag\">/diag endpoint</a> that shows comprehensive HTTP "+
			"information for any received request. Please redirect any request to that endpoint.", http.StatusInternalServerError)
		return
	}
}

func StatusGetHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//checks the status code for the URLs of the third party APIs
	respCountry, err := http.Get("https://restcountries.com/")
	if err != nil {
		log.Printf("Error when getting country API: %v %s", respCountry, err)
		http.Error(w, "Error when getting country API", http.StatusInternalServerError)
		return
	}
	//closes its body at the end of the function
	defer respCountry.Body.Close()

	respNotif, err := http.Get("https://console.firebase.google.com/project/group66assignment2/firestore/data/~2Fwebhooks~2F45NVjImCjz2pGEDAPupD")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer respNotif.Body.Close()

	ctx := context.Background()
	client, err := firebase.CreateFirestoreClient(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	numWh, err := firebase.GetNumWebhooks(ctx, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//fills the struct
	output := structures.Info{
		RESTStatus:  respCountry.StatusCode,
		NotifStatus: respNotif.StatusCode,
		Webhooks:    numWh,
		Version:     strings.Split(r.URL.Path, "/")[2],
		Uptime:      time.Now().Sub(startTime).Seconds(),
	}

	//and prints it out
	functions.PrintData(w, output)
}
