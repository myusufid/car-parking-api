package usecase

import (
	"car-parking-api/internal/parking/domain"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// MockCarRepository is a mock implementation of the CarRepository
type MockCarRepository struct {
	mock.Mock
}

func (m *MockCarRepository) Save(ctx context.Context, car *domain.Car) error {
	args := m.Called(ctx, car)
	return args.Error(0)
}

func (m *MockCarRepository) FindByPlateNumber(ctx context.Context, plateNumber string) (*domain.Car, error) {
	args := m.Called(ctx, plateNumber)
	return args.Get(0).(*domain.Car), args.Error(1)
}

func (m *MockCarRepository) Delete(ctx context.Context, plateNumber string) error {
	args := m.Called(ctx, plateNumber)
	return args.Error(0)
}

func (m *MockCarRepository) CountCarByType(ctx context.Context, carType string) int {
	args := m.Called(ctx, carType)
	return args.Int(0)
}

func (m *MockCarRepository) GetCarsByColor(ctx context.Context, color string) ([]string, error) {
	args := m.Called(ctx, color)
	return args.Get(0).([]string), args.Error(1)
}

// MockParkingSlotRepository is a mock implementation of the ParkingSlotRepository
type MockParkingSlotRepository struct {
	mock.Mock
}

func (m *MockParkingSlotRepository) GetAvailableSlot(ctx context.Context) (*domain.ParkingSlot, error) {
	args := m.Called(ctx)
	return args.Get(0).(*domain.ParkingSlot), args.Error(1)
}

func (m *MockParkingSlotRepository) FindSlotByID(ctx context.Context, slotID primitive.ObjectID) (*domain.ParkingSlot, error) {
	args := m.Called(ctx, slotID)
	return args.Get(0).(*domain.ParkingSlot), args.Error(1)
}

func (m *MockParkingSlotRepository) OccupySlot(ctx context.Context, slotID primitive.ObjectID) error {
	args := m.Called(ctx, slotID)
	return args.Error(0)
}

func (m *MockParkingSlotRepository) VacateSlot(ctx context.Context, slotID primitive.ObjectID) error {
	args := m.Called(ctx, slotID)
	return args.Error(0)
}

// MockTicketRepository is a mock implementation of the TicketRepository
type MockTicketRepository struct {
	mock.Mock
}

func (m *MockTicketRepository) Create(ctx context.Context, ticket *domain.Ticket) error {
	args := m.Called(ctx, ticket)
	return args.Error(0)
}

func (m *MockTicketRepository) FindByCarPlateNumber(ctx context.Context, plateNumber string) (*domain.Ticket, error) {
	args := m.Called(ctx, plateNumber)
	return args.Get(0).(*domain.Ticket), args.Error(1)
}

func (m *MockTicketRepository) UpdateFee(ctx context.Context, ticketID primitive.ObjectID, fee int, exitTime *time.Time) error {
	args := m.Called(ctx, ticketID, fee, exitTime)
	return args.Error(0)
}

func TestRegisterCar(t *testing.T) {
	ctx := context.TODO()

	mockCarRepo := new(MockCarRepository)
	mockParkingSlotRepo := new(MockParkingSlotRepository)
	mockTicketRepo := new(MockTicketRepository)

	usecase := NewParkingUsecase(mockCarRepo, mockParkingSlotRepo, mockTicketRepo)

	slotID := primitive.NewObjectID()
	mockParkingSlot := &domain.ParkingSlot{ID: slotID, SlotNumber: "A1", Occupied: false}
	mockParkingSlotRepo.On("GetAvailableSlot", ctx).Return(mockParkingSlot, nil)
	mockParkingSlotRepo.On("OccupySlot", ctx, slotID).Return(nil)
	mockParkingSlotRepo.On("FindSlotByID", ctx, slotID).Return(mockParkingSlot, nil)

	mockCarRepo.On("Save", ctx, mock.MatchedBy(func(car *domain.Car) bool {
		return car.PlateNumber == "123ABC" && car.Color == "Red" && car.Type == "SUV" && car.SlotID == slotID
	})).Return(nil)

	mockTicketRepo.On("Create", ctx, mock.MatchedBy(func(ticket *domain.Ticket) bool {
		return ticket.CarPlateNumber == "123ABC" && ticket.SlotID == slotID
	})).Return(nil)

	resultTicket, resultSlotNumber, err := usecase.RegisterCar(ctx, "123ABC", "Red", "SUV")
	assert.NoError(t, err)
	assert.NotNil(t, resultTicket)
	assert.Equal(t, "A1", resultSlotNumber)
	mockCarRepo.AssertExpectations(t)
	mockParkingSlotRepo.AssertExpectations(t)
	mockTicketRepo.AssertExpectations(t)
}

func TestExitCar(t *testing.T) {
	ctx := context.TODO()

	mockCarRepo := new(MockCarRepository)
	mockParkingSlotRepo := new(MockParkingSlotRepository)
	mockTicketRepo := new(MockTicketRepository)

	usecase := NewParkingUsecase(mockCarRepo, mockParkingSlotRepo, mockTicketRepo)

	slotID := primitive.NewObjectID()
	car := &domain.Car{PlateNumber: "123ABC", Color: "Red", Type: "SUV", EntryTime: time.Now().Add(-2 * time.Hour).UTC(), SlotID: slotID}
	mockCarRepo.On("FindByPlateNumber", ctx, "123ABC").Return(car, nil)

	ticketID := primitive.NewObjectID()
	ticket := &domain.Ticket{ID: ticketID, CarPlateNumber: car.PlateNumber, SlotID: car.SlotID, EntryTime: &car.EntryTime}
	mockTicketRepo.On("FindByCarPlateNumber", ctx, "123ABC").Return(ticket, nil)

	exitTime := time.Now().UTC()
	expectedFee := calculateParkingFee(car.Type, car.EntryTime, exitTime)

	mockTicketRepo.On("UpdateFee", ctx, ticketID, expectedFee, mock.AnythingOfType("*time.Time")).Return(nil)
	mockCarRepo.On("Delete", ctx, "123ABC").Return(nil)
	mockParkingSlotRepo.On("VacateSlot", ctx, slotID).Return(nil)

	resultTicket, err := usecase.ExitCar(ctx, "123ABC")
	assert.NoError(t, err)
	assert.NotNil(t, resultTicket)
	assert.Equal(t, expectedFee, resultTicket.Fee)

	mockCarRepo.AssertExpectations(t)
	mockParkingSlotRepo.AssertExpectations(t)
	mockTicketRepo.AssertExpectations(t)
}
