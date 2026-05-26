package config

import (
	"strings"

	"github.com/spf13/viper"
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
	JWTSecret      string
	JWTExpiryHours int
	UploadDir      string
}

// Load は viper を使って環境変数・.env ファイルから Config を読み込む。
// 優先順位: 環境変数 > .env ファイル > デフォルト値
func Load() *Config {
	v := viper.New()

	// デフォルト値（ローカル開発用）
	v.SetDefault("PORT", "1323")
	v.SetDefault("CORS_ORIGIN", "http://localhost:3001")
	v.SetDefault("JWT_EXPIRY_HOURS", 2)
	v.SetDefault("UPLOAD_DIR", "uploads")

	// .env ファイルの自動読み込み（source .env 不要）
	v.SetConfigName(".env")
	v.SetConfigType("env")
	v.AddConfigPath(".")      // カレントディレクトリ
	v.AddConfigPath("../../") // monorepo ルート（packages/backend から実行時）
	_ = v.ReadInConfig()      // .env が無くても環境変数で動作

	// 環境変数の自動バインド（OS 環境変数が .env より優先）
	v.AutomaticEnv()

	// CORS_ORIGIN はカンマ区切りで複数オリジンを指定可能。
	// 例: "https://tiara.example.com,https://preview.tiara.example.com"
	corsOrigins := parseCORSOrigins(v.GetString("CORS_ORIGIN"))

	return &Config{
		DatabaseURL:    v.GetString("DATABASE_URL"),
		Port:           v.GetString("PORT"),
		CORSOrigins:    corsOrigins,
		JWTSecret:      v.GetString("JWT_SECRET"),
		JWTExpiryHours: v.GetInt("JWT_EXPIRY_HOURS"),
		UploadDir:      v.GetString("UPLOAD_DIR"),
	}
}

// parseCORSOrigins はカンマ区切りの CORS_ORIGIN 文字列をスライスに変換する。
func parseCORSOrigins(raw string) []string {
	var origins []string
	for _, o := range strings.Split(raw, ",") {
		if trimmed := strings.TrimSpace(o); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	if len(origins) == 0 {
		return []string{"http://localhost:3001"}
	}
	return origins
}
