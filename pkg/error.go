package pkg

import (
	"encoding/json"
	"net/http"
)

type Error struct {
	Message string `json:"error"`
	Status  int    `json:"-"`
}

func SetJSONHeader(w http.ResponseWriter) {
	w.Header().Set("Content-Type", "application/json")
}

func SendJSONError(w http.ResponseWriter, err *Error) {
	SetJSONHeader(w)
	w.WriteHeader(err.Status)
	json.NewEncoder(w).Encode(err)
}

func SendJSON(w http.ResponseWriter, statusCode int, content any) {
	SetJSONHeader(w)
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(content)
}
