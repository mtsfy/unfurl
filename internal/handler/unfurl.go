package handler

import (
	"encoding/json"
	"net/http"

	"github.com/mtsfy/unfurl/internal/service"
)

type InputData struct {
	Url string `json:"url"`
}

func PostUnfurl(w http.ResponseWriter, r *http.Request) {
	var input InputData
	err := json.NewDecoder(r.Body).Decode(&input)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusBadRequest), http.StatusBadRequest)
		return
	}

	html, err := service.Fetch(input.Url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	extracted, err := service.Extract(html, input.Url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data, err := json.Marshal(extracted)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
