# Nebula Live

🚀 一个现代化的 Go 后端 API 服务，基于领域驱动设计(DDD)架构构建，支持多种数据库和完整的 RESTful API。

## ✨ 特性

- 🏗️ **领域驱动设计 (DDD)** - 清晰的架构分层和领域模型
- 🔥 **高性能框架** - 基于 Fiber v2.52.9 构建
- 🗄️ **多数据库支持** - PostgreSQL 和 SQLite
- 🔍 **ORM 集成** - EntGo v0.14.1 提供类型安全的数据访问
- 🔧 **依赖注入** - Uber Fx 实现模块化架构
- 📝 **结构化日志** - Zap 日志库，支持全局和依赖注入
- ⚡ **热重载** - Air 支持开发环境热重载
- 🐳 **容器化** - Docker 和 Docker Compose 支持
- 🔒 **统一错误处理** - APIError 标准化错误响应
- ✅ **健康检查** - 内置健康检查端点
- 🔐 **JWT 认证** - 基于 JWT 的用户认证和授权系统
- 🛡️ **中间件保护** - 路由级别的认证中间件
- 🔒 **RBAC 权限控制** - 基于角色的访问控制系统
- 👥 **角色管理** - 完整的角色和权限管理 API
- 🔑 **细粒度权限** - 资源级别的权限控制 (resource:action)

## 🏛️ 架构设计

```
nebula-live/
├── cmd/server/           # 应用程序入口
├── internal/
│   ├── app/             # 应用层 - Fiber 应用配置
│   ├── domain/          # 领域层
│   │   ├── entity/      # 领域实体
│   │   ├── repository/  # 仓储接口
│   │   └── service/     # 领域服务
│   └── infrastructure/  # 基础设施层
│       ├── config/      # 配置管理
│       ├── persistence/ # 数据持久化
│       └── web/         # Web层 (处理器、路由、中间件)
├── pkg/                 # 公共包
├── ent/                 # EntGo 生成的代码
└── configs/             # 配置文件
```

## 🛠️ 技术栈

| 组件 | 技术选型 | 版本 |
|------|----------|------|
| **Web框架** | Fiber | v2.52.9 |
| **ORM** | EntGo | v0.14.1 |
| **依赖注入** | Uber Fx | v1.24.0 |
| **日志** | Zap | v1.28.0 |
| **配置** | Viper | v1.20.0 |
| **CLI** | Cobra | v1.8.1 |
| **JWT认证** | golang-jwt/jwt | v5.3.0 |
| **密码哈希** | Argon2id | - |
| **数据库** | PostgreSQL / SQLite | - |
| **容器化** | Docker | - |

## 🚀 快速开始

### 前置要求

- Go 1.22+
- Make (推荐)
- Docker & Docker Compose (可选)

### 使用 Makefile (推荐)

项目提供了完整的 Makefile 来简化开发流程：

```bash
# 查看所有可用命令
make help

# 快速开始开发
make install-tools  # 安装开发工具 (Air, golangci-lint)
make db-sqlite      # 切换到 SQLite 配置
make dev           # 启动热重载开发服务器

# 检查服务健康状态
make health
```

### 手动安装 (不使用 Make)

1. **克隆项目**
```bash
git clone <repository-url>
cd nebula-live
```

2. **安装依赖**
```bash
go mod download
```

3. **配置数据库**
```bash
# 使用 SQLite (推荐开发环境)
cp configs/config-sqlite.yaml configs/config.yaml

# 或使用 PostgreSQL
# 确保 PostgreSQL 运行在 localhost:5432
```

4. **启动服务**
```bash
# 直接运行
go run ./cmd/server

# 或使用热重载 (需要先安装 Air)
go install github.com/cosmtrek/air@latest
air
```

5. **验证服务**
```bash
curl http://localhost:8080/health
```

### Docker 部署

#### 使用 Makefile (推荐)
```bash
# 开发环境 (热重载)
make docker-run-dev

# 生产环境
make compose-up

# 完整服务栈 (包含数据库)
make compose-up-full

# 查看日志
make compose-logs

# 停止服务
make compose-down
```

#### 手动 Docker 命令
```bash
# 开发环境 (热重载)
docker-compose -f docker-compose.dev.yml up app-dev

# 生产环境
docker-compose up app

# 完整服务栈 (包含数据库)
docker-compose --profile postgres --profile redis up
```

## 📚 API 文档

### 健康检查
```http
GET /health
```

### 用户认证

#### 用户注册
```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123",
  "nickname": "John"
}
```

**响应示例：**
```json
{
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "nickname": "John",
    "avatar": "",
    "status": "active",
    "created_at": "2025-08-09T10:30:00Z",
    "updated_at": "2025-08-09T10:30:00Z"
  },
  "message": "User registered successfully"
}
```

#### 用户登录
```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "username": "john_doe",
  "password": "password123"
}
```

**响应示例：**
```json
{
  "user": {
    "id": 1,
    "username": "john_doe",
    "email": "john@example.com",
    "nickname": "John",
    "avatar": "",
    "status": "active",
    "created_at": "2025-08-09T10:30:00Z",
    "updated_at": "2025-08-09T10:30:00Z"
  },
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_at": 1691578200,
  "message": "Login successful"
}
```

#### 获取当前用户信息
```http
GET /api/v1/auth/me
Authorization: Bearer {access_token}
```

#### 刷新访问令牌
```http
POST /api/v1/auth/refresh
Content-Type: application/json

{
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}
```

**响应示例：**
```json
{
  "access_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "token_type": "Bearer",
  "expires_at": 1691578200,
  "message": "Token refreshed successfully"
}
```

### 用户管理

⚠️ **所有用户管理API都需要JWT认证和管理员权限** - 请在请求头中包含 `Authorization: Bearer {access_token}`，且当前用户必须拥有`admin`角色

#### 创建用户
```http
POST /api/v1/users
Content-Type: application/json

{
  "username": "john_doe",
  "email": "john@example.com",
  "password": "password123",
  "nickname": "John"
}
```

#### 获取用户
```http
GET /api/v1/users/{id}
```

#### 更新用户
```http
PUT /api/v1/users/{id}
Content-Type: application/json

{
  "nickname": "John Smith",
  "avatar": "https://example.com/avatar.jpg"
}
```

#### 删除用户
```http
DELETE /api/v1/users/{id}
```

#### 用户列表
```http
GET /api/v1/users?page=1&limit=10
```

#### 用户状态管理
```http
POST /api/v1/users/{id}/activate    # 激活用户
POST /api/v1/users/{id}/deactivate  # 停用用户
POST /api/v1/users/{id}/ban         # 禁用用户
```

### RBAC 角色权限管理

⚠️ **所有RBAC管理API都需要JWT认证和管理员权限** - 请在请求头中包含 `Authorization: Bearer {access_token}`，且当前用户必须拥有`admin`角色

#### 角色管理

##### 创建角色
```http
POST /api/v1/roles
Content-Type: application/json

{
  "name": "editor",
  "display_name": "编辑员",
  "description": "拥有内容编辑权限的角色"
}
```

##### 获取角色
```http
GET /api/v1/roles/{id}
```

##### 更新角色
```http
PUT /api/v1/roles/{id}
Content-Type: application/json

{
  "display_name": "高级编辑员",
  "description": "拥有高级编辑权限的角色"
}
```

##### 删除角色
```http
DELETE /api/v1/roles/{id}
```

##### 角色列表
```http
GET /api/v1/roles?page=1&limit=10
```

##### 为用户分配角色
```http
POST /api/v1/roles/{id}/assign
Content-Type: application/json

{
  "user_id": 1
}
```

##### 移除用户角色
```http
DELETE /api/v1/roles/{id}/users/{userId}
```

##### 获取用户角色
```http
GET /api/v1/roles/users/{userId}
```

#### 权限管理

##### 创建权限
```http
POST /api/v1/permissions
Content-Type: application/json

{
  "name": "content:read",
  "display_name": "查看内容",
  "description": "查看内容的权限",
  "resource": "content",
  "action": "read"
}
```

##### 获取权限
```http
GET /api/v1/permissions/{id}
```

##### 更新权限
```http
PUT /api/v1/permissions/{id}
Content-Type: application/json

{
  "display_name": "查看所有内容",
  "description": "查看系统中所有内容的权限"
}
```

##### 删除权限
```http
DELETE /api/v1/permissions/{id}
```

##### 权限列表
```http
GET /api/v1/permissions?page=1&limit=10
```

##### 为角色分配权限
```http
POST /api/v1/permissions/{id}/assign
Content-Type: application/json

{
  "role_id": 1
}
```

##### 移除角色权限
```http
DELETE /api/v1/permissions/{id}/roles/{roleId}
```

##### 获取角色权限
```http
GET /api/v1/permissions/roles/{roleId}
```

##### 获取用户权限
```http
GET /api/v1/permissions/users/{userId}
```

#### 系统预定义角色和权限

**系统角色：**
- `admin` - 管理员，拥有所有系统权限
- `user` - 普通用户，拥有基本查看权限

**系统权限：**
- 用户管理：`user:read`, `user:write`, `user:delete`, `user:manage`
- 角色管理：`role:read`, `role:write`, `role:delete`, `role:manage`
- 权限管理：`permission:read`, `permission:write`, `permission:delete`, `permission:manage`
- 系统管理：`system:manage`

### 错误响应格式
```json
{
  "code": 400,
  "error": "Bad Request",
  "message": "Invalid request body"
}
```

## ⚙️ 配置说明

### 数据库配置

#### SQLite (开发推荐)
```yaml
database:
  driver: "sqlite"
  database: "data/nebula_live.db"  # 或 ":memory:" 内存数据库
```

#### PostgreSQL (生产推荐)
```yaml
database:
  driver: "postgres"
  host: "localhost"
  port: 5432
  username: "postgres"
  password: "password"
  database: "nebula_live"
  ssl_mode: "disable"
```

### 日志配置
```yaml
log:
  level: "info"
  format: "json"
  output: "logs/app.log"
  enable_console: true
  enable_file: true
```

### 服务配置
```yaml
server:
  host: "0.0.0.0"
  port: 8080
  read_timeout: 30s
  write_timeout: 30s
```

### JWT 配置
```yaml
jwt:
  secret: "your-secret-key-change-this-in-production"
  access_token_ttl: "15m"     # 访问令牌过期时间
  refresh_token_ttl: "168h"   # 刷新令牌过期时间 (7天)
  issuer: "nebula-live"       # JWT 发行者
```

### 认证中间件

项目包含完整的JWT认证中间件系统：
- **RequireAuth**: 强制认证中间件，所有用户管理API使用
- **OptionalAuth**: 可选认证中间件，可用于需要识别用户但不强制登录的场景

## 🔧 开发指南

### 项目结构说明

- **cmd/server**: 应用程序启动入口
- **internal/app**: Fiber 应用配置和生命周期管理
- **internal/domain**: 业务核心逻辑，包含实体、服务和仓储接口
- **internal/infrastructure**: 基础设施实现，包含数据库、配置、HTTP处理
- **pkg**: 可重用的工具包
- **ent**: EntGo ORM 自动生成的代码

### 添加新功能

1. **定义领域实体** (internal/domain/entity)
2. **创建仓储接口** (internal/domain/repository)  
3. **实现领域服务** (internal/domain/service)
4. **实现仓储** (internal/infrastructure/persistence)
5. **创建HTTP处理器** (internal/infrastructure/web/handler)
6. **注册路由** (internal/infrastructure/web/router)
7. **配置依赖注入模块**

### 数据库迁移
```bash
# EntGo 会自动处理模式迁移
# 应用启动时自动运行 client.Schema.Create()
```

### 日志使用

#### 全局日志
```go
import "nebula-live/pkg/logger"

logger.Info("操作成功", zap.String("key", "value"))
logger.Error("操作失败", zap.Error(err))
```

#### 依赖注入日志
```go
// 在构造函数中注入
func NewService(logger *zap.Logger) Service {
    return &service{logger: logger}
}
```

## 🧪 测试

### 使用 Makefile
```bash
# 运行所有测试
make test

# 运行测试并生成覆盖率报告
make test-coverage

# 运行基准测试
make bench

# 运行所有代码检查 (格式化、检查、测试)
make check
```

### 手动命令
```bash
# 运行所有测试
go test ./...

# 运行测试并生成覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 🚀 部署

### 使用 Makefile

#### 本地构建
```bash
# 构建应用
make build

# 运行应用
make run

# 构建多平台发布版本
make release
```

#### Docker 部署
```bash
# 构建生产镜像
make docker-build

# 运行 Docker 容器
make docker-run

# 停止 Docker 容器
make docker-stop
```

#### 数据库管理
```bash
# 切换到 SQLite
make db-sqlite

# 重置数据库
make db-reset
```

### 手动命令

#### 本地构建
```bash
go build -o server ./cmd/server
./server
```

#### Docker 部署
```bash
# 构建镜像
docker build -t nebula-live .

# 运行容器
docker run -p 8080:8080 nebula-live
```

## 🛠️ Makefile 命令参考

项目提供了完整的 Makefile，包含以下命令类别：

### 开发命令
```bash
make help          # 显示所有可用命令
make build         # 构建应用
make run           # 运行应用
make dev           # 启动热重载开发服务器
make clean         # 清理构建产物
```

### 代码质量
```bash
make test          # 运行测试
make test-coverage # 生成测试覆盖率报告
make bench         # 运行基准测试
make format        # 格式化代码
make vet           # 运行 go vet
make lint          # 运行 golangci-lint
make check         # 运行所有检查
```

### 依赖管理
```bash
make deps          # 下载依赖
make tidy          # 清理依赖
make install-tools # 安装开发工具
```

### Docker 操作
```bash
make docker-build     # 构建生产镜像
make docker-build-dev # 构建开发镜像
make docker-run       # 运行生产容器
make docker-run-dev   # 运行开发容器
make compose-up       # 启动服务栈
make compose-up-full  # 启动完整服务栈
make compose-down     # 停止服务
```

### 数据库和监控
```bash
make db-sqlite     # 切换到 SQLite
make db-reset      # 重置数据库
make health        # 检查应用健康状态
make logs          # 查看应用日志
make info          # 显示项目信息
```

## 📄 许可证

MIT License

## 🤝 贡献

欢迎贡献代码！请确保：

1. 遵循项目的代码风格
2. 添加适当的测试
3. 更新相关文档
4. 提交前运行所有测试
5. 遵循 Git 提交规范

### Git 提交规范

项目采用 [Conventional Commits](https://www.conventionalcommits.org/) 规范：

#### 提交格式
```
<type>[optional scope]: <description>

[optional body]

[optional footer(s)]
```

#### 提交类型 (Type)
- `feat`: 新功能
- `fix`: 修复 Bug
- `docs`: 文档更新
- `style`: 代码格式调整（不影响功能）
- `refactor`: 代码重构
- `perf`: 性能优化
- `test`: 测试相关
- `chore`: 构建工具、依赖管理等
- `ci`: CI/CD 配置
- `build`: 构建系统相关

#### 作用域 (Scope) - 可选
- `api`: API 相关
- `web`: Web 层相关
- `domain`: 领域层相关
- `infra`: 基础设施层相关
- `config`: 配置相关
- `db`: 数据库相关
- `docker`: Docker 相关
- `deps`: 依赖相关

#### 提交示例
```bash
# 新功能
git commit -m "feat(api): add user authentication endpoint"

# 修复 Bug
git commit -m "fix(db): resolve connection timeout issue"

# 文档更新
git commit -m "docs: update API documentation for user endpoints"

# 重构
git commit -m "refactor(domain): extract user validation logic to service"

# 性能优化
git commit -m "perf(db): optimize user query with database indexes"

# 配置变更
git commit -m "chore(docker): update Docker compose configuration"

# 破坏性变更
git commit -m "feat(api)!: change user API response format

BREAKING CHANGE: user API now returns different response structure"
```

#### 提交规则
- **描述**: 使用祈使语气，首字母小写，结尾不加句号
- **长度**: 描述部分不超过 50 个字符
- **语言**: 统一使用英文
- **破坏性变更**: 在类型后添加 `!` 或在正文中使用 `BREAKING CHANGE:`

## 📞 支持

如有问题或建议，请提交 Issue 或联系项目维护者。

---

⭐ 如果这个项目对你有帮助，请给我们一个 star！