package handler

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/pierre-teodoresco/calculator-api-go/pkg"
)

type Args struct {
	A int `json:"a"`
	B int `json:"b"`
}

type Result struct {
	Value int `json:"result"`
}

func ParseRequest(r *http.Request) (Args, *pkg.Error) {
	if r.Header.Get("Content-Type") != "" && !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		return Args{}, &pkg.Error{Message: "Content-Type must be application/json", Status: http.StatusUnsupportedMediaType}
	}

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var args Args
	err := decoder.Decode(&args)
	if err != nil {
		// Check if the error is due to invalid JSON structure
		if _, ok := err.(*json.SyntaxError); ok {
			return Args{}, &pkg.Error{Message: "Invalid JSON format: malformed JSON structure", Status: http.StatusBadRequest}
		}
		if typeErr, ok := err.(*json.UnmarshalTypeError); ok {
			return Args{}, &pkg.Error{Message: "Invalid field type: field '" + typeErr.Field + "' must be an integer", Status: http.StatusBadRequest}
		}
		// Check for unknown fields
		if strings.Contains(err.Error(), "unknown field") {
			return Args{}, &pkg.Error{Message: "Invalid request format: only 'a' and 'b' fields are allowed", Status: http.StatusBadRequest}
		}
		return Args{}, &pkg.Error{Message: "Invalid request format: expected JSON with 'a' and 'b' integer fields", Status: http.StatusBadRequest}
	}

	return args, nil
}

func AddHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and verify JSON from request
	args, err := ParseRequest(r)
	if err != nil {
		log.Println("[ADD] Parsing error:", err.Message)
		pkg.SendJSONError(w, err)
		return
	}

	// Compute addition
	result := Result{args.A + args.B}

	// Send response
	log.Printf("[ADD] %d + %d = %d", args.A, args.B, result.Value)
	pkg.SendJSON(w, http.StatusOK, result)
}

func MultiplyHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and verify JSON from request
	args, err := ParseRequest(r)
	if err != nil {
		log.Println("[MULTIPLY] Parsing error:", err.Message)
		pkg.SendJSONError(w, err)
		return
	}

	// Compute multiplication
	result := Result{args.A * args.B}

	// Send response
	log.Printf("[MULTIPLY] %d * %d = %d", args.A, args.B, result.Value)
	pkg.SendJSON(w, http.StatusOK, result)
}

func SubtractHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and verify JSON from request
	args, err := ParseRequest(r)
	if err != nil {
		log.Println("[SUBTRACT] Parsing error:", err.Message)
		pkg.SendJSONError(w, err)
		return
	}

	// Compute subtract
	result := Result{args.A - args.B}

	// Send response
	log.Printf("[SUBTRACT] %d - %d = %d", args.A, args.B, result.Value)
	pkg.SendJSON(w, http.StatusOK, result)
}

func DivideHandler(w http.ResponseWriter, r *http.Request) {
	// Parse and verify JSON from request
	args, err := ParseRequest(r)
	if err != nil {
		log.Println("[DIVIDE] Parsing error:", err.Message)
		pkg.SendJSONError(w, err)
		return
	}

	// Division by 0
	if args.B == 0 {
		err = &pkg.Error{Message: "Can't divide by 0", Status: http.StatusBadRequest}
		log.Println("[DIVIDE] Entry issue:", err.Message)
		pkg.SendJSONError(w, err)
		return
	}

	// Compute division
	result := Result{args.A / args.B}

	// Send response
	log.Printf("[DIVIDE] %d - %d = %d", args.A, args.B, result.Value)
	pkg.SendJSON(w, http.StatusOK, result)
}
