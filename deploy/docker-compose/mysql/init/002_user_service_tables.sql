USE smart_community;

-- AUTH/RBAC: users, roles, menus and API permissions
CREATE TABLE IF NOT EXISTS `sys_user` (
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

CREATE TABLE IF NOT EXISTS `sys_role` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `name` VARCHAR(64) NOT NULL DEFAULT '',
  `code` VARCHAR(64) NOT NULL,
  `remark` VARCHAR(255) NOT NULL DEFAULT '',
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  UNIQUE KEY `uk_sys_role_code` (`code`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `sys_menu` (
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

CREATE TABLE IF NOT EXISTS `sys_user_role` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `user_id` BIGINT NOT NULL,
  `role_id` BIGINT NOT NULL,
  UNIQUE KEY `uk_sys_user_role` (`user_id`, `role_id`),
  KEY `idx_sys_user_role_role` (`role_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `sys_role_menu` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `role_id` BIGINT NOT NULL,
  `menu_id` BIGINT NOT NULL,
  UNIQUE KEY `uk_sys_role_menu` (`role_id`, `menu_id`),
  KEY `idx_sys_role_menu_menu` (`menu_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS `sys_permission` (
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

CREATE TABLE IF NOT EXISTS `sys_role_permission` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `role_id` BIGINT NOT NULL,
  `permission_id` BIGINT NOT NULL,
  UNIQUE KEY `uk_role_permission` (`role_id`, `permission_id`),
  KEY `idx_sys_role_permission_permission` (`permission_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- LOG-001: User login logs
CREATE TABLE IF NOT EXISTS `user_login_logs` (
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

-- LOG-002: Admin login logs
CREATE TABLE IF NOT EXISTS `admin_login_logs` (
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

-- AUTH-003: Password reset codes (audit trail)
CREATE TABLE IF NOT EXISTS `password_reset_codes` (
  `id` BIGINT PRIMARY KEY AUTO_INCREMENT,
  `mobile` VARCHAR(20) NOT NULL,
  `code_hash` VARCHAR(255) NOT NULL DEFAULT '',
  `expires_at` DATETIME NOT NULL,
  `used_at` DATETIME DEFAULT NULL,
  `created_at` DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
  INDEX `idx_mobile` (`mobile`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;
