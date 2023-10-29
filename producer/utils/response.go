package utils

import (
	"encoding/json"
	"net/http"
	"producer/dto"
)

// Response returns an http response with the payload provided
func Response(w http.ResponseWriter, res *dto.Response) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(res.Code)
	err := json.NewEncoder(w).Encode(res)
	if err != nil {
		return
	}
}
