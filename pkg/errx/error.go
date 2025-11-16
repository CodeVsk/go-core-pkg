package errx

import "errors"

var (
	ErrNotFound     = errors.New("not found")
	ErrBadRequest   = errors.New("bad request")
	ErrConflict     = errors.New("conflict")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrInternal     = errors.New("internal error")
)

type Error struct {
	Kind    error
	Message string
	Details error
}

func NewNotFound(message string) *Error {
	return &Error{
		Kind:    ErrNotFound,
		Message: message,
	}
}

func NewErrBadRequest(message string) *Error {
	return &Error{
		Kind: ErrBadRequest,
	}
}

func NewErrConflict(message string) *Error {
	return &Error{
		Kind:    ErrConflict,
		Message: message,
	}
}

func NewErrUnauthorized(message string, details error) *Error {
	return &Error{
		Kind:    ErrUnauthorized,
		Message: message,
		Details: details,
	}
}

func NewErrForbidden(message string, details error) *Error {
	return &Error{
		Kind:    ErrForbidden,
		Message: message,
		Details: details,
	}
}

func NewErrInternal(message string, details error) *Error {
	return &Error{
		Kind:    ErrInternal,
		Message: message,
		Details: details,
	}
}
