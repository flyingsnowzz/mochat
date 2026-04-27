# 路由重构指南

## 重构概述

本次重构将原有的单一路由文件 `router.go` 拆分为多个模块化的文件，提高了代码的可维护性和可读性。

## 文件结构

```
internal/router/
├── router.go              # 原有文件（保留作为参考）
├── base.go                # Router 结构体和基础方法
├── setup.go               # 兼容旧接口的初始化函数
├── dashboard.go            # Dashboard 路由注册
├── dashboard_routes.go     # Dashboard 详细路由配置
├── sidebar.go             # Sidebar 路由注册
├── operation.go           # 运营操作路由
├── whitelist.go           # 白名单路由配置
└── handlers.go            # 处理器初始化
```

## 主要改进

### 1. 职责分离

- **base.go**: 定义 Router 结构体和核心配置方法
- **dashboard.go**: Dashboard 路由的顶层注册
- **dashboard_routes.go**: Dashboard 各模块的详细路由配置
- **sidebar.go**: Sidebar 路由注册
- **operation.go**: 运营操作路由
- **whitelist.go**: 白名单路由集中管理
- **handlers.go**: 统一的处理器初始化

### 2. 面向对象设计

```go
type Router struct {
    engine          *gin.Engine
    config          *config.Config
    db              *gorm.DB
    // 中间件
    dashboardAuth   gin.HandlerFunc
    sidebarAuth     gin.HandlerFunc
    permission      gin.HandlerFunc
    // 各种 Handler
    corpHandler     *dashboard.CorpHandler
    userHandler     *dashboard.UserHandler
    // ...
}

func (r *Router) Setup() *gin.Engine {
    r.setMode()
    r.useMiddlewares()
    r.registerStaticRoutes()
    r.registerHealthRoutes()
    r.initHandlers()
    r.registerDashboardRoutes()
    r.registerSidebarRoutes()
    r.registerOperationRoutes()
    return r.engine
}
```

### 3. 向后兼容

保留原有的 `SetupRouter` 函数签名，确保现有代码无需修改：

```go
func SetupRouter(cfg *config.Config, db *gorm.DB) *gin.Engine {
    router := NewRouter(cfg, db)
    return router.Setup()
}
```

## 使用方法

### 方式一：使用新接口（推荐）

```go
router := router.NewRouter(cfg, db)
engine := router.Setup()
```

### 方式二：使用兼容旧接口

```go
engine := router.SetupRouter(cfg, db)
```

## 扩展指南

### 添加新的路由组

1. 在对应的路由文件中添加注册方法：

```go
// dashboard_routes.go
func (r *Router) registerNewFeatureRoutes(group *gin.RouterGroup) {
    group.GET("/index", r.newFeatureHandler.Index)
    group.POST("/store", r.newFeatureHandler.Store)
}
```

2. 在 dashboard.go 中注册：

```go
// 注册新功能路由
r.registerNewFeatureRoutes(group.Group("/newFeature"))
```

3. 在 handlers.go 中初始化 Handler：

```go
func (r *Router) initHandlers() {
    // ...
    r.newFeatureHandler = dashboard.NewNewFeatureHandler(r.db)
}
```

### 添加新的处理器类型

1. 在 base.go 的 Router 结构体中添加字段
2. 在 handlers.go 中初始化
3. 在对应的路由文件中注册路由

## 优势

1. **可维护性**: 代码按功能模块分离，易于查找和修改
2. **可测试性**: 每个路由组可以独立测试
3. **可扩展性**: 添加新功能只需在相应文件中添加代码
4. **可读性**: 文件职责清晰，代码结构一目了然
5. **向后兼容**: 不影响现有代码的使用

## 迁移建议

1. **短期**: 可以同时保留 router.go 和新的结构，逐步迁移
2. **中期**: 将 router.go 标记为 deprecated，引导使用新的结构
3. **长期**: 移除 router.go，完全使用新的模块化结构

## 注意事项

1. 所有路由配置保持与原文件一致
2. 中间件逻辑完全保留
3. 白名单路由集中管理在 whitelist.go
4. Handler 初始化逻辑集中在 handlers.go
