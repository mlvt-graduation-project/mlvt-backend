package migration

import (
	"database/sql"
	"fmt"
)

// Migration defines a database migration.
type Migration struct {
	ID   int
	Name string
	SQL  string
}

// Migrate runs all pending migrations.
func Migrate(db *sql.DB) error {
	// Ensure the migrations table exists
	if err := ensureMigrationsTable(db); err != nil {
		return err
	}

	// Define migrations
	migrations := []Migration{
		{
			ID:   1,
			Name: "create_users_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS users (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					first_name TEXT NOT NULL,
					last_name TEXT NOT NULL,
					email TEXT NOT NULL UNIQUE,
					password TEXT NOT NULL,
					status INTEGER NOT NULL,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
				);`,
		},
		{
			ID:   2,
			Name: "create_videos_table",
			SQL: `
				CREATE TABLE IF NOT EXISTS videos (
					id INTEGER PRIMARY KEY AUTOINCREMENT,
					user_id INTEGER NOT NULL,
					duration INTEGER NOT NULL,
					link TEXT NOT NULL,
					created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
					FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
				);`,
		},
	}

	// Apply pending migrations
	for _, migration := range migrations {
		if err := applyMigration(db, migration); err != nil {
			return err
		}
	}

	return nil
}

// ensureMigrationsTable ensures that the migrations table exists.
func ensureMigrationsTable(db *sql.DB) error {
	createTableSQL := `
	CREATE TABLE IF NOT EXISTS migrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL UNIQUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);
	`
	_, err := db.Exec(createTableSQL)
	return err
}

// applyMigration applies a single migration if it hasn't been applied yet.
func applyMigration(db *sql.DB, migration Migration) error {
	// Check if the migration has already been applied
	var count int
	err := db.QueryRow("SELECT COUNT(*) FROM migrations WHERE name = ?", migration.Name).Scan(&count)
	if err != nil {
		return err
	}

	if count > 0 {
		// Migration already applied
		fmt.Printf("Migration '%s' already applied\n", migration.Name)
		return nil
	}

	// Apply the migration
	fmt.Printf("Applying migration '%s'\n", migration.Name)
	_, err = db.Exec(migration.SQL)
	if err != nil {
		return err
	}

	// Record the migration as applied
	_, err = db.Exec("INSERT INTO migrations (name) VALUES (?)", migration.Name)
	return err
}
