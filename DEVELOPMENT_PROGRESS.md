# ChatNio 开发建设进度文档

## 项目概述
ChatNio 是一个新一代 AIGC 一站式商业解决方案，结合强大的 API 分发系统和丰富的用户界面设计，为 B2B 和 B2C 用户提供 AI 聊天服务。

## 技术栈

### 后端技术栈
- **核心语言**: Go 1.20
- **Web 框架**: Gin
- **数据库**:
  - MySQL (主数据库)
  - Redis (缓存和会话管理)
- **配置管理**: Viper
- **日志系统**: Logrus + Lumberjack
- **认证**: JWT
- **文档处理**:
  - PDF/DOCX/PPTX/Excel 解析
  - 图片处理 (WebP, PNG, JPG)
- **AI 模型集成**:
  - OpenAI GPT 系列
  - Anthropic Claude
  - Google Gemini
  - Cloudflare Workers AI
  - SiliconFlow AI
  - 火山引擎豆包
  - 10+ 其他 AI 提供商

### 前端技术栈
- **核心框架**: React 18.3.1
- **构建工具**: Vite 4.5.3
- **包管理**: pnpm
- **UI 库**:
  - Radix UI (核心组件)
  - Tailwind CSS (样式)
  - Tremor React (图表)
  - Lucide React (图标)
- **状态管理**: Redux Toolkit
- **路由**: React Router DOM
- **Markdown 渲染**:
  - React Markdown
  - Remark/Rehype 插件生态
  - KaTeX (数学公式)
- **国际化**: i18next + react-i18next
- **主题**: next-themes
- **特殊功能**:
  - Mermaid 图表渲染
  - 语法高亮 (react-syntax-highlighter)
  - HTML5 Canvas (图像编辑)
  - PWA 支持
  - WebSocket (实时通信)

### 基础设施
- **容器化**: Docker + Docker Compose
- **镜像仓库**: Harbor
- **开发环境**:
  - nodemon (后端热重载)
  - Vite HMR (前端热重载)

## 快速开始指南

### 环境要求

- **Node.js** >= 16.0.0
- **Go** >= 1.19
- **MySQL** >= 8.0
- **Redis** >= 6.0
- **Docker** & **Docker Compose**（推荐）

### 本地开发环境搭建

#### 1. 克隆项目

```bash
git clone https://github.com/Deeptrain-Community/chatnio.git
cd coai
```

#### 2. 使用 Docker Compose（推荐）

```bash
# 启动所有服务
docker-compose up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f
```

#### 3. 手动安装

**前端安装：**
```bash
cd app
npm install -g pnpm
pnpm install
pnpm dev
```

**后端安装：**
```bash
# 配置数据库连接
cp config/config.example.yaml config/config.yaml
# 编辑 config/config.yaml 配置数据库信息

# 编译运行
go build -o chatnio
./chatnio
```

### 项目结构

```
coai/
├── app/                    # 前端代码
│   ├── public/            # 静态资源
│   ├── src/               # 源代码
│   │   ├── components/    # React 组件
│   │   ├── pages/         # 页面组件
│   │   ├── store/         # Redux store
│   │   ├── utils/         # 工具函数
│   │   └── assets/        # 资源文件
│   ├── package.json       # 前端依赖
│   └── tailwind.config.js # Tailwind 配置
├── adapter/               # AI 模型适配器
├── admin/                 # 管理后台
├── auth/                  # 认证模块
├── channel/               # 渠道管理
├── config/                # 配置文件
│   └── config.yaml        # 主配置文件
├── docker/               # Docker 相关文件
├── utils/                 # 工具函数
├── go.mod                # Go 依赖
├── main.go               # Go 入口文件
├── docker-compose.yaml   # Docker Compose 配置
└── README.md            # 项目说明
```

## 开发规范

### 代码规范

#### 前端代码规范
- 使用 **ESLint** 和 **Prettier** 进行代码格式化
- 组件命名使用 **PascalCase**
- 文件命名使用 **kebab-case**
- 变量命名使用 **camelCase**
- 常量命名使用 **UPPER_CASE**

```javascript
// 组件示例
const ChatComponent = () => {
  const [messages, setMessages] = useState([]);

  return (
    <div className="chat-container">
      {/* 组件内容 */}
    </div>
  );
};
```

#### 后端代码规范
- 使用 **gofmt** 进行代码格式化
- 包命名使用 **lowercase**
- 函数命名使用 **CamelCase**
- 变量命名使用 **camelCase**
- 常量命名使用 **UPPER_SNAKE_CASE**

```go
// 控制器示例
func HandleChatRequest(c *gin.Context) {
    var req ChatRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }

    // 处理逻辑
    c.JSON(http.StatusOK, gin.H{"data": response})
}
```

### 分支管理

- `main` - 主分支，用于生产环境
- `develop` - 开发分支
- `feature/*` - 功能分支
- `hotfix/*` - 热修复分支
- `release/*` - 发布分支

### 提交规范

使用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

类型说明：
- `feat`: 新功能
- `fix`: 修复 bug
- `docs`: 文档更新
- `style`: 代码格式（不影响功能）
- `refactor`: 重构
- `test`: 测试相关
- `chore`: 构建过程或辅助工具的变动

示例：
```
feat(chat): add message pagination support

- Add pagination controls to chat interface
- Implement backend API for message pagination
- Update database schema for message indexing

Closes #123
```

### 开发流程

1. **创建功能分支**
```bash
git checkout develop
git pull origin develop
git checkout -b feature/your-feature-name
```

2. **开发和测试**
```bash
# 运行前端开发服务器
cd app && pnpm dev

# 运行后端开发服务器
go run main.go

# 运行测试
pnpm test        # 前端测试
go test ./...    # 后端测试
```

3. **提交代码**
```bash
git add .
git commit -m "feat: add new feature"
git push origin feature/your-feature-name
```

4. **创建 Pull Request**
- 在 GitHub/GitLab 上创建 PR
- 填写 PR 模板
- 等待代码审查
- 合并到 `develop` 分支

## 配置管理

### 环境变量

主要环境变量：

```bash
# 数据库配置
MYSQL_HOST=localhost
MYSQL_PORT=3306
MYSQL_DB=chatnio
MYSQL_USER=root
MYSQL_PASSWORD=your_password

# Redis 配置
REDIS_HOST=localhost
REDIS_PORT=6379

# 应用配置
SECRET=your_jwt_secret_key
SERVE_STATIC=true
PORT=8094
```

### 配置文件

主配置文件 `config/config.yaml`：

```yaml
app:
  name: "CoAI"
  version: "1.0.0"
  port: 8094
  debug: true

mysql:
  host: "${MYSQL_HOST}"
  port: "${MYSQL_PORT}"
  database: "${MYSQL_DB}"
  username: "${MYSQL_USER}"
  password: "${MYSQL_PASSWORD}"

redis:
  host: "${REDIS_HOST}"
  port: "${REDIS_PORT}"

jwt:
  secret: "${SECRET}"
  expire: 24h
```

## 测试

### 前端测试

```bash
cd app

# 运行所有测试
pnpm test

# 运行测试并生成覆盖率报告
pnpm test:coverage

# 监听模式运行测试
pnpm test:watch
```

### 后端测试

```bash
# 运行所有测试
go test ./...

# 运行测试并生成覆盖率报告
go test -cover ./...

# 运行特定包的测试
go test ./pkg/handler
```

### API 测试

使用 Postman 或 curl 测试 API：

```bash
# 健康检查
curl http://localhost:8094/api/health

# 用户登录
curl -X POST http://localhost:8094/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"username":"root","password":"chatnio123456"}'
```

## 部署

### Docker 部署

```bash
# 构建镜像
docker build -t coai:latest .

# 运行容器
docker run -d --name coai \
  --network host \
  -v ~/config:/config \
  -v ~/logs:/logs \
  coai:latest
```

### Docker Compose 部署

```bash
# 生产环境部署
docker-compose -f docker-compose.prod.yaml up -d

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f coai
```

### 生产环境配置

1. **反向代理配置（Nginx）**
```nginx
server {
    listen 80;
    server_name your-domain.com;

    location / {
        proxy_pass http://localhost:8094;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
    }

    location /ws {
        proxy_pass http://localhost:8094;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }
}
```

2. **SSL 证书配置**
```bash
# 使用 Let's Encrypt
certbot --nginx -d your-domain.com
```

## 故障排除

### 常见问题

1. **数据库连接失败**
```bash
# 检查 MySQL 服务状态
systemctl status mysql

# 检查连接配置
mysql -h localhost -u root -p

# 查看错误日志
tail -f logs/error.log
```

2. **前端构建失败**
```bash
# 清除缓存
pnpm store prune

# 重新安装依赖
rm -rf node_modules package-lock.json
pnpm install
```

3. **Redis 连接问题**
```bash
# 检查 Redis 服务
redis-cli ping

# 查看连接配置
redis-cli -h localhost -p 6379 info
```

### 日志查看

```bash
# 应用日志
tail -f logs/app.log

# 错误日志
tail -f logs/error.log

# Docker 日志
docker-compose logs -f coai
```

## 性能优化

### 前端优化
- 使用 React.memo 优化组件渲染
- 实现虚拟滚动处理大列表
- 使用 Web Workers 处理重计算
- 启用 Gzip 压缩

### 后端优化
- 使用 Redis 缓存热点数据
- 优化数据库查询和索引
- 实现连接池管理
- 使用 goroutine 处理并发

## 安全最佳实践

1. **输入验证** - 严格验证所有用户输入
2. **SQL 注入防护** - 使用参数化查询
3. **XSS 防护** - 对输出进行转义
4. **CSRF 防护** - 使用 CSRF token
5. **密码安全** - 使用 bcrypt 加密存储
6. **JWT 安全** - 设置合理的过期时间

## 贡献指南

1. Fork 项目仓库
2. 创建功能分支
3. 编写代码和测试
4. 确保所有测试通过
5. 提交 Pull Request
6. 等待代码审查

## 主要功能模块

### 1. 核心聊天系统
- **多模型支持**: 支持 15+ AI 模型提供商
- **流式传输**: WebSocket 实时流式响应
- **会话管理**: 跨设备同步，本地存储
- **消息格式**: 支持文本、图片、文件等多媒体输入

### 2. 图像生成与编辑 (重点开发)
- **基础图像生成**:
  - 文生图 (Text-to-Image)
  - 图生图 (Image-to-Image)
- **高级图像编辑**:
  - **Inpainting (图像修复)**: 交互式遮罩绘制
  - **多格式支持**: Cloudflare Workers AI 的 JSON/Binary 响应格式
  - **模型优化**: Flux, Leonardo, Phoenix 系列模型
- **交互式工具**:
  - HTML5 Canvas 绘制界面
  - 画笔/橡皮擦工具
  - 实时预览和调整

### 3. 文件处理系统
- **支持格式**: PDF, DOCX, PPTX, Excel, 图片
- **智能解析**: 文本提取、结构化处理
- **多媒体集成**: 图片和文档混合处理

### 4. 管理后台
- **用户管理**: 注册、认证、权限控制
- **渠道管理**: AI 模型提供商配置和负载均衡
- **计费系统**: 订阅制和弹性计费模式
- **模型市场**: 自定义模型市场配置

### 5. 企业级功能
- **多租户支持**: B2B 和 B2C 双模式
- **API 兼容性**: OpenAI 兼容的 API 端点
- **高可用性**: 自动故障转移、负载均衡
- **监控日志**: 详细的使用统计和错误追踪

## 近期重大开发进展 (2025-09)

### Cloudflare Workers AI 图像生成增强
**开发时间**: 2025年9月17日
**开发者**: Claude Code AI Assistant

#### 问题背景
用户报告三个 Cloudflare 图像模型无法正常显示生成的图片：
- `@cf/leonardo/lucid-origin`
- `@cf/leonardo/phoenix-1.0`
- `@cf/black-forest-labs/flux-1-schnell`

#### 技术调查与解决

**1. API 响应格式分析**
- **发现**: 不同模型返回不同的响应格式
  - Flux/Leonardo 模型: JSON 格式 `{"result":{"image":"<base64>"}}`
  - Phoenix 模型: 直接返回二进制数据
  - 原代码只处理二进制格式

**2. 后端架构改进**
- **文件**: `adapter/cloudflare/image.go`, `adapter/cloudflare/types.go`
- **实现**: 双格式检测和处理机制
```go
if len(buffer) > 0 && buffer[0] == '{' {
    // JSON 格式处理
    var nestedResp ImageResponseNested
    // ... 处理逻辑
} else {
    // 二进制格式处理 (原有逻辑)
}
```

**3. Inpainting 交互式功能开发**

**需求**: 用户反馈需要更友好的图像修复工具，不想手动上传遮罩图像

**技术实现**:
- **前端组件开发**:
  - `InpaintingCanvas.tsx`: 完整的 Canvas 绘制工具
    - HTML5 Canvas API
    - 画笔/橡皮擦工具切换
    - 可调节笔刷大小 (5-50px)
    - 实时绘制预览
    - 遮罩生成算法
  - `InpaintingTrigger.tsx`: 触发组件
    - 自动检测单图像 + inpainting 模型
    - 一键启动绘制界面
    - 结果自动填充到输入框

- **后端逻辑增强**:
  - 智能检测 inpainting 模型
  - 单图像自动触发绘制界面
  - 生成特殊 Markdown 响应格式
```go
return fmt.Sprintf(`请使用涂鸦工具标记要修复的区域：
![需要修复的图片](%s)
<div class="inpainting-trigger" data-image="%s" data-model="%s" data-prompt="%s">
点击开始涂鸦标记
</div>`, inputImage, inputImage, props.Model, prompt), nil
```

- **Markdown 渲染集成**:
  - 扩展 `Markdown.tsx` 支持自定义 div 组件
  - 条件启用 HTML 渲染
  - 数据属性传递

### 关键 Bug 修复

**1. 重复图像提取问题**
- **问题**: `utils/char.go` 中的 `ExtractBase64Images` 函数对同一图像进行重复提取
- **原因**: 两个正则表达式模式匹配到相同内容
  - Pattern 1: 直接 base64 匹配
  - Pattern 2: 文件块内 base64 匹配
- **解决**: 添加去重逻辑
```go
// Remove duplicates
uniqueMatches := make([]string, 0)
seen := make(map[string]bool)
for _, match := range cleanedMatches {
    if !seen[match] {
        seen[match] = true
        uniqueMatches = append(uniqueMatches, match)
    }
}
```

**2. API 参数名称错误**
- **问题**: Cloudflare Inpainting API 要求 `mask_image` 参数，但代码使用 `mask_b64`
- **解决**: 创建专用的 `InpaintingRequest` 结构体
```go
type InpaintingRequest struct {
    // ... 其他字段
    ImageB64  string `json:"image_b64"`
    MaskImage string `json:"mask_image"`   // 正确的参数名
}
```

**3. 大小写敏感问题**
- **问题**: 模型名称检测时大小写敏感导致 inpainting 检测失败
- **解决**: 统一转换为小写进行比较
```go
isInpainting := strings.Contains(strings.ToLower(props.Model), "inpainting")
```

### 用户追踪功能增强 (2025-09-17)

**开发背景**: 为了让上游中间件能够正确识别用户和收集统计信息，需要在 Cloudflare 图像生成请求中添加用户标识字段。

**技术实现**:

**1. 结构体字段添加**
- **文件**: `adapter/cloudflare/types.go`
- **修改**: 为 `ImageRequest` 和 `InpaintingRequest` 添加用户字段
```go
type ImageRequest struct {
    // ... 其他字段
    User      interface{} `json:"user,omitempty"`      // 用户标识
    Userip    string      `json:"user_ip,omitempty"`   // 用户IP
}

type InpaintingRequest struct {
    // ... 其他字段
    User      interface{} `json:"user,omitempty"`      // 用户标识
    Userip    string      `json:"user_ip,omitempty"`   // 用户IP
}
```

**2. 数据传递链路**
- **文件**: `adapter/cloudflare/image.go`
- **修改**: 完善数据传递链路
```go
// ImageProps 结构体添加字段
type ImageProps struct {
    // ... 其他字段
    User       interface{}            // 用户标识
    Userip     string                 // 用户IP
}

// CreateImage 函数传递用户信息
base64Data, err := c.CreateImageRequest(ImageProps{
    // ... 其他字段
    User:       props.User,     // 从 ChatProps 获取
    Userip:     props.Ip,       // 从 ChatProps 获取
})

// CreateImageRequest 函数设置请求体
requestBody := ImageRequest{
    // ... 其他字段
    User:     props.User,
    Userip:   props.Userip,
}
```

**3. 上游兼容性**
- **JSON 输出格式**:
```json
{
  "prompt": "用户提示词",
  "user": "用户标识",
  "user_ip": "用户IP地址"
}
```
- **上游获取方式**:
```javascript
const ip = req.body.user_ip || req.headers['x-user-ip'] || req.ip;
const userId = req.body.user || req.headers['x-user-id'];
```

**4. 一致性设计**
- 参考 OpenAI adapter 的实现模式
- 保持字段名称和数据类型一致
- 确保所有 AI 提供商的用户追踪功能统一

## SiliconFlow AI 图像生成集成 (2025-09-17)

**开发背景**: 应用户需求，集成 SiliconFlow AI 的图像生成模型，包括文生图和图像编辑功能，进一步丰富平台的 AI 图像生成能力。

### 技术实现

**1. 完整适配器开发**
- **文件结构**: 遵循 Cloudflare 适配器模式
  - `adapter/siliconflow/types.go`: 数据结构定义
  - `adapter/siliconflow/image.go`: 核心图像生成逻辑
  - `adapter/siliconflow/chat.go`: 流式聊天接口
- **集成注册**: 在 `adapter/adapter.go` 中注册 SiliconFlow 工厂

**2. 支持的模型**
```go
var SiliconFlowImageModels = []string{
    "Qwen/Qwen-Image",        // 文生图模型
    "Qwen/Qwen-Image-Edit",   // 图像编辑模型
    "Kwai-Kolors/Kolors",     // 高级文生图模型
}
```

**3. 技术特性**

**图像编辑支持**:
- `Qwen/Qwen-Image-Edit` 注册为视觉模型，支持图片上传
- 自动检测图像编辑模型，要求必须提供输入图片
- 完整的 data URI 格式支持

**模型参数优化**:
```go
// Qwen 模型专用配置
if strings.Contains(strings.ToLower(props.Model), "qwen") {
    if strings.Contains(strings.ToLower(props.Model), "edit") {
        cfg = 4.0 // 图像编辑默认 CFG
        if 高质量提示 {
            cfg = 6.0 // 提升质量时使用更高 CFG
        }
    } else {
        cfg = 7.5 // 常规 Qwen 模型 CFG
    }
}

// Kolors 模型专用配置
if strings.Contains(strings.ToLower(props.Model), "kolors") {
    guidanceScale = 5.0 // 不同的引导比例
    steps = 25          // 增加推理步数
}
```

**智能尺寸检测**:
```go
// 根据提示词自动调整图像尺寸
if strings.Contains(lowerPrompt, "宽屏") || strings.Contains(lowerPrompt, "landscape") {
    imageSize = "1280x960"
}
if strings.Contains(lowerPrompt, "竖屏") || strings.Contains(lowerPrompt, "portrait") {
    imageSize = "960x1280"
}
```

**4. 用户追踪集成**
```go
type ImageRequest struct {
    Model              string      `json:"model"`
    Prompt             string      `json:"prompt"`
    Image              string      `json:"image,omitempty"`     // 图像编辑输入
    User               interface{} `json:"user,omitempty"`      // 用户标识
    UserIP             string      `json:"user_ip,omitempty"`   // 用户IP
    // ... 其他参数
}
```

**5. 错误处理增强**
- 详细的调试日志输出
- 图像处理错误的中文提示
- Base64 图像数据完整性检查

**6. 前端管理界面**
- **文件**: `app/src/admin/channel.ts`
- **配置**: 完整的 SiliconFlow 渠道类型定义
```typescript
siliconflow: {
    endpoint: "https://api.siliconflow.cn/v1",
    format: "<api-key>",
    models: ["Qwen/Qwen-Image", "Qwen/Qwen-Image-Edit", "Kwai-Kolors/Kolors"],
    usage: `
1. 获取 SiliconFlow API Key
2. 模型说明：
   - Qwen/Qwen-Image: 通义千问文生图模型
   - Qwen/Qwen-Image-Edit: 图像编辑模型（需要上传图片）
   - Kwai-Kolors/Kolors: 快手可图大模型
    `,
}
```

### 开发流程与验证

**1. 构建测试**
```bash
# 后端编译验证
go build -o chatnio

# 前端构建验证
cd app && npm run build
```

**2. Docker 镜像构建**
```bash
# 构建本地镜像
docker build -t chatniolocal .

# 标记并推送到 Harbor
docker tag chatniolocal:latest harbor.ipv6.liujiarong.top:8024/library/chatniolocal:test
docker push harbor.ipv6.liujiarong.top:8024/library/chatniolocal:test
```

**3. 功能验证**
- ✅ 后端 Go 代码编译通过
- ✅ 前端 TypeScript 编译通过
- ✅ 三个模型配置正确加载
- ✅ 管理面板渠道类型显示正常
- ✅ 用户追踪字段正确传递

### 架构设计亮点

**1. 代码复用性**
- 严格遵循现有 Cloudflare 适配器设计模式
- 统一的用户追踪实现
- 一致的错误处理机制

**2. 扩展性考虑**
- 模块化的模型配置
- 灵活的参数优化策略
- 支持未来新模型的快速接入

**3. 用户体验优化**
- 智能提示词尺寸检测
- 中文错误提示
- 详细的使用说明文档

### 技术亮点

**1. 错误追踪与调试**
- 实施了全面的调试日志系统
- API 请求/响应详细记录
- 参数验证和错误上下文

**2. 向后兼容性**
- 保持原有图像生成功能完全不变
- 渐进式增强用户体验
- 新功能可选启用

**3. 性能优化**
- base64 图像数据去重
- 条件 HTML 渲染 (仅在需要时启用)
- Canvas 操作优化

### 开发工具与流程

**1. 构建和部署**
```bash
# 前端构建
cd app && npm run build

# Docker 镜像构建
docker build -t chatniolocal .

# 推送到 Harbor 仓库
docker tag chatniolocal:latest harbor.ipv6.liujiarong.top:8024/library/chatniolocal:test
docker push harbor.ipv6.liujiarong.top:8024/library/chatniolocal:test
```

**2. 测试策略**
- 本地开发服务器实时测试
- Docker 容器化部署测试
- Harbor 仓库版本管理

## Markdown 视频播放功能集成 (2025-10-17)

**开发背景**: 用户需求在 Markdown 渲染的消息中自动显示视频播放器，实现类似图片的内嵌播放体验，而非仅显示视频链接。

### 技术实现

**1. 视频播放器组件开发**
- **文件**: `app/src/components/plugins/video.tsx`
- **核心功能**:
  - 创建 `MarkdownVideo` 组件用于视频渲染
  - 支持多种视频格式: `.mp4`, `.webm`, `.mov`, `.avi`, `.m4v`, `.ogv`
  - 自动 MIME 类型检测和适配
  - HTML5 原生视频控制器（播放、暂停、进度、音量、全屏）
  - 播放前显示播放按钮 overlay 视觉提示
  - 响应式设计，自适应容器宽度

**2. Markdown 渲染器增强**
- **文件**: `app/src/components/Markdown.tsx`
- **修改内容**:
  - 导入视频播放器组件
  - 在 `img` 标签处理器中添加视频 URL 检测逻辑
  - 自动识别视频链接并渲染为 `MarkdownVideo` 组件
  ```tsx
  img: (props: { src?: string; alt?: string }) => {
    const videoExtensions = ['.mp4', '.webm', '.mov', '.avi', '.m4v', '.ogv'];
    const isVideo = props.src && videoExtensions.some(ext =>
      props.src!.toLowerCase().includes(ext) || props.src!.toLowerCase().endsWith(ext)
    );

    if (isVideo) {
      return <MarkdownVideo src={props.src} alt={props.alt} />;
    }
    // ... 原图片处理逻辑
  }
  ```

**3. 链接组件优化**
- **文件**: `app/src/components/markdown/Link.tsx`
- **增强功能**:
  - 检测 `<a>` 标签中的视频 URL
  - 在链接文本下方自动嵌入视频播放器
  - 保留原始链接文本（如"点击查看视频"）
  ```tsx
  if (isVideo) {
    return (
      <div className="video-link-wrapper">
        <a href={url} target="_blank" rel="noopener noreferrer">
          {children}
        </a>
        <MarkdownVideo src={url} alt={children?.toString()} />
      </div>
    );
  }
  ```

**4. 样式系统集成**
- **文件**: `app/src/assets/markdown/theme.less`
- **新增样式**:
  ```less
  .markdown-body .markdown-video-wrapper {
    position: relative;
    display: block;
    margin: 1rem 0;
    border-radius: 0.5rem;
    overflow: hidden;
    box-shadow: 0 1px 3px 0 rgba(0, 0, 0, 0.1);

    video {
      display: block;
      width: 100%;
      background-color: #000;
      border-radius: 0.5rem;
    }
  }
  ```
- **暗色主题适配**: 完整支持明暗主题切换
- **Hover 效果**: 阴影增强交互反馈

### 使用方式

用户可通过以下 Markdown 语法嵌入视频：

```markdown
# 方式 1: 图片语法（自动识别为视频）
![视频描述](https://example.com/video.mp4)

# 方式 2: 链接语法（链接下方显示播放器）
[点击查看视频](https://harbor.example.com/video.mp4)
```

### 技术特性

- **🎬 自动识别**: 通过文件扩展名智能检测视频 URL
- **📱 响应式设计**: 视频播放器自适应容器宽度，支持移动端
- **🎨 UI 一致性**: 与图片渲染样式保持一致的视觉风格
- **🌙 主题支持**: 完整适配明暗主题切换
- **⚡ 性能优化**:
  - 使用 `preload="metadata"` 仅预加载元数据
  - 懒加载策略，避免影响页面性能
- **🔊 原生控制**:
  - 浏览器原生 `<video>` 控件
  - 支持所有标准功能（播放/暂停/进度/音量/画中画/全屏）
  - 跨平台兼容性良好

### 关键问题修复

**问题**: 播放按钮 overlay 阻挡视频控件点击

**原因**: 装饰性的播放按钮 overlay 层覆盖在视频元素上方，拦截了所有鼠标事件。

**解决方案**: 添加 CSS `pointer-events-none` 属性
```tsx
<div className="... pointer-events-none">
  <div className="w-16 h-16 rounded-full ...">
    <Play className="w-8 h-8 text-black ml-1" />
  </div>
</div>
```

**效果**: 鼠标事件穿透 overlay 层，直接触达底层视频控件，保持视觉提示的同时确保交互正常。

### 构建与部署

**前端构建**:
```bash
cd app && npm run build
# ✓ 6538 modules transformed
# ✓ built in 21.19s
```

**验证状态**:
- ✅ TypeScript 编译通过，无类型错误
- ✅ 视频播放器组件正常渲染
- ✅ 自动视频格式检测功能正常
- ✅ 明暗主题样式适配完成
- ✅ 交互功能测试通过

### 架构亮点

**1. 组件化设计**
- 独立的 `MarkdownVideo` 组件，职责单一
- 可复用性强，易于维护和扩展
- 符合 React 组件最佳实践

**2. 渐进增强**
- 向后兼容，不影响现有功能
- 自动检测，无需手动配置
- 优雅降级，不支持的格式回退到链接

**3. 用户体验优化**
- 零配置，开箱即用
- 视觉反馈丰富（播放按钮提示、hover 效果）
- 移动端友好，触控操作流畅

### 开发工具链

- **开发环境**: React 18 + TypeScript + Vite HMR
- **构建工具**: Vite 4.5.3 生产构建优化
- **样式处理**: Tailwind CSS + Less 预处理器
- **图标库**: Lucide React (Play 图标)

## Linux.do OAuth2 单点登录集成 (2025-10-20)

### 功能概述

完整集成 Linux.do OAuth2 单点登录功能,用户可使用 Linux.do 账号一键登录 CoAI 系统。

### 技术实现

**后端开发** (Go):
- **核心文件**: `auth/oauth.go` (新建)
  - OAuth2 授权流程实现
  - CSRF 防护 (state 参数验证)
  - 用户信息获取和解析
  - 自动注册/登录逻辑
- **路由注册**: `auth/router.go`
  - `GET /oauth/linux-do/login` - 生成授权 URL
  - `GET /callback` - OAuth 回调处理
- **配置管理**: `config/config.yaml`
  - Linux.do OAuth 应用配置
  - Client ID/Secret 管理
  - 回调地址配置

**前端开发** (React + TypeScript):
- **登录页面**: `app/src/routes/Auth.tsx`
  - 添加 "使用 Linux.do 登录" 按钮
  - OAuth 授权流程触发
  - 错误处理和用户提示
- **回调页面**: `app/src/routes/OAuthCallback.tsx` (新建)
  - Token 接收和验证
  - 自动跳转到首页
- **路由配置**: `app/src/router.tsx`
  - `/oauth-success` 路由注册
- **样式设计**: `app/src/assets/pages/auth.less`
  - OAuth 按钮样式
  - 分隔线设计

### 关键技术特性

**1. 安全机制**
- **CSRF 防护**: state 参数随机生成,存储在 Redis 10分钟,单次使用后删除
- **Token 验证**: 完整的 OAuth2 授权码流程
- **用户追踪**: token 字段存储 `linux_do:{user_id}` 标识

**2. 用户管理**
- **自动注册**: 新用户自动创建账号,无需手动注册
- **邮箱绑定**: 通过邮箱识别已有账号,自动登录
- **用户名生成**: `linux_do_{username}` 格式,冲突时自动添加随机后缀

**3. 数据结构**
```go
type LinuxDoUserInfo struct {
    ID         int64  `json:"id"`        // Linux.do 用户 ID (数字类型)
    Username   string `json:"username"`  // 用户名
    Name       string `json:"name"`      // 显示名称
    Email      string `json:"email"`     // 电子邮箱
    AvatarURL  string `json:"avatar_url"` // 头像 URL
}
```

**4. OAuth2 配置**
```yaml
system:
  oauth:
    linux_do:
      enabled: true
      client_id: "iqrFITZB00hqvayAvSsDIhP9cykpjpqG"
      client_secret: "6EgtXxqMWg3PZehaEIIT1JmXLSCfL0ax"
      redirect_url: "https://coai.liujiarong.me/api/callback"
```

### 开发过程中解决的问题

**1. 双重 API 前缀问题**
- **问题**: 请求变成 `/api/api/oauth/linux-do/login`
- **原因**: axios baseURL 设置为 `/api`,前端请求又加了 `/api`
- **解决**: 前端请求路径改为 `/oauth/linux-do/login`,由 axios 自动加 baseURL

**2. 回调 404 错误**
- **问题**: 访问 `/callback` 返回 404
- **原因**: 前后端分离部署,前端 Router 优先匹配
- **解决**: 修改回调地址为 `/api/callback`,确保请求到达后端

**3. redirect_uri 不匹配**
- **问题**: `invalid_request` 错误
- **原因**: 配置文件中的回调地址与 Linux.do 注册的不一致
- **解决**: 统一回调地址为 `https://coai.liujiarong.me/api/callback`

**4. 用户信息解析失败**
- **问题**: `json: cannot unmarshal number into Go struct field`
- **原因**: Linux.do 返回的 `id` 是数字类型,代码定义为字符串
- **解决**: 将 `LinuxDoUserInfo.ID` 类型从 `string` 改为 `int64`

### 登录流程

```
1. 用户点击 "使用 Linux.do 登录"
   ↓
2. 前端调用 /api/oauth/linux-do/login
   ↓
3. 后端生成授权 URL (包含 state 参数)
   ↓
4. 跳转到 Linux.do 授权页面
   ↓
5. 用户同意授权
   ↓
6. Linux.do 回调到 /api/callback?code=xxx&state=xxx
   ↓
7. 后端验证 state 参数
   ↓
8. 用 code 换取 access_token
   ↓
9. 调用 /api/user 获取用户信息
   ↓
10. 检查邮箱是否存在
    ├─ 存在 → 直接登录
    └─ 不存在 → 自动注册新用户
   ↓
11. 生成 JWT token
   ↓
12. 重定向到 /oauth-success?token=xxx
   ↓
13. 前端验证 token 并跳转到首页
```

### 配置要求

**Linux.do SSO 应用配置**:
- **应用名**: coai
- **应用主页**: https://coai.liujiarong.me
- **回调地址**: https://coai.liujiarong.me/api/callback
- **权限范围**: openid, profile, email

**后端配置** (`config/config.yaml`):
- 设置 `system.oauth.linux_do.enabled: true`
- 配置 Client ID 和 Client Secret
- 设置正确的 redirect_url (必须包含 `/api`)

**依赖库**:
- Go: `golang.org/x/oauth2` (v0.32.0)
- React: axios (已有)

### 文件清单

| 类型 | 文件路径 | 说明 |
|------|---------|------|
| **后端** | `auth/oauth.go` | OAuth 核心逻辑 (新建) |
| | `auth/router.go` | 路由注册 (修改) |
| | `config/config.yaml` | OAuth 配置 (修改) |
| | `config.example.yaml` | 配置模板 (修改) |
| **前端** | `app/src/routes/Auth.tsx` | 登录页面 (修改) |
| | `app/src/routes/OAuthCallback.tsx` | 回调页面 (新建) |
| | `app/src/router.tsx` | 路由配置 (修改) |
| | `app/src/assets/pages/auth.less` | 样式文件 (修改) |
| **文档** | `LINUX_DO_OAUTH_SETUP.md` | 配置指南 (新建) |

### 验证状态

- ✅ 后端 OAuth2 授权流程
- ✅ 前端 UI 集成
- ✅ CSRF 防护机制
- ✅ 用户自动注册/登录
- ✅ 错误处理和日志
- ✅ 生产环境部署测试
- ✅ 数据类型兼容性修复

### 技术亮点

1. **完整的 OAuth2 实现**: 严格遵循 RFC 6749 标准
2. **安全性优先**: state 参数、单次使用、过期控制
3. **用户体验优化**: 一键登录、自动注册、错误提示
4. **生产环境适配**: 前后端分离部署支持
5. **详细的调试日志**: 方便问题排查和维护
6. **向后兼容**: 不影响现有用户名密码登录方式

## 待开发功能规划

### 短期目标 (1-2 周)
- [x] Markdown 视频播放功能 (已完成 2025-10-17)
- [x] Linux.do OAuth2 单点登录集成 (已完成 2025-10-20)
- [ ] Inpainting 功能用户测试和优化
- [ ] 更多 OAuth 提供商支持 (GitHub、Google)
- [ ] 图像生成参数微调界面
- [ ] 错误处理和用户反馈改进

### 中期目标 (1-2 月)
- [ ] 图像编辑历史记录
- [ ] 批量图像处理
- [ ] 高级画笔工具 (渐变、纹理)
- [ ] 模型性能监控面板
- [ ] 视频播放器增强（字幕支持、播放速度控制）

### 长期目标 (3-6 月)
- [ ] 视频生成集成
- [ ] 3D 模型生成支持
- [ ] 插件化架构
- [ ] 移动端优化

## 技术债务与改进点

### 代码质量
- [ ] 添加单元测试覆盖
- [ ] TypeScript 严格模式启用
- [ ] 错误边界组件完善
- [ ] 性能监控集成

### 架构优化
- [ ] 微服务拆分评估
- [ ] 数据库连接池优化
- [ ] 缓存策略完善
- [ ] API 限流机制

### 用户体验
- [ ] 加载状态优化
- [ ] 离线模式支持
- [ ] 无障碍访问改进
- [ ] 移动端响应式完善

## 开发团队与贡献

### 核心开发
- **主要开发者**: Claude Code AI Assistant
- **技术顾问**: 项目维护者
- **测试反馈**: 用户社区

### 开发方法
- **敏捷开发**: 快速迭代和用户反馈驱动
- **AI 辅助开发**: 使用 Claude Code 进行代码生成和优化
- **问题驱动**: 基于实际用户问题进行功能开发

## 许可证

本项目采用 Apache-2.0 开源许可证，详情请参阅 [LICENSE](LICENSE) 文件。

---

**最后更新**: 2025年10月20日
**文档版本**: v1.4 (新增 Linux.do OAuth2 单点登录)
**下次更新**: 根据开发进展实时更新