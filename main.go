package main

import (
	"encoding/json" // Package for JSON encoding/decoding
	"log"           // Package for logging messages
	"math/rand"     // Package for generating pseudo-random numbers
	"net/http"      // Package for HTTP client and server implementations
	"time"          // Package for time-related functions
)

// DataResponse defines the structure of the JSON response.
type DataResponse struct {
	Timestamp   time.Time `json:"timestamp"`    // Current timestamp
	RandomValue int       `json:"random_value"` // Randomly generated integer
}

// dataHandler is the HTTP handler for the /data endpoint.
func dataHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is GET.
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Get the current timestamp.
	currentTime := time.Now()

	// Generate a random integer.
	// Using time.Now().UnixNano() to seed the random number generator
	// ensures different random values on each run.
	// For production, consider a more robust random source if cryptographic
	// randomness is required (e.g., crypto/rand).
	source := rand.NewSource(currentTime.UnixNano())
	randomGenerator := rand.New(source)
	randomValue := randomGenerator.Intn(1000) // Generates a random integer between 0 and 999

	// Create an instance of the DataResponse struct.
	response := DataResponse{
		Timestamp:   currentTime,
		RandomValue: randomValue,
	}

	// Set the Content-Type header to application/json.
	w.Header().Set("Content-Type", "application/json")

	// Encode the response struct to JSON and write it to the HTTP response writer.
	// json.NewEncoder writes directly to the writer, which is efficient.
	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Served /data request: Timestamp=%s, RandomValue=%d", currentTime.Format(time.RFC3339), randomValue)
}

func main() {
	// Create a new ServeMux to handle HTTP requests.
	// This is explicitly used when creating a custom http.Server.
	mux := http.NewServeMux()

	// Register the dataHandler function to handle requests for the /data path on the mux.
	mux.HandleFunc("/data", dataHandler)

	// Define the port to listen on.
	port := ":8080"

	// Create a custom HTTP server with specified properties.
	server := &http.Server{
		Addr:           port,             // Address to listen on (e.g., ":8080")
		Handler:        mux,              // The handler to use (our ServeMux)
		ReadTimeout:    10 * time.Second, // Maximum duration for reading the entire request, including the body.
		WriteTimeout:   10 * time.Second, // Maximum duration before timing out writes of the response.
		MaxHeaderBytes: 1 << 20,          // Maximum number of bytes the server will read parsing the request header. (1 MB)
	}

	log.Printf("Starting HTTP server on port %s with custom timeouts", port)

	// Start the HTTP server. ListenAndServe blocks until the server stops or an error occurs.
	// It returns an error if the server fails to start (e.g., port already in use).
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
