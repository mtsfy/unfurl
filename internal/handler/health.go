package handler

import (
	"encoding/json"
	"net/http"
	"time"
)

type Health struct {
	Status    string `json:"status"`
	Timestamp string `json:"timestamp"`
}

func GetHealth(w http.ResponseWriter, r *http.Request) {
	h := Health{
		Status:    "OK!",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	data, err := json.Marshal(h)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
