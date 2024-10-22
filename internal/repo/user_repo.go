package repo

import (
	"database/sql"
	"mlvt/internal/entity"
	"time"
)

type UserRepository interface {
	CreateUser(user *entity.User) error
	GetUserByEmail(email string) (*entity.User, error)
	GetUserByID(userID uint64) (*entity.User, error)
	UpdateUser(user *entity.User) error
	DeleteUser(userID uint64) error
	GetAllUsers() ([]entity.User, error)
	UpdateUserPassword(userID uint64, hashedPassword string) error
	UpdateUserAvatar(userID uint64, avatarPath, avatarFolder string) error
	GetUsersByEmailSuffix(suffix string) ([]entity.User, error)
}

type userRepo struct {
	db *sql.DB
}

func NewUserRepo(db *sql.DB) UserRepository {
	return &userRepo{db: db}
}

// CreateUser inserts a new user into the database
func (r *userRepo) CreateUser(user *entity.User) error {
	query := `
		INSERT INTO users (first_name, last_name, username, email, password, status, premium, role, avatar, avatar_folder, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	_, err := r.db.Exec(query, user.FirstName, user.LastName, user.UserName, user.Email, user.Password, user.Status,
		user.Premium, user.Role, user.Avatar, user.AvatarFolder, user.CreatedAt, user.UpdatedAt)
	return err
}

// GetUserByEmail retrieves a user by their email address
func (r *userRepo) GetUserByEmail(email string) (*entity.User, error) {
	query := `SELECT id, first_name, last_name, username, email, password, status, premium, role, avatar, avatar_folder, created_at, updated_at
	          FROM users WHERE email = ?`
	row := r.db.QueryRow(query, email)

	user := &entity.User{}
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.UserName, &user.Email, &user.Password,
		&user.Status, &user.Premium, &user.Role, &user.Avatar, &user.AvatarFolder, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// GetUserByID retrieves a user by their ID
func (r *userRepo) GetUserByID(userID uint64) (*entity.User, error) {
	query := `SELECT id, first_name, last_name, username, email, password, status, premium, role, avatar, avatar_folder, created_at, updated_at
	          FROM users WHERE id = ?`
	row := r.db.QueryRow(query, userID)

	user := &entity.User{}
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.UserName, &user.Email, &user.Password,
		&user.Status, &user.Premium, &user.Role, &user.Avatar, &user.AvatarFolder, &user.CreatedAt, &user.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return user, err
}

// UpdateUser updates user information
func (r *userRepo) UpdateUser(user *entity.User) error {
	query := `
		UPDATE users
		SET first_name = ?, last_name = ?, username = ?, email = ?, status = ?, premium = ?, role = ?, updated_at = ?
		WHERE id = ?`
	_, err := r.db.Exec(query, user.FirstName, user.LastName, user.UserName, user.Email, user.Status, user.Premium, user.Role, user.UpdatedAt, user.ID)
	return err
}

// DeleteUser performs a soft delete by updating the status of a user to "deleted"
func (r *userRepo) DeleteUser(userID uint64) error {
	query := `UPDATE users SET status = ? WHERE id = ?`
	_, err := r.db.Exec(query, entity.UserStatusDeleted, userID)
	return err
}

// UpdateUserPassword updates the hashed password for a user
func (r *userRepo) UpdateUserPassword(userID uint64, hashedPassword string) error {
	query := `UPDATE users SET password = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, hashedPassword, time.Now(), userID)
	return err
}

// UpdateUserAvatar updates the user's avatar and avatar folder
func (r *userRepo) UpdateUserAvatar(userID uint64, avatarPath, avatarFolder string) error {
	query := `UPDATE users SET avatar = ?, avatar_folder = ?, updated_at = ? WHERE id = ?`
	_, err := r.db.Exec(query, avatarPath, avatarFolder, time.Now(), userID)
	return err
}

// GetAllUsers retrieves all users
func (r *userRepo) GetAllUsers() ([]entity.User, error) {
	query := `SELECT id, first_name, last_name, username, email, password, status, premium, role, avatar, avatar_folder, created_at, updated_at
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
			&user.Status, &user.Premium, &user.Role, &user.Avatar, &user.AvatarFolder, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *userRepo) GetUsersByEmailSuffix(suffix string) ([]entity.User, error) {
	query := `SELECT id, first_name, last_name, username, email, password, status, premium, role, avatar, avatar_folder, created_at, updated_at FROM users WHERE email LIKE ? AND deleted_at IS NULL`
	likePattern := "%" + suffix
	rows, err := r.db.Query(query, likePattern)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []entity.User
	for rows.Next() {
		var user entity.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.UserName, &user.Email, &user.Password, &user.Status, &user.Premium, &user.Role, &user.Avatar, &user.AvatarFolder, &user.CreatedAt, &user.UpdatedAt); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}
