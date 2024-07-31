package repository

import (
	"database/sql"
	"errors"
	"mlvt/internal/models"
)

type UserRepository interface {
	GetByEmail(email string) (*models.User, error)
	GetByID(id uint64) (*models.User, error)
	Create(user *models.User) error
	Update(user *models.User) error
	Delete(user *models.User) error
}

type PostgresUserRepository struct {
	db *sql.DB
}

func NewPostgresUserRepository(db *sql.DB) *PostgresUserRepository {
	return &PostgresUserRepository{
		db: db,
	}
}

func (repo *PostgresUserRepository) GetByEmail(email string) (*models.User, error) {
	var user models.User
	query := `
	SELECT id, first_name, last_name, email, password, create_at, update_at 
	FROM users 
	WHERE email = $1
    `
	row := repo.db.QueryRow(query, email)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreateAt, &user.UpdatedAt)
	if err != nil {
		// Cannot found any user with this Email
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repo *PostgresUserRepository) GetByID(id uint64) (*models.User, error) {
	var user models.User
	query := `
	SELECT id, first_name, last_name, email, password, create_at, update_at 
	FROM users 
	WHERE id = $1
	`
	row := repo.db.QueryRow(query, id)
	err := row.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.CreateAt, &user.UpdatedAt)

	if err != nil {
		// Cannot found any user with this ID
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (repo *PostgresUserRepository) Create(user *models.User) error {
	// Add user into database
	query := `
	INSERT INTO users (first_name, last_name, email, password, create_at, update_at) 
	VALUES ($1, $2, $3, $4, $5, $6) 
    `
	_, err := repo.db.Exec(query, user.FirstName, user.LastName, user.Email, user.Password, user.CreateAt, user.UpdatedAt)
	return err
}

func (repo *PostgresUserRepository) Update(user *models.User) error {
	query := `
	UPDATE users 
	SET first_name = $1, last_name = $2, email = $3, password = $4, create_at = $5, update_at = $6 
	WHERE id = $7
	`
	result, err := repo.db.Exec(query, user.FirstName, user.LastName, user.Email, user.Password, user.CreateAt, user.UpdatedAt, user.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("No rows affected")
	}
	return nil
}

func (repo *PostgresUserRepository) Delete(user *models.User) error {
	query := `
	DELETE FROM users 
	WHERE id = $1
	`
	result, err := repo.db.Exec(query, user.ID)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		return errors.New("No rows affected")
	}
	return nil
}
