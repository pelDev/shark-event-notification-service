package httphandler

import (
	"encoding/json"
	"net/http"
)

func writeJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

// writeError writes an error response
func writeError(w http.ResponseWriter, statusCode int, message string, err error) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)

	errorResponse := map[string]interface{}{
		"error": message,
		"code":  statusCode,
	}

	if err != nil && statusCode == http.StatusInternalServerError {
		// Don't expose internal errors in production
		errorResponse["details"] = err.Error()
	}

	return json.NewEncoder(w).Encode(errorResponse)
}
