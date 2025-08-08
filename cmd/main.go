package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/mtsfy/unfurl/internal/router"
)

func main() {
	mux := http.NewServeMux()
	router.SetupRoutes(mux)

	fmt.Println("Server running on port :8080")
	if err := http.ListenAndServe(":8080", mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
