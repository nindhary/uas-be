package middleware

import (
	"os"
	"strings"
	"uas/app/repository"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

type JWTMiddleware struct {
	userRepo repository.UserRepository
}

func NewJWTMiddleware(repo repository.UserRepository) *JWTMiddleware {
	return &JWTMiddleware{repo}
}

func (mw *JWTMiddleware) RequireAuth(c *fiber.Ctx) error {
	auth := c.Get("Authorization")
	if auth == "" {
		return c.Status(401).JSON(fiber.Map{"message": "unauthorized"})
	}

	tokenString := strings.TrimPrefix(auth, "Bearer ")

	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("JWT_SECRET")), nil
	})

	if err != nil || !token.Valid {
		return c.Status(401).JSON(fiber.Map{"message": "invalid token"})
	}

	claims := token.Claims.(jwt.MapClaims)

	c.Locals("user_id", claims["id"])
	c.Locals("role", claims["role"])

	return c.Next()
}
