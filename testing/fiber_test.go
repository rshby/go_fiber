package testing

import (
	"encoding/json"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"go_fiber/Routes"
	handler "go_fiber/handler"
	"io"
	"net/http"
	"net/http/httptest"
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
	handler := handler.NewTestHandler()

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
	Routes.NewTestRoutes(app)

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
	Routes.NewTestRoutes(app)

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
