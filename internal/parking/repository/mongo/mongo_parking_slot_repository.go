package mongo

import (
	"car-parking-api/internal/parking/domain"
	"car-parking-api/internal/parking/repository"
	"context"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type ParkingSlotRepository struct {
	collection *mongo.Collection
}

func NewParkingSlotRepository(db *mongo.Database) repository.ParkingSlotRepository {
	return &ParkingSlotRepository{
		collection: db.Collection("parking_slots"),
	}
}

func (r *ParkingSlotRepository) GetAvailableSlot(ctx context.Context) (*domain.ParkingSlot, error) {
	var slot domain.ParkingSlot
	err := r.collection.FindOne(ctx, bson.M{"occupied": false}).Decode(&slot)
	return &slot, err
}

func (r *ParkingSlotRepository) FindSlotByID(ctx context.Context, slotID primitive.ObjectID) (*domain.ParkingSlot, error) {
	var slot domain.ParkingSlot
	err := r.collection.FindOne(ctx, bson.M{"_id": slotID}).Decode(&slot)
	return &slot, err
}

func (r *ParkingSlotRepository) OccupySlot(ctx context.Context, slotID primitive.ObjectID) error {
	result, err := r.collection.UpdateOne(ctx, bson.M{"_id": slotID}, bson.M{"$set": bson.M{"occupied": true}})
	if err != nil {
		fmt.Println("Error updating document:", err)
		return err
	}
	fmt.Printf("Matched %v documents and updated %v documents.\n", result.MatchedCount, result.ModifiedCount)
	return err
}

func (r *ParkingSlotRepository) VacateSlot(ctx context.Context, slotID primitive.ObjectID) error {
	_, err := r.collection.UpdateOne(ctx, bson.M{"_id": slotID}, bson.M{"$set": bson.M{"occupied": false}})
	return err
}
