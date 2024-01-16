package testing

import (
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"testing"
)

func TestFiberClient(t *testing.T) {
	client := fiber.AcquireClient()
	defer fiber.ReleaseClient(client)

	agent := client.Get("https://www.example.com")
	statusCode, response, err := agent.String()

	assert.Nil(t, err)
	assert.Equal(t, http.StatusOK, statusCode)
	assert.Contains(t, response, "Example Domain")
}
