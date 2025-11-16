package apperr

import (
	"errors"
)

func (e *AppErr) Error() string {
	return e.Message
}

func NewNotFound(message string) *AppErr {
	return &AppErr{
		Kind:    ErrNotFound,
		Message: message,
	}
}

func NewErrBadRequest(message string) *AppErr {
	return &AppErr{
		Kind:    ErrBadRequest,
		Message: message,
	}
}

func NewErrConflict(message string) *AppErr {
	return &AppErr{
		Kind:    ErrConflict,
		Message: message,
	}
}

func NewErrUnauthorized(message string, details error) *AppErr {
	return &AppErr{
		Kind:    ErrUnauthorized,
		Message: message,
		Details: details,
	}
}

func NewErrForbidden(message string, details error) *AppErr {
	return &AppErr{
		Kind:    ErrForbidden,
		Message: message,
		Details: details,
	}
}

func NewErrInternal(message string, details error) *AppErr {
	return &AppErr{
		Kind:    ErrInternal,
		Message: message,
		Details: details,
	}
}

func GetError(appError error) *AppErr {
	var err *AppErr

	if errors.As(appError, &err) {
		return err
	}

	return &AppErr{
		Kind:    ErrInternal,
		Message: ErrInternal.Error(),
		Details: errors.New("the received error is not compatible with AppErr"),
	}
}
