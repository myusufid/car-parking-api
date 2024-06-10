package memory

import (
	"car-parking-api/internal/parking/domain"
	"car-parking-api/internal/parking/repository"
	"context"
	"errors"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type TicketRepository struct {
	mu      sync.RWMutex
	tickets map[primitive.ObjectID]*domain.Ticket
}

func NewTicketRepository() repository.TicketRepository {
	return &TicketRepository{
		tickets: make(map[primitive.ObjectID]*domain.Ticket),
	}
}

func (r *TicketRepository) Create(ctx context.Context, ticket *domain.Ticket) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ticket.ID = primitive.NewObjectID()
	r.tickets[ticket.ID] = ticket
	return nil
}

func (r *TicketRepository) FindByCarPlateNumber(ctx context.Context, plateNumber string) (*domain.Ticket, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, ticket := range r.tickets {
		if ticket.CarPlateNumber == plateNumber {
			return ticket, nil
		}
	}
	return nil, errors.New("ticket not found")
}

func (r *TicketRepository) UpdateFee(ctx context.Context, ticketID primitive.ObjectID, fee int, exitTime *time.Time) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	ticket, exists := r.tickets[ticketID]
	if !exists {
		return errors.New("ticket not found")
	}
	ticket.Fee = fee
	ticket.ExitTime = exitTime
	return nil
}
