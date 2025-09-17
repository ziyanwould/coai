# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Chat Nio is a Next Generation AIGC One-Stop Business Solution that combines a powerful API distribution system with a rich user interface design. It's built as a Go backend with a React frontend, providing both B2B and B2C AI chat services.

## Development Commands

### Frontend (React + Vite)
Located in the `app/` directory:

- `npm run dev` - Start development server (runs on host 0.0.0.0)
- `npm run build` - Full TypeScript compilation and build
- `npm run fast-build` - Quick build without TypeScript compilation
- `npm run lint` - Run ESLint with TypeScript support
- `npm run prettier` - Format code with Prettier
- `npm run preview` - Preview built application

### Backend (Go)
Main entry point is `main.go`:

- `go run main.go` - Start development server
- `go build -o chatnio` - Build binary
- Default port: 8094 (configurable via config.yaml)

### Docker Development
- `docker-compose up -d` - Run full stack (latest)
- `docker-compose -f docker-compose.stable.yaml up -d` - Stable version
- `docker-compose -f docker-compose.watch.yaml up -d` - With auto-updates
- `docker build -t chatniolocal .` - Build local Docker image
- `docker tag chatniolocal:latest harbor.ipv6.liujiarong.top:8024/library/chatniolocal:test` - Tag for Harbor
- `docker push harbor.ipv6.liujiarong.top:8024/library/chatniolocal:test` - Push to Harbor registry

## Architecture Overview

### Backend Structure (Go)
- **main.go**: Entry point that initializes the Gin engine and registers route groups
- **adapter/**: API adapters for different AI model providers
- **admin/**: Backend admin panel functionality and user management
- **auth/**: Authentication and authorization systems
- **channel/**: Channel management for multiple AI model providers
- **manager/**: Core business logic and conversation management
- **globals/**: Global configuration and shared state
- **utils/**: Utility functions for encryption, networking, file handling, etc.
- **middleware/**: HTTP middleware for CORS, rate limiting, etc.

### Frontend Structure (React)
- **src/main.tsx**: React application entry point
- **src/App.tsx**: Main app component with Redux Provider
- **src/router.tsx**: React Router configuration
- **src/store/**: Redux store and state management
- **src/components/**: Reusable UI components using Radix UI and Tailwind
- **src/assets/**: Stylesheets (Less files) and static assets

### Key Technical Details

1. **Multi-Provider AI Support**: The system supports OpenAI, Anthropic Claude, Google Gemini, Midjourney, and 10+ other AI providers through the `adapter/` and `channel/` modules.

2. **Channel Management**: Advanced load balancing with priority and weight settings, automatic failover, and enterprise-level channel status management.

3. **Billing System**: Supports both subscription and elastic billing models with token counting and usage tracking.

4. **Configuration**: Uses Viper for configuration management with YAML files. Main config example in `config.example.yaml`.

5. **Database**: MySQL primary, Redis for caching and session management.

6. **API Compatibility**: Provides OpenAI-compatible API endpoints for seamless integration.

## Configuration

- **config.yaml**: Main configuration file (copy from `config.example.yaml`)
- Environment variables override config file settings (e.g., `MYSQL_HOST` overrides `mysql.host`)
- Default admin credentials: username `root`, password `chatnio123456`

## Database Setup

The application requires:
- MySQL database (configured in config.yaml)
- Redis instance for caching and real-time features
- Database migrations are handled automatically on startup

## Key Features to Understand

- **Model Market**: Customizable model marketplace in admin panel
- **Conversation Sync**: Zero-cost cross-device conversation synchronization
- **File Parsing**: Built-in support for PDF, DOCX, PPTX, Excel, and images
- **PWA Support**: Progressive Web App capabilities
- **Internationalization**: Multi-language support (CN, US, JP, RU)
- **Preset System**: Custom AI model presets with cloud sync
- **Image Generation**: Advanced support for Cloudflare Workers AI image models including inpainting functionality
- **Interactive Inpainting**: Built-in canvas drawing tool for mask creation in image editing workflows

## Recent Major Updates (2025-09)

### Cloudflare Workers AI Image Generation Enhancement
- **Multi-Format Response Handling**: Added support for both JSON and binary response formats from Cloudflare API
- **Model Support Expansion**: Fixed display issues for Flux, Leonardo Phoenix, and Leonardo Lucid-Origin models
- **Interactive Inpainting Workflow**: Complete implementation of canvas-based mask drawing for inpainting models

### Technical Implementations
- **Backend** (`adapter/cloudflare/`):
  - Dual response format detection (JSON vs binary)
  - Enhanced parameter handling for inpainting models
  - Improved base64 image processing with deduplication
  - Model-specific parameter optimization

- **Frontend** (`app/src/components/plugins/`):
  - `InpaintingCanvas.tsx`: Full-featured drawing canvas with brush/eraser tools
  - `InpaintingTrigger.tsx`: UI component for launching inpainting workflow
  - Enhanced Markdown rendering for custom inpainting triggers

### Bug Fixes
- Fixed duplicate image extraction in `utils/char.go`
- Corrected API parameter naming (`mask_image` vs `mask_b64`)
- Resolved case-sensitive model detection issues
- Enhanced error handling and debug logging

### User Tracking Enhancement (2025-09-17)
- **User Identification**: Added `user` and `user_ip` fields to Cloudflare image generation requests
- **Middleware Compatibility**: Ensures upstream middleware can properly track users and collect statistics
- **Consistent Implementation**: Matches OpenAI adapter pattern for unified user tracking across providers

## Testing and Quality

- ESLint configuration with TypeScript support
- Prettier for code formatting
- No explicit test commands found - verify testing approach in project
- Check Docker health and database connectivity before deployment
- **Development Testing**: Use `harbor.ipv6.liujiarong.top:8024/library/chatniolocal:test` for latest features