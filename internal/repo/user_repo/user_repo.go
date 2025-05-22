package user_repo

import (
	"database/sql"
	"fmt"
	"mlvt/internal/entity"
	"time"

	"github.com/jmoiron/sqlx"
)

type UserRepository interface {
	CreateUser(user *entity.User) error
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(userID uint64) (*entity.User, error)
	GetUserByCondition(user *entity.User) (*entity.User, error)
	UpdateUser(user *entity.User) error
	SoftDeleteUser(userID uint64) error
	DeleteUser(userID uint64) error
	GetAllUsers() ([]entity.User, error)
	UpdateUserPassword(userID uint64, hashedPassword string) error
	UpdateUserAvatar(userID uint64, avatarPath, avatarFolder string) error
	GetUsersByEmailSuffix(suffix string) ([]entity.User, error)
}

type userRepo struct {
	db *sqlx.DB
}

func NewUserRepo(db *sqlx.DB) UserRepository {
	return &userRepo{db: db}
}

// CreateUser inserts a new user into the database
func (r *userRepo) CreateUser(user *entity.User) error {
	query := `
		INSERT INTO users (first_name, last_name, username, email, password, status, role, avatar, avatar_folder, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`
	_, err := r.db.Exec(query, user.FirstName, user.LastName, user.UserName, user.Email, user.Password, user.Status,
		user.Role, user.Avatar, user.AvatarFolder, user.CreatedAt, user.UpdatedAt)
	return err
}

// GetUserByEmail retrieves a user by their email address
func (r *userRepo) GetUserByEmail(email string) (*entity.User, error) {
	query := `SELECT id, first_name, last_name, username, email, password, status, role, avatar, avatar_folder, created_at, updated_at
	          FROM users WHERE email = $1`
	row := r.db.QueryRow(query, email)

	user := &entity.User{}
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.UserName, &user.Email, &user.Password,
		&user.Status, &user.Role, &user.Avatar, &user.AvatarFolder, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// GetUserByUserName retrieves a user by their username
func (r *userRepo) GetUserByCondition(user *entity.User) (*entity.User, error) {
	if user == nil {
		return nil, fmt.Errorf("user condition is nil")
	}

	var args []interface{}
	query := `SELECT id, first_name, last_name, username, email, password, status, role, avatar, avatar_folder, created_at, updated_at
	          FROM users WHERE 1=1`

	if user.UserName != "" {
		query += " AND username = ?"
		args = append(args, user.UserName)
	}
	if user.Role != "" {
		query += " AND role = ?"
		args = append(args, user.Role)
	}
	if user.Email != "" {
		query += " AND email = ?"
		args = append(args, user.Email)
	}
	if user.Status != "" {
		query += " AND status = ?"
		args = append(args, user.Status)
	}

	query = sqlx.Rebind(sqlx.DOLLAR, query)

	row := r.db.QueryRow(query, args...)
	result := &entity.User{}
	err := row.Scan(&result.ID, &result.FirstName, &result.LastName, &result.UserName, &result.Email, &result.Password,
		&result.Status, &result.Role, &result.Avatar, &result.AvatarFolder, &result.CreatedAt, &result.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("db error: %w", err)
	}
	return result, nil
}

// GetUserByID retrieves a user by their ID
func (r *userRepo) GetUserByID(userID uint64) (*entity.User, error) {
	query := `SELECT id, first_name, last_name, username, email, password, status, role, avatar, avatar_folder, created_at, updated_at
	          FROM users WHERE id = $1`
	row := r.db.QueryRow(query, userID)

	user := &entity.User{}
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.UserName, &user.Email, &user.Password,
		&user.Status, &user.Role, &user.Avatar, &user.AvatarFolder, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// UpdateUser updates user information
func (r *userRepo) UpdateUser(user *entity.User) error {
	query := `
		UPDATE users
		SET first_name = $1, last_name = $2, username = $3, email = $4, status = $5, role = $6, updated_at = $7
		WHERE id = $8`
	_, err := r.db.Exec(query, user.FirstName, user.LastName, user.UserName, user.Email, user.Status, user.Role, user.UpdatedAt, user.ID)
	return err
}

// SoftDeleteUser performs a soft delete by updating the status of a user to "deleted"
func (r *userRepo) SoftDeleteUser(userID uint64) error {
	query := `UPDATE users SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, entity.UserStatusDeleted, userID)
	return err
}

func (r *userRepo) DeleteUser(userID uint64) error {
	query := "DELETE FROM users WHERE id = $1"
	_, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete user with ID %d: %v", userID, err)
	}
	return nil
}

// UpdateUserPassword updates the hashed password for a user
func (r *userRepo) UpdateUserPassword(userID uint64, hashedPassword string) error {
	query := `UPDATE users SET password = $1, updated_at = $2 WHERE id = $3`
	_, err := r.db.Exec(query, hashedPassword, time.Now(), userID)
	return err
}

// UpdateUserAvatar updates the user's avatar and avatar folder
func (r *userRepo) UpdateUserAvatar(userID uint64, avatarPath, avatarFolder string) error {
	query := `UPDATE users SET avatar = $1, avatar_folder = $2, updated_at = $3 WHERE id = $4`
	_, err := r.db.Exec(query, avatarPath, avatarFolder, time.Now(), userID)
	return err
}

// GetAllUsers retrieves all users
func (r *userRepo) GetAllUsers() ([]entity.User, error) {
	query := `SELECT id, first_name, last_name, username, email, password, status, role, avatar, avatar_folder, created_at, updated_at
	          FROM users`
	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.UserName, &user.Email, &user.Password,
			&user.Status, &user.Role, &user.Avatar, &user.AvatarFolder, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepo) GetUsersByEmailSuffix(suffix string) ([]entity.User, error) {
	query := `SELECT id, first_name, last_name, username, email, password, status, role, avatar, avatar_folder, created_at, updated_at FROM users WHERE email LIKE $1` // AND deleted_at IS NULL`
	likePattern := "%" + suffix
	rows, err := r.db.Query(query, likePattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.UserName, &user.Email, &user.Password, &user.Status, &user.Role, &user.Avatar, &user.AvatarFolder, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
