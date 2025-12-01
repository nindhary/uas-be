package models

type Role struct {
	ID          string `db:"id"`
	Name        string `db:"name"`
	Description string `db:"description"`
}
