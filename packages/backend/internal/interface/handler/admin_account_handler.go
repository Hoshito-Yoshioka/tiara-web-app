package handler

import (
	"net/http"
	"time"

	"tiara-web-app/backend/internal/usecase"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

// AdminAccountHandler は管理者によるスタッフアカウント管理のHTTPハンドラー。
type AdminAccountHandler struct {
	accountUsecase usecase.AdminAccountUsecase
}

// NewAdminAccountHandler は新しいAdminAccountHandlerのインスタンスを作成する。
func NewAdminAccountHandler(us usecase.AdminAccountUsecase) *AdminAccountHandler {
	return &AdminAccountHandler{accountUsecase: us}
}

// StaffAccountResponse はスタッフアカウントのレスポンス型。
// パスワードハッシュは含めない。
type StaffAccountResponse struct {
	ID        string `json:"id"`
	StaffID   string `json:"staffId"`
	Username  string `json:"username"`
	CreatedAt string `json:"createdAt"`
	UpdatedAt string `json:"updatedAt"`
}

// ListStaffAccounts は全スタッフアカウント一覧を返す。
func (h *AdminAccountHandler) ListStaffAccounts(c echo.Context) error {
	accounts, err := h.accountUsecase.ListStaffAccounts(c.Request().Context())
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	resp := make([]StaffAccountResponse, len(accounts))
	for i, a := range accounts {
		resp[i] = StaffAccountResponse{
			ID:        a.ID.String(),
			StaffID:   a.StaffID.String(),
			Username:  a.Username,
			CreatedAt: a.CreatedAt.Format(time.RFC3339),
			UpdatedAt: a.UpdatedAt.Format(time.RFC3339),
		}
	}
	return c.JSON(http.StatusOK, resp)
}

// GetStaffAccountByStaffID は特定スタッフのアカウント情報を返す。
func (h *AdminAccountHandler) GetStaffAccountByStaffID(c echo.Context) error {
	staffIDStr := c.Param("staffId")
	staffID, err := uuid.Parse(staffIDStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid staff ID"})
	}

	account, err := h.accountUsecase.GetStaffAccountByStaffID(c.Request().Context(), staffID)
	if err != nil {
		// アカウントが存在しない場合は空レスポンスを返す（404ではなく）
		return c.JSON(http.StatusOK, nil)
	}

	return c.JSON(http.StatusOK, StaffAccountResponse{
		ID:        account.ID.String(),
		StaffID:   account.StaffID.String(),
		Username:  account.Username,
		CreatedAt: account.CreatedAt.Format(time.RFC3339),
		UpdatedAt: account.UpdatedAt.Format(time.RFC3339),
	})
}

// CreateStaffAccountRequest はスタッフアカウント作成リクエスト型。
type CreateStaffAccountRequest struct {
	StaffID  string `json:"staffId"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// CreateStaffAccount は新しいスタッフアカウントを作成する。
func (h *AdminAccountHandler) CreateStaffAccount(c echo.Context) error {
	var req CreateStaffAccountRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	staffID, err := uuid.Parse(req.StaffID)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid staff ID"})
	}

	account, err := h.accountUsecase.CreateStaffAccount(c.Request().Context(), staffID, req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusCreated, StaffAccountResponse{
		ID:        account.ID.String(),
		StaffID:   account.StaffID.String(),
		Username:  account.Username,
		CreatedAt: account.CreatedAt.Format(time.RFC3339),
		UpdatedAt: account.UpdatedAt.Format(time.RFC3339),
	})
}

// UpdateStaffAccountRequest はスタッフアカウント更新リクエスト型。
type UpdateStaffAccountRequest struct {
	Username string `json:"username"`
	Password string `json:"password"` // 空の場合はパスワード変更なし
}

// UpdateStaffAccount はスタッフアカウントを更新する。
func (h *AdminAccountHandler) UpdateStaffAccount(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid account ID"})
	}

	var req UpdateStaffAccountRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	account, err := h.accountUsecase.UpdateStaffAccount(c.Request().Context(), id, req.Username, req.Password)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, StaffAccountResponse{
		ID:        account.ID.String(),
		StaffID:   account.StaffID.String(),
		Username:  account.Username,
		CreatedAt: account.CreatedAt.Format(time.RFC3339),
		UpdatedAt: account.UpdatedAt.Format(time.RFC3339),
	})
}

// DeleteStaffAccount はスタッフアカウントを削除する。
func (h *AdminAccountHandler) DeleteStaffAccount(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid account ID"})
	}

	if err := h.accountUsecase.DeleteStaffAccount(c.Request().Context(), id); err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": err.Error()})
	}

	return c.NoContent(http.StatusNoContent)
}
