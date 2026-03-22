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

// PendingProfileDraftResponse は管理者向けの承認待ちプロフィール下書きレスポンス。
type PendingProfileDraftResponse struct {
	ID                string  `json:"id"`
	StaffID           string  `json:"staffId"`
	Name              string  `json:"name"`
	Role              string  `json:"role"`
	Bio               string  `json:"bio"`
	ImageURL          string  `json:"imageUrl"`
	ImageCropPosition string  `json:"imageCropPosition"`
	Status            string  `json:"status"`
	AdminComment      string  `json:"adminComment"`
	SubmittedAt       *string `json:"submittedAt,omitempty"`
	CreatedAt         string  `json:"createdAt"`
}

// toPendingProfileDraftResponse は domain を管理者向けレスポンスに変換する。
func toPendingProfileDraftResponse(d domain.StaffProfileDraft) PendingProfileDraftResponse {
	resp := PendingProfileDraftResponse{
		ID:                d.ID.String(),
		StaffID:           d.StaffID.String(),
		Name:              d.Name,
		Role:              d.Role,
		Bio:               d.Bio,
		ImageURL:          d.ImageURL,
		ImageCropPosition: d.ImageCropPosition,
		Status:            d.Status,
		AdminComment:      d.AdminComment,
		CreatedAt:         d.CreatedAt.Format(time.RFC3339),
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
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	resp := make([]PendingProfileDraftResponse, len(drafts))
	for i, d := range drafts {
		resp[i] = toPendingProfileDraftResponse(d)
	}
	return c.JSON(http.StatusOK, resp)
}

// ReviewDraftRequest はレビューリクエストのボディ型。
type ReviewDraftRequest struct {
	Status       string `json:"status"`
	AdminComment string `json:"adminComment"`
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
		Status:       req.Status,
		AdminComment: req.AdminComment,
	}

	draft, err := h.reviewUsecase.ReviewProfileDraft(c.Request().Context(), draftID, input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toProfileDraftResponse(draft))
}

// --- Schedule Draft Review ---

// PendingScheduleDraftResponse は管理者向けの承認待ちスケジュール下書きレスポンス。
type PendingScheduleDraftResponse struct {
	ID           string                      `json:"id"`
	StaffID      string                      `json:"staffId"`
	Status       string                      `json:"status"`
	AdminComment string                      `json:"adminComment"`
	SubmittedAt  *string                     `json:"submittedAt,omitempty"`
	CreatedAt    string                      `json:"createdAt"`
	Items        []ScheduleDraftItemResponse `json:"items"`
}

// toPendingScheduleDraftResponse は domain を管理者向けレスポンスに変換する。
func toPendingScheduleDraftResponse(d domain.StaffScheduleDraft) PendingScheduleDraftResponse {
	resp := PendingScheduleDraftResponse{
		ID:           d.ID.String(),
		StaffID:      d.StaffID.String(),
		Status:       d.Status,
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
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	resp := make([]PendingScheduleDraftResponse, len(drafts))
	for i, d := range drafts {
		resp[i] = toPendingScheduleDraftResponse(d)
	}
	return c.JSON(http.StatusOK, resp)
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
		Status:       req.Status,
		AdminComment: req.AdminComment,
	}

	draft, err := h.reviewUsecase.ReviewScheduleDraft(c.Request().Context(), draftID, input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toScheduleDraftResponse(draft))
}
