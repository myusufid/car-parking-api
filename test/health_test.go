package test

import (
	"car-parking-api/config"
	health "car-parking-api/internal/health/delivery/http"
	"car-parking-api/internal/route"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func SetupTestHealth() *fiber.App {
	app := config.NewFiber()

	routeConfig := route.ConfigRoute{
		App:              app,
		HealthController: health.NewHealthController(),
	}
	routeConfig.Setup()
	return app
}

func TestGetHealth(t *testing.T) {
	app := SetupTestHealth()

	req := httptest.NewRequest(http.MethodGet, "/health", nil)
	resp, err := app.Test(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	fmt.Println(string(body))
}
