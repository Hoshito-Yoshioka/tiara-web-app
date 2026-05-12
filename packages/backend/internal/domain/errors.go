package domain

import (
	"errors"
	"fmt"
)

// ドメイン層で共通的に使用するセンチネルエラー。
// usecase 層でこれらのエラーを返し、handler 層で errors.Is() により
// 適切な HTTP ステータスコードへマッピングする。
var (
	// ErrNotFound はリソースが見つからない場合のエラー。
	ErrNotFound = errors.New("resource not found")

	// ErrUnauthorized は認証に失敗した場合のエラー。
	ErrUnauthorized = errors.New("unauthorized")

	// ErrForbidden は権限が不足している場合のエラー。
	ErrForbidden = errors.New("forbidden")

	// ErrConflict はリソースの競合が発生した場合のエラー。
	ErrConflict = errors.New("resource conflict")

	// ErrInvalidInput は入力値が不正な場合のエラー。
	ErrInvalidInput = errors.New("invalid input")

	// ErrInternal は内部エラーが発生した場合のエラー。
	ErrInternal = errors.New("internal error")
)

// WrapNotFound は ErrNotFound をラップしてメッセージを付加する。
func WrapNotFound(msg string) error {
	return fmt.Errorf("%s: %w", msg, ErrNotFound)
}

// WrapInvalidInput は ErrInvalidInput をラップしてメッセージを付加する。
func WrapInvalidInput(msg string) error {
	return fmt.Errorf("%s: %w", msg, ErrInvalidInput)
}
