* [x] Go 项目结构完整创建，包含 cmd/internal/pkg/config 等标准目录

* [x] go.mod 正确初始化，所有依赖声明完整

* [x] Makefile 包含 build/run/test/lint/swag/docker 目标

* [x] Dockerfile 多阶段构建成功，可生成可运行镜像

* [x] docker-compose.yml 可一键启动完整系统（API + MySQL + Redis）

* [x] config.example.yaml 包含所有配置项，与原项目 .env.example 对应

* [x] viper 配置加载正确，支持环境变量覆盖

* [x] zap 日志集成，输出格式化和结构化日志

* [x] 统一错误码体系与原项目 AppErrCode 对应

* [x] 统一响应格式 {code, data, msg} 与原项目一致

* [x] GORM 连接 MySQL 成功，连接池参数与原项目一致

* [x] 表前缀 mc\_ 正确设置

* [x] 分表路由 mc\_work\_unionid\_external\_userid\_mapping\_0\~9 正确实现

* [x] go-redis 连接池配置正确

* [x] 缓存操作封装完整（Get/Set/Del/Remember）

* [x] 所有 35+ 核心表 GORM Model 定义完整，字段与原 SQL 一致

* [x] 所有插件表 GORM Model 定义完整

* [x] Dashboard JWT 认证（User 模型）正确实现

* [x] Sidebar JWT 认证（WorkEmployee 模型）正确实现

* [x] DashboardAuthMiddleware 正确拦截和验证

* [x] SidebarAuthMiddleware 正确拦截和验证

* [x] 认证白名单路由配置正确

* [x] PermissionMiddleware 菜单权限检查正确

* [x] 三级数据权限（全企业/本部门/本人）正确实现

* [x] RBAC 工具函数正确封装

* [x] CoreMiddleware 请求上下文注入正确

* [x] CORS 中间件正确配置

* [x] 企业微信 SDK 封装完整（通讯录/客户/消息/应用）

* [x] 企业微信回调解密/加密正确实现

* [x] WeWorkFactory 按企业创建 SDK 实例正确

* [x] 素材上传（image/voice/video/file）+ 缓存正确实现

* [x] 微信开放平台公众号授权流程正确实现

* [x] asynq 异步队列集成，所有通道配置正确

* [x] 队列投递和消费正确工作

* [x] 失败重试和死信队列正确实现

* [x] robfig/cron 定时任务调度器正确集成

* [x] 所有定时任务正确实现（CorpData/MediaIdUpdate/WorkEmployeeStatistic/SyncWorkAgent）

* [x] 动态 Crontab 注册机制正确实现

* [x] EventBus 事件注册、触发、异步分发正确实现

* [x] 所有事件类型和监听器正确实现

* [x] Storage 接口定义完整

* [x] 所有文件存储驱动正确实现（local/oss/cos/s3/minio/qiniu）

* [x] 异步文件上传队列正确实现

* [x] 企业模块所有 Handler 正确实现，API 响应与原项目一致

* [x] 租户模块所有 Handler 正确实现

* [x] 用户模块所有 Handler 正确实现（登录/CRUD/密码管理）

* [x] RBAC 模块所有 Handler 正确实现（菜单/角色/权限）

* [x] 客户模块所有 Dashboard Handler 正确实现

* [x] 客户模块所有 Sidebar Handler 正确实现

* [x] 部门模块所有 Handler 正确实现

* [x] 员工模块所有 Handler 正确实现

* [x] 客户群模块所有 Handler 正确实现

* [x] 应用模块所有 Handler 正确实现

* [x] 素材库模块所有 Handler 正确实现

* [x] 首页模块所有 Handler 正确实现

* [x] 侧边工具栏模块所有 Handler 正确实现

* [x] 公众号模块所有 Handler 正确实现

* [x] 通用模块所有 Handler 正确实现（上传/JSSDK）

* [x] 渠道活码插件完整实现

* [x] 欢迎语插件完整实现

* [x] 自动拉群插件完整实现

* [x] 标签建群插件完整实现

* [x] 群欢迎语插件完整实现

* [x] 客户消息群发插件完整实现

* [x] 群消息群发插件完整实现

* [x] 客户转接插件完整实现

* [x] 统计插件完整实现

* [x] 裂变插件完整实现

* [x] 所有 Dashboard 路由正确注册并应用中间件

* [x] 所有 Sidebar 路由正确注册并应用中间件

* [x] 公开路由（登录/回调/白名单）正确注册

* [x] 静态文件服务正确配置

* [x] main.go 入口正确初始化所有组件

* [x] HTTP 服务在端口 9501 正确启动

* [x] 优雅关闭正确实现

* [x] 健康检查端点正确响应

* [ ] 单元测试覆盖率 >= 70%

* [ ] 所有 Swagger 注解添加完整

* [ ] OpenAPI 文档可正常访问

* [ ] API 接口集成测试通过

* [ ] 队列处理集成测试通过

* [ ] 事件驱动集成测试通过

