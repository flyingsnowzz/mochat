# Handler 目录重构迁移指南

## 进度说明

已完成：
- ✅ 创建目录结构
- ✅ system 模块（corp、user、role、menu）文件复制
- ✅ contact 模块文件复制
- ✅ organization 模块文件复制（需要更新 package 名称）

待完成：
- ⏳ 更新所有迁移文件的 package 名称
- ⏳ 迁移 content 模块
- ⏳ 迁移 marketing 模块
- ⏳ 迁移 analysis 模块
- ⏳ 迁移 platform 模块
- ⏳ 迁移 sidebar 到 client
- ⏳ 迁移 common 公共接口
- ⏳ 更新 router 引用
- ⏳ 测试验证

## 当前问题

由于 macOS 的 sed 命令与 Linux 有差异，批量替换没有正确执行。需要手动更新 package 名称。

## 手动修正步骤

### 1. 修正 organization 模块
```bash
cd /Users/zhanglei/MyProjects/mochat/api-server-go
```

需要手动修改以下文件：
- `internal/handler/dashboard/organization/employee.go`: package dashboard → package organization
- `internal/handler/dashboard/organization/room.go`: package dashboard → package organization

### 2. 继续迁移其他模块

按照相同模式迁移剩余模块。

## 更新的目录结构（最终）

```
internal/handler/
├── dashboard/                     # 后台管理端
│   ├── system/                    # 系统管理
│   │   ├── corp.go
│   │   ├── user.go
│   │   ├── role.go
│   │   └── menu.go
│   │
│   ├── contact/                   # 客户管理
│   │   ├── contact.go
│   │   ├── contact_field.go
│   │   ├── contact_tag.go
│   │   └── contact_group.go
│   │
│   ├── organization/              # 组织架构
│   │   ├── department.go
│   │   ├── employee.go
│   │   └── room.go
│   │
│   ├── content/                   # 内容管理
│   │   ├── medium.go
│   │   ├── medium_group.go
│   │   └── greeting.go
│   │
│   ├── marketing/                 # 营销工具
│   │   ├── channel_code.go
│   │   ├── room_welcome.go
│   │   ├── auto_pull.go
│   │   ├── batch_send.go
│   │   └── fission.go
│   │
│   ├── analysis/                  # 数据分析
│   │   ├── statistic.go
│   │   └── index.go
│   │
│   ├── platform/                  # 平台配置
│   │   ├── agent.go
│   │   └── official_account.go
│   │
│   └── common/                    # 公共接口
│       └── common.go
│
└── client/                        # 客户端（侧边栏）
    ├── contact.go
    ├── room.go
    ├── agent.go
    ├── medium.go
    └── common.go
```

## Router 引用更新

需要在以下文件中更新 import 路径：
- `internal/router/handlers.go`
- `internal/router/dashboard_routes.go`
- `internal/router/sidebar.go`

### 更新示例

**修改前**:
```go
import (
    "mochat-api-server/internal/handler/dashboard"
    dashboardPlugin "mochat-api-server/internal/handler/dashboard/plugin"
    "mochat-api-server/internal/handler/sidebar"
)
```

**修改后**:
```go
import (
    dashboardSystem "mochat-api-server/internal/handler/dashboard/system"
    dashboardContact "mochat-api-server/internal/handler/dashboard/contact"
    dashboardOrg "mochat-api-server/internal/handler/dashboard/organization"
    // ... 其他模块
    clientContact "mochat-api-server/internal/handler/client/contact"
    // ... 其他客户端模块
)
```

## 测试验证清单

- [ ] 编译通过: `go build ./...`
- [ ] 后台管理端接口可访问
- [ ] 客户端接口可访问
- [ ] 测试用例通过
- [ ] 前端功能正常
