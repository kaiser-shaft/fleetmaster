package render

import (
	"encoding/json"
	"net/http"
)

func JSON(w http.ResponseWriter, code int, body any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)

	json.NewEncoder(w).Encode(body)
}

func Error(w http.ResponseWriter, code int, msg string, err error) {
	data := map[string]string{
		"error":   err.Error(),
		"message": msg,
	}
	JSON(w, code, data)
}
