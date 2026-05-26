package testutil

import (
	"context"
	"os"
	"strings"
	"testing"

	"tiara-web-app/backend/internal/infrastructure/db"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

// defaultTestDSN はテスト用DBのデフォルト接続文字列。
// docker-compose.test.yml で起動した PostgreSQL に接続する。
const defaultTestDSN = "postgres://testuser:testpassword@localhost:5433/tiara_test?sslmode=disable"

// SetupTestDB はテスト用 DB プールを作成し、スキーマを適用する。
// テスト終了時にプールを自動クローズする。
// 環境変数 TEST_DATABASE_URL が設定されていればそちらを使用する。
func SetupTestDB(t *testing.T) *pgxpool.Pool {
	t.Helper()

	dsn := os.Getenv("TEST_DATABASE_URL")
	if dsn == "" {
		dsn = defaultTestDSN
	}

	ctx := context.Background()
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		t.Skipf("テスト用DBに接続できません（docker-compose.test.yml を起動してください）: %v", err)
	}

	if err := pool.Ping(ctx); err != nil {
		pool.Close()
		t.Skipf("テスト用DBに接続できません: %v", err)
	}

	// スキーマを適用（毎回クリーンな状態にする）
	applySchema(t, pool)

	t.Cleanup(func() {
		pool.Close()
	})

	return pool
}

// SetupTestTx はテスト用トランザクションを開始し、sqlc の Queries を返す。
// テスト終了時にトランザクションをロールバックして自動クリーンアップする。
// これにより各テストが完全に分離され、テスト順序に依存しない。
func SetupTestTx(t *testing.T, pool *pgxpool.Pool) (*db.Queries, pgx.Tx) {
	t.Helper()

	ctx := context.Background()
	tx, err := pool.Begin(ctx)
	if err != nil {
		t.Fatalf("トランザクション開始に失敗: %v", err)
	}

	t.Cleanup(func() {
		// テスト終了時に必ずロールバック（コミットしない）
		_ = tx.Rollback(ctx)
	})

	q := db.New(tx)
	return q, tx
}

// applySchema はテスト用DBにスキーマを適用する。
// IF NOT EXISTS 句により、既存テーブルがある場合は何もしない。
func applySchema(t *testing.T, pool *pgxpool.Pool) {
	t.Helper()

	// Go テストの作業ディレクトリはパッケージディレクトリなので、
	// プロジェクトルート（go.mod のあるディレクトリ）からの相対パスを探索する。
	candidates := []string{
		"migrations/schema.sql",
		"../migrations/schema.sql",
		"../../migrations/schema.sql",
		"../../../migrations/schema.sql",
		"../../../../migrations/schema.sql",
	}

	var schema []byte
	var err error
	for _, path := range candidates {
		schema, err = os.ReadFile(path)
		if err == nil {
			break
		}
	}
	if err != nil {
		t.Fatalf("schema.sql の読み込みに失敗（go.mod のあるディレクトリに migrations/schema.sql が必要です）: %v", err)
	}

	ctx := context.Background()
	// PostgreSQL は CREATE TRIGGER IF NOT EXISTS をサポートしないため、
	// CREATE OR REPLACE TRIGGER に変換する（PG14+）
	schemaSQL := strings.ReplaceAll(string(schema), "CREATE TRIGGER IF NOT EXISTS", "CREATE OR REPLACE TRIGGER")
	if _, err := pool.Exec(ctx, schemaSQL); err != nil {
		t.Fatalf("スキーマ適用に失敗: %v", err)
	}
}
