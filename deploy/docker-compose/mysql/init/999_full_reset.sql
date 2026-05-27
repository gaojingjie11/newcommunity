-- ============================================================
-- 智慧社区数据库完整重置脚本
-- 管理员: admin / sdl@admin / 13483000001
-- ============================================================

USE smart_community;

-- ── 1. 删除所有表（按依赖顺序） ──

SET FOREIGN_KEY_CHECKS = 0;

DROP TABLE IF EXISTS `migration_marker`;
DROP TABLE IF EXISTS `community_messages`;
DROP TABLE IF EXISTS `user_balance_logs`;
DROP TABLE IF EXISTS `cms_ai_report`;
DROP TABLE IF EXISTS `workorder_logs`;
DROP TABLE IF EXISTS `workorders`;
DROP TABLE IF EXISTS `complaints`;
DROP TABLE IF EXISTS `repairs`;
DROP TABLE IF EXISTS `property_fee_payments`;
DROP TABLE IF EXISTS `property_fees`;
DROP TABLE IF EXISTS `user_parking_bindings`;
DROP TABLE IF EXISTS `parking_spaces`;
DROP TABLE IF EXISTS `notice_view_logs`;
DROP TABLE IF EXISTS `notices`;
DROP TABLE IF EXISTS `visitors`;
DROP TABLE IF EXISTS `product_view_logs`;
DROP TABLE IF EXISTS `payment_records`;
DROP TABLE IF EXISTS `wallet_transactions`;
DROP TABLE IF EXISTS `wallets`;
DROP TABLE IF EXISTS `pms_favorite`;
DROP TABLE IF EXISTS `pms_product_comment`;
DROP TABLE IF EXISTS `pms_store_product`;
DROP TABLE IF EXISTS `pms_store`;
DROP TABLE IF EXISTS `oms_order_item`;
DROP TABLE IF EXISTS `oms_order`;
DROP TABLE IF EXISTS `oms_cart`;
DROP TABLE IF EXISTS `promotion_products`;
DROP TABLE IF EXISTS `pms_promotion_product`;
DROP TABLE IF EXISTS `pms_promotion`;
DROP TABLE IF EXISTS `pms_product`;
DROP TABLE IF EXISTS `pms_product_category`;
DROP TABLE IF EXISTS `service_areas`;
DROP TABLE IF EXISTS `password_reset_codes`;
DROP TABLE IF EXISTS `admin_login_logs`;
DROP TABLE IF EXISTS `user_login_logs`;
DROP TABLE IF EXISTS `sys_role_permission`;
DROP TABLE IF EXISTS `sys_role_menu`;
DROP TABLE IF EXISTS `sys_user_role`;
DROP TABLE IF EXISTS `sys_permission`;
DROP TABLE IF EXISTS `sys_menu`;
DROP TABLE IF EXISTS `sys_role`;
DROP TABLE IF EXISTS `sys_user`;

SET FOREIGN_KEY_CHECKS = 1;

-- ── 2. 创建所有表 ──

-- 用户与权限
CREATE TABLE `sys_user` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `username` VARCHAR(64) NOT NULL DEFAULT '',
  `password` VARCHAR(255) NOT NULL DEFAULT '',
  `real_name` VARCHAR(64) NOT NULL DEFAULT '',
  `mobile` VARCHAR(20) NOT NULL,
  `age` INT NOT NULL DEFAULT 0,
  `gender` INT NOT NULL DEFAULT 0,
  `email` VARCHAR(128) NOT NULL DEFAULT '',
  `avatar` VARCHAR(255) NOT NULL DEFAULT '',
  `green_points` INT NOT NULL DEFAULT 0,
  `balance` DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  `face_registered` TINYINT(1) NOT NULL DEFAULT 0,
  `face_image_url` VARCHAR(512) NOT NULL DEFAULT '',
  `role` VARCHAR(32) NOT NULL DEFAULT '',
  `status` INT NOT NULL DEFAULT 1,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY `uk_sys_user_mobile` (`mobile`),
  KEY `idx_sys_user_role` (`role`),
  KEY `idx_sys_user_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `sys_role` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL DEFAULT '',
  `code` VARCHAR(64) NOT NULL,
  `remark` VARCHAR(255) NOT NULL DEFAULT '',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY `uk_sys_role_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `sys_menu` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `parent_id` BIGINT NOT NULL DEFAULT 0,
  `name` VARCHAR(64) NOT NULL DEFAULT '',
  `path` VARCHAR(255) NOT NULL DEFAULT '',
  `component` VARCHAR(255) NOT NULL DEFAULT '',
  `sort` INT NOT NULL DEFAULT 0,
  `type` INT NOT NULL DEFAULT 1,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY `idx_sys_menu_parent_sort` (`parent_id`, `sort`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `sys_user_role` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `role_id` BIGINT NOT NULL,
  UNIQUE KEY `uk_sys_user_role` (`user_id`, `role_id`),
  KEY `idx_sys_user_role_role` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `sys_role_menu` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `role_id` BIGINT NOT NULL,
  `menu_id` BIGINT NOT NULL,
  UNIQUE KEY `uk_sys_role_menu` (`role_id`, `menu_id`),
  KEY `idx_sys_role_menu_menu` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `sys_permission` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `code` VARCHAR(128) NOT NULL,
  `name` VARCHAR(128) NOT NULL DEFAULT '',
  `resource` VARCHAR(64) NOT NULL DEFAULT '',
  `method` VARCHAR(16) NOT NULL DEFAULT '',
  `path` VARCHAR(255) NOT NULL DEFAULT '',
  `type` INT NOT NULL DEFAULT 1,
  `status` INT NOT NULL DEFAULT 1,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY `uk_sys_permission_code` (`code`),
  KEY `idx_sys_permission_resource` (`resource`),
  KEY `idx_sys_permission_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `sys_role_permission` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `role_id` BIGINT NOT NULL,
  `permission_id` BIGINT NOT NULL,
  UNIQUE KEY `uk_role_permission` (`role_id`, `permission_id`),
  KEY `idx_sys_role_permission_permission` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `user_login_logs` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `mobile` VARCHAR(20) NOT NULL DEFAULT '',
  `login_time` DATETIME NOT NULL,
  `ip` VARCHAR(64) NOT NULL DEFAULT '',
  `user_agent` VARCHAR(512) NOT NULL DEFAULT '',
  `success` TINYINT(1) NOT NULL DEFAULT 0,
  `failure_reason` VARCHAR(255) NOT NULL DEFAULT '',
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_login_time` (`login_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `admin_login_logs` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `admin_user_id` BIGINT NOT NULL,
  `mobile` VARCHAR(20) NOT NULL DEFAULT '',
  `login_time` DATETIME NOT NULL,
  `ip` VARCHAR(64) NOT NULL DEFAULT '',
  `user_agent` VARCHAR(512) NOT NULL DEFAULT '',
  `success` TINYINT(1) NOT NULL DEFAULT 0,
  `failure_reason` VARCHAR(255) NOT NULL DEFAULT '',
  INDEX `idx_admin_user_id` (`admin_user_id`),
  INDEX `idx_login_time` (`login_time`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `password_reset_codes` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `mobile` VARCHAR(20) NOT NULL,
  `code_hash` VARCHAR(255) NOT NULL DEFAULT '',
  `expires_at` DATETIME NOT NULL,
  `used_at` DATETIME DEFAULT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_mobile` (`mobile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 商城
CREATE TABLE `pms_product_category` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL DEFAULT '',
  `icon` VARCHAR(255) NOT NULL DEFAULT '',
  `sort` INT NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `pms_product` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `category_name` VARCHAR(64) NOT NULL DEFAULT '',
  `name` VARCHAR(128) NOT NULL DEFAULT '',
  `description` TEXT,
  `price` BIGINT NOT NULL DEFAULT 0 COMMENT '售价（分）',
  `original_price` BIGINT NOT NULL DEFAULT 0 COMMENT '原价（分）',
  `stock` INT NOT NULL DEFAULT 0,
  `image_url` VARCHAR(255) NOT NULL DEFAULT '',
  `is_promotion` TINYINT NOT NULL DEFAULT 0,
  `sales` INT NOT NULL DEFAULT 0,
  `status` TINYINT NOT NULL DEFAULT 1,
  `version` INT NOT NULL DEFAULT 0 COMMENT '乐观锁版本',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `category_id` BIGINT NOT NULL DEFAULT 0,
  INDEX `idx_category_id` (`category_id`),
  INDEX `idx_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `pms_promotion` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `title` VARCHAR(128) NOT NULL DEFAULT '',
  `type` TINYINT NOT NULL DEFAULT 0,
  `start_date` DATETIME NOT NULL,
  `end_date` DATETIME NOT NULL,
  `status` TINYINT NOT NULL DEFAULT 1
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `pms_promotion_product` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `promotion_id` BIGINT NOT NULL,
  `product_id` BIGINT NOT NULL,
  INDEX `idx_promotion_id` (`promotion_id`),
  INDEX `idx_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `pms_product_comment` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `product_id` BIGINT NOT NULL,
  `content` TEXT NOT NULL,
  `rating` INT NOT NULL DEFAULT 5,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `oms_cart` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `product_id` BIGINT NOT NULL,
  `quantity` INT NOT NULL DEFAULT 1,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_user_id` (`user_id`),
  UNIQUE INDEX `uk_user_product` (`user_id`, `product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `oms_order` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `order_no` VARCHAR(64) NOT NULL DEFAULT '',
  `user_id` BIGINT NOT NULL,
  `store_id` BIGINT NOT NULL DEFAULT 0,
  `total_amount` BIGINT NOT NULL DEFAULT 0 COMMENT '订单总额（分）',
  `used_points` INT NOT NULL DEFAULT 0,
  `used_balance` BIGINT NOT NULL DEFAULT 0 COMMENT '使用余额（分）',
  `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0=pending_payment,1=paid,2=shipped,3=completed,40=cancelled',
  `expire_at` DATETIME DEFAULT NULL,
  `cancel_reason` VARCHAR(255) DEFAULT '',
  `cancelled_at` DATETIME DEFAULT NULL,
  `paid_at` DATETIME DEFAULT NULL,
  `version` INT NOT NULL DEFAULT 0,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_order_no` (`order_no`),
  INDEX `idx_status` (`status`),
  INDEX `idx_expire_at` (`expire_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `oms_order_item` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `order_id` BIGINT NOT NULL,
  `product_id` BIGINT NOT NULL,
  `price` BIGINT NOT NULL DEFAULT 0 COMMENT '单价（分）',
  `quantity` INT NOT NULL DEFAULT 1,
  `product_snapshot` VARCHAR(512) DEFAULT '',
  INDEX `idx_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `pms_store` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL DEFAULT '',
  `address` VARCHAR(255) NOT NULL DEFAULT '',
  `phone` VARCHAR(32) NOT NULL DEFAULT '',
  `area_id` BIGINT NOT NULL DEFAULT 0,
  `region` VARCHAR(128) NOT NULL DEFAULT '',
  `business_hours` VARCHAR(64) NOT NULL DEFAULT '',
  INDEX `idx_area_id` (`area_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `pms_store_product` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `store_id` BIGINT NOT NULL,
  `product_id` BIGINT NOT NULL,
  `stock` INT NOT NULL DEFAULT 0,
  `locked_stock` INT NOT NULL DEFAULT 0,
  `sold_count` INT NOT NULL DEFAULT 0,
  `version` INT NOT NULL DEFAULT 0,
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '1=on_shelf,0=off_shelf',
  UNIQUE INDEX `uk_store_product` (`store_id`, `product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `pms_favorite` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `product_id` BIGINT NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE INDEX `uk_user_product` (`user_id`, `product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `wallets` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `balance` BIGINT NOT NULL DEFAULT 0 COMMENT '余额（分）',
  `version` INT NOT NULL DEFAULT 0,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE INDEX `uk_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `wallet_transactions` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `type` TINYINT NOT NULL COMMENT '1=order_payment,2=transfer,3=recharge,4=refund,5=property_fee',
  `amount` BIGINT NOT NULL DEFAULT 0 COMMENT '金额（分）',
  `balance_before` BIGINT NOT NULL DEFAULT 0,
  `balance_after` BIGINT NOT NULL DEFAULT 0,
  `related_id` BIGINT NOT NULL DEFAULT 0,
  `remark` VARCHAR(255) NOT NULL DEFAULT '',
  `biz_type` VARCHAR(32) DEFAULT '',
  `biz_id` VARCHAR(64) DEFAULT '',
  `idempotency_key` VARCHAR(64) DEFAULT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_type` (`type`),
  UNIQUE INDEX `uk_idempotency_key` (`idempotency_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `payment_records` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `order_id` BIGINT NOT NULL,
  `order_no` VARCHAR(64) NOT NULL,
  `user_id` BIGINT NOT NULL,
  `amount` BIGINT NOT NULL DEFAULT 0,
  `payment_method` VARCHAR(32) NOT NULL DEFAULT 'wallet',
  `status` INT NOT NULL DEFAULT 0,
  `idempotency_key` VARCHAR(64) NOT NULL,
  `fail_reason` VARCHAR(255) DEFAULT '',
  `paid_at` DATETIME DEFAULT NULL,
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY `uk_idempotency_key` (`idempotency_key`),
  INDEX `idx_order_id` (`order_id`),
  INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `product_view_logs` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `product_id` BIGINT NOT NULL,
  `user_id` BIGINT NOT NULL DEFAULT 0,
  `ip` VARCHAR(64) NOT NULL DEFAULT '',
  `user_agent` VARCHAR(512) NOT NULL DEFAULT '',
  `viewed_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_product_id` (`product_id`),
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_viewed_at` (`viewed_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `service_areas` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL DEFAULT '',
  `sort` INT NOT NULL DEFAULT 0,
  `status` TINYINT NOT NULL DEFAULT 1
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- 社区
CREATE TABLE `notices` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `title` VARCHAR(100) NOT NULL,
  `content` TEXT NOT NULL,
  `publisher` VARCHAR(50) NOT NULL DEFAULT '',
  `view_count` BIGINT NOT NULL DEFAULT 0,
  `status` TINYINT NOT NULL DEFAULT 1,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_notices_status_created` (`status`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `notice_view_logs` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `notice_id` BIGINT NOT NULL,
  `user_id` BIGINT NOT NULL,
  `viewed_at` DATETIME NOT NULL,
  `read_at` DATETIME NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY `idx_notice_user` (`notice_id`, `user_id`),
  KEY `idx_notice_view_logs_user` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `visitors` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `visitor_name` VARCHAR(50) NOT NULL,
  `visitor_phone` VARCHAR(20) NOT NULL,
  `visit_purpose` VARCHAR(255) NOT NULL,
  `release_time` DATETIME NOT NULL,
  `valid_date` DATE NOT NULL,
  `status` TINYINT NOT NULL DEFAULT 0,
  `audit_remark` VARCHAR(255) NOT NULL DEFAULT '',
  `audit_at` DATETIME NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_visitors_user` (`user_id`),
  KEY `idx_visitors_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `parking_spaces` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `parking_no` VARCHAR(50) NOT NULL,
  `status` TINYINT NOT NULL DEFAULT 0,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY `uk_parking_no` (`parking_no`),
  KEY `idx_parking_spaces_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `user_parking_bindings` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `parking_space_id` BIGINT NOT NULL,
  `car_plate` VARCHAR(20) NOT NULL DEFAULT '',
  `status` TINYINT NOT NULL DEFAULT 1,
  `bound_at` DATETIME NOT NULL,
  `unbound_at` DATETIME NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_user_parking_user` (`user_id`, `status`),
  KEY `idx_user_parking_space` (`parking_space_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `property_fees` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `month` VARCHAR(20) NOT NULL,
  `amount` BIGINT NOT NULL,
  `status` TINYINT NOT NULL DEFAULT 0,
  `due_date` DATE NULL,
  `paid_at` DATETIME NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_property_fees_user` (`user_id`, `status`),
  KEY `idx_property_fees_month` (`month`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `property_fee_payments` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `property_fee_id` BIGINT NOT NULL,
  `user_id` BIGINT NOT NULL,
  `amount` BIGINT NOT NULL,
  `wallet_transaction_id` BIGINT NOT NULL DEFAULT 0,
  `idempotency_key` VARCHAR(64) DEFAULT NULL,
  `status` TINYINT NOT NULL DEFAULT 1,
  `paid_at` DATETIME NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY `uk_property_payment_user_idempotency` (`user_id`, `idempotency_key`),
  KEY `idx_property_payments_user` (`user_id`),
  KEY `idx_property_payments_fee` (`property_fee_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- 工单
CREATE TABLE `workorders` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `type` VARCHAR(20) NOT NULL COMMENT 'repair=报事维修, complaint=事项投诉',
  `user_id` BIGINT NOT NULL,
  `category` VARCHAR(50) NOT NULL,
  `description` TEXT NOT NULL,
  `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0=pending,1=processing,2=completed',
  `result` VARCHAR(500) NOT NULL DEFAULT '',
  `processor_id` BIGINT NOT NULL DEFAULT 0,
  `processed_at` DATETIME NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_workorders_user` (`user_id`),
  KEY `idx_workorders_type_status` (`type`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE `workorder_logs` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `target_type` VARCHAR(20) NOT NULL,
  `target_id` BIGINT NOT NULL,
  `from_status` TINYINT NOT NULL DEFAULT -1,
  `to_status` TINYINT NOT NULL,
  `operator_id` BIGINT NOT NULL DEFAULT 0,
  `action` VARCHAR(50) NOT NULL,
  `remark` VARCHAR(500) NOT NULL DEFAULT '',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  KEY `idx_workorder_logs_target` (`target_type`, `target_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- AI 报告
CREATE TABLE `cms_ai_report` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `repair_new_count` BIGINT NOT NULL DEFAULT 0,
  `repair_pending_count` BIGINT NOT NULL DEFAULT 0,
  `visitor_new_count` BIGINT NOT NULL DEFAULT 0,
  `property_paid_count` BIGINT NOT NULL DEFAULT 0,
  `property_paid_amount` DECIMAL(10,2) NOT NULL DEFAULT 0.00,
  `report_summary` VARCHAR(255) DEFAULT '',
  `report_markdown` LONGTEXT,
  `generated_by` BIGINT NOT NULL DEFAULT 0,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE `community_messages` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `content` TEXT NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ── 3. 种子数据 ──

-- 角色
INSERT INTO `sys_role` (`id`, `name`, `code`, `remark`) VALUES
(1, '系统管理员', 'admin', '全局管理权限'),
(2, '物业管理员', 'property', '物业与报修管理'),
(3, '门店管理员', 'store', '门店与商品管理'),
(4, '普通用户', 'user', '居民用户');

-- 菜单
INSERT INTO `sys_menu` (`id`, `parent_id`, `name`, `path`, `component`, `sort`, `type`) VALUES
(1, 0, '仪表盘', '/admin/dashboard', 'views/admin/Dashboard.vue', 1, 1),
(2, 0, '订单管理', '/admin/order', 'views/admin/order/OrderList.vue', 2, 1),
(3, 0, '物业管理', '/admin/property', 'views/admin/property/PropertyFeeList.vue', 3, 1),
(4, 0, 'AI报表', '/admin/ai-report', 'views/admin/AIReport.vue', 4, 1),
(5, 0, '权限管理', '/admin/rbac', 'views/admin/rbac/RbacLayout.vue', 5, 1),
(6, 5, '角色管理', '/admin/rbac/roles', 'views/admin/rbac/RoleList.vue', 1, 2),
(7, 5, '用户管理', '/admin/rbac/users', 'views/admin/rbac/UserList.vue', 2, 2),
(8, 5, '会员管理', '/admin/rbac/members', 'views/admin/rbac/MemberList.vue', 3, 2),
(9, 0, '商品管理', '/admin/mall', 'views/admin/mall/MallLayout.vue', 6, 1),
(10, 9, '商品列表', '/admin/mall/products', 'views/admin/mall/ProductList.vue', 1, 2),
(11, 9, '分类管理', '/admin/mall/categories', 'views/admin/mall/CategoryList.vue', 2, 2),
(12, 9, '营销管理', '/admin/mall/promotions', 'views/admin/mall/PromotionList.vue', 3, 2),
(13, 0, '门店管理', '/admin/store', 'views/admin/store/StoreLayout.vue', 7, 1),
(14, 13, '门店列表', '/admin/store/list', 'views/admin/store/StoreList.vue', 1, 2),
(15, 13, '服务区域', '/admin/store/areas', 'views/admin/store/ServiceAreaList.vue', 2, 2),
(16, 0, '社区管理', '/admin/community', 'views/admin/community/CommunityLayout.vue', 8, 1),
(17, 16, '公告管理', '/admin/community/notices', 'views/admin/community/NoticeList.vue', 1, 2),
(18, 16, '访客管理', '/admin/community/visitors', 'views/admin/community/VisitorList.vue', 2, 2),
(19, 16, '车位管理', '/admin/community/parking', 'views/admin/community/ParkingList.vue', 3, 2),
(20, 0, '工单管理', '/admin/workorder', 'views/admin/workorder/WorkorderLayout.vue', 9, 1),
(21, 20, '报修投诉', '/admin/repairs', 'views/admin/RepairList.vue', 1, 2),
(23, 0, '系统管理', '/admin/system', 'views/admin/system/SystemLayout.vue', 10, 1),
(24, 23, '用户登录日志', '/admin/system/user-login-logs', 'views/admin/system/UserLoginLogs.vue', 1, 2),
(25, 23, '管理员登录日志', '/admin/system/admin-login-logs', 'views/admin/system/AdminLoginLogs.vue', 2, 2);

-- 权限
INSERT INTO `sys_permission` (`id`, `code`, `name`, `resource`, `method`, `path`, `type`, `status`) VALUES
(1,  'rbac:role:create',          '创建角色',       'role',       'POST',   '/api/admin/roles',              2, 1),
(2,  'rbac:role:update',          '更新角色',       'role',       'PUT',    '/api/admin/roles',              2, 1),
(3,  'rbac:role:delete',          '删除角色',       'role',       'DELETE', '/api/admin/roles',              2, 1),
(4,  'rbac:role:list',            '查询角色列表',    'role',       'GET',    '/api/admin/roles',              1, 1),
(5,  'rbac:role:bind_menu',       '角色绑定菜单',    'role',       'POST',   '/api/admin/roles/:id/menus',    2, 1),
(6,  'rbac:role:bind_permission', '角色绑定权限',    'role',       'POST',   '/api/admin/roles/:id/permissions', 2, 1),
(7,  'rbac:role:get_permissions', '查询角色权限',    'role',       'GET',    '/api/admin/roles/:id/permissions', 1, 1),
(8,  'rbac:user:list',            '查询管理员列表',  'user',       'GET',    '/api/admin/users',              1, 1),
(9,  'rbac:user:freeze',          '冻结/解冻用户',  'user',       'POST',   '/api/admin/users/freeze',       2, 1),
(10, 'rbac:user:assign_role',     '分配用户角色',    'user',       'POST',   '/api/admin/users/assign-role',  2, 1),
(11, 'rbac:user:assign_roles',    '分配用户多角色',  'user',       'POST',   '/api/admin/users/:id/roles',    2, 1),
(12, 'rbac:user:get_roles',       '查询用户角色',    'user',       'GET',    '/api/admin/users/:id/roles',    1, 1),
(13, 'rbac:member:list',          '查询会员列表',    'member',     'GET',    '/api/admin/members',            1, 1),
(14, 'rbac:permission:list',      '查询权限列表',    'permission', 'GET',    '/api/admin/permissions',        1, 1),
(15, 'rbac:menu:list',            '查询菜单列表',    'menu',       'GET',    '/api/admin/menus',              1, 1),
(16, 'log:user_login:list',       '查询用户登录日志', 'log',       'GET',    '/api/admin/user-login-logs',    1, 1),
(17, 'log:admin_login:list',      '查询管理员登录日志', 'log',     'GET',    '/api/admin/admin-login-logs',   1, 1),
(18, 'rbac:user:update_balance',  '调整用户余额',    'user',       'POST',   '/api/admin/users/update-balance', 2, 1),
(19, 'mall:product:list',         '查询商品列表',     'mall_product',     'GET',    '/api/admin/mall/products',             1, 1),
(20, 'mall:category:create',      '创建商品分类',     'mall_category',    'POST',   '/api/admin/mall/categories',           2, 1),
(21, 'mall:category:update',      '更新商品分类',     'mall_category',    'PUT',    '/api/admin/mall/categories/:id',       2, 1),
(22, 'mall:category:delete',      '删除商品分类',     'mall_category',    'DELETE', '/api/admin/mall/categories/:id',       2, 1),
(23, 'mall:product:create',       '创建商品',        'mall_product',     'POST',   '/api/admin/mall/products',             2, 1),
(24, 'mall:product:update',       '更新商品',        'mall_product',     'PUT',    '/api/admin/mall/products/:id',         2, 1),
(25, 'mall:product:delete',       '删除商品',        'mall_product',     'DELETE', '/api/admin/mall/products/:id',         2, 1),
(26, 'mall:promotion:create',     '创建促销',        'mall_promotion',   'POST',   '/api/admin/mall/promotions',           2, 1),
(27, 'mall:promotion:update',     '更新促销',        'mall_promotion',   'PUT',    '/api/admin/mall/promotions/:id',       2, 1),
(28, 'mall:promotion:delete',     '删除促销',        'mall_promotion',   'DELETE', '/api/admin/mall/promotions/:id',       2, 1),
(29, 'mall:promotion:bind_product','促销绑定商品',    'mall_promotion',   'POST',   '/api/admin/mall/promotions/:id/products', 2, 1),
(30, 'mall:service_area:create',  '创建服务区域',     'mall_service_area','POST',   '/api/admin/mall/service-areas',        2, 1),
(31, 'mall:service_area:update',  '更新服务区域',     'mall_service_area','PUT',    '/api/admin/mall/service-areas/:id',    2, 1),
(32, 'mall:service_area:delete',  '删除服务区域',     'mall_service_area','DELETE', '/api/admin/mall/service-areas/:id',    2, 1),
(33, 'mall:store:create',         '创建门店',        'mall_store',       'POST',   '/api/admin/mall/stores',               2, 1),
(34, 'mall:store:update',         '更新门店',        'mall_store',       'PUT',    '/api/admin/mall/stores/:id',           2, 1),
(35, 'mall:store:delete',         '删除门店',        'mall_store',       'DELETE', '/api/admin/mall/stores/:id',           2, 1),
(36, 'mall:store_product:bind',   '绑定门店商品',     'mall_store_product','POST',  '/api/admin/mall/store-products',       2, 1),
(37, 'mall:store_product:unbind', '解绑门店商品',     'mall_store_product','DELETE','/api/admin/mall/store-products',       2, 1),
(38, 'mall:store_product:status', '上下架门店商品',   'mall_store_product','PUT',   '/api/admin/mall/store-products/status',2, 1),
(39, 'mall:store_product:stock',  '门店商品库存',     'mall_store_product','PUT',   '/api/admin/mall/store-products/stock', 2, 1),
(40, 'mall:store_product:list',   '查询门店商品',     'mall_store_product','GET',   '/api/admin/mall/store-products/:store_id', 1, 1),
(41, 'mall:order:list',           '查询订单',        'mall_order',       'GET',    '/api/admin/mall/orders',               1, 1),
(42, 'mall:order:ship',           '订单发货',        'mall_order',       'POST',   '/api/admin/mall/orders/:id/ship',      2, 1),
(43, 'mall:order:cancel',         '订单作废',        'mall_order',       'POST',   '/api/admin/mall/orders/:id/cancel',    2, 1),
(50, 'community:notice:list',          '查询公告管理列表', 'community_notice',  'GET',    '/api/admin/community/notices',                  1, 1),
(51, 'community:notice:create',        '发布公告',        'community_notice',  'POST',   '/api/admin/community/notices',                  2, 1),
(52, 'community:notice:delete',        '删除公告',        'community_notice',  'DELETE', '/api/admin/community/notices/:id',              2, 1),
(53, 'community:notice:views',         '查询公告浏览状态', 'community_notice',  'GET',    '/api/admin/community/notices/:id/views',        1, 1),
(54, 'community:visitor:list',         '查询访客记录',     'community_visitor', 'GET',    '/api/admin/community/visitors',                 1, 1),
(55, 'community:visitor:audit',        '访客审核放行',     'community_visitor', 'POST',   '/api/admin/community/visitors/:id/audit',       2, 1),
(56, 'community:parking:list',         '查询车位列表',     'community_parking', 'GET',    '/api/admin/community/parking-spaces',           1, 1),
(57, 'community:parking:create',       '创建车位',        'community_parking', 'POST',   '/api/admin/community/parking-spaces',           2, 1),
(58, 'community:parking:assign',       '分配车位',        'community_parking', 'POST',   '/api/admin/community/parking-spaces/:id/assign',2, 1),
(59, 'community:parking:statistics',   '查询车位统计',     'community_parking', 'GET',    '/api/admin/community/parking-spaces/statistics',1, 1),
(60, 'community:fee:list',             '查询物业费',       'community_fee',     'GET',    '/api/admin/community/property-fees',            1, 1),
(61, 'community:fee:create',           '创建物业费',       'community_fee',     'POST',   '/api/admin/community/property-fees',            2, 1),
(62, 'community:fee:payment_list',     '查询缴费记录',     'community_fee',     'GET',    '/api/admin/community/property-fees/payments',   1, 1),
(70, 'workorder:repair:list',          '查询报修投诉列表', 'workorder',           'GET',  '/api/admin/workorders',                        1, 1),
(71, 'workorder:repair:process',       '处理报修投诉',     'workorder',           'POST', '/api/admin/workorders/:id/process',            2, 1),
(80, 'statistics:product:sales_rank',  '商品销售排行',     'statistics_product',  'GET',  '/api/statistics/products/sales-rank',           1, 1),
(81, 'statistics:product:view_rank',   '商品访客排行',     'statistics_product',  'GET',  '/api/statistics/products/view-rank',            1, 1),
(82, 'statistics:community:overview',  '社区运营概览',     'statistics_community','GET',  '/api/statistics/community/overview',            1, 1),
(83, 'statistics:order:summary',       '订单统计',        'statistics_order',    'GET',  '/api/statistics/orders',                        1, 1),
(84, 'statistics:workorder:summary',   '报修投诉统计',     'statistics_workorder','GET',  '/api/statistics/workorders',                    1, 1),
(85, 'statistics:ai_report:generate',  '生成AI报告',      'statistics_report',   'POST', '/api/statistics/ai-report/generate',            2, 1),
(86, 'statistics:ai_report:read',      '查看AI报告',      'statistics_report',   'GET',  '/api/statistics/ai-report',                     1, 1);

-- 管理员角色拥有所有权限
INSERT INTO `sys_role_permission` (`role_id`, `permission_id`)
SELECT 1, `id` FROM `sys_permission`;

-- 物业角色: 社区 + 工单权限
INSERT INTO `sys_role_permission` (`role_id`, `permission_id`) VALUES
(2, 4), (2, 13), (2, 16), (2, 17),
(2, 50), (2, 51), (2, 52), (2, 53), (2, 54), (2, 55), (2, 56), (2, 57), (2, 58), (2, 59), (2, 60), (2, 61), (2, 62),
(2, 70), (2, 71);

-- 门店角色: 商城权限
INSERT INTO `sys_role_permission` (`role_id`, `permission_id`) VALUES
(3, 19), (3, 20), (3, 21), (3, 22), (3, 23), (3, 24), (3, 25), (3, 26), (3, 27), (3, 28), (3, 29),
(3, 30), (3, 31), (3, 32), (3, 33), (3, 34), (3, 35), (3, 36), (3, 37), (3, 38), (3, 39),
(3, 40), (3, 41), (3, 42), (3, 43);

-- 菜单绑定: 管理员看所有
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`)
SELECT 1, `id` FROM `sys_menu`;

-- 物业看: 仪表盘、物业、社区、工单、系统
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES
(2, 1), (2, 3), (2, 16), (2, 17), (2, 18), (2, 19), (2, 20), (2, 21), (2, 23), (2, 24), (2, 25);

-- 门店看: 仪表盘、订单、商品、门店
INSERT INTO `sys_role_menu` (`role_id`, `menu_id`) VALUES
(3, 1), (3, 2), (3, 9), (3, 10), (3, 11), (3, 12), (3, 13), (3, 14), (3, 15);

-- 管理员用户 (密码: sdl@admin)
INSERT INTO `sys_user` (`id`, `username`, `password`, `real_name`, `mobile`, `age`, `gender`, `email`, `avatar`, `green_points`, `balance`, `role`, `status`) VALUES
(1, 'admin', '$2a$14$7Q2NpdUbkHHlfrfUzK7pq.kyuvvW4M8Va5Abu7dFEWnvZr63wqgPu', '系统管理员', '13483000001', 30, 1, 'admin@community.com', 'https://cube.elemecdn.com/3/7c/3ea6beec64369c2642b92c6726f1epng.png', 5000, 8000.00, 'admin', 1);

-- 管理员用户-角色绑定
INSERT INTO `sys_user_role` (`user_id`, `role_id`) VALUES (1, 1);

-- 管理员钱包
INSERT INTO `wallets` (`user_id`, `balance`) VALUES (1, 100000);

-- 默认公告
INSERT INTO `notices` (`id`, `title`, `content`, `publisher`, `status`) VALUES
(1, '社区服务上线通知', '智慧社区微服务框架已接入公告、访客、车位与物业费基础接口。', '系统管理员', 1);

-- 默认车位
INSERT INTO `parking_spaces` (`id`, `parking_no`, `status`) VALUES
(1, 'A-001', 0), (2, 'A-002', 0), (3, 'B-001', 0);

-- 商品分类
INSERT INTO `pms_product_category` (`name`, `icon`, `sort`) VALUES
('生鲜果蔬', '🍎', 1), ('粮油副食', '🌾', 2), ('日用百货', '🧴', 3), ('家居清洁', '🧹', 4), ('个人护理', '🧴', 5);

-- 服务区域
INSERT INTO `service_areas` (`name`, `sort`, `status`) VALUES
('东软智慧社区A区', 1, 1), ('东软智慧社区B区', 2, 1), ('东软智慧社区C区', 3, 1);

-- 门店
INSERT INTO `pms_store` (`name`, `address`, `phone`, `area_id`, `region`, `business_hours`) VALUES
('社区便民超市A', 'A区1号楼底商', '024-12345678', 1, 'A区', '08:00-22:00'),
('社区便民超市B', 'B区2号楼底商', '024-87654321', 2, 'B区', '08:00-22:00');

-- 示例商品
INSERT INTO `pms_product` (`category_name`, `name`, `description`, `price`, `original_price`, `stock`, `image_url`, `is_promotion`, `sales`, `status`, `category_id`) VALUES
('生鲜果蔬', '新鲜苹果 500g', '产地直供，新鲜脆甜', 1290, 1590, 200, '', 1, 50, 1, 1),
('生鲜果蔬', '有机西红柿 500g', '自然成熟，口感沙甜', 850, 1000, 150, '', 0, 30, 1, 1),
('粮油副食', '东北大米 5kg', '五常稻花香，粒粒饱满', 3990, 4990, 100, '', 1, 80, 1, 2),
('日用百货', '抽纸 3层120抽*10包', '柔软亲肤，不掉屑', 2990, 3500, 300, '', 0, 120, 1, 3),
('家居清洁', '洗洁精 1.5kg', '去油不伤手', 1590, 1890, 200, '', 0, 60, 1, 4),
('个人护理', '牙膏 120g', '清新口气，防蛀固齿', 990, 1200, 250, '', 0, 90, 1, 5);

-- 门店商品绑定
INSERT INTO `pms_store_product` (`store_id`, `product_id`, `stock`, `status`) VALUES
(1, 1, 50, 1), (1, 2, 30, 1), (1, 3, 20, 1), (1, 4, 100, 1),
(2, 1, 40, 1), (2, 3, 25, 1), (2, 5, 50, 1), (2, 6, 80, 1);
