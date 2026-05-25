package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"tiara-web-app/backend/internal/domain"
	"tiara-web-app/backend/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestHandleError(t *testing.T) {
	tests := []struct {
		name           string
		err            error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "ErrNotFound → 404",
			err:            domain.ErrNotFound,
			expectedStatus: http.StatusNotFound,
			expectedBody:   domain.ErrNotFound.Error(),
		},
		{
			name:           "ErrUnauthorized → 401",
			err:            domain.ErrUnauthorized,
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   domain.ErrUnauthorized.Error(),
		},
		{
			name:           "ErrForbidden → 403",
			err:            domain.ErrForbidden,
			expectedStatus: http.StatusForbidden,
			expectedBody:   domain.ErrForbidden.Error(),
		},
		{
			name:           "ErrConflict → 409",
			err:            domain.ErrConflict,
			expectedStatus: http.StatusConflict,
			expectedBody:   domain.ErrConflict.Error(),
		},
		{
			name:           "ErrInvalidInput → 400",
			err:            domain.ErrInvalidInput,
			expectedStatus: http.StatusBadRequest,
			expectedBody:   domain.ErrInvalidInput.Error(),
		},
		{
			name:           "ラップされたドメインエラーも正しくマッピングされる",
			err:            domain.WrapNotFound("staff not found"),
			expectedStatus: http.StatusNotFound,
			expectedBody:   "staff not found: resource not found",
		},
		{
			name:           "未知のエラー → 500",
			err:            errors.New("unexpected db error"),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "internal server error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := testutil.NewEchoContext(http.MethodGet, "/")

			_ = handleError(c, tt.err)

			assert.Equal(t, tt.expectedStatus, rec.Code)

			var body map[string]string
			err := json.Unmarshal(rec.Body.Bytes(), &body)
			assert.NoError(t, err)
			assert.Equal(t, tt.expectedBody, body["error"])
		})
	}
}
