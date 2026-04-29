# 高校信息查询系统 (Go版本)

## 项目简介

这是一个用Go语言开发的高校信息查询系统，基于Gin框架和SQLite数据库。系统提供了完整的学校信息管理、用户认证、权限控制、搜索功能和个性化设置等功能，适用于高等院校信息展示和管理场景。

### 主要功能

- **学校信息管理**：支持学校的增删改查操作，包含详细的学校信息（地址、类别、性质、排名等）
- **用户认证系统**：基于Session的用户登录认证机制，支持密码强度验证
- **角色权限控制**：区分管理员(admin)和普通用户(user)两种角色，管理员拥有全部权限
- **智能搜索功能**：支持按学校ID或名称搜索，提供分页和缓存优化
- **滚动数据展示**：首页自动滚动展示学校信息
- **Logo上传管理**：支持学校Logo的上传和管理
- **操作日志记录**：记录用户的增删改操作，便于审计追踪
- **个性化设置**：支持用户主题、语言、收藏学校等个性化配置

## 技术栈

| 技术 | 版本 | 说明 |
|------|------|------|
| Go | 1.24.0 | 开发语言 |
| Gin | v1.9.1 | Web框架 |
| SQLite | modernc.org/sqlite v1.25.0 | 数据库 |
| bcrypt | - | 密码哈希 |

## 项目结构

```
school-go/
├── cmd/
│   └── server/
│       └── main.go              # 应用入口，负责初始化配置、数据库、路由
├── internal/
│   ├── api/                     # API层
│   │   ├── handlers/            # HTTP请求处理器
│   │   │   ├── university.go    # 学校相关请求处理
│   │   │   ├── user.go          # 用户相关请求处理
│   │   │   ├── settings.go      # 设置相关请求处理
│   │   │   ├── base.go          # 基础处理器
│   │   │   └── operation_log_handler.go  # 操作日志处理
│   │   ├── middleware/          # 中间件
│   │   │   └── auth.go          # 认证和角色验证中间件
│   │   └── routes.go            # 路由注册
│   ├── models/                  # 数据模型
│   │   ├── university.go        # 学校模型
│   │   ├── user.go              # 用户模型
│   │   ├── operation_log.go      # 操作日志模型
│   │   ├── session.go           # 会话模型
│   │   └── settings.go          # 用户设置模型
│   ├── repository/              # 数据访问层
│   │   ├── db.go                # 数据库初始化和连接池管理
│   │   ├── university_repo.go    # 学校数据访问
│   │   ├── user_repo.go         # 用户数据访问
│   │   └── operation_log_repository.go  # 操作日志数据访问
│   ├── service/                 # 业务逻辑层
│   │   ├── university_service.go # 学校业务逻辑
│   │   ├── user_service.go      # 用户业务逻辑
│   │   ├── session_service.go    # 会话业务逻辑
│   │   ├── operation_log_service.go  # 操作日志业务逻辑
│   │   └── settings_service.go   # 设置业务逻辑
│   ├── config/                  # 配置管理
│   │   └── config.go            # 配置加载
│   └── pkg/
│       └── errors/             # 错误处理
│           ├── errors.go        # 自定义错误类型
│           ├── types.go         # 错误类型定义
│           └── middleware.go     # 错误处理中间件
├── static/                      # 静态资源
│   ├── js/
│   │   └── script.js           # 前端JavaScript
│   └── logo/                   # 学校Logo存储目录
├── templates/                   # HTML模板
│   ├── index.html              # 首页
│   ├── university.html         # 学校详情页
│   ├── university_edit.html    # 学校编辑页
│   ├── add_school.html         # 添加学校页
│   ├── account.html            # 账户管理页
│   ├── operation_logs.html     # 操作日志页
│   ├── error.html              # 错误页
│   └── not_found.html          # 404页面
├── scripts/                    # 工具脚本
│   ├── check_db.go            # 数据库检查脚本
│   ├── import_year.go         # 数据导入脚本
│   ├── gen_password.go        # 密码生成工具
│   └── test_auth.go           # 认证测试工具
├── go.mod                     # Go模块文件
├── go.sum                     # 依赖锁定文件
├── .env.example               # 环境变量示例
├── run.sh                     # Linux/Mac运行脚本
├── run.bat                    # Windows运行脚本
├── build.sh                   # 构建脚本
└── README.md                  # 项目说明文档
```

## 数据库设计

### 用户表 (users)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER | 主键，自增 |
| username | TEXT | 用户名，唯一 |
| password | TEXT | 密码（bcrypt哈希） |
| role | TEXT | 角色：admin/user |

### 学校表 (universities)

| 字段 | 类型 | 说明 |
|------|------|------|
| 学校ID | TEXT | 学校唯一标识 |
| 学校名称 | TEXT | 学校名称 |
| 地址 | TEXT | 学校地址 |
| 类别 | TEXT | 类别（本科/专科等） |
| 性质 | TEXT | 性质（公立/私立） |
| 归属部门 | TEXT | 主管部门 |
| 标签 | TEXT | 学校标签 |
| 建校时间 | TEXT | 建校年份 |
| 占地面积 | TEXT | 校园面积 |
| 保研星级 | TEXT | 保研星级 |
| 博士点数量 | TEXT | 博士点数量 |
| 硕士点数量 | TEXT | 硕士点数量 |
| 国家重点学科数量 | TEXT | 国家重点学科数 |
| 软科综合排名 | TEXT | 软科排名 |
| 校友会综合排名 | TEXT | 校友会排名 |
| QS世界排名 | TEXT | QS世界排名 |
| US世界排名 | TEXT | US News排名 |
| 泰晤士排名 | TEXT | 泰晤士排名 |
| 人气值排名 | TEXT | 人气排名 |
| 基本信息 | TEXT | 基本信息介绍 |
| 办学形式 | TEXT | 办学形式 |
| logo_path | TEXT | Logo路径 |

### 操作日志表 (operation_logs)

| 字段 | 类型 | 说明 |
|------|------|------|
| id | INTEGER | 主键，自增 |
| operation_type | TEXT | 操作类型 |
| table_name | TEXT | 操作表名 |
| record_id | TEXT | 记录ID |
| user_id | INTEGER | 用户ID |
| username | TEXT | 用户名 |
| operation_time | DATETIME | 操作时间 |
| old_data | TEXT | 修改前数据 |
| new_data | TEXT | 修改后数据 |
| ip_address | TEXT | IP地址 |
| user_agent | TEXT | 用户代理 |

## 快速开始

### 环境要求

- Go 1.20 或更高版本
- SQLite3 数据库驱动（由 modernc.org/sqlite 提供，纯Go实现）

### 1. 安装依赖

```bash
go mod tidy
```

### 2. 配置环境变量

复制 `.env.example` 为 `.env` 并根据需要修改：

```bash
cp .env.example .env
```

配置项说明：

| 环境变量 | 默认值 | 说明 |
|---------|-------|------|
| PORT | 5000 | 服务器端口 |
| HOST | 0.0.0.0 | 服务器主机 |
| DB_PATH | school.db | 数据库文件路径 |
| FLASK_DEBUG | false | 调试模式 |

### 3. 运行应用

**Linux/Mac:**
```bash
./run.sh
```

**Windows:**
```bash
run.bat
```

**或者直接运行：**
```bash
go run cmd/server/main.go
```

应用将在 `http://localhost:5000` 启动。

### 4. 构建应用

```bash
./build.sh
```

构建完成后，可执行文件为 `school-app`（或 `school-app.exe`）。

## 默认账户

系统会自动创建一个默认管理员账户：

- 用户名：`admin`
- 密码：`admin123`

**请在首次登录后立即修改密码！**

## API文档

### 公开接口（无需登录）

| 方法 | 路径 | 说明 |
|------|------|------|
| GET | `/` | 首页 |
| GET | `/search` | 搜索学校（参数：keyword, page, pageSize） |
| GET | `/university/:id` | 获取学校详情 |
| GET | `/api/scrolling_data` | 获取滚动数据 |
| POST | `/api/scrolling_position` | 更新滚动位置 |

### 认证接口

| 方法 | 路径 | 说明 |
|------|------|------|
| POST | `/login` | 用户登录（参数：username, password） |
| POST | `/logout` | 用户注销 |
| GET | `/check_login` | 检查登录状态 |

### 需要登录的接口

| 方法 | 路径 | 权限 | 说明 |
|------|------|------|------|
| GET | `/account` | 全部用户 | 账户管理页面 |
| GET | `/get_accounts` | 全部用户 | 获取所有账户 |
| POST | `/add_account` | 仅管理员 | 添加账户 |
| POST | `/update_account` | 全部用户 | 更新账户 |
| POST | `/delete_account` | 仅管理员 | 删除账户 |
| GET | `/add_school` | 全部用户 | 添加学校页面 |
| POST | `/upload_logo` | 全部用户 | 上传学校Logo |
| POST | `/add_school` | 全部用户 | 添加学校 |
| GET | `/edit/:id` | 全部用户 | 编辑学校页面 |
| POST | `/save` | 全部用户 | 保存学校信息 |
| POST | `/delete_school` | 仅管理员 | 删除学校 |
| GET | `/operation_logs` | 全部用户 | 操作日志页面 |
| GET | `/get_operation_logs` | 全部用户 | 获取操作日志 |
| POST | `/delete_operation_logs` | 仅管理员 | 删除操作日志 |
| GET | `/api/settings` | 全部用户 | 获取用户设置 |
| POST | `/api/settings` | 全部用户 | 保存用户设置 |
| GET | `/api/favorites` | 全部用户 | 获取收藏学校 |
| POST | `/api/favorites` | 全部用户 | 添加收藏 |
| DELETE | `/api/favorites/:id` | 全部用户 | 移除收藏 |

### 登录请求示例

```json
POST /login
Content-Type: application/json

{
    "username": "admin",
    "password": "admin123"
}
```

### 响应示例

```json
{
    "success": true,
    "message": "登录成功"
}
```

## 权限说明

### 角色类型

- **admin（管理员）**：拥有全部权限，可以管理用户账户、删除学校和操作日志
- **user（普通用户）**：可以查看和编辑学校信息、管理自己的账户

### 权限矩阵

| 功能 | admin | user |
|------|-------|------|
| 查看学校信息 | ✅ | ✅ |
| 添加学校 | ✅ | ✅ |
| 编辑学校 | ✅ | ✅ |
| 删除学校 | ✅ | ❌ |
| 管理账户 | ✅ | 仅自己 |
| 查看操作日志 | ✅ | ✅ |
| 删除操作日志 | ✅ | ❌ |

## 安全特性

### 输入验证

- **XSS防护**：所有用户输入都会经过HTML特殊字符转义处理
- **SQL注入防护**：使用参数化查询，验证输入格式
- **密码强度验证**：密码至少8个字符，必须包含字母和数字

### 会话管理

- 使用安全的Cookie配置（HttpOnly、Secure标志）
- 会话有效期24小时
- 支持从Cookie或请求头获取会话ID

### 错误处理

- 统一的错误处理中间件
- 自定义错误类型（BadRequest、Unauthorized、Forbidden、NotFound、InternalServerError）
- 针对AJAX请求和普通请求返回不同格式的错误响应

## 性能优化

### 数据库优化

- 连接池配置：最大25个连接，空闲10个
- 自动创建常用字段索引（学校ID、学校名称、用户名）
- 搜索结果缓存机制

### 缓存策略

- 静态资源设置24小时缓存（Cache-Control、Expires）
- 搜索结果缓存，减少数据库查询

## 目录结构说明

### Logo存储

学校Logo存储在 `static/logo/` 目录下，文件名格式为 `{学校ID}.png`

### 模板文件

使用Gin框架的HTML模板引擎，模板文件位于 `templates/` 目录

## 工具脚本

### 数据库检查

```bash
go run scripts/check_db.go
```

检查数据库连接和表结构是否正常。

### 数据导入

```bash
go run scripts/import_year.go
```

批量导入学校建校时间数据。

### 密码生成

```bash
go run scripts/gen_password.go
```

生成bcrypt密码哈希。

### 认证测试

```bash
go run scripts/test_auth.go
```

测试用户认证功能。

## 开发指南

### 添加新功能流程

1. 在 `models/` 中定义数据模型
2. 在 `repository/` 中实现数据访问层
3. 在 `service/` 中实现业务逻辑
4. 在 `handlers/` 中实现HTTP处理函数
5. 在 `routes.go` 中注册路由
6. 如需权限控制，使用 `middleware.AuthRequired()` 和 `middleware.RoleRequired()`

### 错误处理规范

使用 `internal/pkg/errors/` 包中的错误类型：

```go
c.Error(errors.BadRequest("错误信息", err))
c.Error(errors.Unauthorized("未授权", nil))
c.Error(errors.Forbidden("权限不足", nil))
c.Error(errors.NotFound("资源不存在", nil))
c.Error(errors.InternalServerError("服务器错误", err))
```

## 注意事项

1. **数据库路径**：默认数据库路径在 `repository/db.go` 中硬编码，首次运行会自动创建
2. **静态文件服务**：Logo和静态文件通过Gin的静态文件服务提供
3. **模板加载**：模板文件路径相对于运行目录，需要确保 `templates/` 目录存在

## 写在最后
项目一开始是用- Flask+SQLite实现的，docker部署到服务器后发现资源消耗有点高，所以用Go重写了一个，两个项目都是AI对话式写的，将就用。


## 许可证

MIT License
