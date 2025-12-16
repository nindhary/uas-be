package service

import (
	"os"
	"path/filepath"
	"time"
	"uas/app/models"
	"uas/app/repository"
	"uas/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type StudentAchievementService interface {
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	Submit(c *fiber.Ctx) error
	GetMyAchievements(c *fiber.Ctx) error
	GetDetail(c *fiber.Ctx) error
	GetHistory(c *fiber.Ctx) error
	UploadAttachment(c *fiber.Ctx) error
}

type studentAchievementService struct {
	repo        repository.AchievementRepository
	studentRepo repository.StudentRepository
	mongo       repository.MongoAchievementRepository
}

func NewStudentAchievementService(
	repo repository.AchievementRepository,
	studentRepo repository.StudentRepository,
	mongo repository.MongoAchievementRepository,
) StudentAchievementService {
	return &studentAchievementService{repo, studentRepo, mongo}
}

func (s *studentAchievementService) Create(c *fiber.Ctx) error {
	user := c.Locals("user").(models.Users)
	student, err := s.studentRepo.FindByUserID(c.Context(), user.ID.String())
	if err != nil {
		return helper.Error(c, 403, "only students can create achievements")
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Level       string `json:"level"`
		EventDate   string `json:"eventDate"`
	}
	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, 400, "invalid request body")
	}
	now := time.Now()

	detail := models.AchievementDetail{
		ID:            primitive.NewObjectID(),
		AchievementID: uuid.New().String(),
		StudentID:     student.ID.String(),
		Title:         req.Title,
		Description:   req.Description,
		Category:      req.Category,
		Level:         req.Level,
		EventDate:     req.EventDate,
		Attachments:   []string{},
		History: []struct {
			Status    string    `bson:"status"`
			Timestamp time.Time `bson:"timestamp"`
			ChangedBy string    `bson:"changedBy"`
		}{
			{Status: "draft", Timestamp: now, ChangedBy: user.ID.String()},
		},
		CreatedAt: now,
		UpdatedAt: now,
	}

	if _, err := s.mongo.Insert(c.Context(), detail); err != nil {
		return helper.Error(c, 500, "failed insert mongo")
	}

	ref := models.AchievementRef{
		ID:                 uuid.New(),
		StudentID:          student.ID,
		MongoAchievementID: detail.ID.Hex(),
		Status:             "draft",
		CreatedAt:          now,
		UpdatedAt:          now,
	}

	if err := s.repo.CreateReference(c.Context(), ref); err != nil {
		return helper.Error(c, 500, "failed insert reference")
	}

	return helper.Success(c, fiber.Map{
		"id":      ref.ID,
		"mongoId": detail.ID.Hex(),
		"status":  "draft",
	})
}

func (s *studentAchievementService) Update(c *fiber.Ctx) error {
	refID := uuid.MustParse(c.Params("id"))

	ref, err := s.repo.FindByID(c.Context(), refID)
	if err != nil {
		return helper.Error(c, 404, "achievement not found")
	}

	if ref.Status != "draft" {
		return helper.Error(c, 400, "only draft can be updated")
	}

	var req struct {
		Title       string `json:"title"`
		Description string `json:"description"`
		Category    string `json:"category"`
		Level       string `json:"level"`
		EventDate   string `json:"eventDate"`
	}

	c.BodyParser(&req)

	update := bson.M{
		"$set": bson.M{
			"title":       req.Title,
			"description": req.Description,
			"category":    req.Category,
			"level":       req.Level,
			"eventDate":   req.EventDate,
			"updatedAt":   time.Now(),
		},
		"$push": bson.M{
			"history": bson.M{
				"status":    "draft-updated",
				"timestamp": time.Now(),
			},
		},
	}

	err = s.mongo.UpdateByHexID(c.Context(), ref.MongoAchievementID, update)
	if err != nil {
		return helper.Error(c, 500, "failed update mongo")
	}

	return helper.Success(c, "updated")
}

func (s *studentAchievementService) Delete(c *fiber.Ctx) error {
	refID := uuid.MustParse(c.Params("id"))

	ref, err := s.repo.FindByID(c.Context(), refID)
	if err != nil {
		return helper.Error(c, 404, "achievement not found")
	}

	if ref.Status != "draft" {
		return helper.Error(c, 400, "only draft can be deleted")
	}

	if err := s.mongo.DeleteByHexID(c.Context(), ref.MongoAchievementID); err != nil {
		return helper.Error(c, 500, "failed delete mongo")
	}

	if err := s.repo.Delete(c.Context(), refID); err != nil {
		return helper.Error(c, 500, "failed delete reference")
	}

	return helper.Success(c, "deleted")
}

func (s *studentAchievementService) Submit(c *fiber.Ctx) error {
	refID := uuid.MustParse(c.Params("id"))

	ref, err := s.repo.FindByID(c.Context(), refID)
	if err != nil {
		return helper.Error(c, 404, "achievement not found")
	}

	if ref.Status != "draft" {
		return helper.Error(c, 400, "only draft can be submitted")
	}

	now := time.Now()

	if err := s.repo.UpdateStatusSubmitted(c.Context(), refID, now); err != nil {
		return helper.Error(c, 500, "failed update status")
	}

	err = s.mongo.PushHistoryByHexID(
		c.Context(),
		ref.MongoAchievementID,
		"submitted",
	)
	if err != nil {
		return helper.Error(c, 500, "failed update history")
	}

	return helper.Success(c, "submitted")
}

func (s *studentAchievementService) GetMyAchievements(c *fiber.Ctx) error {
	user := c.Locals("user").(models.Users)

	student, _ := s.studentRepo.FindByUserID(c.Context(), user.ID.String())

	list, err := s.repo.FindByStudent(c.Context(), student.ID)
	if err != nil {
		return helper.Error(c, 500, "failed load achievements")
	}

	return helper.Success(c, list)
}

func (s *studentAchievementService) GetDetail(c *fiber.Ctx) error {
	refID := c.Params("id")

	uid, err := uuid.Parse(refID)
	if err != nil {
		return helper.Error(c, 400, "invalid id")
	}

	ref, err := s.repo.FindByID(c.Context(), uid)
	if err != nil {
		return helper.Error(c, 404, "achievement not found")
	}

	user := c.Locals("user").(models.Users)
	student, _ := s.studentRepo.FindByUserID(c.Context(), user.ID.String())
	if ref.StudentID != student.ID {
		return helper.Error(c, 403, "forbidden: this achievement is not yours")
	}

	detail, err := s.mongo.FindByHexID(c.Context(), ref.MongoAchievementID)
	if err != nil {
		return helper.Error(c, 500, "failed load mongo detail")
	}

	return helper.Success(c, fiber.Map{
		"reference": ref,
		"detail":    detail,
	})
}

func (s *studentAchievementService) GetHistory(c *fiber.Ctx) error {
	refID := c.Params("id")

	uid, err := uuid.Parse(refID)
	if err != nil {
		return helper.Error(c, 400, "invalid id")
	}

	ref, err := s.repo.FindByID(c.Context(), uid)
	if err != nil {
		return helper.Error(c, 404, "achievement not found")
	}

	user := c.Locals("user").(models.Users)
	student, _ := s.studentRepo.FindByUserID(c.Context(), user.ID.String())
	if ref.StudentID != student.ID {
		return helper.Error(c, 403, "forbidden")
	}

	detail, err := s.mongo.FindByHexID(c.Context(), ref.MongoAchievementID)
	if err != nil {
		return helper.Error(c, 500, "failed load history")
	}

	return helper.Success(c, fiber.Map{
		"id":      ref.ID,
		"status":  ref.Status,
		"history": detail.History,
	})
}

func (s *studentAchievementService) UploadAttachment(c *fiber.Ctx) error {
	refID := uuid.MustParse(c.Params("id"))
	user := c.Locals("user").(models.Users)

	ref, err := s.repo.FindByID(c.Context(), refID)
	if err != nil {
		return helper.Error(c, 404, "achievement not found")
	}

	student, _ := s.studentRepo.FindByUserID(c.Context(), user.ID.String())
	if ref.StudentID != student.ID {
		return helper.Error(c, 403, "forbidden")
	}

	if ref.Status != "draft" {
		return helper.Error(c, 400, "only draft can upload attachments")
	}

	file, err := c.FormFile("file")
	if err != nil {
		return helper.Error(c, 400, "file required")
	}

	ext := filepath.Ext(file.Filename)
	if ext != ".pdf" && ext != ".jpg" && ext != ".jpeg" && ext != ".png" {
		return helper.Error(c, 400, "invalid file type")
	}

	baseDir := "./uploads" + ref.ID.String()
	if err := os.MkdirAll(baseDir, os.ModePerm); err != nil {
		return helper.Error(c, 500, "failed create upload dir")
	}

	filename := uuid.New().String() + ext
	fullPath := baseDir + "/" + filename

	if err := c.SaveFile(file, fullPath); err != nil {
		return helper.Error(c, 500, "failed save file")
	}

	update := bson.M{
		"$push": bson.M{
			"attachments": filename,
			"history": bson.M{
				"status":    "attachment-added",
				"timestamp": time.Now(),
				"changedBy": user.ID.String(),
			},
		},
		"$set": bson.M{
			"updatedAt": time.Now(),
		},
	}

	if err := s.mongo.UpdateByHexID(c.Context(), ref.MongoAchievementID, update); err != nil {
		return helper.Error(c, 500, "failed update mongo")
	}

	return helper.Success(c, fiber.Map{
		"file": filename,
	})
}
