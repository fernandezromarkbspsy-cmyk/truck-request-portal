package users

import (
	"context"
	"strings"
)

// ProcessClerkWebhook contains the business logic for syncing users.
func ProcessClerkWebhook(ctx context.Context, payload *ClerkWebhookPayload) error {
	// Determine if user is FTE or Backroom based on email domain or username pattern
	// For this example, we assume a default role of 'ops_pic' unless specified otherwise.
	// In production, you might check against an allowed list of FTE emails.

	isFTE := strings.HasSuffix(payload.Data.EmailAddress, "@yourcompany.com") // Example logic
	role := "ops_pic"                                                         // Default backroom role
	if isFTE {
		role = "fte_ops" // Default FTE role
	}

	user := &User{
		ClerkID:  payload.Data.ID,
		Name:     payload.Data.FirstName + " " + payload.Data.LastName,
		Email:    payload.Data.EmailAddress,
		OpsID:    payload.Data.Username, // Maps Clerk Username to our ops_id
		Role:     role,
		IsFTE:    isFTE,
		IsActive: true,
	}

	return UpsertUser(ctx, user)
}
