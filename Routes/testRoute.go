package Routes

import (
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go_fiber/handler"
	"net/http"
)

func NewTestRoutes(app *fiber.App, validate *validator.Validate) {
	handler := handler.NewTestHandler(validate)

	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.JSON(map[string]any{
			"status_code": http.StatusOK,
			"message":     "success",
		})
	})

	app.Get("/hello", handler.Hello)
	app.Get("/request", handler.RequestHandler)
	app.Get("/user/:userId/order/:orderId", handler.RouteParameterHandler)
	app.Get("/hello-form", handler.RequestFormHandler)
	app.Post("/upload-file", handler.MultiPartFormHandler)
	app.Post("/login", handler.RequestBodyHandler)
	app.Post("/register", handler.RegisterUserBodyParser)
	app.Get("/response-json", handler.ResponseJsonHandler)
	app.Get("/download", handler.DownloadFile)

	v1 := app.Group("/v1")
	v1.Get("/test", handler.RoutingGroup)

	hello := app.Group("/hello")
	hello.Get("/test", handler.RoutingGroup)
}
