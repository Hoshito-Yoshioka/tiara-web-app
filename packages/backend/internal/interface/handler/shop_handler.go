package handler

import (
	"net/http"
	"tiara-web-app/backend/internal/usecase"

	"github.com/labstack/echo/v4"
)

type ShopHandler struct {
	shopUsecase usecase.ShopUsecase
}

func NewShopHandler(us usecase.ShopUsecase) *ShopHandler {
	return &ShopHandler{
		shopUsecase: us,
	}
}

func (h *ShopHandler) ListShops(c echo.Context) error {
	shops, err := h.shopUsecase.ListShops(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, shops)
}

func (h *ShopHandler) GetShopByID(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	shop, err := h.shopUsecase.GetShopByID(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "shop not found"})
	}
	return c.JSON(http.StatusOK, shop)
}
