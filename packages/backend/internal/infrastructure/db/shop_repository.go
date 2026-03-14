package db

import (
	"context"
	"tiara-web-app/backend/internal/domain"
	"tiara-web-app/backend/internal/usecase"
	"time" // timeパッケージを再度インポート

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgtype" // pgtypeパッケージをインポート
)

// Helper function to convert pgtype.Time to time.Time
// This assumes the date part is irrelevant for TIME type and defaults to 0001-01-01 UTC
func convertPgtypeTimeToTime(pgt pgtype.Time) time.Time {
	if !pgt.Valid {
		return time.Time{} // Return zero time if not valid
	}
	// Microseconds since midnight
	totalMicroseconds := pgt.Microseconds
	hour := totalMicroseconds / (3600 * 1_000_000)
	totalMicroseconds %= (3600 * 1_000_000)
	minute := totalMicroseconds / (60 * 1_000_000)
	totalMicroseconds %= (60 * 1_000_000)
	second := totalMicroseconds / 1_000_000
	microsecond := totalMicroseconds % 1_000_000

	return time.Date(0, 1, 1, int(hour), int(minute), int(second), int(microsecond*1000), time.UTC)
}

// parseTimeToPgtype は "HH:MM" or "HH:MM:SS" 形式の時刻文字列を pgtype.Time に変換する。
// admin API でフロントエンドから受け取った時刻文字列をDBに保存する際に使用。
func parseTimeToPgtype(timeStr string) (pgtype.Time, error) {
	t, err := time.Parse("15:04", timeStr)
	if err != nil {
		t, err = time.Parse("15:04:05", timeStr)
		if err != nil {
			return pgtype.Time{}, err
		}
	}
	microseconds := int64(t.Hour())*3600*1_000_000 + int64(t.Minute())*60*1_000_000 + int64(t.Second())*1_000_000
	return pgtype.Time{
		Microseconds: microseconds,
		Valid:        true,
	}, nil
}

// shopRepository は usecase.ShopRepository インターフェースの実装です。
type shopRepository struct {
	q *Queries
}

// NewShopRepository は新しいShopRepositoryのインスタンスを作成します。
func NewShopRepository(q *Queries) usecase.ShopRepository {
	return &shopRepository{
		q: q,
	}
}

// ListShops はすべての店舗をデータベースから取得します。
func (r *shopRepository) ListShops(ctx context.Context) ([]domain.Shop, error) {
	rows, err := r.q.ListShops(ctx)
	if err != nil {
		return nil, err
	}

	shops := make([]domain.Shop, len(rows))
	for i, row := range rows {
		shops[i] = domain.Shop{
			ID:          uuid.UUID(row.ID.Bytes), // pgtype.UUID to uuid.UUID
			Name:        row.Name,
			Address:     row.Address,
			OpeningTime: convertPgtypeTimeToTime(row.OpeningTime), // pgtype.Time to time.Time
			ClosingTime: convertPgtypeTimeToTime(row.ClosingTime), // pgtype.Time to time.Time
			CreatedAt:   row.CreatedAt.Time,                       // pgtype.Timestamptz to time.Time
			UpdatedAt:   row.UpdatedAt.Time,                       // pgtype.Timestamptz to time.Time
		}
	}
	return shops, nil
}

// GetShopByID は指定されたIDの店舗をデータベースから取得します。
func (r *shopRepository) GetShopByID(ctx context.Context, id string) (domain.Shop, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.Shop{}, err
	}

	pgtypeUUID := pgtype.UUID{
		Bytes: uid,
		Valid: true,
	}

	row, err := r.q.GetShopByID(ctx, pgtypeUUID)
	if err != nil {
		return domain.Shop{}, err
	}

	shop := domain.Shop{
		ID:          uuid.UUID(row.ID.Bytes),
		Name:        row.Name,
		Address:     row.Address,
		OpeningTime: convertPgtypeTimeToTime(row.OpeningTime),
		ClosingTime: convertPgtypeTimeToTime(row.ClosingTime),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}
	return shop, nil
}

// UpdateShop は指定されたIDの店舗情報を更新する。
func (r *shopRepository) UpdateShop(ctx context.Context, id string, input domain.UpdateShopInput) (domain.Shop, error) {
	uid, err := uuid.Parse(id)
	if err != nil {
		return domain.Shop{}, err
	}

	openingTime, err := parseTimeToPgtype(input.OpeningTime)
	if err != nil {
		return domain.Shop{}, err
	}

	closingTime, err := parseTimeToPgtype(input.ClosingTime)
	if err != nil {
		return domain.Shop{}, err
	}

	row, err := r.q.UpdateShop(ctx, UpdateShopParams{
		ID:          pgtype.UUID{Bytes: uid, Valid: true},
		Name:        input.Name,
		Address:     input.Address,
		OpeningTime: openingTime,
		ClosingTime: closingTime,
	})
	if err != nil {
		return domain.Shop{}, err
	}

	return domain.Shop{
		ID:          uuid.UUID(row.ID.Bytes),
		Name:        row.Name,
		Address:     row.Address,
		OpeningTime: convertPgtypeTimeToTime(row.OpeningTime),
		ClosingTime: convertPgtypeTimeToTime(row.ClosingTime),
		CreatedAt:   row.CreatedAt.Time,
		UpdatedAt:   row.UpdatedAt.Time,
	}, nil
}
