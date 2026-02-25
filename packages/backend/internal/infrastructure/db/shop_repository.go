package db

import (
	"context"
	"time" // timeパッケージを再度インポート
	"tiara-web-app/backend/internal/domain"
	"tiara-web-app/backend/internal/usecase"

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
			ID:          uuid.UUID(row.ID.Bytes),                               // pgtype.UUID to uuid.UUID
			Name:        row.Name,
			Address:     row.Address,
			OpeningTime: convertPgtypeTimeToTime(row.OpeningTime), // pgtype.Time to time.Time
			ClosingTime: convertPgtypeTimeToTime(row.ClosingTime), // pgtype.Time to time.Time
			CreatedAt:   row.CreatedAt.Time,                              // pgtype.Timestamptz to time.Time
			UpdatedAt:   row.UpdatedAt.Time,                              // pgtype.Timestamptz to time.Time
		}
	}
	return shops, nil
}
