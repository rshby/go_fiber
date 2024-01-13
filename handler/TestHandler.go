package handler

import (
	"fmt"
	"github.com/gofiber/fiber/v2"
	"net/http"
	"strconv"
)

type TestHandler struct {
}

// function Provider
func NewTestHandler() *TestHandler {
	return &TestHandler{}
}

func (t *TestHandler) Hello(ctx *fiber.Ctx) error {
	// get name query parameter
	var name string = ctx.Query("name", "guest")

	ctx.SendStatus(http.StatusOK)
	return ctx.JSON(map[string]any{
		"status_code": http.StatusOK,
		"message":     fmt.Sprintf("hello %v", name),
	})
}

func (t *TestHandler) RequestHandler(ctx *fiber.Ctx) error {
	// get data from header
	firstName := ctx.Get("firstname", "this")

	// get data from cookies
	lastName := ctx.Cookies("lastname", "guest")

	ctx.SendStatus(http.StatusOK)
	return ctx.JSON(map[string]any{
		"status_code": http.StatusOK,
		"message":     fmt.Sprintf("hello %v %v", firstName, lastName),
	})
}

// handler with url parameter
func (t *TestHandler) RouteParameterHandler(ctx *fiber.Ctx) error {
	// get data from url parameters
	userId, _ := strconv.Atoi(ctx.Params("userId"))
	orderId, _ := strconv.Atoi(ctx.Params("orderId"))

	ctx.SendStatus(http.StatusOK)
	return ctx.JSON(map[string]any{
		"status_code": http.StatusOK,
		"user":        userId,
		"order":       orderId,
	})
}

// handler with reuest http-form
func (t *TestHandler) RequestFormHandler(ctx *fiber.Ctx) error {
	// get name from Form
	name := ctx.FormValue("name", "guest")

	ctx.Status(http.StatusOK)
	return ctx.JSON(map[string]any{
		"status_code": http.StatusOK,
		"message":     fmt.Sprintf("hello %v", name),
	})
}
