package memory

import (
	"car-parking-api/internal/parking/domain"
	"car-parking-api/internal/parking/repository"
	"context"
	"errors"
	"sync"
)

type CarRepository struct {
	mu    sync.RWMutex
	cars  map[string]*domain.Car
	color map[string][]string
}

func NewCarRepository() repository.CarRepository {
	return &CarRepository{
		cars:  make(map[string]*domain.Car),
		color: make(map[string][]string),
	}
}

func (r *CarRepository) Save(ctx context.Context, car *domain.Car) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.cars[car.PlateNumber] = car
	r.color[car.Color] = append(r.color[car.Color], car.PlateNumber)
	return nil
}

func (r *CarRepository) FindByPlateNumber(ctx context.Context, plateNumber string) (*domain.Car, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	car, exists := r.cars[plateNumber]
	if !exists {
		return nil, errors.New("car not found")
	}
	return car, nil
}

func (r *CarRepository) Delete(ctx context.Context, plateNumber string) error {
	r.mu.Lock()
	defer r.mu.Unlock()
	car, exists := r.cars[plateNumber]
	if !exists {
		return errors.New("car not found")
	}
	delete(r.cars, plateNumber)

	// Remove car from the color index
	for i, pn := range r.color[car.Color] {
		if pn == plateNumber {
			r.color[car.Color] = append(r.color[car.Color][:i], r.color[car.Color][i+1:]...)
			break
		}
	}

	return nil
}

func (r *CarRepository) CountCarByType(ctx context.Context, carType string) int {
	r.mu.RLock()
	defer r.mu.RUnlock()
	count := 0
	for _, car := range r.cars {
		if car.Type == carType {
			count++
		}
	}
	return count
}

func (r *CarRepository) GetCarsByColor(ctx context.Context, color string) ([]string, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()
	plateNumbers, exists := r.color[color]
	if !exists {
		return nil, errors.New("no cars found for this color")
	}
	return plateNumbers, nil
}
