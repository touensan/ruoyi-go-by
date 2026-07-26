# RuoYi Vue 兼容管理端

这是 APIAuth 在“白云风格前端”改造前保留的完整若依视觉版本，包含当时已有的授权
管理业务页面。

本目录自 2026-07-26 起进入冻结维护：

- 保留给需要原若依视觉或已有二次开发兼容性的使用者；
- 不再同步 Whiteyun Vue 的常规功能、视觉与体验更新；
- 仅在出现影响构建、严重安全或迁移可用性的缺陷时考虑修复；
- 新项目和生产发布默认使用 `../whiteyun-vue`。

## 运行

```bash
cp .env.example .env.development
npm ci
npm run dev
```

## 构建

```bash
npm run build:prod
```

后端 API 默认通过 `/dev-api` 代理，请根据本地后端地址调整 Vite 配置。
