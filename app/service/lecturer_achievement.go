package service

import (
	"fmt"
	"time"
	"uas/app/models"
	"uas/app/repository"
	"uas/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
)

type LecturerAchievementService interface {
	GetAdviseeAchievements(c *fiber.Ctx) error
	GetDetail(c *fiber.Ctx) error
	Verify(c *fiber.Ctx) error
	Reject(c *fiber.Ctx) error
	GetHistory(c *fiber.Ctx) error
}

type lecturerAchievementService struct {
	repo         repository.AchievementRepository
	studentRepo  repository.StudentRepository
	lecturerRepo repository.LecturerRepository
	mongo        repository.MongoAchievementRepository
}

func NewLecturerAchievementService(
	repo repository.AchievementRepository,
	studentRepo repository.StudentRepository,
	lecturerRepo repository.LecturerRepository,
	mongo repository.MongoAchievementRepository,
) LecturerAchievementService {
	return &lecturerAchievementService{repo, studentRepo, lecturerRepo, mongo}
}

func (s *lecturerAchievementService) GetAdviseeAchievements(c *fiber.Ctx) error {
	user := c.Locals("user").(models.Users)

	lecturer, err := s.lecturerRepo.FindByUserID(c.Context(), user.ID.String())
	if err != nil {
		return helper.Error(c, 403, "hanya dosen wali yang dapat mengakses")
	}

	list, err := s.repo.FindByAdvisor(c.Context(), lecturer.ID)
	if err != nil {
		return helper.Error(c, 500, "gagal memuat prestasi")
	}
	return helper.Success(c, list)
}

func (s *lecturerAchievementService) GetDetail(c *fiber.Ctx) error {
	refID := uuid.MustParse(c.Params("id"))
	user := c.Locals("user").(models.Users)

	lecturer, err := s.lecturerRepo.FindByUserID(
		c.Context(),
		user.ID.String(),
	)
	if err != nil {
		return helper.Error(c, 403, "bukan dosen")
	}

	ref, err := s.repo.FindByID(c.Context(), refID)
	if err != nil {
		return helper.Error(c, 404, "prestasi tidak ditemukan")
	}

	student, err := s.studentRepo.FindByID(
		c.Context(),
		ref.StudentID,
	)
	if err != nil {
		return helper.Error(c, 404, "student tidak ditemukan")
	}

	if student.AdvisorID != lecturer.ID {
		return helper.Error(c, 403, "bukan mahasiswa bimbingan")
	}

	detail, err := s.mongo.FindByHexID(
		c.Context(),
		ref.MongoAchievementID,
	)
	if err != nil {
		return helper.Error(c, 500, "gagal mengambil data dari mongo")
	}

	return helper.Success(c, detail)
}

func (s *lecturerAchievementService) Verify(c *fiber.Ctx) error {
	refIDStr := c.Params("id")

	refID, err := uuid.Parse(refIDStr)
	if err != nil {
		return helper.Error(c, 400, "invalid id")
	}

	user := c.Locals("user").(models.Users)

	lecturer, err := s.lecturerRepo.FindByUserID(c.Context(), user.ID.String())
	if err != nil {
		return helper.Error(c, 403, "bukan dosen")
	}
	fmt.Println("Lecturer ID:", lecturer.ID)

	ref, err := s.repo.FindByID(c.Context(), refID)
	if err != nil {
		return helper.Error(c, 404, "prestasi tidak ditemukan")
	}

	if ref.Status != "submitted" {
		return helper.Error(c, 400, "hanya item tersubmit yang dapat diverifikasi")
	}

	student, err := s.studentRepo.FindByID(c.Context(), ref.StudentID)
	if err != nil {
		return helper.Error(c, 404, "mahasiswa tidak ditemukan")
	}

	fmt.Println("Student AdvisorID:", student.AdvisorID)
	fmt.Println("Lecturer ID:", lecturer.ID)

	if student.AdvisorID != lecturer.ID {
		return helper.Error(c, 403, "bukan mahasiswa bimbingan")
	}
	now := time.Now()

	err = s.repo.UpdateStatusVerified(c.Context(), refID, user.ID, now)
	if err != nil {
		return helper.Error(c, 500, "gagal diverifikasi")
	}

	err = s.mongo.PushHistoryByHexID(
		c.Context(),
		ref.MongoAchievementID,
		"verified",
	)
	if err != nil {
		fmt.Println("MONGO HISTORY ERROR:", err)
	}

	fmt.Println("VERIFY SUCCESS")
	return helper.Success(c, "verified")
}

func (s *lecturerAchievementService) Reject(c *fiber.Ctx) error {
	refID := uuid.MustParse(c.Params("id"))
	user := c.Locals("user").(models.Users)
	now := time.Now()

	var req struct {
		Note string `json:"note"`
	}
	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, 400, "invalid request body")
	}

	lecturer, _ := s.lecturerRepo.FindByUserID(c.Context(), user.ID.String())
	ref, _ := s.repo.FindByID(c.Context(), refID)

	if ref.Status != "submitted" {
		return helper.Error(c, 400, "hanya item tersubmit yang dapat ditolak")
	}

	student, _ := s.studentRepo.FindByID(c.Context(), ref.StudentID)
	if student.AdvisorID != lecturer.ID {
		return helper.Error(c, 403, "bukan mahasiswa bimbingan")
	}

	err := s.repo.UpdateStatusRejected(
		c.Context(),
		refID,
		user.ID,
		req.Note,
		now,
	)
	if err != nil {
		return helper.Error(c, 500, "gagal untuk reject")
	}

	s.mongo.PushHistoryByHexID(
		c.Context(),
		ref.MongoAchievementID,
		"rejected",
	)

	return helper.Success(c, "rejected")
}

func (s *lecturerAchievementService) GetHistory(c *fiber.Ctx) error {
	refID := uuid.MustParse(c.Params("id"))
	user := c.Locals("user").(models.Users)

	lecturer, err := s.lecturerRepo.FindByUserID(c.Context(), user.ID.String())
	if err != nil {
		return helper.Error(c, 403, "bukan dosen")
	}

	ref, err := s.repo.FindByID(c.Context(), refID)
	if err != nil {
		return helper.Error(c, 404, "prestasi tidak ditemukan")
	}

	student, err := s.studentRepo.FindByID(c.Context(), ref.StudentID)
	if err != nil {
		return helper.Error(c, 404, "mahasiswa tidak ditemukan")
	}

	if student.AdvisorID != lecturer.ID {
		return helper.Error(c, 403, "bukan mahasiswa bimbingan")
	}

	history, err := s.mongo.GetHistory(c.Context(), ref.MongoAchievementID)
	if err != nil {
		return helper.Error(c, 500, "gagal mendapatkan history")
	}

	return helper.Success(c, history)
}
