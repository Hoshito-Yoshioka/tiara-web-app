package usecase

import (
	"context"
	"errors"
	"testing"
	"time"

	"tiara-web-app/backend/internal/domain"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/crypto/bcrypt"
)

// mockAdminRepository は AdminRepository のテスト用モック。
type mockAdminRepository struct {
	admin domain.Admin
	err   error
}

func (m *mockAdminRepository) GetAdminByUsername(_ context.Context, _ string) (domain.Admin, error) {
	return m.admin, m.err
}

func TestAuthUsecase_Login_Success(t *testing.T) {
	password := "test-password"
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	require.NoError(t, err)

	admin := domain.Admin{
		ID:           uuid.New(),
		Username:     "admin",
		PasswordHash: string(hash),
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	repo := &mockAdminRepository{admin: admin}
	uc := NewAuthUsecase(repo)

	result, err := uc.Login(context.Background(), "admin", password)

	assert.NoError(t, err)
	assert.Equal(t, admin.ID, result.ID)
	assert.Equal(t, admin.Username, result.Username)
}

func TestAuthUsecase_Login_UserNotFound(t *testing.T) {
	repo := &mockAdminRepository{err: errors.New("not found")}
	uc := NewAuthUsecase(repo)

	_, err := uc.Login(context.Background(), "unknown", "password")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrUnauthorized))
}

func TestAuthUsecase_Login_WrongPassword(t *testing.T) {
	hash, err := bcrypt.GenerateFromPassword([]byte("correct-password"), bcrypt.DefaultCost)
	require.NoError(t, err)

	admin := domain.Admin{
		ID:           uuid.New(),
		Username:     "admin",
		PasswordHash: string(hash),
	}

	repo := &mockAdminRepository{admin: admin}
	uc := NewAuthUsecase(repo)

	_, err = uc.Login(context.Background(), "admin", "wrong-password")

	assert.Error(t, err)
	assert.True(t, errors.Is(err, domain.ErrUnauthorized))
}
