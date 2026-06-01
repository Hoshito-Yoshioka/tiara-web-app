CREATE EXTENSION IF NOT EXISTS pgcrypto;

-- ==============================
-- トリガー関数（updated_at 自動更新用）
-- 全テーブルで共有するため、最初に定義する
-- ==============================
CREATE OR REPLACE FUNCTION update_timestamp()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- ==============================
-- admins テーブル
-- 管理者認証用。パスワードは bcrypt ハッシュで保存。
-- ==============================
CREATE TABLE IF NOT EXISTS admins (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

DROP TRIGGER IF EXISTS update_admins_updated_at ON admins;
CREATE TRIGGER update_admins_updated_at
BEFORE UPDATE ON admins
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

CREATE TABLE IF NOT EXISTS shops (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    address TEXT NOT NULL,
    opening_time TIME NOT NULL,
    closing_time TIME NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

DROP TRIGGER IF EXISTS update_shops_updated_at ON shops;
CREATE TRIGGER update_shops_updated_at
BEFORE UPDATE ON shops
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- ==============================
-- staffs テーブル
-- ==============================
CREATE TABLE IF NOT EXISTS staffs (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    shop_id UUID NOT NULL REFERENCES shops(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(100) NOT NULL DEFAULT '',
    bio TEXT NOT NULL DEFAULT '',
    image_url TEXT NOT NULL DEFAULT '',
    external_schedule_url TEXT NOT NULL DEFAULT '',
    image_crop_position VARCHAR(20) NOT NULL DEFAULT '50 50',
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE staffs
    ADD COLUMN IF NOT EXISTS external_schedule_url TEXT NOT NULL DEFAULT '';

DROP TRIGGER IF EXISTS update_staffs_updated_at ON staffs;
CREATE TRIGGER update_staffs_updated_at
BEFORE UPDATE ON staffs
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- ==============================
-- staff_accounts テーブル
-- スタッフ専用ログインアカウント（admin とは別系統）
-- staff_id は staffs テーブルへの参照（1:1）
-- ==============================
CREATE TABLE IF NOT EXISTS staff_accounts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    staff_id UUID NOT NULL UNIQUE REFERENCES staffs(id) ON DELETE CASCADE,
    username VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

DROP TRIGGER IF EXISTS update_staff_accounts_updated_at ON staff_accounts;
CREATE TRIGGER update_staff_accounts_updated_at
BEFORE UPDATE ON staff_accounts
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- ==============================
-- staff_profile_drafts テーブル
-- スタッフが提出するプロフィール変更の下書き/承認待ち
-- status: draft / pending / approved / rejected
-- ==============================
CREATE TABLE IF NOT EXISTS staff_profile_drafts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    staff_id UUID NOT NULL REFERENCES staffs(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    role VARCHAR(100) NOT NULL DEFAULT '',
    bio TEXT NOT NULL DEFAULT '',
    image_url TEXT NOT NULL DEFAULT '',
    external_schedule_url TEXT NOT NULL DEFAULT '',
    image_crop_position VARCHAR(20) NOT NULL DEFAULT '50 50',
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'pending', 'approved', 'rejected')),
    admin_comment TEXT NOT NULL DEFAULT '',
    submitted_at TIMESTAMP WITH TIME ZONE,
    reviewed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

ALTER TABLE staff_profile_drafts
    ADD COLUMN IF NOT EXISTS external_schedule_url TEXT NOT NULL DEFAULT '';

DROP TRIGGER IF EXISTS update_staff_profile_drafts_updated_at ON staff_profile_drafts;
CREATE TRIGGER update_staff_profile_drafts_updated_at
BEFORE UPDATE ON staff_profile_drafts
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- ==============================
-- staff_schedule_drafts テーブル
-- スタッフが提出するシフト変更の下書き/承認待ち
-- status: draft / pending / approved / rejected
-- 承認時に staff_schedules へ反映される
-- ==============================
CREATE TABLE IF NOT EXISTS staff_schedule_drafts (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    staff_id UUID NOT NULL REFERENCES staffs(id) ON DELETE CASCADE,
    status VARCHAR(20) NOT NULL DEFAULT 'draft' CHECK (status IN ('draft', 'pending', 'approved', 'rejected')),
    admin_comment TEXT NOT NULL DEFAULT '',
    submitted_at TIMESTAMP WITH TIME ZONE,
    reviewed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

DROP TRIGGER IF EXISTS update_staff_schedule_drafts_updated_at ON staff_schedule_drafts;
CREATE TRIGGER update_staff_schedule_drafts_updated_at
BEFORE UPDATE ON staff_schedule_drafts
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- ==============================
-- staff_schedule_draft_items テーブル
-- シフト下書きの各曜日データ
-- ==============================
CREATE TABLE IF NOT EXISTS staff_schedule_draft_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    draft_id UUID NOT NULL REFERENCES staff_schedule_drafts(id) ON DELETE CASCADE,
    day_of_week INT NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL
);

-- ==============================
-- staff_images テーブル
-- スタッフの画像（メイン1枚+サブ複数枚）
-- is_main: メイン画像フラグ（スタッフ一覧等で使用）
-- ==============================
CREATE TABLE IF NOT EXISTS staff_images (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    staff_id UUID NOT NULL REFERENCES staffs(id) ON DELETE CASCADE,
    image_url TEXT NOT NULL,
    is_main BOOLEAN NOT NULL DEFAULT false,
    sort_order INT NOT NULL DEFAULT 0,
    crop_position VARCHAR(20) NOT NULL DEFAULT '50 50',
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

DROP TRIGGER IF EXISTS update_staff_images_updated_at ON staff_images;
CREATE TRIGGER update_staff_images_updated_at
BEFORE UPDATE ON staff_images
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- ==============================
-- staff_schedules テーブル
-- スタッフの出勤スケジュール（曜日ベース）
-- day_of_week: 0=日, 1=月, 2=火, ..., 6=土
-- ==============================
CREATE TABLE IF NOT EXISTS staff_schedules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    staff_id UUID NOT NULL REFERENCES staffs(id) ON DELETE CASCADE,
    day_of_week INT NOT NULL CHECK (day_of_week BETWEEN 0 AND 6),
    start_time TIME NOT NULL,
    end_time TIME NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (staff_id, day_of_week)
);

DROP TRIGGER IF EXISTS update_staff_schedules_updated_at ON staff_schedules;
CREATE TRIGGER update_staff_schedules_updated_at
BEFORE UPDATE ON staff_schedules
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- ==============================
-- menu_categories テーブル
-- PRICEページに表示するメニューのカテゴリ（例: System, Drinks, Bottle）
-- ==============================
CREATE TABLE IF NOT EXISTS menu_categories (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

DROP TRIGGER IF EXISTS update_menu_categories_updated_at ON menu_categories;
CREATE TRIGGER update_menu_categories_updated_at
BEFORE UPDATE ON menu_categories
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();

-- ==============================
-- menu_items テーブル
-- 各カテゴリに属するメニュー品目と価格
-- ==============================
CREATE TABLE IF NOT EXISTS menu_items (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    category_id UUID NOT NULL REFERENCES menu_categories(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    price VARCHAR(100) NOT NULL,
    description TEXT NOT NULL DEFAULT '',
    sort_order INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    UNIQUE (category_id, name)
);

DROP TRIGGER IF EXISTS update_menu_items_updated_at ON menu_items;
CREATE TRIGGER update_menu_items_updated_at
BEFORE UPDATE ON menu_items
FOR EACH ROW
EXECUTE FUNCTION update_timestamp();
