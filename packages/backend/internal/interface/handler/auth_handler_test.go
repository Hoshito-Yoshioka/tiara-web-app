package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"tiara-web-app/backend/internal/domain"
	"tiara-web-app/backend/internal/testutil"

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
		admin := testutil.NewAdmin()
		mock := &mockAuthUsecase{admin: admin}
		h := NewAuthHandler(mock, "test-secret", 2)

		c, rec := testutil.NewEchoContext(http.MethodPost, "/auth/login", `{"username":"admin","password":"pass123"}`)

		err := h.Login(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp LoginResponse
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		assert.NotEmpty(t, resp.Token)
	})

	t.Run("異常系: リクエストボディが不正", func(t *testing.T) {
		h := NewAuthHandler(&mockAuthUsecase{}, "test-secret", 2)

		c, rec := testutil.NewEchoContext(http.MethodPost, "/auth/login", "invalid")

		_ = h.Login(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: username が空", func(t *testing.T) {
		h := NewAuthHandler(&mockAuthUsecase{}, "test-secret", 2)

		c, rec := testutil.NewEchoContext(http.MethodPost, "/auth/login", `{"username":"","password":"pass123"}`)

		_ = h.Login(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)

		var resp map[string]string
		err := json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "username and password are required", resp["error"])
	})

	t.Run("異常系: password が空", func(t *testing.T) {
		h := NewAuthHandler(&mockAuthUsecase{}, "test-secret", 2)

		c, rec := testutil.NewEchoContext(http.MethodPost, "/auth/login", `{"username":"admin","password":""}`)

		_ = h.Login(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 認証失敗 → 401", func(t *testing.T) {
		mock := &mockAuthUsecase{err: domain.ErrUnauthorized}
		h := NewAuthHandler(mock, "test-secret", 2)

		c, rec := testutil.NewEchoContext(http.MethodPost, "/auth/login", `{"username":"admin","password":"wrong"}`)

		_ = h.Login(c)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)
	})
}

func TestAuthHandler_Verify(t *testing.T) {
	t.Run("正常系: OK を返す", func(t *testing.T) {
		h := NewAuthHandler(&mockAuthUsecase{}, "test-secret", 2)

		c, rec := testutil.NewEchoContext(http.MethodGet, "/auth/verify")

		err := h.Verify(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp map[string]string
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "ok", resp["status"])
	})
}
