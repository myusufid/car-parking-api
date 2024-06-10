package repository

import (
	"car-parking-api/internal/parking/domain"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ParkingSlotRepository interface {
	GetAvailableSlot(ctx context.Context) (*domain.ParkingSlot, error)
	OccupySlot(ctx context.Context, slotID primitive.ObjectID) error
	VacateSlot(ctx context.Context, slotID primitive.ObjectID) error
	FindSlotByID(ctx context.Context, slotID primitive.ObjectID) (*domain.ParkingSlot, error)
}
