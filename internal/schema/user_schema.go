package schema

import "time"

// RegisterUserRequest is used for registering a new user
type RegisterUserRequest struct {
	FirstName string `json:"first_name" validate:"required"`
	LastName  string `json:"last_name" validate:"required"`
	Email     string `json:"email" validate:"required,email"`
	Password  string `json:"password" validate:"required,gte=8"`
}

// LoginUserRequest is used for user login
type LoginUserRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

// GetUserResponse represents the response structure for fetching user details
type GetUserResponse struct {
	ID        uint64    `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Status    int       `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// LoginResponse represents the response structure after successful login
type LoginResponse struct {
	ID        uint64    `json:"id"`
	Email     string    `json:"email"`
	Token     string    `json:"token"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
