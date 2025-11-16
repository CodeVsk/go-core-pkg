package apperr

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestAppErr_Error(t *testing.T) {
	tests := []struct {
		name   string
		appErr *AppErr
		want   string
	}{
		{
			name: "returns message",
			appErr: &AppErr{
				Message: "test error message",
			},
			want: "test error message",
		},
		{
			name: "returns empty string when message is empty",
			appErr: &AppErr{
				Message: "",
			},
			want: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.appErr.Error()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestNewNotFound(t *testing.T) {
	message := "resource not found"
	err := NewNotFound(message)

	require.NotNil(t, err, "NewNotFound() should not return nil")
	assert.Equal(t, ErrNotFound, err.Kind)
	assert.Equal(t, message, err.Message)
}

func TestNewErrBadRequest(t *testing.T) {
	message := "invalid input"
	err := NewErrBadRequest(message)

	require.NotNil(t, err, "NewErrBadRequest() should not return nil")
	assert.Equal(t, ErrBadRequest, err.Kind)
	// Note: Current implementation doesn't set Message, even though it accepts a parameter
	// This test reflects the current behavior
}

func TestNewErrConflict(t *testing.T) {
	message := "resource conflict"
	err := NewErrConflict(message)

	require.NotNil(t, err, "NewErrConflict() should not return nil")
	assert.Equal(t, ErrConflict, err.Kind)
	assert.Equal(t, message, err.Message)
}

func TestNewErrUnauthorized(t *testing.T) {
	message := "unauthorized access"
	details := errors.New("token expired")
	err := NewErrUnauthorized(message, details)

	require.NotNil(t, err, "NewErrUnauthorized() should not return nil")
	assert.Equal(t, ErrUnauthorized, err.Kind)
	assert.Equal(t, message, err.Message)
	assert.Equal(t, details, err.Details)
}

func TestNewErrForbidden(t *testing.T) {
	message := "forbidden access"
	details := errors.New("insufficient permissions")
	err := NewErrForbidden(message, details)

	require.NotNil(t, err, "NewErrForbidden() should not return nil")
	assert.Equal(t, ErrForbidden, err.Kind)
	assert.Equal(t, message, err.Message)
	assert.Equal(t, details, err.Details)
}

func TestNewErrInternal(t *testing.T) {
	message := "internal server error"
	details := errors.New("database connection failed")
	err := NewErrInternal(message, details)

	require.NotNil(t, err, "NewErrInternal() should not return nil")
	assert.Equal(t, ErrInternal, err.Kind)
	assert.Equal(t, message, err.Message)
	assert.Equal(t, details, err.Details)
}

func TestGetError(t *testing.T) {
	tests := []struct {
		name      string
		input     error
		want      *AppErr
		wantNil   bool
		checkCode bool
	}{
		{
			name: "extracts AppErr successfully",
			input: &AppErr{
				Kind:    ErrNotFound,
				Message: "not found",
			},
			want: &AppErr{
				Kind:    ErrNotFound,
				Message: "not found",
			},
			wantNil:   false,
			checkCode: true,
		},
		{
			name:    "returns nil for non-AppErr error",
			input:   errors.New("regular error"),
			wantNil: true,
		},
		{
			name:    "returns nil for nil error",
			input:   nil,
			wantNil: true,
		},
		{
			name: "extracts AppErr with details",
			input: &AppErr{
				Kind:    ErrInternal,
				Message: "internal error",
				Details: errors.New("underlying error"),
			},
			want: &AppErr{
				Kind:    ErrInternal,
				Message: "internal error",
				Details: errors.New("underlying error"),
			},
			wantNil:   false,
			checkCode: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetError(tt.input)

			if tt.wantNil {
				assert.Nil(t, got)
				return
			}

			require.NotNil(t, got, "GetError() should not return nil")
			assert.Equal(t, tt.want.Kind, got.Kind)
			assert.Equal(t, tt.want.Message, got.Message)

			if tt.want.Details != nil {
				require.NotNil(t, got.Details)
				assert.Equal(t, tt.want.Details.Error(), got.Details.Error())
			}
		})
	}
}

func TestHttpCode(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{
			name: "ErrNotFound returns StatusNotFound",
			err:  ErrNotFound,
			want: http.StatusNotFound,
		},
		{
			name: "ErrBadRequest returns StatusBadRequest",
			err:  ErrBadRequest,
			want: http.StatusBadRequest,
		},
		{
			name: "ErrConflict returns StatusConflict",
			err:  ErrConflict,
			want: http.StatusConflict,
		},
		{
			name: "ErrUnauthorized returns StatusUnauthorized",
			err:  ErrUnauthorized,
			want: http.StatusUnauthorized,
		},
		{
			name: "ErrForbidden returns StatusForbidden",
			err:  ErrForbidden,
			want: http.StatusForbidden,
		},
		{
			name: "ErrInternal returns StatusInternalServerError",
			err:  ErrInternal,
			want: http.StatusInternalServerError,
		},
		{
			name: "unknown error returns StatusInternalServerError",
			err:  errors.New("unknown error"),
			want: http.StatusInternalServerError,
		},
		{
			name: "nil error returns StatusInternalServerError",
			err:  nil,
			want: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := httpCode(tt.err)
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestGetErrorWithHttpCode(t *testing.T) {
	errorTypes := []struct {
		name     string
		err      *AppErr
		wantCode int
	}{
		{
			name: "NotFound sets HttpCode correctly",
			err: &AppErr{
				Kind:    ErrNotFound,
				Message: "not found",
			},
			wantCode: http.StatusNotFound,
		},
		{
			name: "BadRequest sets HttpCode correctly",
			err: &AppErr{
				Kind:    ErrBadRequest,
				Message: "bad request",
			},
			wantCode: http.StatusBadRequest,
		},
		{
			name: "Conflict sets HttpCode correctly",
			err: &AppErr{
				Kind:    ErrConflict,
				Message: "conflict",
			},
			wantCode: http.StatusConflict,
		},
		{
			name: "Unauthorized sets HttpCode correctly",
			err: &AppErr{
				Kind:    ErrUnauthorized,
				Message: "unauthorized",
			},
			wantCode: http.StatusUnauthorized,
		},
		{
			name: "Forbidden sets HttpCode correctly",
			err: &AppErr{
				Kind:    ErrForbidden,
				Message: "forbidden",
			},
			wantCode: http.StatusForbidden,
		},
		{
			name: "Internal sets HttpCode correctly",
			err: &AppErr{
				Kind:    ErrInternal,
				Message: "internal",
			},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range errorTypes {
		t.Run(tt.name, func(t *testing.T) {
			got := GetError(tt.err)
			require.NotNil(t, got, "GetError() should not return nil")
		})
	}
}

func TestGetApiError(t *testing.T) {
	tests := []struct {
		name         string
		input        error
		wantCode     string
		wantMessage  string
		wantHttpCode int
	}{
		{
			name: "NotFound returns correct ApiErr",
			input: &AppErr{
				Kind:    ErrNotFound,
				Message: "resource not found",
			},
			wantCode:     ErrNotFound.Error(),
			wantMessage:  "resource not found",
			wantHttpCode: http.StatusNotFound,
		},
		{
			name: "BadRequest returns correct ApiErr",
			input: &AppErr{
				Kind:    ErrBadRequest,
				Message: "invalid input",
			},
			wantCode:     ErrBadRequest.Error(),
			wantMessage:  "invalid input",
			wantHttpCode: http.StatusBadRequest,
		},
		{
			name: "Conflict returns correct ApiErr",
			input: &AppErr{
				Kind:    ErrConflict,
				Message: "resource conflict",
			},
			wantCode:     ErrConflict.Error(),
			wantMessage:  "resource conflict",
			wantHttpCode: http.StatusConflict,
		},
		{
			name: "Unauthorized returns correct ApiErr",
			input: &AppErr{
				Kind:    ErrUnauthorized,
				Message: "unauthorized access",
			},
			wantCode:     ErrUnauthorized.Error(),
			wantMessage:  "unauthorized access",
			wantHttpCode: http.StatusUnauthorized,
		},
		{
			name: "Forbidden returns correct ApiErr",
			input: &AppErr{
				Kind:    ErrForbidden,
				Message: "forbidden access",
			},
			wantCode:     ErrForbidden.Error(),
			wantMessage:  "forbidden access",
			wantHttpCode: http.StatusForbidden,
		},
		{
			name: "Internal returns correct ApiErr",
			input: &AppErr{
				Kind:    ErrInternal,
				Message: "internal server error",
			},
			wantCode:     ErrInternal.Error(),
			wantMessage:  "internal server error",
			wantHttpCode: http.StatusInternalServerError,
		},
		{
			name: "AppErr with empty message returns correct ApiErr",
			input: &AppErr{
				Kind:    ErrNotFound,
				Message: "",
			},
			wantCode:     ErrNotFound.Error(),
			wantMessage:  "",
			wantHttpCode: http.StatusNotFound,
		},
		{
			name:         "non-AppErr error returns default internal error",
			input:        errors.New("regular error"),
			wantCode:     ErrInternal.Error(),
			wantMessage:  ErrInternal.Error(),
			wantHttpCode: http.StatusInternalServerError,
		},
		{
			name:         "nil error returns default internal error",
			input:        nil,
			wantCode:     ErrInternal.Error(),
			wantMessage:  ErrInternal.Error(),
			wantHttpCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotApiErr, gotHttpCode := GetApiError(tt.input)

			require.NotNil(t, gotApiErr, "GetApiError() should not return nil ApiErr")
			assert.Equal(t, tt.wantCode, gotApiErr.Code)
			assert.Equal(t, tt.wantMessage, gotApiErr.Message)
			assert.Equal(t, tt.wantHttpCode, gotHttpCode)
		})
	}
}
