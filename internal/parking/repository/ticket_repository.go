package repository

import (
	"car-parking-api/internal/parking/domain"
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type TicketRepository interface {
	Create(ctx context.Context, ticket *domain.Ticket) error
	FindByCarPlateNumber(ctx context.Context, plateNumber string) (*domain.Ticket, error)
	UpdateFee(ctx context.Context, ticketID primitive.ObjectID, fee int, exitTime *time.Time) error
}
