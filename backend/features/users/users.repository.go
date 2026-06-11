package users

import (
	"context"
	"truck-request-portal/pkg/database"
)

// UpsertUser creates or updates a user in Supabase.
// We use ON CONFLICT to handle existing users safely.
func UpsertUser(ctx context.Context, user *User) error {
	query := `
		INSERT INTO users (clerk_id, name, email, ops_id, role, is_fte, is_active)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		ON CONFLICT (clerk_id) 
		DO UPDATE SET 
			name = EXCLUDED.name,
			email = EXCLUDED.email,
			ops_id = EXCLUDED.ops_id,
			role = EXCLUDED.role,
			is_fte = EXCLUDED.is_fte,
			is_active = EXCLUDED.is_active;
	`
	_, err := database.DB.Exec(ctx, query,
		user.ClerkID, user.Name, user.Email, user.OpsID, user.Role, user.IsFTE, user.IsActive)

	return err
}
