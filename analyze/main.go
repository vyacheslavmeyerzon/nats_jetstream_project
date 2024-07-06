package main

import (
	"encoding/json" // For JSON encoding and decoding
	"log"           // For logging

	"github.com/nats-io/nats.go" // NATS client for Go
)

// KafkaState represents the state of Kafka
type KafkaState struct {
	ID        int    `json:"id"`         // Unique identifier for the state
	Name      string `json:"name"`       // Name of the state
	IsHealthy bool   `json:"is_healthy"` // Whether the state is healthy or not
	Message   string `json:"message"`    // Additional message about the state
}

// AnalysisResult represents the result of analyzing a Kafka state
type AnalysisResult struct {
	ID      int    `json:"id"`      // Identifier of the analyzed state
	Status  string `json:"status"`  // Status of the analysis (ok or issue)
	Message string `json:"message"` // Detailed message about the analysis result
}

// analyzeState performs analysis on a given Kafka state
func analyzeState(state KafkaState) AnalysisResult {
	if state.IsHealthy {
		return AnalysisResult{
			ID:      state.ID,
			Status:  "ok",
			Message: "ok no changes needed",
		}
	}
	return AnalysisResult{
		ID:      state.ID,
		Status:  "issue",
		Message: state.Message,
	}
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

	// Create a stream for analyzed states if it doesn't exist
	_, err = js.AddStream(&nats.StreamConfig{
		Name:     "ANALYZED_STREAM",
		Subjects: []string{"analyzed_state"},
	})
	if err != nil {
		log.Fatalf("Error creating stream: %v", err)
	}

	log.Println("Successfully connected to NATS and JetStream")

	// Subscribe to the kafka_state subject
	_, err = js.Subscribe("kafka_state", func(msg *nats.Msg) {
		// Declare a variable to store the received state
		var state KafkaState

		// Unmarshal the message data into the state variable
		if err := json.Unmarshal(msg.Data, &state); err != nil {
			log.Printf("Error decoding message: %v", err)
			return
		}

		log.Printf("Received Kafka state: %v", state)

		// Analyze the received state
		result := analyzeState(state)

		// Marshal the analysis result into JSON
		resultJSON, err := json.Marshal(result)
		if err != nil {
			log.Printf("Error encoding analysis result: %v", err)
			return
		}

		// Publish the analysis result to the ANALYZED_STREAM
		if _, err := js.Publish("analyzed_state", resultJSON); err != nil {
			log.Printf("Error publishing analysis result to JetStream: %v", err)
		} else {
			log.Printf("Analysis result published: %v", result)
		}
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to kafka_state channel: %v", err)
	}

	log.Println("Successfully subscribed to kafka_state channel")

	// Keep the application running
	select {}
}
