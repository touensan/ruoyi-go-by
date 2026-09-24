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

## 2026-09-05｜源码与历史发布状态纠正

- 核对当前源码、历史标签、README/前端说明、功能/部署文档和 GitHub Release。新增 `VERSIONS.md`，明确 main 与标签/附件不会自动同步；历史功能限制保留所属版本，避免将“未联调”误写为“未实现”。
- 审查覆盖全部六条历史 Release，当前源码与历史标签分别说明。本次不创建新标签、不替换下载包，也不改变既有附件校验和；发布页说明与本文档配套维护。
- 发现并修正定时任务、代码生成和公告占位接口误报成功：29 个未实现入口改为 HTTP/业务码 501，顶部通知保留明确 `supported=false` 的空列表兼容。数据库表/列只读查看保留。新增 29 个入口子用例和 1 项顶部通知测试，Go test/vet/build 通过；本轮不依赖数据库，已有 RBAC 数据库用例按未设置 DSN 跳过，不能当作本轮 MySQL 回归。
- 使用 remote-build 资源探测后在独立构建目录以 3 个编译任务验证，仅传源码。版本与依赖声明未改，原有前端无需重建。

## 2026-09-08｜v1.1.3 原生主线 Release 验证

- Go 1.26.3：`go test -p 3 ./...`、`go vet -p 3 ./...` 与 `CGO_ENABLED=0 go build -trimpath` 通过；本轮未设置 MySQL DSN，因此数据库用例跳过，同一运行代码此前已完成真实 MySQL 5.7 的 5 个 RBAC 子场景。
- 前端将 Axios、ECharts、js-cookie 与 glob 更新到修复版本；Node 22.22.0 上 7 项 RBAC 测试及生产构建通过，`npm audit --omit=dev` 为 0 项报告。
- Linux x86_64 后端为静态、去除调试符号的 ELF。运行包包含该二进制、构建后的前端、示例配置、SQL 与文档，不含 `application.yaml`、真实凭据、上传、日志、Node 依赖目录或项目记忆。
- 运行包按源码和构建规格生成不可变缓存，再以只读凭据取回核对；GitHub 附件另行下载验证。v1.1.3 延续原生若依 v1.1.2 主线，v1.2.0–v1.4.0 仍是 Whiteyun 历史支线；发布不代表任何业务站点已部署。

## 2026-09-24｜单机平滑更新源码

- 新增 `deploy/smooth-release.py`、双端口配置模板和 `SMOOTH_RELEASE.md`；Go 候选可关闭启动时的自动建表，`/internal/ready` 检查 MySQL/Redis，SIGTERM 等待在途 HTTP 请求。修正旧 `service.sh` 强杀及重启路径，并修正示例 YAML 中 `user.password` 的层级，保证示例可解析。
- 打包脚本追加 `FILE_SHA256SUMS`；工具对归档、逐文件、Linux x86_64 二进制、静态资源同名内容及上传目录执行检查。旧版本页面的 hash 静态资源保留。旧 v1.1.3 包未改，生产业务站未部署。
- 深圳隔离构建节点检查前磁盘可用约 118 GiB/inode 17%、可用内存约 10.5 GiB；本机仅约 2.8 GiB available 且 swap 占用高，因此使用隔离 Go 缓存与测试目录。`RUOYI_CONFIG_FILE` 指向仓库空测试配置时 `go test -p 2 ./...` 通过；示例 YAML 解析测试 `go test ./config` 通过。未设置测试 MySQL DSN，数据库 RBAC 用例未作真实数据库回归。
- Python 发布工具 3 项完整性/准备测试通过；隔离 Nginx 片段 `nginx -t` 通过；隔离双后端切流演练中，一条旧慢请求正常结束，新请求抵达候选。尚未针对真实站点、支付回调或海量并发做生产验收。

## 2026-09-24｜v1.1.4 Release

- v1.1.4 将上述单机平滑更新能力作为新的主线补丁版本发布；历史 v1.1.3 与 Whiteyun 支线标签保持不变。本次只发布 GitHub 运行包，不升级业务站。
- 后端从 v1.1.4 标签源码在深圳隔离构建节点以 Go 1.26.3、`CGO_ENABLED=0`、Linux amd64 构建；前端源码相对 v1.1.3 未变化，复用并核对 v1.1.3 已发布前端文件摘要。最终包包含逐文件摘要、构建来源、最新发布工具和文档。
