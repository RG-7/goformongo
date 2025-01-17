package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/RG-7/goformonogo/controllers"
	"github.com/joho/godotenv"
	"github.com/julienschmidt/httprouter"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	// Access environment variables
	mongoURI := os.Getenv("MONGO_URI")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	client := getMongoClient(mongoURI)

	r := httprouter.New()

	uc := controllers.NewUserController(client)

	// Define routes
	r.GET("/user/:id", uc.GetUser)
	r.POST("/user", uc.CreateUser)
	r.DELETE("/user/:id", uc.DeleteUser)

	// Start the server
	log.Printf("Server running on port %s...", port)
	log.Fatal(http.ListenAndServe(":"+port, r))
}

func getMongoClient(uri string) *mongo.Client {
	clientOptions := options.Client().ApplyURI(uri)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		log.Fatalf("Failed to connect to MongoDB: %v", err)
	}

	// Check the connection
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Failed to ping MongoDB: %v", err)
	}

	log.Println("Connected to MongoDB!")
	return client
}
