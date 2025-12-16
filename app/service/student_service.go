package service

import (
	"uas/app/models"
	"uas/app/repository"
	"uas/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type StudentService interface {
	CreateProfile(c *fiber.Ctx) error
	GetAll(c *fiber.Ctx) error
	GetDetail(c *fiber.Ctx) error
	GetMyAchievements(c *fiber.Ctx) error
	UpdateAdvisor(c *fiber.Ctx) error
}

type studentService struct {
	repo            repository.StudentRepository
	achievementRepo repository.AchievementRepository
}

func NewStudentService(
	repo repository.StudentRepository,
	achievementRepo repository.AchievementRepository,
) StudentService {
	return &studentService{repo, achievementRepo}
}

func (s *studentService) CreateProfile(c *fiber.Ctx) error {
	var req struct {
		UserID       string `json:"user_id"`
		StudentID    string `json:"student_id"`
		ProgramStudy string `json:"program_study"`
		AcademicYear string `json:"academic_year"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, 400, "invalid body")
	}

	student := models.Student{
		ID:           uuid.New(),
		UserID:       uuid.MustParse(req.UserID),
		StudentID:    req.StudentID,
		ProgramStudy: req.ProgramStudy,
		AcademicYear: req.AcademicYear,
	}

	if err := s.repo.Create(c.Context(), student); err != nil {
		return helper.Error(c, 500, "failed create student profile")
	}

	return helper.Success(c, "student profile created")
}

func (s *studentService) GetAll(c *fiber.Ctx) error {
	list, err := s.repo.FindAll(c.Context())
	if err != nil {
		return helper.Error(c, 500, "gagal menampilkan mahasiswa")
	}
	return helper.Success(c, list)
}

func (s *studentService) GetDetail(c *fiber.Ctx) error {
	id := uuid.MustParse(c.Params("id"))
	user := c.Locals("user").(models.Users)

	student, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		return helper.Error(c, 404, "mahasiswa tidak ditemukan")
	}

	if student.UserID != user.ID {
		return helper.Error(c, 403, "forbidden")
	}

	return helper.Success(c, student)
}

func (s *studentService) GetMyAchievements(c *fiber.Ctx) error {
	id := uuid.MustParse(c.Params("id"))
	user := c.Locals("user").(models.Users)

	student, _ := s.repo.FindByID(c.Context(), id)
	if student.UserID != user.ID {
		return helper.Error(c, 403, "forbidden")
	}

	list, err := s.achievementRepo.FindByStudent(c.Context(), student.ID)
	if err != nil {
		return helper.Error(c, 500, "gagal menampilkan achievements")
	}

	return helper.Success(c, list)
}

func (s *studentService) UpdateAdvisor(c *fiber.Ctx) error {
	studentID := uuid.MustParse(c.Params("id"))

	var req struct {
		AdvisorID uuid.UUID `json:"advisor_id"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, 400, "invalid body")
	}

	if req.AdvisorID == uuid.Nil {
		return helper.Error(c, 400, "advisor_id required")
	}

	if _, err := s.repo.FindByID(c.Context(), studentID); err != nil {
		return helper.Error(c, 404, "mahasiswa tidak ditemukan")
	}

	if err := s.repo.UpdateAdvisor(c.Context(), studentID, req.AdvisorID); err != nil {
		return helper.Error(c, 500, "gagal update advisor")
	}

	return helper.Success(c, "advisor updated")
}
