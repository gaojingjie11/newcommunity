USE smart_community;

-- Default product categories
INSERT INTO `pms_product_category` (`name`, `icon`, `sort`) VALUES
('生鲜果蔬', '🍎', 1),
('粮油副食', '🌾', 2),
('日用百货', '🧴', 3),
('家居清洁', '🧹', 4),
('个人护理', '🧴', 5)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- Default service areas
INSERT INTO `service_areas` (`name`, `sort`, `status`) VALUES
('东软智慧社区A区', 1, 1),
('东软智慧社区B区', 2, 1),
('东软智慧社区C区', 3, 1)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- Default stores
INSERT INTO `pms_store` (`name`, `address`, `phone`, `area_id`, `region`, `business_hours`) VALUES
('社区便民超市A', 'A区1号楼底商', '024-12345678', 1, 'A区', '08:00-22:00'),
('社区便民超市B', 'B区2号楼底商', '024-87654321', 2, 'B区', '08:00-22:00')
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- Sample products
INSERT INTO `pms_product` (`category_name`, `name`, `description`, `price`, `original_price`, `stock`, `image_url`, `is_promotion`, `sales`, `status`, `category_id`) VALUES
('生鲜果蔬', '新鲜苹果 500g', '产地直供，新鲜脆甜', 1290, 1590, 200, '', 1, 50, 1, 1),
('生鲜果蔬', '有机西红柿 500g', '自然成熟，口感沙甜', 850, 1000, 150, '', 1, 30, 1, 1),
('粮油副食', '东北大米 5kg', '五常稻花香，粒粒饱满', 3990, 4990, 100, '', 1, 80, 1, 2),
('日用百货', '抽纸 3层120抽*10包', '柔软亲肤，不掉屑', 2990, 3500, 300, '', 1, 120, 1, 3),
('家居清洁', '洗洁精 1.5kg', '去油不伤手', 1590, 1890, 200, '', 1, 60, 1, 4),
('个人护理', '牙膏 120g', '清新口气，防蛀固齿', 990, 1200, 250, '', 1, 90, 1, 5)
ON DUPLICATE KEY UPDATE `name` = VALUES(`name`);

-- Sample store-product bindings
INSERT INTO `pms_store_product` (`store_id`, `product_id`, `stock`, `status`) VALUES
(1, 1, 50, 1),
(1, 2, 30, 1),
(1, 3, 20, 1),
(1, 4, 100, 1),
(2, 1, 40, 1),
(2, 3, 25, 1),
(2, 5, 50, 1),
(2, 6, 80, 1)
ON DUPLICATE KEY UPDATE `stock` = VALUES(`stock`);

-- Initialize wallets for existing users (admin + sample users)
INSERT INTO `wallets` (`user_id`, `balance`) VALUES
(1, 100000)
ON DUPLICATE KEY UPDATE `balance` = VALUES(`balance`);
