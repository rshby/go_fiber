package handler

import (
	"github.com/gofiber/fiber/v2"
	"go_fiber/model/dto"
	"net/http"
)

type ErrorHandler struct{}

// function provider
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{}
}

// method error
func (e *ErrorHandler) ErrorHandler(ctx *fiber.Ctx, err error) error {
	ctx.Status(http.StatusInternalServerError)
	return ctx.JSON(&dto.ApiResponse{
		StatusCode: http.StatusInternalServerError,
		Status:     "internal server error",
		Message:    err.Error(),
	})
}
