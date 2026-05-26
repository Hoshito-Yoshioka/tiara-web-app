package testutil

import (
	"context"
	"testing"

	"tiara-web-app/backend/internal/infrastructure/db"

	"github.com/jackc/pgx/v5/pgtype"
)

// InsertShop はテスト用の Shop レコードを DB に挿入する。
// トランザクション内で呼び出すことで、テスト終了時にロールバックで自動削除される。
func InsertShop(t *testing.T, q *db.Queries) db.Shop {
	t.Helper()

	shop, err := q.CreateShop(context.Background(), db.CreateShopParams{
		Name:    "Test Shop",
		Address: "Test Address",
		OpeningTime: pgtype.Time{
			Microseconds: 20 * 3600 * 1_000_000, // 20:00
			Valid:        true,
		},
		ClosingTime: pgtype.Time{
			Microseconds: 2 * 3600 * 1_000_000, // 02:00
			Valid:        true,
		},
	})
	if err != nil {
		t.Fatalf("Shop の挿入に失敗: %v", err)
	}
	return shop
}

// InsertStaff はテスト用の Staff レコードを DB に挿入する。
// shopID は事前に InsertShop で作成した Shop の ID を渡す。
func InsertStaff(t *testing.T, q *db.Queries, shopID pgtype.UUID) db.Staff {
	t.Helper()

	staff, err := q.CreateStaff(context.Background(), db.CreateStaffParams{
		ShopID:            shopID,
		Name:              "Test Staff",
		Role:              "キャスト",
		Bio:               "テスト用スタッフです。",
		ImageUrl:          "/uploads/staff/test.jpg",
		ImageCropPosition: "50 50",
		SortOrder:         1,
	})
	if err != nil {
		t.Fatalf("Staff の挿入に失敗: %v", err)
	}
	return staff
}

// InsertStaffWithOrder は指定した sortOrder で Staff レコードを DB に挿入する。
func InsertStaffWithOrder(t *testing.T, q *db.Queries, shopID pgtype.UUID, name string, sortOrder int32) db.Staff {
	t.Helper()

	staff, err := q.CreateStaff(context.Background(), db.CreateStaffParams{
		ShopID:            shopID,
		Name:              name,
		Role:              "キャスト",
		Bio:               "",
		ImageUrl:          "",
		ImageCropPosition: "50 50",
		SortOrder:         sortOrder,
	})
	if err != nil {
		t.Fatalf("Staff の挿入に失敗: %v", err)
	}
	return staff
}

// InsertSchedule はテスト用の StaffSchedule レコードを DB に挿入する。
func InsertSchedule(t *testing.T, q *db.Queries, staffID pgtype.UUID, dayOfWeek int32) db.StaffSchedule {
	t.Helper()

	schedule, err := q.CreateSchedule(context.Background(), db.CreateScheduleParams{
		StaffID:   staffID,
		DayOfWeek: dayOfWeek,
		StartTime: pgtype.Time{Microseconds: 20 * 3600 * 1_000_000, Valid: true},
		EndTime:   pgtype.Time{Microseconds: 2 * 3600 * 1_000_000, Valid: true},
	})
	if err != nil {
		t.Fatalf("Schedule の挿入に失敗: %v", err)
	}
	return schedule
}
