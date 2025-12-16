package repository

import (
	"context"
	"database/sql"
	"uas/app/models"

	"github.com/google/uuid"
)

type StudentRepository interface {
	FindByUserID(ctx context.Context, userID string) (models.Student, error)
	FindByID(ctx context.Context, id uuid.UUID) (models.Student, error)
}

type studentRepo struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) StudentRepository {
	return &studentRepo{db}
}

func (r *studentRepo) FindByUserID(ctx context.Context, userID string) (models.Student, error) {
	var s models.Student

	query := `
        SELECT id, user_id, student_id, program_study, academic_year, advisor_id
        FROM students
        WHERE user_id = $1
        LIMIT 1
    `

	err := r.db.QueryRowContext(ctx, query, userID).Scan(
		&s.ID, &s.UserID, &s.StudentID, &s.ProgramStudy, &s.AcademicYear, &s.AdvisorID,
	)

	return s, err
}

func (r *studentRepo) FindByID(ctx context.Context, id uuid.UUID) (models.Student, error) {

	var s models.Student

	err := r.db.QueryRowContext(ctx, `
		SELECT 
			id,
			user_id,
			student_id,
			program_study,
			academic_year,
			advisor_id
		FROM students
		WHERE id = $1
	`, id).Scan(
		&s.ID,
		&s.UserID,
		&s.StudentID,
		&s.ProgramStudy,
		&s.AcademicYear,
		&s.AdvisorID,
	)

	if err != nil {
		return s, err
	}

	return s, nil
}
