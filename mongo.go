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

const mongoURI ="mongodb+srv://user:qweasdzxc1@cluster0.gfhon.mongodb.net/?retryWrites=true&w=majority&appName=Cluster0"

type User struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
	Age  string `json:"age" bson:"age"`
}

// Connect to MongoDB Atlas
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

// Fetch users from MongoDB Atlas
func getUsers(w http.ResponseWriter, r *http.Request) {
	// Connect to MongoDB Atlas
	client, err := connectDB()
	if err != nil {
		http.Error(w, "Failed to connect to database", http.StatusInternalServerError)
		log.Println("Error connecting to MongoDB:", err)
		return
	}
	defer client.Disconnect(context.Background())

	// Access the "customer" database and "users" collection
	collection := client.Database("customer").Collection("users")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Fetch all documents
	cursor, err := collection.Find(ctx, bson.M{})
	if err != nil {
		http.Error(w, "Failed to retrieve data", http.StatusInternalServerError)
		log.Println("Error fetching data:", err)
		return
	}
	defer cursor.Close(ctx)

	var users []User
	for cursor.Next(ctx) {
		var user User
		if err := cursor.Decode(&user); err != nil {
			log.Println("Error decoding user:", err)
			continue
		}
		users = append(users, user)
	}

	// If no users found
	if len(users) == 0 {
		log.Println("No users found in the collection")
	}

	// Log users for debugging
	log.Println("Fetched users:", users)

	// Convert users to JSON and send the response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(users); err != nil {
		log.Println("Error encoding JSON:", err)
		http.Error(w, "Error encoding data to JSON", http.StatusInternalServerError)
	}
}


func main() {
	http.HandleFunc("/", getUsers)
	fmt.Println("Server is running on http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", nil))
}