package main

import (
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/pierre-teodoresco/calculator-api-go/internal/handler"
)

func main() {
	// Env
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Routes
	port := os.Getenv("PORT")

	http.HandleFunc("GET /health", handler.HealthHandler)
	http.HandleFunc("POST /add", handler.AddHandler)
	http.HandleFunc("POST /multiply", handler.MultiplyHandler)
	http.HandleFunc("POST /subtract", handler.SubtractHandler)
	http.HandleFunc("POST /divide", handler.DivideHandler)

	log.Println("API is listening on port", port)
	log.Fatal(http.ListenAndServe(port, nil))
}
