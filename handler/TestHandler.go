package handler

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go_fiber/model/dto"
	"net/http"
	"strconv"
	"strings"
)

type TestHandler struct {
	Validate *validator.Validate
}

// function Provider
func NewTestHandler(validate *validator.Validate) *TestHandler {
	return &TestHandler{
		Validate: validate,
	}
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

// handler with request http-form
func (t *TestHandler) RequestFormHandler(ctx *fiber.Ctx) error {
	// get name from Form
	name := ctx.FormValue("name", "guest")

	ctx.Status(http.StatusOK)
	return ctx.JSON(map[string]any{
		"status_code": http.StatusOK,
		"message":     fmt.Sprintf("hello %v", name),
	})
}

// hander with request MultiPart Form
func (t *TestHandler) MultiPartFormHandler(ctx *fiber.Ctx) error {
	// get data from multipart form
	file, err := ctx.FormFile("file")

	// jika error ketika read file
	if err != nil {
		ctx.SendStatus(http.StatusInternalServerError)
		return ctx.JSON(map[string]any{
			"status_code": http.StatusInternalServerError,
			"message":     err.Error(),
		})
	}

	// save file to target folder
	err = ctx.SaveFile(file, "C:/Users/HP/Documents/go/src/go_fiber/multipart/target/"+file.Filename)

	// jika error ketika save file
	if err != nil {
		ctx.SendStatus(http.StatusInternalServerError)
		return ctx.JSON(map[string]any{
			"status_code": http.StatusInternalServerError,
			"message":     err.Error(),
		})
	}

	// sucess save file
	ctx.SendStatus(http.StatusOK)
	return ctx.JSON(map[string]any{
		"status_code": http.StatusOK,
		"message":     "success upload file",
	})
}

// handler with request body
func (t *TestHandler) RequestBodyHandler(ctx *fiber.Ctx) error {
	// get data from request_body
	body := ctx.Body()
	requestBody := dto.LoginRequest{}
	if err := json.Unmarshal(body, &requestBody); err != nil {
		ctx.Status(http.StatusInternalServerError)
		return ctx.JSON(&dto.ApiResponse{
			StatusCode: http.StatusInternalServerError,
			Status:     "internal server error",
			Message:    err.Error(),
		})
	}

	// validate
	if err := t.Validate.StructCtx(ctx.Context(), &requestBody); err != nil {
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			var errorMessage []string
			for _, fieldError := range validationErrors {
				msg := fmt.Sprintf("error on field [%v] with tag [%v]", fieldError.Field(), fieldError.Tag())
				errorMessage = append(errorMessage, msg)
			}

			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(&dto.ApiResponse{
				StatusCode: http.StatusBadRequest,
				Status:     "bad request",
				Message:    strings.Join(errorMessage, ", "),
			})
		}
	}

	// success get request body
	ctx.Status(http.StatusOK)
	return ctx.JSON(&dto.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     "ok",
		Message:    "success login",
		Data: map[string]string{
			"email":    requestBody.Email,
			"password": requestBody.Password,
		},
	})
}

// handelr register menggunakan Body Parser
func (t *TestHandler) RegisterUserBodyParser(ctx *fiber.Ctx) error {
	// ambil request body
	request := dto.RegisterUser{}
	if err := ctx.BodyParser(&request); err != nil {
		// error bad request
		ctx.Status(http.StatusBadRequest)
		return ctx.JSON(&dto.ApiResponse{
			StatusCode: http.StatusBadRequest,
			Status:     "bad request",
			Message:    err.Error(),
		})
	}

	// validasi
	if err := t.Validate.StructCtx(ctx.Context(), &request); err != nil {
		// error validasi
		validationErrors, ok := err.(validator.ValidationErrors)
		if ok {
			var errorMessage []string
			for _, fieldError := range validationErrors {
				message := fmt.Sprintf("error in field [%v] with tag [%v]", fieldError.Field(), fieldError.Tag())
				errorMessage = append(errorMessage, message)
			}

			// bad request
			ctx.Status(http.StatusBadRequest)
			return ctx.JSON(&dto.ApiResponse{
				StatusCode: http.StatusBadRequest,
				Status:     "bad request",
				Message:    strings.Join(errorMessage, ", "),
			})
		}
	}

	// success
	ctx.Status(http.StatusOK)
	return ctx.JSON(&dto.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     "ok",
		Message:    "success",
		Data: map[string]any{
			"username": request.Username,
			"password": request.Password,
			"name":     request.Name,
		},
	})
}

// handler HTTP Response
func (t *TestHandler) ResponseJsonHandler(ctx *fiber.Ctx) error {
	// get query params
	name := ctx.Query("name", "guest")

	// send to response
	ctx.Status(http.StatusOK)
	return ctx.JSON(&dto.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     "ok",
		Message:    "success send response json",
		Data:       fmt.Sprintf("your name is [%v]", name),
	})
}

// handler untuk download file
func (t *TestHandler) DownloadFile(ctx *fiber.Ctx) error {
	return ctx.Download("C:\\Users\\HP\\Documents\\go\\src\\go_fiber\\multipart\\source\\contoh.txt", "contoh2.txt")
}

// handler routing group
func (t *TestHandler) RoutingGroup(ctx *fiber.Ctx) error {
	return ctx.JSON(&dto.ApiResponse{
		StatusCode: http.StatusOK,
		Status:     "ok",
		Message:    "success routing group",
	})
}
