package v1

import (
	"database/sql"
	"fmt"
	"mlvt/internal/infra/zap-logging/log"
	"strings"
)

// MigrateV1Force attempts to add the new columns. If a duplicate column error occurs, it will
// recreate the table without that column, preserving other columns and data, then re-add the column.
func MigrateV1Force(db *sql.DB) error {
	if err := addColumnForce(db, "videos", "original_video_id", "INTEGER"); err != nil {
		return fmt.Errorf("adding original_video_id to videos: %w", err)
	}
	if err := addColumnForce(db, "videos", "audio_id", "INTEGER"); err != nil {
		return fmt.Errorf("adding audio_id to videos: %w", err)
	}
	if err := addColumnForce(db, "videos", "status", "TEXT"); err != nil {
		return fmt.Errorf("adding status to videos: %w", err)
	}
	// Set status to raw
	if _, err := db.Exec(`UPDATE videos SET status = ?;`, StatusRaw); err != nil {
		return fmt.Errorf("updating videos status: %w", err)
	}

	// Transcriptions
	if err := addColumnForce(db, "transcriptions", "original_transcription_id", "INTEGER"); err != nil {
		return fmt.Errorf("adding original_transcription_id to transcriptions: %w", err)
	}
	if err := addColumnForce(db, "transcriptions", "status", "TEXT"); err != nil {
		return fmt.Errorf("adding status to transcriptions: %w", err)
	}
	if _, err := db.Exec(`UPDATE transcriptions SET status = ?;`, StatusRaw); err != nil {
		return fmt.Errorf("updating transcriptions status: %w", err)
	}

	// Audios
	if err := addColumnForce(db, "audios", "transcription_id", "INTEGER"); err != nil {
		return fmt.Errorf("adding transcription_id to audios: %w", err)
	}
	if err := addColumnForce(db, "audios", "status", "TEXT"); err != nil {
		return fmt.Errorf("adding status to audios: %w", err)
	}
	if _, err := db.Exec(`UPDATE audios SET status = ?;`, StatusRaw); err != nil {
		return fmt.Errorf("updating audios status: %w", err)
	}

	return nil
}

// addColumnForce tries to add a column. If the column already exists,
// it will recreate the table without dropping existing data (except for the conflicting column).
func addColumnForce(db *sql.DB, tableName, columnName, columnType string) error {
	stmt := fmt.Sprintf(`ALTER TABLE %s ADD COLUMN %s %s;`, tableName, columnName, columnType)
	_, err := db.Exec(stmt)
	if err != nil {
		if strings.Contains(err.Error(), "duplicate column name") {
			// The column already exists. Let's recreate the table without losing data.
			log.Warnf("Column %s already exists in table %s. Recreating table to ensure correct schema...", columnName, tableName)
			if err := recreateTableWithoutColumn(db, tableName, columnName); err != nil {
				return fmt.Errorf("recreating table %s: %w", tableName, err)
			}
			// After recreation, try adding the column again.
			_, err = db.Exec(stmt)
			if err != nil {
				return fmt.Errorf("adding column %s after recreation: %w", columnName, err)
			}
		} else {
			return err
		}
	}
	return nil
}

// recreateTableWithoutColumn recreates the given table without the specified column.
// This preserves other columns and data. For SQLite, we must:
// 1. Get the current columns (excluding the one we want to remove).
// 2. Create a temp table with the desired columns.
// 3. Copy data from the old table to the temp table.
// 4. Drop old table.
// 5. Rename temp table to the old table.
func recreateTableWithoutColumn(db *sql.DB, tableName, removeColumn string) error {
	// Get existing columns
	cols, err := getTableColumns(db, tableName)
	if err != nil {
		return fmt.Errorf("getting columns: %w", err)
	}

	// Filter out the column we want to remove
	newCols := []string{}
	for _, c := range cols {
		if c != removeColumn {
			newCols = append(newCols, c)
		}
	}
	if len(newCols) == 0 {
		return fmt.Errorf("no columns left after removing %s from %s", removeColumn, tableName)
	}

	// Start a transaction for safety
	tx, err := db.Begin()
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		if p := recover(); p != nil {
			_ = tx.Rollback()
			panic(p)
		}
	}()

	colDefs, err := getColumnDefinitions(db, tableName, removeColumn)
	if err != nil {
		tx.Rollback()
		return fmt.Errorf("getting column definitions: %w", err)
	}

	tempTableName := tableName + "_temp"
	createStmt := fmt.Sprintf("CREATE TABLE %s (%s);", tempTableName, colDefs)
	if _, err := tx.Exec(createStmt); err != nil {
		tx.Rollback()
		return fmt.Errorf("creating temp table: %w", err)
	}

	// Copy data
	srcCols := strings.Join(newCols, ", ")
	destCols := srcCols
	copyStmt := fmt.Sprintf("INSERT INTO %s (%s) SELECT %s FROM %s;", tempTableName, destCols, srcCols, tableName)
	if _, err := tx.Exec(copyStmt); err != nil {
		tx.Rollback()
		return fmt.Errorf("copying data: %w", err)
	}

	// Drop old table
	dropStmt := fmt.Sprintf("DROP TABLE %s;", tableName)
	if _, err := tx.Exec(dropStmt); err != nil {
		tx.Rollback()
		return fmt.Errorf("dropping old table: %w", err)
	}

	// Rename temp to old
	renameStmt := fmt.Sprintf("ALTER TABLE %s RENAME TO %s;", tempTableName, tableName)
	if _, err := tx.Exec(renameStmt); err != nil {
		tx.Rollback()
		return fmt.Errorf("renaming temp table: %w", err)
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

// getTableColumns returns a list of column names for the given table.
func getTableColumns(db *sql.DB, tableName string) ([]string, error) {
	rows, err := db.Query("PRAGMA table_info(" + tableName + ");")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var cols []string
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dfltValue interface{}
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err != nil {
			return nil, err
		}
		cols = append(cols, name)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}
	return cols, nil
}

// getColumnDefinitions constructs a column definition string (e.g., "col1 TEXT, col2 INTEGER")
// suitable for CREATE TABLE, excluding the removeColumn.
func getColumnDefinitions(db *sql.DB, tableName, removeColumn string) (string, error) {
	rows, err := db.Query("PRAGMA table_info(" + tableName + ");")
	if err != nil {
		return "", err
	}
	defer rows.Close()

	var defs []string
	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull, pk int
		var dfltValue interface{}
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err != nil {
			return "", err
		}
		if name == removeColumn {
			continue
		}

		// Build column definition
		def := fmt.Sprintf("%q %s", name, ctype)
		if notnull == 1 {
			def += " NOT NULL"
		}
		if pk == 1 {
			def += " PRIMARY KEY"
		}
		if dfltValue != nil {
			def += fmt.Sprintf(" DEFAULT %v", dfltValue)
		}

		defs = append(defs, def)
	}
	if err := rows.Err(); err != nil {
		return "", err
	}

	return strings.Join(defs, ", "), nil
}
