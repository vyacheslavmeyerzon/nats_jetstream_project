package main

import (
	"encoding/json" // For JSON encoding and decoding
	"log"           // For logging
	"math/rand"     // For generating random numbers
	"net/http"      // For HTTP server functionality
	"time"          // For time-related operations
)

// KafkaState represents the state of Kafka
type KafkaState struct {
	ID        int    `json:"id"`         // Unique identifier for the state
	Name      string `json:"name"`       // Name of the state
	IsHealthy bool   `json:"is_healthy"` // Whether the state is healthy or not
	Message   string `json:"message"`    // Additional message about the state
}

// Predefined list of possible Kafka states
var kafkaStates = []KafkaState{
	{ID: 1, Name: "Healthy State 1", IsHealthy: true, Message: "ok no changes needed"},
	{ID: 2, Name: "Healthy State 3", IsHealthy: true, Message: "ok no changes needed"},
	{ID: 3, Name: "Healthy State 5", IsHealthy: true, Message: "ok no changes needed"},
	{ID: 4, Name: "Healthy State 7", IsHealthy: true, Message: "ok no changes needed"},
	{ID: 5, Name: "Healthy State 9", IsHealthy: true, Message: "ok no changes needed"},
	{ID: 6, Name: "Healthy State 11", IsHealthy: true, Message: "ok no changes needed"},
	{ID: 7, Name: "Healthy State 13", IsHealthy: true, Message: "ok no changes needed"},
	{ID: 8, Name: "Healthy State 15", IsHealthy: true, Message: "ok no changes needed"},
	{ID: 9, Name: "Healthy State 17", IsHealthy: true, Message: "ok no changes needed"},
	{ID: 10, Name: "Healthy State 19", IsHealthy: true, Message: "ok no changes needed"},
	{ID: 11, Name: "Healthy State 21", IsHealthy: true, Message: "ok no changes needed"},
	{ID: 12, Name: "Healthy State 23", IsHealthy: true, Message: "ok no changes needed"},
	{ID: 13, Name: "Healthy State 25", IsHealthy: true, Message: "ok no changes needed"},
	{ID: 14, Name: "Healthy State 27", IsHealthy: true, Message: "ok no changes needed"},
	{ID: 15, Name: "Healthy State 29", IsHealthy: true, Message: "ok no changes needed"},
	{ID: 16, Name: "Healthy State 31", IsHealthy: true, Message: "ok no changes needed"},
	{ID: 17, Name: "Healthy State 33", IsHealthy: true, Message: "ok no changes needed"},
	{ID: 18, Name: "Healthy State 35", IsHealthy: true, Message: "ok no changes needed"},
	{ID: 19, Name: "Healthy State 37", IsHealthy: true, Message: "ok no changes needed"},
	{ID: 20, Name: "Healthy State 39", IsHealthy: true, Message: "ok no changes needed"},
	{ID: 21, Name: "Unhealthy State 2", IsHealthy: false, Message: "configuration issue need to change following configuration AAA"},
	{ID: 22, Name: "Unhealthy State 4", IsHealthy: false, Message: "configuration issue need to change following configuration BBB"},
	{ID: 23, Name: "Unhealthy State 6", IsHealthy: false, Message: "configuration issue need to change following configuration CCC"},
	{ID: 24, Name: "Unhealthy State 8", IsHealthy: false, Message: "configuration issue need to change following configuration DDD"},
	{ID: 25, Name: "Unhealthy State 10", IsHealthy: false, Message: "configuration issue need to change following configuration EEE"},
	{ID: 26, Name: "Unhealthy State 12", IsHealthy: false, Message: "configuration issue need to change following configuration FFF"},
	{ID: 27, Name: "Unhealthy State 14", IsHealthy: false, Message: "configuration issue need to change following configuration GGG"},
	{ID: 28, Name: "Unhealthy State 16", IsHealthy: false, Message: "configuration issue need to change following configuration HHH"},
	{ID: 29, Name: "Unhealthy State 18", IsHealthy: false, Message: "configuration issue need to change following configuration LLL"},
	{ID: 30, Name: "Unhealthy State 20", IsHealthy: false, Message: "configuration issue need to change following configuration MMM"},
	{ID: 31, Name: "Unhealthy State 22", IsHealthy: false, Message: "configuration issue need to change following configuration NNN"},
	{ID: 32, Name: "Unhealthy State 24", IsHealthy: false, Message: "configuration issue need to change following configuration OOO"},
	{ID: 33, Name: "Unhealthy State 26", IsHealthy: false, Message: "configuration issue need to change following configuration PPP"},
	{ID: 34, Name: "Unhealthy State 28", IsHealthy: false, Message: "configuration issue need to change following configuration YYY"},
	{ID: 35, Name: "Unhealthy State 30", IsHealthy: false, Message: "configuration issue need to change following configuration TTT"},
	{ID: 36, Name: "Unhealthy State 32", IsHealthy: false, Message: "configuration issue need to change following configuration WWW"},
	{ID: 37, Name: "Unhealthy State 34", IsHealthy: false, Message: "configuration issue need to change following configuration QQQ"},
	{ID: 38, Name: "Unhealthy State 36", IsHealthy: false, Message: "configuration issue need to change following configuration AAA2"},
	{ID: 39, Name: "Unhealthy State 38", IsHealthy: false, Message: "configuration issue need to change following configuration AAA3"},
	{ID: 40, Name: "Unhealthy State 40", IsHealthy: false, Message: "configuration issue need to change following configuration BBB2"},
	{ID: 41, Name: "Unhealthy State 42", IsHealthy: false, Message: "configuration issue need to change following configuration BBB3"},
	{ID: 42, Name: "Unhealthy State 44", IsHealthy: false, Message: "configuration issue need to change following configuration ZZZ"},
	{ID: 43, Name: "Unhealthy State 46", IsHealthy: false, Message: "configuration issue need to change following configuration VVV"},
	{ID: 44, Name: "Unhealthy State 48", IsHealthy: false, Message: "configuration issue need to change following configuration XXX"},
	{ID: 45, Name: "Unhealthy State 50", IsHealthy: false, Message: "configuration issue need to change following configuration ZZZ7"},
	{ID: 46, Name: "Unhealthy State 52", IsHealthy: false, Message: "configuration issue need to change following configuration RRR9"},
	{ID: 47, Name: "Unhealthy State 54", IsHealthy: false, Message: "configuration issue need to change following configuration UUU3"},
	{ID: 48, Name: "Unhealthy State 56", IsHealthy: false, Message: "configuration issue need to change following configuration ABC8"},
	{ID: 49, Name: "Unhealthy State 58", IsHealthy: false, Message: "configuration issue need to change following configuration SSS1"},
	{ID: 50, Name: "Unhealthy State 60", IsHealthy: false, Message: "configuration issue need to change following configuration ASW"},
}

// Create a new random number generator
var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

// getRandomState handles GET requests to return a random Kafka state
func getRandomState(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is GET
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Select a random state from the predefined list
	state := kafkaStates[rng.Intn(len(kafkaStates))]

	// Set the response content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Encode the selected state as JSON and write it to the response
	json.NewEncoder(w).Encode(state)
}

// postAnalysisResult handles POST requests to receive analysis results
func postAnalysisResult(w http.ResponseWriter, r *http.Request) {
	// Check if the request method is POST
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Declare a variable to store the received result
	var result map[string]interface{}

	// Decode the JSON body of the request into the result variable
	if err := json.NewDecoder(r.Body).Decode(&result); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Log the received analysis result
	log.Printf("Received analysis result: %v", result)

	// Send a 200 OK status code in response
	w.WriteHeader(http.StatusOK)
}

func main() {
	// Set up HTTP handlers for our endpoints
	http.HandleFunc("/state", getRandomState)
	http.HandleFunc("/result", postAnalysisResult)

	// Define the port to run the server on
	port := ":8080"

	// Log that the server is starting
	log.Printf("Server started on port %s", port)

	// Start the HTTP server and log any fatal errors
	log.Fatal(http.ListenAndServe(port, nil))
}
