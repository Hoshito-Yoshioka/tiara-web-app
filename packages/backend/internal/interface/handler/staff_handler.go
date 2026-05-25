package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"tiara-web-app/backend/internal/domain"
	"tiara-web-app/backend/internal/usecase"

	"github.com/google/uuid"
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
// page クエリパラメータが指定された場合はページネーション付きレスポンスを返す。
// 未指定の場合は全件を配列で返す（後方互換）。
func (h *StaffHandler) ListStaffs(c echo.Context) error {
	pageStr := c.QueryParam("page")
	if pageStr != "" {
		page, err := strconv.Atoi(pageStr)
		if err != nil || page < 1 {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "page must be a positive integer"})
		}

		const perPage = 10
		result, err := h.staffUsecase.ListStaffsPaginated(c.Request().Context(), page, perPage)
		if err != nil {
			return handleError(c, err)
		}
		return c.JSON(http.StatusOK, result)
	}

	staffs, err := h.staffUsecase.ListStaffs(c.Request().Context())
	if err != nil {
		return handleError(c, err)
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
		return handleError(c, err)
	}
	return c.JSON(http.StatusOK, result)
}

// ListAllStaffsWithSchedules は全スタッフの出勤スケジュール一覧を返すハンドラー。
// Schedule ページ用。
func (h *StaffHandler) ListAllStaffsWithSchedules(c echo.Context) error {
	result, err := h.staffUsecase.ListAllStaffsWithSchedules(c.Request().Context())
	if err != nil {
		return handleError(c, err)
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
	ShopID            string          `json:"shopId"`
	Name              string          `json:"name"`
	Role              string          `json:"role"`
	Bio               string          `json:"bio"`
	ImageURL          string          `json:"imageUrl"`
	ImageCropPosition string          `json:"imageCropPosition"`
	SortOrder         int             `json:"sortOrder"`
	Schedules         []ScheduleInput `json:"schedules"`
}

// UpdateStaffRequest はスタッフ更新リクエストのボディ型。
type UpdateStaffRequest struct {
	Name              string          `json:"name"`
	Role              string          `json:"role"`
	Bio               string          `json:"bio"`
	ImageURL          string          `json:"imageUrl"`
	ImageCropPosition string          `json:"imageCropPosition"`
	SortOrder         int             `json:"sortOrder"`
	Schedules         []ScheduleInput `json:"schedules"`
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
		ShopID:            req.ShopID,
		Name:              req.Name,
		Role:              req.Role,
		Bio:               req.Bio,
		ImageURL:          req.ImageURL,
		ImageCropPosition: req.ImageCropPosition,
		SortOrder:         req.SortOrder,
		Schedules:         toScheduleInputs(req.Schedules),
	}

	result, err := h.staffUsecase.CreateStaff(c.Request().Context(), input)
	if err != nil {
		return handleError(c, err)
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
		Name:              req.Name,
		Role:              req.Role,
		Bio:               req.Bio,
		ImageURL:          req.ImageURL,
		ImageCropPosition: req.ImageCropPosition,
		SortOrder:         req.SortOrder,
		Schedules:         toScheduleInputs(req.Schedules),
	}

	result, err := h.staffUsecase.UpdateStaff(c.Request().Context(), id, input)
	if err != nil {
		return handleError(c, err)
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
		return handleError(c, err)
	}
	return c.NoContent(http.StatusNoContent)
}

// UploadStaffImage はスタッフ画像をアップロードするハンドラー。
// multipart/form-data で受け取り、ファイルを uploads/ に保存する。
func (h *StaffHandler) UploadStaffImage(c echo.Context) error {
	staffID := c.Param("id")
	if staffID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "staff id is required"})
	}

	file, err := c.FormFile("image")
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "image file is required"})
	}

	// ファイルバリデーション（サイズ・拡張子・MIMEタイプ）
	if err := validateImageFile(file); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	src, err := file.Open()
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to read file"})
	}
	defer src.Close() //nolint:errcheck // best-effort close

	// uploads ディレクトリ作成
	uploadDir := "uploads/staff"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to create upload directory"})
	}

	// ユニークファイル名生成（拡張子を正規化）
	ext := strings.ToLower(filepath.Ext(file.Filename))
	filename := fmt.Sprintf("%s%s", uuid.New().String(), ext)
	dstPath := filepath.Join(uploadDir, filename)

	dst, err := os.Create(dstPath)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to save file"})
	}
	defer dst.Close() //nolint:errcheck // best-effort close

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to write file"})
	}

	// DB に画像レコードを作成
	isMain := c.FormValue("isMain") == "true"
	imageURL := fmt.Sprintf("/uploads/staff/%s", filename)

	image, err := h.staffUsecase.UploadStaffImage(c.Request().Context(), staffID, imageURL, isMain, 0)
	if err != nil {
		// ファイルを削除
		_ = os.Remove(dstPath)
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, image)
}

// DeleteStaffImage はスタッフ画像を削除するハンドラー。
func (h *StaffHandler) DeleteStaffImage(c echo.Context) error {
	imageID := c.Param("imageId")
	if imageID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "image id is required"})
	}

	err := h.staffUsecase.DeleteStaffImage(c.Request().Context(), imageID)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.NoContent(http.StatusNoContent)
}

// SetMainImageRequest はメイン画像設定リクエストのボディ型。
type SetMainImageRequest struct {
	ImageID string `json:"imageId"`
}

// SetMainImage はスタッフのメイン画像を設定するハンドラー。
func (h *StaffHandler) SetMainImage(c echo.Context) error {
	staffID := c.Param("id")
	if staffID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "staff id is required"})
	}

	var req SetMainImageRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	image, err := h.staffUsecase.SetMainImage(c.Request().Context(), staffID, req.ImageID)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.JSON(http.StatusOK, image)
}

// UpdateImageCropPositionRequest はクロップ位置更新リクエストのボディ型。
type UpdateImageCropPositionRequest struct {
	CropPosition string `json:"cropPosition"`
}

// UpdateImageCropPosition は画像のクロップ位置を更新するハンドラー。
func (h *StaffHandler) UpdateImageCropPosition(c echo.Context) error {
	imageID := c.Param("imageId")
	if imageID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "image id is required"})
	}

	var req UpdateImageCropPositionRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	image, err := h.staffUsecase.UpdateImageCropPosition(c.Request().Context(), imageID, req.CropPosition)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.JSON(http.StatusOK, image)
}
