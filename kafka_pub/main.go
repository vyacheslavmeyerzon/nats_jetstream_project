package main

import (
	"encoding/json" // For JSON encoding and decoding
	"log"           // For logging
	"net/http"      // For making HTTP requests
	"time"          // For time-related operations

	"github.com/nats-io/nats.go" // NATS client for Go
)

// KafkaState represents the state of Kafka
type KafkaState struct {
	ID        int    `json:"id"`         // Unique identifier for the state
	Name      string `json:"name"`       // Name of the state
	IsHealthy bool   `json:"is_healthy"` // Whether the state is healthy or not
	Message   string `json:"message"`    // Additional message about the state
}

// getKafkaState retrieves a Kafka state from the mock service
func getKafkaState() (KafkaState, error) {
	// Send a GET request to the mock service
	resp, err := http.Get("http://localhost:8080/state")
	if err != nil {
		return KafkaState{}, err
	}
	defer resp.Body.Close()

	// Declare a variable to store the retrieved state
	var state KafkaState

	// Decode the JSON response into the state variable
	if err := json.NewDecoder(resp.Body).Decode(&state); err != nil {
		return KafkaState{}, err
	}

	return state, nil
}

func main() {
	// NATS server URL
	natsURL := "nats://localhost:4222"

	// Connect to the NATS server
	nc, err := nats.Connect(natsURL)
	if err != nil {
		log.Fatalf("Failed to connect to NATS: %v", err)
	}
	defer nc.Close()

	// Get a JetStream context
	js, err := nc.JetStream()
	if err != nil {
		log.Fatalf("Error connecting to JetStream: %v", err)
	}

	// Create a stream for Kafka states if it doesn't exist
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "KAFKA_STREAM",
		Subjects: []string{"kafka_state"},
	})
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	log.Println("Successfully connected to NATS and JetStream")

	// Create a ticker that triggers every 10 seconds
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// Main loop
	for range ticker.C {
		// Get the current Kafka state
		state, err := getKafkaState()
		if err != nil {
			log.Printf("Error getting Kafka state: %v", err)
			continue
		}

		// Marshal the state into JSON
		stateJSON, err := json.Marshal(state)
		if err != nil {
			log.Printf("Error encoding Kafka state: %v", err)
			continue
		}

		// Publish the state to the KAFKA_STREAM
		_, err = js.Publish("kafka_state", stateJSON)
		if err != nil {
			log.Printf("Error publishing Kafka state to JetStream: %v", err)
		} else {
			log.Printf("Kafka state published: %v", state)
		}
	}
}
