package models

import "github.com/google/uuid"

type Student struct {
	ID           uuid.UUID `json:"id"`
	UserID       uuid.UUID `json:"userId"`
	StudentID    string    `json:"studentId"`
	ProgramStudy string    `json:"programStudy"`
	AcademicYear string    `json:"academicYear"`
	AdvisorID    uuid.UUID `json:"advisorId"`
	CreatedAt    string    `json:"createdAt"`
}
