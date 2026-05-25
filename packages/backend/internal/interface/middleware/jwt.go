package middleware

import (
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// TokenClaims はパース済み JWT トークンのクレームを保持する。
type TokenClaims struct {
	Subject  string
	Username string
	Type     string // "admin" or "staff"
}

// parseToken は Authorization ヘッダーから JWT トークンをパースし、クレームを返す。
// jwtSecret は外部（Config）から注入される。
func parseToken(c echo.Context, jwtSecret string) (*TokenClaims, error) {
	authHeader := c.Request().Header.Get("Authorization")
	if authHeader == "" {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "authorization header is required")
	}

	parts := strings.SplitN(authHeader, " ", 2)
	if len(parts) != 2 || parts[0] != "Bearer" {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid authorization header format")
	}

	token, err := jwt.Parse(parts[1], func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, echo.NewHTTPError(http.StatusUnauthorized, "unexpected signing method")
		}
		return []byte(jwtSecret), nil
	})

	if err != nil || !token.Valid {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid or expired token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "invalid token claims")
	}

	sub, ok := claims["sub"].(string)
	if !ok || sub == "" {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid sub claim")
	}
	tokenType, ok := claims["type"].(string)
	if !ok || tokenType == "" {
		return nil, echo.NewHTTPError(http.StatusUnauthorized, "missing or invalid type claim")
	}
	username, _ := claims["username"].(string)

	return &TokenClaims{
		Subject:  sub,
		Username: username,
		Type:     tokenType,
	}, nil
}

// JWTAuth は管理者用 JWT トークンを検証する Echo ミドルウェア。
// jwtSecret を外部から受け取ることで os.Getenv への依存を排除する。
func JWTAuth(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tc, err := parseToken(c, jwtSecret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
			}

			// スタッフトークンによる管理者APIへのアクセスを拒否
			if tc.Type != "admin" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token type"})
			}

			c.Set("admin_id", tc.Subject)
			c.Set("admin_username", tc.Username)
			return next(c)
		}
	}
}

// StaffJWTAuth はスタッフ用 JWT トークンを検証する Echo ミドルウェア。
// 管理者トークンと区別するため type="staff" クレームを検証する。
func StaffJWTAuth(jwtSecret string) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			tc, err := parseToken(c, jwtSecret)
			if err != nil {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": err.Error()})
			}

			if tc.Type != "staff" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token type"})
			}

			c.Set("staff_id", tc.Subject)
			c.Set("staff_username", tc.Username)
			return next(c)
		}
	}
}
