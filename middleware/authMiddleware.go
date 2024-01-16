package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(ctx *fiber.Ctx) error {
	fmt.Println("------ masuk ke auth middleware")
	err := ctx.Next()
	fmt.Println("------ keluar dari auth middleware")

	return err
}
