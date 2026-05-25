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

// ShopResponse は店舗情報のレスポンス型。
// OpeningTime/ClosingTime は time.Time を "HH:MM" 形式にフォーマットして返す。
type ShopResponse struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Address     string `json:"address"`
	OpeningTime string `json:"openingTime"`
	ClosingTime string `json:"closingTime"`
}

// toShopResponse は domain.Shop をレスポンス型に変換する。
// time.Time の OpeningTime/ClosingTime を "15:04" フォーマットで文字列化する。
func toShopResponse(s domain.Shop) ShopResponse {
	return ShopResponse{
		ID:          s.ID.String(),
		Name:        s.Name,
		Address:     s.Address,
		OpeningTime: s.OpeningTime.Format("15:04"),
		ClosingTime: s.ClosingTime.Format("15:04"),
	}
}

func (h *ShopHandler) ListShops(c echo.Context) error {
	shops, err := h.shopUsecase.ListShops(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}

	resp := make([]ShopResponse, len(shops))
	for i, s := range shops {
		resp[i] = toShopResponse(s)
	}
	return c.JSON(http.StatusOK, resp)
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
	return c.JSON(http.StatusOK, toShopResponse(shop))
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
	return c.JSON(http.StatusOK, toShopResponse(shop))
}
