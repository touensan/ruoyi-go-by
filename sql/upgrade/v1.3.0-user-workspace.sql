-- RuoYi-Go BY v1.3.0
-- 将默认信息架构升级为：
--   1. 普通用户只拥有共享用户工作台；
--   2. 超级管理员拥有用户工作台和全部管理菜单；
--   3. 自定义管理员通过 common 基础角色叠加管理角色。
--
-- 兼容 MySQL 5.7+，可重复执行。

SET NAMES utf8mb4;

START TRANSACTION;

INSERT INTO sys_role (
    role_name, role_key, role_sort, data_scope,
    menu_check_strictly, dept_check_strictly, status,
    create_by, create_time, update_by, update_time, delete_time, remark
)
SELECT
    '普通用户', 'common', 2, '5',
    1, 1, '0',
    'system', NOW(), '', NULL, NULL, '所有登录用户的基础角色'
FROM DUAL
WHERE NOT EXISTS (
    SELECT 1 FROM sys_role WHERE role_key = 'common' AND delete_time IS NULL
);

SET @common_role_id := (
    SELECT role_id
    FROM sys_role
    WHERE role_key = 'common' AND delete_time IS NULL
    ORDER BY role_id
    LIMIT 1
);

UPDATE sys_role
SET role_name = '普通用户',
    data_scope = '5',
    remark = '所有登录用户的基础角色',
    update_by = 'system',
    update_time = NOW()
WHERE role_id = @common_role_id;

DELETE FROM sys_role_dept
WHERE role_id = @common_role_id;

INSERT INTO sys_menu (
    menu_id, menu_name, parent_id, order_num, path, component, query, route_name,
    is_frame, is_cache, menu_type, visible, perms, icon, status,
    create_by, create_time, update_by, update_time, delete_time, remark
)
VALUES
    (2, '我的账户', 0, 1, 'account', NULL, '', '', 1, 0, 'M', '0', '', 'user', '0',
     'system', NOW(), '', NULL, NULL, '所有登录用户共享的账户目录'),
    (200, '个人资料', 2, 1, 'profile', 'system/user/profile/index', '', 'UserProfile',
     1, 0, 'C', '0', '', 'user', '0',
     'system', NOW(), '', NULL, NULL, '用户维护自己的资料、头像和密码')
ON DUPLICATE KEY UPDATE
    menu_name = VALUES(menu_name),
    parent_id = VALUES(parent_id),
    order_num = VALUES(order_num),
    path = VALUES(path),
    component = VALUES(component),
    query = VALUES(query),
    route_name = VALUES(route_name),
    is_frame = VALUES(is_frame),
    is_cache = VALUES(is_cache),
    menu_type = VALUES(menu_type),
    visible = VALUES(visible),
    perms = VALUES(perms),
    icon = VALUES(icon),
    status = VALUES(status),
    update_by = 'system',
    update_time = NOW(),
    delete_time = NULL,
    remark = VALUES(remark);

UPDATE sys_menu
SET order_num = 90,
    remark = '仅管理角色可见的系统管理目录',
    update_by = 'system',
    update_time = NOW()
WHERE menu_id = 1;

-- common 角色不再继承系统管理目录及其完整后代。
DELETE role_menu
FROM sys_role_menu AS role_menu
JOIN sys_menu AS menu ON menu.menu_id = role_menu.menu_id
LEFT JOIN sys_menu AS parent_menu ON parent_menu.menu_id = menu.parent_id
LEFT JOIN sys_menu AS grandparent_menu ON grandparent_menu.menu_id = parent_menu.parent_id
WHERE role_menu.role_id = @common_role_id
  AND (
      menu.menu_id = 1
      OR menu.parent_id = 1
      OR parent_menu.parent_id = 1
      OR grandparent_menu.parent_id = 1
  );

INSERT IGNORE INTO sys_role_menu (role_id, menu_id)
VALUES (@common_role_id, 2), (@common_role_id, 200);

-- 旧版本自助注册可能产生无角色账号；升级时补齐基础用户角色。
INSERT IGNORE INTO sys_user_role (user_id, role_id)
SELECT user_id, @common_role_id
FROM sys_user
WHERE delete_time IS NULL
  AND user_id <> 1
  AND NOT EXISTS (
      SELECT 1
      FROM sys_user_role
      WHERE sys_user_role.user_id = sys_user.user_id
  );

COMMIT;
