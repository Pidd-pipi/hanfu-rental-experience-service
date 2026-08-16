# Hanfu Rental 汉服租赁与文化体验平台后端

汉服租赁与文化体验平台后端 API，提供汉服目录、会员卡、租赁订单、活动、文章与用户押金等能力。

## 技术栈
- Go 1.22
- Gin + GORM
- PostgreSQL 15
- JWT + RBAC

## 标准命令
```bash
go build ./...        # 编译
go test ./...         # 运行测试
go run ./cmd/server   # 启动 HTTP 服务
```

## 环境变量
| 变量 | 说明 | 默认值 |
| --- | --- | --- |
| DB_HOST | PostgreSQL 主机 | localhost |
| DB_PORT | PostgreSQL 端口 | 5432 |
| DB_NAME | 数据库名 | hanfu_rental |
| DB_USER | 数据库用户 | hanfu_user |
| DB_PASSWORD | 数据库密码 | hanfu_pwd |
| JWT_SECRET | JWT 签名密钥 | change_me_to_a_long_random_string |
| APP_CORS_ORIGINS | 允许跨域来源 | http://localhost:28504 |
