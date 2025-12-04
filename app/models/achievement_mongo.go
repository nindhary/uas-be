package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type AchievementDetail struct {
	ID            primitive.ObjectID `bson:"_id,omitempty"`
	AchievementID string             `bson:"achievementId"`
	StudentID     string             `bson:"studentId"`

	Title       string   `bson:"title"`
	Description string   `bson:"description"`
	Category    string   `bson:"category"`
	Level       string   `bson:"level"`
	EventDate   string   `bson:"eventDate"`
	Attachments []string `bson:"attachments"`

	History []struct {
		Status    string    `bson:"status"`
		Timestamp time.Time `bson:"timestamp"`
		ChangedBy string    `bson:"changedBy"`
	} `bson:"history"`

	CreatedAt time.Time `bson:"createdAt"`
	UpdatedAt time.Time `bson:"updatedAt"`
}
