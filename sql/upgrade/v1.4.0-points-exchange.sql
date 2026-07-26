-- RuoYi-Go BY v1.4.0：通用积分、充值、任务兑换中心、可选雨云核验。
-- MySQL 5.7+，可重复执行。
SET NAMES utf8mb4;
START TRANSACTION;

CREATE TABLE IF NOT EXISTS platform_point_accounts (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    system_user_id BIGINT NOT NULL,
    balance_minor BIGINT UNSIGNED NOT NULL DEFAULT 0 COMMENT '百分之一积分',
    status VARCHAR(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT 'ACTIVE',
    lock_version BIGINT UNSIGNED NOT NULL DEFAULT 1,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    PRIMARY KEY (id),
    UNIQUE KEY uk_platform_point_account_user (system_user_id),
    CONSTRAINT fk_platform_point_account_user FOREIGN KEY (system_user_id) REFERENCES sys_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS platform_point_ledger (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    public_id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    account_id BIGINT UNSIGNED NOT NULL,
    entry_type VARCHAR(24) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    amount_minor BIGINT NOT NULL,
    balance_after_minor BIGINT UNSIGNED NOT NULL,
    reference_type VARCHAR(32) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    reference_id VARCHAR(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    idempotency_key VARCHAR(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    description VARCHAR(255) NOT NULL DEFAULT '',
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    PRIMARY KEY (id),
    UNIQUE KEY uk_platform_point_ledger_public (public_id),
    UNIQUE KEY uk_platform_point_ledger_idempotency (account_id, idempotency_key),
    KEY idx_platform_point_ledger_account_time (account_id, created_at),
    CONSTRAINT fk_platform_point_ledger_account FOREIGN KEY (account_id) REFERENCES platform_point_accounts (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS platform_point_recharge_orders (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    public_id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    order_no VARCHAR(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    account_id BIGINT UNSIGNED NOT NULL,
    system_user_id BIGINT NOT NULL,
    points BIGINT UNSIGNED NOT NULL,
    amount_minor BIGINT UNSIGNED NOT NULL,
    pay_type VARCHAR(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    status VARCHAR(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT 'PENDING',
    idempotency_key VARCHAR(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    external_order_id VARCHAR(191) CHARACTER SET ascii COLLATE ascii_bin NULL,
    paid_at DATETIME(6) NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    PRIMARY KEY (id),
    UNIQUE KEY uk_platform_recharge_public (public_id),
    UNIQUE KEY uk_platform_recharge_order (order_no),
    UNIQUE KEY uk_platform_recharge_idempotency (account_id, idempotency_key),
    KEY idx_platform_recharge_user (system_user_id, created_at),
    KEY idx_platform_recharge_status (status, created_at),
    CONSTRAINT fk_platform_recharge_account FOREIGN KEY (account_id) REFERENCES platform_point_accounts (id),
    CONSTRAINT fk_platform_recharge_user FOREIGN KEY (system_user_id) REFERENCES sys_user (user_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS platform_reward_tasks (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    public_id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    task_code VARCHAR(64) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    name VARCHAR(120) NOT NULL,
    summary VARCHAR(500) NOT NULL DEFAULT '',
    provider VARCHAR(32) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT 'MANUAL',
    verify_mode VARCHAR(32) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT 'MANUAL',
    reward_points BIGINT UNSIGNED NOT NULL,
    action_url VARCHAR(500) NOT NULL DEFAULT '',
    requirements VARCHAR(1000) NOT NULL DEFAULT '',
    display_order INT NOT NULL DEFAULT 100,
    status VARCHAR(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT 'ACTIVE',
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    PRIMARY KEY (id),
    UNIQUE KEY uk_platform_reward_task_public (public_id),
    UNIQUE KEY uk_platform_reward_task_code (task_code),
    KEY idx_platform_reward_task_status_order (status, display_order)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS platform_exchange_codes (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    public_id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    code_hmac BINARY(32) NOT NULL,
    code_mask VARCHAR(40) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    reward_points BIGINT UNSIGNED NOT NULL,
    source_type VARCHAR(32) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    source_id BIGINT UNSIGNED NULL,
    owner_system_user_id BIGINT NULL,
    status VARCHAR(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT 'UNUSED',
    description VARCHAR(255) NOT NULL DEFAULT '',
    expires_at DATETIME(6) NULL,
    redeemed_account_id BIGINT UNSIGNED NULL,
    redeemed_system_user_id BIGINT NULL,
    redeemed_at DATETIME(6) NULL,
    ledger_id BIGINT UNSIGNED NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    PRIMARY KEY (id),
    UNIQUE KEY uk_platform_exchange_code_public (public_id),
    UNIQUE KEY uk_platform_exchange_code_hmac (code_hmac),
    UNIQUE KEY uk_platform_exchange_code_source (source_type, source_id),
    KEY idx_platform_exchange_code_status (status, expires_at),
    CONSTRAINT fk_platform_exchange_owner FOREIGN KEY (owner_system_user_id) REFERENCES sys_user (user_id),
    CONSTRAINT fk_platform_exchange_account FOREIGN KEY (redeemed_account_id) REFERENCES platform_point_accounts (id),
    CONSTRAINT fk_platform_exchange_user FOREIGN KEY (redeemed_system_user_id) REFERENCES sys_user (user_id),
    CONSTRAINT fk_platform_exchange_ledger FOREIGN KEY (ledger_id) REFERENCES platform_point_ledger (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS platform_reward_claims (
    id BIGINT UNSIGNED NOT NULL AUTO_INCREMENT,
    public_id CHAR(36) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    task_id BIGINT UNSIGNED NOT NULL,
    account_id BIGINT UNSIGNED NOT NULL,
    system_user_id BIGINT NOT NULL,
    provider_subject VARCHAR(128) CHARACTER SET ascii COLLATE ascii_bin NOT NULL,
    status VARCHAR(16) CHARACTER SET ascii COLLATE ascii_bin NOT NULL DEFAULT 'PENDING',
    exchange_code_id BIGINT UNSIGNED NULL,
    review_note VARCHAR(500) NOT NULL DEFAULT '',
    verified_at DATETIME(6) NULL,
    issued_at DATETIME(6) NULL,
    created_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6),
    updated_at DATETIME(6) NOT NULL DEFAULT CURRENT_TIMESTAMP(6) ON UPDATE CURRENT_TIMESTAMP(6),
    PRIMARY KEY (id),
    UNIQUE KEY uk_platform_reward_claim_public (public_id),
    UNIQUE KEY uk_platform_reward_claim_account_task (account_id, task_id),
    UNIQUE KEY uk_platform_reward_claim_subject (task_id, provider_subject),
    UNIQUE KEY uk_platform_reward_claim_code (exchange_code_id),
    KEY idx_platform_reward_claim_status (status, created_at),
    CONSTRAINT fk_platform_reward_claim_task FOREIGN KEY (task_id) REFERENCES platform_reward_tasks (id),
    CONSTRAINT fk_platform_reward_claim_account FOREIGN KEY (account_id) REFERENCES platform_point_accounts (id),
    CONSTRAINT fk_platform_reward_claim_user FOREIGN KEY (system_user_id) REFERENCES sys_user (user_id),
    CONSTRAINT fk_platform_reward_claim_code FOREIGN KEY (exchange_code_id) REFERENCES platform_exchange_codes (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

INSERT IGNORE INTO platform_point_accounts (system_user_id, balance_minor, status, lock_version)
SELECT user_id, 0, 'ACTIVE', 1 FROM sys_user WHERE delete_time IS NULL;

INSERT INTO platform_reward_tasks (
    public_id, task_code, name, summary, provider, verify_mode, reward_points,
    action_url, requirements, display_order, status
) SELECT UUID(), 'RAINYUN_INVITE', '注册雨云，领取积分',
    '通过项目配置的雨云邀请关系注册，核验通过后获得可转让积分兑换码。',
    'RAINYUN', 'RAINYUN_SUBUSER', 20, 'https://www.rainyun.com/',
    '管理员可配置活动入口、奖励积分与雨云 X-Api-Key；同一账号和雨云 UID 只能领取一次。',
    10, 'ACTIVE'
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM platform_reward_tasks WHERE task_code = 'RAINYUN_INVITE');

INSERT INTO sys_system_setting (
    setting_key, setting_group, setting_value, remark, create_by, create_time
) SELECT 'platform.exchange', 'platform',
    '{"enabled":true,"rainyunEnabled":false,"rainyunApiBaseUrl":"https://api.rainyun.com","rainyunInviteUrl":"https://www.rainyun.com/","rainyunApiKey":""}',
    '通用兑换中心设置；雨云任务默认关闭，密钥不回显', 'migration-v1.4.0', NOW()
FROM DUAL WHERE NOT EXISTS (SELECT 1 FROM sys_system_setting WHERE setting_key = 'platform.exchange');

INSERT INTO sys_menu (
    menu_id, menu_name, parent_id, order_num, path, component, query, route_name,
    is_frame, is_cache, menu_type, visible, perms, icon, status,
    create_by, create_time, update_by, update_time, delete_time, remark
) VALUES
    (6200, '控制台', 2, 1, 'console', 'account/console/index', '', 'AccountConsole', 1, 0, 'C', '0', 'account:console:view', 'dashboard', '0', 'migration-v1.4.0', NOW(), '', NULL, NULL, '用户积分与任务概览'),
    (6201, '服务市场', 2, 2, 'marketplace', 'account/marketplace/index', '', 'AccountMarketplace', 1, 0, 'C', '0', 'account:marketplace:view', 'shopping', '0', 'migration-v1.4.0', NOW(), '', NULL, NULL, '业务项目扩展自己的商品与服务'),
    (6202, '兑换中心', 2, 3, 'exchange', 'account/exchange/index', '', 'AccountExchange', 1, 0, 'C', '0', 'account:exchange:view', 'ticket', '0', 'migration-v1.4.0', NOW(), '', NULL, NULL, '任务与积分兑换码'),
    (6203, '积分与订单', 2, 4, 'points', 'account/points/index', '', 'AccountPoints', 1, 0, 'C', '0', 'account:points:view', 'money', '0', 'migration-v1.4.0', NOW(), '', NULL, NULL, '积分充值、订单与流水'),
    (6000, '积分运营', 0, 80, 'platform', NULL, '', 'PlatformOperations', 1, 0, 'M', '0', '', 'money', '0', 'migration-v1.4.0', NOW(), '', NULL, NULL, '管理员积分与兑换运营'),
    (6001, '积分管理', 6000, 1, 'points', 'platform/points/index', '', 'PlatformPoints', 1, 0, 'C', '0', 'platform:points:list', 'money', '0', 'migration-v1.4.0', NOW(), '', NULL, NULL, '积分账户与人工调整'),
    (6002, '兑换中心', 6000, 2, 'exchange', 'platform/exchange/index', '', 'PlatformExchange', 1, 0, 'C', '0', 'platform:exchange:list', 'ticket', '0', 'migration-v1.4.0', NOW(), '', NULL, NULL, '任务、雨云与兑换码')
ON DUPLICATE KEY UPDATE
    menu_name = VALUES(menu_name), parent_id = VALUES(parent_id), order_num = VALUES(order_num),
    path = VALUES(path), component = VALUES(component), route_name = VALUES(route_name),
    visible = VALUES(visible), perms = VALUES(perms), icon = VALUES(icon), status = VALUES(status),
    update_by = 'migration-v1.4.0', update_time = NOW(), delete_time = NULL, remark = VALUES(remark);

UPDATE sys_menu SET order_num = 5 WHERE menu_id = 200;
UPDATE sys_menu SET order_num = 99, remark = '侧边栏作为“我的服务”分界线展开为直接菜单' WHERE menu_id = 2;

INSERT INTO sys_menu (
    menu_id, menu_name, parent_id, order_num, path, component, query, route_name,
    is_frame, is_cache, menu_type, visible, perms, icon, status,
    create_by, create_time, update_by, update_time, delete_time, remark
) VALUES
    (16000, '积分账户查询', 6001, 1, '#', '', '', '', 1, 0, 'F', '1', 'platform:points:list', '#', '0', 'migration-v1.4.0', NOW(), '', NULL, NULL, ''),
    (16001, '积分人工调整', 6001, 2, '#', '', '', '', 1, 0, 'F', '1', 'platform:points:adjust', '#', '0', 'migration-v1.4.0', NOW(), '', NULL, NULL, ''),
    (16010, '兑换中心查询', 6002, 1, '#', '', '', '', 1, 0, 'F', '1', 'platform:exchange:list', '#', '0', 'migration-v1.4.0', NOW(), '', NULL, NULL, ''),
    (16011, '兑换中心维护', 6002, 2, '#', '', '', '', 1, 0, 'F', '1', 'platform:exchange:edit', '#', '0', 'migration-v1.4.0', NOW(), '', NULL, NULL, '')
ON DUPLICATE KEY UPDATE perms = VALUES(perms), status = VALUES(status), delete_time = NULL;

SET @common_role_id := (SELECT role_id FROM sys_role WHERE role_key = 'common' AND delete_time IS NULL ORDER BY role_id LIMIT 1);
INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT @common_role_id, menu_id FROM sys_menu WHERE menu_id IN (2, 200, 6200, 6201, 6202, 6203);

SET @admin_role_id := (SELECT role_id FROM sys_role WHERE role_key = 'admin' AND delete_time IS NULL ORDER BY role_id LIMIT 1);
INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
SELECT @admin_role_id, menu_id
FROM sys_menu
WHERE menu_id IN (2, 200, 6000, 6001, 6002, 6200, 6201, 6202, 6203, 16000, 16001, 16010, 16011);

COMMIT;
