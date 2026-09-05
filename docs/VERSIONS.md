# 源码、标签与下载包状态

核对日期：2026-09-05。功能描述按所属版本解读，历史发布说明不代表当前源码。

当前 `main` 是若依原生 Vue 3 主线，包含统一后台 RBAC；没有发布包含这些最新变更的新标签或运行包。`application-example.yaml` 中的 `1.1.2` 是可配置的基线版本字段，不能证明工作目录等同于 `v1.1.2` 标签。

| 来源 | 实际含义 |
| --- | --- |
| `main` | 当前源码：原生若依统一后台、RBAC 与明确的未实现接口响应 |
| `v1.1.2` | 历史原生基线；源码 ZIP，不含后续 RBAC 和能力状态修正 |
| `v1.2.0` / `v1.3.0` / `v1.4.0` | Whiteyun Vue 历史支线；功能、目录和升级脚本不代表当前 main |

## 获取当前功能

使用 [main 源码](https://github.com/touensan/ruoyi-go-by/tree/main)，记录 `git rev-parse HEAD` 的完整提交号，并按 [README](../README.md) 构建/安装。已有部署升级前先备份并阅读迁移说明。

[Releases](https://github.com/touensan/ruoyi-go-by/releases) 中的附件和 GitHub 自动生成的 Source code ZIP/TAR 均对应其标签，不能因发布页位于列表首位、标记 Latest 或版本号较大，就当成 main 的最新代码。

本次更正文档和发布说明，不移动旧标签、不替换旧包、不改变其校验和。旧包内的文档也属于该版本的历史快照；当前状态请以本页为准。

## 已核对的标签与附件

| 标签 | 对应源码提交 | 附件 |
| --- | --- | --- |
| [v1.1.2](https://github.com/touensan/ruoyi-go-by/releases/tag/v1.1.2) | [`d13494d`](https://github.com/touensan/ruoyi-go-by/tree/d13494d0584da7808621b00051e466d1234b9489) | ruoyi-go-by-v1.1.2-source.zip |
| [v1.2.0](https://github.com/touensan/ruoyi-go-by/releases/tag/v1.2.0) | [`17e4b5f`](https://github.com/touensan/ruoyi-go-by/tree/17e4b5f0e610f9204f3806c3ffa994622b0205e7) | 无附加资产；可下载 GitHub 标签源码 |
| [v1.3.0](https://github.com/touensan/ruoyi-go-by/releases/tag/v1.3.0) | [`b1326c1`](https://github.com/touensan/ruoyi-go-by/tree/b1326c141c6980e6a0bc337ea643cc76f26e6941) | 无附加资产；可下载 GitHub 标签源码 |
| [v1.4.0](https://github.com/touensan/ruoyi-go-by/releases/tag/v1.4.0) | [`a5965a9`](https://github.com/touensan/ruoyi-go-by/tree/a5965a9dc30f085f405963f0c2e257dc9e7c96d1) | 无附加资产；可下载 GitHub 标签源码 |

## 维护约定

新增、移植或撤销功能时，同步 README、FEATURES、开发记忆与本页。发布新包时必须记录源码提交和校验和，并确认二进制、前端与该提交一致。未发布新包时明确标注“源码已实现，旧下载包不包含”，不能把“未真实供应商联调”写成“未实现”，也不能把旧版本局限写成整个语言框架的当前局限。
