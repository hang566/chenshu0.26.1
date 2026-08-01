# 云舒 (YunShu) 全栈服务系统

基于 Go + Gin + GORM + SQLite 构建的后端服务，配合前端实现用户脱机/联机双模式及数据继承功能。

## 系统架构

### 后端服务 (bg/)
- **Go 1.21+** 实现的 RESTful API 服务
- 使用 **Gin** 框架处理 HTTP 请求
- 使用 **GORM** ORM 操作 **SQLite** 数据库
- JWT Token 认证机制
- 支持 CORS 跨域访问

### 前端应用 (service/)
- 纯 HTML/CSS/JavaScript 实现
- 自动检测服务器连接状态
- 智能切换联机/脱机模式
- 脱机账号升级为联机账号时数据自动继承

## 核心功能

### 1. 双模式账号系统

#### 脱机模式 (Offline Mode)
- 无需后端服务器
- 数据存储在浏览器 localStorage
- 支持账号注册、登录（本地验证）
- 方便离线环境使用

#### 联机模式 (Online Mode)
- 连接后端服务器
- 数据存储在 SQLite 数据库
- JWT Token 安全认证
- 支持跨设备数据同步

### 2. 账号升级与数据继承

#### 脱机升级为联机
1. 在脱机模式下使用系统，所有数据保存本地
2. 连接服务器后，在注册页注册新账号
3. 系统自动检测脱机账号并提示数据继承
4. 注册成功后，原脱机数据合并至新联机账号

#### 登录时数据继承
1. 已有联机账号，在其他设备或清除缓存后登录
2. 如果本地有脱机账号，系统提示是否迁移数据
3. 用户确认后，脱机数据合并到联机账号

## 后端 API 接口

### 用户模块

#### 注册 (无需认证)
- **POST** `/api/user/register`
```json
// 请求
{
    "username": "string",      // 必填，2-50字符
    "emailOrPhone": "string",  // 必填，邮箱或电话
    "password": "string"       // 必填，最少6位
}

// 响应
{
    "code": 200,
    "message": "注册成功",
    "data": {
        "token": "jwt-token",
        "user": {
            "id": 1,
            "username": "testuser",
            "emailOrPhone": "test@test.com",
            "code": "1785570141056401800",
            "type": "online",
            "avatar": "",
            "status": 1,
            "createdAt": "2026-08-01T15:42:21Z"
        }
    }
}
```

#### 登录 (无需认证)
- **POST** `/api/user/login`
```json
// 请求
{
    "loginId": "string",  // 用户名或专属编码
    "password": "string"
}

// 响应
{
    "code": 200,
    "message": "登录成功",
    "data": {
        "token": "jwt-token",
        "user": { ... }
    }
}
```

#### 获取用户信息 (需认证)
- **GET** `/api/user/info`
- **请求头**: `Authorization: <token>`

#### 更新用户信息 (需认证)
- **PUT** `/api/user/info`
```json
{
    "username": "string",
    "emailOrPhone": "string",
    "avatar": "string"
}
```

#### 修改密码 (需认证)
- **PUT** `/api/user/password`
```json
{
    "oldPassword": "string",
    "newPassword": "string"
}
```

#### 登出 (需认证)
- **POST** `/api/user/logout`

## 快速开始

### 1. 启动后端服务

```bash
# 进入后端目录
cd bg

# 编译 (Windows)
go build -o yunshu-bg.exe .

# 编译 (Linux/macOS)
CGO_ENABLED=0 go build -o yunshu-bg .

# 运行服务
./yunshu-bg
```

服务默认启动在 `http://localhost:8080`

### 2. 配置环境变量 (可选)

```bash
# Windows PowerShell
$env:SERVER_PORT = "8080"
$env:DB_PATH = "./data/yunshu.db"
$env:JWT_SECRET = "your-secret-key"
$env:JWT_EXPIRY = "24"

# 启动服务
.\yunshu-bg.exe
```

### 3. 访问前端

直接用浏览器打开 `service/` 下的 HTML 文件即可，系统会自动检测后端服务状态。

## 项目结构

```
web w.0.26.1/
├── bg/                           # 后端服务
│   ├── main.go                   # 入口文件
│   ├── config/                   # 配置管理
│   │   └── config.go
│   ├── database/                 # 数据库初始化
│   │   └── database.go
│   ├── handlers/                 # 请求处理
│   │   └── user_handler.go
│   ├── middleware/               # 中间件
│   │   └── auth.go
│   ├── models/                   # 数据模型
│   │   └── user.go
│   ├── routes/                   # 路由配置
│   │   └── routes.go
│   ├── go.mod                    # Go 模块
│   └── README.md                 # 后端文档
├── service/                      # 前端服务页面
│   ├── user/                     # 用户中心
│   ├── register/                 # 注册页
│   ├── login/                    # 登录页
│   └── js/                       # 公共 JS
│       └── api.js                # API 服务封装
└── UI/                           # UI 组件库
```

## 数据模型

### User 用户表

| 字段 | 类型 | 说明 |
|------|------|------|
| id | uint | 主键自增 |
| username | string | 用户名（唯一索引） |
| email_or_phone | string | 邮箱或电话 |
| password | string | bcrypt 加密密码 |
| code | string | 专属编码（唯一索引） |
| type | string | online/offline |
| avatar | string | 头像 URL |
| status | int | 1 正常 / 0 禁用 |
| created_at | datetime | 创建时间 |
| updated_at | datetime | 更新时间 |
| deleted_at | datetime | 软删除时间 |

## 技术细节

### 纯 Go SQLite 驱动

使用 `github.com/glebarez/sqlite` 作为 GORM 的 SQLite 驱动，无需 CGO 依赖，可在任何平台编译。

### JWT Token 认证

- Token 有效期默认 24 小时
- 支持 Bearer Token 认证
- Token 存储在浏览器 localStorage
- 认证中间件自动验证请求

### 跨域支持

后端默认启用 CORS 中间件，允许所有来源访问。生产环境建议限制具体域名。

## 开发说明

### 添加新的 API 接口

1. 在 `models/` 中定义数据模型
2. 在 `handlers/` 中实现请求处理逻辑
3. 在 `routes/` 中注册路由
4. 在 `main.go` 中配置路由分组

### 扩展功能

- **邮件/短信服务**: 可扩展通知模块
- **支付集成**: 可接入第三方支付
- **数据导出**: 可添加 Excel/PDF 导出
- **权限管理**: 可增加角色权限系统

## 注意事项

1. 生产环境请更换 JWT_SECRET 为强随机密钥
2. 建议将数据目录 (./data) 添加到 .gitignore
3. SQLite 适合单机使用，多用户场景请迁移 MySQL/PostgreSQL
4. 生产环境请配置正确的 CORS 白名单域名
