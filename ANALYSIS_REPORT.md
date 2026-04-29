# api-server(PHP) vs api-server-go(Go) 接口一致性分析报告

> 生成日期: 2026-04-29
> 范围: Dashboard + Sidebar + Operation 全端接口

---

## 目录

1. [整体评估](#1-整体评估)
2. [路由覆盖对比](#2-路由覆盖对比)
3. [核心模块实现差异详解](#3-核心模块实现差异详解)
4. [全部缺失路由清单](#4-全部缺失路由清单)
5. [路由路径不一致清单](#5-路由路径不一致清单)
6. [Go端新增路由](#6-go端新增路由)
7. [响应格式差异](#7-响应格式差异)
8. [中间件/权限差异](#8-中间件权限差异)
9. [综合建议](#9-综合建议)

---

## 1. 整体评估

| 维度 | 评估 |
|------|------|
| **Dashboard 路由覆盖** | ~70%（约 70 条 vs PHP 125 条） |
| **Sidebar 路由覆盖** | ~55%（约 13 条 vs PHP 20 条） |
| **Plugin 模块覆盖** | ~15%（仅 4 个模块部分实现，PHP 有 9 个插件模块） |
| **Operation 路由覆盖** | ~15%（1 条 vs PHP 7+ 条） |
| **实现逻辑一致度** | ~50%（大量核心逻辑缺失或简化） |
| **整体完成度评估** | **约 50-60%** |

---

## 2. 路由覆盖对比

### 2.1 Dashboard 端

| 模块 | PHP 路由数 | Go 路由数 | 覆盖度 | 备注 |
|------|-----------|-----------|--------|------|
| Corp (企业管理) | 7 | 9 | ~100% | 路由数量超 PHP，但实现逻辑有差异 |
| User (用户管理) | 10 | 10 | ~100% | |
| Menu (菜单管理) | 8 | 9 | ~100% | |
| Role (角色管理) | 11 | 12 | ~100% | |
| WorkContact (客户管理) | 20 | 16 | ~80% | 缺少部分标签/字段路由 |
| WorkDepartment (部门) | 4 | 4 | ~100% | |
| WorkEmployee (员工) | 4 | 5 | ~100% | |
| WorkRoom (客户群) | 6 | 5 | ~80% | |
| Medium (素材) | 10 | 10 | ~100% | |
| Index (首页) | 2 | 2 | ~100% | 但实现逻辑差异大 |
| ChatTool (聊天工具) | 1 | 1 | ~100% | |
| OfficialAccount (公众号) | 6 | 4 | ~67% | 缺少 messageEventCallback, authRedirect |
| Agent (第三方应用) | 3 | 3 | ~100% | |
| Greeting (欢迎语) | 5 | 5 | ~100% | |
| ChannelCode (渠道码) | 11 | 9 | ~82% | 缺少 group/detail, group/move |
| Statistic (数据统计) | 5 | 5 | ~100% | |
| RoomWelcome (入群欢迎语) | 6 | 5 | ~83% | 缺少 select 路由 |
| RoomAutoPull (自动拉群) | 4 | 4 | ~100% | 但 corpID 硬编码为1 |
| RoomTagPull (标签建群) | 9 | 4 | ~44% | 仅实现最基本的 4 条 |
| **Plugin 缺失模块** | | | | |
| ContactMessageBatchSend (客户群发) | 8 | 0 | 0% | 完全缺失 |
| ContactTransfer (客户继承) | 7 | 0 | 0% | 完全缺失 |
| RoomMessageBatchSend (群群发) | 7 | 0 | 0% | 完全缺失 |
| WorkFission (任务宝) | 18 | 0 | 0% | 完全缺失 |

### 2.2 Sidebar 端

| 路由 | PHP | Go | 状态 |
|------|-----|----|------|
| `/sidebar/workContact/show` | ✅ | ✅ | 一致 |
| `/sidebar/workContact/detail` | ✅ | ✅ | 一致 |
| `/sidebar/workContact/update` | ✅ | ✅ | 一致 |
| `/sidebar/workContact/track` | ✅ | ✅ | 一致 |
| `/sidebar/workRoom/roomManage` | ✅ | ✅ | 一致 |
| `/sidebar/agent/auth` | ✅ | ✅ | 一致 |
| `/sidebar/agent/oauth` | ✅ | ✅ | 一致 |
| `/sidebar/agent/jssdkConfig` | ✅ | ✅ | 一致 |
| `/sidebar/medium/index` | ✅ | ✅ | 一致 |
| `/sidebar/medium/mediaIdUpdate` | ✅ | ✅ | 一致 |
| `/sidebar/mediumGroup/index` | ✅ | ✅ | 一致 |
| `/sidebar/common/upload` | ✅ | ✅ | 一致 |
| `/sidebar/wxJSSDK/config` | ✅ | ✅ | 一致 |
| `/sidebar/workContact/tag/allTag` | ✅ | ❌ | **缺失** |
| `/sidebar/workContact/tagGroup/index` | ✅ | ❌ | **缺失** |
| `/sidebar/workContact/fieldPivot/index` | ✅ | ❌ | **缺失** |
| `/sidebar/workContact/fieldPivot/update` | ✅ | ❌ | **缺失** |
| `/sidebar/workContact/processStatus/index` | ✅ | ❌ | **缺失** |
| `/sidebar/workContact/processStatus/update` | ✅ | ❌ | **缺失** |
| `/sidebar/medium/group/index` | ✅ | ❌ | **缺失** |

### 2.3 Operation 端

| 路由 | PHP | Go | 状态 |
|------|-----|----|------|
| `/operation/officialAccount/authRedirect` | ✅ | ✅ | 一致 |
| `/operation/workFission/*` (6条) | ✅ | ❌ | **完全缺失** |

---

## 3. 核心模块实现差异详解

### 3.1 Corp (企业管理)

#### Index - 获取企业列表

| 维度 | PHP | Go |
|------|-----|----|
| **搜索** | 支持 `corpName` 模糊搜索 | ❌ **缺少** |
| **排序** | `id desc` | ❌ **缺少** |
| **用户过滤** | 非超级管理员限制 `corpIds` | ❌ **缺少** |
| **分页参数** | `page`, `perPage`, 默认10 | `page`, `pageSize`, 默认20 |
| **字段映射** | wxCorpId, corpId, corpName, chatStatus 等 | 返回原始模型 |
| **响应格式** | 自定义 page 对象 + list | PageResult 通用格式 |

#### Store - 创建企业

| 维度 | PHP | Go |
|------|-----|----|
| **重复校验** | 检查最多只能添加一个企业 | ❌ **缺少** |
| **token 生成** | 自动生成随机回调 token | ❌ **缺少** |
| **encoding_aes_key** | 自动生成 43 位随机密钥 | ❌ **缺少** |
| **回调 URL** | 自动构建 event_callback URL | ❌ **缺少** |
| **DB 事务** | 使用数据库事务 | ✅ 无事务直接写入 |
| **Redis 绑定** | 存储 `mc:user.{id}` 用户-企业关系 | ❌ **缺少** |
| **队列触发** | 触发 EmployeeApply 同步通讯录 | ❌ **缺少** |

#### Bind - 绑定企业 (**逻辑完全不同**)

| 维度 | PHP | Go |
|------|-----|----|
| **语义** | **用户-企业绑定**: 选择企业后存入 Redis `mc:user.{id}` = `{corpId}-{employeeId}` | **企业微信配置更新**: wx_corpid, employee_secret 等 |
| **参数** | `corpId` (用户选择的企微ID) | `wxCorpid`, `employeeSecret`, `contactSecret` 等 |
| **HTTP 路径** | `POST /dashboard/corp/bind` | `POST /dashboard/corp/bind/:id` |
| **结论** | 🔴 **严重错误** | Go 把 **用户切换企业** 实现成了 **更新企业微信配置** |

> **注意**: Go 在 `POST /dashboard/user/corp/bind` 实现了正确的用户-企业绑定功能，但路径与 PHP 不同。

#### Update - 更新企业

| 维度 | PHP | Go |
|------|-----|----|
| **参数** | `corpId`, `corpName`, `wxCorpId`, `employeeSecret`, `contactSecret` | `id`(路径参数) + 任意 JSON map |
| **缓存清理** | 调用 `WeWorkFactory.unbindApp()` 清理应用实例 | ❌ **缺少** |
| **验证** | 检查企业是否存在 | ✅ 直接更新 |

#### Select - 企业选择列表

| 维度 | PHP | Go |
|------|-----|----|
| **超级管理员** | 根据 tenantId 查询 | ✅ 类似逻辑 |
| **非超级管理员** | 通过 `WorkEmployeeService` 过滤员工所属企业 | ❌ **直接返回所有企业，缺少过滤** |
| **搜索** | 支持 `corpName` 搜索 | ❌ **缺少** |

---

### 3.2 User (用户管理)

| 方法 | 差异 |
|------|------|
| **Auth** | PHP 返回 `{token, expire}`；Go 返回 `{token, userInfo}`，响应格式不一致 |
| **Store** | PHP 有复杂密码加密；Go 直接存储原始密码字符串 |
| **PasswordUpdate** | PHP 验证旧密码；Go 直接更新新密码 |
| **Index** | PHP 支持 `phone`, `status` 搜索；Go 只简单分页 |

---

### 3.3 Menu (菜单管理)

| 方法 | 差异 |
|------|------|
| **Index** | PHP 支持名称模糊搜索并递归关联父级菜单下的所有子菜单；Go 仅按名称精确搜索 |
| **Index 响应** | PHP 返回 page/list 自定义格式；Go 使用通用 PageResult 格式 |
| **Index 字段** | PHP 返回 `menuId`, `menuPath`, `level`, `levelName`, `parentId`, `icon`, `status`, `operateName`, `updatedAt`, `children`；Go 类似但 `status` 返回字符串 |

PHP 的 Menu/Index 搜索逻辑更复杂：根据名称模糊匹配到菜单后，还会向上递归找出所有祖先菜单，再向下找出所有子孙菜单，形成一个完整的子树。

---

### 3.4 Role (角色管理)

| 方法 | 差异 |
|------|------|
| **Index** | PHP 通过 `WorkEmployeeService` 统计每个角色下的员工数量 (`employeeNum`)；Go ❌ **缺少员工数量统计** |
| **Index 搜索** | PHP 支持 `name` 搜索；Go ❌ **缺少** |
| **ShowEmployee** | PHP 调用 `RbacUserRole` 关联查询，返回完整的用户-角色关联表数据；Go 直接查 `RbacUserRole` 表 |

---

### 3.5 Medium (素材管理)

| 方法 | 差异 |
|------|------|
| **Store** | PHP 复杂内容构建（含 media_id 上传）；Go 简单 JSON 存储 |
| **Update** | PHP 检查素材是否存在；Go 直接更新（虽有 GetByID 但无存在性校验返回错误） |
| **GroupUpdate** | PHP `/dashboard/medium/groupUpdate` - 通过 `id`, `mediumGroupId` 单独更新分组；Go 实现了一致逻辑 |

---

### 3.6 Greeting (欢迎语)

| 方法 | 差异 |
|------|------|
| **Index 响应格式** | PHP 返回 `{page: {perPage, total, totalPage}, list: []}`；Go 返回 `{list: [], hadGeneral, hadEmployees, page: {perPage, total, totalPage, currentPage}, currentPage, pageSize, total}`（Go 额外返回了 hadGeneral, hadEmployees 和兼容字段） |
| **权限中间件** | PHP 有 `@Middleware(PermissionMiddleware::class)` 权限校验；Go ❌ **未添加 permission 中间件** |
| **参数校验** | PHP 使用 `ValidateSceneTrait` 验证；Go 手动校验 |

---

### 3.7 ChannelCode (渠道码)

| 方法 | 差异 |
|------|------|
| **Store** | PHP 调用企业微信 API 创建渠道码并保存返回的 `qrcode_url` 和 `config_id`；Go **使用占位符 URL** `https://dummyimage.com/240x240/...`，❌ **无实际企业微信 API 调用** |
| **Index** | 实现逻辑基本一致（分页、搜索、分组筛选） |
| **Group/Detail** | PHP 有 `/dashboard/channelCode/group/detail`；Go ❌ **缺失** |
| **Group/Move** | PHP 有 `/dashboard/channelCode/group/move`；Go 实现了但命名不同 |

---

### 3.8 RoomWelcome (入群欢迎语)

| 方法 | 差异 |
|------|------|
| **整体逻辑** | PHP 使用 `plugin/mochat/room-welcome` 模块，有完整的 CRUD + select 路由 |
| **CorpID** | Go 中 **硬编码为 `uint(1)`**（见 `room_welcome.go:31`），❌ **未从认证上下文获取** |
| **状态码** | Go 使用魔数（400, 500）而非定义常量 |
| **Select** | PHP 有 `/dashboard/roomWelcome/select` 路由；Go ❌ **缺失** |

---

### 3.9 WorkRoomAutoPull (自动拉群)

| 方法 | 差异 |
|------|------|
| **CorpID** | Go 中 **硬编码为 `uint(1)`**（注释说明"暂时使用默认值1"）|
| **Show** | Go 中 employees/tags/rooms 字段解析 **返回空数组**（标注"暂时返回空数组"）|
| **PHP 对比** | PHP 有完整的 JSON 字段解析和复杂的队列处理逻辑，Go 仅为简单 CURD |

---

### 3.10 RoomTagPull (标签建群)

| 维度 | PHP | Go |
|------|-----|----|
| **路由数量** | 9 条（index/show/store/destroy/chooseContact/filterContact/showContact/roomList/remindSend） | 4 条（index/create/detail/contactDetail） |
| **路由路径** | `/dashboard/roomTagPull/{action}` | `/dashboard/roomTagPull/{action}` |
| **CorpID** | 正常从上下文获取 | 硬编码为 1 |
| **路由方法** | GET/POST | GET/POST |

#### PHP 有但 Go 缺失的 5 条路由：

| 路由 | 方法 | 功能 |
|------|------|------|
| `/dashboard/roomTagPull/chooseContact` | POST | 选择联系人 |
| `/dashboard/roomTagPull/filterContact` | POST | 筛选联系人 |
| `/dashboard/roomTagPull/showContact` | GET | 查看已选联系人 |
| `/dashboard/roomTagPull/roomList` | GET | 群列表 |
| `/dashboard/roomTagPull/remindSend` | POST | 提醒发送 |

---

### 3.11 Dashboard Index (首页数据)

| 维度 | PHP | Go |
|------|-----|----|
| **数据来源** | `CorpDayData` 日汇总表（定时任务预先统计） | 实时查询原始业务表 |
| **总客户数** | `weChatContactNum` - 通过 `WorkContactEmployee` 表统计 | `totalContact` - 通过 `WorkContact` 表统计 |
| **总群聊数** | `weChatRoomNum` ✅ | `totalRoom` ✅ |
| **总群成员数** | `roomMemberNum` ✅ | ❌ **缺少** |
| **总员工数** | `corpMemberNum` ✅ | ❌ **缺少** |
| **今日新增客户** | `addContactNum` ✅ | `todayAddContact` ✅ |
| **昨日新增客户** | `lastAddContactNum` ✅ | ❌ **缺少** |
| **今日流失客户** | `lossContactNum` ✅ | `todayLossContact` ✅ |
| **昨日流失客户** | `lastLossContactNum` ✅ | ❌ **缺少** |
| **本月新增客户** | `addFriendsNum` ✅ | ❌ **缺少** |
| **上月新增客户** | `lastAddFriendsNum` ✅ | ❌ **缺少** |
| **同步时间** | `updateTime` ✅ | ❌ **缺少** |
| **LineChat** | ✅ 折线图数据 | ✅ 折线图数据（最近7天） |

**结论**: Go 实现的功能仅为 PHP 的 **30%**，缺少月度聚合、环比对比、总成员数等关键指标。

---

### 3.12 Contact (客户管理)

| 方法 | 差异 |
|------|------|
| **Index** | Go 响应字段丰富（含头像、员工名称、标签、群聊等）；PHP 使用 Logic 层处理复杂关联 |
| **SynContact** | PHP 调用企业微信 API 同步；Go **仅留空 goroutine 注释** |
| **Source** | PHP 支持 `corpId` 筛选；Go 逻辑基本一致 |
| **LossContact** | PHP 通过 `employee_status` 判断；Go 通过 `deleted_at` 判断（可能结果不同） |

---

### 3.13 Agent (第三方应用)

| 方法 | 差异 |
|------|------|
| **PHP 全量路由** | `/dashboard/workAgent/store`, `/dashboard/workAgent/txtVerifyShow`, `/dashboard/workAgent/txtVerifyUpload` |
| **Go 路由** | `/dashboard/agent/store`, `/dashboard/agent/txtVerifyShow`, `/dashboard/agent/txtVerifyUpload` |
| **路径问题** | Go 去掉了 `work` 前缀，导致与 PHP 不兼容 |

---

### 3.14 OfficialAccount (公众号)

| 缺失路由 | 说明 |
|---------|------|
| `/dashboard/officialAccount/messageEventCallback` | 消息事件回调 |
| `/dashboard/officialAccount/authRedirect` | 授权重定向 |
| `/dashboard/officialAccount/set` | Go 已实现 `POST /dashboard/officialAccount/set` ✅ |

---

## 4. 全部缺失路由清单

### 4.1 Sidebar 缺失（7条）

```http
GET    /sidebar/workContact/tag/allTag          # 客户标签列表
GET    /sidebar/workContact/tagGroup/index       # 客户标签组列表
GET    /sidebar/workContact/fieldPivot/index     # 自定义字段值列表
PUT    /sidebar/workContact/fieldPivot/update    # 更新自定义字段值
GET    /sidebar/workContact/processStatus/index  # 跟进状态列表
PUT    /sidebar/workContact/processStatus/update # 更新跟进状态
GET    /sidebar/medium/group/index               # 素材分组列表
```

### 4.2 Plugin 模块缺失（约45条）

#### contact-message-batch-send (客户群发) - 8条

```http
GET    /dashboard/contactMessageBatchSend/index
POST   /dashboard/contactMessageBatchSend/store
GET    /dashboard/contactMessageBatchSend/show
DELETE /dashboard/contactMessageBatchSend/destroy
POST   /dashboard/contactMessageBatchSend/remind
GET    /dashboard/contactMessageBatchSend/showRoom
GET    /dashboard/contactMessageBatchSend/employeeSendIndex
GET    /dashboard/contactMessageBatchSend/contactReceiveIndex
```

#### contact-transfer (客户继承) - 7条

```http
GET    /dashboard/contactTransfer/index
GET    /dashboard/contactTransfer/info
GET    /dashboard/contactTransfer/log
GET    /dashboard/contactTransfer/room
GET    /dashboard/contactTransfer/unassignedList
POST   /dashboard/contactTransfer/saveUnassignedList
POST   /dashboard/contactTransfer/transferRoom
```

#### room-message-batch-send (群群发) - 7条

```http
GET    /dashboard/roomMessageBatchSend/index
POST   /dashboard/roomMessageBatchSend/store
GET    /dashboard/roomMessageBatchSend/show
DELETE /dashboard/roomMessageBatchSend/destroy
POST   /dashboard/roomMessageBatchSend/remind
GET    /dashboard/roomMessageBatchSend/roomOwnerSendIndex
GET    /dashboard/roomMessageBatchSend/roomReceiveIndex
```

#### room-tag-pull (标签建群) - 5条

```http
POST   /dashboard/roomTagPull/chooseContact
POST   /dashboard/roomTagPull/filterContact
GET    /dashboard/roomTagPull/showContact
GET    /dashboard/roomTagPull/roomList
POST   /dashboard/roomTagPull/remindSend
```

#### work-fission (任务宝) - 18条

```http
# Dashboard
GET    /dashboard/workFission/index
GET    /dashboard/workFission/show
POST   /dashboard/workFission/store
PUT    /dashboard/workFission/update
DELETE /dashboard/workFission/destroy
GET    /dashboard/workFission/statistics
GET    /dashboard/workFission/info
GET    /dashboard/workFission/invite
GET    /dashboard/workFission/inviteData
GET    /dashboard/workFission/inviteDetail
POST   /dashboard/workFission/chooseContact

# Operation
GET    /operation/workFission/auth
POST   /operation/workFission/inviteFriends
POST   /operation/workFission/taskData
POST   /operation/workFission/receive
GET    /operation/workFission/poster
GET    /operation/workFission/openUserInfo
```

#### OfficialAccount - 2条

```http
POST   /dashboard/officialAccount/messageEventCallback  # 消息事件回调
GET    /dashboard/officialAccount/authRedirect           # 授权重定向
```

---

## 5. 路由路径不一致清单

约有 7 组路由路径在 PHP 和 Go 之间不一致：

| PHP 路径 | Go 路径 | 影响 |
|----------|---------|------|
| `/dashboard/workContact/field/*` | `/dashboard/contactField/*` | 前端需适配 |
| `/dashboard/workContact/tag/*` | `/dashboard/workContactTag/*` | 前端需适配 |
| `/dashboard/workContact/tagGroup/*` | `/dashboard/workContactTagGroup/*` | 前端需适配 |
| `/dashboard/workContact/room/*` | `/dashboard/workContactRoom/*` | 前端需适配 |
| `/dashboard/workContact/fieldPivot/*` | `/dashboard/contactFieldPivot/*` | 前端需适配 |
| `/dashboard/workAgent/*` | `/dashboard/agent/*` | 前端需适配 |
| `/dashboard/roomAutoPull/*` | `/dashboard/workRoomAutoPull/*` | 前端需适配 |

> **建议**: 将 Go 路径统一为 PHP 现有格式，保持前后端兼容。

---

## 6. Go端新增路由

以下路由在 PHP 端不存在，仅在 Go 端有实现（可能是未来的扩展）：

| Go 路由 | 方法 | 说明 |
|---------|------|------|
| `/dashboard/health` | GET | 健康检查 |
| `/health` | GET | 全局健康检查，返回 `{"status":"ok"}` |
| `/dashboard/user/corp/bind` | POST | 用户-企业绑定（PHP 在 `/dashboard/corp/bind`） |
| contactField 的额外方法 | GET/PUT | 如 `getByWxExternalUserID`, `listByOffset` 等 |
| contactTag 的额外方法 | GET/PUT | 如 `updateByID`, `listByOrder` 等 |

---

## 7. 响应格式差异

### 7.1 分页响应

| 维度 | PHP | Go |
|------|-----|----|
| **通用格式** | `{page: {perPage, total, totalPage}, list: []}` | `{list: [], total, page, pageSize}` (PageResult) |
| **Menu/Index** | 自定义 tree 结构 + 自定义 page 对象 | PageResult 通用格式 |
| **Greeting/Index** | 自定义 page 对象 | 额外返回 hadGeneral, hadEmployees 等字段 |

### 7.2 User/Auth 登录响应

| 字段 | PHP | Go |
|------|-----|----|
| token | ✅ | ✅ |
| expire | ✅ | ❌ **缺少** |
| userInfo | ❌ | ✅ **额外返回** |

---

## 8. 中间件/权限差异

| 维度 | PHP | Go |
|------|-----|----|
| **认证中间件** | `@Middleware(DashboardAuthMiddleware::class)` | JWT 中间件（whitelist 机制） |
| **权限中间件** | `@Middleware(PermissionMiddleware::class)` | permission 中间件（仅部分路由使用） |
| **白名单** | `config/autoload/framework.php` 中 `auth_white_routes` | `whitelist.go` 中 `getDashboardWhiteRoutes()` |
| **覆盖范围** | 全量标注在注解上 | 仅为部分路由组添加 permission |

### 权限中间件不一致的地方

以下 PHP 有 `@Middleware(PermissionMiddleware::class)` 但 Go **缺少** permission 中间件：

| 路由 | 说明 |
|------|------|
| `/dashboard/greeting/*` | 欢迎语管理 |
| `/dashboard/channelCode/*` | 渠道码管理 |
| `/dashboard/roomWelcome/*` | 进群欢迎语 |
| `/dashboard/workRoomAutoPull/*` | 自动拉群 |
| `/dashboard/roomTagPull/*` | 标签建群 |
| `/dashboard/statistic/*` | 数据统计 |
| `/dashboard/medium/*` | 素材管理 |

---

## 9. 综合建议

### 9.1 优先级排序

| 优先级 | 问题 | 影响 |
|--------|------|------|
| **P0** | Corp/Bind 逻辑完全错误 | 用户无法正常切换企业 |
| **P0** | Corp/Store 缺少核心逻辑 | 企业创建流程中断 |
| **P0** | 全部 Plugin 模块缺失 | 群发、继承、裂变等功能不可用 |
| **P0** | RoomWelcome/AutoPull/RoomTagPull 中 CorpID 硬编码 | 多企业场景无法使用 |
| **P1** | Corp/Index/Select 缺少搜索和过滤 | 功能不全 |
| **P1** | Dashboard/Index 数据统计功能严重不足 | 首页数据不准确 |
| **P1** | Contact/SynContact 空实现 | 客户同步功能不可用 |
| **P2** | 路由路径不一致（7组） | 前端需额外适配 |
| **P2** | Sidebar 缺失 7 条路由 | 侧边栏功能不全 |
| **P2** | 权限中间件缺失 | 安全性不足 |
| **P3** | 响应格式差异 | 前端兼容性问题 |
| **P3** | Corp/Update 缺少缓存清理 | 企业配置更新后缓存脏数据 |

### 9.2 修复建议

1. **P0 问题优先修复**：修正 Corp/Bind 逻辑，补充 Corp/Store 核心逻辑，移植 Plugin 模块
2. **统一 CorpID 获取方式**：废弃硬编码，从 JWT 中间件统一获取
3. **统一路由路径**：Go 端路径对齐 PHP，或保持 Go 路径并提供兼容映射
4. **补充权限校验**：为所有需要权限的接口添加 permission 中间件
5. **对齐响应格式**：统一分页和错误响应格式

---

*本报告基于代码静态分析生成，建议结合实际运行测试进行验证。*
