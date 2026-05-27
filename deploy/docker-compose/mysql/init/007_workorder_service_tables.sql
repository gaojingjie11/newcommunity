USE smart_community;

CREATE TABLE IF NOT EXISTS `workorders` (
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

CREATE TABLE IF NOT EXISTS `workorder_logs` (
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
