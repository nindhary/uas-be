package service

import (
	"uas/app/models"
	"uas/app/repository"
	"uas/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LecturerService interface {
	CreateProfile(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	GetAdvisees(c *fiber.Ctx) error
}

type lecturerService struct {
	repo repository.LecturerRepository
}

func NewLecturerService(repo repository.LecturerRepository) LecturerService {
	return &lecturerService{repo}
}

func (s *lecturerService) CreateProfile(c *fiber.Ctx) error {
	var req struct {
		UserID     string `json:"user_id"`
		LecturerID string `json:"lecturer_id"`
		Department string `json:"department"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, 400, "invalid body")
	}

	lecturer := models.Lecturer{
		ID:         uuid.New(),
		UserID:     uuid.MustParse(req.UserID),
		LecturerID: req.LecturerID,
		Department: req.Department,
	}

	if err := s.repo.Create(c.Context(), lecturer); err != nil {
		return helper.Error(c, 500, "failed create lecturer profile")
	}

	return helper.Success(c, "lecturer profile created")
}

func (s *lecturerService) GetAll(c *fiber.Ctx) error {
	list, err := s.repo.FindAll(c.Context())
	if err != nil {
		return helper.Error(c, 500, "gagal menampilkan dosen")
	}
	return helper.Success(c, list)
}

func (s *lecturerService) GetAdvisees(c *fiber.Ctx) error {
	user := c.Locals("user").(models.Users)

	lecturer, _ := s.repo.FindByUserID(c.Context(), user.ID.String())

	list, err := s.repo.FindAdvisees(c.Context(), lecturer.ID)
	if err != nil {
		return helper.Error(c, 500, "failed load mahasiswa bimbingan")
	}

	return helper.Success(c, list)
}
