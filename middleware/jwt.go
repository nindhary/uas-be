package middleware

import (
	"fmt"
	"os"
	"strings"
	"uas/app/repository"
	"uas/helper"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JWTMiddleware struct {
	userRepo repository.UserRepository
}

func NewJWTMiddleware(repo repository.UserRepository) *JWTMiddleware {
	return &JWTMiddleware{repo}
}

func (m *JWTMiddleware) RequireAuth(c *fiber.Ctx) error {
	authHeader := c.Get("Authorization")

	fmt.Println("=== DEBUG MIDDLEWARE ===")
	fmt.Println("Authorization Header:", authHeader)

	if authHeader == "" {
		return helper.Error(c, 401, "missing Authorization header")
	}

	tokenStr := strings.Replace(authHeader, "Bearer ", "", 1)
	fmt.Println("Extracted Token:", tokenStr)

	if IsTokenBlacklisted(tokenStr) {
		fmt.Println("TOKEN IS BLACKLISTED")
		return helper.Error(c, 401, "token already logged out")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "DEFAULT_SECRET"
	}

	token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		fmt.Println("JWT Parse Error:", err)
		return helper.Error(c, 401, "invalid or expired token")
	}

	if !token.Valid {
		fmt.Println("Token invalid")
		return helper.Error(c, 401, "invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("ERROR: invalid claims type")
		return helper.Error(c, 401, "invalid token claims")
	}

	fmt.Println("Claims:", claims)

	userID, ok := claims["id"].(string)
	if !ok {
		fmt.Println("ERROR: id claim missing or not string")
		return helper.Error(c, 401, "invalid token claims")
	}

	fmt.Println("UserID from claims:", userID)

	user, err := m.userRepo.FindByID(c.Context(), userID)
	if err != nil {
		fmt.Println("ERROR DB FindByID:", err)
		return helper.Error(c, 401, "user not found")
	}

	fmt.Println("Loaded User:", user.Username)

	c.Locals("userID", userID)
	c.Locals("user", user)
	c.Locals("claims", claims)

	return c.Next()
}
