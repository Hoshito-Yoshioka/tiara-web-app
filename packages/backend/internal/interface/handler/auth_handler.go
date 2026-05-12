package handler

import (
	"net/http"
	"time"

	"tiara-web-app/backend/internal/usecase"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// AuthHandler は認証関連のHTTPハンドラー。
type AuthHandler struct {
	authUsecase    usecase.AuthUsecase
	jwtSecret      string
	jwtExpiryHours int
}

// NewAuthHandler は新しいAuthHandlerのインスタンスを作成する。
func NewAuthHandler(us usecase.AuthUsecase, jwtSecret string, jwtExpiryHours int) *AuthHandler {
	return &AuthHandler{authUsecase: us, jwtSecret: jwtSecret, jwtExpiryHours: jwtExpiryHours}
}

// LoginRequest はログインリクエストのボディ型。
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// LoginResponse はログイン成功時のレスポンス型。
type LoginResponse struct {
	Token string `json:"token"`
}

// Login はユーザー名とパスワードで認証し、JWTトークンを返す。
func (h *AuthHandler) Login(c echo.Context) error {
	var req LoginRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "invalid request body"})
	}

	if req.Username == "" || req.Password == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{"error": "username and password are required"})
	}

	admin, err := h.authUsecase.Login(c.Request().Context(), req.Username, req.Password)
	if err != nil {
		return handleError(c, err)
	}

	// JWT トークンを生成
	claims := jwt.MapClaims{
		"sub":      admin.ID.String(),
		"username": admin.Username,
		"type":     "admin", // スタッフトークンと区別するためのクレーム
		"exp":      time.Now().Add(time.Duration(h.jwtExpiryHours) * time.Hour).Unix(),
		"iat":      time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(h.jwtSecret))
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{"error": "failed to generate token"})
	}

	return c.JSON(http.StatusOK, LoginResponse{Token: tokenString})
}

// Verify はJWTトークンの有効性を確認する。
// JWTAuth ミドルウェアを通過できた時点でトークンは有効。
func (h *AuthHandler) Verify(c echo.Context) error {
	return c.JSON(http.StatusOK, map[string]string{"status": "ok"})
}
