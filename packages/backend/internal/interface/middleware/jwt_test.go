package middleware

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
)

func makeTestJWT(t *testing.T, secret, tokenType string) string {
	t.Helper()

	claims := jwt.MapClaims{
		"sub":      "user-123",
		"username": "tester",
		"type":     tokenType,
		"exp":      time.Now().Add(1 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	assert.NoError(t, err)
	return tokenString
}

func TestJWTAuth_AuthorizationByRole(t *testing.T) {
	const secret = "test-secret"
	e := echo.New()

	adminHandler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"result": "admin operation allowed"})
	}

	t.Run("admin操作はadminトークンで実行可能", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin/reviews/schedules/1/publish", nil)
		req.Header.Set("Authorization", "Bearer "+makeTestJWT(t, secret, "admin"))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := JWTAuth(secret)(adminHandler)(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "user-123", c.Get("admin_id"))
	})

	t.Run("一般ユーザー(staffトークン)はadmin操作を実行不可", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/admin/reviews/schedules/1/publish", nil)
		req.Header.Set("Authorization", "Bearer "+makeTestJWT(t, secret, "staff"))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := JWTAuth(secret)(adminHandler)(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var body map[string]string
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
		assert.Equal(t, "invalid token type", body["error"])
	})
}

func TestStaffJWTAuth_AuthorizationByRole(t *testing.T) {
	const secret = "test-secret"
	e := echo.New()

	staffHandler := func(c echo.Context) error {
		return c.JSON(http.StatusOK, map[string]string{"result": "staff operation allowed"})
	}

	t.Run("staff操作はstaffトークンで実行可能", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/portal/profile", nil)
		req.Header.Set("Authorization", "Bearer "+makeTestJWT(t, secret, "staff"))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := StaffJWTAuth(secret)(staffHandler)(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
		assert.Equal(t, "user-123", c.Get("staff_id"))
	})

	t.Run("adminトークンはstaff操作を実行不可", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/portal/profile", nil)
		req.Header.Set("Authorization", "Bearer "+makeTestJWT(t, secret, "admin"))
		rec := httptest.NewRecorder()
		c := e.NewContext(req, rec)

		err := StaffJWTAuth(secret)(staffHandler)(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusUnauthorized, rec.Code)

		var body map[string]string
		assert.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
		assert.Equal(t, "invalid token type", body["error"])
	})
}
