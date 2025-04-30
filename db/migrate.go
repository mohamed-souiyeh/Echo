package db

import (
	"database/sql"
	"embed"
	"errors"
	"net/http"

	"github.com/charmbracelet/log"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/golang-migrate/migrate/v4/source/httpfs"
	_ "modernc.org/sqlite"
)

//go:embed migrations
var migrations embed.FS

func RunMigration(db *sql.DB) {

	sourceInstance, err := httpfs.New(http.FS(migrations), "migrations")
	if err != nil {
		log.Fatalf("Failed to create httpfs source instance: %v", err)
	}

	defer func() {
		sourceErr := sourceInstance.Close()
		if sourceErr != nil {
			log.Warnf("Error closing migrate source: %v", sourceErr)
		}
	}()

	targetInstance, err := sqlite.WithInstance(db, new(sqlite.Config))
	if err != nil {
		log.Fatalf("Failed to create sqlite database instance for migrate: %v", err)
	}

	m, err := migrate.NewWithInstance("httpfs", sourceInstance, "sqlite", targetInstance)
	if err != nil {
		log.Fatalf("Failed to create migrate instance: %v", err)
	}


	// --- Run Migrations ---
	// Choose ONE of the following migration commands:

	// Option 1: Apply all pending "up" migrations (most common)
	log.Info("Applying all pending UP migrations...")
	err = m.Up()

	// Option 2: Roll back the last migration
	// log.Println("Rolling back the last migration...")
	// err = m.Down()

	// Option 3: Apply a specific number of migrations
	// log.Println("Applying next 1 migration...")
	// err = m.Steps(1)

	// Option 4: Roll back a specific number of migrations
	// log.Println("Rolling back 1 migration...")
	// err = m.Steps(-1)

	// Option 5: Migrate to a specific version
	// targetVersion := uint(1) // Version number from filename (e.g., 001_...)
	// log.Printf("Migrating to specific version %d...\n", targetVersion)
	// err = m.Migrate(targetVersion)

	if err != nil {
		// Check specifically for ErrNoChange, which is not a failure
		if errors.Is(err, migrate.ErrNoChange) {
			log.Info("No schema changes needed. Database is up-to-date.")
		} else {
			// Use Fatalf for actual migration errors, as the schema might be inconsistent
			log.Fatalf("Migration failed: %v", err)
		}
	} else {
		log.Info("Migrations applied successfully!")
	}
}
