package handlers

import (
	"encoding/json"
	"net/http"
	"time"
)

func GetTime(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"time": time.Now().Format(time.RFC3339),
	})
}
