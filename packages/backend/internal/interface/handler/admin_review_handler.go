package handler

import (
	"net/http"
	"time"

	"tiara-web-app/backend/internal/domain"
	"tiara-web-app/backend/internal/usecase"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// AdminReviewHandler は管理者による下書きレビューのHTTPハンドラー。
type AdminReviewHandler struct {
	reviewUsecase usecase.AdminReviewUsecase
}

// NewAdminReviewHandler は新しいAdminReviewHandlerのインスタンスを作成する。
func NewAdminReviewHandler(us usecase.AdminReviewUsecase) *AdminReviewHandler {
	return &AdminReviewHandler{reviewUsecase: us}
}

// --- Profile Draft Review ---

// StaffImageJSON はスタッフ画像のJSON表現。
type StaffImageJSON struct {
	ID           string `json:"id"`
	StaffID      string `json:"staffId"`
	ImageURL     string `json:"imageUrl"`
	IsMain       bool   `json:"isMain"`
	SortOrder    int    `json:"sortOrder"`
	CropPosition string `json:"cropPosition"`
}

// PendingProfileDraftResponse は管理者向けの承認待ちプロフィール下書きレスポンス。
type PendingProfileDraftResponse struct {
	ID                  string           `json:"id"`
	StaffID             string           `json:"staffId"`
	StaffName           string           `json:"staffName"`
	Name                string           `json:"name"`
	Role                string           `json:"role"`
	Bio                 string           `json:"bio"`
	ImageURL            string           `json:"imageUrl"`
	ExternalScheduleURL string           `json:"externalScheduleUrl"`
	ImageCropPosition   string           `json:"imageCropPosition"`
	Status              string           `json:"status"`
	AdminComment        string           `json:"adminComment,omitempty"`
	SubmittedAt         *string          `json:"submittedAt,omitempty"`
	CreatedAt           string           `json:"createdAt"`
	Images              []StaffImageJSON `json:"images"`
}

// toPendingProfileDraftResponse は domain を管理者向けレスポンスに変換する。
func toPendingProfileDraftResponse(d domain.StaffProfileDraft) PendingProfileDraftResponse {
	resp := PendingProfileDraftResponse{
		ID:                  d.ID.String(),
		StaffID:             d.StaffID.String(),
		Name:                d.Name,
		Role:                d.Role,
		Bio:                 d.Bio,
		ImageURL:            d.ImageURL,
		ExternalScheduleURL: d.ExternalScheduleURL,
		ImageCropPosition:   d.ImageCropPosition,
		Status:              string(d.Status),
		AdminComment:        d.AdminComment,
		CreatedAt:           d.CreatedAt.Format(time.RFC3339),
	}
	if d.SubmittedAt != nil {
		s := d.SubmittedAt.Format(time.RFC3339)
		resp.SubmittedAt = &s
	}
	return resp
}

// ListPendingProfileDrafts は承認待ちのプロフィール下書き一覧を返す。
func (h *AdminReviewHandler) ListPendingProfileDrafts(c echo.Context) error {
	drafts, err := h.reviewUsecase.ListPendingProfileDrafts(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}

	resp := make([]PendingProfileDraftResponse, len(drafts))
	for i, d := range drafts {
		resp[i] = toPendingProfileDraftResponse(d)
		// スタッフ名を取得して設定
		staffName, _ := h.reviewUsecase.GetStaffName(c.Request().Context(), d.StaffID)
		resp[i].StaffName = staffName
		// スタッフ画像一覧を取得して設定
		images, _ := h.reviewUsecase.ListImagesByStaffID(c.Request().Context(), d.StaffID)
		resp[i].Images = toStaffImageJSONList(images)
	}
	return c.JSON(http.StatusOK, resp)
}

// toStaffImageJSONList は domain.StaffImage のスライスを JSON 応答用に変換する。
func toStaffImageJSONList(images []domain.StaffImage) []StaffImageJSON {
	result := make([]StaffImageJSON, len(images))
	for i, img := range images {
		result[i] = StaffImageJSON{
			ID:           img.ID.String(),
			StaffID:      img.StaffID.String(),
			ImageURL:     img.ImageURL,
			IsMain:       img.IsMain,
			SortOrder:    img.SortOrder,
			CropPosition: img.CropPosition,
		}
	}
	return result
}

// ReviewDraftRequest はレビューリクエストのボディ型。
type ReviewDraftRequest struct {
	Status       string `json:"status"`
	AdminComment string `json:"adminComment"`
	UpdatedAt    string `json:"updatedAt"`
}

// ReviewProfileDraft は管理者がプロフィール下書きをレビューする。
func (h *AdminReviewHandler) ReviewProfileDraft(c echo.Context) error {
	draftIDStr := c.Param("id")
	draftID, err := uuid.Parse(draftIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid draft ID"})
	}

	var req ReviewDraftRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	input := domain.ReviewDraftInput{
		Status:       domain.DraftStatus(req.Status),
		AdminComment: req.AdminComment,
	}

	updatedAt, _ := time.Parse(time.RFC3339, req.UpdatedAt)
	draft, err := h.reviewUsecase.ReviewProfileDraft(c.Request().Context(), draftID, input, updatedAt)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, toProfileDraftResponse(draft))
}

// --- Schedule Draft Review ---

// PendingScheduleDraftResponse は管理者向けの承認待ちスケジュール下書きレスポンス。
type PendingScheduleDraftResponse struct {
	ID           string                      `json:"id"`
	StaffID      string                      `json:"staffId"`
	StaffName    string                      `json:"staffName"`
	Status       string                      `json:"status"`
	AdminComment string                      `json:"adminComment,omitempty"`
	SubmittedAt  *string                     `json:"submittedAt,omitempty"`
	CreatedAt    string                      `json:"createdAt"`
	Items        []ScheduleDraftItemResponse `json:"items"`
}

// toPendingScheduleDraftResponse は domain を管理者向けレスポンスに変換する。
func toPendingScheduleDraftResponse(d domain.StaffScheduleDraft) PendingScheduleDraftResponse {
	resp := PendingScheduleDraftResponse{
		ID:           d.ID.String(),
		StaffID:      d.StaffID.String(),
		Status:       string(d.Status),
		AdminComment: d.AdminComment,
		CreatedAt:    d.CreatedAt.Format(time.RFC3339),
		Items:        make([]ScheduleDraftItemResponse, len(d.Items)),
	}
	if d.SubmittedAt != nil {
		s := d.SubmittedAt.Format(time.RFC3339)
		resp.SubmittedAt = &s
	}
	for i, item := range d.Items {
		resp.Items[i] = ScheduleDraftItemResponse{
			DayOfWeek: item.DayOfWeek,
			StartTime: item.StartTime.Format("15:04"),
			EndTime:   item.EndTime.Format("15:04"),
		}
		if item.ID != uuid.Nil {
			resp.Items[i].ID = item.ID.String()
		}
	}
	return resp
}

// ListPendingScheduleDrafts は承認待ちのスケジュール下書き一覧を返す。
func (h *AdminReviewHandler) ListPendingScheduleDrafts(c echo.Context) error {
	drafts, err := h.reviewUsecase.ListPendingScheduleDrafts(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}

	resp := make([]PendingScheduleDraftResponse, len(drafts))
	for i, d := range drafts {
		resp[i] = toPendingScheduleDraftResponse(d)
		// スタッフ名を取得して設定
		staffName, _ := h.reviewUsecase.GetStaffName(c.Request().Context(), d.StaffID)
		resp[i].StaffName = staffName
	}
	return c.JSON(http.StatusOK, resp)
}

// ListApprovedScheduleDrafts は承認済み（未反映）のスケジュール下書き一覧を返す。
func (h *AdminReviewHandler) ListApprovedScheduleDrafts(c echo.Context) error {
	drafts, err := h.reviewUsecase.ListApprovedScheduleDrafts(c.Request().Context())
	if err != nil {
		return handleError(c, err)
	}

	resp := make([]PendingScheduleDraftResponse, len(drafts))
	for i, d := range drafts {
		resp[i] = toPendingScheduleDraftResponse(d)
		staffName, _ := h.reviewUsecase.GetStaffName(c.Request().Context(), d.StaffID)
		resp[i].StaffName = staffName
	}
	return c.JSON(http.StatusOK, resp)
}

// PublishScheduleDraft は承認済みスケジュール下書きをライブデータに反映する。
func (h *AdminReviewHandler) PublishScheduleDraft(c echo.Context) error {
	draftIDStr := c.Param("id")
	draftID, err := uuid.Parse(draftIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid draft ID"})
	}

	err = h.reviewUsecase.PublishScheduleDraft(c.Request().Context(), draftID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]string{"message": "スケジュールを店舗ページに反映しました"})
}

// ReviewScheduleDraft は管理者がスケジュール下書きをレビューする。
func (h *AdminReviewHandler) ReviewScheduleDraft(c echo.Context) error {
	draftIDStr := c.Param("id")
	draftID, err := uuid.Parse(draftIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid draft ID"})
	}

	var req ReviewDraftRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	input := domain.ReviewDraftInput{
		Status:       domain.DraftStatus(req.Status),
		AdminComment: req.AdminComment,
	}

	schedUpdatedAt, _ := time.Parse(time.RFC3339, req.UpdatedAt)
	draft, err := h.reviewUsecase.ReviewScheduleDraft(c.Request().Context(), draftID, input, schedUpdatedAt)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, toScheduleDraftResponse(draft))
}

// --- 単体取得・内容修正 ---

// GetProfileDraft は管理者がプロフィール下書きを単体で取得する。
func (h *AdminReviewHandler) GetProfileDraft(c echo.Context) error {
	draftIDStr := c.Param("id")
	draftID, err := uuid.Parse(draftIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid draft ID"})
	}

	draft, err := h.reviewUsecase.GetProfileDraft(c.Request().Context(), draftID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, toProfileDraftResponse(draft))
}

// UpdateProfileDraftContent は管理者がプロフィール下書きの内容を修正する（ステータス変更なし）。
func (h *AdminReviewHandler) UpdateProfileDraftContent(c echo.Context) error {
	draftIDStr := c.Param("id")
	draftID, err := uuid.Parse(draftIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid draft ID"})
	}

	var req SaveProfileDraftRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	input := domain.SaveProfileDraftInput{
		Name:                req.Name,
		Role:                req.Role,
		Bio:                 req.Bio,
		ImageURL:            req.ImageURL,
		ExternalScheduleURL: req.ExternalScheduleURL,
		ImageCropPosition:   req.ImageCropPosition,
	}

	contentUpdatedAt, _ := time.Parse(time.RFC3339, req.UpdatedAt)
	draft, err := h.reviewUsecase.UpdateProfileDraftContent(c.Request().Context(), draftID, input, contentUpdatedAt)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, toProfileDraftResponse(draft))
}

// GetScheduleDraft は管理者がスケジュール下書きを単体で取得する。
func (h *AdminReviewHandler) GetScheduleDraft(c echo.Context) error {
	draftIDStr := c.Param("id")
	draftID, err := uuid.Parse(draftIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid draft ID"})
	}

	draft, err := h.reviewUsecase.GetScheduleDraft(c.Request().Context(), draftID)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, toScheduleDraftResponse(draft))
}

// UpdateScheduleDraftContentRequest はスケジュール下書き内容修正リクエストのボディ型。
type UpdateScheduleDraftContentRequest struct {
	Items     []ScheduleDraftItemRequest `json:"items"`
	UpdatedAt string                     `json:"updatedAt"`
}

// UpdateScheduleDraftContent は管理者がスケジュール下書きの内容を修正する（ステータス変更なし）。
func (h *AdminReviewHandler) UpdateScheduleDraftContent(c echo.Context) error {
	draftIDStr := c.Param("id")
	draftID, err := uuid.Parse(draftIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid draft ID"})
	}

	var req UpdateScheduleDraftContentRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	items := make([]domain.ScheduleDraftItem, len(req.Items))
	for i, item := range req.Items {
		startTime, err := time.Parse("15:04", item.StartTime)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid startTime format, use HH:MM"})
		}
		endTime, err := time.Parse("15:04", item.EndTime)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid endTime format, use HH:MM"})
		}
		items[i] = domain.ScheduleDraftItem{
			DayOfWeek: item.DayOfWeek,
			StartTime: startTime,
			EndTime:   endTime,
		}
	}

	schedContentUpdatedAt, _ := time.Parse(time.RFC3339, req.UpdatedAt)
	draft, err := h.reviewUsecase.UpdateScheduleDraftContent(c.Request().Context(), draftID, items, schedContentUpdatedAt)
	if err != nil {
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, toScheduleDraftResponse(draft))
}
