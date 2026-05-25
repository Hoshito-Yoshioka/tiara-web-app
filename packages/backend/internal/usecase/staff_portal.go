package usecase

import (
	"context"
	"fmt"
	"tiara-web-app/backend/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// --- Repository Interfaces ---

// StaffAccountRepository はスタッフアカウントの永続化を抽象化するインターフェース。
type StaffAccountRepository interface {
	GetStaffAccountByUsername(ctx context.Context, username string) (domain.StaffAccount, error)
	GetStaffAccountByStaffID(ctx context.Context, staffID uuid.UUID) (domain.StaffAccount, error)
	ListStaffAccounts(ctx context.Context) ([]domain.StaffAccount, error)
	CreateStaffAccount(ctx context.Context, staffID uuid.UUID, username, passwordHash string) (domain.StaffAccount, error)
	UpdateStaffAccount(ctx context.Context, id uuid.UUID, username, passwordHash string) (domain.StaffAccount, error)
	DeleteStaffAccount(ctx context.Context, id uuid.UUID) error
}

// StaffDraftRepository はプロフィール/スケジュール下書きの永続化を抽象化するインターフェース。
type StaffDraftRepository interface {
	// Profile Draft
	GetProfileDraftByStaffID(ctx context.Context, staffID uuid.UUID) (domain.StaffProfileDraft, error)
	GetProfileDraftByID(ctx context.Context, id uuid.UUID) (domain.StaffProfileDraft, error)
	ListPendingProfileDrafts(ctx context.Context) ([]domain.StaffProfileDraft, error)
	CreateProfileDraft(ctx context.Context, staffID uuid.UUID, input domain.SaveProfileDraftInput) (domain.StaffProfileDraft, error)
	UpdateProfileDraft(ctx context.Context, id uuid.UUID, input domain.SaveProfileDraftInput) (domain.StaffProfileDraft, error)
	SubmitProfileDraft(ctx context.Context, id uuid.UUID) (domain.StaffProfileDraft, error)
	ReviewProfileDraft(ctx context.Context, id uuid.UUID, input domain.ReviewDraftInput) (domain.StaffProfileDraft, error)
	DeleteProfileDraft(ctx context.Context, id uuid.UUID) error
	// Schedule Draft
	GetScheduleDraftByStaffID(ctx context.Context, staffID uuid.UUID) (domain.StaffScheduleDraft, error)
	GetScheduleDraftByID(ctx context.Context, id uuid.UUID) (domain.StaffScheduleDraft, error)
	ListPendingScheduleDrafts(ctx context.Context) ([]domain.StaffScheduleDraft, error)
	ListApprovedScheduleDrafts(ctx context.Context) ([]domain.StaffScheduleDraft, error)
	CreateScheduleDraft(ctx context.Context, staffID uuid.UUID, items []domain.ScheduleDraftItem) (domain.StaffScheduleDraft, error)
	UpdateScheduleDraftItems(ctx context.Context, draftID uuid.UUID, items []domain.ScheduleDraftItem) (domain.StaffScheduleDraft, error)
	SubmitScheduleDraft(ctx context.Context, id uuid.UUID) (domain.StaffScheduleDraft, error)
	ReviewScheduleDraft(ctx context.Context, id uuid.UUID, input domain.ReviewDraftInput) (domain.StaffScheduleDraft, error)
	DeleteScheduleDraft(ctx context.Context, id uuid.UUID) error
	// Admin用（ステータス変更なし）
	UpdateProfileDraftContent(ctx context.Context, id uuid.UUID, input domain.SaveProfileDraftInput) (domain.StaffProfileDraft, error)
	ReplaceScheduleDraftItems(ctx context.Context, draftID uuid.UUID, items []domain.ScheduleDraftItem) (domain.StaffScheduleDraft, error)
}

// --- Staff Auth Usecase ---

// StaffAuthUsecase はスタッフ認証のビジネスロジックを定義するインターフェース。
type StaffAuthUsecase interface {
	Login(ctx context.Context, username, password string) (domain.StaffAccount, error)
}

type staffAuthUsecase struct {
	accountRepo StaffAccountRepository
}

// NewStaffAuthUsecase は新しいStaffAuthUsecaseのインスタンスを作成する。
func NewStaffAuthUsecase(repo StaffAccountRepository) StaffAuthUsecase {
	return &staffAuthUsecase{accountRepo: repo}
}

// Login はユーザー名とパスワードでスタッフを認証する。
func (u *staffAuthUsecase) Login(ctx context.Context, username, password string) (domain.StaffAccount, error) {
	account, err := u.accountRepo.GetStaffAccountByUsername(ctx, username)
	if err != nil {
		return domain.StaffAccount{}, fmt.Errorf("invalid credentials: %w", domain.ErrUnauthorized)
	}

	err = bcrypt.CompareHashAndPassword([]byte(account.PasswordHash), []byte(password))
	if err != nil {
		return domain.StaffAccount{}, fmt.Errorf("invalid credentials: %w", domain.ErrUnauthorized)
	}

	return account, nil
}

// --- Staff Portal Usecase ---

// StaffPortalUsecase はスタッフポータルのビジネスロジックを定義するインターフェース。
// スタッフが自分のプロフィール・スケジュールの下書きを管理する。
type StaffPortalUsecase interface {
	// Profile Draft
	GetMyProfileDraft(ctx context.Context, staffID uuid.UUID) (domain.StaffProfileDraft, error)
	SaveProfileDraft(ctx context.Context, staffID uuid.UUID, input domain.SaveProfileDraftInput) (domain.StaffProfileDraft, error)
	SubmitProfileDraft(ctx context.Context, staffID uuid.UUID, draftID uuid.UUID) (domain.StaffProfileDraft, error)
	// Schedule Draft
	GetMyScheduleDraft(ctx context.Context, staffID uuid.UUID) (domain.StaffScheduleDraft, error)
	SaveScheduleDraft(ctx context.Context, staffID uuid.UUID, items []domain.ScheduleDraftItem) (domain.StaffScheduleDraft, error)
	SubmitScheduleDraft(ctx context.Context, staffID uuid.UUID, draftID uuid.UUID) (domain.StaffScheduleDraft, error)
}

type staffPortalUsecase struct {
	draftRepo StaffDraftRepository
	staffRepo StaffRepository
}

// NewStaffPortalUsecase は新しいStaffPortalUsecaseのインスタンスを作成する。
func NewStaffPortalUsecase(draftRepo StaffDraftRepository, staffRepo StaffRepository) StaffPortalUsecase {
	return &staffPortalUsecase{draftRepo: draftRepo, staffRepo: staffRepo}
}

// GetMyProfileDraft は自分のアクティブなプロフィール下書きを取得する。
// 下書きが存在しない場合は、現在のスタッフ情報から初期データを生成する。
func (u *staffPortalUsecase) GetMyProfileDraft(ctx context.Context, staffID uuid.UUID) (domain.StaffProfileDraft, error) {
	draft, err := u.draftRepo.GetProfileDraftByStaffID(ctx, staffID)
	if err != nil {
		// 下書きが存在しない場合、現在のスタッフ情報を元に返す（保存はしない）
		staff, staffErr := u.staffRepo.GetStaffByID(ctx, staffID.String())
		if staffErr != nil {
			return domain.StaffProfileDraft{}, staffErr
		}
		return domain.StaffProfileDraft{
			StaffID:           staffID,
			Name:              staff.Name,
			Role:              staff.Role,
			Bio:               staff.Bio,
			ImageURL:          staff.ImageURL,
			ImageCropPosition: staff.ImageCropPosition,
			Status:            "", // 下書き未作成を示す
		}, nil
	}
	return draft, nil
}

// SaveProfileDraft はプロフィール下書きを保存する（作成または更新）。
func (u *staffPortalUsecase) SaveProfileDraft(ctx context.Context, staffID uuid.UUID, input domain.SaveProfileDraftInput) (domain.StaffProfileDraft, error) {
	// 既存の下書きを探す
	existing, err := u.draftRepo.GetProfileDraftByStaffID(ctx, staffID)
	if err != nil {
		// 存在しない場合は新規作成
		return u.draftRepo.CreateProfileDraft(ctx, staffID, input)
	}

	// 既存の下書きを更新（pending でも更新可能 → ステータスは draft に戻る）
	return u.draftRepo.UpdateProfileDraft(ctx, existing.ID, input)
}

// SubmitProfileDraft はプロフィール下書きを承認申請する。
// 自分の下書きのみ申請可能。
func (u *staffPortalUsecase) SubmitProfileDraft(ctx context.Context, staffID uuid.UUID, draftID uuid.UUID) (domain.StaffProfileDraft, error) {
	draft, err := u.draftRepo.GetProfileDraftByID(ctx, draftID)
	if err != nil {
		return domain.StaffProfileDraft{}, err
	}

	// 自分の下書きのみ操作可能
	if draft.StaffID != staffID {
		return domain.StaffProfileDraft{}, fmt.Errorf("他のスタッフの下書きは操作できません: %w", domain.ErrForbidden)
	}

	if !draft.Status.IsEditable() {
		return domain.StaffProfileDraft{}, fmt.Errorf("draft または rejected 状態の下書きのみ申請できます: %w", domain.ErrInvalidInput)
	}

	return u.draftRepo.SubmitProfileDraft(ctx, draftID)
}

// GetMyScheduleDraft は自分のアクティブなスケジュール下書きを取得する。
func (u *staffPortalUsecase) GetMyScheduleDraft(ctx context.Context, staffID uuid.UUID) (domain.StaffScheduleDraft, error) {
	draft, err := u.draftRepo.GetScheduleDraftByStaffID(ctx, staffID)
	if err != nil {
		// 下書きが存在しない場合、現在のスケジュール情報を元に返す
		schedules, schedErr := u.staffRepo.ListSchedulesByStaffID(ctx, staffID.String())
		if schedErr != nil {
			return domain.StaffScheduleDraft{}, schedErr
		}
		items := make([]domain.ScheduleDraftItem, len(schedules))
		for i, s := range schedules {
			items[i] = domain.ScheduleDraftItem{
				DayOfWeek: s.DayOfWeek,
				StartTime: s.StartTime,
				EndTime:   s.EndTime,
			}
		}
		return domain.StaffScheduleDraft{
			StaffID: staffID,
			Status:  "", // 下書き未作成
			Items:   items,
		}, nil
	}
	return draft, nil
}

// SaveScheduleDraft はスケジュール下書きを保存する（作成または更新）。
func (u *staffPortalUsecase) SaveScheduleDraft(ctx context.Context, staffID uuid.UUID, items []domain.ScheduleDraftItem) (domain.StaffScheduleDraft, error) {
	existing, err := u.draftRepo.GetScheduleDraftByStaffID(ctx, staffID)
	if err != nil {
		// 存在しない場合は新規作成
		return u.draftRepo.CreateScheduleDraft(ctx, staffID, items)
	}

	// 既存の下書きを更新（pending でも更新可能 → ステータスは draft に戻る）
	return u.draftRepo.UpdateScheduleDraftItems(ctx, existing.ID, items)
}

// SubmitScheduleDraft はスケジュール下書きを承認申請する。
func (u *staffPortalUsecase) SubmitScheduleDraft(ctx context.Context, staffID uuid.UUID, draftID uuid.UUID) (domain.StaffScheduleDraft, error) {
	draft, err := u.draftRepo.GetScheduleDraftByID(ctx, draftID)
	if err != nil {
		return domain.StaffScheduleDraft{}, err
	}

	if draft.StaffID != staffID {
		return domain.StaffScheduleDraft{}, fmt.Errorf("他のスタッフの下書きは操作できません: %w", domain.ErrForbidden)
	}

	if !draft.Status.IsEditable() {
		return domain.StaffScheduleDraft{}, fmt.Errorf("draft または rejected 状態の下書きのみ申請できます: %w", domain.ErrInvalidInput)
	}

	return u.draftRepo.SubmitScheduleDraft(ctx, draftID)
}

// --- Admin Review Usecase ---

// AdminReviewUsecase は管理者による下書きレビューのビジネスロジック。
// 承認時にライブデータ（staffs, staff_schedules）に反映する。
type AdminReviewUsecase interface {
	ListPendingProfileDrafts(ctx context.Context) ([]domain.StaffProfileDraft, error)
	ListPendingScheduleDrafts(ctx context.Context) ([]domain.StaffScheduleDraft, error)
	ListApprovedScheduleDrafts(ctx context.Context) ([]domain.StaffScheduleDraft, error)
	ReviewProfileDraft(ctx context.Context, draftID uuid.UUID, input domain.ReviewDraftInput) (domain.StaffProfileDraft, error)
	ReviewScheduleDraft(ctx context.Context, draftID uuid.UUID, input domain.ReviewDraftInput) (domain.StaffScheduleDraft, error)
	PublishScheduleDraft(ctx context.Context, draftID uuid.UUID) error
	// 単体取得
	GetProfileDraft(ctx context.Context, draftID uuid.UUID) (domain.StaffProfileDraft, error)
	GetScheduleDraft(ctx context.Context, draftID uuid.UUID) (domain.StaffScheduleDraft, error)
	// Admin内容修正（ステータス変更なし）
	UpdateProfileDraftContent(ctx context.Context, draftID uuid.UUID, input domain.SaveProfileDraftInput) (domain.StaffProfileDraft, error)
	UpdateScheduleDraftContent(ctx context.Context, draftID uuid.UUID, items []domain.ScheduleDraftItem) (domain.StaffScheduleDraft, error)
	// スタッフ名取得
	GetStaffName(ctx context.Context, staffID uuid.UUID) (string, error)
	// スタッフ画像取得（staff_images テーブルから複数画像を返す）
	ListImagesByStaffID(ctx context.Context, staffID uuid.UUID) ([]domain.StaffImage, error)
}

type adminReviewUsecase struct {
	draftRepo StaffDraftRepository
	staffRepo StaffRepository
}

// NewAdminReviewUsecase は新しいAdminReviewUsecaseのインスタンスを作成する。
func NewAdminReviewUsecase(draftRepo StaffDraftRepository, staffRepo StaffRepository) AdminReviewUsecase {
	return &adminReviewUsecase{draftRepo: draftRepo, staffRepo: staffRepo}
}

// ListPendingProfileDrafts は承認待ちのプロフィール下書き一覧を返す。
func (u *adminReviewUsecase) ListPendingProfileDrafts(ctx context.Context) ([]domain.StaffProfileDraft, error) {
	return u.draftRepo.ListPendingProfileDrafts(ctx)
}

// ListImagesByStaffID はスタッフの画像一覧を返す（staff_images テーブル）。
func (u *adminReviewUsecase) ListImagesByStaffID(ctx context.Context, staffID uuid.UUID) ([]domain.StaffImage, error) {
	return u.staffRepo.ListImagesByStaffID(ctx, staffID.String())
}

// ListPendingScheduleDrafts は承認待ちのスケジュール下書き一覧を返す。
func (u *adminReviewUsecase) ListPendingScheduleDrafts(ctx context.Context) ([]domain.StaffScheduleDraft, error) {
	return u.draftRepo.ListPendingScheduleDrafts(ctx)
}

func (u *adminReviewUsecase) ListApprovedScheduleDrafts(ctx context.Context) ([]domain.StaffScheduleDraft, error) {
	return u.draftRepo.ListApprovedScheduleDrafts(ctx)
}

// ReviewProfileDraft は管理者がプロフィール下書きをレビューする。
// 承認時はライブの staffs テーブルに反映する。
func (u *adminReviewUsecase) ReviewProfileDraft(ctx context.Context, draftID uuid.UUID, input domain.ReviewDraftInput) (domain.StaffProfileDraft, error) {
	if input.Status != domain.DraftStatusApproved && input.Status != domain.DraftStatusRejected {
		return domain.StaffProfileDraft{}, fmt.Errorf("ステータスは approved または rejected を指定してください: %w", domain.ErrInvalidInput)
	}

	draft, err := u.draftRepo.GetProfileDraftByID(ctx, draftID)
	if err != nil {
		return domain.StaffProfileDraft{}, fmt.Errorf("下書きが見つかりません: %w", domain.ErrNotFound)
	}

	if draft.Status != domain.DraftStatusPending {
		return domain.StaffProfileDraft{}, fmt.Errorf("pending 状態の下書きのみレビューできます: %w", domain.ErrInvalidInput)
	}

	// レビュー結果を記録
	reviewed, err := u.draftRepo.ReviewProfileDraft(ctx, draftID, input)
	if err != nil {
		return domain.StaffProfileDraft{}, err
	}

	// 承認の場合、ライブデータに反映
	if input.Status == domain.DraftStatusApproved {
		_, err := u.staffRepo.UpdateStaff(ctx, draft.StaffID.String(), domain.UpdateStaffInput{
			Name:              draft.Name,
			Role:              draft.Role,
			Bio:               draft.Bio,
			ImageURL:          draft.ImageURL,
			ImageCropPosition: draft.ImageCropPosition,
		})
		if err != nil {
			return domain.StaffProfileDraft{}, fmt.Errorf("ライブデータへの反映に失敗しました: %w", domain.ErrInternal)
		}
	}

	return reviewed, nil
}

// ReviewScheduleDraft は管理者がスケジュール下書きをレビューする。
// 承認してもライブデータには即時反映せず、別途 PublishScheduleDraft で反映する。
func (u *adminReviewUsecase) ReviewScheduleDraft(ctx context.Context, draftID uuid.UUID, input domain.ReviewDraftInput) (domain.StaffScheduleDraft, error) {
	if input.Status != domain.DraftStatusApproved && input.Status != domain.DraftStatusRejected {
		return domain.StaffScheduleDraft{}, fmt.Errorf("ステータスは approved または rejected を指定してください: %w", domain.ErrInvalidInput)
	}

	draft, err := u.draftRepo.GetScheduleDraftByID(ctx, draftID)
	if err != nil {
		return domain.StaffScheduleDraft{}, fmt.Errorf("下書きが見つかりません: %w", domain.ErrNotFound)
	}

	if draft.Status != domain.DraftStatusPending {
		return domain.StaffScheduleDraft{}, fmt.Errorf("pending 状態の下書きのみレビューできます: %w", domain.ErrInvalidInput)
	}

	// 承認しても店舗ページには即時反映しない。
	// 管理者が別途 PublishScheduleDraft を呼び出して反映する。
	reviewed, err := u.draftRepo.ReviewScheduleDraft(ctx, draftID, input)
	if err != nil {
		return domain.StaffScheduleDraft{}, err
	}

	return reviewed, nil
}

// PublishScheduleDraft は承認済みスケジュール下書きをライブデータに反映し、下書きを削除する。
// 2段階フロー: 承認 → 店舗反映 の「店舗反映」ステップ。
func (u *adminReviewUsecase) PublishScheduleDraft(ctx context.Context, draftID uuid.UUID) error {
	draft, err := u.draftRepo.GetScheduleDraftByID(ctx, draftID)
	if err != nil {
		return fmt.Errorf("下書きが見つかりません: %w", domain.ErrNotFound)
	}
	if draft.Status != domain.DraftStatusApproved {
		return fmt.Errorf("approved 状態の下書きのみ反映できます: %w", domain.ErrInvalidInput)
	}
	scheduleInputs := make([]domain.ScheduleInput, len(draft.Items))
	for i, item := range draft.Items {
		scheduleInputs[i] = domain.ScheduleInput{
			DayOfWeek: item.DayOfWeek,
			StartTime: item.StartTime.Format("15:04"),
			EndTime:   item.EndTime.Format("15:04"),
		}
	}
	_, err = u.staffRepo.ReplaceSchedules(ctx, draft.StaffID.String(), scheduleInputs)
	if err != nil {
		return fmt.Errorf("ライブスケジュールへの反映に失敗しました: %w", domain.ErrInternal)
	}
	// 反映後、下書きを削除
	return u.draftRepo.DeleteScheduleDraft(ctx, draftID)
}

// GetProfileDraft は管理者がプロフィール下書きを単体で取得する。
func (u *adminReviewUsecase) GetProfileDraft(ctx context.Context, draftID uuid.UUID) (domain.StaffProfileDraft, error) {
	return u.draftRepo.GetProfileDraftByID(ctx, draftID)
}

// GetScheduleDraft は管理者がスケジュール下書きを単体で取得する。
func (u *adminReviewUsecase) GetScheduleDraft(ctx context.Context, draftID uuid.UUID) (domain.StaffScheduleDraft, error) {
	return u.draftRepo.GetScheduleDraftByID(ctx, draftID)
}

// UpdateProfileDraftContent は管理者がプロフィール下書きの内容のみ更新する（ステータス変更なし）。
func (u *adminReviewUsecase) UpdateProfileDraftContent(ctx context.Context, draftID uuid.UUID, input domain.SaveProfileDraftInput) (domain.StaffProfileDraft, error) {
	return u.draftRepo.UpdateProfileDraftContent(ctx, draftID, input)
}

// UpdateScheduleDraftContent は管理者がスケジュール下書きの内容のみ更新する（ステータス変更なし）。
func (u *adminReviewUsecase) UpdateScheduleDraftContent(ctx context.Context, draftID uuid.UUID, items []domain.ScheduleDraftItem) (domain.StaffScheduleDraft, error) {
	return u.draftRepo.ReplaceScheduleDraftItems(ctx, draftID, items)
}

// GetStaffName はスタッフIDからスタッフ名を取得する。
func (u *adminReviewUsecase) GetStaffName(ctx context.Context, staffID uuid.UUID) (string, error) {
	staff, err := u.staffRepo.GetStaffByID(ctx, staffID.String())
	if err != nil {
		return "", err
	}
	return staff.Name, nil
}

// --- Admin Account Management Usecase ---

// AdminAccountUsecase は管理者によるスタッフアカウント管理のビジネスロジック。
type AdminAccountUsecase interface {
	ListStaffAccounts(ctx context.Context) ([]domain.StaffAccount, error)
	GetStaffAccountByStaffID(ctx context.Context, staffID uuid.UUID) (domain.StaffAccount, error)
	CreateStaffAccount(ctx context.Context, staffID uuid.UUID, username, password string) (domain.StaffAccount, error)
	UpdateStaffAccount(ctx context.Context, id uuid.UUID, username, password string) (domain.StaffAccount, error)
	DeleteStaffAccount(ctx context.Context, id uuid.UUID) error
}

type adminAccountUsecase struct {
	accountRepo StaffAccountRepository
}

// NewAdminAccountUsecase は新しいAdminAccountUsecaseのインスタンスを作成する。
func NewAdminAccountUsecase(repo StaffAccountRepository) AdminAccountUsecase {
	return &adminAccountUsecase{accountRepo: repo}
}

// ListStaffAccounts は全スタッフアカウント一覧を返す。
func (u *adminAccountUsecase) ListStaffAccounts(ctx context.Context) ([]domain.StaffAccount, error) {
	return u.accountRepo.ListStaffAccounts(ctx)
}

// GetStaffAccountByStaffID はスタッフIDでアカウントを返す。
func (u *adminAccountUsecase) GetStaffAccountByStaffID(ctx context.Context, staffID uuid.UUID) (domain.StaffAccount, error) {
	return u.accountRepo.GetStaffAccountByStaffID(ctx, staffID)
}

// CreateStaffAccount は新しいスタッフアカウントを作成する（パスワードをbcryptハッシュ化）。
func (u *adminAccountUsecase) CreateStaffAccount(ctx context.Context, staffID uuid.UUID, username, password string) (domain.StaffAccount, error) {
	if username == "" || password == "" {
		return domain.StaffAccount{}, fmt.Errorf("ユーザー名とパスワードは必須です: %w", domain.ErrInvalidInput)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return domain.StaffAccount{}, fmt.Errorf("パスワードのハッシュ化に失敗しました: %w", domain.ErrInternal)
	}

	return u.accountRepo.CreateStaffAccount(ctx, staffID, username, string(hash))
}

// UpdateStaffAccount はスタッフアカウントのユーザー名とパスワードを更新する。
// password が空の場合はパスワードを変更せず、現在のハッシュを維持する。
func (u *adminAccountUsecase) UpdateStaffAccount(ctx context.Context, id uuid.UUID, username, password string) (domain.StaffAccount, error) {
	if username == "" {
		return domain.StaffAccount{}, fmt.Errorf("ユーザー名は必須です: %w", domain.ErrInvalidInput)
	}

	var passwordHash string
	if password != "" {
		hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return domain.StaffAccount{}, fmt.Errorf("パスワードのハッシュ化に失敗しました: %w", domain.ErrInternal)
		}
		passwordHash = string(hash)
	} else {
		// パスワード未変更 — 既存のハッシュを取得して維持する
		accounts, err := u.accountRepo.ListStaffAccounts(ctx)
		if err != nil {
			return domain.StaffAccount{}, err
		}
		for _, a := range accounts {
			if a.ID == id {
				passwordHash = a.PasswordHash
				break
			}
		}
		if passwordHash == "" {
			return domain.StaffAccount{}, fmt.Errorf("アカウントが見つかりません: %w", domain.ErrNotFound)
		}
	}

	return u.accountRepo.UpdateStaffAccount(ctx, id, username, passwordHash)
}

// DeleteStaffAccount はスタッフアカウントを削除する。
func (u *adminAccountUsecase) DeleteStaffAccount(ctx context.Context, id uuid.UUID) error {
	return u.accountRepo.DeleteStaffAccount(ctx, id)
}
