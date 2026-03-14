package domain

import (
	"time"

	"github.com/google/uuid"
)

// MenuCategory はメニューカテゴリ（例: Cocktails, Whisky & Spirits）を表す。
type MenuCategory struct {
	ID          uuid.UUID
	Name        string
	Description string
	SortOrder   int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// MenuItem はメニュー品目（名前・価格・説明）を表す。
type MenuItem struct {
	ID          uuid.UUID
	CategoryID  uuid.UUID
	Name        string
	Price       string
	Description string
	SortOrder   int32
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// MenuCategoryWithItems はカテゴリとそれに属するアイテムをまとめた集約型。
type MenuCategoryWithItems struct {
	Category MenuCategory
	Items    []MenuItem
}

// CreateMenuCategoryInput はカテゴリ新規作成時の入力型。
type CreateMenuCategoryInput struct {
	Name        string
	Description string
	SortOrder   int32
}

// UpdateMenuCategoryInput はカテゴリ更新時の入力型。
type UpdateMenuCategoryInput struct {
	Name        string
	Description string
	SortOrder   int32
}

// CreateMenuItemInput はメニュー品目新規作成時の入力型。
type CreateMenuItemInput struct {
	CategoryID  string
	Name        string
	Price       string
	Description string
	SortOrder   int32
}

// UpdateMenuItemInput はメニュー品目更新時の入力型。
type UpdateMenuItemInput struct {
	Name        string
	Price       string
	Description string
	SortOrder   int32
}
