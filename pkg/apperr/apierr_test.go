package apperr

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

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
			gotApiErr := GetApiError(tt.input)

			require.NotNil(t, gotApiErr, "GetApiError() should not return nil ApiErr")
			assert.Equal(t, tt.wantCode, gotApiErr.Code)
			assert.Equal(t, tt.wantMessage, gotApiErr.Message)
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
