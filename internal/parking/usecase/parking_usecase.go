package usecase

import (
	"car-parking-api/internal/parking/domain"
	"car-parking-api/internal/parking/repository"
	"context"
	"errors"
	"time"
)

type ParkingUsecase struct {
	carRepo        repository.CarRepository
	parkingLotRepo repository.ParkingSlotRepository
	ticketRepo     repository.TicketRepository
}

func NewParkingUsecase(carRepo repository.CarRepository, parkingLotRepo repository.ParkingSlotRepository, ticketRepo repository.TicketRepository) *ParkingUsecase {
	return &ParkingUsecase{
		carRepo:        carRepo,
		parkingLotRepo: parkingLotRepo,
		ticketRepo:     ticketRepo,
	}
}

func (u *ParkingUsecase) RegisterCar(ctx context.Context, plateNumber, color, carType string) (*domain.Ticket, string, error) {
	slot, err := u.parkingLotRepo.GetAvailableSlot(ctx)
	if err != nil {
		return nil, "", err
	}

	if slot == nil {
		return nil, "", errors.New("no available parking slots")
	}

	car := &domain.Car{
		PlateNumber: plateNumber,
		Color:       color,
		Type:        carType,
		EntryTime:   time.Now().UTC(),
		SlotID:      slot.ID,
	}

	err = u.carRepo.Save(ctx, car)
	if err != nil {
		return nil, "", err
	}

	err = u.parkingLotRepo.OccupySlot(ctx, slot.ID)
	if err != nil {
		return nil, "", err
	}

	ticket := &domain.Ticket{
		CarPlateNumber: car.PlateNumber,
		SlotID:         car.SlotID,
		EntryTime:      &car.EntryTime,
	}

	err = u.ticketRepo.Create(ctx, ticket)
	if err != nil {
		return nil, "", err
	}

	parkingSlot, err := u.parkingLotRepo.FindSlotByID(ctx, ticket.SlotID)
	if err != nil {
		return nil, "", err
	}

	return ticket, parkingSlot.SlotNumber, nil
}

func (u *ParkingUsecase) ExitCar(ctx context.Context, plateNumber string) (*domain.Ticket, error) {
	car, err := u.carRepo.FindByPlateNumber(ctx, plateNumber)
	if err != nil {
		return nil, err
	}

	ticket, err := u.ticketRepo.FindByCarPlateNumber(ctx, plateNumber)
	if err != nil {
		return nil, err
	}

	exitTime := time.Now().UTC()
	fee := calculateParkingFee(car.Type, car.EntryTime, exitTime)

	ticket.ExitTime = &exitTime
	ticket.Fee = fee

	err = u.ticketRepo.UpdateFee(ctx, ticket.ID, fee, ticket.ExitTime)
	if err != nil {
		return nil, err
	}

	err = u.carRepo.Delete(ctx, plateNumber)
	if err != nil {
		return nil, err
	}

	err = u.parkingLotRepo.VacateSlot(ctx, car.SlotID)
	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func calculateParkingFee(carType string, entryTime, exitTime time.Time) int {
	var initialFee, additionalFee int
	if carType == "SUV" {
		initialFee = 25000
		additionalFee = int(float64(initialFee) * 0.20)
	} else if carType == "MPV" {
		initialFee = 35000
		additionalFee = int(float64(initialFee) * 0.20)
	}
	duration := exitTime.Sub(entryTime)
	hours := int(duration.Hours()) + 1

	if hours <= 1 {
		return initialFee
	}
	return initialFee + (additionalFee * (hours - 1))
}

func (u *ParkingUsecase) GetCarCountByType(ctx context.Context, carType string) (int, error) {
	count := u.carRepo.CountCarByType(ctx, carType)
	return count, nil
}

func (u *ParkingUsecase) GetCarsByColor(ctx context.Context, color string) ([]string, error) {
	platNumbers, _ := u.carRepo.GetCarsByColor(ctx, color)
	return platNumbers, nil
}
