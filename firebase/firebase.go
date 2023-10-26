package firebase

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"groupXX/structures"
	"log"
	"net/http"
	"time"
	"strings"

	"cloud.google.com/go/firestore"
	"google.golang.org/api/iterator"
	"google.golang.org/api/option"
)

func CreateFirestoreClient(ctx context.Context) (*firestore.Client, error) {
	//sets up the Firestore client
	projectID := "group66assignment2"
	//refers to the credentials JSON file
	credentials := "./.secrets/group66assignment2-firebase-adminsdk-pv6iv-c0b5b34aeb.json"
	opt := option.WithCredentialsFile(credentials)
	//creates a client eith the context, projectID and credentials
	client, err := firestore.NewClient(ctx, projectID, opt)
	if err != nil {
		return nil, err
	}

	return client, nil
}

// stores a new webhook to the firestore
func StoreWebhooks(ctx context.Context, client *firestore.Client, webhook structures.Webhook) (string, error) {
	doc, _, err := client.Collection("webhooks").Add(ctx, webhook)
	if err != nil {
		return "", err
	}
	return doc.ID, nil
}

// deletes a webhook from the firestore
func DeleteWebhook(ctx context.Context, client *firestore.Client, id string) error {
	_, err := client.Collection("webhooks").Doc(id).Delete(ctx)
	return err
}

// retrieves a webhook from the firestore
func GetWebhook(ctx context.Context, client *firestore.Client, id string) (structures.Webhook, error) {
	wh := structures.Webhook{}
	snapshot, err := client.Collection("webhooks").Doc(id).Get(ctx)
	if err != nil {
		return wh, err
	}
	err = snapshot.DataTo(&wh)
	if err != nil {
		return wh, err
	}
	return wh, nil
}

func GetNumWebhooks(ctx context.Context, client *firestore.Client) (int, error) {
	webhooks, err := client.Collection("webhooks").Documents(ctx).GetAll()
	if err != nil {
		return 0, err
	}
	return len(webhooks), nil
}

var calls = make(map[string]int)

// Updates call count for country and checks if any of the webhooks are to be invocated
func UpdateCalls(country string) error {
	//increments call by 1
	calls[country] += 1
	ctx := context.Background()
	client, err := CreateFirestoreClient(ctx)
	if err != nil {
		return err
	}
	iter := client.Collection("webhooks").Documents(ctx)
	defer iter.Stop()
	//infinite for loop to loop until break
	for {
		doc, err := iter.Next()
		//if iterated trough break the loop
		if err == iterator.Done {
			break
		}
		if err != nil {
			return err
		}
		//struct for webhook
		wh := structures.Webhook{}
		//parse data into struct
		err = doc.DataTo(&wh)
		if err != nil {
			return err
		}
		//if country name matches and mod is = 0
		if wh.Country == country && calls[country]%wh.Calls == 0 {
			jsonData, err := json.Marshal(wh)
			if err != nil {
				return err
			}
			_, err = http.Post(wh.URL, "application/json", bytes.NewBuffer(jsonData))
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// if the data is not found on the stack this function will be called to place it there
func SetCachedData(ctx context.Context, cacheKey string, data []structures.DataEntry) error {
	client, err := CreateFirestoreClient(ctx)
	if err != nil {
		return fmt.Errorf("Error creating Firestore client: %v", err)
	}
	defer client.Close()

	cacheKey = strings.ToLower(cacheKey)

	// uses the cache collection in the firestore project
	cacheCollectionRef := client.Collection("cache")
	// sets key
	query := cacheCollectionRef.Where("key", "==", cacheKey).Limit(1)
	iter := query.Documents(ctx)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("Error marshaling data to JSON: %v", err)
	}

	// iterates to next instance
	doc, err := iter.Next()

	if err == iterator.Done {
		newCacheEntry := map[string]interface{}{
			"key":       cacheKey,
			"data":      string(jsonData),
			"timestamp": time.Now(),
			"hits":      1,
		}

		// add new cache
		_, _, err = cacheCollectionRef.Add(ctx, newCacheEntry)
		if err != nil {
			return fmt.Errorf("Error adding cache entry to Firestore: %v", err)
		}
	} else if err != nil {
		return fmt.Errorf("Error iterating Firestore documents: %v", err)
		// updates the found cached data
	} else {
		docRef := doc.Ref
		_, err = docRef.Update(ctx, []firestore.Update{
			{Path: "data", Value: string(jsonData)},
			{Path: "timestamp", Value: time.Now()},
		})
		if err != nil {
			return fmt.Errorf("Error updating cache entry in Firestore: %v", err)
		}
	}

	return nil
}

// gets cached data
func GetCachedData(ctx context.Context, cacheKey string) ([]structures.DataEntry, error) {
	client, err := CreateFirestoreClient(ctx)
	if err != nil {
		return nil, fmt.Errorf("Error creating Firestore client: %v", err)
	}
	defer client.Close()

	cacheKey = strings.ToLower(cacheKey)

	// from the cache collection
	cacheCollectionRef := client.Collection("cache")
	// where the key matches the cacheKey (user inputted) with a limit of 1
	query := cacheCollectionRef.Where("key", "==", cacheKey).Limit(1)
	iter := query.Documents(ctx)

	doc, err := iter.Next()

	// if iterated through and nothing is found return nil
	if err == iterator.Done {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("Error iterating Firestore documents: %v", err)
	}

	// updating cached data with incrementing hit count if found
	docRef := doc.Ref
	_, err = docRef.Update(ctx, []firestore.Update{
		{Path: "hits", Value: doc.Data()["hits"].(int64) + 1},
	})
	if err != nil {
		return nil, fmt.Errorf("Error updating hit count in Firestore: %v", err)
	}

	// Retrieve the JSON string from Firestore
	jsonString, ok := doc.Data()["data"].(string)
	if !ok {
		return nil, fmt.Errorf("Error converting data to string")
	}

	// Unmarshal the JSON string into a []structures.DataEntry
	var data []structures.DataEntry
	err = json.Unmarshal([]byte(jsonString), &data)
	if err != nil {
		return nil, fmt.Errorf("Error unmarshaling JSON to DataEntry: %v", err)
	}

	return data, nil
}

// deletes caches older than a specific treshold
func PurgeOldCacheEntries(ctx context.Context, client *firestore.Client, daysThreshold int) {
	cacheCollectionRef := client.Collection("cache")
	query := cacheCollectionRef.Where("timestamp", "<", time.Now().AddDate(0, 0, -daysThreshold))
	iter := query.Documents(ctx)

	for {
		doc, err := iter.Next()
		//iterated trough them all
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Printf("Error iterating Firestore documents: %v", err)
			return
		}

		//delete the document
		_, err = doc.Ref.Delete(ctx)
		if err != nil {
			log.Printf("Error deleting Firestore document: %v", err)
		} else {
			log.Printf("Treshold met, deleted document with ID: %s", doc.Ref.ID)
		}
	}
}