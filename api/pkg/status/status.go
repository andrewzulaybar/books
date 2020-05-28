package status

import (
	"errors"
	"fmt"
)

// List of HTTP status codes.
const (
	OK        int = 200
	Created   int = 201
	NoContent int = 204

	BadRequest          int = 400
	Unauthorized        int = 401
	DoesNotExist        int = 404
	MethodNotAllowed    int = 405
	Conflict            int = 409
	UnprocessableEntity int = 422
	TooManyRequests     int = 429

	InternalServerError int = 500
	NotImplemented      int = 501
)

// Status contains information about a response, namely the status code and a message.
type Status struct {
	code    int
	message string
}

// New returns a pointer to a new Status struct instantiated with the given values.
func New(c int, msg string) *Status {
	return &Status{code: c, message: msg}
}

// Newf calls New(int, string) with the formatted string as the message.
func Newf(c int, format string, a ...interface{}) *Status {
	return New(c, fmt.Sprintf(format, a...))
}

// Code returns the receiver's code field.
func (s *Status) Code() int {
	return s.code
}

// Message returns the receiver's message field.
func (s *Status) Message() string {
	return s.message
}

// Err returns an error with the receiver's message as the error text.
// If the receiver's status code is of the form 2xx (success), it returns nil.
func (s *Status) Err() error {
	class := s.code / 100
	switch class {
	case 2:
		return nil
	default:
		return errors.New(s.message)
	}
}
