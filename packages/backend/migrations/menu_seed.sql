INSERT INTO menu_categories (id, name, description, sort_order) VALUES
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'System', 'セット料金・指名料・席料', 1),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'Drinks', 'ハウスボトル・各種ドリンク', 2),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'Bottle', 'シャンパン・ワイン・各種ボトル', 3),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'Food', 'フード・おつまみ', 4);

INSERT INTO menu_items (category_id, name, price, description, sort_order) VALUES
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'セット料金（60分）', '¥5,000', 'ハウスボトル飲み放題', 1),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', '延長（30分）', '¥2,500', '', 2),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', '指名料', '¥1,000', '', 3),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', 'VIP席料金', '¥3,000', 'セット料金に加算', 4),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a01', '同伴料', '¥3,000', '', 5),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'ハウスボトル', '飲み放題', 'ウイスキー・焼酎・ブランデー等', 1),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'ソフトドリンク', '¥500', '', 2),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a02', 'キャストドリンク', '¥1,000', '', 3),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'シャンパン（グラス）', '¥2,000〜', '', 1),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'シャンパン（ボトル）', '¥10,000〜', '', 2),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'ワイン（ボトル）', '¥8,000〜', '', 3),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a03', 'プレミアムボトル', '¥15,000〜', 'ブランデー・高級ウイスキー等', 4),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'フルーツ盛り合わせ', '¥2,000', '', 1),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'ミックスナッツ', '¥500', '', 2),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'チーズ盛り合わせ', '¥1,200', '', 3),
('c1eebc99-9c0b-4ef8-bb6d-6bb9bd380a04', 'チョコレートアソート', '¥800', '', 4);
