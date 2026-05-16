package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"tiara-web-app/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

// mockAuthUsecase は AuthUsecase のテスト用モック。
type mockAuthUsecase struct {
	admin domain.Admin
	err   error
}

func (m *mockAuthUsecase) Login(_ context.Context, _, _ string) (domain.Admin, error) {
	return m.admin, m.err
}

func TestAuthHandler_Login(t *testing.T) {
	t.Run("正常系: ログイン成功で JWT トークンが返る", func(t *testing.T) {
		mock := &mockAuthUsecase{
			admin: domain.Admin{ID: uuid.New(), Username: "admin"},
		}
		h := NewAuthHandler(mock, "test-secret", 2)

		body := `{"username":"admin","password":"pass123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := h.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp LoginResponse
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.NotEmpty(t, resp.Token)
	})

	t.Run("異常系: リクエストボディが不正", func(t *testing.T) {
		h := NewAuthHandler(&mockAuthUsecase{}, "test-secret", 2)

		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader("invalid"))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		_ = h.Login(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: username が空", func(t *testing.T) {
		h := NewAuthHandler(&mockAuthUsecase{}, "test-secret", 2)

		body := `{"username":"","password":"pass123"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		_ = h.Login(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "username and password are required", resp["error"])
	})

	t.Run("異常系: password が空", func(t *testing.T) {
		h := NewAuthHandler(&mockAuthUsecase{}, "test-secret", 2)

		body := `{"username":"admin","password":""}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		_ = h.Login(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 認証失敗 → 401", func(t *testing.T) {
		mock := &mockAuthUsecase{err: domain.ErrUnauthorized}
		h := NewAuthHandler(mock, "test-secret", 2)

		body := `{"username":"admin","password":"wrong"}`
		req := httptest.NewRequest(http.MethodPost, "/auth/login", strings.NewReader(body))
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		_ = h.Login(c)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestAuthHandler_Verify(t *testing.T) {
	t.Run("正常系: OK を返す", func(t *testing.T) {
		h := NewAuthHandler(&mockAuthUsecase{}, "test-secret", 2)

		req := httptest.NewRequest(http.MethodGet, "/auth/verify", nil)
		rec := httptest.NewRecorder()
		c := echo.New().NewContext(req, rec)

		err := h.Verify(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "ok", resp["status"])
	})
}
