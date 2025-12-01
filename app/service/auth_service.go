package service

import (
	"fmt"
	"os"
	"time"

	"uas/app/repository"
	"uas/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	LoginHandler(c *fiber.Ctx) error
}

type authService struct {
	users repository.UserRepository
}

func NewAuthService(users repository.UserRepository) AuthService {
	return &authService{users}
}

func (s *authService) LoginHandler(c *fiber.Ctx) error {

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := c.BodyParser(&req); err != nil {
		return helper.Error(c, 400, "invalid request body")
	}

	// debug input
	fmt.Println("========== LOGIN DEBUG ==========")
	fmt.Println("Input Username:", req.Username)
	fmt.Println("Input Password:", req.Password)

	// find user
	user, err := s.users.FindByUsername(c.Context(), req.Username)
	if err != nil {
		fmt.Println("ERROR: User not found in DB")
		return helper.Error(c, 401, "invalid credentials")
	}

	// debug user
	fmt.Println("DB Username:", user.Username)
	fmt.Println("DB PasswordHash:", user.PasswordHash)
	fmt.Println("DB RoleID:", user.RoleID)
	fmt.Println("DB IsActive:", user.IsActive)

	// cek status aktif
	if !user.IsActive {
		fmt.Println("ERROR: User inactive")
		return helper.Error(c, 403, "account is not active")
	}

	// bcrypt pass check
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password))
	if err != nil {
		fmt.Println("Bcrypt Result: FAILED")
		fmt.Println("Error:", err)
		return helper.Error(c, 401, "invalid credentials")
	}

	fmt.Println("Bcrypt Result: SUCCESS")
	fmt.Println("User authenticated successfully!")

	// generate token jwt
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "DEFAULT_SECRET"
	}

	claims := jwt.MapClaims{
		"id":     user.ID.String(),
		"roleId": user.RoleID.String(),
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := token.SignedString([]byte(secret))
	if err != nil {
		fmt.Println("ERROR: JWT generation failed")
		return helper.Error(c, 500, "failed to generate token")
	}

	fmt.Println("JWT Generated Successfully:", signedToken)
	fmt.Println("=============== END DEBUG ================")

	// respon sukses
	return helper.Success(c, fiber.Map{
		"token": signedToken,
		"user": fiber.Map{
			"id":       user.ID,
			"username": user.Username,
			"email":    user.Email,
			"fullName": user.FullName,
			"roleId":   user.RoleID,
			"isActive": user.IsActive,
		},
	})
}
