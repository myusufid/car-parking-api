package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type ParkingSlot struct {
	ID         primitive.ObjectID `bson:"_id"`
	Block      string             `bson:"block"`
	SlotNumber string             `bson:"slot_number"`
	Occupied   bool               `bson:"occupied"`
}
