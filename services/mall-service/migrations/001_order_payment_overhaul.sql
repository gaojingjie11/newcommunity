-- =============================================================
-- Mall-Service: Order Payment Concurrency Overhaul Migration
-- =============================================================
-- This migration converts float64 columns to bigint (cents),
-- adds concurrency control fields, and creates the payment_records table.
--
-- IMPORTANT: Run this BEFORE deploying the new application code.
-- AutoMigrate will handle new tables/columns, but ALTERs for
-- existing columns should be done manually to avoid data loss.
-- =============================================================

-- 1. Orders table: float64 → bigint, add new columns
ALTER TABLE `oms_order`
    MODIFY COLUMN `total_amount` BIGINT NOT NULL DEFAULT 0 COMMENT '订单总额（分）',
    MODIFY COLUMN `used_balance` BIGINT NOT NULL DEFAULT 0 COMMENT '使用余额（分）',
    ADD COLUMN IF NOT EXISTS `expire_at` DATETIME DEFAULT NULL COMMENT '订单过期时间' AFTER `status`,
    ADD COLUMN IF NOT EXISTS `cancel_reason` VARCHAR(255) DEFAULT '' COMMENT '取消原因' AFTER `expire_at`,
    ADD COLUMN IF NOT EXISTS `cancelled_at` DATETIME DEFAULT NULL COMMENT '取消时间' AFTER `cancel_reason`,
    ADD COLUMN IF NOT EXISTS `paid_at` DATETIME DEFAULT NULL COMMENT '支付时间' AFTER `cancelled_at`,
    ADD COLUMN IF NOT EXISTS `version` INT NOT NULL DEFAULT 0 COMMENT '乐观锁版本' AFTER `paid_at`,
    ADD COLUMN IF NOT EXISTS `updated_at` DATETIME DEFAULT NULL COMMENT '更新时间' AFTER `version`,
    ADD INDEX IF NOT EXISTS `idx_expire_at` (`expire_at`);

-- 2. Order items: price float64 → bigint
ALTER TABLE `oms_order_item`
    MODIFY COLUMN `price` BIGINT NOT NULL DEFAULT 0 COMMENT '单价（分）',
    ADD COLUMN IF NOT EXISTS `product_snapshot` VARCHAR(512) DEFAULT '' COMMENT '下单时商品快照' AFTER `quantity`;

-- 3. Products: price float64 → bigint, add version
ALTER TABLE `pms_product`
    MODIFY COLUMN `price` BIGINT NOT NULL DEFAULT 0 COMMENT '售价（分）',
    MODIFY COLUMN `original_price` BIGINT NOT NULL DEFAULT 0 COMMENT '原价（分）',
    ADD COLUMN IF NOT EXISTS `version` INT NOT NULL DEFAULT 0 COMMENT '乐观锁版本';

-- 4. Wallets: balance float64 → bigint, add version
ALTER TABLE `wallets`
    MODIFY COLUMN `balance` BIGINT NOT NULL DEFAULT 0 COMMENT '余额（分）',
    ADD COLUMN IF NOT EXISTS `version` INT NOT NULL DEFAULT 0 COMMENT '乐观锁版本';

-- 5. Wallet transactions: amount float64 → bigint, add new fields
ALTER TABLE `wallet_transactions`
    MODIFY COLUMN `amount` BIGINT NOT NULL DEFAULT 0 COMMENT '金额（分）',
    ADD COLUMN IF NOT EXISTS `balance_before` BIGINT NOT NULL DEFAULT 0 COMMENT '变动前余额（分）',
    ADD COLUMN IF NOT EXISTS `balance_after` BIGINT NOT NULL DEFAULT 0 COMMENT '变动后余额（分）',
    ADD COLUMN IF NOT EXISTS `biz_type` VARCHAR(32) DEFAULT '' COMMENT '业务类型',
    ADD COLUMN IF NOT EXISTS `biz_id` VARCHAR(64) DEFAULT '' COMMENT '业务ID',
    ADD COLUMN IF NOT EXISTS `idempotency_key` VARCHAR(64) DEFAULT '' COMMENT '幂等键',
    ADD UNIQUE INDEX IF NOT EXISTS `uk_idempotency_key` (`idempotency_key`);

-- 6. Store products: add locked_stock, sold_count, version
ALTER TABLE `pms_store_product`
    ADD COLUMN IF NOT EXISTS `locked_stock` INT NOT NULL DEFAULT 0 COMMENT '锁定库存',
    ADD COLUMN IF NOT EXISTS `sold_count` INT NOT NULL DEFAULT 0 COMMENT '已售数量',
    ADD COLUMN IF NOT EXISTS `version` INT NOT NULL DEFAULT 0 COMMENT '乐观锁版本';

-- 7. Payment records table (new)
CREATE TABLE IF NOT EXISTS `payment_records` (
    `id`              BIGINT AUTO_INCREMENT PRIMARY KEY,
    `order_id`        BIGINT NOT NULL COMMENT '订单ID',
    `order_no`        VARCHAR(64) NOT NULL COMMENT '订单号',
    `user_id`         BIGINT NOT NULL COMMENT '用户ID',
    `amount`          BIGINT NOT NULL DEFAULT 0 COMMENT '支付金额（分）',
    `payment_method`  VARCHAR(32) NOT NULL DEFAULT 'wallet' COMMENT '支付方式',
    `status`          INT NOT NULL DEFAULT 0 COMMENT '状态: 0=init, 1=success, 2=failed',
    `idempotency_key` VARCHAR(64) NOT NULL COMMENT '幂等键',
    `fail_reason`     VARCHAR(255) DEFAULT '' COMMENT '失败原因',
    `paid_at`         DATETIME DEFAULT NULL COMMENT '支付成功时间',
    `created_at`      DATETIME DEFAULT CURRENT_TIMESTAMP,
    `updated_at`      DATETIME DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP,
    UNIQUE KEY `uk_idempotency_key` (`idempotency_key`),
    INDEX `idx_order_id` (`order_id`),
    INDEX `idx_user_id` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='支付记录表';
