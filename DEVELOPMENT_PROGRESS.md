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

### 基础设施
- **容器化**: Docker + Docker Compose
- **镜像仓库**: Harbor
- **开发环境**:
  - nodemon (后端热重载)
  - Vite HMR (前端热重载)

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

## 待开发功能规划

### 短期目标 (1-2 周)
- [ ] Inpainting 功能用户测试和优化
- [ ] 更多 Cloudflare 模型支持测试
- [ ] 图像生成参数微调界面
- [ ] 错误处理和用户反馈改进

### 中期目标 (1-2 月)
- [ ] 图像编辑历史记录
- [ ] 批量图像处理
- [ ] 高级画笔工具 (渐变、纹理)
- [ ] 模型性能监控面板

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

---

**最后更新**: 2025年9月17日
**文档版本**: v1.1
**下次更新**: 根据开发进展实时更新