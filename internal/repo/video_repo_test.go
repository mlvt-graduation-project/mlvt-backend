package repo

import (
	"fmt"
	"testing"

	"mlvt/internal/entity"

	"github.com/DATA-DOG/go-sqlmock"
)

func TestUpdateVideoStatus_Success(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	videoRepo := NewVideoRepo(db)

	videoID := uint64(1)
	newStatus := entity.StatusProcessing

	// Define the expected query and result
	query := `
		UPDATE videos
		SET status = \?, updated_at = \?
		WHERE id = \?`

	mock.ExpectExec(query).
		WithArgs(newStatus, sqlmock.AnyArg(), videoID).
		WillReturnResult(sqlmock.NewResult(0, 1)) // 1 row affected

	// Call the method
	err = videoRepo.UpdateVideoStatus(videoID, newStatus)
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestUpdateVideoStatus_NoRows(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	videoRepo := NewVideoRepo(db)

	videoID := uint64(2)
	newStatus := entity.StatusFailed

	// Define the expected query and result
	query := `
		UPDATE videos
		SET status = \?, updated_at = \?
		WHERE id = \?`

	mock.ExpectExec(query).
		WithArgs(newStatus, sqlmock.AnyArg(), videoID).
		WillReturnResult(sqlmock.NewResult(0, 0)) // 0 rows affected

	// Call the method
	err = videoRepo.UpdateVideoStatus(videoID, newStatus)
	if err == nil {
		t.Errorf("expected error, got none")
	}
	expectedErr := fmt.Sprintf("no video found with id %d", videoID)
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}

func TestUpdateVideoStatus_QueryError(t *testing.T) {
	// Initialize sqlmock
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock database: %v", err)
	}
	defer db.Close()

	videoRepo := NewVideoRepo(db)

	videoID := uint64(3)
	newStatus := entity.StatusSuccess

	// Define the expected query and result
	query := `
		UPDATE videos
		SET status = \?, updated_at = \?
		WHERE id = \?`

	mock.ExpectExec(query).
		WithArgs(newStatus, sqlmock.AnyArg(), videoID).
		WillReturnError(fmt.Errorf("some database error"))

	// Call the method
	err = videoRepo.UpdateVideoStatus(videoID, newStatus)
	if err == nil {
		t.Errorf("expected error, got none")
	}
	expectedErr := "failed to update video status: some database error"
	if err.Error() != expectedErr {
		t.Errorf("expected error '%s', got '%s'", expectedErr, err.Error())
	}

	// Ensure all expectations were met
	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %v", err)
	}
}
