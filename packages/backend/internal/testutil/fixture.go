package testutil

import (
	"tiara-web-app/backend/internal/domain"

	"github.com/google/uuid"
)

// NewStaff はデフォルト値を持つ Staff テストデータを生成する。
func NewStaff() domain.Staff {
	return domain.Staff{
		ID:   uuid.New(),
		Name: "Test Staff",
		Role: "キャスト",
	}
}

// NewShop はデフォルト値を持つ Shop テストデータを生成する。
func NewShop() domain.Shop {
	return domain.Shop{
		ID:      uuid.New(),
		Name:    "Tiara",
		Address: "Tokyo",
	}
}

// NewMenuCategory はデフォルト値を持つ MenuCategory テストデータを生成する。
func NewMenuCategory() domain.MenuCategory {
	return domain.MenuCategory{
		ID:        uuid.New(),
		Name:      "Cocktails",
		SortOrder: 1,
	}
}

// NewMenuItem はデフォルト値を持つ MenuItem テストデータを生成する。
func NewMenuItem() domain.MenuItem {
	return domain.MenuItem{
		ID:    uuid.New(),
		Name:  "Mojito",
		Price: "800",
	}
}

// NewStaffAccount はデフォルト値を持つ StaffAccount テストデータを生成する。
func NewStaffAccount() domain.StaffAccount {
	return domain.StaffAccount{
		ID:       uuid.New(),
		StaffID:  uuid.New(),
		Username: "testuser",
	}
}

// NewStaffProfileDraft はデフォルト値を持つ StaffProfileDraft テストデータを生成する。
func NewStaffProfileDraft() domain.StaffProfileDraft {
	return domain.StaffProfileDraft{
		ID:      uuid.New(),
		StaffID: uuid.New(),
		Name:    "Test Staff",
		Status:  domain.DraftStatusDraft,
	}
}

// NewStaffScheduleDraft はデフォルト値を持つ StaffScheduleDraft テストデータを生成する。
func NewStaffScheduleDraft() domain.StaffScheduleDraft {
	return domain.StaffScheduleDraft{
		ID:      uuid.New(),
		StaffID: uuid.New(),
		Status:  domain.DraftStatusDraft,
	}
}

// NewStaffImage はデフォルト値を持つ StaffImage テストデータを生成する。
func NewStaffImage() domain.StaffImage {
	return domain.StaffImage{
		ID:       uuid.New(),
		StaffID:  uuid.New(),
		ImageURL: "/uploads/staff/test.jpg",
		IsMain:   true,
	}
}

// NewAdmin はデフォルト値を持つ Admin テストデータを生成する。
func NewAdmin() domain.Admin {
	return domain.Admin{
		ID:       uuid.New(),
		Username: "admin",
	}
}
