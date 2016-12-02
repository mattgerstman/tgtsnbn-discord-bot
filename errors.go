package main

import "errors"

// Error Types
const (
	ErrorDatabase     = "ErrorDatabase"
	ErrorInvalidHouse = "ErrorInvalidHouse"
	ErrorRoleNotFound = "ErrorRoleNotFound"
	ErrorFetchRoles   = "ErrorFetchRoles"
)

// ApplicationError contains information about errors
// that arise from application level logic.
type ApplicationError struct {
	Message string
	Error   error
	Type    string
}

// Converts an error to an Application Error with a user facing message.
func NewApplicationError(
	Message string, Error error, Type string,
) *ApplicationError {
	return &ApplicationError{
		Message: Message,
		Error:   Error,
		Type:    Type,
	}
}

// Shorthand for making a new application error without an actual error.
func NewApplicationErrorWithoutError(
	Message string, Type string,
) *ApplicationError {
	return NewApplicationError(Message, errors.New(Message), Type)
}
