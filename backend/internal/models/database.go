package models

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Global database connections
var (
	PgPool *pgxpool.Pool
	MongoDB *mongo.Database
	MongoClient *mongo.Client
)

// InitDatabases initializes all database connections
func InitDatabases() error {
	var err error

	// Initialize PostgreSQL connection
	if err = initPostgres(); err != nil {
		return fmt.Errorf("failed to initialize PostgreSQL: %w", err)
	}

	// Initialize MongoDB connection
	if err = initMongoDB(); err != nil {
		return fmt.Errorf("failed to initialize MongoDB: %w", err)
	}

	return nil
}

// initPostgres establishes connection to PostgreSQL
func initPostgres() error {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbName := os.Getenv("DB_NAME")

	connString := fmt.Sprintf("postgres://%s:%s@%s:5432/%s", 
		dbUser, dbPassword, dbHost, dbName)

	config, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return fmt.Errorf("unable to parse connection string: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	PgPool, err = pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("unable to create connection pool: %w", err)
	}

	// Test the connection
	if err := PgPool.Ping(ctx); err != nil {
		return fmt.Errorf("could not ping database: %w", err)
	}

	log.Println("Successfully connected to PostgreSQL")
	return nil
}

// initMongoDB establishes connection to MongoDB
func initMongoDB() error {
	mongoURI := os.Getenv("MONGO_URI")
	
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Create a MongoDB client
	clientOptions := options.Client().ApplyURI(mongoURI)
	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return fmt.Errorf("failed to connect to MongoDB: %w", err)
	}

	// Test the connection
	if err := client.Ping(ctx, nil); err != nil {
		return fmt.Errorf("could not ping MongoDB: %w", err)
	}

	// Set the global variables
	MongoClient = client
	MongoDB = client.Database("trendscout")

	log.Println("Successfully connected to MongoDB")
	return nil
}

// CloseDatabases closes all database connections
func CloseDatabases() {
	// Close PostgreSQL connection
	if PgPool != nil {
		PgPool.Close()
		log.Println("PostgreSQL connection closed")
	}

	// Close MongoDB connection
	if MongoClient != nil {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		
		if err := MongoClient.Disconnect(ctx); err != nil {
			log.Printf("Error disconnecting from MongoDB: %v", err)
		} else {
			log.Println("MongoDB connection closed")
		}
	}
} 