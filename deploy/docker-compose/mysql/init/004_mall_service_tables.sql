USE smart_community;

-- MALL-001~005: Product and category tables
CREATE TABLE IF NOT EXISTS `pms_product_category` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL DEFAULT '',
  `icon` VARCHAR(255) NOT NULL DEFAULT '',
  `sort` INT NOT NULL DEFAULT 0
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `pms_product` (
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

-- MALL-003/ADMIN-MALL-007: Promotion tables
CREATE TABLE IF NOT EXISTS `pms_promotion` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `title` VARCHAR(128) NOT NULL DEFAULT '',
  `type` TINYINT NOT NULL DEFAULT 0,
  `start_date` DATETIME NOT NULL,
  `end_date` DATETIME NOT NULL,
  `status` TINYINT NOT NULL DEFAULT 1
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `pms_promotion_product` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `promotion_id` BIGINT NOT NULL,
  `product_id` BIGINT NOT NULL,
  INDEX `idx_promotion_id` (`promotion_id`),
  INDEX `idx_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `pms_product_comment` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `product_id` BIGINT NOT NULL,
  `content` TEXT NOT NULL,
  `rating` INT NOT NULL DEFAULT 5,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_product_id` (`product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- MALL-007/010/011: Cart
CREATE TABLE IF NOT EXISTS `oms_cart` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `product_id` BIGINT NOT NULL,
  `quantity` INT NOT NULL DEFAULT 1,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_user_id` (`user_id`),
  UNIQUE INDEX `uk_user_product` (`user_id`, `product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- MALL-012~017/ADMIN-MALL-011: Order tables
CREATE TABLE IF NOT EXISTS `oms_order` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `order_no` VARCHAR(64) NOT NULL DEFAULT '',
  `user_id` BIGINT NOT NULL,
  `store_id` BIGINT NOT NULL DEFAULT 0,
  `total_amount` BIGINT NOT NULL DEFAULT 0 COMMENT '订单总额（分）',
  `used_points` INT NOT NULL DEFAULT 0,
  `used_balance` BIGINT NOT NULL DEFAULT 0 COMMENT '使用余额（分）',
  `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0=pending_payment,1=paid,2=shipped,3=completed,40=cancelled',
  `expire_at` DATETIME DEFAULT NULL COMMENT '订单过期时间',
  `cancel_reason` VARCHAR(255) DEFAULT '' COMMENT '取消原因',
  `cancelled_at` DATETIME DEFAULT NULL COMMENT '取消时间',
  `paid_at` DATETIME DEFAULT NULL COMMENT '支付时间',
  `version` INT NOT NULL DEFAULT 0 COMMENT '乐观锁版本',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_order_no` (`order_no`),
  INDEX `idx_status` (`status`),
  INDEX `idx_expire_at` (`expire_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `oms_order_item` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `order_id` BIGINT NOT NULL,
  `product_id` BIGINT NOT NULL,
  `price` BIGINT NOT NULL DEFAULT 0 COMMENT '单价（分）',
  `quantity` INT NOT NULL DEFAULT 1,
  `product_snapshot` VARCHAR(512) DEFAULT '' COMMENT '下单时商品快照',
  INDEX `idx_order_id` (`order_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- MALL-006/ADMIN-MALL-009: Store tables
CREATE TABLE IF NOT EXISTS `pms_store` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL DEFAULT '',
  `address` VARCHAR(255) NOT NULL DEFAULT '',
  `phone` VARCHAR(32) NOT NULL DEFAULT '',
  `area_id` BIGINT NOT NULL DEFAULT 0,
  `region` VARCHAR(128) NOT NULL DEFAULT '',
  `business_hours` VARCHAR(64) NOT NULL DEFAULT '',
  INDEX `idx_area_id` (`area_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ADMIN-MALL-010: Store-product binding
CREATE TABLE IF NOT EXISTS `pms_store_product` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `store_id` BIGINT NOT NULL,
  `product_id` BIGINT NOT NULL,
  `stock` INT NOT NULL DEFAULT 0,
  `locked_stock` INT NOT NULL DEFAULT 0 COMMENT '锁定库存',
  `sold_count` INT NOT NULL DEFAULT 0 COMMENT '已售数量',
  `version` INT NOT NULL DEFAULT 0 COMMENT '乐观锁版本',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '1=on_shelf,0=off_shelf',
  UNIQUE INDEX `uk_store_product` (`store_id`, `product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- MALL-008/009/018: Favorites
CREATE TABLE IF NOT EXISTS `pms_favorite` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `product_id` BIGINT NOT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE INDEX `uk_user_product` (`user_id`, `product_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- MALL-019~021: Wallet
CREATE TABLE IF NOT EXISTS `wallets` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `balance` BIGINT NOT NULL DEFAULT 0 COMMENT '余额（分）',
  `version` INT NOT NULL DEFAULT 0 COMMENT '乐观锁版本',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE INDEX `uk_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `wallet_transactions` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `type` TINYINT NOT NULL COMMENT '1=order_payment,2=transfer,3=recharge,4=refund,5=property_fee',
  `amount` BIGINT NOT NULL DEFAULT 0 COMMENT '金额（分）',
  `balance_before` BIGINT NOT NULL DEFAULT 0 COMMENT '变动前余额（分）',
  `balance_after` BIGINT NOT NULL DEFAULT 0 COMMENT '变动后余额（分）',
  `related_id` BIGINT NOT NULL DEFAULT 0,
  `remark` VARCHAR(255) NOT NULL DEFAULT '',
  `biz_type` VARCHAR(32) DEFAULT '' COMMENT '业务类型',
  `biz_id` VARCHAR(64) DEFAULT '' COMMENT '业务ID',
  `idempotency_key` VARCHAR(64) DEFAULT NULL COMMENT '幂等键',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_type` (`type`),
  UNIQUE INDEX `uk_idempotency_key` (`idempotency_key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Payment records (idempotent audit trail)
CREATE TABLE IF NOT EXISTS `payment_records` (
  `id` BIGINT AUTO_INCREMENT PRIMARY KEY,
  `order_id` BIGINT NOT NULL COMMENT '订单ID',
  `order_no` VARCHAR(64) NOT NULL COMMENT '订单号',
  `user_id` BIGINT NOT NULL COMMENT '用户ID',
  `amount` BIGINT NOT NULL DEFAULT 0 COMMENT '支付金额（分）',
  `payment_method` VARCHAR(32) NOT NULL DEFAULT 'wallet' COMMENT '支付方式',
  `status` INT NOT NULL DEFAULT 0 COMMENT '状态: 0=init, 1=success, 2=failed',
  `idempotency_key` VARCHAR(64) NOT NULL COMMENT '幂等键',
  `fail_reason` VARCHAR(255) DEFAULT '' COMMENT '失败原因',
  `paid_at` DATETIME DEFAULT NULL COMMENT '支付成功时间',
  `created_at` DATETIME DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY `uk_idempotency_key` (`idempotency_key`),
  INDEX `idx_order_id` (`order_id`),
  INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- STAT-002: Product view logs (used by statistics-service for visitor ranking)
CREATE TABLE IF NOT EXISTS `product_view_logs` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `product_id` BIGINT NOT NULL,
  `user_id` BIGINT NOT NULL DEFAULT 0 COMMENT '0=anonymous',
  `ip` VARCHAR(64) NOT NULL DEFAULT '',
  `user_agent` VARCHAR(512) NOT NULL DEFAULT '',
  `viewed_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_product_id` (`product_id`),
  INDEX `idx_user_id` (`user_id`),
  INDEX `idx_viewed_at` (`viewed_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- ADMIN-MALL-008: Service areas
CREATE TABLE IF NOT EXISTS `service_areas` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(128) NOT NULL DEFAULT '',
  `sort` INT NOT NULL DEFAULT 0,
  `status` TINYINT NOT NULL DEFAULT 1
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
