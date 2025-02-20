//checking user authentication using mongo db


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

// Define the structure for user credentials
type Credentials struct {
	Username string `json:"username" bson:"username"`
	Password string `json:"password" bson:"password"`
}

// Connect to MongoDB
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

// Login handler - Checks if user exists in MongoDB
func loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	var input Credentials


	if r.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		input.Username = r.FormValue("username")
		input.Password = r.FormValue("password")
	} else {
		// Default JSON parsing
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			http.Error(w, "Invalid JSON format", http.StatusBadRequest)
			return
		}
	}

	client, err := connectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		log.Println("Error connecting to MongoDB:", err)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("userdata").Collection("credentials")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var user Credentials
	err = collection.FindOne(ctx, bson.M{"username": input.Username, "password": input.Password}).Decode(&user)

	w.Header().Set("Content-Type", "application/json")
	if err != nil {
		json.NewEncoder(w).Encode(map[string]string{"message": "Invalid credentials "})
	} else {
		json.NewEncoder(w).Encode(map[string]string{"message": "User logged in successfully ! Welcome  " +user.Username})
	}
}

// Get all users from the database
func getUsersHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Only GET method is allowed", http.StatusMethodNotAllowed)
		return
	}

	client, err := connectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		log.Println("Error connecting to MongoDB:", err)
		return
	}
	defer client.Disconnect(context.Background())

	collection := client.Database("userdata").Collection("credentials")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch all users
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		log.Println("Error fetching users:", err)
		return
	}
	defer cursor.Close(ctx)

	var users []Credentials
	for cursor.Next(ctx) {
		var user Credentials
		if err := cursor.Decode(&user); err != nil {
			log.Println("Error decoding user:", err)
			continue
		}
		users = append(users, user)
	}

	// If no users found
	if len(users) == 0 {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(map[string]string{"message": "No users found"})
		return
	}

	// Return users as JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/users", getUsersHandler) // New GET endpoint

	fmt.Println("Server is running on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}
