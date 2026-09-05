package security_test

import (
	"os"
	"strings"
	"testing"

	drivermysql "github.com/go-sql-driver/mysql"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"ruoyi-go/app/model"
	"ruoyi-go/app/security"
	"ruoyi-go/app/service"
	"ruoyi-go/framework/dal"
)

// Only a dedicated empty test database may be used; use RUOYI_CONFIG_FILE with the empty testdata/rbac.yaml fixture.
func TestUnifiedBackendRBAC(t *testing.T) {
	dsn := os.Getenv("RBAC_TEST_MYSQL_DSN")
	if dsn == "" {
		t.Skip("set RBAC_TEST_MYSQL_DSN for isolated MySQL regression")
	}
	cfg, err := drivermysql.ParseDSN(dsn)
	if err != nil || cfg.DBName != "ruoyi_go_rbac_test" {
		t.Fatal("requires database ruoyi_go_rbac_test")
	}
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatal("cannot connect to isolated MySQL")
	}
	sqlDB, err := db.DB()
	if err != nil {
		t.Fatal(err)
	}
	defer sqlDB.Close()
	for _, table := range []string{"sys_role", "sys_menu", "sys_user_role", "sys_role_menu"} {
		if db.Migrator().HasTable(table) {
			t.Fatal("requires empty test database")
		}
	}
	if err := db.AutoMigrate(&model.SysRole{}, &model.SysMenu{}, &model.SysUserRole{}, &model.SysRoleMenu{}); err != nil {
		t.Fatal(err)
	}
	defer db.Migrator().DropTable(&model.SysRoleMenu{}, &model.SysUserRole{}, &model.SysMenu{}, &model.SysRole{})
	old := dal.Gorm
	dal.Gorm = db
	defer func() { dal.Gorm = old }()
	exec := func(sql string, args ...interface{}) {
		t.Helper()
		if err := db.Exec(sql, args...).Error; err != nil {
			t.Fatal(err)
		}
	}
	exec("INSERT INTO sys_role(role_id,role_name,role_key,status,data_scope) VALUES (1,'Admin','admin','0','5'),(2,'User','common','0','5'),(3,'Other','other','0','5')")
	exec("INSERT INTO sys_user_role(user_id,role_id) VALUES (1,1),(9,1),(10,2),(11,3)")
	exec("INSERT INTO sys_menu(menu_id,menu_name,parent_id,order_num,path,component,menu_type,perms,status) VALUES (101,'User',0,1,'member','member/index','C','member:read,member:edit','0'),(102,'Admin',0,2,'manage','manage/index','C','system:manage','0')")
	exec("INSERT INTO sys_role_menu(role_id,menu_id) VALUES (1,102),(2,101)")
	expect := func(ok bool, message string) {
		t.Helper()
		if !ok {
			t.Error(message)
		}
	}
	routes := func(id int) string {
		var paths []string
		for _, menu := range (&service.MenuService{}).GetMenuMCListByUserId(id) {
			paths = append(paths, menu.Path)
		}
		return strings.Join(paths, ",")
	}
	t.Run("same_roles_menus_and_api_permissions", func(t *testing.T) {
		expect(security.HasRole(9, "common"), "admin must inherit common role")
		expect(security.HasAnyRoles(9, []string{"missing", "common"}), "any-role must check roles")
		expect(security.HasPerm(9, "member:read") && security.HasPerm(9, "member:edit"), "admin must inherit comma-separated user permissions")
		expect(security.HasPerm(9, "system:manage"), "admin retains direct management grant")
		expect(routes(9) == "member,manage", "admin menus must include both functions")
		expect(len((&service.RoleService{}).GetRoleListByUserId(9)) == 1, "inherited roles must not become explicit assignments")
		expect(!security.HasRole(9, "other"), "ordinary administrator is not a wildcard superuser")
	})
	t.Run("no_reverse_or_anonymous_access", func(t *testing.T) {
		expect(security.HasPerm(10, "member:read"), "user retains own permission")
		expect(!security.HasPerm(10, "system:manage") && !security.HasRole(10, "admin"), "user must not inherit admin")
		expect(routes(10) == "member", "user menus exclude management")
		expect(routes(0) == "" && !security.HasPerm(0, "member:read"), "anonymous users must not get menus or permissions")
		expect(!security.HasPerm(11, "member:read"), "unrelated roles do not inherit common")
		expect(!security.HasAnyRoles(9, []string{}) && !security.HasPerm(1, ""), "empty authorization checks fail closed")
	})
	t.Run("superuser_covers_user_functions", func(t *testing.T) {
		expect(security.HasRole(1, "common") && security.HasAnyRoles(1, []string{"other"}), "superuser role checks")
		expect(security.HasPerm(1, "member:read") && security.HasPerm(1, "unassigned:permission"), "superuser permission checks")
		expect(routes(1) == "member,manage", "superuser sees user menus")
	})
	t.Run("disabled_and_deleted_roles_revoke_immediately", func(t *testing.T) {
		exec("UPDATE sys_role SET status='1' WHERE role_id=2")
		expect(!security.HasRole(9, "common") && !security.HasPerm(9, "member:read"), "disabled common role revokes inheritance")
		expect(routes(9) == "manage", "disabled common removes inherited route")
		exec("UPDATE sys_role SET status='0',delete_time=NOW() WHERE role_id=2")
		expect(!security.HasRole(9, "common") && !security.HasPerm(9, "member:read"), "deleted common role revokes inheritance")
		exec("UPDATE sys_role SET delete_time=NULL WHERE role_id=2")
		exec("UPDATE sys_role SET status='1' WHERE role_id=1")
		expect(!security.HasRole(9, "admin") && !security.HasPerm(9, "member:read") && routes(9) == "", "disabled admin loses direct and inherited grants")
		exec("UPDATE sys_role SET status='0' WHERE role_id=1")
	})
	t.Run("disabled_and_deleted_menus_revoke_immediately", func(t *testing.T) {
		exec("UPDATE sys_menu SET status='1' WHERE menu_id=101")
		expect(!security.HasPerm(9, "member:read") && routes(9) == "manage", "disabled menu revokes permission and route")
		exec("UPDATE sys_menu SET status='0',delete_time=NOW() WHERE menu_id=101")
		expect(!security.HasPerm(9, "member:read") && routes(9) == "manage", "deleted menu revokes permission and route")
		exec("UPDATE sys_menu SET delete_time=NULL WHERE menu_id=101")
		exec("DELETE FROM sys_role_menu WHERE role_id=2")
		expect(!security.HasPerm(9, "member:read"), "removed grant revokes permission without token refresh")
	})
}
