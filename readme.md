# RuoYi-Go BY

一个基于 Go、Gin、GORM 和 Vue 3 的 RuoYi 风格后台管理系统示例。

本仓库是一个可独立运行的开源代码副本，默认配置仅用于本地开发。部署前请修改数据库、Redis、JWT 密钥和所有默认账号配置。

## 技术栈

- 后端：Go 1.25+、Gin、GORM、MySQL、Redis、JWT、excelize
- 前端：Vue 3、Vite、TypeScript、Element Plus、Pinia、Vue Router、Axios

## 后端运行

```bash
cp application-example.yaml application.yaml
# 编辑 application.yaml，填写本地数据库和 Redis 配置
go mod tidy
go run main.go
```

## 前端运行

```bash
cd frontend/RuoYi-Vue3-ts
cp .env.example .env.development
npm install
npm run dev
```

## 构建

```bash
go build -o ruoyi-go-by main.go
cd frontend/RuoYi-Vue3-ts
npm run build:prod
```

## 数据库

`ruoyi.sql` 包含表结构和本地开发所需的示例初始化数据。请不要把真实用户、支付、邮件或生产配置写入 SQL 文件。

## 安全说明

- `application.yaml` 已被 Git 忽略，请勿提交。
- 生产环境必须替换示例 JWT 密钥、数据库密码、Redis 密码和默认账号密码。
- 发现安全问题请参考 [SECURITY.md](SECURITY.md)。

## 许可证

本项目使用 MIT License。第三方依赖仍以各自许可证为准。
