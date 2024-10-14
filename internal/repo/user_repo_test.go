package repo

import (
	"regexp"
	"testing"
	"time"

	"mlvt/internal/entity"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestCreateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)

	user := &entity.User{
		FirstName:    "John",
		LastName:     "Doe",
		UserName:     "johndoe",
		Email:        "john@example.com",
		Password:     "hashedpassword",
		Status:       entity.UserStatusAvailable,
		Premium:      false,
		Role:         "user",
		Avatar:       "avatar.jpg",
		AvatarFolder: "avatars",
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	// Expect the INSERT query
	mock.ExpectExec(regexp.QuoteMeta(`
		INSERT INTO users (first_name, last_name, username, email, password, status, premium, role, avatar, avatar_folder, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`)).
		WithArgs(user.FirstName, user.LastName, user.UserName, user.Email, user.Password, user.Status,
			user.Premium, user.Role, user.Avatar, user.AvatarFolder, user.CreatedAt, user.UpdatedAt).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = repo.CreateUser(user)
	assert.NoError(t, err)

	// Ensure all expectations were met
	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetUserByEmail(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)
	email := "john@example.com"

	rows := sqlmock.NewRows([]string{
		"id", "first_name", "last_name", "username", "email", "password",
		"status", "premium", "role", "avatar", "avatar_folder", "created_at", "updated_at",
	}).AddRow(
		1, "John", "Doe", "johndoe", email, "hashedpassword",
		entity.UserStatusAvailable, false, "user", "avatar.jpg", "avatars",
		time.Now(), time.Now(),
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, first_name, last_name, username, email, password, status, premium, role, avatar, avatar_folder, created_at, updated_at
		          FROM users WHERE email = ?`)).
		WithArgs(email).
		WillReturnRows(rows)

	user, err := repo.GetUserByEmail(email)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, email, user.Email)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetUserByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)
	userID := uint64(1)

	rows := sqlmock.NewRows([]string{
		"id", "first_name", "last_name", "username", "email", "password",
		"status", "premium", "role", "avatar", "avatar_folder", "created_at", "updated_at",
	}).AddRow(
		userID, "John", "Doe", "johndoe", "john@example.com", "hashedpassword",
		entity.UserStatusAvailable, false, "user", "avatar.jpg", "avatars",
		time.Now(), time.Now(),
	)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, first_name, last_name, username, email, password, status, premium, role, avatar, avatar_folder, created_at, updated_at
		          FROM users WHERE id = ?`)).
		WithArgs(userID).
		WillReturnRows(rows)

	user, err := repo.GetUserByID(userID)
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, userID, user.ID)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)

	user := &entity.User{
		ID:        1,
		FirstName: "Jane",
		LastName:  "Doe",
		UserName:  "janedoe",
		Email:     "jane@example.com",
		Status:    entity.UserStatusAvailable,
		Premium:   true,
		Role:      "admin",
		UpdatedAt: time.Now(),
	}

	mock.ExpectExec(regexp.QuoteMeta(`
		UPDATE users
		SET first_name = ?, last_name = ?, username = ?, email = ?, status = ?, premium = ?, role = ?, updated_at = ?
		WHERE id = ?`)).
		WithArgs(user.FirstName, user.LastName, user.UserName, user.Email, user.Status, user.Premium, user.Role, user.UpdatedAt, user.ID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateUser(user)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestDeleteUser(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)
	userID := uint64(1)

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET status = ? WHERE id = ?`)).
		WithArgs(entity.UserStatusDeleted, userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.DeleteUser(userID)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateUserPassword(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)
	userID := uint64(1)
	hashedPassword := "newhashedpassword"

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET password = ?, updated_at = ? WHERE id = ?`)).
		WithArgs(hashedPassword, sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateUserPassword(userID, hashedPassword)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestUpdateUserAvatar(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)
	userID := uint64(1)
	avatarPath := "avatar_new.jpg"
	avatarFolder := "avatars_new"

	mock.ExpectExec(regexp.QuoteMeta(`UPDATE users SET avatar = ?, avatar_folder = ?, updated_at = ? WHERE id = ?`)).
		WithArgs(avatarPath, avatarFolder, sqlmock.AnyArg(), userID).
		WillReturnResult(sqlmock.NewResult(0, 1))

	err = repo.UpdateUserAvatar(userID, avatarPath, avatarFolder)
	assert.NoError(t, err)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}

func TestGetAllUsers(t *testing.T) {
	db, mock, err := sqlmock.New()
	assert.NoError(t, err)
	defer db.Close()

	repo := NewUserRepo(db)

	rows := sqlmock.NewRows([]string{
		"id", "first_name", "last_name", "username", "email", "password",
		"status", "premium", "role", "avatar", "avatar_folder", "created_at", "updated_at",
	}).
		AddRow(
			1, "John", "Doe", "johndoe", "john@example.com", "hashedpassword",
			entity.UserStatusAvailable, false, "user", "avatar.jpg", "avatars",
			time.Now(), time.Now(),
		).
		AddRow(
			2, "Jane", "Smith", "janesmith", "jane@example.com", "hashedpassword2",
			entity.UserStatusAvailable, true, "admin", "avatar2.jpg", "avatars",
			time.Now(), time.Now(),
		)

	mock.ExpectQuery(regexp.QuoteMeta(`
		SELECT id, first_name, last_name, username, email, password, status, premium, role, avatar, avatar_folder, created_at, updated_at
		          FROM users`)).
		WillReturnRows(rows)

	users, err := repo.GetAllUsers()
	assert.NoError(t, err)
	assert.Len(t, users, 2)
	assert.Equal(t, "john@example.com", users[0].Email)
	assert.Equal(t, "jane@example.com", users[1].Email)

	err = mock.ExpectationsWereMet()
	assert.NoError(t, err)
}
