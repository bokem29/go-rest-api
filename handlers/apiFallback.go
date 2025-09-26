package handlers

import (
	"encoding/json"
	"net/http"
)

type errorResponse struct {
	Error   string `json:"error"`
	Message string `json:"message"`
}

// ApiNotFoundHandler returns JSON 404 for unknown /api/* routes
func ApiNotFoundHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusNotFound)
	json.NewEncoder(w).Encode(errorResponse{
		Error:   "not_found",
		Message: "API route not found",
	})
}
