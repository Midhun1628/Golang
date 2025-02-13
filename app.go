package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)


type Message struct {
	Greeting string `json:"greetings"`

}

func handler(w http.ResponseWriter, r *http.Request) {

	currentHour := time.Now().Hour()

	var greetingss string
	if currentHour < 12 {
		greetingss = "Good Morning"
	} else {
		greetingss = "Good Afternoon"
	}

message := Message{Greeting: greetingss}

	
	jsonMessage, err := json.Marshal(message)
	if err != nil {
		http.Error(w, "Error generating JSON", http.StatusInternalServerError)
		return
	}

	
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonMessage)
}

func main() {
	
	http.HandleFunc("/", handler)

	fmt.Println("Server is running at 3000")
	http.ListenAndServe(":3000", nil)
}