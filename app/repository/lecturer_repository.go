package repository

import (
	"context"
	"database/sql"
	"uas/app/models"
)

type LecturerRepository interface {
	FindByUserID(ctx context.Context, userID string) (models.Lecturer, error)
}

type lecturerRepo struct {
	db *sql.DB
}

func NewLecturerRepository(db *sql.DB) LecturerRepository {
	return &lecturerRepo{db}
}

func (r *lecturerRepo) FindByUserID(ctx context.Context, userID string) (models.Lecturer, error) {
	var l models.Lecturer

	query := `
        SELECT id, user_id, lecturer_id, department
        FROM lecturers
        WHERE user_id = $1
        LIMIT 1
    `

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&l.ID, &l.UserID, &l.LecturerID, &l.Department,
	)

	return l, err
}
