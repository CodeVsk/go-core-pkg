package apperr

import (
	"errors"
	"net/http"
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
		Kind: ErrBadRequest,
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

	return nil
}

func GetApiError(appError error) *ApiErr {
	err := GetError(appError)
	if err != nil {
		return &ApiErr{
			Code:       err.Kind.Error(),
			Message:    err.Message,
			statusCode: httpCode(err.Kind),
		}
	}

	return &ApiErr{
		Code:       ErrInternal.Error(),
		Message:    ErrInternal.Error(),
		statusCode: httpCode(ErrInternal),
	}
}

func (err *ApiErr) StatusCode() int {
	return err.statusCode
}

func httpCode(err error) int {
	switch err {
	case ErrNotFound:
		return http.StatusNotFound
	case ErrBadRequest:
		return http.StatusBadRequest
	case ErrConflict:
		return http.StatusConflict
	case ErrUnauthorized:
		return http.StatusUnauthorized
	case ErrForbidden:
		return http.StatusForbidden
	default:
		return http.StatusInternalServerError
	}
}
