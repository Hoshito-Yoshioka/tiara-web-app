// Package testutil はテスト用の共通ヘルパーを提供する。
package testutil

import (
	"net/http"
	"net/http/httptest"
	"strings"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// NewEchoContext は Echo コンテキストとレスポンスレコーダーを生成するヘルパー。
// body を渡すと Content-Type: application/json を自動設定する。
func NewEchoContext(method, path string, body ...string) (echo.Context, *httptest.ResponseRecorder) {
	var req *http.Request
	if len(body) > 0 && body[0] != "" {
		req = httptest.NewRequest(method, path, strings.NewReader(body[0]))
		req.Header.Set("Content-Type", "application/json")
	} else {
		req = httptest.NewRequest(method, path, nil)
	}
	rec := httptest.NewRecorder()
	c := echo.New().NewContext(req, rec)
	return c, rec
}

// SetPathParams は Echo コンテキストにパスパラメータをセットする。
// pairs は "name", "value" の順で交互に渡す。
func SetPathParams(c echo.Context, pairs ...string) {
	names := make([]string, 0, len(pairs)/2)
	values := make([]string, 0, len(pairs)/2)
	for i := 0; i < len(pairs)-1; i += 2 {
		names = append(names, pairs[i])
		values = append(values, pairs[i+1])
	}
	c.SetParamNames(names...)
	c.SetParamValues(values...)
}

// SetStaffID は Echo コンテキストに staff_id をセットする。
func SetStaffID(c echo.Context, staffID uuid.UUID) {
	c.Set("staff_id", staffID.String())
}
