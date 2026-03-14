package usecase

import (
	"context"
	"tiara-web-app/backend/internal/domain"
)

// ShopRepository は店舗データの永続化を抽象化するインターフェースです。
// これはinfrastructure層によって実装されます。
type ShopRepository interface {
	ListShops(ctx context.Context) ([]domain.Shop, error)
	GetShopByID(ctx context.Context, id string) (domain.Shop, error)
	UpdateShop(ctx context.Context, id string, input domain.UpdateShopInput) (domain.Shop, error)
}

// ShopUsecase は店舗に関するビジネスロジックを定義するインターフェースです。
type ShopUsecase interface {
	ListShops(ctx context.Context) ([]domain.Shop, error)
	GetShopByID(ctx context.Context, id string) (domain.Shop, error)
	UpdateShop(ctx context.Context, id string, input domain.UpdateShopInput) (domain.Shop, error)
}

// shopUsecase はShopUsecaseインターフェースの実装です。
type shopUsecase struct {
	shopRepo ShopRepository
}

// NewShopUsecase は新しいShopUsecaseのインスタンスを作成します。
func NewShopUsecase(repo ShopRepository) ShopUsecase {
	return &shopUsecase{
		shopRepo: repo,
	}
}

// ListShops はすべての店舗を取得します。
func (u *shopUsecase) ListShops(ctx context.Context) ([]domain.Shop, error) {
	shops, err := u.shopRepo.ListShops(ctx)
	if err != nil {
		return nil, err
	}
	return shops, nil
}

// GetShopByID は指定されたIDの店舗を取得します。
func (u *shopUsecase) GetShopByID(ctx context.Context, id string) (domain.Shop, error) {
	shop, err := u.shopRepo.GetShopByID(ctx, id)
	if err != nil {
		return domain.Shop{}, err
	}
	return shop, nil
}

// UpdateShop は指定されたIDの店舗情報を更新します。
func (u *shopUsecase) UpdateShop(ctx context.Context, id string, input domain.UpdateShopInput) (domain.Shop, error) {
	return u.shopRepo.UpdateShop(ctx, id, input)
}
