package handler

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"tiara-web-app/backend/internal/domain"
	"tiara-web-app/backend/internal/usecase"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// StaffPortalHandler はスタッフポータル関連のHTTPハンドラー。
type StaffPortalHandler struct {
	authUsecase    usecase.StaffAuthUsecase
	portalUsecase  usecase.StaffPortalUsecase
	staffUsecase   usecase.StaffUsecase
	jwtSecret      string
	jwtExpiryHours int
	uploadDir      string
}

// NewStaffPortalHandler は新しいStaffPortalHandlerのインスタンスを作成する。
func NewStaffPortalHandler(auth usecase.StaffAuthUsecase, portal usecase.StaffPortalUsecase, staff usecase.StaffUsecase, jwtSecret string, jwtExpiryHours int, uploadDir string) *StaffPortalHandler {
	return &StaffPortalHandler{
		authUsecase:    auth,
		portalUsecase:  portal,
		staffUsecase:   staff,
		jwtSecret:      jwtSecret,
		jwtExpiryHours: jwtExpiryHours,
		uploadDir:      uploadDir,
	}
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
		return handleError(c, err)
	}

	jwtClaims := jwt.MapClaims{
		"sub":      account.StaffID.String(),
		"username": account.Username,
		"type":     "staff", // 管理者トークンと区別するためのクレーム
		"exp":      time.Now().Add(time.Duration(h.jwtExpiryHours) * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwtClaims)
	tokenString, err := token.SignedString([]byte(h.jwtSecret))
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
		return handleError(c, err)
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
		return handleError(c, err)
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
		return handleError(c, err)
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
		return handleError(c, err)
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
		return handleError(c, err)
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
		return handleError(c, err)
	}

	return c.JSON(http.StatusOK, toScheduleDraftResponse(draft))
}

// --- Image Management ---
// スタッフが自分自身の画像を管理するためのハンドラー群。
// JWT から取得した staff_id を使い、自分の画像のみ操作可能にする。

// ListMyImages は自分の画像一覧を取得する。
func (h *StaffPortalHandler) ListMyImages(c echo.Context) error {
	staffID, err := getStaffIDFromContext(c)
	if err != nil {
		return err
	}

	result, err := h.staffUsecase.GetStaffWithSchedules(c.Request().Context(), staffID.String())
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusOK, result.Images)
}

// UploadMyImage は自分のスタッフ画像をアップロードする。
func (h *StaffPortalHandler) UploadMyImage(c echo.Context) error {
	staffID, err := getStaffIDFromContext(c)
	if err != nil {
		return err
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
	defer src.Close()

	// uploads ディレクトリ作成
	uploadDir := filepath.Join(h.uploadDir, "staff")
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
	defer dst.Close()

	if _, err = io.Copy(dst, src); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to write file"})
	}

	isMain := c.FormValue("isMain") == "true"
	imageURL := fmt.Sprintf("/uploads/staff/%s", filename)

	image, err := h.staffUsecase.UploadStaffImage(c.Request().Context(), staffID.String(), imageURL, isMain, 0)
	if err != nil {
		os.Remove(dstPath)
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}

	return c.JSON(http.StatusCreated, image)
}

// verifyImageOwnership は指定された画像がスタッフの所有物かを検証する。
// スタッフの全画像を取得し、対象の画像IDが含まれるかチェックする。
func (h *StaffPortalHandler) verifyImageOwnership(c echo.Context, staffID uuid.UUID, imageID string) error {
	result, err := h.staffUsecase.GetStaffWithSchedules(c.Request().Context(), staffID.String())
	if err != nil {
		return fmt.Errorf("failed to verify image ownership")
	}
	for _, img := range result.Images {
		if img.ID.String() == imageID {
			return nil
		}
	}
	return fmt.Errorf("image does not belong to this staff")
}

// DeleteMyImage は自分のスタッフ画像を削除する。
func (h *StaffPortalHandler) DeleteMyImage(c echo.Context) error {
	staffID, err := getStaffIDFromContext(c)
	if err != nil {
		return err
	}

	imageID := c.Param("imageId")
	if imageID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "image id is required"})
	}

	// 所有権チェック
	if err := h.verifyImageOwnership(c, staffID, imageID); err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}

	err = h.staffUsecase.DeleteStaffImage(c.Request().Context(), imageID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.NoContent(http.StatusNoContent)
}

// SetMyMainImage は自分のメイン画像を設定する。
func (h *StaffPortalHandler) SetMyMainImage(c echo.Context) error {
	staffID, err := getStaffIDFromContext(c)
	if err != nil {
		return err
	}

	var req struct {
		ImageID string `json:"imageId"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	image, err := h.staffUsecase.SetMainImage(c.Request().Context(), staffID.String(), req.ImageID)
	if err != nil {
		c.Logger().Error(err)
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.JSON(http.StatusOK, image)
}

// UpdateMyImageCropPosition は自分の画像のクロップ位置を更新する。
func (h *StaffPortalHandler) UpdateMyImageCropPosition(c echo.Context) error {
	staffID, err := getStaffIDFromContext(c)
	if err != nil {
		return err
	}

	imageID := c.Param("imageId")
	if imageID == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "image id is required"})
	}

	// 所有権チェック
	if err := h.verifyImageOwnership(c, staffID, imageID); err != nil {
		return c.JSON(http.StatusForbidden, map[string]string{"error": "forbidden"})
	}

	var req struct {
		CropPosition string `json:"cropPosition"`
	}
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	image, err := h.staffUsecase.UpdateImageCropPosition(c.Request().Context(), imageID, req.CropPosition)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "internal server error"})
	}
	return c.JSON(http.StatusOK, image)
}
