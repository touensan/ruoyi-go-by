# Whiteyun Vue 管理端

RuoYi-Go BY 默认管理前端，中文名称为“白云风格前端”。

它仍然继承 RuoYi Vue 3 的 Vue、Vite、Element Plus、动态路由与 RBAC 工程基础，
但视觉体系、主题变量、导航矩阵、响应式布局和登录页已经由本项目独立
维护。后续功能、兼容性和界面更新均以本目录为准。

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

旧版若依视觉保留在 `../ruoyi-vue`，仅作为兼容选择，不再接收常规功能和视觉更新。
