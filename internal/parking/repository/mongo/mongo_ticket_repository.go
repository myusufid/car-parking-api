package mongo

import (
	"car-parking-api/internal/parking/domain"
	"car-parking-api/internal/parking/repository"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
)

type TicketRepository struct {
	collection *mongo.Collection
}

func NewTicketRepository(db *mongo.Database) repository.TicketRepository {
	return &TicketRepository{
		collection: db.Collection("tickets"),
	}
}

func (r *TicketRepository) Create(ctx context.Context, ticket *domain.Ticket) error {
	_, err := r.collection.InsertOne(ctx, ticket)
	return err
}

func (r *TicketRepository) FindByCarPlateNumber(ctx context.Context, plateNumber string) (*domain.Ticket, error) {
	var ticket domain.Ticket
	err := r.collection.FindOne(ctx, bson.M{"car_plate_number": plateNumber}).Decode(&ticket)
	return &ticket, err
}

func (r *TicketRepository) UpdateFee(ctx context.Context, ticketID primitive.ObjectID, fee int, exitTime *time.Time) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": ticketID}, bson.M{"$set": bson.M{"fee": fee, "exit_time": exitTime}})
	return err
}
