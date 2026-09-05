# RuoYi-Go BY 统一后台前端

普通用户与管理员使用同一套前端、登录和页面，按 RBAC 动态授权。当前功能与下载版本见 [项目说明](../../README.md) 和 [版本状态](../../docs/VERSIONS.md)。

基于 Vue 3、TypeScript、Vite 和 Element Plus 的统一后台前端。

## 运行

```bash
cp .env.example .env.development
npm install
npm run dev
```

## 构建

```bash
npm run build:prod
```

后端 API 默认通过 `/dev-api` 代理，请根据本地后端地址调整 Vite 配置。
