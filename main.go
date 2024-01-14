package main

import (
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/template/mustache/v2"
	"github.com/spf13/viper"
	"go_fiber/Routes"
	"go_fiber/handler"
	"log"
	"time"
)

func main() {
	// load env
	var config *viper.Viper = viper.New()
	config.SetConfigFile("config.json")
	config.AddConfigPath("./")
	if err := config.ReadInConfig(); err != nil {
		log.Fatalf("error cant load config.json")
	}

	// instance validate
	validate := validator.New()
	errorHandler := handler.NewErrorHandler()

	// engine template mustache
	engineView := mustache.New("./view", ".mustache")

	// create instance app fiber
	app := fiber.New(fiber.Config{
		IdleTimeout:  3 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		Prefork:      true,
		Views:        engineView,
		ErrorHandler: errorHandler.ErrorHandler, // override default error handler
	})

	if fiber.IsChild() {
		log.Println("im child process")
	} else {
		log.Println("im parent process")
	}

	// use logger to log HTTP request
	app.Use(logger.New())

	app.Get("/test", func(ctx *fiber.Ctx) error {
		return fiber.NewError(500, "error internal")
	})

	// routes
	Routes.NewTestRoutes(app, validate)

	addr := fmt.Sprintf("%v:%v", config.GetString("app.host"), config.GetString("app.port"))
	if err := app.Listen(addr); err != nil {
		log.Fatalf(err.Error())
	}
}
