package http

import (
	"car-parking-api/internal/parking/usecase"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"time"
)

type ParkingController struct {
	parkingUsecase *usecase.ParkingUsecase
}

func NewParkingController(parkingUsecase *usecase.ParkingUsecase) *ParkingController {
	return &ParkingController{
		parkingUsecase: parkingUsecase,
	}
}

type RegisterCarResponse struct {
	PlatNomor    string `json:"plat_nomor"`
	ParkingLot   string `json:"parking_lot"`
	TanggalMasuk string `json:"tanggal_masuk"`
}

func (h *ParkingController) RegisterCar(c *fiber.Ctx) error {
	var req struct {
		PlateNumber string `json:"plat_nomor"`
		Color       string `json:"warna"`
		Type        string `json:"tipe"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ticket, slotNumber, err := h.parkingUsecase.RegisterCar(c.Context(), req.PlateNumber, req.Color, req.Type)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		fmt.Println("Error loading location:", err)
	}

	response := RegisterCarResponse{
		PlatNomor:    ticket.CarPlateNumber,
		ParkingLot:   slotNumber,
		TanggalMasuk: ticket.EntryTime.In(loc).Format("2006-01-02 15:04"),
	}

	return c.JSON(response)
}

type ExitCarResponse struct {
	PlatNomor     string `json:"plat_nomor"`
	TanggalMasuk  string `json:"tanggal_masuk"`
	TanggalKeluar string `json:"tanggal_keluar"`
	JumlahBayar   int    `json:"jumlah_bayar"`
}

func (h *ParkingController) ExitCar(c *fiber.Ctx) error {
	var req struct {
		PlateNumber string `json:"plat_nomor"`
	}

	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	ticket, err := h.parkingUsecase.ExitCar(c.Context(), req.PlateNumber)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	loc, err := time.LoadLocation("Asia/Jakarta")
	if err != nil {
		fmt.Println("Error loading location:", err)
	}

	response := ExitCarResponse{
		PlatNomor:     ticket.CarPlateNumber,
		TanggalMasuk:  ticket.EntryTime.In(loc).Format("2006-01-02 15:04"),
		TanggalKeluar: ticket.ExitTime.In(loc).Format("2006-01-02 15:04"),
		JumlahBayar:   ticket.Fee,
	}

	return c.JSON(response)
}

func (h *ParkingController) GetCarCountByType(c *fiber.Ctx) error {
	carType := c.Query("tipe", "")
	count, _ := h.parkingUsecase.GetCarCountByType(c.Context(), carType)

	return c.JSON(fiber.Map{
		"jumlah_kendaraan": count,
	})
}

func (h *ParkingController) GetCarsByColor(c *fiber.Ctx) error {
	carColor := c.Query("warna", "")
	licenseNumbers, _ := h.parkingUsecase.GetCarsByColor(c.Context(), carColor)

	return c.JSON(fiber.Map{
		"plat_nomor": licenseNumbers,
	})
}
