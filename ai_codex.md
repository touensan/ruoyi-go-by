# RuoYi-Go BY 长期记忆

公开仓库：`touensan/ruoyi-go-by`。后端 Go/Gin/GORM，MySQL/Redis；前端为 `frontend/RuoYi-Vue3-ts` 的原生若依 Vue 3。真实部署配置不纳入仓库。开始工作先阅读本文件与 `docs/ARCHITECTURE.md`，完成关键阶段后及时记录验证结果。

## 2026-09-05｜统一后台与角色继承

用户明确要求整个 `ruoyi-xxx-by` 系列不拆分用户端和管理端，采用同一后台及 RBAC。`admin` 单向继承启用的 `common` 功能授权；超级管理员全权限与显式角色数据范围保持原边界。生效授权和显式角色分配分开，前端、菜单/路由及接口规则一致。详细契约见 `docs/ARCHITECTURE.md`；实现及验证已完成，详细结果与命令见 `docs/DEVELOPMENT.md`。

## 2026-09-05｜版本与功能说明纠正

当前源码与旧 Release/附件分开记录在 `docs/VERSIONS.md`。历史发布说明仅适用于其标签；未真实供应商联调不等于未实现。以后每次改变功能必须同步版本状态与文档，旧标签和下载附件不得被文案修改冒充为新版本。

## 2026-09-08｜v1.1.3 原生主线 Release

用户要求发布最新 main 对应下载包。为延续 v1.1.2 原生若依主线并避免与 v1.2.0–v1.4.0 Whiteyun 历史支线混淆，新版本定为 v1.1.3。发布前更新 Axios、ECharts、js-cookie 与 glob，以清除本轮发现的前端生产依赖公告。运行包包含 Linux x86_64 后端、构建后的 Vue 前端、示例配置和数据库脚本；包内 `BUILD_INFO.json` 记录源码与产物摘要。实际发布地址、附件 SHA-256 和本轮验证结果以 GitHub Release 与 `docs/DEVELOPMENT.md` 最新记录为准；本次不部署业务站。

## 2026-09-24｜单机平滑更新源码

用户要求参考 MOX 项目的蓝绿发布方式，改造 Go/Rust 公开仓库；PHP 仓库仅拉取核对。当前 Go main 已加入回环私有就绪检查、SIGTERM 优雅退出、候选端口覆盖及关闭候选自动建表的开关；发布工具与兼容边界见 `docs/SMOOTH_RELEASE.md`。旧 v1.1.3 附件未改，新流程尚未在业务站点部署。验证结论见 `docs/DEVELOPMENT.md`，不得把 Python 工具自测等同于正式生产切流验收。
