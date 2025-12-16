package repository

import (
	"context"
	"database/sql"
	"uas/app/models"

	"github.com/google/uuid"
)

type StudentRepository interface {
	Create(ctx context.Context, s models.Student) error
	FindByUserID(ctx context.Context, userID string) (models.Student, error)
	FindByID(ctx context.Context, id uuid.UUID) (models.Student, error)
	FindAll(ctx context.Context) ([]models.Student, error)
	UpdateAdvisor(ctx context.Context, studentID uuid.UUID, advisorID uuid.UUID) error
}

type studentRepo struct {
	db *sql.DB
}

func NewStudentRepository(db *sql.DB) StudentRepository {
	return &studentRepo{db}
}

func (r *studentRepo) Create(ctx context.Context, s models.Student) error {
	query := `
		INSERT INTO students (
			id, user_id, student_id,
			program_study, academic_year, advisor_id
		) VALUES ($1, $2, $3, $4, $5, $6)
	`

	_, err := r.db.ExecContext(ctx, query,
		s.ID,
		s.UserID,
		s.StudentID,
		s.ProgramStudy,
		s.AcademicYear,
		s.AdvisorID,
	)

	return err
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
func (r *studentRepo) FindAll(ctx context.Context) ([]models.Student, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT id, user_id, student_id, program_study, academic_year, advisor_id
		FROM students
	`)
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

func (r *studentRepo) UpdateAdvisor(
	ctx context.Context,
	studentID uuid.UUID,
	advisorID uuid.UUID,
) error {

	query := `
		UPDATE students
		SET advisor_id = $1
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, advisorID, studentID)
	return err
}
