# RuoYi-Go BY 管理前端

项目保留两套可选管理前端，二者共享同一套后端 API、动态菜单和权限模型。

| 前端 | 对应版本 | 目录 | 状态 | 适用场景 |
| --- | --- | --- | --- | --- |
| Whiteyun Vue（白云风格前端） | `v1.2.0` 起 | `whiteyun-vue` | 默认、持续维护 | 新部署、生产环境和后续功能开发 |
| RuoYi Vue | `v1.1.2` 基线 | `ruoyi-vue` | 兼容保留、停止常规维护 | 需要旧若依视觉或旧二次开发兼容 |

Whiteyun Vue 不是从零替换技术栈：它继续使用 RuoYi Vue 3 的 Vue 3、Vite、
Element Plus、动态路由和 RBAC 基础，但已经形成独立的白云视觉、主题与响应式体系。

## 构建选择

默认构建 Whiteyun Vue：

```bash
frontend/build-admin.sh
```

显式构建旧版 RuoYi Vue：

```bash
APIAUTH_ADMIN_FRONTEND=ruoyi-vue frontend/build-admin.sh
```

允许值只有 `whiteyun-vue` 和 `ruoyi-vue`。完整发布脚本使用相同的
`APIAUTH_ADMIN_FRONTEND` 变量；未指定时始终选择 `whiteyun-vue`。

为避免旧脚本和已有二次开发立即失效，历史路径 `RuoYi-Vue3-ts` 暂时作为兼容链接
指向 `whiteyun-vue`。它不代表冻结旧版；需要旧视觉时必须显式选择 `ruoyi-vue`。
