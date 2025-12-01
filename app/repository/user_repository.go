package repository

import (
	"context"
	"database/sql"
	"errors"
	"uas/app/models"
)

type UserRepository interface {
	FindByUsername(ctx context.Context, username string) (models.Users, error)
}

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (models.Users, error) {

	var u models.Users

	query := `
		SELECT 
			id,
			username,
			email,
			password_hash,
			full_name,
			role_id,
			is_active
		FROM users
		WHERE username = $1
		LIMIT 1
	`

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&u.ID,
		&u.Username,
		&u.Email,
		&u.PasswordHash,
		&u.FullName,
		&u.RoleID,
		&u.IsActive,
	)

	if err != nil {
		return u, errors.New("user not found")
	}

	return u, nil
}
