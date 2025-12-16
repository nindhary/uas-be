package service

import (
	"fmt"
	"os"
	"strings"
	"time"

	"uas/app/repository"
	"uas/helper"
	"uas/middleware"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	LoginHandler(c *fiber.Ctx) error
	RefreshHandler(c *fiber.Ctx) error
	ProfileHandler(c *fiber.Ctx) error
	Logout(c *fiber.Ctx) error
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

	user, err := s.users.FindByUsername(c.Context(), req.Username)
	if err != nil {
		return helper.Error(c, 401, "invalid credentials")
	}

	if !user.IsActive {
		return helper.Error(c, 403, "account is not active")
	}

	// cek password
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return helper.Error(c, 401, "invalid credentials")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "DEFAULT_SECRET"
	}

	accessClaims := jwt.MapClaims{
		"id":     user.ID.String(),
		"roleId": user.RoleID.String(),
		"exp":    time.Now().Add(time.Hour * 24).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	signedAccess, err := accessToken.SignedString([]byte(secret))
	if err != nil {
		return helper.Error(c, 500, "failed to generate access token")
	}

	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if refreshSecret == "" {
		refreshSecret = "DEFAULT_REFRESH_SECRET"
	}

	refreshClaims := jwt.MapClaims{
		"id":  user.ID.String(),
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(),
	}

	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)
	signedRefresh, err := refreshToken.SignedString([]byte(refreshSecret))
	if err != nil {
		return helper.Error(c, 500, "failed to generate refresh token")
	}

	return helper.Success(c, fiber.Map{
		"accessToken":  signedAccess,
		"refreshToken": signedRefresh,
	})
}

func (s *authService) RefreshHandler(c *fiber.Ctx) error {
	fmt.Println("=== DEBUG REFRESH TOKEN ===")

	var req struct {
		RefreshToken string `json:"refreshToken"`
	}

	if err := c.BodyParser(&req); err != nil {
		fmt.Println("BodyParser error:", err)
		return helper.Error(c, 400, "invalid request body")
	}

	fmt.Println("Input RefreshToken:", req.RefreshToken)

	if req.RefreshToken == "" {
		fmt.Println("ERROR: Token kosong")
		return helper.Error(c, 400, "missing refresh token")
	}

	// load refresh secret
	refreshSecret := os.Getenv("JWT_REFRESH_SECRET")
	if refreshSecret == "" {
		refreshSecret = "DEFAULT_REFRESH_SECRET"
	}

	token, err := jwt.Parse(req.RefreshToken, func(t *jwt.Token) (interface{}, error) {
		return []byte(refreshSecret), nil
	})

	if err != nil {
		fmt.Println("JWT Parse Error:", err)
		return helper.Error(c, 401, "invalid or expired refresh token")
	}

	if !token.Valid {
		fmt.Println("Token INVALID")
		return helper.Error(c, 401, "invalid refresh token")
	}

	claims := token.Claims.(jwt.MapClaims)
	fmt.Println("Claims:", claims)

	userID, ok := claims["id"].(string)
	if !ok {
		fmt.Println("ERROR: Missing `id` claim")
		return helper.Error(c, 401, "invalid refresh token format")
	}

	fmt.Println("UserID dari token:", userID)

	user, err := s.users.FindByID(c.Context(), userID)
	if err != nil {
		fmt.Println("DB error:", err)
		return helper.Error(c, 401, "user not found")
	}

	fmt.Println("User ditemukan:", user.Username)

	// access token baru
	accessClaims := jwt.MapClaims{
		"id":     user.ID.String(),
		"roleId": user.RoleID.String(),
		"exp":    time.Now().Add(24 * time.Hour).Unix(),
	}

	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)

	signedAccess, err := accessToken.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		fmt.Println("ERROR ACCESS TOKEN:", err)
		return helper.Error(c, 500, "failed to generate access token")
	}

	fmt.Println("Access Token Baru:", signedAccess)

	// generate refresh token
	refreshClaims := jwt.MapClaims{
		"id":  user.ID.String(),
		"exp": time.Now().Add(7 * 24 * time.Hour).Unix(), // 7 hari
	}

	newRefresh := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	signedRefresh, err := newRefresh.SignedString([]byte(refreshSecret))
	if err != nil {
		fmt.Println("ERROR REFRESH TOKEN:", err)
		return helper.Error(c, 500, "failed to generate refresh token")
	}

	fmt.Println("Refresh Token Baru:", signedRefresh)
	fmt.Println("=== END DEBUG REFRESH TOKEN ===")

	return helper.Success(c, fiber.Map{
		"accessToken":  signedAccess,
		"refreshToken": signedRefresh,
	})
}

func (s *authService) ProfileHandler(c *fiber.Ctx) error {
	user := c.Locals("user")
	if user == nil {
		return helper.Error(c, 500, "failed to load user profile")
	}

	return helper.Success(c, user)
}

func (s *authService) Logout(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")
	if authHeader == "" {
		return helper.Error(c, 400, "missing token")
	}

	tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)

	claims := c.Locals("claims").(jwt.MapClaims)
	expUnix := int64(claims["exp"].(float64))
	exp := time.Unix(expUnix, 0)

	middleware.BlacklistToken(tokenStr, exp)

	return helper.Success(c, "logged out")
}
