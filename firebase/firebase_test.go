package firebase

import (
	"context"
	"testing"
	"net/http/httptest"
	"net/http"
	"os"
	"time"
	"strings"
	"reflect"

	"github.com/stretchr/testify/assert"
	"github.com/google/uuid"

	"groupXX/structures"
)

func TestUpdateCalls(t *testing.T) {
	// Create a new test server
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Assert that the correct URL was called
		assert.Equal(t, "http://testwebhook.com", r.URL.String())
	}))

	// Ensure that the test server is closed after the test
	defer ts.Close()

	// Set the environment variable to the test server URL
	os.Setenv("WEBHOOK_TEST_URL", ts.URL)

	// Call the function that we want to test
	//ctx := context.Background()
	err := UpdateCalls("germany")
	assert.NoError(t, err)
}

func TestStoreWebhooks(t *testing.T) {
    // create a context with a 5-second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // create a new httptest server to simulate the API
    srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusCreated)
    }))
    defer srv.Close()

    // create a new webhook to store
    webhook := structures.Webhook{
        URL:      uuid.NewString(),
        Country: "germany",
        Calls:   2,
    }

    // create a new Firestore client
    client, err := CreateFirestoreClient(ctx)
    if err != nil {
        t.Fatalf("CreateFirestoreClient() returned an error: %v", err)
    }
    defer client.Close()

    // store the webhook
    _, err = StoreWebhooks(ctx, client, webhook)
    if err != nil {
        t.Fatalf("StoreWebhooks() returned an error: %v", err)
    }

    // retrieve the stored webhook
    storedWebhook, err := GetWebhook(ctx, client, webhook.URL)
    if err != nil {
        t.Fatalf("GetWebhook() returned an error: %v", err)
    }

    // check if the stored webhook is the same as the original webhook
    if !reflect.DeepEqual(storedWebhook, webhook) {
        t.Errorf("Stored webhook %+v does not match original webhook %+v", storedWebhook, webhook)
    }
}

func TestGetWebhook(t *testing.T) {
    // create a context with a 5-second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // create a new Firestore client
    client, err := CreateFirestoreClient(ctx)
    if err != nil {
        t.Fatalf("CreateFirestoreClient() returned an error: %v", err)
    }
    defer client.Close()

    // create a new webhook to store
    webhook := structures.Webhook{
        URL:      uuid.NewString(),
        Country: "germany",
        Calls:   2,
    }

    // store the webhook
    _, err = StoreWebhooks(ctx, client, webhook)
    if err != nil {
        t.Fatalf("StoreWebhooks() returned an error: %v", err)
    }

    // get the stored webhook
    storedWebhook, err := GetWebhook(ctx, client, webhook.URL)
    if err != nil {
        t.Fatalf("GetWebhook() returned an error: %v", err)
    }

    // check if the stored webhook is the same as the original webhook
    if !reflect.DeepEqual(storedWebhook, webhook) {
        t.Errorf("Stored webhook %+v does not match original webhook %+v", storedWebhook, webhook)
    }

    // delete the stored webhook
    err = DeleteWebhook(ctx, client, webhook.URL)
    if err != nil {
        t.Fatalf("DeleteWebhook() returned an error: %v", err)
    }
}

func TestGetWebhookNotFound(t *testing.T) {
    // create a context with a 5-second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // create a new Firestore client
    client, err := CreateFirestoreClient(ctx)
    if err != nil {
        t.Fatalf("CreateFirestoreClient() returned an error: %v", err)
    }
    defer client.Close()

    // try to get a non-existent webhook
    _, err = GetWebhook(ctx, client, "non-existent-url")
    if !strings.Contains(err.Error(), "webhook not found") {
        t.Fatalf("GetWebhook() returned unexpected error: %v", err)
    }
}

func TestDeleteWebhookNotFound(t *testing.T) {
    // create a context with a 5-second timeout
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()

    // create a new Firestore client
    client, err := CreateFirestoreClient(ctx)
    if err != nil {
        t.Fatalf("CreateFirestoreClient() returned an error: %v", err)
    }
    defer client.Close()

    // try to delete a non-existent webhook
    err = DeleteWebhook(ctx, client, "non-existent-url")
    if !strings.Contains(err.Error(), "webhook not found") {
        t.Fatalf("DeleteWebhook() returned unexpected error: %v", err)
    }
}

func TestGetNumWebhooks(t *testing.T) {
	// create a context with a 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// create a new Firestore client
	client, err := CreateFirestoreClient(ctx)
	if err != nil {
		t.Fatalf("CreateFirestoreClient() returned an error: %v", err)
	}
	defer client.Close()

	// store some webhooks
	webhooks := []structures.Webhook{
		{URL: "webhook1", Country: "germany", Calls: 2},
		{URL: "webhook2", Country: "germany", Calls: 4},
		{URL: "webhook3", Country: "spain", Calls: 1},
	}
	for _, webhook := range webhooks {
		if _, err := StoreWebhooks(ctx, client, webhook); err != nil {
			t.Fatalf("StoreWebhooks() returned an error: %v", err)
		}
	}

	// check that the correct number of webhooks are returned
	numWebhooks, err := GetNumWebhooks(ctx, client)
	if err != nil {
		t.Fatalf("GetNumWebhooks() returned an error: %v", err)
	}
	if numWebhooks != len(webhooks) {
		t.Errorf("GetNumWebhooks() returned incorrect number of webhooks: got %d, want %d", numWebhooks, len(webhooks))
	}
}
