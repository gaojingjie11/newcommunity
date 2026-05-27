USE smart_community;

CREATE TABLE IF NOT EXISTS `notices` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `title` VARCHAR(100) NOT NULL,
  `content` TEXT NOT NULL,
  `publisher` VARCHAR(50) NOT NULL DEFAULT '',
  `view_count` BIGINT NOT NULL DEFAULT 0,
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '1=published,0=hidden/deleted',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_notices_status_created` (`status`, `created_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `notice_view_logs` (
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

CREATE TABLE IF NOT EXISTS `visitors` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `visitor_name` VARCHAR(50) NOT NULL,
  `visitor_phone` VARCHAR(20) NOT NULL,
  `visit_purpose` VARCHAR(255) NOT NULL,
  `release_time` DATETIME NOT NULL,
  `valid_date` DATE NOT NULL,
  `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0=pending,1=approved,2=rejected',
  `audit_remark` VARCHAR(255) NOT NULL DEFAULT '',
  `audit_at` DATETIME NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_visitors_user` (`user_id`),
  KEY `idx_visitors_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `parking_spaces` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `parking_no` VARCHAR(50) NOT NULL,
  `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0=free,1=bound',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY `uk_parking_no` (`parking_no`),
  KEY `idx_parking_spaces_status` (`status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `user_parking_bindings` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `parking_space_id` BIGINT NOT NULL,
  `car_plate` VARCHAR(20) NOT NULL DEFAULT '',
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '1=active,0=inactive',
  `bound_at` DATETIME NOT NULL,
  `unbound_at` DATETIME NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_user_parking_user` (`user_id`, `status`),
  KEY `idx_user_parking_space` (`parking_space_id`, `status`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `property_fees` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `month` VARCHAR(20) NOT NULL,
  `amount` BIGINT NOT NULL COMMENT 'amount in cents',
  `status` TINYINT NOT NULL DEFAULT 0 COMMENT '0=unpaid,1=paid',
  `due_date` DATE NULL,
  `paid_at` DATETIME NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  KEY `idx_property_fees_user` (`user_id`, `status`),
  KEY `idx_property_fees_month` (`month`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

CREATE TABLE IF NOT EXISTS `property_fee_payments` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `property_fee_id` BIGINT NOT NULL,
  `user_id` BIGINT NOT NULL,
  `amount` BIGINT NOT NULL COMMENT 'amount in cents',
  `wallet_transaction_id` BIGINT NOT NULL DEFAULT 0,
  `idempotency_key` VARCHAR(64) DEFAULT NULL,
  `status` TINYINT NOT NULL DEFAULT 1 COMMENT '1=success',
  `paid_at` DATETIME NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `updated_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
  UNIQUE KEY `uk_property_payment_user_idempotency` (`user_id`, `idempotency_key`),
  KEY `idx_property_payments_user` (`user_id`),
  KEY `idx_property_payments_fee` (`property_fee_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

INSERT INTO `notices` (`id`, `title`, `content`, `publisher`, `status`, `created_at`, `updated_at`) VALUES
(1, '社区服务上线通知', '智慧社区微服务框架已接入公告、访客、车位与物业费基础接口。', '系统管理员', 1, NOW(), NOW())
ON DUPLICATE KEY UPDATE `title`=VALUES(`title`), `content`=VALUES(`content`), `status`=VALUES(`status`);

INSERT INTO `parking_spaces` (`id`, `parking_no`, `status`, `created_at`, `updated_at`) VALUES
(1, 'A-001', 0, NOW(), NOW()),
(2, 'A-002', 0, NOW(), NOW()),
(3, 'B-001', 0, NOW(), NOW())
ON DUPLICATE KEY UPDATE `parking_no`=VALUES(`parking_no`);
