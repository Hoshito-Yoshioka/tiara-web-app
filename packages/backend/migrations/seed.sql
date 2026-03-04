-- 管理者データ（パスワード: admin123）
INSERT INTO admins (username, password_hash) VALUES
('admin', crypt('admin123', gen_salt('bf', 10)));

INSERT INTO shops (id, name, address, opening_time, closing_time) VALUES
('a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'TEST BAR Tiara', '東京都渋谷区テスト1-2-3', '18:00:00', '02:00:00');

-- スタッフデータ
INSERT INTO staffs (id, shop_id, name, role, bio, image_url, sort_order) VALUES
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Yuki Tanaka', 'バーテンダー', '10年以上の経験を持つベテランバーテンダー。クラシックカクテルを得意とし、お客様一人一人に合わせたオリジナルカクテルの提案が好評。', 'https://images.unsplash.com/photo-1507003211169-0a1dd7228f2d?w=400&h=400&fit=crop&crop=face', 1),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Rina Suzuki', 'チーフバーテンダー', '国内外のカクテルコンペティションで入賞経験あり。フルーツを活かした華やかなカクテルが得意。', 'https://images.unsplash.com/photo-1494790108377-be9c29b29330?w=400&h=400&fit=crop&crop=face', 2),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'a0eebc99-9c0b-4ef8-bb6d-6bb9bd380a11', 'Kenji Yamamoto', 'ホールスタッフ', 'ソムリエ資格保有。ワインとフードのペアリング提案に定評があり、落ち着いた接客スタイルが特徴。', 'https://images.unsplash.com/photo-1472099645785-5658abf4ff4e?w=400&h=400&fit=crop&crop=face', 3);

-- スタッフ出勤スケジュール
-- Yuki Tanaka: 月・水・金・土
INSERT INTO staff_schedules (staff_id, day_of_week, start_time, end_time) VALUES
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 1, '18:00:00', '02:00:00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 3, '18:00:00', '02:00:00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 5, '18:00:00', '02:00:00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 6, '19:00:00', '03:00:00');

-- Rina Suzuki: 火・木・金・土
INSERT INTO staff_schedules (staff_id, day_of_week, start_time, end_time) VALUES
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 2, '18:00:00', '02:00:00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 4, '18:00:00', '02:00:00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 5, '19:00:00', '03:00:00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 6, '19:00:00', '03:00:00');

-- Kenji Yamamoto: 月・火・水・木・土
INSERT INTO staff_schedules (staff_id, day_of_week, start_time, end_time) VALUES
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 1, '18:00:00', '02:00:00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 2, '18:00:00', '02:00:00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 3, '18:00:00', '02:00:00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 4, '18:00:00', '02:00:00'),
('b1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 6, '19:00:00', '03:00:00');
