package db_test

import (
	"context"
	"testing"

	"tiara-web-app/backend/internal/domain"
	"tiara-web-app/backend/internal/infrastructure/db"
	"tiara-web-app/backend/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==========================================================
// ShopRepository 統合テスト
// - テスト用DB（docker-compose.test.yml）に対して実行
// - トランザクションロールバックで各テストを分離
// - Fixture ヘルパーでテストデータを投入
// ==========================================================

func TestShopRepository_ListShops(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	q, _ := testutil.SetupTestTx(t, pool)
	repo := db.NewShopRepository(q)

	// Fixture でテストデータを挿入
	testutil.InsertShop(t, q)

	ctx := context.Background()
	shops, err := repo.ListShops(ctx)

	require.NoError(t, err)
	assert.Len(t, shops, 1)
	assert.Equal(t, "Test Shop", shops[0].Name)
	assert.Equal(t, "Test Address", shops[0].Address)
}

func TestShopRepository_GetShopByID(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	q, _ := testutil.SetupTestTx(t, pool)
	repo := db.NewShopRepository(q)

	// Fixture でテストデータを挿入
	inserted := testutil.InsertShop(t, q)
	shopID := uuid.UUID(inserted.ID.Bytes).String()

	ctx := context.Background()

	t.Run("正常系: 存在する Shop を取得", func(t *testing.T) {
		shop, err := repo.GetShopByID(ctx, shopID)

		require.NoError(t, err)
		assert.Equal(t, "Test Shop", shop.Name)
		assert.Equal(t, "Test Address", shop.Address)
	})

	t.Run("異常系: 存在しない ID → エラー", func(t *testing.T) {
		_, err := repo.GetShopByID(ctx, "00000000-0000-0000-0000-000000000000")

		assert.Error(t, err)
	})
}

func TestShopRepository_UpdateShop(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	q, _ := testutil.SetupTestTx(t, pool)
	repo := db.NewShopRepository(q)

	// Fixture でテストデータを挿入
	inserted := testutil.InsertShop(t, q)
	shopID := uuid.UUID(inserted.ID.Bytes).String()

	ctx := context.Background()

	updated, err := repo.UpdateShop(ctx, shopID, domain.UpdateShopInput{
		Name:        "Updated Shop",
		Address:     "Updated Address",
		OpeningTime: "19:00",
		ClosingTime: "03:00",
	})

	require.NoError(t, err)
	assert.Equal(t, "Updated Shop", updated.Name)
	assert.Equal(t, "Updated Address", updated.Address)
}

// TestShopRepository_Isolation はトランザクションロールバックによるテスト分離を検証する。
// 他のテストで挿入されたデータがここでは見えないことを確認する。
func TestShopRepository_Isolation(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	q, _ := testutil.SetupTestTx(t, pool)
	repo := db.NewShopRepository(q)

	ctx := context.Background()

	// データを挿入していないので、一覧は空であるべき
	shops, err := repo.ListShops(ctx)
	require.NoError(t, err)
	assert.Empty(t, shops, "トランザクションロールバックにより他テストのデータは見えないはず")
}
