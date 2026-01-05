// Package response provides types and functions for HTTP API responses.
package response

import (
	"fmt"
	"strings"

	"github.com/go-playground/validator/v10"
)

// Response represents a standard API response structure.
type Response struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

const (
	// StatusOk indicates a successful response.
	StatusOk = "Ok"
	// StatusError indicates an error response.
	StatusError = "Error"
)

// Error creates an error response with the given error message.
func Error(errMsg string) Response {
	return Response{
		Status: StatusError,
		Error:  errMsg,
	}
}

// Ok creates a successful response.
func Ok() Response {
	return Response{
		Status: StatusOk,
	}
}

// ValidationError creates an error response from validation errors.
func ValidationError(errs validator.ValidationErrors) Response {
	var errMsg []string

	for _, err := range errs {
		switch err.ActualTag() {
		case "required":
			errMsg = append(errMsg, fmt.Sprintf("field %s is a required field", err.Field()))
		case "url":
			errMsg = append(errMsg, fmt.Sprintf("field %s is not a valid URL", err.Field()))
		default:
			errMsg = append(errMsg, fmt.Sprintf("field %s is not valid", err.Field()))
		}
	}

	return Error(strings.Join(errMsg, ", "))
}
