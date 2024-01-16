package middleware

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
)

func OnlyV1Middleware(ctx *fiber.Ctx) error {
	fmt.Println(" ======== masuk middleware only V1")
	err := ctx.Next()
	fmt.Println(" ======== keluar middleware only V1")
	return err
}
