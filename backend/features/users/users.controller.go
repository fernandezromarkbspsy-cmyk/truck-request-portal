package users

import (
	"encoding/json"
	"net/http"
)

// HandleClerkWebhook receives the HTTP request from Clerk
func HandleClerkWebhook(w http.ResponseWriter, r *http.Request) {
	var payload ClerkWebhookPayload
	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid payload", http.StatusBadRequest)
		return
	}

	// In a real production environment, you MUST verify the Clerk webhook signature here
	// using the Clerk Go SDK to ensure the request actually came from Clerk.

	ctx := r.Context()
	err := ProcessClerkWebhook(ctx, &payload)
	if err != nil {
		http.Error(w, "Failed to sync user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("User synced successfully"))
}
