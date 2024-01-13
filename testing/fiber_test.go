package testing

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go_fiber/Routes"
	handler "go_fiber/handler"
	"go_fiber/model/dto"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// test routing fiber
func TestFiberRouting(t *testing.T) {
	app := fiber.New()
	app.Get("/", func(ctx *fiber.Ctx) error {
		return ctx.SendString("hello world!")
	})

	// create request
	request := httptest.NewRequest(http.MethodGet, "/", nil)
	response, err := app.Test(request)
	assert.NotNil(t, response)
	assert.Equal(t, http.StatusOK, response.StatusCode)
	assert.Nil(t, err)

	bytes, err := io.ReadAll(response.Body)
	assert.Nil(t, err)
	assert.Equal(t, "hello world!", string(bytes))
}

func TestHello(t *testing.T) {
	app := fiber.New()
	validate := validator.New()
	handler := handler.NewTestHandler(validate)

	// register endpoint
	app.Get("/hello", handler.Hello)

	// test with query params
	t.Run("test handler success", func(t *testing.T) {
		// create request
		request, _ := http.NewRequest(http.MethodGet, "/hello", nil)

		// add query parameter
		q := request.URL.Query()
		q.Add("name", "reo")
		request.URL.RawQuery = q.Encode()

		// hit and receive response
		response, err := app.Test(request)

		assert.Nil(t, err)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// read body
		body, _ := io.ReadAll(response.Body)
		var bodyJson map[string]any
		json.Unmarshal(body, &bodyJson)
		assert.Equal(t, "hello reo", bodyJson["message"].(string))
	})

	// test without query params -> default value
	t.Run("test handler without query params", func(t *testing.T) {
		// create request
		request := httptest.NewRequest(http.MethodGet, "/hello", nil)

		// hit and receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// get response body
		body, _ := io.ReadAll(response.Body)
		bodyJson := map[string]any{}
		json.Unmarshal(body, &bodyJson)

		assert.Equal(t, "hello guest", bodyJson["message"].(string))
	})
}

// test Http Request read header and cookies
func TestGetHttpRequest(t *testing.T) {
	app := fiber.New()
	validate := validator.New()
	Routes.NewTestRoutes(app, validate)

	// test with add header and cookies
	t.Run("test with header and cookies", func(t *testing.T) {
		// create request
		request := httptest.NewRequest(http.MethodGet, "/request", nil)

		// add header and cookies
		request.Header.Add("firstname", "reo")
		request.AddCookie(&http.Cookie{
			Name:  "lastname",
			Value: "sahobby",
		})

		// hit and receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// get response body
		body, err := io.ReadAll(response.Body)
		bodyJson := map[string]any{}
		json.Unmarshal(body, &bodyJson)

		assert.Equal(t, "hello reo sahobby", bodyJson["message"].(string))
	})

	// test without header and cookies -> default value
	t.Run("test without header and cookies", func(t *testing.T) {
		// create request
		request := httptest.NewRequest(http.MethodGet, "/request", nil)

		// hit and receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// get response body
		body, err := io.ReadAll(response.Body)
		bodyJson := map[string]any{}
		json.Unmarshal(body, &bodyJson)

		assert.Equal(t, "hello this guest", bodyJson["message"].(string))
	})
}

// test with URL parameter
func TestGetValueURLParams(t *testing.T) {
	app := fiber.New()
	validate := validator.New()
	Routes.NewTestRoutes(app, validate)

	t.Run("test with url parameter", func(t *testing.T) {
		// create request
		request := httptest.NewRequest(http.MethodGet, "/user/1/order/2", nil)

		// hit and receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// get response body
		body, _ := io.ReadAll(response.Body)
		bodyJson := map[string]any{}
		json.Unmarshal(body, &bodyJson)

		assert.Equal(t, 1, int(bodyJson["user"].(float64)))
		assert.Equal(t, 2, int(bodyJson["order"].(float64)))
	})
}

func TestFormParameter(t *testing.T) {
	app := fiber.New()
	validate := validator.New()
	Routes.NewTestRoutes(app, validate)

	t.Run("test with form parameter", func(t *testing.T) {
		// create request
		requestBody := strings.NewReader("name=reo")
		request := httptest.NewRequest(http.MethodGet, "/hello-form", requestBody)
		request.Header.Add("content-type", "application/x-www-form-urlencoded")

		// hit and receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// get response body
		body, _ := io.ReadAll(response.Body)
		bodyJson := map[string]any{}
		json.Unmarshal(body, &bodyJson)

		assert.Equal(t, "hello reo", bodyJson["message"].(string))
	})

	// test without parameter form -> using default value
	t.Run("test without form params", func(t *testing.T) {
		// create request
		request := httptest.NewRequest(http.MethodGet, "/hello-form", nil)

		// hit and receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// get response body
		body, _ := io.ReadAll(response.Body)
		bodyJson := map[string]any{}
		json.Unmarshal(body, &bodyJson)

		assert.Equal(t, "hello guest", bodyJson["message"].(string))
	})
}

// test hit endpoint menggunakan request_body
func TestRequestBody(t *testing.T) {
	app := fiber.New()
	validate := validator.New()
	Routes.NewTestRoutes(app, validate)

	// test success
	t.Run("test request body login success", func(t *testing.T) {
		// create request body
		req := map[string]any{
			"email":    "reoshby@gmail.com",
			"password": "123456",
		}

		reqJson, err := json.Marshal(&req)
		if err != nil {
			fmt.Println(err)
			t.Fail()
			return
		}
		requestBody := strings.NewReader(string(reqJson))

		// create request
		request := httptest.NewRequest(http.MethodPost, "/login", requestBody)
		request.Header.Add("content-type", "application/json")

		// hit and receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// get response body
		all, _ := io.ReadAll(response.Body)
		responseBody := map[string]any{}
		json.Unmarshal(all, &responseBody)

		assert.Equal(t, req["email"], responseBody["data"].(map[string]any)["email"].(string))
		assert.Equal(t, req["password"], responseBody["data"].(map[string]any)["password"].(string))
	})

	t.Run("test request body login failed", func(t *testing.T) {
		// create request
		req := dto.LoginRequest{
			Email:    "reoshby",
			Password: "12qw",
		}

		marshal, err := json.Marshal(&req)
		if err != nil {
			fmt.Println(err)
			t.Fail()
			return
		}

		requestBody := strings.NewReader(string(marshal))

		// create request
		request := httptest.NewRequest(http.MethodPost, "/login", requestBody)
		request.Header.Add("content-type", "application/json")

		// hit and receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})
}
