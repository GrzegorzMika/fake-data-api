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

// MassiveDataItem defines the structure of a single item in the /data-massive response.
type MassiveDataItem struct {
	Timestamp     time.Time `json:"timestamp"`       // Current timestamp
	RandomNumber1 int       `json:"random_number_1"` // First randomly generated integer
	RandomNumber2 int       `json:"random_number_2"` // Second randomly generated integer
	RandomNumber3 int       `json:"random_number_3"` // Third randomly generated integer
	RandomNumber4 int       `json:"random_number_4"` // Fourth randomly generated integer
	RandomText    string    `json:"random_text"`     // Randomly generated text string
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

// generateRandomString generates a random string of a given length using a character set.
func generateRandomString(length int, r *rand.Rand) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 "
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

// dataMassiveHandler is the HTTP handler for the /data-massive endpoint.
func dataMassiveHandler(w http.ResponseWriter, r *http.Request) {
	// Ensure the request method is GET.
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Seed the random number generator for this handler.
	// We'll use a new source for each request to ensure different results per request.
	source := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(source)

	var massiveData []MassiveDataItem
	for i := 0; i < 100; i++ {
		item := MassiveDataItem{
			Timestamp:     time.Now(),
			RandomNumber1: randomGenerator.Intn(10000), // Random number between 0 and 9999
			RandomNumber2: randomGenerator.Intn(10000),
			RandomNumber3: randomGenerator.Intn(10000),
			RandomNumber4: randomGenerator.Intn(10000),
			RandomText:    generateRandomString(20, randomGenerator), // Random text of length 20
		}
		massiveData = append(massiveData, item)
	}

	// Set the Content-Type header to application/json.
	w.Header().Set("Content-Type", "application/json")

	// Encode the response slice to JSON.
	if err := json.NewEncoder(w).Encode(massiveData); err != nil {
		log.Printf("Error encoding JSON response for /data-massive: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	log.Printf("Served /data-massive request: Generated %d items", len(massiveData))
}

func main() {
	// Create a new ServeMux to handle HTTP requests.
	// This is explicitly used when creating a custom http.Server.
	mux := http.NewServeMux()

	// Register the dataHandler function to handle requests for the /data path on the mux.
	mux.HandleFunc("/data", dataHandler)
	// Register the dataMassiveHandler function for /data-massive.
	mux.HandleFunc("/data-massive", dataMassiveHandler)

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
