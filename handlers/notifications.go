package handlers

import (
	"context"
	"encoding/json"
	"net/http"

	"groupXX/firebase"
	"groupXX/structures"
	"groupXX/functions"
)

func NotificationsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	//provides different methods for the notification handler
	case http.MethodGet:
		NotificationsGetRequest(w, r)
	case http.MethodPost:
		NotificationsPostRequest(w, r)
	case http.MethodDelete:
		NotificationsDeleteRequest(w, r)
	default:
		http.Error(w, "REST Method '"+r.Method+"' not supported. Currently only '"+http.MethodGet+","+
			http.MethodDelete+" and "+http.MethodPost+"' are supported.", http.StatusNotImplemented)
		return
	}
}

//get request
func NotificationsGetRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	basePath := structures.NOTIFICATIONS_PATH
	//retrieves the user inputted id
	id := r.URL.Path[len(basePath):]

	//opens firebase project and creates a client
	ctx := context.Background()
	client, err := firebase.CreateFirestoreClient(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//gets the webhook based on the id
	wh, err := firebase.GetWebhook(ctx, client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//and returns it
	functions.PrintData(w, wh)
}

//post request
func NotificationsPostRequest(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	//creates instance of webhook
	wh := structures.Webhook{}
	err := json.NewDecoder(r.Body).Decode(&wh)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
	}

	//opens the firestore project for uses and Stores a webhook there with the given id
	ctx := context.Background()
	client, err := firebase.CreateFirestoreClient(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
	id, err := firebase.StoreWebhooks(ctx, client, wh)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//registers with the registration struct
	resp := structures.WebhookRegistration{
		ID:      id,
		Webhook: wh,
	}

	functions.PrintData(w, resp)
}

//delete request
func NotificationsDeleteRequest(w http.ResponseWriter, r *http.Request) {
	basePath := structures.NOTIFICATIONS_PATH
	//user inputted id which they want to delete
	id := r.URL.Path[len(basePath):]

	ctx := context.Background()
	client, err := firebase.CreateFirestoreClient(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	//deletes webhook based on id
	err = firebase.DeleteWebhook(ctx, client, id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}
