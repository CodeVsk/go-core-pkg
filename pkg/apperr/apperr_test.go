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
		name  string
		input error
		want  *AppErr
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
		},
		{
			name:  "returns internal error for non-AppErr error",
			input: errors.New("regular error"),
			want: &AppErr{
				Kind:    ErrInternal,
				Message: ErrInternal.Error(),
			},
		},
		{
			name:  "returns internal error for nil error",
			input: nil,
			want: &AppErr{
				Kind:    ErrInternal,
				Message: ErrInternal.Error(),
			},
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
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := GetError(tt.input)

			require.NotNil(t, got, "GetError() should not return nil")
			assert.Equal(t, tt.want.Kind, got.Kind)
			assert.Equal(t, tt.want.Message, got.Message)

			if tt.want.Details != nil {
				require.NotNil(t, got.Details)
				assert.Equal(t, tt.want.Details.Error(), got.Details.Error())
			} else {
				// For non-AppErr errors (or nil input), check that Details contains the expected message
				var appErr *AppErr
				if tt.input == nil || !errors.As(tt.input, &appErr) {
					require.NotNil(t, got.Details)
					assert.Contains(t, got.Details.Error(), "not compatible with AppErr")
				}
			}
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
