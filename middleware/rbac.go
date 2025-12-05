package middleware

import (
	// "context"

	"uas/app/repository"

	"github.com/gofiber/fiber/v2"
)

type RBACMiddleware struct {
	users repository.AdminUserRepository
}

func NewRBACMiddleware(userRepo repository.AdminUserRepository) *RBACMiddleware {
	return &RBACMiddleware{users: userRepo}
}

// RequirePermission("user:manage")
func (m *RBACMiddleware) RequirePermission(perm string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userIDValue := c.Locals("userID")
		userID, ok := userIDValue.(string)
		if !ok || userID == "" {
			return c.Status(401).JSON(fiber.Map{
				"status":  "error",
				"message": "unauthorized: missing user id",
			})
		}

		perms, err := m.users.GetUserPermissions(c.Context(), userID)
		if err != nil {
			return c.Status(500).JSON(fiber.Map{
				"status":  "error",
				"message": "failed to load permissions",
			})
		}

		hasPerm := false
		for _, p := range perms {
			if p == perm {
				hasPerm = true
				break
			}
		}

		if !hasPerm {
			return c.Status(403).JSON(fiber.Map{
				"status":  "error",
				"message": "forbidden: missing permission " + perm,
			})
		}

		return c.Next()
	}
}
