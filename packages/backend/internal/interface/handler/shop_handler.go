package handler

import (
	"net/http"
	"tiara-web-app/backend/internal/domain"
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
		return handleError(c, err)
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
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, shop)
}

// UpdateShopRequest は店舗更新リクエストのボディ型。
type UpdateShopRequest struct {
	Name        string `json:"name"`
	Address     string `json:"address"`
	OpeningTime string `json:"openingTime"`
	ClosingTime string `json:"closingTime"`
}

// UpdateShop は店舗情報を更新するハンドラー。
func (h *ShopHandler) UpdateShop(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	var req UpdateShopRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	input := domain.UpdateShopInput{
		Name:        req.Name,
		Address:     req.Address,
		OpeningTime: req.OpeningTime,
		ClosingTime: req.ClosingTime,
	}

	shop, err := h.shopUsecase.UpdateShop(c.Request().Context(), id, input)
	if err != nil {
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, shop)
}
