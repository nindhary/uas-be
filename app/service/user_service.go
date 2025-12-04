package service

import (
	"uas/app/models"
	"uas/app/repository"
	"uas/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AdminUserService interface {
	GetAll(c *fiber.Ctx) error
	GetByID(c *fiber.Ctx) error
	Create(c *fiber.Ctx) error
	Update(c *fiber.Ctx) error
	Delete(c *fiber.Ctx) error
	UpdateRole(c *fiber.Ctx) error
}

type adminUserService struct {
	repo repository.AdminUserRepository
}

func NewAdminUserService(repo repository.AdminUserRepository) AdminUserService {
	return &adminUserService{repo}
}

func (s *adminUserService) GetAll(c *fiber.Ctx) error {
	users, err := s.repo.FindAll(c.Context())
	if err != nil {
		return helper.Error(c, 500, "failed fetch users")
	}
	return helper.Success(c, users)
}

func (s *adminUserService) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")
	user, err := s.repo.FindByID(c.Context(), id)
	if err != nil {
		return helper.Error(c, 404, "user not found")
	}
	return helper.Success(c, user)
}

func (s *adminUserService) Create(c *fiber.Ctx) error {
	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
		FullName string `json:"fullName"`
		RoleID   string `json:"roleId"`
	}
	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, 400, "invalid body")
	}

	hash, _ := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)

	user := models.Users{
		ID:           uuid.New(),
		Username:     req.Username,
		Email:        req.Email,
		PasswordHash: string(hash),
		FullName:     req.FullName,
		RoleID:       uuid.MustParse(req.RoleID),
		IsActive:     true,
	}

	if err := s.repo.Create(c.Context(), user); err != nil {
		return helper.Error(c, 500, "failed create user")
	}

	return helper.Success(c, "user created")
}

func (s *adminUserService) Update(c *fiber.Ctx) error {
	id := c.Params("id")

	var req struct {
		Username string `json:"username"`
		Email    string `json:"email"`
		FullName string `json:"fullName"`
		RoleID   string `json:"roleId"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, 400, "invalid body")
	}

	user := models.Users{
		ID:       uuid.MustParse(id),
		Username: req.Username,
		Email:    req.Email,
		FullName: req.FullName,
		RoleID:   uuid.MustParse(req.RoleID),
	}

	if err := s.repo.Update(c.Context(), user); err != nil {
		return helper.Error(c, 500, "failed update user")
	}

	return helper.Success(c, "user updated")
}

func (s *adminUserService) Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if err := s.repo.SoftDelete(c.Context(), id); err != nil {
		return helper.Error(c, 500, "failed to delete user")
	}

	return helper.Success(c, "user soft deleted")
}

func (s *adminUserService) UpdateRole(c *fiber.Ctx) error {
	id := c.Params("id")

	var req struct {
		RoleID string `json:"roleId"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, 400, "invalid body")
	}

	roleUUID := uuid.MustParse(req.RoleID)

	if err := s.repo.UpdateRole(c.Context(), id, roleUUID); err != nil {
		return helper.Error(c, 500, "failed to update role")
	}

	return helper.Success(c, "role updated")
}
