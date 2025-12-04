package repository

import (
	"context"
	"database/sql"
	"uas/app/models"

	"github.com/google/uuid"
)

type AdminUserRepository interface {
	FindAll(ctx context.Context) ([]models.Users, error)
	FindByID(ctx context.Context, id string) (models.Users, error)
	Create(ctx context.Context, user models.Users) error
	Update(ctx context.Context, user models.Users) error
	SoftDelete(ctx context.Context, id string) error
	UpdateRole(ctx context.Context, id string, roleID uuid.UUID) error
	GetUserPermissions(ctx context.Context, userID string) ([]string, error)
}

type adminUserRepo struct {
	db *sql.DB
}

func NewAdminUserRepository(db *sql.DB) AdminUserRepository {
	return &adminUserRepo{db}
}

func (r *adminUserRepo) FindAll(ctx context.Context) ([]models.Users, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, username, email, full_name, role_id, is_active
		FROM users
		WHERE is_active = true
	`)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var users []models.Users
	for rows.Next() {
		var u models.Users
		rows.Scan(&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.IsActive)
		users = append(users, u)
	}
	return users, nil
}

func (r *adminUserRepo) FindByID(ctx context.Context, id string) (models.Users, error) {
	var u models.Users
	err := r.db.QueryRowContext(ctx, `
		SELECT id, username, email, full_name, role_id, is_active
		FROM users WHERE id = $1
	`, id).Scan(
		&u.ID, &u.Username, &u.Email, &u.FullName, &u.RoleID, &u.IsActive,
	)
	return u, err
}

func (r *adminUserRepo) Create(ctx context.Context, u models.Users) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO users (id, username, email, password_hash, full_name, role_id)
		VALUES ($1,$2,$3,$4,$5,$6)
	`,
		u.ID, u.Username, u.Email, u.PasswordHash, u.FullName, u.RoleID)
	return err
}

func (r *adminUserRepo) Update(ctx context.Context, u models.Users) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET username=$2, email=$3, full_name=$4, role_id=$5
		WHERE id=$1
	`,
		u.ID, u.Username, u.Email, u.FullName, u.RoleID)
	return err
}

func (r *adminUserRepo) SoftDelete(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET is_active=false WHERE id=$1
	`, id)
	return err
}

func (r *adminUserRepo) UpdateRole(ctx context.Context, id string, roleID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE users SET role_id=$2 WHERE id=$1
	`, id, roleID)
	return err
}

func (r *adminUserRepo) GetUserPermissions(ctx context.Context, userID string) ([]string, error) {
	query := `
		SELECT p.name
		FROM role_permissions rp
		JOIN permissions p ON p.id = rp.permission_id
		JOIN users u ON u.role_id = rp.role_id
		WHERE u.id = $1
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var perms []string

	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err != nil {
			return nil, err
		}
		perms = append(perms, name)
	}

	return perms, nil
}
