package models

import "github.com/google/uuid"

type Lecturer struct {
	ID         uuid.UUID `json:"id"`
	UserID     uuid.UUID `json:"userId"`
	LecturerID string    `json:"lecturerId"`
	Department string    `json:"department"`
	CreatedAt  string    `json:"createdAt"`
}
