package apperr

import "net/http"

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
