package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
)

// ValidationError represents a validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// ErrorResponse represents an API error response
type ErrorResponse struct {
	Error   string            `json:"error"`
	Message string            `json:"message,omitempty"`
	Details []ValidationError `json:"details,omitempty"`
}

// WriteErrorResponse writes a JSON error response
func WriteErrorResponse(w http.ResponseWriter, statusCode int, err string, message string, details []ValidationError) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	
	response := ErrorResponse{
		Error:   err,
		Message: message,
		Details: details,
	}
	
	json.NewEncoder(w).Encode(response)
}

// ValidateVoteDirection validates vote direction parameter
func ValidateVoteDirection(direction string) error {
	if direction == "" {
		return fmt.Errorf("vote direction is required")
	}
	if direction != "up" && direction != "down" {
		return fmt.Errorf("vote direction must be 'up' or 'down'")
	}
	return nil
}

// ValidatePokemonID validates Pokemon ID parameter
func ValidatePokemonID(idStr string) (int, error) {
	if idStr == "" {
		return 0, fmt.Errorf("Pokemon ID is required")
	}
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, fmt.Errorf("Pokemon ID must be a valid integer")
	}
	
	if id <= 0 {
		return 0, fmt.Errorf("Pokemon ID must be a positive integer")
	}
	
	return id, nil
}

// ValidateQueryParams validates common query parameters
func ValidateQueryParams(r *http.Request, required []string) []ValidationError {
	var errors []ValidationError
	params := r.URL.Query()
	
	for _, param := range required {
		if params.Get(param) == "" {
			errors = append(errors, ValidationError{
				Field:   param,
				Message: fmt.Sprintf("%s is required", param),
			})
		}
	}
	
	return errors
}