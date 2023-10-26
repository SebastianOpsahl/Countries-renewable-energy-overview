package main

import (
	"time"
	"context"
	"log"
	"net/http"
	"os"

	"groupXX/handlers"
	"groupXX/structures"
	"groupXX/firebase"
)

func main() {
	// Create a Firestore client
	ctx := context.Background()
	client, err := firebase.CreateFirestoreClient(ctx)
	if err != nil {
		log.Fatalf("Error creating Firestore client: %v", err)
	}
	defer client.Close()

	// Create a ticker to purge old cache entries every daysThreshold days
	purgeInterval := time.Duration(structures.DAYSTHRESHOLD) * 24 * time.Hour
	ticker := time.NewTicker(purgeInterval)

	// Run the PurgeOldCacheEntries function in the background
	go func() {
		for range ticker.C {
			firebase.PurgeOldCacheEntries(ctx, client, structures.DAYSTHRESHOLD)
		}
	}()
	

	port := os.Getenv("PORT")
	//if no specified port
	if port == "" {
		log.Println("$PORT has not been set. Default: 8080")
		//sets 8080 to default if nothing else specified
		port = "8080"
	}

	//based on the path (defined in constans.go), it forwards it to the corresponding function
	http.HandleFunc(structures.DEFAULT_PATH, handlers.DefaultHandler)
	http.HandleFunc(structures.RENEWABLECURRENT_PATH, handlers.CurrentHandler)
	http.HandleFunc(structures.RENEWABLEHISTORY_PATH, handlers.HistoryHandler)
	http.HandleFunc(structures.NOTIFICATIONS_PATH, handlers.NotificationsHandler)
	http.HandleFunc(structures.STATUS_PATH, handlers.StatusHandler)
	http.HandleFunc(structures.INFO_PATH, handlers.InfoHandler)

	//starts server
	log.Println("Starting server on port " + port + " ...")
	log.Fatal(http.ListenAndServe(":"+port, nil))
}