package handler

import (
	"net/http"
	"os"
	"time"

	"tiara-web-app/backend/internal/domain"
	"tiara-web-app/backend/internal/usecase"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// StaffPortalHandler はスタッフポータル関連のHTTPハンドラー。
type StaffPortalHandler struct {
	authUsecase   usecase.StaffAuthUsecase
	portalUsecase usecase.StaffPortalUsecase
}

// NewStaffPortalHandler は新しいStaffPortalHandlerのインスタンスを作成する。
func NewStaffPortalHandler(auth usecase.StaffAuthUsecase, portal usecase.StaffPortalUsecase) *StaffPortalHandler {
	return &StaffPortalHandler{authUsecase: auth, portalUsecase: portal}
}

// --- Helper ---

// getStaffIDFromContext はミドルウェアでセットされたスタッフIDを取得する。
func getStaffIDFromContext(c echo.Context) (uuid.UUID, error) {
	staffIDStr, ok := c.Get("staff_id").(string)
	if !ok {
		return uuid.UUID{}, echo.NewHTTPError(http.StatusUnauthorized, "staff_id not found in context")
	}
	return uuid.Parse(staffIDStr)
}

// --- Auth ---

// StaffLoginRequest はスタッフログインリクエストのボディ型。
type StaffLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// StaffLoginResponse はスタッフログイン成功時のレスポンス型。
type StaffLoginResponse struct {
	Token   string `json:"token"`
	StaffID string `json:"staffId"`
}

// Login はスタッフのユーザー名とパスワードで認証し、JWTトークンを返す。
// 管理者トークンと区別するため、type: "staff" クレームを含む。
func (h *StaffPortalHandler) Login(c echo.Context) error {
	var req StaffLoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "username and password are required"})
	}

	account, err := h.authUsecase.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid credentials"})
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "tiara-dev-secret-key"
	}

	claims := jwt.MapClaims{
		"sub":      account.StaffID.String(),
		"username": account.Username,
		"type":     "staff", // 管理者トークンと区別するためのクレーム
		"exp":      time.Now().Add(24 * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
	}

	return c.JSON(http.StatusOK, StaffLoginResponse{
		Token:   tokenString,
		StaffID: account.StaffID.String(),
	})
}

// Verify はスタッフJWTトークンの有効性を確認する。
func (h *StaffPortalHandler) Verify(c echo.Context) error {
	staffID, err := getStaffIDFromContext(c)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token"})
	}
	return c.JSON(http.StatusOK, map[string]string{
		"status":  "ok",
		"staffId": staffID.String(),
	})
}

// --- Profile Draft ---

// ProfileDraftResponse はプロフィール下書きのレスポンス型。
type ProfileDraftResponse struct {
	ID                string  `json:"id,omitempty"`
	StaffID           string  `json:"staffId"`
	Name              string  `json:"name"`
	Role              string  `json:"role"`
	Bio               string  `json:"bio"`
	ImageURL          string  `json:"imageUrl"`
	ImageCropPosition string  `json:"imageCropPosition"`
	Status            string  `json:"status"`
	AdminComment      string  `json:"adminComment"`
	SubmittedAt       *string `json:"submittedAt,omitempty"`
	ReviewedAt        *string `json:"reviewedAt,omitempty"`
	CreatedAt         string  `json:"createdAt,omitempty"`
	UpdatedAt         string  `json:"updatedAt,omitempty"`
}

// toProfileDraftResponse は domain.StaffProfileDraft をレスポンス型に変換する。
func toProfileDraftResponse(d domain.StaffProfileDraft) ProfileDraftResponse {
	resp := ProfileDraftResponse{
		StaffID:           d.StaffID.String(),
		Name:              d.Name,
		Role:              d.Role,
		Bio:               d.Bio,
		ImageURL:          d.ImageURL,
		ImageCropPosition: d.ImageCropPosition,
		Status:            d.Status,
		AdminComment:      d.AdminComment,
	}
	if d.ID != uuid.Nil {
		resp.ID = d.ID.String()
	}
	if !d.CreatedAt.IsZero() {
		resp.CreatedAt = d.CreatedAt.Format(time.RFC3339)
	}
	if !d.UpdatedAt.IsZero() {
		resp.UpdatedAt = d.UpdatedAt.Format(time.RFC3339)
	}
	if d.SubmittedAt != nil {
		s := d.SubmittedAt.Format(time.RFC3339)
		resp.SubmittedAt = &s
	}
	if d.ReviewedAt != nil {
		s := d.ReviewedAt.Format(time.RFC3339)
		resp.ReviewedAt = &s
	}
	return resp
}

// GetMyProfileDraft は自分のプロフィール下書きを取得する。
func (h *StaffPortalHandler) GetMyProfileDraft(c echo.Context) error {
	staffID, err := getStaffIDFromContext(c)
	if err != nil {
		return err
	}

	draft, err := h.portalUsecase.GetMyProfileDraft(c.Request().Context(), staffID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toProfileDraftResponse(draft))
}

// SaveProfileDraftRequest はプロフィール下書き保存リクエストのボディ型。
type SaveProfileDraftRequest struct {
	Name              string `json:"name"`
	Role              string `json:"role"`
	Bio               string `json:"bio"`
	ImageURL          string `json:"imageUrl"`
	ImageCropPosition string `json:"imageCropPosition"`
}

// SaveMyProfileDraft はプロフィール下書きを保存する。
func (h *StaffPortalHandler) SaveMyProfileDraft(c echo.Context) error {
	staffID, err := getStaffIDFromContext(c)
	if err != nil {
		return err
	}

	var req SaveProfileDraftRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	input := domain.SaveProfileDraftInput{
		Name:              req.Name,
		Role:              req.Role,
		Bio:               req.Bio,
		ImageURL:          req.ImageURL,
		ImageCropPosition: req.ImageCropPosition,
	}

	draft, err := h.portalUsecase.SaveProfileDraft(c.Request().Context(), staffID, input)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toProfileDraftResponse(draft))
}

// SubmitMyProfileDraft はプロフィール下書きを承認申請する。
func (h *StaffPortalHandler) SubmitMyProfileDraft(c echo.Context) error {
	staffID, err := getStaffIDFromContext(c)
	if err != nil {
		return err
	}

	draftIDStr := c.Param("id")
	draftID, err := uuid.Parse(draftIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid draft ID"})
	}

	draft, err := h.portalUsecase.SubmitProfileDraft(c.Request().Context(), staffID, draftID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toProfileDraftResponse(draft))
}

// --- Schedule Draft ---

// ScheduleDraftItemRequest はスケジュール下書きアイテムのリクエスト型。
type ScheduleDraftItemRequest struct {
	DayOfWeek int    `json:"dayOfWeek"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// SaveScheduleDraftRequest はスケジュール下書き保存リクエストのボディ型。
type SaveScheduleDraftRequest struct {
	Items []ScheduleDraftItemRequest `json:"items"`
}

// ScheduleDraftItemResponse はスケジュール下書きアイテムのレスポンス型。
type ScheduleDraftItemResponse struct {
	ID        string `json:"id,omitempty"`
	DayOfWeek int    `json:"dayOfWeek"`
	StartTime string `json:"startTime"`
	EndTime   string `json:"endTime"`
}

// ScheduleDraftResponse はスケジュール下書きのレスポンス型。
type ScheduleDraftResponse struct {
	ID           string                      `json:"id,omitempty"`
	StaffID      string                      `json:"staffId"`
	Status       string                      `json:"status"`
	AdminComment string                      `json:"adminComment"`
	SubmittedAt  *string                     `json:"submittedAt,omitempty"`
	ReviewedAt   *string                     `json:"reviewedAt,omitempty"`
	CreatedAt    string                      `json:"createdAt,omitempty"`
	UpdatedAt    string                      `json:"updatedAt,omitempty"`
	Items        []ScheduleDraftItemResponse `json:"items"`
}

// toScheduleDraftResponse は domain.StaffScheduleDraft をレスポンス型に変換する。
func toScheduleDraftResponse(d domain.StaffScheduleDraft) ScheduleDraftResponse {
	resp := ScheduleDraftResponse{
		StaffID:      d.StaffID.String(),
		Status:       d.Status,
		AdminComment: d.AdminComment,
		Items:        make([]ScheduleDraftItemResponse, len(d.Items)),
	}
	if d.ID != uuid.Nil {
		resp.ID = d.ID.String()
	}
	if !d.CreatedAt.IsZero() {
		resp.CreatedAt = d.CreatedAt.Format(time.RFC3339)
	}
	if !d.UpdatedAt.IsZero() {
		resp.UpdatedAt = d.UpdatedAt.Format(time.RFC3339)
	}
	if d.SubmittedAt != nil {
		s := d.SubmittedAt.Format(time.RFC3339)
		resp.SubmittedAt = &s
	}
	if d.ReviewedAt != nil {
		s := d.ReviewedAt.Format(time.RFC3339)
		resp.ReviewedAt = &s
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

// GetMyScheduleDraft は自分のスケジュール下書きを取得する。
func (h *StaffPortalHandler) GetMyScheduleDraft(c echo.Context) error {
	staffID, err := getStaffIDFromContext(c)
	if err != nil {
		return err
	}

	draft, err := h.portalUsecase.GetMyScheduleDraft(c.Request().Context(), staffID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toScheduleDraftResponse(draft))
}

// SaveMyScheduleDraft はスケジュール下書きを保存する。
func (h *StaffPortalHandler) SaveMyScheduleDraft(c echo.Context) error {
	staffID, err := getStaffIDFromContext(c)
	if err != nil {
		return err
	}

	var req SaveScheduleDraftRequest
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

	draft, err := h.portalUsecase.SaveScheduleDraft(c.Request().Context(), staffID, items)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toScheduleDraftResponse(draft))
}

// SubmitMyScheduleDraft はスケジュール下書きを承認申請する。
func (h *StaffPortalHandler) SubmitMyScheduleDraft(c echo.Context) error {
	staffID, err := getStaffIDFromContext(c)
	if err != nil {
		return err
	}

	draftIDStr := c.Param("id")
	draftID, err := uuid.Parse(draftIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid draft ID"})
	}

	draft, err := h.portalUsecase.SubmitScheduleDraft(c.Request().Context(), staffID, draftID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, toScheduleDraftResponse(draft))
}
