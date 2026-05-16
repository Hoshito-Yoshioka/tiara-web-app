-- ============================================================
-- ⭐ 注意: このファイルは開発環境専用のテストデータです。
-- 本番環境では絶対に実行しないでください。
-- パスワードは平文で記載されており、セキュリティ上のリスクがあります。
-- ============================================================

-- 管理者データ（パスワード: admin123）
INSERT INTO admins (username, password_hash) VALUES
('admin', crypt('admin123', gen_salt('bf', 10)));

INSERT INTO shops (id, name, address, opening_time, closing_time) VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'TEST New Club TIARA', '北海道函館市本町１−２８ 第５大栄ビル', '20:00:00', '02:00:00');

-- スタッフデータ
INSERT INTO staffs (id, shop_id, name, role, bio, image_url, sort_order) VALUES
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Yuki Tanaka', 'キャスト', '入店５年目のベテランキャスト。明るい笑顔と細やかな気配りで、お客様一人一人に対して特別なひとときをお届けします。', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=400&h=400&fit=crop&crop=face', 1),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Rina Suzuki', 'チーフキャスト', '店内No.1の実力派。華やかな雰囲気と知的な会話で、多くのお客様から指名をいただいています。', 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=400&h=400&fit=crop&crop=face', 2),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Kenji Yamamoto', 'フロアスタッフ', 'ホール業務を担当。丁寧な接客とスムーズな案内で、お客様に快適な空間を提供します。', 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=400&h=400&fit=crop&crop=face', 3);

-- スタッフ出勤スケジュール
-- Yuki Tanaka: 月・水・金・土
INSERT INTO staff_schedules (staff_id, day_of_week, start_time, end_time) VALUES
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 1, '20:00:00', '02:00:00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 3, '20:00:00', '02:00:00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 5, '20:00:00', '02:00:00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 6, '20:00:00', '03:00:00');

-- Rina Suzuki: 火・木・金・土
INSERT INTO staff_schedules (staff_id, day_of_week, start_time, end_time) VALUES
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 2, '20:00:00', '02:00:00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 4, '20:00:00', '02:00:00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 5, '20:00:00', '03:00:00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 6, '20:00:00', '03:00:00');

-- Kenji Yamamoto: 月・火・水・木・土
INSERT INTO staff_schedules (staff_id, day_of_week, start_time, end_time) VALUES
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 1, '20:00:00', '02:00:00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 2, '20:00:00', '02:00:00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 3, '20:00:00', '02:00:00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 4, '20:00:00', '02:00:00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 6, '20:00:00', '03:00:00');

-- ==============================
-- メニューカテゴリとメニュー品目
-- ==============================
INSERT INTO menu_categories (id, name, description, sort_order) VALUES
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'System', 'セット料金・指名料・席料', 1),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'Drinks', 'ハウスボトル・各種ドリンク', 2),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'Bottle', 'シャンパン・ワイン・各種ボトル', 3),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'Food', 'フード・おつまみ', 4);

INSERT INTO menu_items (category_id, name, price, description, sort_order) VALUES
-- System
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'セット料金（60分）', '¥5,000', 'ハウスボトル飲み放題', 1),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', '延長（30分）', '¥2,500', '', 2),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', '指名料', '¥1,000', '', 3),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'VIP席料金', '¥3,000', 'セット料金に加算', 4),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', '同伴料', '¥3,000', '', 5),
-- Drinks
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'ハウスボトル', '飲み放題', 'ウイスキー・焼酎・ブランデー等', 1),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'ソフトドリンク', '¥500', '', 2),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'キャストドリンク', '¥1,000', '', 3),
-- Bottle
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'シャンパン（グラス）', '¥2,000〜', '', 1),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'シャンパン（ボトル）', '¥10,000〜', '', 2),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'ワイン（ボトル）', '¥8,000〜', '', 3),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'プレミアムボトル', '¥15,000〜', 'ブランデー・高級ウイスキー等', 4),
-- Food
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'フルーツ盛り合わせ', '¥2,000', '', 1),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'ミックスナッツ', '¥500', '', 2),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'チーズ盛り合わせ', '¥1,200', '', 3),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'チョコレートアソート', '¥800', '', 4);

-- ==============================
-- スタッフアカウント（パスワード: staff123）
-- ==============================
INSERT INTO staff_accounts (staff_id, username, password_hash) VALUES
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'yuki', crypt('staff123', gen_salt('bf', 10))),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'rina', crypt('staff123', gen_salt('bf', 10))),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'kenji', crypt('staff123', gen_salt('bf', 10)));
