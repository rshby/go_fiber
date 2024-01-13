package Routes

import (
	"github.com/gofiber/fiber/v2"
	"go_fiber/handler"
	"net/http"
)

func NewTestRoutes(app *fiber.App) {
	handler := handler.NewTestHandler()

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(map[string]any{
			"status_code": http.StatusOK,
			"message":     "success",
		})
	})

	app.Get("/hello", handler.Hello)
	app.Get("/request", handler.RequestHandler)
	app.Get("/user/:userId/order/:orderId", handler.RouteParameterHandler)
}
