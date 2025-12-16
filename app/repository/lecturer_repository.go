package repository

import (
	"context"
	"database/sql"
	"uas/app/models"

	"github.com/google/uuid"
)

type LecturerRepository interface {
	Create(ctx context.Context, l models.Lecturer) error
	FindByUserID(ctx context.Context, userID string) (models.Lecturer, error)
	FindAll(ctx context.Context) ([]models.Lecturer, error)
	FindAdvisees(ctx context.Context, lecturerID uuid.UUID) ([]models.Student, error)
}

type lecturerRepo struct {
	db *sql.DB
}

func NewLecturerRepository(db *sql.DB) LecturerRepository {
	return &lecturerRepo{db}
}

func (r *lecturerRepo) Create(ctx context.Context, l models.Lecturer) error {
	query := `
		INSERT INTO lecturers (
			id, user_id, lecturer_id, department
		) VALUES ($1, $2, $3, $4)
	`

	_, err := r.db.ExecContext(ctx, query,
		l.ID,
		l.UserID,
		l.LecturerID,
		l.Department,
	)

	return err
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

func (r *lecturerRepo) FindAll(ctx context.Context) ([]models.Lecturer, error) {
	rows, err := r.db.QueryContext(ctx, `
	SELECT id, user_id, lecturer_id, department
FROM lecturers
`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Lecturer
	for rows.Next() {
		var l models.Lecturer
		if err := rows.Scan(
			&l.ID, &l.UserID, &l.LecturerID, &l.Department,
		); err != nil {
			return nil, err
		}
		list = append(list, l)
	}
	return list, nil
}

func (r *lecturerRepo) FindAdvisees(ctx context.Context, lecturerID uuid.UUID,
) ([]models.Student, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id
		FROM students
		WHERE advisor_id = $1
	`, lecturerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var list []models.Student
	for rows.Next() {
		var s models.Student
		if err := rows.Scan(
			&s.ID, &s.UserID, &s.StudentID,
			&s.ProgramStudy, &s.AcademicYear, &s.AdvisorID,
		); err != nil {
			return nil, err
		}
		list = append(list, s)
	}
	return list, nil
}
