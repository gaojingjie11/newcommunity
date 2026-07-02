-- Development-only user ID resequencing script for MySQL 8.x
-- Purpose:
-- 1. Keep sys_user.id = 1 unchanged if it exists
-- 2. Resequence remaining users by current id order: 100,101,102 -> 2,3,4
-- 3. Rewrite all known user-id references in related tables
--
-- Run only after taking a full database backup.
-- After running, clear Redis login / RBAC caches before letting users log in again.

START TRANSACTION;

DROP TEMPORARY TABLE IF EXISTS tmp_user_id_map;
CREATE TEMPORARY TABLE tmp_user_id_map (
  old_id BIGINT PRIMARY KEY,
  new_id BIGINT NOT NULL UNIQUE,
  temp_id BIGINT NOT NULL UNIQUE
);

INSERT INTO tmp_user_id_map (old_id, new_id, temp_id)
SELECT id, 1, 1000001
FROM sys_user
WHERE id = 1;

SET @next_user_id := (
  SELECT CASE WHEN EXISTS (SELECT 1 FROM sys_user WHERE id = 1) THEN 1 ELSE 0 END
);

INSERT INTO tmp_user_id_map (old_id, new_id, temp_id)
SELECT src.id,
       (@next_user_id := @next_user_id + 1) AS new_id,
       1000000 + @next_user_id AS temp_id
FROM (
  SELECT id
  FROM sys_user
  WHERE id <> 1
  ORDER BY id ASC
) src;

DROP PROCEDURE IF EXISTS remap_user_ref_if_exists;
DELIMITER //
CREATE PROCEDURE remap_user_ref_if_exists(
  IN p_table VARCHAR(64),
  IN p_column VARCHAR(64),
  IN p_from_col VARCHAR(16),
  IN p_to_col VARCHAR(16),
  IN p_where_clause VARCHAR(255)
)
BEGIN
  IF EXISTS (
    SELECT 1
    FROM information_schema.columns
    WHERE table_schema = DATABASE()
      AND table_name = p_table
      AND column_name = p_column
  ) THEN
    SET @sql = CONCAT(
      'UPDATE `', p_table, '` t ',
      'JOIN tmp_user_id_map m ON t.`', p_column, '` = m.`', p_from_col, '` ',
      'SET t.`', p_column, '` = m.`', p_to_col, '`'
    );

    IF p_where_clause IS NOT NULL AND p_where_clause <> '' THEN
      SET @sql = CONCAT(@sql, ' WHERE ', p_where_clause);
    END IF;

    PREPARE stmt FROM @sql;
    EXECUTE stmt;
    DEALLOCATE PREPARE stmt;
  END IF;
END//
DELIMITER ;

-- Stage 1: move all user ids and references to temporary ids to avoid collisions.
CALL remap_user_ref_if_exists('sys_user', 'id', 'old_id', 'temp_id', '');
CALL remap_user_ref_if_exists('sys_user_role', 'user_id', 'old_id', 'temp_id', '');
CALL remap_user_ref_if_exists('user_login_logs', 'user_id', 'old_id', 'temp_id', '');
CALL remap_user_ref_if_exists('admin_login_logs', 'admin_user_id', 'old_id', 'temp_id', '');
CALL remap_user_ref_if_exists('oms_cart', 'user_id', 'old_id', 'temp_id', '');
CALL remap_user_ref_if_exists('oms_order', 'user_id', 'old_id', 'temp_id', '');
CALL remap_user_ref_if_exists('pms_favorite', 'user_id', 'old_id', 'temp_id', '');
CALL remap_user_ref_if_exists('pms_product_comment', 'user_id', 'old_id', 'temp_id', '');
CALL remap_user_ref_if_exists('wallets', 'user_id', 'old_id', 'temp_id', '');
CALL remap_user_ref_if_exists('wallet_transactions', 'user_id', 'old_id', 'temp_id', '');
CALL remap_user_ref_if_exists('wallet_transactions', 'related_id', 'old_id', 'temp_id', '`type` = 2 OR `biz_type` IN (''admin_recharge'', ''admin_deduct'')');
CALL remap_user_ref_if_exists('payment_records', 'user_id', 'old_id', 'temp_id', '');
CALL remap_user_ref_if_exists('product_view_logs', 'user_id', 'old_id', 'temp_id', '`user_id` > 0');
CALL remap_user_ref_if_exists('notice_view_logs', 'user_id', 'old_id', 'temp_id', '');
CALL remap_user_ref_if_exists('visitors', 'user_id', 'old_id', 'temp_id', '');
CALL remap_user_ref_if_exists('user_parking_bindings', 'user_id', 'old_id', 'temp_id', '');
CALL remap_user_ref_if_exists('property_fees', 'user_id', 'old_id', 'temp_id', '');
CALL remap_user_ref_if_exists('property_fee_payments', 'user_id', 'old_id', 'temp_id', '');
CALL remap_user_ref_if_exists('workorders', 'user_id', 'old_id', 'temp_id', '');
CALL remap_user_ref_if_exists('workorders', 'processor_id', 'old_id', 'temp_id', '`processor_id` > 0');
CALL remap_user_ref_if_exists('workorder_logs', 'operator_id', 'old_id', 'temp_id', '`operator_id` > 0');
CALL remap_user_ref_if_exists('community_messages', 'user_id', 'old_id', 'temp_id', '');
CALL remap_user_ref_if_exists('cms_ai_report', 'generated_by', 'old_id', 'temp_id', '`generated_by` > 0');
CALL remap_user_ref_if_exists('pms_user_store', 'user_id', 'old_id', 'temp_id', '');
CALL remap_user_ref_if_exists('sys_user_conversation', 'user_id', 'old_id', 'temp_id', '');
CALL remap_user_ref_if_exists('sys_user_chat_message', 'user_id', 'old_id', 'temp_id', '');
CALL remap_user_ref_if_exists('agent_action_approval', 'user_id', 'old_id', 'temp_id', '');

-- Stage 2: move temporary ids to final compact ids.
CALL remap_user_ref_if_exists('sys_user', 'id', 'temp_id', 'new_id', '');
CALL remap_user_ref_if_exists('sys_user_role', 'user_id', 'temp_id', 'new_id', '');
CALL remap_user_ref_if_exists('user_login_logs', 'user_id', 'temp_id', 'new_id', '');
CALL remap_user_ref_if_exists('admin_login_logs', 'admin_user_id', 'temp_id', 'new_id', '');
CALL remap_user_ref_if_exists('oms_cart', 'user_id', 'temp_id', 'new_id', '');
CALL remap_user_ref_if_exists('oms_order', 'user_id', 'temp_id', 'new_id', '');
CALL remap_user_ref_if_exists('pms_favorite', 'user_id', 'temp_id', 'new_id', '');
CALL remap_user_ref_if_exists('pms_product_comment', 'user_id', 'temp_id', 'new_id', '');
CALL remap_user_ref_if_exists('wallets', 'user_id', 'temp_id', 'new_id', '');
CALL remap_user_ref_if_exists('wallet_transactions', 'user_id', 'temp_id', 'new_id', '');
CALL remap_user_ref_if_exists('wallet_transactions', 'related_id', 'temp_id', 'new_id', '`type` = 2 OR `biz_type` IN (''admin_recharge'', ''admin_deduct'')');
CALL remap_user_ref_if_exists('payment_records', 'user_id', 'temp_id', 'new_id', '');
CALL remap_user_ref_if_exists('product_view_logs', 'user_id', 'temp_id', 'new_id', '`user_id` > 0');
CALL remap_user_ref_if_exists('notice_view_logs', 'user_id', 'temp_id', 'new_id', '');
CALL remap_user_ref_if_exists('visitors', 'user_id', 'temp_id', 'new_id', '');
CALL remap_user_ref_if_exists('user_parking_bindings', 'user_id', 'temp_id', 'new_id', '');
CALL remap_user_ref_if_exists('property_fees', 'user_id', 'temp_id', 'new_id', '');
CALL remap_user_ref_if_exists('property_fee_payments', 'user_id', 'temp_id', 'new_id', '');
CALL remap_user_ref_if_exists('workorders', 'user_id', 'temp_id', 'new_id', '');
CALL remap_user_ref_if_exists('workorders', 'processor_id', 'temp_id', 'new_id', '`processor_id` > 0');
CALL remap_user_ref_if_exists('workorder_logs', 'operator_id', 'temp_id', 'new_id', '`operator_id` > 0');
CALL remap_user_ref_if_exists('community_messages', 'user_id', 'temp_id', 'new_id', '');
CALL remap_user_ref_if_exists('cms_ai_report', 'generated_by', 'temp_id', 'new_id', '`generated_by` > 0');
CALL remap_user_ref_if_exists('pms_user_store', 'user_id', 'temp_id', 'new_id', '');
CALL remap_user_ref_if_exists('sys_user_conversation', 'user_id', 'temp_id', 'new_id', '');
CALL remap_user_ref_if_exists('sys_user_chat_message', 'user_id', 'temp_id', 'new_id', '');
CALL remap_user_ref_if_exists('agent_action_approval', 'user_id', 'temp_id', 'new_id', '');

SET @next_auto_id = (SELECT COALESCE(MAX(id), 0) + 1 FROM sys_user);
SET @alter_sql = CONCAT('ALTER TABLE `sys_user` AUTO_INCREMENT = ', @next_auto_id);
PREPARE stmt_auto FROM @alter_sql;
EXECUTE stmt_auto;
DEALLOCATE PREPARE stmt_auto;

DROP PROCEDURE IF EXISTS remap_user_ref_if_exists;

SELECT old_id, new_id
FROM tmp_user_id_map
ORDER BY new_id;

COMMIT;

