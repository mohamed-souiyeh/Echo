package db

import (
	"context"
	"fmt"

	"github.com/charmbracelet/log"

	repo "echo/db/repository"
)

// TODO: maybe make the seed generate random data and use an argumen tas the amount to generate
func RunUserSeed(ctx context.Context, repo repo.UserRepository) error {
	log.Info("Attempting to seed users...")

	usersToSeed := []struct {
		Username string
		Password string // Plain text password for seeding input
	}{
		{"alice", "password123"},
		{"bob", "securepass"},
		{"charlie", "qwerty"},
	}

	createdCount := 0
	skippedCount := 0

	for _, userData := range usersToSeed {
		// Optional: Check if user exists first (requires GetUserByUsername)
		// _, err := repo.GetUserByUsername(ctx, userData.Username)
		// if err == nil {
		//  log.Printf("User %s already exists, skipping.", userData.Username)
		//  skippedCount++
		//  continue
		// }
		// Since we don't have GetUserByUsername yet, we just try creating.
		// CreateUser should handle potential UNIQUE constraint errors gracefully.

		log.Debugf("Attempting to create user: %s", userData.Username)
		_, err := repo.CreateUser(ctx, userData.Username, userData.Password)
		if err != nil {
			// Check if the error is because the user already exists (e.g., UNIQUE constraint)
			// This check is driver-specific, sqlite often returns 'UNIQUE constraint failed'
			// For now, we just log any error and continue
			log.Warn("WARN: Failed to create user (might already exist): ", "username", userData.Username, "error", err)
			skippedCount++ // Assume failure means skipped/already exists for this simple seed
		} else {
			log.Info("Successfully created the user: ", "user", userData.Username)
			createdCount++
		}
	}

	log.Info("Seeding finished.", "Created", createdCount, "Skipped/Existing", skippedCount)

	// --- Verification Step ---
	log.Info("Verifying seeded data by fetching all users...")
	allUsers, err := repo.GetAllUsers(ctx)
	if err != nil {
		log.Error("ERROR: Failed to fetch users after seeding", "error", err)
		return fmt.Errorf("failed to verify seeding: %w", err)
	}

	log.Info("Found users in the database:", "found", len(allUsers))
	for _, u := range allUsers {
		// Be careful not to log sensitive info like passwords, even hashed ones, in real logs
		log.Debugf("  - ID: %d,\n  - Username: %s,\n  - hashedPassword: %s,\n  - CreatedAt: %s", u.ID, u.Username, u.Password, u.CreatedAt)
	}

	return nil // Seeding (and verification fetch) completed
}
