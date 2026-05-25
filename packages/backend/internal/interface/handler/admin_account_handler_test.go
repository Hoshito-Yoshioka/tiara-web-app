package handler

import (
	"context"
	"encoding/json"
	"net/http"
	"testing"

	"tiara-web-app/backend/internal/domain"
	"tiara-web-app/backend/internal/testutil"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

// mockAdminAccountUsecase は AdminAccountUsecase のテスト用モック。
type mockAdminAccountUsecase struct {
	accounts []domain.StaffAccount
	account  domain.StaffAccount
	err      error
}

func (m *mockAdminAccountUsecase) ListStaffAccounts(_ context.Context) ([]domain.StaffAccount, error) {
	return m.accounts, m.err
}
func (m *mockAdminAccountUsecase) GetStaffAccountByStaffID(_ context.Context, _ uuid.UUID) (domain.StaffAccount, error) {
	return m.account, m.err
}
func (m *mockAdminAccountUsecase) CreateStaffAccount(_ context.Context, _ uuid.UUID, _, _ string) (domain.StaffAccount, error) {
	return m.account, m.err
}
func (m *mockAdminAccountUsecase) UpdateStaffAccount(_ context.Context, _ uuid.UUID, _, _ string) (domain.StaffAccount, error) {
	return m.account, m.err
}
func (m *mockAdminAccountUsecase) DeleteStaffAccount(_ context.Context, _ uuid.UUID) error {
	return m.err
}

func TestAdminAccountHandler_ListStaffAccounts(t *testing.T) {
	t.Run("正常系: アカウント一覧を返す", func(t *testing.T) {
		account := testutil.NewStaffAccount()
		account.Username = "staff1"
		mock := &mockAdminAccountUsecase{accounts: []domain.StaffAccount{account}}
		h := NewAdminAccountHandler(mock)

		c, rec := testutil.NewEchoContext(http.MethodGet, "/admin/staff-accounts")

		err := h.ListStaffAccounts(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)

		var resp []StaffAccountResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Len(t, resp, 1)
		assert.Equal(t, "staff1", resp[0].Username)
	})

	t.Run("異常系: usecase エラー → 500", func(t *testing.T) {
		mock := &mockAdminAccountUsecase{err: domain.ErrInternal}
		h := NewAdminAccountHandler(mock)

		c, rec := testutil.NewEchoContext(http.MethodGet, "/admin/staff-accounts")

		_ = h.ListStaffAccounts(c)
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
	})
}

func TestAdminAccountHandler_GetStaffAccountByStaffID(t *testing.T) {
	staffID := uuid.New()

	t.Run("正常系: スタッフIDでアカウントを返す", func(t *testing.T) {
		account := testutil.NewStaffAccount()
		account.StaffID = staffID
		account.Username = "staff1"
		mock := &mockAdminAccountUsecase{account: account}
		h := NewAdminAccountHandler(mock)

		c, rec := testutil.NewEchoContext(http.MethodGet, "/admin/staff-accounts/staff/:staffId")
		testutil.SetPathParams(c, "staffId", staffID.String())

		err := h.GetStaffAccountByStaffID(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: 不正なスタッフID → 400", func(t *testing.T) {
		h := NewAdminAccountHandler(&mockAdminAccountUsecase{})

		c, rec := testutil.NewEchoContext(http.MethodGet, "/admin/staff-accounts/staff/:staffId")
		testutil.SetPathParams(c, "staffId", "invalid-uuid")

		_ = h.GetStaffAccountByStaffID(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: アカウント未存在 → 空レスポンス（200）", func(t *testing.T) {
		mock := &mockAdminAccountUsecase{err: domain.ErrNotFound}
		h := NewAdminAccountHandler(mock)

		c, rec := testutil.NewEchoContext(http.MethodGet, "/admin/staff-accounts/staff/:staffId")
		testutil.SetPathParams(c, "staffId", uuid.New().String())

		err := h.GetStaffAccountByStaffID(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestAdminAccountHandler_CreateStaffAccount(t *testing.T) {
	staffID := uuid.New()

	t.Run("正常系: アカウント作成 → 201", func(t *testing.T) {
		account := testutil.NewStaffAccount()
		account.StaffID = staffID
		account.Username = "newuser"
		mock := &mockAdminAccountUsecase{account: account}
		h := NewAdminAccountHandler(mock)

		body := `{"staffId":"` + staffID.String() + `","username":"newuser","password":"pass1234"}`
		c, rec := testutil.NewEchoContext(http.MethodPost, "/admin/staff-accounts", body)

		err := h.CreateStaffAccount(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusCreated, rec.Code)

		var resp StaffAccountResponse
		err = json.Unmarshal(rec.Body.Bytes(), &resp)
		assert.NoError(t, err)
		assert.Equal(t, "newuser", resp.Username)
	})

	t.Run("異常系: 不正なボディ → 400", func(t *testing.T) {
		h := NewAdminAccountHandler(&mockAdminAccountUsecase{})

		c, rec := testutil.NewEchoContext(http.MethodPost, "/admin/staff-accounts", "bad")

		_ = h.CreateStaffAccount(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 不正なスタッフID → 400", func(t *testing.T) {
		h := NewAdminAccountHandler(&mockAdminAccountUsecase{})

		c, rec := testutil.NewEchoContext(http.MethodPost, "/admin/staff-accounts", `{"staffId":"invalid","username":"user","password":"pass"}`)

		_ = h.CreateStaffAccount(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: usecase エラー（重複）→ 409", func(t *testing.T) {
		mock := &mockAdminAccountUsecase{err: domain.ErrConflict}
		h := NewAdminAccountHandler(mock)

		body := `{"staffId":"` + staffID.String() + `","username":"dup","password":"pass"}`
		c, rec := testutil.NewEchoContext(http.MethodPost, "/admin/staff-accounts", body)

		_ = h.CreateStaffAccount(c)
		assert.Equal(t, http.StatusConflict, rec.Code)
	})
}

func TestAdminAccountHandler_UpdateStaffAccount(t *testing.T) {
	accountID := uuid.New()

	t.Run("正常系: アカウント更新 → 200", func(t *testing.T) {
		mock := &mockAdminAccountUsecase{
			account: domain.StaffAccount{ID: accountID, Username: "updated"},
		}
		h := NewAdminAccountHandler(mock)

		c, rec := testutil.NewEchoContext(http.MethodPut, "/admin/staff-accounts/:id", `{"username":"updated","password":"newpass"}`)
		testutil.SetPathParams(c, "id", accountID.String())

		err := h.UpdateStaffAccount(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusOK, rec.Code)
	})

	t.Run("異常系: 不正なID → 400", func(t *testing.T) {
		h := NewAdminAccountHandler(&mockAdminAccountUsecase{})

		c, rec := testutil.NewEchoContext(http.MethodPut, "/admin/staff-accounts/:id", `{"username":"u","password":"p"}`)
		testutil.SetPathParams(c, "id", "bad-uuid")

		_ = h.UpdateStaffAccount(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 不正なボディ → 400", func(t *testing.T) {
		h := NewAdminAccountHandler(&mockAdminAccountUsecase{})

		c, rec := testutil.NewEchoContext(http.MethodPut, "/admin/staff-accounts/:id", "bad")
		testutil.SetPathParams(c, "id", accountID.String())

		_ = h.UpdateStaffAccount(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})
}

func TestAdminAccountHandler_DeleteStaffAccount(t *testing.T) {
	accountID := uuid.New()

	t.Run("正常系: アカウント削除 → 204", func(t *testing.T) {
		mock := &mockAdminAccountUsecase{}
		h := NewAdminAccountHandler(mock)

		c, rec := testutil.NewEchoContext(http.MethodDelete, "/admin/staff-accounts/:id")
		testutil.SetPathParams(c, "id", accountID.String())

		err := h.DeleteStaffAccount(c)
		assert.NoError(t, err)
		assert.Equal(t, http.StatusNoContent, rec.Code)
	})

	t.Run("異常系: 不正なID → 400", func(t *testing.T) {
		h := NewAdminAccountHandler(&mockAdminAccountUsecase{})

		c, rec := testutil.NewEchoContext(http.MethodDelete, "/admin/staff-accounts/:id")
		testutil.SetPathParams(c, "id", "bad")

		_ = h.DeleteStaffAccount(c)
		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("異常系: 存在しない → 404", func(t *testing.T) {
		mock := &mockAdminAccountUsecase{err: domain.ErrNotFound}
		h := NewAdminAccountHandler(mock)

		c, rec := testutil.NewEchoContext(http.MethodDelete, "/admin/staff-accounts/:id")
		testutil.SetPathParams(c, "id", uuid.New().String())

		_ = h.DeleteStaffAccount(c)
		assert.Equal(t, http.StatusNotFound, rec.Code)
	})
}
