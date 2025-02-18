// second time commiting on this file

package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// MongoDB connection URI
const mongoURI = "mongodb://localhost:27017"
// Define a struct for user credentials
type Credentials struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

// Function to connect to MongoDB
func connectDB() (*mongo.Client, error) {
	clientOptions := options.Client().ApplyURI(mongoURI)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, clientOptions)
	if err != nil {
		return nil, err
	}
	return client, nil
}

// Function to check login credentials
func loginHandler(w http.ResponseWriter, r *http.Request) {
	// Only allow POST requests
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode incoming JSON request
	var input Credentials
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Connect to MongoDB
	client, err := connectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		log.Println("Error connecting to MongoDB:", err)
		return
	}
	defer client.Disconnect(context.Background())

	// Access database and collection
	collection := client.Database("userdata").Collection("credentials")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Search for matching credentials in MongoDB
	var user Credentials
	err = collection.FindOne(ctx, bson.M{"username": input.Username, "password": input.Password}).Decode(&user)

	// Prepare response
	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		// If credentials are not found
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid credentials"})
	} else {
		// If credentials exist
		json.NewEncoder(w).Encode(map[string]string{"message": "User logged in successfully"})
	}
}

func main() {
	// Define login API endpoint
	http.HandleFunc("/", loginHandler)

	// Start server
	fmt.Println("Server is running on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
