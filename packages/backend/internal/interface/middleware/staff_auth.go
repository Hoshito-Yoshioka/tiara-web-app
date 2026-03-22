package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

// StaffJWTAuth はスタッフ用 JWT トークンを検証する Echo ミドルウェア。
// 管理者用 JWTAuth() と同じ署名キーを使用するが、
// クレームに staff_id / staff_username を格納する点が異なる。
func StaffJWTAuth() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			authHeader := c.Request().Header.Get("Authorization")
			if authHeader == "" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "authorization header is required"})
			}

			parts := strings.SplitN(authHeader, " ", 2)
			if len(parts) != 2 || parts[0] != "Bearer" {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid authorization header format"})
			}

			tokenString := parts[1]

			jwtSecret := os.Getenv("JWT_SECRET")
			if jwtSecret == "" {
				jwtSecret = "tiara-dev-secret-key"
			}

			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, echo.NewHTTPError(http.StatusUnauthorized, "unexpected signing method")
				}
				return []byte(jwtSecret), nil
			})

			if err != nil || !token.Valid {
				return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid or expired token"})
			}

			// スタッフ用クレームをコンテキストに格納
			if claims, ok := token.Claims.(jwt.MapClaims); ok {
				// トークン種別を検証（スタッフトークンであることを確認）
				tokenType, _ := claims["type"].(string)
				if tokenType != "staff" {
					return c.JSON(http.StatusUnauthorized, map[string]string{"error": "invalid token type"})
				}
				c.Set("staff_id", claims["sub"])
				c.Set("staff_username", claims["username"])
			}

			return next(c)
		}
	}
}
