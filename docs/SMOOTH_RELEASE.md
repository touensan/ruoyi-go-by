# 单机平滑更新

当前 `main` 提供单机双端口更新工具；v1.1.3 及以前的下载包没有此能力。它只处理应用版本更新，不保证机器、MySQL、Redis 故障时继续服务。

## 原理与首次接入

同一数据库、Redis、JWT 密钥和共享上传目录由新旧进程同时使用。候选版本先在回环备用端口启动，`/internal/ready` 检查数据库和 Redis；Nginx 配置验证后 reload，新连接走候选，旧 worker 继续完成旧连接。旧进程只在旧 worker 退出后由操作员执行 `retire`，systemd 使用 SIGTERM 和 90 秒退出窗口，不使用强杀。`/admin/assets/` 由共享目录提供，按文件摘要拒绝同名不同内容，旧 HTML 引用的资源保留。

首次接入必须监督进行：旧版没有就绪接口，也不能由工具判定后台任务是否排空。先确认旧进程在 `ports[0]`，把旧运行目录和上传目录分开；配置原有站点 server block 包含 `nginx_include`，移除与模板冲突的旧 `location /`、`/api/`、`/admin/assets/` 等代理/静态规则，保留现有域名、TLS、自定义限制和访问日志。旧站点继续运行时执行 `adopt`，工具生成指向原端口的配置并 reload。核对域名、API、登录、上传与旧静态资源后才准备新版本。宝塔站点同样在对应 Nginx 虚拟主机内 include，不修改别的站点。

## 配置与命令

复制 `deploy/smooth-release.example.json` 到主机私有路径，按现场端口、Nginx 二进制/配置/PID、域名、服务用户、目录填写。`probe_origin` 必须是能命中该虚拟主机的本机 HTTP 地址，默认 `http://127.0.0.1`；HTTPS 站点可另开仅回环 HTTP 监听供发布探测。`env_file` 是 systemd 私有环境文件，应包含 `RUOYI_CONFIG_FILE` 指向**外部持久化**的 `application.yaml`；不要把配置或密钥放入发布包。配置中的 `server.host` 必须是 `127.0.0.1`。`shared_uploads` 必须先创建并授予服务用户写权限。系统内存要能同时容纳两个 Go 进程；端口只监听回环。

```sh
python3 deploy/smooth-release.py --config /private/release.json adopt legacy /old/runtime
python3 deploy/smooth-release.py --config /private/release.json prepare release-yyyymmdd /build/ruoyi-go-by.tar.gz ARCHIVE_SHA256 1.1.3
python3 deploy/smooth-release.py --config /private/release.json start release-yyyymmdd
python3 deploy/smooth-release.py --config /private/release.json switch release-yyyymmdd
python3 deploy/smooth-release.py --config /private/release.json status
python3 deploy/smooth-release.py --config /private/release.json retire previous-managed-release
```

`prepare` 只接受经 SHA256 核对的 Linux x86_64 包，还逐文件核对包内 `FILE_SHA256SUMS`，拒绝符号链接、路径穿越和上传数据。运行包由隔离构建节点生成、SSH/SCP 直传并再次校验；正式站不编译。`start` 只启动新实例并检查版本。`switch` 记录恢复状态、reload、核对实际端口和 `smoke_paths`；失败恢复旧路由。中断的切换在下次执行命令时优先恢复旧路由。`rollback` 切回仍运行的上一进程；已 `retire` 的托管进程须先 `start`。旧版首次回退仅用 `/admin/` 探测，须人工监督。切流后先观察真实业务日志和错误率，再排空旧进程；长请求未结束时 `retire` 会拒绝停止。首次 legacy 进程必须确认无后台任务后人工优雅停止，工具故意不自动退休它。工具不自动删除旧版本或共享静态资源。

## 升级边界

`RUOYI_AUTO_SETUP=false` 使候选启动时不执行 `EnsureSystemSettingBaseline`/GORM AutoMigrate；首次安装仍可显式使用旧默认行为。每次数据库改动须先审查并按“新增兼容结构→新旧共同运行→回填→旧版退出后清理”分阶段做。旧浏览器发出的 API 请求必须仍被新后端接受；支付回调、投票等写请求不会被 Nginx 重放。一次切流不覆盖破坏性 schema 或 API 改动，也不能保证客户端在切流后仍持有旧版本页面时的兼容性。

该工具未在任何业务站点自动启用。生产启用前须针对现场 Nginx 配置、上传路径、鉴权/支付回调和较长请求做一次监督演练。`service.sh restart` 已禁用作为生产升级方式。
