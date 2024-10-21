// internal/infra/seeder/user_video_seeder.go

package seeder

import (
	"crypto/md5"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"mlvt/internal/entity"
	"mlvt/internal/infra/aws"
	"mlvt/internal/infra/env"
	"mlvt/internal/repo"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// UserVideoSeeder handles the seeding of users and videos from folders of images and videos.
type UserVideoSeeder struct {
	userRepo  repo.UserRepository
	videoRepo repo.VideoRepository
	s3Client  aws.S3ClientInterface
}

// NewUserVideoSeeder initializes a new UserVideoSeeder.
func NewUserVideoSeeder(userRepo repo.UserRepository, videoRepo repo.VideoRepository, s3Client aws.S3ClientInterface) *UserVideoSeeder {
	return &UserVideoSeeder{
		userRepo:  userRepo,
		videoRepo: videoRepo,
		s3Client:  s3Client,
	}
}

// SeedUsersAndVideosFromFolders seeds users from the avatars folder and videos from the videos folder.
func (s *UserVideoSeeder) SeedUsersAndVideosFromFolders(avatarsFolder, videosFolder string) error {
	// Seed Users
	if err := s.SeedUsersFromFolder(avatarsFolder); err != nil {
		return fmt.Errorf("failed to seed users: %v", err)
	}

	// Seed Videos
	if err := s.SeedVideosFromFolder(videosFolder); err != nil {
		return fmt.Errorf("failed to seed videos: %v", err)
	}

	return nil
}

// SeedUsersFromFolder reads images from the avatars folder, creates users, uploads avatars to S3, and updates user records.
func (s *UserVideoSeeder) SeedUsersFromFolder(avatarsFolder string) error {
	files, err := ioutil.ReadDir(avatarsFolder)
	if err != nil {
		return fmt.Errorf("failed to read avatars folder: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories
		}

		// Check if the file is an image based on extension
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ext != ".jpg" && ext != ".jpeg" && ext != ".png" && ext != ".gif" {
			continue // Skip non-image files
		}

		// Generate unique identifier for the user
		uniqueID := uuid.New().String()

		// Generate user data
		firstName := strings.Title(strings.TrimSuffix(file.Name(), ext))
		lastName := "Seeder"
		username := fmt.Sprintf("%s_%s", strings.ToLower(firstName), "seeder")
		email := fmt.Sprintf("%s@seeder.com", username)
		password := generateRandomPassword(12) // Generate a random password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			log.Printf("Failed to hash password for user %s: %v", username, err)
			continue
		}

		// Check if user already exists
		existingUser, err := s.userRepo.GetUserByEmail(email)
		if err != nil {
			log.Printf("Failed to check existing user for email %s: %v", email, err)
			continue
		}
		if existingUser != nil {
			log.Printf("User with email %s already exists, skipping", email)
			continue
		}

		// Create user entity
		user := &entity.User{
			FirstName:    firstName,
			LastName:     lastName,
			UserName:     username,
			Email:        email,
			Password:     string(hashedPassword),
			Status:       entity.UserStatusAvailable,
			Premium:      false,
			Role:         "User",
			Avatar:       "", // To be updated after upload
			AvatarFolder: "", // To be updated after upload
			CreatedAt:    time.Now(),
			UpdatedAt:    time.Now(),
		}

		// Insert user into the database
		err = s.userRepo.CreateUser(user)
		if err != nil {
			log.Printf("Failed to create user %s: %v", username, err)
			continue
		}

		// Upload avatar to S3
		avatarFolder := env.EnvConfig.AvatarFolder // Use configured avatar folder
		filePath := filepath.Join(avatarsFolder, file.Name())
		fileType := mime.TypeByExtension(ext)
		if fileType == "" {
			fileType = "image/jpeg" // Default to JPEG if unknown
		}

		// Generate a unique filename to prevent collisions
		uniqueFileName := fmt.Sprintf("%s_%s", uniqueID, file.Name())

		// Read the file data
		fileData, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to read file %s: %v", filePath, err)
			continue
		}

		// Upload the avatar directly to S3
		err = s.s3Client.UploadFile(avatarFolder, uniqueFileName, fileType, fileData)
		if err != nil {
			log.Printf("Failed to upload avatar for user %s: %v", username, err)
			continue
		}

		// Update user with avatar information
		user.Avatar = uniqueFileName
		user.AvatarFolder = avatarFolder
		user.UpdatedAt = time.Now()

		err = s.userRepo.UpdateUserAvatar(user.ID, uniqueFileName, avatarFolder)
		if err != nil {
			log.Printf("Failed to update avatar info for user %s: %v", username, err)
			continue
		}

		log.Printf("Successfully created user %s with avatar %s", username, uniqueFileName)
	}

	return nil
}

// SeedVideosFromFolder reads videos and their corresponding frames from the videos folder,
// creates video records, uploads videos and frames to S3, and updates video records.
func (s *UserVideoSeeder) SeedVideosFromFolder(videosFolder string) error {
	files, err := ioutil.ReadDir(videosFolder)
	if err != nil {
		return fmt.Errorf("failed to read videos folder: %v", err)
	}

	for _, file := range files {
		if file.IsDir() {
			continue // Skip directories
		}

		// Check if the file is a video based on extension
		ext := strings.ToLower(filepath.Ext(file.Name()))
		if ext != ".mp4" && ext != ".avi" && ext != ".mov" && ext != ".mkv" {
			continue // Skip non-video files
		}

		// Generate unique identifier for the video
		uniqueID := uuid.New().String()

		// Generate video data
		title := strings.Title(strings.TrimSuffix(file.Name(), ext))
		duration := 120 // Placeholder duration in seconds; replace with actual duration if available
		description := "Seeder video description"
		fileName := file.Name()
		folder := env.EnvConfig.VideosFolder
		image := fmt.Sprintf("%s_thumbnail.jpg", strings.TrimSuffix(file.Name(), ext)) // Thumbnail image filename

		// Path to the video file
		filePath := filepath.Join(videosFolder, fileName)

		// Path to the corresponding frame image
		framePath := filepath.Join(videosFolder, image)

		// Check if the frame file exists
		if _, err := os.Stat(framePath); os.IsNotExist(err) {
			log.Printf("Frame image %s does not exist for video %s, skipping video", image, fileName)
			continue
		}

		// Assign the video to user ID 1 or 2 randomly
		userID := s.assignToUser()

		// Check if video already exists by unique ID
		existingVideo, err := s.videoRepo.GetVideoByID(s.uniqueIDToUint64(uniqueID))
		if err != nil {
			log.Printf("Failed to check existing video for ID %s: %v", uniqueID, err)
			continue
		}
		if existingVideo != nil {
			log.Printf("Video with ID %s already exists, skipping", uniqueID)
			continue
		}

		// Create video entity
		video := &entity.Video{
			Title:       title,
			Duration:    duration,
			Description: description,
			FileName:    fileName,
			Folder:      folder,
			Image:       image,
			Status:      entity.StatusRaw,
			UserID:      userID,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		}

		// Insert video into the database
		err = s.videoRepo.CreateVideo(video)
		if err != nil {
			log.Printf("Failed to create video %s: %v", title, err)
			continue
		}

		// Upload video to S3
		fileType := mime.TypeByExtension(ext)
		if fileType == "" {
			fileType = "video/mp4" // Default to MP4 if unknown
		}

		// Generate a unique filename to prevent collisions
		uniqueFileName := fmt.Sprintf("%s_%s", uniqueID, fileName)

		// Read the video file data
		videoData, err := ioutil.ReadFile(filePath)
		if err != nil {
			log.Printf("Failed to read video file %s: %v", filePath, err)
			continue
		}

		// Upload the video directly to S3
		err = s.s3Client.UploadFile(folder, uniqueFileName, fileType, videoData)
		if err != nil {
			log.Printf("Failed to upload video %s to S3: %v", uniqueFileName, err)
			continue
		}

		// Upload the corresponding frame image to S3
		frameExt := strings.ToLower(filepath.Ext(image))
		frameType := mime.TypeByExtension(frameExt)
		if frameType == "" {
			frameType = "image/jpeg" // Default to JPEG if unknown
		}

		// Read the frame image data
		frameData, err := ioutil.ReadFile(framePath)
		if err != nil {
			log.Printf("Failed to read frame file %s: %v", framePath, err)
			continue
		}

		// Generate a unique filename for the frame to prevent collisions
		uniqueFrameName := fmt.Sprintf("%s_%s", uniqueID, image)

		// Upload the frame image directly to S3
		err = s.s3Client.UploadFile(env.EnvConfig.VideoFramesFolder, uniqueFrameName, frameType, frameData)
		if err != nil {
			log.Printf("Failed to upload frame %s to S3: %v", uniqueFrameName, err)
			continue
		}

		// Update video record with uploaded file information
		video.FileName = uniqueFileName
		video.Image = uniqueFrameName
		video.UpdatedAt = time.Now()

		err = s.videoRepo.UpdateVideo(video)
		if err != nil {
			log.Printf("Failed to update video record for %s: %v", title, err)
			continue
		}

		log.Printf("Successfully created video %s assigned to user ID %d with file %s and frame %s", title, userID, uniqueFileName, uniqueFrameName)
	}

	return nil
}

// assignToUser randomly assigns a video to user ID 1 or 2.
func (s *UserVideoSeeder) assignToUser() uint64 {
	// Generate a random number between 0 and 1
	b := make([]byte, 1)
	rand.Read(b)
	if b[0]%2 == 0 {
		return 1
	}
	return 2
}

// selectRandomUser selects a random user from the database to associate with a video.
// You can modify this function to implement different selection strategies.
// func (s *UserVideoSeeder) selectRandomUser() (*entity.User, error) {
// 	users, err := s.userRepo.GetAllUsers()
// 	if err != nil {
// 		return nil, err
// 	}
// 	if len(users) == 0 {
// 		return nil, nil
// 	}
// 	// Simple random selection
// 	index := randInt(0, len(users))
// 	return &users[index], nil
// }

// uniqueIDToUint64 converts a unique ID string to uint64.
// This is a placeholder; ensure that your unique ID generation aligns with your database's ID type.
func (s *UserVideoSeeder) uniqueIDToUint64(id string) uint64 {
	// Example conversion using MD5 hash
	// Note: This is simplistic and not collision-resistant. Adjust as needed.
	hash := md5.Sum([]byte(id))
	return uint64(hash[0]) // Simplistic; replace with proper mapping if necessary
}

// generateRandomPassword creates a random password of the given length.
func (s *UserVideoSeeder) generateRandomPassword(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}
