#!/usr/bin/env python3
"""Verify the v1.3+ shared user workspace invariants without a live database."""

from pathlib import Path
import re


ROOT = Path(__file__).resolve().parents[1]


def require(condition: bool, message: str) -> None:
    if not condition:
        raise SystemExit(f"ERROR: {message}")


seed = (ROOT / "ruoyi.sql").read_text(encoding="utf-8")
upgrade = (ROOT / "sql/upgrade/v1.3.0-user-workspace.sql").read_text(encoding="utf-8")
auth_controller = (ROOT / "app/controller/auth_controller.go").read_text(encoding="utf-8")
constants = (ROOT / "common/types/constant/constant.go").read_text(encoding="utf-8")
router = (ROOT / "frontend/whiteyun-vue/src/router/index.ts").read_text(encoding="utf-8")
package_json = (ROOT / "frontend/whiteyun-vue/package.json").read_text(encoding="utf-8")
application_example = (ROOT / "application-example.yaml").read_text(encoding="utf-8")

common_menu_ids = re.findall(
    r"insert into sys_role_menu values \('2', '([0-9]+)'\);",
    seed,
)
require(
    sorted(common_menu_ids) == ["2", "200"],
    f"common 角色应只绑定用户菜单 2/200，实际为 {common_menu_ids}",
)
require(
    "insert into sys_role values('2', '普通用户',    'common', 2, 5" in seed
    and "data_scope = '5'" in upgrade
    and "DELETE FROM sys_role_dept" in upgrade
    and "insert into sys_role_dept values ('2'," not in seed,
    "common 角色必须使用“仅本人”数据范围",
)

require(
    "insert into sys_menu values('2', '我的账户'" in seed,
    "全新安装 SQL 缺少“我的账户”目录",
)
require(
    seed.startswith("SET NAMES utf8mb4;")
    and "SET NAMES utf8mb4;" in upgrade,
    "初始化和升级 SQL 必须显式声明 utf8mb4 客户端字符集",
)
require(
    "insert into sys_menu values('200',  '个人资料'" in seed,
    "全新安装 SQL 缺少“个人资料”菜单",
)
require(
    "DEFAULT_USER_ROLE_KEY = \"common\"" in constants,
    "默认注册角色常量未固定为 common",
)
require(
    "GetRoleByRoleKey(constant.DEFAULT_USER_ROLE_KEY)" in auth_controller
    and "[]int{defaultRoleID}" in auth_controller,
    "自助注册没有校验并绑定 common 角色",
)
require(
    "menu.menu_id = 1" in upgrade
    and "grandparent_menu.parent_id = 1" in upgrade
    and "INSERT IGNORE INTO sys_role_menu" in upgrade
    and "NOT EXISTS (" in upgrade,
    "升级脚本没有覆盖管理菜单回收、用户菜单授权和无角色账号修复",
)
require(
    "title: '工作台'" in router,
    "Whiteyun Vue 的共享首页未命名为“工作台”",
)
frontend_version = re.search(r'"version": "([0-9]+)\.([0-9]+)\.([0-9]+)"', package_json)
backend_version = re.search(r"version: ([0-9]+)\.([0-9]+)\.([0-9]+)", application_example)
require(
    frontend_version is not None
    and backend_version is not None
    and frontend_version.groups() == backend_version.groups()
    and tuple(map(int, frontend_version.groups())) >= (1, 3, 0),
    "v1.3+ 版本号未在前后端示例配置中同步",
)

print("user workspace architecture verification passed")
