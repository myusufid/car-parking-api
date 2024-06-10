package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Ticket struct {
	ID             primitive.ObjectID `bson:"_id,omitempty"`
	CarPlateNumber string             `bson:"car_plate_number"`
	SlotID         primitive.ObjectID `bson:"slot_id"`
	EntryTime      *time.Time         `bson:"entry_time"`
	ExitTime       *time.Time         `bson:"exit_time"`
	Fee            int                `bson:"fee"`
}
