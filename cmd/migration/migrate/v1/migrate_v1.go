package v1

import (
	"database/sql"
	"fmt"
)

// MigrateV1 adds new columns to the tables and sets the status to "raw".
func MigrateV1(db *sql.DB) error {
	// Add columns to videos
	if _, err := db.Exec(`ALTER TABLE videos ADD COLUMN original_video_id INTEGER;`); err != nil {
		return fmt.Errorf("adding original_video_id to videos: %w", err)
	}
	if _, err := db.Exec(`ALTER TABLE videos ADD COLUMN audio_id INTEGER;`); err != nil {
		return fmt.Errorf("adding audio_id to videos: %w", err)
	}
	if _, err := db.Exec(`ALTER TABLE videos ADD COLUMN status TEXT;`); err != nil {
		return fmt.Errorf("adding status to videos: %w", err)
	}
	// Set status to raw for all existing video rows
	if _, err := db.Exec(`UPDATE videos SET status = ?;`, StatusRaw); err != nil {
		return fmt.Errorf("updating videos status to raw: %w", err)
	}

	// Add columns to transcriptions
	if _, err := db.Exec(`ALTER TABLE transcriptions ADD COLUMN original_transcription_id INTEGER;`); err != nil {
		return fmt.Errorf("adding original_transcription_id to transcriptions: %w", err)
	}
	if _, err := db.Exec(`ALTER TABLE transcriptions ADD COLUMN status TEXT;`); err != nil {
		return fmt.Errorf("adding status to transcriptions: %w", err)
	}
	// Set status to raw for all existing transcription rows
	if _, err := db.Exec(`UPDATE transcriptions SET status = ?;`, StatusRaw); err != nil {
		return fmt.Errorf("updating transcriptions status to raw: %w", err)
	}

	// Add columns to audios
	if _, err := db.Exec(`ALTER TABLE audios ADD COLUMN transcription_id INTEGER;`); err != nil {
		return fmt.Errorf("adding transcription_id to audios: %w", err)
	}
	if _, err := db.Exec(`ALTER TABLE audios ADD COLUMN status TEXT;`); err != nil {
		return fmt.Errorf("adding status to audios: %w", err)
	}
	// Set status to raw for all existing audio rows
	if _, err := db.Exec(`UPDATE audios SET status = ?;`, StatusRaw); err != nil {
		return fmt.Errorf("updating audios status to raw: %w", err)
	}

	return nil
}
