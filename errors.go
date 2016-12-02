package main

// Error Types
const (
	ErrorInvalidHouse = "ErrorInvalidHouse"
)

// ApplicationError contains information about errors
// that arise from application level logic.
type ApplicationError struct {
	Message string
	Error   error
	Type    string
}
