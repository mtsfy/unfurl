package handler

import (
	"encoding/json"
	"net/http"
)

type ExtractedData struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Site        string `json:"site"`
}

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

	resp := ExtractedData{
		Title:       "TODO: Title",
		Description: "TODO: Description",
		Image:       "TODO: Image",
		Site:        "TODO: Site",
	}

	data, err := json.Marshal(resp)
	if err != nil {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(data)
}
