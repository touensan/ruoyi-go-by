# 开发与验证记录

开始前阅读 `ai_codex.md`、`AGENTS.md` 与 `docs/ARCHITECTURE.md`。只记录已验证事实，不将计划或测试跳过写成成功。

## 2026-09-05｜统一后台 RBAC（实现与验证完成）

统一架构契约见 [ARCHITECTURE.md](ARCHITECTURE.md)。`admin` 单向继承启用的 `common` 功能授权，服务端菜单、路由、接口与前端权限规则一致；前端全权限判断来自服务端 `*:*:*`，角色全部撤销时同步清空旧权限。

- Go 1.26.2：`go test -p 3 ./...`、`go vet -p 3 ./...` 和 `go build -p 3` 通过。真实 MySQL 5.7.43 上的 RBAC 测试含 5 个子场景（继承/菜单与权限、反向和匿名拒绝、超级管理员、停用/删除角色、停用/删除菜单及撤销授权）。
- 前端 Node 22.22.3：依照现有无锁文件安装流程执行 `npm install`，7 项 `npm run test:rbac` 测试、`npm run build:prod` 通过。没有改动依赖版本声明。
- 修复 `HasAnyRoles` 错调用权限查询；服务端权限、路由采用有效角色并集，过滤停用/软删除角色和菜单，匿名 ID 不再落入全菜单分支。
- `RUOYI_CONFIG_FILE` 可指定配置路径，默认仍为 `application.yaml`。测试使用空 `app/security/testdata/rbac.yaml`，不读取部署配置；保持原有启动期初始化顺序，避免 Redis key 的包级初始化先于配置加载。

运行回归时，先提供专用空 MySQL 库 `ruoyi_go_rbac_test`（测试拒绝已有目标表的数据库），设置 `RBAC_TEST_MYSQL_DSN`，再执行：

```bash
RUOYI_CONFIG_FILE="$PWD/app/security/testdata/rbac.yaml" go test -p 3 ./...
cd frontend/RuoYi-Vue3-ts
npm run test:rbac
npm run build:prod
```

未设置 `RBAC_TEST_MYSQL_DSN` 时数据库用例会跳过；跳过不能视为数据库验证通过。
本次调整公开框架源码与开发约定，不涉及其他业务站点的生产部署。Go/PHP 前端通过权限函数回归和构建；两类账号的真实浏览器场景在同系列 Rust 测试实例执行，不把它描述为三个后端均运行了浏览器验收。
