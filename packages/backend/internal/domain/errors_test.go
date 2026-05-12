package domain

import (
	"errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDomainErrors_Is(t *testing.T) {
	// WrapNotFound でラップしたエラーが ErrNotFound として判定できることを検証
	wrapped := WrapNotFound("shop not found")
	assert.True(t, errors.Is(wrapped, ErrNotFound))
	assert.False(t, errors.Is(wrapped, ErrUnauthorized))
	assert.Contains(t, wrapped.Error(), "shop not found")

	// WrapInvalidInput でラップしたエラーが ErrInvalidInput として判定できることを検証
	wrapped2 := WrapInvalidInput("name is required")
	assert.True(t, errors.Is(wrapped2, ErrInvalidInput))
	assert.False(t, errors.Is(wrapped2, ErrNotFound))
}

func TestDomainErrors_FmtErrorf(t *testing.T) {
	// fmt.Errorf で %w を使った場合も errors.Is で正しく判定できることを検証
	err := fmt.Errorf("invalid credentials: %w", ErrUnauthorized)
	assert.True(t, errors.Is(err, ErrUnauthorized))
	assert.Contains(t, err.Error(), "invalid credentials")
}
