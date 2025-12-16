package service

import (
	"fmt"
	"strconv"
	"uas/app/repository"
	"uas/helper"

	"github.com/gofiber/fiber/v2"
)

type AdminAchievementService interface {
	GetAll(c *fiber.Ctx) error
}

type adminAchievementService struct {
	pgRepo    repository.AchievementRepository
	mongoRepo repository.MongoAchievementRepository
}

func NewAdminAchievementService(
	pgRepo repository.AchievementRepository,
	mongoRepo repository.MongoAchievementRepository,
) AdminAchievementService {
	return &adminAchievementService{pgRepo, mongoRepo}
}

func (s *adminAchievementService) GetAll(c *fiber.Ctx) error {
	fmt.Println("=== ADMIN GET ALL ACHIEVEMENTS ===")

	status := c.Query("status", "")
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))

	refs, err := s.pgRepo.FindAll(
		c.Context(),
		status,
		limit,
		offset,
	)

	if err != nil {
		fmt.Println("ERROR PG FIND ALL:", err)
		return helper.Error(c, 500, "failed get achievement references")
	}

	fmt.Println("TOTAL REFERENCES:", len(refs))

	results := make([]fiber.Map, 0)

	for _, ref := range refs {
		fmt.Println("FETCH MONGO ID:", ref.MongoAchievementID)

		detail, err := s.mongoRepo.FindByHexID(
			c.Context(),
			ref.MongoAchievementID,
		)
		if err != nil {
			fmt.Println("MONGO ERROR ID:", ref.MongoAchievementID, err)
			continue
		}

		results = append(results, fiber.Map{
			"reference": ref,
			"detail":    detail,
		})
	}

	fmt.Println("RETURN DATA COUNT:", len(results))

	return helper.Success(c, results)
}
