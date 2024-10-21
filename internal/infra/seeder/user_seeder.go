package seeder

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"
	"mime"
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

// UserSeeder handles the seeding of users from a folder of images.
type UserSeeder struct {
	userRepo repo.UserRepository
	s3Client aws.S3ClientInterface
}

// NewUserSeeder initializes a new UserSeeder.
func NewUserSeeder(userRepo repo.UserRepository, s3Client aws.S3ClientInterface) *UserSeeder {
	return &UserSeeder{
		userRepo: userRepo,
		s3Client: s3Client,
	}
}

// SeedUsersFromFolder reads images from the specified folder, creates users, uploads avatars to S3, and updates user records.
func (s *UserSeeder) SeedUsersFromFolder(imagesFolder string) error {
	files, err := ioutil.ReadDir(imagesFolder)
	if err != nil {
		return fmt.Errorf("failed to read images folder: %v", err)
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
		filePath := filepath.Join(imagesFolder, file.Name())
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

		// Upload the file directly to S3
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

// generateRandomPassword creates a random password of the given length.
func generateRandomPassword(length int) string {
	bytes := make([]byte, length)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)[:length]
}
