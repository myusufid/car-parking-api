package config

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func NewDatabase(viper *viper.Viper) *mongo.Database {
	user := viper.GetString("DATABASE_USER")
	password := viper.GetString("DATABASE_PASSWORD")
	host := viper.GetString("DATABASE_HOST")
	port := viper.GetString("DATABASE_PORT")
	name := viper.GetString("DATABASE_NAME")
	authSource := "admin"

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s",
		user, password, host, port, name, authSource)

	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB: ", err)
	}

	log.Println("Connected to MongoDB")

	return client.Database(name)
}
