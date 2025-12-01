package models

import "github.com/google/uuid"

type Users struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	FullName     string    `json:"fullName"`
	RoleID       uuid.UUID `json:"roleId"`
	IsActive     bool      `json:"isActive"`
}
