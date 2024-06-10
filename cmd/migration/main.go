package main

import (
	"car-parking-api/config"
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/mongodb"
	"github.com/golang-migrate/migrate/v4/source/file"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

func main() {
	// Load configuration
	viperConfig := config.NewViper()

	// Get MongoDB URI and database name from configuration
	user := viperConfig.GetString("DATABASE_USER")
	password := viperConfig.GetString("DATABASE_PASSWORD")
	host := viperConfig.GetString("DATABASE_HOST")
	port := viperConfig.GetString("DATABASE_PORT")
	name := viperConfig.GetString("DATABASE_NAME")
	authSource := "admin"

	uri := fmt.Sprintf("mongodb://%s:%s@%s:%s/%s?authSource=%s",
		user, password, host, port, name, authSource)

	fmt.Println(uri, name)

	// Create a new MongoDB client
	mongoClient, err := NewMongoClient(uri)
	if err != nil {
		log.Fatalf("Failed to create MongoDB client: %v", err)
	}

	// Run migrations
	err = runMigrations(mongoClient, name, "./migrations")
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}
}

func runMigrations(client *mongo.Client, dbName, migrationsPath string) error {
	// Ensure the client is connected
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := client.Ping(ctx, nil); err != nil {
		return err
	}

	// Create a MongoDB driver
	driver, err := mongodb.WithInstance(client, &mongodb.Config{DatabaseName: dbName})
	if err != nil {
		return err
	}

	// Create a file source for the migrations
	sourceDriver, err := (&file.File{}).Open(migrationsPath)
	if err != nil {
		return err
	}

	// Create a new migrate instance
	m, err := migrate.NewWithInstance("file", sourceDriver, dbName, driver)
	if err != nil {
		return err
	}

	// Run the migrations
	if err := m.Up(); err != nil && err != migrate.ErrNoChange {
		return err
	}

	log.Println("Migrations ran successfully")
	return nil
}

func NewMongoClient(uri string) (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(uri)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatal("Failed to connect to MongoDB: ", err)
		return nil, err
	}

	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("Failed to ping MongoDB: ", err)
		return nil, err
	}

	log.Println("Connected to MongoDB")
	return client, nil
}
