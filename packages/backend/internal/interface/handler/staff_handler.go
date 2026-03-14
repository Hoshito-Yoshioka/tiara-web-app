package handler

import (
	"net/http"
	"tiara-web-app/backend/internal/domain"
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

// ScheduleInput はスケジュール入力のJSONマッピング型。
type ScheduleInput struct {
	DayOfWeek int    `json:"dayOfWeek"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// CreateStaffRequest はスタッフ作成リクエストのボディ型。
type CreateStaffRequest struct {
	ShopID    string          `json:"shopId"`
	Name      string          `json:"name"`
	Role      string          `json:"role"`
	Bio       string          `json:"bio"`
	ImageURL  string          `json:"imageUrl"`
	SortOrder int             `json:"sortOrder"`
	Schedules []ScheduleInput `json:"schedules"`
}

// UpdateStaffRequest はスタッフ更新リクエストのボディ型。
type UpdateStaffRequest struct {
	Name      string          `json:"name"`
	Role      string          `json:"role"`
	Bio       string          `json:"bio"`
	ImageURL  string          `json:"imageUrl"`
	SortOrder int             `json:"sortOrder"`
	Schedules []ScheduleInput `json:"schedules"`
}

func toScheduleInputs(inputs []ScheduleInput) []domain.ScheduleInput {
	result := make([]domain.ScheduleInput, len(inputs))
	for i, s := range inputs {
		result[i] = domain.ScheduleInput{
			DayOfWeek: s.DayOfWeek,
			StartTime: s.StartTime,
			EndTime:   s.EndTime,
		}
	}
	return result
}

// CreateStaff は新しいスタッフを作成するハンドラー。
func (h *StaffHandler) CreateStaff(c echo.Context) error {
	var req CreateStaffRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	input := domain.CreateStaffInput{
		ShopID:    req.ShopID,
		Name:      req.Name,
		Role:      req.Role,
		Bio:       req.Bio,
		ImageURL:  req.ImageURL,
		SortOrder: req.SortOrder,
		Schedules: toScheduleInputs(req.Schedules),
	}

	result, err := h.staffUsecase.CreateStaff(c.Request().Context(), input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusCreated, result)
}

// UpdateStaff はスタッフ情報を更新するハンドラー。
func (h *StaffHandler) UpdateStaff(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	var req UpdateStaffRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	input := domain.UpdateStaffInput{
		Name:      req.Name,
		Role:      req.Role,
		Bio:       req.Bio,
		ImageURL:  req.ImageURL,
		SortOrder: req.SortOrder,
		Schedules: toScheduleInputs(req.Schedules),
	}

	result, err := h.staffUsecase.UpdateStaff(c.Request().Context(), id, input)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.JSON(http.StatusOK, result)
}

// DeleteStaff はスタッフを削除するハンドラー。
func (h *StaffHandler) DeleteStaff(c echo.Context) error {
	id := c.Param("id")
	if id == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "id is required"})
	}

	err := h.staffUsecase.DeleteStaff(c.Request().Context(), id)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}
	return c.NoContent(http.StatusNoContent)
}
