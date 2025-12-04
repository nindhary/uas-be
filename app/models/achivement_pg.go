package models

import (
	"time"

	"github.com/google/uuid"
)

type AchievementRef struct {
	ID                 uuid.UUID `json:"id"`
	StudentID          uuid.UUID `json:"studentId"`
	MongoAchievementID string    `json:"mongoAchievementId"`
	Status             string    `json:"status"`

	SubmittedAt *time.Time `json:"submittedAt"`
	VerifiedAt  *time.Time `json:"verifiedAt"`
	VerifiedBy  *uuid.UUID `json:"verifiedBy"`

	RejectionNote *string   `json:"rejectionNote"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
}
