package handler

import (
	"net/http"
	"tiara-web-app/backend/internal/domain"
	"tiara-web-app/backend/internal/usecase"

	"github.com/labstack/echo/v4"
)

type MenuHandler struct {
	menuUsecase usecase.MenuUsecase
}

func NewMenuHandler(us usecase.MenuUsecase) *MenuHandler {
	return &MenuHandler{menuUsecase: us}
}

// GET /menus — カテゴリ＋アイテム一覧（公開）
func (h *MenuHandler) ListMenuCategoriesWithItems(c echo.Context) error {
	result, err := h.menuUsecase.ListMenuCategoriesWithItems(c.Request().Context())
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.JSON(http.StatusOK, result)
}

// ============================================================
// Admin: MenuCategory CRUD
// ============================================================

type CreateMenuCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SortOrder   int32  `json:"sortOrder"`
}

// POST /admin/menu/categories
func (h *MenuHandler) CreateMenuCategory(c echo.Context) error {
	var req CreateMenuCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	cat, err := h.menuUsecase.CreateMenuCategory(c.Request().Context(), domain.CreateMenuCategoryInput{
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	})
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.JSON(http.StatusCreated, cat)
}

type UpdateMenuCategoryRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	SortOrder   int32  `json:"sortOrder"`
}

// PUT /admin/menu/categories/:id
func (h *MenuHandler) UpdateMenuCategory(c echo.Context) error {
	id := c.Param("id")
	var req UpdateMenuCategoryRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	cat, err := h.menuUsecase.UpdateMenuCategory(c.Request().Context(), id, domain.UpdateMenuCategoryInput{
		Name:        req.Name,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	})
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.JSON(http.StatusOK, cat)
}

// DELETE /admin/menu/categories/:id
func (h *MenuHandler) DeleteMenuCategory(c echo.Context) error {
	id := c.Param("id")
	if err := h.menuUsecase.DeleteMenuCategory(c.Request().Context(), id); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.NoContent(http.StatusNoContent)
}

// ============================================================
// Admin: MenuItem CRUD
// ============================================================

type CreateMenuItemRequest struct {
	CategoryID  string `json:"categoryId"`
	Name        string `json:"name"`
	Price       string `json:"price"`
	Description string `json:"description"`
	SortOrder   int32  `json:"sortOrder"`
}

// POST /admin/menu/items
func (h *MenuHandler) CreateMenuItem(c echo.Context) error {
	var req CreateMenuItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	item, err := h.menuUsecase.CreateMenuItem(c.Request().Context(), domain.CreateMenuItemInput{
		CategoryID:  req.CategoryID,
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	})
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.JSON(http.StatusCreated, item)
}

type UpdateMenuItemRequest struct {
	Name        string `json:"name"`
	Price       string `json:"price"`
	Description string `json:"description"`
	SortOrder   int32  `json:"sortOrder"`
}

// PUT /admin/menu/items/:id
func (h *MenuHandler) UpdateMenuItem(c echo.Context) error {
	id := c.Param("id")
	var req UpdateMenuItemRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}
	item, err := h.menuUsecase.UpdateMenuItem(c.Request().Context(), id, domain.UpdateMenuItemInput{
		Name:        req.Name,
		Price:       req.Price,
		Description: req.Description,
		SortOrder:   req.SortOrder,
	})
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.JSON(http.StatusOK, item)
}

// DELETE /admin/menu/items/:id
func (h *MenuHandler) DeleteMenuItem(c echo.Context) error {
	id := c.Param("id")
	if err := h.menuUsecase.DeleteMenuItem(c.Request().Context(), id); err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.NoContent(http.StatusNoContent)
}
