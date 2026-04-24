# MoChat PHP-to-Go 迁移规格说明

## Why

MoChat 当前后端基于 PHP + Swoole + Hyperf 框架，存在 PHP 版本兼容性限制（要求 PHP 7.4，不支持 PHP 8+）、协程支持非原生、类型系统较弱等问题。迁移至 Go 语言可获得原生并发优势、严格类型安全、单二进制部署便利性以及更优的长期维护性。

## What Changes

- **BREAKING**: 整个后端从 PHP/Hyperf 重写为 Go 语言实现
- 使用 Go 生态框架和库替代 PHP 依赖（Gin/Echo → Hyperf, GORM → Hyperf Database, asynq → Hyperf AsyncQueue 等）
- 保持与原项目完全一致的 API 接口路径、请求/响应格式，确保前端代码零修改
- 保持与原项目完全一致的数据库表结构（mc_ 前缀），确保数据平滑迁移
- 实现双轨 JWT 认证体系（Dashboard + Sidebar）
- 实现基于 RBAC 的权限控制（含三级数据权限）
- 实现企业微信全量事件回调处理
- 实现异步队列、定时任务、事件驱动等异步处理机制
- 实现多驱动文件存储（local/oss/cos/s3/minio/qiniu）
- 实现所有插件功能（渠道活码、欢迎语、自动拉群、标签建群、群欢迎语、消息群发、客户转接、统计、裂变）

## Impact

- Affected code: `/Users/zhanglei/GitHub/mochat/api-server` → `/Users/zhanglei/GitHub/mochat/api-server-go`
- Affected systems: MySQL (mc_ 前缀表)、Redis（缓存+队列）、文件存储、企业微信 API 集成
- 前端项目（dashboard/sidebar/operation）无需修改，API 接口保持兼容

## 技术选型

| 分类 | PHP 原技术 | Go 替代技术 | 说明 |
|------|-----------|------------|------|
| Web 框架 | Hyperf 2.2 | Gin v1.9+ | 高性能 HTTP 框架，社区活跃 |
| ORM | Hyperf Database (Eloquent) | GORM v2 | Go 生态最成熟的 ORM |
| 缓存 | Hyperf Redis (Swoole Coroutine) | go-redis/v9 | Redis 客户端，支持连接池 |
| 队列 | Hyperf AsyncQueue (Redis) | asynq | 基于 Redis 的异步任务队列 |
| 定时任务 | Hyperf Crontab | robfig/cron/v3 | Go 标准 cron 库 |
| JWT 认证 | 96qbhy/hyperf-auth | golang-jwt/jwt/v5 | JWT 库 |
| 配置管理 | Hyperf Config | viper | 配置管理，支持多格式 |
| 日志 | Hyperf Logger (Monolog) | zap | 高性能结构化日志 |
| 企业微信 SDK | EasyWeChat 5.x | silenceper/wechat/v2 | Go 企业微信 SDK |
| 文件存储 | Hyperf Filesystem (Flysystem) | 对各 SDK 直接封装 | OSS/COS/S3/七牛 |
| 验证 | Hyperf Validation | go-playground/validator/v10 | 请求参数验证 |
| 事件系统 | 自定义 EventDispatcher | 自行实现 EventBus | 观察者模式 |
| HTTP 客户端 | Hyperf Guzzle | resty | HTTP 请求库 |
| Excel | phpspreadsheet | excelize | Excel 读写 |
| FFmpeg | php-ffmpeg | go-ffmpeg / exec ffmpeg | 音视频处理 |

## 项目目录结构

```
api-server-go/
├── cmd/
│   └── server/
│       └── main.go                 # 入口
├── internal/
│   ├── config/                     # 配置加载
│   ├── middleware/                  # 中间件（认证、权限、CORS等）
│   ├── handler/                    # HTTP 处理器（对应 PHP Action）
│   │   ├── dashboard/              # 后台接口
│   │   └── sidebar/                # 侧边栏接口
│   ├── logic/                      # 业务逻辑层
│   ├── service/                    # 服务层（数据访问）
│   ├── model/                      # 数据模型（GORM Model）
│   ├── queue/                      # 异步队列处理
│   ├── task/                       # 定时任务
│   ├── event/                      # 事件定义与调度
│   ├── listener/                   # 事件监听器
│   └── pkg/                        # 内部公共包
│       ├── wechat/                 # 企业微信封装
│       ├── storage/                # 文件存储
│       ├── rbac/                   # RBAC 权限
│       └── response/               # 统一响应格式
├── pkg/                            # 可导出公共包
│   └── utils/                      # 工具函数
├── config/                         # 配置文件
│   ├── config.yaml                 # 主配置
│   └── config.example.yaml         # 示例配置
├── migrations/                     # 数据库迁移
├── api/                            # API 文档（OpenAPI/Swagger）
├── scripts/                        # 脚本
├── Dockerfile
├── docker-compose.yml
├── go.mod
├── go.sum
└── Makefile
```

## ADDED Requirements

### Requirement: Go 项目初始化

系统 SHALL 创建完整的 Go 项目结构，包含 go.mod、Makefile、Dockerfile、配置文件模板等基础设施。

#### Scenario: 项目初始化成功
- **WHEN** 执行 `go build ./cmd/server/`
- **THEN** 成功编译生成可执行文件

#### Scenario: 配置加载
- **WHEN** 启动服务且配置文件存在
- **THEN** 成功加载所有配置项（数据库、Redis、JWT、文件存储、企业微信等）

### Requirement: 数据库模型与连接

系统 SHALL 使用 GORM 连接 MySQL 数据库，定义与原项目完全一致的表结构模型（mc_ 前缀），支持连接池和模型缓存。

#### Scenario: 数据库连接
- **WHEN** 启动服务且数据库配置正确
- **THEN** 成功建立连接池，可执行 CRUD 操作

#### Scenario: 模型定义
- **WHEN** 定义所有 35+ 张表的 GORM Model
- **THEN** 字段名、类型、索引与原项目 SQL 完全一致

#### Scenario: 分表支持
- **WHEN** 访问 mc_work_unionid_external_userid_mapping 表
- **THEN** 根据规则路由到 _0 ~ _9 分表

### Requirement: Redis 缓存与队列

系统 SHALL 使用 go-redis 连接 Redis，实现缓存操作和异步队列。

#### Scenario: Redis 连接
- **WHEN** 启动服务且 Redis 配置正确
- **THEN** 成功建立连接池，可执行缓存操作

#### Scenario: 异步队列
- **WHEN** 投递异步任务到队列
- **THEN** worker 进程正确消费并处理任务

### Requirement: HTTP 路由与 API 接口

系统 SHALL 使用 Gin 框架注册所有 API 路由，路径和请求方法与原项目完全一致。

#### Scenario: 路由注册
- **WHEN** 启动服务
- **THEN** 所有 100+ 个 API 端点正确注册并可访问

#### Scenario: 请求响应格式
- **WHEN** 客户端发送 API 请求
- **THEN** 响应格式与原项目一致（`{code: 0, data: {}, msg: ""}`）

### Requirement: 双轨 JWT 认证

系统 SHALL 实现两套 JWT 认证体系：Dashboard 认证（User 模型）和 Sidebar 认证（WorkEmployee 模型）。

#### Scenario: Dashboard 认证
- **WHEN** 后台用户通过手机号+密码登录
- **THEN** 签发 jwt guard 的 JWT Token，后续请求通过 DashboardAuthMiddleware 验证

#### Scenario: Sidebar 认证
- **WHEN** 侧边栏用户通过企业微信 OAuth 登录
- **THEN** 签发 sidebar guard 的 JWT Token，后续请求通过 SidebarAuthMiddleware 验证

#### Scenario: 白名单路由
- **WHEN** 请求路径在白名单中
- **THEN** 跳过认证，直接处理请求

### Requirement: RBAC 权限控制

系统 SHALL 实现基于角色的访问控制，包括菜单权限和数据权限（全企业/本部门/本人）。

#### Scenario: 菜单权限
- **WHEN** 用户访问需权限的接口
- **THEN** PermissionMiddleware 检查用户角色是否拥有对应菜单权限

#### Scenario: 数据权限 - 全企业
- **WHEN** 角色数据权限为"全企业"
- **THEN** 可查看所有数据

#### Scenario: 数据权限 - 本部门
- **WHEN** 角色数据权限为"本部门"
- **THEN** 仅可查看本部门及下级部门数据

#### Scenario: 数据权限 - 本人
- **WHEN** 角色数据权限为"本人"
- **THEN** 仅可查看自己创建的数据

### Requirement: 企业微信集成

系统 SHALL 实现企业微信 API 集成，包括通讯录同步、客户管理、消息推送、事件回调等。

#### Scenario: 通讯录同步
- **WHEN** 触发通讯录同步
- **THEN** 调用企业微信 API 获取部门和员工数据并更新本地数据库

#### Scenario: 客户事件回调
- **WHEN** 接收企业微信客户事件回调（添加/删除/更新客户）
- **THEN** 解密回调数据，触发内部事件，更新本地数据库

#### Scenario: 消息推送
- **WHEN** 发送企业微信消息
- **THEN** 调用企业微信 API 成功推送消息

### Requirement: 异步队列处理

系统 SHALL 实现基于 Redis 的异步队列，支持原项目的所有队列通道和任务类型。

#### Scenario: 队列投递与消费
- **WHEN** 业务逻辑投递异步任务
- **THEN** 对应通道的 worker 正确消费任务并执行

#### Scenario: 队列失败重试
- **WHEN** 任务执行失败
- **THEN** 按配置进行重试，超过重试次数后记录失败日志

### Requirement: 定时任务

系统 SHALL 实现所有原项目的定时任务，包括企业数据统计、素材 mediaId 更新、员工统计同步、应用同步等。

#### Scenario: 定时任务执行
- **WHEN** 到达定时任务触发时间
- **THEN** 正确执行任务逻辑

#### Scenario: 动态 Crontab
- **WHEN** 企业配置了同步数据任务
- **THEN** 动态注册对应的定时任务

### Requirement: 事件驱动系统

系统 SHALL 实现事件总线，支持事件定义、监听器注册和事件分发。

#### Scenario: 事件触发与监听
- **WHEN** 业务逻辑触发事件
- **THEN** 所有已注册的监听器按顺序执行

### Requirement: 多驱动文件存储

系统 SHALL 支持多种文件存储驱动：local、阿里云 OSS、腾讯云 COS、AWS S3/MinIO、七牛。

#### Scenario: 文件上传
- **WHEN** 上传文件
- **THEN** 根据配置的驱动存储文件，返回文件 URL

#### Scenario: 文件 URL 获取
- **WHEN** 获取文件完整 URL
- **THEN** 根据驱动类型返回正确的访问地址

### Requirement: 完整业务模块实现

系统 SHALL 完整实现以下所有业务模块：

1. **企业模块 (corp)**: 企业 CRUD、企业微信绑定、回调处理
2. **租户模块 (tenant)**: 租户管理
3. **用户模块 (user)**: 用户 CRUD、登录认证、密码管理
4. **RBAC 模块 (rbac)**: 菜单管理、角色管理、权限分配
5. **客户模块 (work-contact)**: 客户管理、标签管理、高级属性、跟进记录、互动轨迹、流失提醒
6. **部门模块 (work-department)**: 部门管理、通讯录同步
7. **员工模块 (work-employee)**: 员工管理、通讯录同步、统计数据
8. **客户群模块 (work-room)**: 客户群管理、群分组、群统计
9. **应用模块 (work-agent)**: 企业应用管理、OAuth 授权、JSSDK
10. **素材库模块 (medium)**: 素材管理、素材分组、mediaId 同步
11. **欢迎语模块 (greeting)**: 欢迎语管理
12. **首页模块 (index)**: 首页数据、折线图
13. **侧边工具栏模块 (chat-tool)**: 工具栏配置
14. **公众号模块 (official-account)**: 公众号授权管理

### Requirement: 完整插件实现

系统 SHALL 完整实现以下所有插件功能：

1. **渠道活码 (channel-code)**: 活码 CRUD、分组管理、统计
2. **客户消息群发 (contact-message-batch-send)**: 批量消息发送、状态追踪
3. **客户转接 (contact-transfer)**: 在职转接、离职继承
4. **自动拉群 (room-auto-pull)**: 自动拉群配置
5. **标签建群 (room-tag-pull)**: 按标签筛选建群
6. **群欢迎语 (room-welcome)**: 入群欢迎语模板
7. **群消息群发 (room-message-batch-send)**: 群消息批量发送
8. **统计 (statistic)**: 员工统计、排行榜
9. **裂变 (work-fission)**: 裂变活动、海报、邀请、推送

### Requirement: 单元测试

系统 SHALL 为所有核心模块编写单元测试，覆盖率不低于 70%。

#### Scenario: 测试运行
- **WHEN** 执行 `go test ./...`
- **THEN** 所有测试通过，覆盖率 >= 70%

### Requirement: 部署支持

系统 SHALL 提供 Dockerfile、docker-compose.yml 和 Makefile，支持容器化部署。

#### Scenario: Docker 构建
- **WHEN** 执行 `docker build -t mochat-api .`
- **THEN** 成功构建 Docker 镜像

#### Scenario: Docker Compose 启动
- **WHEN** 执行 `docker-compose up -d`
- **THEN** 完整系统（API + MySQL + Redis）成功启动

### Requirement: API 文档

系统 SHALL 提供 OpenAPI/Swagger 格式的 API 文档。

#### Scenario: 文档访问
- **WHEN** 访问 /swagger/index.html
- **THEN** 显示完整的 API 文档

## MODIFIED Requirements

无（此为全新项目，非修改现有项目）

## REMOVED Requirements

无（此为全新项目，非删除现有功能）
