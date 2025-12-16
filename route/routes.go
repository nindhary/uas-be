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
	studentAch service.StudentAchievementService,
	lecturerAch service.LecturerAchievementService,
	studentSvc service.StudentService,
	lecturerSvc service.LecturerService,
	adminAchievementSvc service.AdminAchievementService,
	reportSvc service.ReportService,
) {

	api := app.Group("/app")

	// auth
	authRoute := api.Group("/auth")
	authRoute.Post("/login", auth.LoginHandler)
	authRoute.Post("/refresh", auth.RefreshHandler)
	authRoute.Get("/profile", jwt.RequireAuth, auth.ProfileHandler)
	authRoute.Post("/logout", jwt.RequireAuth, auth.Logout)

	// users
	users := api.Group("/users")

	users.Get("/", jwt.RequireAuth, rbac.RequirePermission("user:manage"), adminUser.GetAll)
	users.Get("/:id", jwt.RequireAuth, rbac.RequirePermission("user:manage"), adminUser.GetByID)
	users.Post("/", jwt.RequireAuth, rbac.RequirePermission("user:manage"), adminUser.Create)
	users.Put("/:id", jwt.RequireAuth, rbac.RequirePermission("user:manage"), adminUser.Update)
	users.Delete("/:id", jwt.RequireAuth, rbac.RequirePermission("user:manage"), adminUser.Delete)
	users.Put("/:id/role", jwt.RequireAuth, rbac.RequirePermission("user:manage"), adminUser.UpdateRole)

	// achievements
	achievement := api.Group("student/achievements", jwt.RequireAuth)

	achievement.Get("/", rbac.RequirePermission("achievement:create"), studentAch.GetMyAchievements)
	achievement.Post("/", rbac.RequirePermission("achievement:create"), studentAch.Create)
	achievement.Get("/:id", rbac.RequirePermission("achievement:create"), studentAch.GetDetail)
	achievement.Get("/:id/history", rbac.RequirePermission("achievement:create"), studentAch.GetHistory)
	achievement.Put("/:id", rbac.RequirePermission("achievement:update"), studentAch.Update)
	achievement.Delete("/:id", rbac.RequirePermission("achievement:delete"), studentAch.Delete)
	achievement.Post("/:id/submit", rbac.RequirePermission("achievement:submit"), studentAch.Submit)
	achievement.Post("/:id/attachments", rbac.RequirePermission("achievement:upload"), studentAch.UploadAttachment)

	lecturer := api.Group("lecturer/achievements", jwt.RequireAuth)
	lecturer.Get("/", rbac.RequirePermission("achievement:read_advisee"), lecturerAch.GetAdviseeAchievements)
	lecturer.Get("/:id", rbac.RequirePermission("achievement:read_advisee"), lecturerAch.GetDetail)
	lecturer.Get("/:id/history", rbac.RequirePermission("achievement:read_advisee"), lecturerAch.GetHistory)
	lecturer.Post("/:id/verify", rbac.RequirePermission("achievement:verify"), lecturerAch.Verify)
	lecturer.Post("/:id/reject", rbac.RequirePermission("achievement:reject"), lecturerAch.Reject)

	// students
	students := api.Group("/students", jwt.RequireAuth)

	students.Get("/", studentSvc.GetAll)
	students.Get("/:id", studentSvc.GetDetail)
	students.Get("/:id/achievements", studentSvc.GetMyAchievements)
	students.Put("/:id/advisor", jwt.RequireAuth, rbac.RequirePermission("user:manage"), studentSvc.UpdateAdvisor)
	students.Post("/profile", jwt.RequireAuth, rbac.RequirePermission("user:manage"), studentSvc.CreateProfile)

	// lecturers
	lecturers := api.Group("/lecturers", jwt.RequireAuth)

	lecturers.Get("/", lecturerSvc.GetAll)
	lecturers.Get("/:id/advisees", rbac.RequirePermission("achievement:read_advisee"), lecturerSvc.GetAdvisees)
	lecturers.Post("/profile", jwt.RequireAuth, rbac.RequirePermission("user:manage"), lecturerSvc.CreateProfile)

	// admin
	admin := api.Group("/admin", jwt.RequireAuth)
	admin.Get("/achievements", rbac.RequirePermission("user:manage"), adminAchievementSvc.GetAll)

	// reports
	reports := api.Group("/reports", jwt.RequireAuth)

	reports.Get("/statistics", rbac.RequirePermission("user:manage"), reportSvc.GetStatistics)
	reports.Get("/student/:id", reportSvc.GetStudentStatistics)

}
