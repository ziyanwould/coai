# Linux.do OAuth 单点登录配置指南

## 📋 功能说明

本项目已集成 Linux.do OAuth2 单点登录功能,用户可以使用 Linux.do 账号一键登录 CoAI 系统。

## 🎯 主要特性

- ✅ 标准 OAuth2 授权码流程
- ✅ 自动用户注册/登录
- ✅ CSRF 防护 (state 参数验证)
- ✅ 用户名冲突自动处理
- ✅ 完整的错误处理和用户提示
- ✅ 美观的 UI 设计

## 🔧 配置步骤

### 1. 在 Linux.do 创建 OAuth 应用

访问 https://connect.linux.do/dash/sso/new 创建新应用:

```
应用名: coai
应用主页: https://coai.drawaspark.com
应用描述: coai
回调地址: https://coai.drawaspark.com/callback
```

创建后获取:
- **Client ID**: `iqrFITZB00hqvayAvSsDIhP9cykpjpqG`
- **Client Secret**: `6EgtXxqMWg3PZehaEIIT1JmXLSCfL0ax`

### 2. 配置后端

编辑 `config.yaml` (从 `config.example.yaml` 复制):

```yaml
system:
  oauth:
    linux_do:
      enabled: true
      client_id: "iqrFITZB00hqvayAvSsDIhP9cykpjpqG"
      client_secret: "6EgtXxqMWg3PZehaEIIT1JmXLSCfL0ax"
      redirect_url: "https://coai.drawaspark.com/callback"
```

**重要提示**:
- `redirect_url` 必须与 Linux.do 应用配置中的回调地址完全一致
- 生产环境请使用环境变量或密钥管理系统存储敏感信息

### 3. 构建和部署

**后端构建**:
```bash
go build -o chatnio
```

**前端构建**:
```bash
cd app
npm run build
```

**Docker 部署** (推荐):
```bash
docker build -t chatniolocal .
docker tag chatniolocal:latest harbor.ipv6.liujiarong.top:8024/library/chatniolocal:test
docker push harbor.ipv6.liujiarong.top:8024/library/chatniolocal:test
```

## 📖 使用说明

### 用户登录流程

1. 用户访问 `/login` 登录页面
2. 点击 "使用 Linux.do 登录" 按钮
3. 跳转到 Linux.do 授权页面
4. 用户同意授权后,Linux.do 回调到 `/callback`
5. 系统自动处理:
   - 如果邮箱已存在 → 直接登录
   - 如果是新用户 → 自动注册并登录
6. 登录成功,跳转到首页

### 用户名生成规则

新用户注册时,用户名生成规则如下:

1. 基础用户名: `linux_do_{linux.do用户名}`
2. 如果用户名冲突,添加随机后缀: `linux_do_{用户名}_{随机码}`
3. 极端情况下添加时间戳: `linux_do_{用户名}_{时间戳}`

示例:
- Linux.do 用户名 `zhangsan` → CoAI 用户名 `linux_do_zhangsan`
- 如果冲突 → `linux_do_zhangsan_A3B9`

### 安全机制

**CSRF 防护**:
- 使用 `state` 参数防止跨站请求伪造攻击
- state 在 Redis 中存储 10 分钟,验证后立即删除

**用户追踪**:
- 登录成功后,`token` 字段存储 `linux_do:{user_id}`
- 可用于区分 OAuth 用户和普通注册用户

## 🗂️ 文件变更列表

### 后端文件

| 文件 | 说明 |
|------|------|
| `auth/oauth.go` | OAuth 核心逻辑 (新建) |
| `auth/router.go` | 添加 OAuth 路由 |
| `config.example.yaml` | 添加 OAuth 配置示例 |

### 前端文件

| 文件 | 说明 |
|------|------|
| `app/src/routes/OAuthCallback.tsx` | OAuth 回调页面 (新建) |
| `app/src/routes/Auth.tsx` | 添加 OAuth 登录按钮 |
| `app/src/router.tsx` | 添加回调路由 |
| `app/src/assets/pages/auth.less` | 添加 OAuth 样式 |

## 🔍 API 端点

### 后端 API

- `GET /api/oauth/linux-do/login` - 跳转到 Linux.do 授权页
  - 返回: `{"status": true, "url": "https://connect.linux.do/oauth2/authorize?..."}`

- `GET /api/callback` - OAuth 回调处理
  - 参数: `code`, `state`, `error`
  - 成功: 重定向到 `/oauth-success?token={jwt_token}`
  - 失败: 重定向到 `/login?error={error_message}`

### 前端路由

- `/login` - 登录页面 (包含 OAuth 按钮)
- `/oauth-success` - OAuth 回调处理页面

## 🛠️ 故障排除

### 常见问题

**1. 点击登录按钮无反应**
- 检查 `config.yaml` 中 `enabled: true`
- 检查 `client_id` 和 `client_secret` 是否正确配置

**2. 回调时提示 "invalid_state"**
- Redis 连接问题,检查 Redis 服务状态
- state 可能已过期 (10分钟),重新登录

**3. 提示 "本站暂未开放注册"**
- 检查 `globals.CloseRegistration` 配置
- 如需允许 OAuth 用户注册,需要在代码中添加例外逻辑

**4. 用户名冲突**
- 系统会自动添加随机后缀,无需手动处理

### 调试日志

OAuth 相关日志前缀: `[oauth]`

查看日志:
```bash
# 应用日志
tail -f logs/app.log | grep oauth

# 错误日志
tail -f logs/error.log | grep oauth
```

## 📊 数据库变更

**无需修改数据库结构**

OAuth 用户信息复用现有 `auth` 表字段:
- `username` - 自动生成的唯一用户名
- `email` - Linux.do 邮箱
- `token` - 存储 `linux_do:{user_id}` 标识
- `password` - 随机生成的哈希密码 (64位)

## 🔐 安全建议

1. **生产环境配置**:
   - 使用环境变量存储敏感信息
   - 启用 HTTPS (OAuth 必须使用 HTTPS)
   - 定期轮换 Client Secret

2. **监控和日志**:
   - 监控 OAuth 登录失败率
   - 记录异常的 state 验证失败
   - 追踪用户注册来源

3. **用户隐私**:
   - 仅请求必要的 OAuth scope (openid, profile, email)
   - 明确告知用户数据使用方式

## 📚 参考资料

- [Linux.do OAuth 文档](https://connect.linux.do/dash/sso)
- [OAuth 2.0 RFC](https://datatracker.ietf.org/doc/html/rfc6749)
- [golang.org/x/oauth2 文档](https://pkg.go.dev/golang.org/x/oauth2)

## 🎉 完成检查清单

- [x] 后端 OAuth 逻辑实现
- [x] 前端 UI 集成
- [x] 配置文件更新
- [x] 路由注册
- [x] 样式适配
- [x] 编译测试通过
- [x] 安全机制实现 (CSRF 防护)
- [x] 错误处理完善
- [x] 使用文档编写

---

**版本**: v1.0
**更新日期**: 2025-10-20
**开发者**: Claude Code AI Assistant
