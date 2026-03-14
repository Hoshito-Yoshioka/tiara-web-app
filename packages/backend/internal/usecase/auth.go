package usecase

import (
	"context"
	"errors"
	"tiara-web-app/backend/internal/domain"

	"golang.org/x/crypto/bcrypt"
)

// AdminRepository は管理者データの永続化を抽象化するインターフェース。
// クリーンアーキテクチャの原則に従い、usecase層にインターフェースを定義し、
// infrastructure層で実装する（DIP: 依存性逆転の原則）。
type AdminRepository interface {
	GetAdminByUsername(ctx context.Context, username string) (domain.Admin, error)
}

// AuthUsecase は認証に関するビジネスロジックを定義するインターフェース。
type AuthUsecase interface {
	Login(ctx context.Context, username, password string) (domain.Admin, error)
}

type authUsecase struct {
	adminRepo AdminRepository
}

// NewAuthUsecase は新しいAuthUsecaseのインスタンスを作成する。
func NewAuthUsecase(repo AdminRepository) AuthUsecase {
	return &authUsecase{adminRepo: repo}
}

// Login はユーザー名とパスワードで管理者を認証する。
// bcrypt を使用してパスワードハッシュを検証する。
func (u *authUsecase) Login(ctx context.Context, username, password string) (domain.Admin, error) {
	admin, err := u.adminRepo.GetAdminByUsername(ctx, username)
	if err != nil {
		return domain.Admin{}, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(admin.PasswordHash), []byte(password))
	if err != nil {
		return domain.Admin{}, errors.New("invalid credentials")
	}

	return admin, nil
}
