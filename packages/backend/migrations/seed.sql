-- 管理者データ（パスワード: admin123）
INSERT INTO admins (username, password_hash) VALUES
('admin', crypt('admin123', gen_salt('bf', 10)));

INSERT INTO shops (id, name, address, opening_time, closing_time) VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'TEST BAR Tiara', '東京都渋谷区テスト1-2-3', '20:00:00', '02:00:00');

-- スタッフデータ
INSERT INTO staffs (id, shop_id, name, role, bio, image_url, sort_order) VALUES
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Yuki Tanaka', 'バーテンダー', '10年以上の経験を持つベテランバーテンダー。クラシックカクテルを得意とし、お客様一人一人に合わせたオリジナルカクテルの提案が好評。', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=400&h=400&fit=crop&crop=face', 1),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Rina Suzuki', 'チーフバーテンダー', '国内外のカクテルコンペティションで入賞経験あり。フルーツを活かした華やかなカクテルが得意。', 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=400&h=400&fit=crop&crop=face', 2),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Kenji Yamamoto', 'ホールスタッフ', 'ソムリエ資格保有。ワインとフードのペアリング提案に定評があり、落ち着いた接客スタイルが特徴。', 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=400&h=400&fit=crop&crop=face', 3);

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
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'System', 'チャージ・セット料金', 1),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'Cocktails', '', 2),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'Whisky & Spirits', '', 3),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'Beer & Wine', '', 4);

INSERT INTO menu_items (category_id, name, price, description, sort_order) VALUES
-- System
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'チャージ', '¥1,000', 'お一人様あたり', 1),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'セット料金（60分）', '¥3,000', '', 2),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', '延長（30分）', '¥1,500', '', 3),
-- Cocktails
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'スタンダードカクテル', '¥800〜', '', 1),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'プレミアムカクテル', '¥1,200〜', '', 2),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'オリジナルカクテル', '¥1,000〜', '', 3),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'ノンアルコールカクテル', '¥700〜', '', 4),
-- Whisky & Spirits
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'ハウスウイスキー', '¥800', '', 1),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'プレミアムウイスキー', '¥1,200〜', '', 2),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'ブランデー', '¥1,000〜', '', 3),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'ジン / ウォッカ / ラム', '¥800〜', '', 4),
-- Beer & Wine
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', '生ビール', '¥700', '', 1),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'クラフトビール', '¥900〜', '', 2),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'グラスワイン（赤・白）', '¥800〜', '', 3),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'シャンパン（グラス）', '¥1,500〜', '', 4),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'シャンパン（ボトル）', '¥8,000〜', '', 5);

-- ==============================
-- スタッフアカウント（パスワード: staff123）
-- ==============================
INSERT INTO staff_accounts (staff_id, username, password_hash) VALUES
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'yuki', crypt('staff123', gen_salt('bf', 10))),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'rina', crypt('staff123', gen_salt('bf', 10))),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'kenji', crypt('staff123', gen_salt('bf', 10)));
