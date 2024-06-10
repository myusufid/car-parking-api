package mongo

import (
	"car-parking-api/internal/parking/domain"
	"car-parking-api/internal/parking/repository"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type CarRepository struct {
	collection *mongo.Collection
}

func NewCarRepository(db *mongo.Database) repository.CarRepository {
	return &CarRepository{
		collection: db.Collection("cars"),
	}
}

func (r *CarRepository) Save(ctx context.Context, car *domain.Car) error {
	_, err := r.collection.InsertOne(ctx, car)
	return err
}

func (r *CarRepository) FindByPlateNumber(ctx context.Context, plateNumber string) (*domain.Car, error) {
	var car domain.Car
	err := r.collection.FindOne(ctx, bson.M{"plate_number": plateNumber}).Decode(&car)
	return &car, err
}

func (r *CarRepository) Delete(ctx context.Context, plateNumber string) error {
	_, err := r.collection.DeleteOne(ctx, bson.M{"plate_number": plateNumber})
	return err
}

func (r *CarRepository) CountCarByType(ctx context.Context, carType string) int {
	count, err := r.collection.CountDocuments(ctx, bson.M{"type": carType})
	if err != nil {
		return 0
	}
	return int(count)
}

func (r *CarRepository) GetCarsByColor(ctx context.Context, color string) ([]string, error) {
	cursor, err := r.collection.Find(ctx, bson.M{"color": color})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)

	var cars []domain.Car
	err = cursor.All(ctx, &cars)
	if err != nil {
		return nil, err
	}

	var plateNumbers []string
	for _, car := range cars {
		plateNumbers = append(plateNumbers, car.PlateNumber)
	}
	return plateNumbers, nil
}
