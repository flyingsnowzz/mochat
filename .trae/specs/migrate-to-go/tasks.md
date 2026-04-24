# Tasks

## 阶段一：项目基础设施

- [x] Task 1: 初始化 Go 项目结构
  - [x] 1.1 创建 api-server-go 目录及完整子目录结构
  - [x] 1.2 初始化 go.mod
  - [x] 1.3 创建 Makefile
  - [x] 1.4 创建 Dockerfile
  - [x] 1.5 创建 docker-compose.yml
  - [x] 1.6 创建 .gitignore

- [x] Task 2: 配置管理模块
  - [x] 2.1 定义配置结构体
  - [x] 2.2 使用 viper 加载配置
  - [x] 2.3 支持环境变量覆盖配置
  - [x] 2.4 创建 config.example.yaml

- [x] Task 3: 日志与错误处理
  - [x] 3.1 集成 zap 日志库
  - [x] 3.2 定义统一错误码体系
  - [x] 3.3 定义统一响应格式

- [x] Task 4: 数据库连接与 GORM 初始化
  - [x] 4.1 配置 MySQL 连接池
  - [x] 4.2 初始化 GORM，设置表前缀 mc_
  - [x] 4.3 启用模型缓存
  - [x] 4.4 实现分表路由逻辑

- [x] Task 5: Redis 连接与缓存
  - [x] 5.1 配置 go-redis 连接池
  - [x] 5.2 封装缓存操作

## 阶段二：核心数据模型

- [x] Task 6: 定义所有 GORM 数据模型
  - [x] 6.1 企业与租户模型
  - [x] 6.2 用户与权限模型
  - [x] 6.3 员工与部门模型
  - [x] 6.4 客户模型
  - [x] 6.5 客户标签模型
  - [x] 6.6 客户群模型
  - [x] 6.7 应用与工具模型
  - [x] 6.8 系统模型

- [x] Task 7: 定义插件数据模型
  - [x] 7.1 渠道活码模型
  - [x] 7.2 消息群发模型
  - [x] 7.3 客户转接模型
  - [x] 7.4 自动拉群模型
  - [x] 7.5 标签建群模型
  - [x] 7.6 群欢迎语模型
  - [x] 7.7 裂变模型

## 阶段三：认证与中间件

- [x] Task 8: JWT 双轨认证实现
  - [x] 8.1 Dashboard JWT 认证
  - [x] 8.2 Sidebar JWT 认证
  - [x] 8.3 DashboardAuthMiddleware
  - [x] 8.4 SidebarAuthMiddleware
  - [x] 8.5 认证白名单路由配置

- [x] Task 9: RBAC 权限中间件
  - [x] 9.1 PermissionMiddleware
  - [x] 9.2 三级数据权限
  - [x] 9.3 RBAC 工具函数

- [x] Task 10: 通用中间件
  - [x] 10.1 CoreMiddleware
  - [x] 10.2 CORS 中间件
  - [x] 10.3 请求日志中间件
  - [x] 10.4 限流中间件

## 阶段四：企业微信集成

- [x] Task 11: 企业微信 SDK 封装
  - [x] 11.1 封装企业微信 API 客户端
  - [x] 11.2 回调解密/加密
  - [x] 11.3 WeWorkFactory
  - [x] 11.4 素材上传

- [x] Task 12: 微信开放平台集成
  - [x] 12.1 公众号授权流程
  - [x] 12.2 授权事件回调处理
  - [x] 12.3 消息事件回调处理

## 阶段五：异步处理基础设施

- [x] Task 13: 异步队列系统
  - [x] 13.1 asynq 集成
  - [x] 13.2 定义队列通道
  - [x] 13.3 队列投递接口
  - [x] 13.4 worker 进程管理
  - [x] 13.5 失败重试和死信队列

- [x] Task 14: 定时任务系统
  - [x] 14.1 robfig/cron 调度器
  - [x] 14.2 CorpData 任务
  - [x] 14.3 MediaIdUpdate 任务
  - [x] 14.4 WorkEmployeeStatistic 任务
  - [x] 14.5 SyncWorkAgent 任务
  - [x] 14.6 动态 Crontab 注册机制

- [x] Task 15: 事件驱动系统
  - [x] 15.1 EventBus 实现
  - [x] 15.2 所有事件类型
  - [x] 15.3 所有事件监听器

## 阶段六：文件存储

- [x] Task 16: 多驱动文件存储
  - [x] 16.1 Storage 接口定义
  - [x] 16.2 local 驱动
  - [x] 16.3 阿里云 OSS 驱动
  - [x] 16.4 腾讯云 COS 驱动
  - [x] 16.5 S3/MinIO 驱动
  - [x] 16.6 七牛驱动
  - [x] 16.7 异步文件上传队列

## 阶段七：核心业务模块 Service + Logic + Handler

- [x] Task 17: 企业模块 (corp)
- [x] Task 18: 租户模块 (tenant)
- [x] Task 19: 用户模块 (user)
- [x] Task 20: RBAC 模块 (rbac)
- [x] Task 21: 客户模块 (work-contact)
- [x] Task 22: 部门模块 (work-department)
- [x] Task 23: 员工模块 (work-employee)
- [x] Task 24: 客户群模块 (work-room)
- [x] Task 25: 应用模块 (work-agent)
- [x] Task 26: 素材库模块 (medium)
- [x] Task 27: 首页模块 (index)
- [x] Task 28: 侧边工具栏模块 (chat-tool)
- [x] Task 29: 公众号模块 (official-account)
- [x] Task 30: 通用模块 (common)

## 阶段八：插件模块

- [x] Task 31: 渠道活码插件 (channel-code)
- [x] Task 32: 欢迎语插件 (greeting)
- [x] Task 33: 自动拉群插件 (room-auto-pull)
- [x] Task 34: 标签建群插件 (room-tag-pull)
- [x] Task 35: 群欢迎语插件 (room-welcome)
- [x] Task 36: 客户消息群发插件 (contact-message-batch-send)
- [x] Task 37: 群消息群发插件 (room-message-batch-send)
- [x] Task 38: 客户转接插件 (contact-transfer)
- [x] Task 39: 统计插件 (statistic)
- [x] Task 40: 裂变插件 (work-fission)

## 阶段九：路由注册与服务启动

- [x] Task 41: 路由注册
- [x] Task 42: 服务启动与优雅关闭

## 阶段十：测试与文档

- [ ] Task 43: 单元测试
  - [ ] 43.1 Service 层单元测试
  - [ ] 43.2 Logic 层单元测试
  - [ ] 43.3 中间件单元测试
  - [ ] 43.4 工具函数单元测试

- [ ] Task 44: API 文档
  - [ ] 44.1 集成 swaggo/swag
  - [ ] 44.2 Swagger 注解
  - [ ] 44.3 OpenAPI 文档

- [ ] Task 45: 集成测试
  - [ ] 45.1 API 接口集成测试
  - [ ] 45.2 队列处理集成测试
  - [ ] 45.3 事件驱动集成测试

# 编译验证

- [x] go build ./... 编译通过
