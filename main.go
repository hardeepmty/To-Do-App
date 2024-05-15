package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Todo struct represents a single Todo item
type Todo struct {
	ID        string `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

var client *mongo.Client

func main() {
	// Connect to MongoDB
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	var err error
	client, err = mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	// Create a new router
	r := mux.NewRouter()

	// Enable CORS middleware
	r.Use(corsMiddleware)

	// Define routes
	r.HandleFunc("/todos", CreateTodoHandler).Methods("POST")

	// Start the HTTP server
	fmt.Println("Server is running on port 8000...")
	http.ListenAndServe(":5500", r)
}

// CreateTodoHandler handles POST requests to create a new Todo item
func CreateTodoHandler(w http.ResponseWriter, r *http.Request) {
	// Decode the request body into a Todo struct
	var todo Todo
	err := json.NewDecoder(r.Body).Decode(&todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Insert Todo item into MongoDB
	collection := client.Database("todo").Collection("todos")
	_, err = collection.InsertOne(context.Background(), todo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Send success response
	w.WriteHeader(http.StatusCreated)
}

// CORS middleware to allow requests from any origin
func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
