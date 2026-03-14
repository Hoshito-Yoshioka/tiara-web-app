INSERT INTO menu_categories (id, name, description, sort_order) VALUES
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'System', 'チャージ・セット料金', 1),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'Cocktails', '', 2),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'Whisky & Spirits', '', 3),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'Beer & Wine', '', 4);

INSERT INTO menu_items (category_id, name, price, description, sort_order) VALUES
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'チャージ', '¥1,000', 'お一人様あたり', 1),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'セット料金（60分）', '¥3,000', 'チャージ＋ドリンク2杯', 2),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', '延長（30分）', '¥1,500', 'ドリンク1杯付き', 3),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'スタンダードカクテル', '¥800〜', '', 1),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'プレミアムカクテル', '¥1,200〜', '', 2),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'オリジナルカクテル', '¥1,000〜', '', 3),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'ノンアルコールカクテル', '¥700〜', '', 4),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'ハウスウイスキー', '¥800', '', 1),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'プレミアムウイスキー', '¥1,200〜', '', 2),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'ブランデー', '¥1,000〜', '', 3),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'ジン / ウォッカ / ラム', '¥800〜', '', 4),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', '生ビール', '¥700', '', 1),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'クラフトビール', '¥900〜', '', 2),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'グラスワイン（赤・白）', '¥800〜', '', 3),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'シャンパン（グラス）', '¥1,500〜', '', 4),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'シャンパン（ボトル）', '¥8,000〜', '', 5);
