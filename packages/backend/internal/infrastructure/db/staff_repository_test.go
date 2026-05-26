package db_test

import (
	"context"
	"testing"

	"tiara-web-app/backend/internal/infrastructure/db"
	"tiara-web-app/backend/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// ==========================================================
// StaffRepository 統合テスト
// - テスト用DB（docker-compose.test.yml）に対して実行
// - トランザクションロールバックで各テストを分離
// - Fixture ヘルパーでテストデータを投入
// ==========================================================

func TestStaffRepository_ListStaffs(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	q, _ := testutil.SetupTestTx(t, pool)
	// StaffRepository は pool を必要とするが、ListStaffs は pool を使わないため nil で可
	repo := db.NewStaffRepository(q, nil)

	// Fixture: Shop → Staff の順で挿入（FK制約）
	shop := testutil.InsertShop(t, q)
	testutil.InsertStaff(t, q, shop.ID)

	ctx := context.Background()
	staffs, err := repo.ListStaffs(ctx)

	require.NoError(t, err)
	assert.Len(t, staffs, 1)
	assert.Equal(t, "Test Staff", staffs[0].Name)
	assert.Equal(t, "キャスト", staffs[0].Role)
}

func TestStaffRepository_ListStaffsPaginated(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	q, _ := testutil.SetupTestTx(t, pool)
	repo := db.NewStaffRepository(q, nil)

	// Fixture: 3名のスタッフを挿入
	shop := testutil.InsertShop(t, q)
	testutil.InsertStaffWithOrder(t, q, shop.ID, "Staff A", 1)
	testutil.InsertStaffWithOrder(t, q, shop.ID, "Staff B", 2)
	testutil.InsertStaffWithOrder(t, q, shop.ID, "Staff C", 3)

	ctx := context.Background()

	t.Run("1ページ目 (limit=2)", func(t *testing.T) {
		staffs, err := repo.ListStaffsPaginated(ctx, 2, 0)
		require.NoError(t, err)
		assert.Len(t, staffs, 2)
		assert.Equal(t, "Staff A", staffs[0].Name)
		assert.Equal(t, "Staff B", staffs[1].Name)
	})

	t.Run("2ページ目 (limit=2, offset=2)", func(t *testing.T) {
		staffs, err := repo.ListStaffsPaginated(ctx, 2, 2)
		require.NoError(t, err)
		assert.Len(t, staffs, 1)
		assert.Equal(t, "Staff C", staffs[0].Name)
	})
}

func TestStaffRepository_CountStaffs(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	q, _ := testutil.SetupTestTx(t, pool)
	repo := db.NewStaffRepository(q, nil)

	shop := testutil.InsertShop(t, q)
	testutil.InsertStaffWithOrder(t, q, shop.ID, "Staff 1", 1)
	testutil.InsertStaffWithOrder(t, q, shop.ID, "Staff 2", 2)

	ctx := context.Background()
	count, err := repo.CountStaffs(ctx)

	require.NoError(t, err)
	assert.Equal(t, 2, count)
}

func TestStaffRepository_GetStaffByID(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	q, _ := testutil.SetupTestTx(t, pool)
	repo := db.NewStaffRepository(q, nil)

	shop := testutil.InsertShop(t, q)
	inserted := testutil.InsertStaff(t, q, shop.ID)
	staffID := uuid.UUID(inserted.ID.Bytes).String()

	ctx := context.Background()

	t.Run("正常系: 存在するスタッフを取得", func(t *testing.T) {
		staff, err := repo.GetStaffByID(ctx, staffID)

		require.NoError(t, err)
		assert.Equal(t, "Test Staff", staff.Name)
		assert.Equal(t, "キャスト", staff.Role)
	})

	t.Run("異常系: 存在しない ID → エラー", func(t *testing.T) {
		_, err := repo.GetStaffByID(ctx, "00000000-0000-0000-0000-000000000000")

		assert.Error(t, err)
	})
}

func TestStaffRepository_ListSchedulesByStaffID(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	q, _ := testutil.SetupTestTx(t, pool)
	repo := db.NewStaffRepository(q, nil)

	shop := testutil.InsertShop(t, q)
	staff := testutil.InsertStaff(t, q, shop.ID)
	// 月曜・水曜のスケジュールを追加
	testutil.InsertSchedule(t, q, staff.ID, 1) // 月曜
	testutil.InsertSchedule(t, q, staff.ID, 3) // 水曜

	staffID := uuid.UUID(staff.ID.Bytes).String()
	ctx := context.Background()

	schedules, err := repo.ListSchedulesByStaffID(ctx, staffID)

	require.NoError(t, err)
	assert.Len(t, schedules, 2)
}

// TestStaffRepository_Isolation はテスト分離を検証する。
func TestStaffRepository_Isolation(t *testing.T) {
	pool := testutil.SetupTestDB(t)
	q, _ := testutil.SetupTestTx(t, pool)
	repo := db.NewStaffRepository(q, nil)

	ctx := context.Background()

	staffs, err := repo.ListStaffs(ctx)
	require.NoError(t, err)
	assert.Empty(t, staffs, "トランザクションロールバックにより他テストのデータは見えないはず")

	count, err := repo.CountStaffs(ctx)
	require.NoError(t, err)
	assert.Equal(t, 0, count)
}
