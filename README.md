# 全屋智能家居中央控制系统

基于 Vue3 + Tailwind CSS + Go (Gin) + PostgreSQL 的全屋智能家居中央控制系统。

## 项目结构

```
cx1/
├── backend/                    # 后端 Go 服务
│   ├── cmd/
│   │   └── server/            # 服务入口
│   │       └── main.go
│   ├── internal/
│   │   ├── config/            # 配置管理
│   │   ├── models/            # 数据模型
│   │   ├── database/          # 数据库连接与初始化
│   │   ├── repositories/      # 数据访问层
│   │   ├── services/          # 业务逻辑层
│   │   ├── handlers/          # API 控制器
│   │   ├── middleware/        # 中间件
│   │   ├── websocket/         # WebSocket 实时推送
│   │   └── gateway/           # 网关模拟器（数据上报）
│   ├── pkg/
│   │   └── utils/
│   ├── go.mod
│   ├── go.sum
│   └── .env
├── frontend/                   # 前端 Vue3 应用
│   ├── src/
│   │   ├── api/               # API 封装
│   │   ├── router/            # 路由配置
│   │   ├── stores/            # 状态管理 (Pinia)
│   │   ├── types/             # TypeScript 类型定义
│   │   ├── views/
│   │   │   ├── mobile/        # 移动端页面
│   │   │   └── dashboard/     # 监控大屏页面
│   │   ├── components/        # 公共组件
│   │   ├── utils/             # 工具函数
│   │   └── assets/            # 静态资源
│   ├── package.json
│   ├── vite.config.ts
│   ├── tailwind.config.js
│   └── tsconfig.json
├── docker-compose.yml
└── README.md
```

## 核心功能

### 1. 设备联动规则配置模块
- 手机网页端适配，支持按房间管理联动规则
- 支持多种触发条件（温度、湿度、PM2.5、甲醛、定时、手动）
- 支持多种执行动作（控制灯光、空调、净化器、发送通知）
- 规则的增删改查、启用/禁用切换

### 2. 环境数据实时监控大屏
- WebSocket 实时数据推送（每秒刷新）
- 网关模拟高频上报各房间传感器数据
- 温度、湿度、PM2.5、甲醛四项指标实时监控
- ECharts 趋势图表、仪表盘展示
- 各房间数据卡片、状态颜色指示
- 全屋平均数据统计

## 技术栈

### 后端
- **语言**: Go 1.21
- **框架**: Gin v1.9.1
- **数据库**: PostgreSQL 15
- **ORM**: GORM v1.25.5
- **WebSocket**: Gorilla WebSocket v1.5.1
- **配置**: godotenv v1.5.1

### 前端
- **框架**: Vue 3.4 (Composition API)
- **构建工具**: Vite 5.0
- **类型系统**: TypeScript 5.3
- **样式**: Tailwind CSS 3.3
- **路由**: Vue Router 4.2
- **状态管理**: Pinia 2.1
- **图表**: ECharts 5.4
- **HTTP 客户端**: Axios 1.6

## 快速开始

### 1. 启动数据库

```bash
docker-compose up -d
```

### 2. 启动后端服务

```bash
cd backend
go mod download
go run cmd/server/main.go
```

后端服务将在 `http://localhost:8080` 启动

### 3. 启动前端服务

```bash
cd frontend
npm install
npm run dev
```

前端服务将在 `http://localhost:5173` 启动

## 访问地址

- 移动端联动规则配置: `http://localhost:5173/mobile/rules`
- 环境数据监控大屏: `http://localhost:5173/dashboard`

## API 接口

### 房间管理
- `GET /api/v1/rooms` - 获取所有房间
- `GET /api/v1/rooms/:id` - 获取房间详情
- `POST /api/v1/rooms` - 创建房间
- `PUT /api/v1/rooms/:id` - 更新房间
- `DELETE /api/v1/rooms/:id` - 删除房间

### 设备管理
- `GET /api/v1/devices` - 获取设备列表（支持按房间过滤）
- `PATCH /api/v1/devices/:id/toggle` - 切换设备状态

### 联动规则管理
- `GET /api/v1/linkage-rules` - 获取规则列表（支持按房间过滤）
- `GET /api/v1/linkage-rules/:id` - 获取规则详情
- `POST /api/v1/linkage-rules` - 创建规则
- `PUT /api/v1/linkage-rules/:id` - 更新规则
- `PATCH /api/v1/linkage-rules/:id/toggle` - 切换规则启用状态
- `DELETE /api/v1/linkage-rules/:id` - 删除规则

### 环境数据
- `GET /api/v1/environment/latest` - 获取所有房间最新环境数据
- `GET /api/v1/environment/latest/:room_id` - 获取指定房间最新环境数据
- `GET /api/v1/environment/history/:room_id` - 获取指定房间历史数据

### WebSocket
- `ws://localhost:8080/ws` - 环境数据实时推送

## 系统特性

1. **前后端完全解耦**: 独立的 Go 后端 API 服务和 Vue3 前端应用
2. **多文件分层架构**: Repository → Service → Handler 三层架构
3. **实时数据推送**: WebSocket 实现每秒数据刷新
4. **网关模拟**: 内置网关模拟器，模拟高频传感器数据上报
5. **响应式设计**: 移动端页面适配手机屏幕，大屏页面适配大尺寸显示器
6. **数据持久化**: PostgreSQL 存储所有历史环境数据
