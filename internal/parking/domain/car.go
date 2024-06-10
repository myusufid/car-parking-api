package domain

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
	"time"
)

type Car struct {
	PlateNumber string             `bson:"plate_number"`
	Color       string             `bson:"color"`
	Type        string             `bson:"type"`
	EntryTime   time.Time          `bson:"entry_time"`
	SlotID      primitive.ObjectID `bson:"slot_id"`
}
