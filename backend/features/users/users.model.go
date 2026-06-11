package users

import "time"

type User struct {
	ID        string    `json:"id"`
	ClerkID   string    `json:"clerk_id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	OpsID     string    `json:"ops_id"`
	Role      string    `json:"role"` // ops_pic, fte_ops, fte_mm, dock_officer, doc_officer
	IsFTE     bool      `json:"is_fte"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

type ClerkWebhookPayload struct {
	Data struct {
		ID           string `json:"id"`
		Username     string `json:"username"` // We will use this for ops_id
		EmailAddress string `json:"email_address"`
		FirstName    string `json:"first_name"`
		LastName     string `json:"last_name"`
	} `json:"data"`
	Type string `json:"type"` // e.g., "user.created"
}
