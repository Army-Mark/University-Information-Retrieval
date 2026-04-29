package repository

import (
	"database/sql"
	"school-go/internal/models"
)

type UserRepository interface {
	GetAll() ([]models.User, error)
	GetByUsername(username string) (*models.User, error)
	Create(user *models.User) error
	Update(oldUsername, newUsername, newPassword, newRole string) error
	Delete(username string) error
	SaveAll(users []models.User) error
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository() UserRepository {
	return &userRepository{
		db: GetDB(),
	}
}

func (r *userRepository) GetAll() ([]models.User, error) {
	rows, err := r.db.Query(`SELECT username, password, role FROM users`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var u models.User
		if err := rows.Scan(&u.Username, &u.Password, &u.Role); err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepository) GetByUsername(username string) (*models.User, error) {
	row := r.db.QueryRow(`SELECT username, password, role FROM users WHERE username = ?`, username)
	var u models.User
	if err := row.Scan(&u.Username, &u.Password, &u.Role); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *userRepository) Create(user *models.User) error {
	_, err := r.db.Exec(`INSERT INTO users (username, password, role) VALUES (?, ?, ?)`, user.Username, user.Password, user.Role)
	return err
}

func (r *userRepository) Update(oldUsername, newUsername, newPassword, newRole string) error {
	if _, err := r.db.Exec(`DELETE FROM users WHERE username = ?`, oldUsername); err != nil {
		return err
	}
	_, err := r.db.Exec(`INSERT INTO users (username, password, role) VALUES (?, ?, ?)`, newUsername, newPassword, newRole)
	return err
}

func (r *userRepository) Delete(username string) error {
	_, err := r.db.Exec(`DELETE FROM users WHERE username = ?`, username)
	return err
}

func (r *userRepository) SaveAll(users []models.User) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	if _, err := tx.Exec(`DELETE FROM users`); err != nil {
		tx.Rollback()
		return err
	}
	for _, u := range users {
		if _, err := tx.Exec(`INSERT INTO users (username, password, role) VALUES (?, ?, ?)`, u.Username, u.Password, u.Role); err != nil {
			tx.Rollback()
			return err
		}
	}
	return tx.Commit()
}
