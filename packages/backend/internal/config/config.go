package config

import (
	"os"
	"strconv"
)

// Config はアプリケーション全体の設定を保持する構造体。
// 環境変数から一度だけ読み込み、各レイヤーに注入する。
type Config struct {
	// Database
	DatabaseURL string

	// Server
	Port        string
	CORSOrigins []string

	// JWT
	JWTSecret         string
	JWTExpiryHours    int
	UploadDir         string
}

// Load は環境変数からConfigを読み込む。
// デフォルト値はローカル開発用。
// JWT_SECRET は必須（デフォルト値なし）。
func Load() *Config {
	return &Config{
		DatabaseURL:    getEnv("DATABASE_URL", ""),
		Port:           getEnv("PORT", "1323"),
		CORSOrigins:    []string{getEnv("CORS_ORIGIN", "http://localhost:3001")},
		JWTSecret:      getEnv("JWT_SECRET", ""),
		JWTExpiryHours: getEnvInt("JWT_EXPIRY_HOURS", 2),
		UploadDir:      getEnv("UPLOAD_DIR", "uploads"),
	}
}

// getEnv は環境変数を取得し、未設定ならデフォルト値を返す。
func getEnv(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}

// getEnvInt は環境変数を整数で取得する。
func getEnvInt(key string, defaultVal int) int {
	if v := os.Getenv(key); v != "" {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultVal
}
