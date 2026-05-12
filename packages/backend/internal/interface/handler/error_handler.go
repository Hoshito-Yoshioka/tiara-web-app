package handler

import (
	"errors"
	"net/http"

	"tiara-web-app/backend/internal/domain"

	"github.com/labstack/echo/v4"
)

// handleError はドメインエラーを適切な HTTP ステータスコードにマッピングして返す。
// usecase 層が返すセンチネルエラーを errors.Is() で判定し、
// handler 層のエラーレスポンスを統一する。
func handleError(c echo.Context, err error) error {
	switch {
	case errors.Is(err, domain.ErrNotFound):
		return c.JSON(http.StatusNotFound, map[string]string{"error": err.Error()})
	case errors.Is(err, domain.ErrUnauthorized):
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
	case errors.Is(err, domain.ErrForbidden):
		return c.JSON(http.StatusForbidden, map[string]string{"error": err.Error()})
	case errors.Is(err, domain.ErrConflict):
		return c.JSON(http.StatusConflict, map[string]string{"error": err.Error()})
	case errors.Is(err, domain.ErrInvalidInput):
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	default:
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
}
