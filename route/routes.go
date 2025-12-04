package route

import (
	"uas/app/service"
	"uas/middleware"

	"github.com/gofiber/fiber/v2"
)

func RegisterRoutes(
	app *fiber.App,
	auth service.AuthService,
	adminUser service.AdminUserService,
	jwt *middleware.JWTMiddleware,
	rbac *middleware.RBACMiddleware,
) {

	api := app.Group("/app")

	// auth
	authRoute := api.Group("/auth")
	authRoute.Post("/login", auth.LoginHandler)
	authRoute.Post("/refresh", auth.RefreshHandler)
	authRoute.Get("/profile", jwt.RequireAuth, auth.ProfileHandler)

	// users
	users := api.Group("/users")

	users.Get("/", jwt.RequireAuth, rbac.RequirePermission("user:manage"), adminUser.GetAll)
	users.Get("/:id", jwt.RequireAuth, rbac.RequirePermission("user:manage"), adminUser.GetByID)
	users.Post("/", jwt.RequireAuth, rbac.RequirePermission("user:manage"), adminUser.Create)
	users.Put("/:id", jwt.RequireAuth, rbac.RequirePermission("user:manage"), adminUser.Update)
	users.Delete("/:id", jwt.RequireAuth, rbac.RequirePermission("user:manage"), adminUser.Delete)
	users.Put("/:id/role", jwt.RequireAuth, rbac.RequirePermission("user:manage"), adminUser.UpdateRole)

}
