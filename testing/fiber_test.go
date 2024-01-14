package testing

import (
	"encoding/json"
	"fmt"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/mustache/v2"
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

// test initial
func TestInitial(t *testing.T) {
	app := fiber.New(fiber.Config{
		Prefork: true,
	})
	validate := validator.New()
	Routes.NewTestRoutes(app, validate)

	// test endpoint /
	t.Run("test initial endpoint", func(t *testing.T) {
		// create request
		request := httptest.NewRequest(http.MethodGet, "/", nil)

		// hit and receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// get response body
		body, _ := io.ReadAll(response.Body)
		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, "success", responseBody["message"].(string))
	})
}

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

// test endpoint menggunakan body parser otomatis
func TestBodyParserRequest(t *testing.T) {
	app := fiber.New()
	validate := validator.New()
	Routes.NewTestRoutes(app, validate)

	// test menggunakan json
	t.Run("test with json request body", func(t *testing.T) {
		// create request_body
		req := dto.RegisterUser{
			Username: "reoshby@gmail.com",
			Password: "123456",
			Name:     "Reo Sahobby",
		}
		marshal, err := json.Marshal(&req)
		if err != nil {
			fmt.Println(err.Error())
			t.Fail()
			return
		}
		requestBody := strings.NewReader(string(marshal))

		// create HTTP Request
		request := httptest.NewRequest(http.MethodPost, "/register", requestBody)
		request.Header.Add("content-type", "application/json")

		// hit and get response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// get request body
		body, _ := io.ReadAll(response.Body)
		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusOK, int(responseBody["status_code"].(float64)))
		assert.Equal(t, req.Username, responseBody["data"].(map[string]any)["username"].(string))
		assert.Equal(t, req.Password, responseBody["data"].(map[string]any)["password"].(string))
		assert.Equal(t, req.Name, responseBody["data"].(map[string]any)["name"].(string))
	})

	// test menggunakan form request
	t.Run("test with json form body", func(t *testing.T) {
		// create request
		requestForm := strings.NewReader("username=reoshby@gmail.com&password=123456&name=Reo Sahobby")

		// create request
		request := httptest.NewRequest(http.MethodPost, "/register", requestForm)
		request.Header.Add("content-type", "application/x-www-form-urlencoded")

		// hit and receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// get response_body
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println(err.Error())
			t.Fail()
			return
		}
		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusOK, int(responseBody["status_code"].(float64)))
		assert.Equal(t, "reoshby@gmail.com", responseBody["data"].(map[string]any)["username"].(string))
		assert.Equal(t, "123456", responseBody["data"].(map[string]any)["password"].(string))
		assert.Equal(t, "Reo Sahobby", responseBody["data"].(map[string]any)["name"].(string))
	})

	// test bad request parsing
	t.Run("test with json failed parsing", func(t *testing.T) {
		// create request_body
		requestBody := strings.NewReader(`name=reo`)

		// create request
		request := httptest.NewRequest(http.MethodPost, "/register", requestBody)
		request.Header.Add("content-type", "application/json")

		// hit and receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})

	// test failed bad request
	t.Run("test with json bad request", func(t *testing.T) {
		// create request_body
		req := dto.RegisterUser{
			Username: "reo",
		}
		marshal, err := json.Marshal(&req)
		if err != nil {
			fmt.Println(err.Error())
			t.Fail()
			return
		}
		requestBody := strings.NewReader(string(marshal))

		// create request
		request := httptest.NewRequest(http.MethodPost, "/register", requestBody)
		request.Header.Add("content-type", "application/json")

		// hit and receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusBadRequest, response.StatusCode)
	})
}

// test handler response-json
func TestResponseJson(t *testing.T) {
	app := fiber.New()
	validate := validator.New()
	Routes.NewTestRoutes(app, validate)

	// test with query parameters
	t.Run("response json with parameter", func(t *testing.T) {
		// create request
		request, _ := http.NewRequest(http.MethodGet, "/response-json", nil)
		q := request.URL.Query()
		q.Add("name", "reo")
		request.URL.RawQuery = q.Encode()

		// hit and receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// get response body
		body, _ := io.ReadAll(response.Body)
		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusOK, int(responseBody["status_code"].(float64)))
		assert.Equal(t, "your name is [reo]", responseBody["data"].(string))
	})

	// test without query parameter -> default value
	t.Run("response json without parameter", func(t *testing.T) {
		// create request
		request := httptest.NewRequest(http.MethodGet, "/response-json", nil)

		// hit and receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// get response body
		body, _ := io.ReadAll(response.Body)
		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusOK, int(responseBody["status_code"].(float64)))
		assert.Equal(t, "your name is [guest]", responseBody["data"].(string))
	})
}

// test handler download
func TestDownloadFile(t *testing.T) {
	app := fiber.New()
	validate := validator.New()
	Routes.NewTestRoutes(app, validate)

	// test download file
	t.Run("test download file", func(t *testing.T) {
		// create request
		request := httptest.NewRequest(http.MethodGet, "/download", nil)

		// hit and receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)
		assert.Equal(t, "attachment; filename=\"contoh2.txt\"", response.Header.Get("Content-Disposition"))

		// get response body
		body, err := io.ReadAll(response.Body)
		assert.Nil(t, err)
		assert.Contains(t, "test\r\noke", string(body))
	})
}

// test handler routing group
func TestRoutingGroup(t *testing.T) {
	app := fiber.New()
	validate := validator.New()
	Routes.NewTestRoutes(app, validate)

	// test routing v1
	t.Run("test routing v1", func(t *testing.T) {
		// create request
		request := httptest.NewRequest(http.MethodGet, "/v1/test", nil)

		// hit and receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// get response body
		body, _ := io.ReadAll(response.Body)
		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, "success routing group", responseBody["message"].(string))
	})

	// test routing hello
	t.Run("test routing hello", func(t *testing.T) {
		// create request
		request := httptest.NewRequest(http.MethodGet, "/hello/test", nil)

		// hit and receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// get response body
		body, _ := io.ReadAll(response.Body)
		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, "success routing group", responseBody["message"].(string))
	})
}

// test endpoint file static
func TestEndpointStatic(t *testing.T) {
	app := fiber.New()
	validate := validator.New()
	Routes.NewTestRoutes(app, validate)

	// test access static file
	t.Run("test static success", func(t *testing.T) {
		// create request
		request := httptest.NewRequest(http.MethodGet, "/public/contoh.txt", nil)

		// hit and receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// get responsebody
		body, err := io.ReadAll(response.Body)
		assert.Nil(t, err)
		assert.Contains(t, "test\r\noke", string(body))
	})

	// test access not found
	t.Run("test static not found", func(t *testing.T) {
		// create request
		request := httptest.NewRequest(http.MethodGet, "/public", nil)

		// hit and receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusNotFound, response.StatusCode)
	})
}

// test default error handler
func TestDefaultErrorHandler(t *testing.T) {
	errorHandler := handler.NewErrorHandler()
	app := fiber.New(fiber.Config{
		Prefork:      true,
		ErrorHandler: errorHandler.ErrorHandler,
	})

	validate := validator.New()
	Routes.NewTestRoutes(app, validate)

	// test internal server error
	t.Run("test error handler internal server error", func(t *testing.T) {
		// create request
		request := httptest.NewRequest(http.MethodGet, "/v1/wasd", nil)

		// hit and receive response
		respoonse, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, respoonse)
		assert.Equal(t, http.StatusInternalServerError, respoonse.StatusCode)

		// get response body
		body, _ := io.ReadAll(respoonse.Body)
		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Equal(t, http.StatusInternalServerError, int(responseBody["status_code"].(float64)))
		assert.Contains(t, responseBody["message"].(string), "Cannot")
	})
}

// test render template mustache
func TestRenderTemplateView(t *testing.T) {
	engine := mustache.New("C:/Users/HP/Documents/go/src/go_fiber/view", ".mustache")
	app := fiber.New(fiber.Config{Prefork: true, Views: engine})
	validate := validator.New()
	Routes.NewTestRoutes(app, validate)

	// test endpoint render template
	t.Run("test render template view", func(t *testing.T) {
		// create request
		request := httptest.NewRequest(http.MethodGet, "/v1/view", nil)

		// hit and receive response
		response, err := app.Test(request)
		assert.Nil(t, err)
		assert.NotNil(t, response)
		assert.Equal(t, http.StatusOK, response.StatusCode)

		// get response body
		body, _ := io.ReadAll(response.Body)
		responseBody := map[string]any{}
		json.Unmarshal(body, &responseBody)

		assert.Contains(t, string(body), "GOlang")
		assert.Contains(t, string(body), "web")
	})
}
