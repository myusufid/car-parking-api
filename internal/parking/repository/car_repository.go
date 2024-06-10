package repository

import (
	"car-parking-api/internal/parking/domain"
	"context"
)

type CarRepository interface {
	Save(ctx context.Context, car *domain.Car) error
	FindByPlateNumber(ctx context.Context, plateNumber string) (*domain.Car, error)
	Delete(ctx context.Context, plateNumber string) error
	CountCarByType(ctx context.Context, carType string) int
	GetCarsByColor(ctx context.Context, color string) ([]string, error)
}
