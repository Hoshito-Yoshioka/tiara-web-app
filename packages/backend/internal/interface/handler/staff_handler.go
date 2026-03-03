package handler

import (
	"net/http"
	"tiara-web-app/backend/internal/usecase"

	"github.com/labstack/echo/v4"
)

// StaffHandler はスタッフ関連のHTTPハンドラー。
type StaffHandler struct {
	staffUsecase usecase.StaffUsecase
}

// NewStaffHandler は新しいStaffHandlerのインスタンスを作成する。
func NewStaffHandler(us usecase.StaffUsecase) *StaffHandler {
	return &StaffHandler{staffUsecase: us}
}

// ListStaffs はスタッフ一覧を返すハンドラー。
func (h *StaffHandler) ListStaffs(c echo.Context) error {
	staffs, err := h.staffUsecase.ListStaffs(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, staffs)
}

// GetStaffWithSchedules はスタッフ詳細（出勤スケジュール付き）を返すハンドラー。
func (h *StaffHandler) GetStaffWithSchedules(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	result, err := h.staffUsecase.GetStaffWithSchedules(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusNotFound, map[string]string{"error": "staff not found"})
	}
	return c.JSON(http.StatusOK, result)
}

// ListAllStaffsWithSchedules は全スタッフの出勤スケジュール一覧を返すハンドラー。
// Schedule ページ用。
func (h *StaffHandler) ListAllStaffsWithSchedules(c echo.Context) error {
	result, err := h.staffUsecase.ListAllStaffsWithSchedules(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}
