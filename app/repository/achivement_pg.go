package repository

import (
	"context"
	"database/sql"
	"time"
	"uas/app/models"

	"github.com/google/uuid"
)

type AchievementRepository interface {
	CreateReference(ctx context.Context, ref models.AchievementRef) error
	FindByID(ctx context.Context, id uuid.UUID) (models.AchievementRef, error)
	FindByStudent(ctx context.Context, studentID uuid.UUID) ([]models.AchievementRef, error)
	UpdateStatusSubmitted(ctx context.Context, id uuid.UUID, submittedAt time.Time) error
	Delete(ctx context.Context, id uuid.UUID) error
}

type achievementRepo struct {
	db *sql.DB
}

func NewAchievementRepository(db *sql.DB) AchievementRepository {
	return &achievementRepo{db}
}

func (r *achievementRepo) CreateReference(ctx context.Context, ref models.AchievementRef) error {
	_, err := r.db.ExecContext(ctx, `
		INSERT INTO achievement_references (
			id, student_id, mongo_achievement_id, status, 
			created_at, updated_at
		)
		VALUES ($1, $2, $3, $4, NOW(), NOW())
	`,
		ref.ID,
		ref.StudentID,
		ref.MongoAchievementID,
		ref.Status,
	)
	return err
}

func (r *achievementRepo) FindByID(ctx context.Context, id uuid.UUID) (models.AchievementRef, error) {
	var ref models.AchievementRef

	err := r.db.QueryRowContext(ctx, `
		SELECT
			id, student_id, mongo_achievement_id, status,
			submitted_at, verified_at, verified_by,
			rejection_note, created_at, updated_at
		FROM achievement_references
		WHERE id = $1
	`, id).Scan(
		&ref.ID,
		&ref.StudentID,
		&ref.MongoAchievementID,
		&ref.Status,
		&ref.SubmittedAt,
		&ref.VerifiedAt,
		&ref.VerifiedBy,
		&ref.RejectionNote,
		&ref.CreatedAt,
		&ref.UpdatedAt,
	)

	return ref, err
}

func (r *achievementRepo) FindByStudent(ctx context.Context, studentID uuid.UUID) ([]models.AchievementRef, error) {
	rows, err := r.db.QueryContext(ctx, `
		SELECT
			id, student_id, mongo_achievement_id, status,
			submitted_at, verified_at, verified_by, rejection_note,
			created_at, updated_at
		FROM achievement_references
		WHERE student_id = $1
		ORDER BY created_at DESC
	`, studentID)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var list []models.AchievementRef

	for rows.Next() {
		var ref models.AchievementRef
		err := rows.Scan(
			&ref.ID,
			&ref.StudentID,
			&ref.MongoAchievementID,
			&ref.Status,
			&ref.SubmittedAt,
			&ref.VerifiedAt,
			&ref.VerifiedBy,
			&ref.RejectionNote,
			&ref.CreatedAt,
			&ref.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		list = append(list, ref)
	}

	return list, nil
}

func (r *achievementRepo) UpdateStatusSubmitted(ctx context.Context, id uuid.UUID, submittedAt time.Time) error {
	_, err := r.db.ExecContext(ctx, `
		UPDATE achievement_references
		SET status = 'submitted',
			submitted_at = $2,
			updated_at = NOW()
		WHERE id = $1
	`, id, submittedAt)

	return err
}

func (r *achievementRepo) Delete(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx, `
		DELETE FROM achievement_references
		WHERE id = $1
	`, id)

	return err
}
