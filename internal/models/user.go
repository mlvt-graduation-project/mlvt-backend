package models

import "time"

type User struct {
	ID        int       `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreateAt  time.Time `json:"create_at"`
	UpdatedAt time.Time `json:"update_at"`
	Videos    []Video   `json:"videos"` // List of videos uploaded by the user
}
