package parking_test

import (
	"bytes"
	parking "car-parking-api/internal/parking/delivery/http"
	parkingRepoMemory "car-parking-api/internal/parking/repository/memory"
	parkingUseCase "car-parking-api/internal/parking/usecase"
	"car-parking-api/internal/route"
	"car-parking-api/test"
	"encoding/json"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

func setupParkingController() *fiber.App {
	app, _, _, _ := test.SetupIntegration()

	carRepo := parkingRepoMemory.NewCarRepository()
	parkingLotRepo := parkingRepoMemory.NewParkingSlotRepository()
	ticketRepo := parkingRepoMemory.NewTicketRepository()

	parkUseCase := parkingUseCase.NewParkingUsecase(carRepo, parkingLotRepo, ticketRepo)
	parkingController := parking.NewParkingController(parkUseCase)

	routeConfig := route.ConfigRoute{
		App:               app,
		ParkingController: parkingController,
	}
	routeConfig.Setup()
	return app
}

func TestRegisterCar(t *testing.T) {
	app := setupParkingController()

	jsonBody, err := json.Marshal(map[string]interface{}{
		"plat_nomor": "BG23102PX",
		"warna":      "hitam",
		"tipe":       "SUV",
	})

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
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

func TestExitCar(t *testing.T) {
	app := setupParkingController()

	jsonBody, err := json.Marshal(map[string]interface{}{
		"plat_nomor": "BG23102PX",
		"warna":      "hitam",
		"tipe":       "SUV",
	})

	req := httptest.NewRequest(http.MethodPost, "/register", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read the response body
	body, err := ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	fmt.Println(string(body))

	jsonBody, err = json.Marshal(map[string]interface{}{
		"plat_nomor": "BG23102PX",
	})

	req = httptest.NewRequest(http.MethodPost, "/exit", bytes.NewBuffer(jsonBody))
	req.Header.Set("Content-Type", "application/json")
	resp, err = app.Test(req)
	assert.NoError(t, err)
	defer resp.Body.Close()

	// Check the status code
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Read the response body
	body, err = ioutil.ReadAll(resp.Body)
	assert.NoError(t, err)

	fmt.Println(string(body))
}

func TestGetTotalCar(t *testing.T) {
	app := setupParkingController()

	req := httptest.NewRequest(http.MethodGet, "/total_car?tipe=SUV", nil)
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

func TestGetLicenceByColor(t *testing.T) {
	app := setupParkingController()

	req := httptest.NewRequest(http.MethodGet, "/license_by_color?warna=hitam", nil)
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
