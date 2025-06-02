package main

import (
	"encoding/json"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

type DataResponse struct {
	Timestamp   time.Time `json:"timestamp"`
	RandomValue int       `json:"random_value"`
}

type MassiveDataItem struct {
	Timestamp     time.Time `json:"timestamp"`
	RandomNumber1 int       `json:"random_number_1"`
	RandomNumber2 int       `json:"random_number_2"`
	RandomNumber3 int       `json:"random_number_3"`
	RandomNumber4 int       `json:"random_number_4"`
	RandomText    string    `json:"random_text"`
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	user, password, ok := r.BasicAuth()
	log.Printf("Uesr: %s, password: %s, ok: %t", user, password, ok)
	log.Printf("Message headers: %s", r.Header)
	log.Printf("Message URL: %s", r.URL.String())

	currentTime := time.Now()

	source := rand.NewSource(currentTime.UnixNano())
	randomGenerator := rand.New(source)
	randomValue := randomGenerator.Intn(1000)
	response := DataResponse{
		Timestamp:   currentTime,
		RandomValue: randomValue,
	}

	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Error encoding JSON response: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func generateRandomString(length int, r *rand.Rand) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789 "
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

func dataMassiveHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	log.Printf("Message headers: %s", r.Header)
	log.Printf("Message URL: %s", r.URL.String())

	defaultSize := 100
	requestedSize := defaultSize

	sizeStr := r.URL.Query().Get("size")
	if sizeStr != "" {
		size, err := strconv.Atoi(sizeStr)
		if err != nil {
			http.Error(w, "Invalid 'size' parameter. Must be an integer.", http.StatusBadRequest)
			log.Printf("Invalid 'size' parameter: %v", err)
			return
		}
		if size < 0 {
			http.Error(w, "Invalid 'size' parameter. Must be a non-negative integer.", http.StatusBadRequest)
			log.Printf("Invalid 'size' parameter: %d (negative)", size)
			return
		}
		requestedSize = size
	}

	source := rand.NewSource(time.Now().UnixNano())
	randomGenerator := rand.New(source)

	var massiveData []MassiveDataItem
	for i := 0; i < requestedSize; i++ {
		item := MassiveDataItem{
			Timestamp:     time.Now(),
			RandomNumber1: randomGenerator.Intn(10000),
			RandomNumber2: randomGenerator.Intn(10000),
			RandomNumber3: randomGenerator.Intn(10000),
			RandomNumber4: randomGenerator.Intn(10000),
			RandomText:    generateRandomString(20, randomGenerator),
		}
		massiveData = append(massiveData, item)
	}
	w.Header().Set("Content-Type", "application/json")

	data, err := json.Marshal(massiveData)
	if err != nil {
		log.Printf("Error marshaling JSON for /data-massive: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Length", strconv.Itoa(len(data)))
	_, err = w.Write(data)
	if err != nil {
		log.Printf("Error writing response for /data-massive: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}

func health(w http.ResponseWriter, _ *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func main() {

	mux := http.NewServeMux()

	mux.Handle("/data", AuthMiddleware(http.HandlerFunc(dataHandler)))
	mux.Handle("/data-massive", AuthMiddleware(http.HandlerFunc(dataMassiveHandler)))
	mux.HandleFunc("/health", health)

	port := ":8080"

	server := &http.Server{
		Addr:           port,
		Handler:        mux,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	log.Printf("Starting HTTP server on port %s with custom timeouts", port)

	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Could not start server: %s\n", err)
	}
}
