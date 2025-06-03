package lookergo

import (
	"context"
	"fmt"
	"net/http"
	"reflect"
	"testing"
)

func TestApiCredentialsResourceOp_Get(t *testing.T) {
	setup()
	defer teardown()

	userID := 18
	credentialsID := "8"

	// Setup mux handler to mock API response
	mux.HandleFunc(fmt.Sprintf("/4.0/users/%d/credentials_api3/%s", userID, credentialsID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodGet)
		fmt.Fprint(w, `{
			"id": "8",
			"client_id": "yyy",
			"type": "api3",
			"is_disabled": false,
			"url": "https://localhost:19999/api/4.0/users/18/credentials_api3/8"
		}`)
	})

	cred, resp, err := client.ApiCredentials.Get(context.Background(), userID, credentialsID)
	if err != nil {
		t.Fatalf("ApiCredentials.Get returned error: %v", err)
	}
	_ = resp

	expected := &ApiCredential{
		ID:         "8",
		ClientID:   "yyy",
		Type:       "api3",
		IsDisabled: false,
		URL:        "https://localhost:19999/api/4.0/users/18/credentials_api3/8",
	}

	if !reflect.DeepEqual(cred, expected) {
		t.Errorf("ApiCredentials.Get returned\n got: %+v\nwant: %+v", cred, expected)
	}
}

func TestApiCredentialsResourceOp_Create(t *testing.T) {
	setup()
	defer teardown()

	userID := 18

	newCredential := &ApiCredential{
		ClientID:   "zzz",
		Type:       "api3",
		IsDisabled: false,
	}

	mux.HandleFunc(fmt.Sprintf("/4.0/users/%d/credentials_api3", userID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodPost)

		fmt.Fprint(w, `{
			"id": "9",
			"client_id": "zzz",
			"client_secret": "supersecretvalue",
			"type": "api3",
			"is_disabled": false,
			"url": "https://localhost:19999/api/4.0/users/18/credentials_api3/9"
		}`)
	})

	cred, resp, err := client.ApiCredentials.Create(context.Background(), userID, newCredential)
	if err != nil {
		t.Fatalf("ApiCredentials.Create returned error: %v", err)
	}
	_ = resp

	expected := &ApiCredential{
		ID:           "9",
		ClientID:     "zzz",
		ClientSecret: "supersecretvalue",
		Type:         "api3",
		IsDisabled:   false,
		URL:          "https://localhost:19999/api/4.0/users/18/credentials_api3/9",
	}

	if !reflect.DeepEqual(cred, expected) {
		t.Errorf("ApiCredentials.Create returned\n got: %+v\nwant: %+v", cred, expected)
	}
}

func TestApiCredentialsResourceOp_Delete(t *testing.T) {
	setup()
	defer teardown()

	userID := 18
	credentialsID := "9"

	mux.HandleFunc(fmt.Sprintf("/4.0/users/%d/credentials_api3/%s", userID, credentialsID), func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, http.MethodDelete)
		// Return 204 No Content
		w.WriteHeader(http.StatusNoContent)
	})

	resp, err := client.ApiCredentials.Delete(context.Background(), userID, credentialsID)
	if err != nil {
		t.Fatalf("ApiCredentials.Delete returned error: %v", err)
	}

	if resp == nil {
		t.Fatalf("ApiCredentials.Delete response is nil")
	}

	if resp.StatusCode != http.StatusNoContent {
		t.Errorf("ApiCredentials.Delete returned status %d, want %d", resp.StatusCode, http.StatusNoContent)
	}
}
