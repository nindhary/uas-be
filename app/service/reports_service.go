package service

import (
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"

	"uas/app/models"
	"uas/app/repository"
	"uas/helper"
)

type ReportService interface {
	GetStatistics(c *fiber.Ctx) error
	GetStudentStatistics(c *fiber.Ctx) error
}

type reportService struct {
	achievementRepo repository.AchievementRepository
	mongoRepo       repository.MongoAchievementRepository
	studentRepo     repository.StudentRepository
	lecturerRepo    repository.LecturerRepository
}

func NewReportService(
	achievementRepo repository.AchievementRepository,
	mongoRepo repository.MongoAchievementRepository,
	studentRepo repository.StudentRepository,
	lecturerRepo repository.LecturerRepository,
) ReportService {
	return &reportService{
		achievementRepo,
		mongoRepo,
		studentRepo,
		lecturerRepo,
	}
}

func (s *reportService) GetStatistics(c *fiber.Ctx) error {
	user := c.Locals("user").(models.Users)

	var refs []models.AchievementRef
	var err error

	if lecturer, errLect := s.lecturerRepo.FindByUserID(c.Context(), user.ID.String()); errLect == nil {
		refs, err = s.achievementRepo.FindByAdvisor(c.Context(), lecturer.ID)

	} else if student, errStd := s.studentRepo.FindByUserID(c.Context(), user.ID.String()); errStd == nil {
		refs, err = s.achievementRepo.FindByStudent(c.Context(), student.ID)

	} else {
		refs, err = s.achievementRepo.FindAll(
			c.Context(),
			"",
			1000,
			0,
		)
	}

	if err != nil {
		return helper.Error(c, 500, "failed load achievements")
	}

	typeCount := map[string]int{}
	levelCount := map[string]int{}
	periodCount := map[string]int{}
	studentCount := map[string]int{}

	for _, ref := range refs {
		if ref.Status != "verified" {
			continue
		}

		detail, err := s.mongoRepo.FindByHexID(
			c.Context(),
			ref.MongoAchievementID,
		)
		if err != nil {
			continue
		}

		typeCount[detail.Category]++
		levelCount[detail.Level]++

		if len(detail.EventDate) >= 4 {
			year := detail.EventDate[:4]
			periodCount[year]++
		}

		studentCount[ref.StudentID.String()]++
	}

	return helper.Success(c, fiber.Map{
		"total_by_type":   typeCount,
		"total_by_level":  levelCount,
		"total_by_period": periodCount,
		"top_students":    studentCount,
	})
}

func (s *reportService) GetStudentStatistics(c *fiber.Ctx) error {
	user := c.Locals("user").(models.Users)
	studentID := uuid.MustParse(c.Params("id"))

	if lecturer, errLect := s.lecturerRepo.FindByUserID(c.Context(), user.ID.String()); errLect == nil {

		student, err := s.studentRepo.FindByID(c.Context(), studentID)
		if err != nil || student.AdvisorID != lecturer.ID {
			return helper.Error(c, 403, "forbidden")
		}

	} else if student, errStd := s.studentRepo.FindByUserID(c.Context(), user.ID.String()); errStd == nil {

		if student.ID != studentID {
			return helper.Error(c, 403, "forbidden")
		}

	}

	refs, err := s.achievementRepo.FindByStudent(c.Context(), studentID)
	if err != nil {
		return helper.Error(c, 500, "failed load achievements")
	}

	typeCount := map[string]int{}
	levelCount := map[string]int{}
	periodCount := map[string]int{}

	for _, ref := range refs {
		if ref.Status != "verified" {
			continue
		}

		detail, err := s.mongoRepo.FindByHexID(
			c.Context(),
			ref.MongoAchievementID,
		)
		if err != nil {
			continue
		}

		typeCount[detail.Category]++
		levelCount[detail.Level]++

		if len(detail.EventDate) >= 4 {
			year := detail.EventDate[:4]
			periodCount[year]++
		}
	}

	return helper.Success(c, fiber.Map{
		"student_id":      studentID,
		"total_by_type":   typeCount,
		"total_by_level":  levelCount,
		"total_by_period": periodCount,
	})
}
