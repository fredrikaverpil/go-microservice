package domain

import "fmt"

type ErrorType int

const (
	NotFound ErrorType = iota
	AlreadyExists
	InvalidInput
	Internal
	Timeout
	Unavailable
	ResourceExhausted
)

type Error struct {
	Type    ErrorType
	Message string
	Err     error
}

func (e *Error) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.Message, e.Err)
	}
	return e.Message
}

func (e *Error) Unwrap() error {
	return e.Err
}

func NewErrorNotFound(message string, err error) error {
	return &Error{Type: NotFound, Message: message, Err: err}
}

func NewErrorAlreadyExists(message string, err error) error {
	return &Error{Type: AlreadyExists, Message: message, Err: err}
}

func NewErrorInvalidInput(message string, err error) error {
	return &Error{Type: InvalidInput, Message: message, Err: err}
}

func NewErrorInternal(message string, err error) error {
	return &Error{Type: Internal, Message: message, Err: err}
}

func NewErrorTimeout(message string, err error) error {
	return &Error{Type: Timeout, Message: message, Err: err}
}

func NewErrorUnavailable(message string, err error) error {
	return &Error{Type: Unavailable, Message: message, Err: err}
}

func NewErrorResourceExhausted(message string, err error) error {
	return &Error{Type: ResourceExhausted, Message: message, Err: err}
}
