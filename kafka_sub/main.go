package main

import (
	"bytes"         // For creating a bytes buffer
	"encoding/json" // For JSON encoding and decoding
	"fmt"           // For formatting strings
	"log"           // For logging
	"net/http"      // For making HTTP requests

	"github.com/nats-io/nats.go" // NATS client for Go
)

// AnalysisResult represents the result of analyzing a Kafka state
type AnalysisResult struct {
	ID      int    `json:"id"`      // Identifier of the analyzed state
	Status  string `json:"status"`  // Status of the analysis (ok or issue)
	Message string `json:"message"` // Detailed message about the analysis result
}

// sendAnalysisResult sends the analysis result to the mock service
func sendAnalysisResult(result AnalysisResult) error {
	// Marshal the result into JSON
	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Errorf("marshaling error: %w", err)
	}

	// Send a POST request to the mock service with the result
	resp, err := http.Post("http://localhost:8080/result", "application/json", bytes.NewBuffer(resultJSON))
	if err != nil {
		return fmt.Errorf("error sending request: %w", err)
	}
	defer resp.Body.Close()

	// Check if the response status is OK
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected response status: %d", resp.StatusCode)
	}

	return nil
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

	// Subscribe to the analyzed_state subject
	_, err = js.Subscribe("analyzed_state", func(msg *nats.Msg) {
		// Declare a variable to store the received result
		var result AnalysisResult

		// Unmarshal the message data into the result variable
		if err := json.Unmarshal(msg.Data, &result); err != nil {
			log.Printf("Error decoding message: %v", err)
			return
		}

		log.Printf("Received analysis result: %v", result)

		// Send the analysis result to the mock service
		if err := sendAnalysisResult(result); err != nil {
			log.Printf("Error sending analysis result: %v", err)
		} else {
			log.Println("Analysis result successfully sent to Mock microservice")
		}
	})
	if err != nil {
		log.Fatalf("Failed to subscribe to analyzed_state channel: %v", err)
	}

	log.Println("Successfully subscribed to analyzed_state channel")

	// Keep the application running
	select {}
}
