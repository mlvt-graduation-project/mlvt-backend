package entity

import "time"

// UserStatus constants
const (
	UserStatusAvailable = 1
	UserStatusSuspended = 9
	UserStatusDeleted   = 10
)

type UserStatus string
type UserPermission string

const (
	UserStatusActive   UserStatus = "active"
	UserStatusInactive UserStatus = "inactive"
    UserStatusPending UserStatus = "pending"
	AdminRole UserPermission = "admin"
	UserRole  UserPermission = "user"
)

// User represents the schema for user data
type User struct {
	ID            uint64         `json:"id"`         // Unique identifier for the user
	FirstName     string         `json:"first_name"` // User's first name
	LastName      string         `json:"last_name"`  // User's last name\
	UserName      string         `json:"username"`
	Email         string         `json:"email"`         // User's email address
	Password      string         `json:"password"`      // User's hashed password
	Status        UserStatus     `json:"status"`        // Status of the user (available, suspended, deleted)
	Role          UserPermission `json:"role"`          // Role of the user (User, Admin, etc.)
	Avatar        string         `json:"avatar"`        // file name
	AvatarFolder  string         `json:"avatar_folder"` // Folder that contain the avatar image on s3
	WalletBalance int64          `json:"wallet_balance"`
	CreatedAt     time.Time      `json:"created_at"` // Timestamp of when the user was created
	UpdatedAt     time.Time      `json:"updated_at"` // Timestamp of the last update to the user's data
}
