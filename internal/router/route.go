package router

import (
	"net/http"

	"github.com/mtsfy/unfurl/internal/handler"
)

func SetupRoutes(router *http.ServeMux) {
	router.HandleFunc("GET /api/v1/health", handler.GetHealth)
	router.HandleFunc("POST /api/v1/unfurl", handler.PostUnfurl)
}
