package main

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strconv"
	"sync"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

const Port = 8080

// ResponseObject represents the schema of the response objects
type ResponseObject struct {
	UUID        string   `json:"uuid"`
	Author      string   `json:"author"`
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Categories  []string `json:"categories"`
	Number      string   `json:"number"`
	Timestamp   string   `json:"timestamp"`
}

func generateRandomObject() ResponseObject {
	author := gofakeit.Name()
	number := rand.Intn(101)

	title := gofakeit.Sentence(3)
	description := ""
	if rand.Intn(2) == 1 {
		description = gofakeit.Paragraph(1, 3, 5, " ")
	}

	categories := []string{}
	if rand.Intn(2) == 1 {
		for range rand.Intn(5) + 1 {
			categories = append(categories, gofakeit.BeerName())
		}
	}

	return ResponseObject{
		UUID:        uuid.NewString(),
		Author:      author,
		Title:       title,
		Description: description,
		Categories:  categories,
		Number:      strconv.Itoa(number),
		Timestamp:   gofakeit.Date().String(),
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	message := struct {
		Api    string `json:"api"`
		Route  string `json:"route"`
		Status string `json:"status"`
	}{
		Api:    "Test-Data API",
		Route:  "/test/{amount} for test data",
		Status: "Healthy",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(message)
}

func testHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	amountStr := vars["amount"]
	log.Printf("\tPath parameter 'amount': %s", amountStr)
	amount, err := strconv.Atoi(amountStr)

	if err != nil || amount < 1 {
		http.Error(w, "Invalid amount", http.StatusBadRequest)
		log.Printf("\tResponded with status: %d", http.StatusBadRequest)
		return
	}

	responseObjects := make([]ResponseObject, amount)
	var wg sync.WaitGroup
	for i := 0; i < amount; i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			responseObjects[i] = generateRandomObject()
		}(i)
	}

	wg.Wait()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(responseObjects)
	log.Printf("\tResponded with %d objects\n", amount)
	log.Printf("\tResponded with status: %d", http.StatusOK)
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Received request:")
		log.Printf("\tMethod: %s", r.Method)
		log.Printf("\tURL: %s", r.URL.String())
		log.Printf("\tHeaders: %v", r.Header)
		log.Printf("\tClient IP: %s", r.RemoteAddr)
		next.ServeHTTP(w, r)
	})
}

func main() {
	r := mux.NewRouter()

	r.Use(loggingMiddleware)
	r.HandleFunc("/", rootHandler).Methods("GET")
	r.HandleFunc("/test/{amount}", testHandler).Methods("GET")
	log.Println("Starting Server...")

	// Start the server
	err := http.ListenAndServe(fmt.Sprintf(":%d", Port), r)
	if err != nil {
		fmt.Println(err.Error())

	}
	log.Println("Closing Server...")
}
