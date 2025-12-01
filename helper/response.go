package helper

import "github.com/gofiber/fiber/v2"

func Success(c *fiber.Ctx, data any) error {
	return c.Status(200).JSON(fiber.Map{
		"status": "success",
		"data":   data,
	})
}

func Error(c *fiber.Ctx, code int, msg string) error {
	return c.Status(code).JSON(fiber.Map{
		"status":  "error",
		"message": msg,
	})
}
