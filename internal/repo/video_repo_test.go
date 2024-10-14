package repo

import (
	"database/sql"
	"mlvt/internal/entity"
	"testing"
	"time"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func setupTestDB() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", ":memory:")
	if err != nil {
		return nil, err
	}

	createTableQuery := `
	CREATE TABLE videos (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		duration INTEGER,
		description TEXT,
		file_name TEXT,
		folder TEXT,
		image TEXT,
		status TEXT,
		user_id INTEGER,
		created_at DATETIME,
		updated_at DATETIME
	);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func TestCreateVideo(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)
	defer db.Close()

	videoRepo := NewVideoRepo(db)

	video := &entity.Video{
		Title:       "Test Video",
		Duration:    120,
		Description: "Test Description",
		FileName:    "test.mp4",
		Folder:      "test_folder",
		Image:       "test_image.jpg",
		Status:      entity.StatusRaw,
		UserID:      1,
	}

	err = videoRepo.CreateVideo(video)
	assert.NoError(t, err)

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM videos WHERE title = ?", video.Title).Scan(&count)
	assert.NoError(t, err)
	assert.Equal(t, 1, count)
}

func TestGetVideoByID(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)
	defer db.Close()

	videoRepo := NewVideoRepo(db)

	video := &entity.Video{
		Title:       "Test Video",
		Duration:    120,
		Description: "Test Description",
		FileName:    "test.mp4",
		Folder:      "test_folder",
		Image:       "test_image.jpg",
		Status:      entity.StatusRaw,
		UserID:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = videoRepo.CreateVideo(video)
	assert.NoError(t, err)

	result, err := videoRepo.GetVideoByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, video.Title, result.Title)
}

func TestGetVideoByIDNotFound(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)
	defer db.Close()

	videoRepo := NewVideoRepo(db)

	result, err := videoRepo.GetVideoByID(1)
	assert.NoError(t, err)
	assert.Nil(t, result)
}

func TestListVideosByUserID(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)
	defer db.Close()

	videoRepo := NewVideoRepo(db)

	video1 := &entity.Video{
		Title:       "Test Video 1",
		Duration:    120,
		Description: "Test Description 1",
		FileName:    "test1.mp4",
		Folder:      "test_folder_1",
		Image:       "test_image_1.jpg",
		Status:      entity.StatusRaw,
		UserID:      1,
	}
	video2 := &entity.Video{
		Title:       "Test Video 2",
		Duration:    150,
		Description: "Test Description 2",
		FileName:    "test2.mp4",
		Folder:      "test_folder_2",
		Image:       "test_image_2.jpg",
		Status:      entity.StatusProcessing,
		UserID:      1,
	}

	err = videoRepo.CreateVideo(video1)
	assert.NoError(t, err)
	err = videoRepo.CreateVideo(video2)
	assert.NoError(t, err)

	result, err := videoRepo.ListVideosByUserID(1)
	assert.NoError(t, err)
	assert.Len(t, result, 2)
}

func TestUpdateVideo(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)
	defer db.Close()

	videoRepo := NewVideoRepo(db)

	video := &entity.Video{
		Title:       "Test Video",
		Duration:    120,
		Description: "Test Description",
		FileName:    "test.mp4",
		Folder:      "test_folder",
		Image:       "test_image.jpg",
		Status:      entity.StatusRaw,
		UserID:      1,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = videoRepo.CreateVideo(video)
	assert.NoError(t, err)

	// Retrieve the video to get the assigned ID
	savedVideo, err := videoRepo.GetVideoByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, savedVideo)

	savedVideo.Title = "Updated Test Video"
	savedVideo.UpdatedAt = time.Now()
	err = videoRepo.UpdateVideo(savedVideo)
	assert.NoError(t, err)

	updatedVideo, err := videoRepo.GetVideoByID(savedVideo.ID)
	assert.NoError(t, err)
	assert.NotNil(t, updatedVideo)
	assert.Equal(t, "Updated Test Video", updatedVideo.Title)
	assert.WithinDuration(t, savedVideo.UpdatedAt, updatedVideo.UpdatedAt, time.Second)
}

func TestDeleteVideo(t *testing.T) {
	db, err := setupTestDB()
	assert.NoError(t, err)
	defer db.Close()

	videoRepo := NewVideoRepo(db)

	video := &entity.Video{
		Title:       "Test Video",
		Duration:    120,
		Description: "Test Description",
		FileName:    "test.mp4",
		Folder:      "test_folder",
		Image:       "test_image.jpg",
		Status:      entity.StatusRaw,
		UserID:      1,
	}

	err = videoRepo.CreateVideo(video)
	assert.NoError(t, err)

	err = videoRepo.DeleteVideo(1)
	assert.NoError(t, err)

	deletedVideo, err := videoRepo.GetVideoByID(1)
	assert.NoError(t, err)
	assert.Nil(t, deletedVideo)
}
