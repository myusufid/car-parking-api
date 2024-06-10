package memory

import (
	"car-parking-api/internal/parking/domain"
	"car-parking-api/internal/parking/repository"
	"context"
	"errors"
	"sync"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type ParkingSlotRepository struct {
	mu    sync.RWMutex
	slots map[primitive.ObjectID]*domain.ParkingSlot
}

func NewParkingSlotRepository() repository.ParkingSlotRepository {
	repo := &ParkingSlotRepository{
		slots: make(map[primitive.ObjectID]*domain.ParkingSlot),
	}

	for i := 1; i <= 10; i++ {
		slotID := primitive.NewObjectID()
		slot := &domain.ParkingSlot{
			ID:         slotID,
			Block:      "A",
			SlotNumber: "A" + string(rune('0'+i)),
			Occupied:   false,
		}
		repo.slots[slotID] = slot
	}

	for i := 1; i <= 5; i++ {
		slotID := primitive.NewObjectID()
		slot := &domain.ParkingSlot{
			ID:         slotID,
			Block:      "B",
			SlotNumber: "B" + string(rune('0'+i)),
			Occupied:   false,
		}
		repo.slots[slotID] = slot
	}

	return repo
}

func (r *ParkingSlotRepository) GetAvailableSlot(ctx context.Context) (*domain.ParkingSlot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	for _, slot := range r.slots {
		if !slot.Occupied {
			return slot, nil
		}
	}
	return nil, errors.New("no available parking slots")
}

func (r *ParkingSlotRepository) FindSlotByID(ctx context.Context, slotID primitive.ObjectID) (*domain.ParkingSlot, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	slot, exists := r.slots[slotID]
	if !exists {
		return nil, errors.New("slot not found")
	}
	return slot, nil
}

func (r *ParkingSlotRepository) OccupySlot(ctx context.Context, slotID primitive.ObjectID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	slot, exists := r.slots[slotID]
	if !exists {
		return errors.New("slot not found")
	}
	if slot.Occupied {
		return errors.New("slot already occupied")
	}
	slot.Occupied = true
	return nil
}

func (r *ParkingSlotRepository) VacateSlot(ctx context.Context, slotID primitive.ObjectID) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	slot, exists := r.slots[slotID]
	if !exists {
		return errors.New("slot not found")
	}
	slot.Occupied = false
	return nil
}
