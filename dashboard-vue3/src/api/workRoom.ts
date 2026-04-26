// 客户群模块
import request from '@/utils/http'

// 客户群分组列表
export function workRoomGroupList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    roomCount: number
    createTime: string
  }>>({
    url: '/workRoomGroup/index',
    params
  })
}

// 删除分组
export function deleteGroup(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/workRoomGroup/destroy',
    data: params
  })
}

// 新建分组
export function createGroup(params: { name: string }) {
  return request.post<{ success: boolean }>({
    url: '/workRoomGroup/store',
    data: params
  })
}

// 更新分组
export function updateGroup(params: { id: number; name: string }) {
  return request.put<{ success: boolean }>({
    url: '/workRoomGroup/update',
    data: params
  })
}

// 客户群列表
export function workRoomList(params: {
  groupId?: number
  name?: string
  employeeId?: number
} & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    avatar: string
    owner: { id: number; name: string; avatar: string }
    memberCount: number
    createTime: string
  }>>({
    url: '/workRoom/index',
    params
  })
}

// 同步群
export function synList() {
  return request.put<{ success: boolean; message?: string }>({
    url: '/workRoom/syn'
  })
}

// 批量修改群
export function batchUpdate(params: { ids: number[]; groupId?: number; employeeIds?: number[] }) {
  return request.put<{ success: boolean }>({
    url: '/workRoom/batchUpdate',
    data: params
  })
}

// 客户群成员
export function workContactRoom(params: { roomId: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    avatar: string
    phone: string
    joinTime: string
    type: string
  }>>({
    url: '/workContactRoom/index',
    params
  })
}

// 部门成员列表
export function workDepartmentList(params: { departmentId: number }) {
  return request.get<{ id: number; name: string; avatar: string; position: string }[]>({
    url: '/workEmployeeDepartment/memberIndex',
    params
  })
}

// 部门列表
export function departmentList() {
  return request.get<{ id: number; name: string; parentId: number }[]>({
    url: '/workDepartment/index'
  })
}

// 统计分页数据
export function statisticsIndex(params: {
  roomId?: number
  startDate?: string
  endDate?: string
} & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    roomId: number
    roomName: string
    date: string
    addMemberNum: number
    quitMemberNum: number
    totalMemberNum: number
  }>>({
    url: '/workRoom/statisticsIndex',
    params
  })
}

// 折线图
export function statistics(params: { roomId?: number; startDate: string; endDate: string }) {
  return request.get<{
    dates: string[]
    addMembers: number[]
    quitMembers: number[]
    totalMembers: number[]
  }>({
    url: '/workRoom/statistics',
    params
  })
}

// 自动拉群列表
export function workRoomAutoPullList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    room: { id: number; name: string }
    tags: { id: number; name: string }[]
    employee: { id: number; name: string }
    createTime: string
  }>>({
    url: '/workRoomAutoPull/index',
    params
  })
}

// 更新
export function autoPullUpdate(params: {
  id: number
  name: string
  roomId: number
  tagIds: number[]
  employeeId: number
}) {
  return request.put<{ success: boolean }>({
    url: '/workRoomAutoPull/update',
    data: params
  })
}

// 创建
export function autoPullCreate(params: {
  name: string
  roomId: number
  tagIds: number[]
  employeeId: number
}) {
  return request.post<{ success: boolean }>({
    url: '/workRoomAutoPull/store',
    data: params
  })
}

// 移动
export function autoPullMove(params: { id: number; targetRoomId: number }) {
  return request.put<{ success: boolean }>({
    url: '/workRoomAutoPull/move',
    data: params
  })
}

// 详情
export function autoPullShow(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    room: { id: number; name: string }
    tags: { id: number; name: string }[]
    employee: { id: number; name: string }
  }>({
    url: '/workRoomAutoPull/show',
    params
  })
}

// 标签分组
export function workContactTagGroup() {
  return request.get<{ id: number; name: string; tags: { id: number; name: string; color: string }[] }[]>({
    url: '/workContactTagGroup/index'
  })
}

// 新建标签
export function addWorkContactTag(params: { groupId: number; name: string; color?: string }) {
  return request.post<{ success: boolean }>({
    url: '/workContactTag/store',
    data: params
  })
}

// 标签列表
export function tagList() {
  return request.get<{ id: number; name: string; color: string; groupId: number }[]>({
    url: '/workContactTag/allTag'
  })
}

// 群聊
export function roomList() {
  return request.get<{ id: number; name: string; memberCount: number }[]>({
    url: '/workRoom/roomIndex'
  })
}

// 选择群聊
export function optCroup(params: { tagIds?: number[] }) {
  return request.get<{ id: number; name: string }[]>({
    url: '/roomTagPull/roomList',
    params
  })
}

// 标签建群获取列表
export function tagGetList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    room: { id: number; name: string }
    tags: { id: number; name: string }[]
    employee: { id: number; name: string }
    contactCount: number
    status: string
    createTime: string
  }>>({
    url: '/roomTagPull/index',
    params
  })
}

// 创建标签建群邀请
export function addGroup(params: {
  name: string
  roomId: number
  tagIds: number[]
  employeeId: number
  contactIds?: number[]
}) {
  return request.post<{ success: boolean }>({
    url: '/roomTagPull/store',
    data: params
  })
}

// 标签建群详情
export function labelShow(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    room: { id: number; name: string }
    tags: { id: number; name: string }[]
    employee: { id: number; name: string }
    contactCount: number
    status: string
  }>({
    url: '/roomTagPull/show',
    params
  })
}

// 删除标签建群
export function delRoomTag(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/roomTagPull/destroy',
    data: params
  })
}

// 标签建群提醒发送
export function remindRoomTag(params: { id: number }) {
  return request.get<{ success: boolean }>({
    url: '/roomTagPull/remindSend',
    params
  })
}

// 标签建群筛选客户
export function chooseContactRoomTag(params: { id: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    avatar: string
    phone: string
    tags: string[]
  }>>({
    url: '/roomTagPull/chooseContact',
    params
  })
}

// 标签建群过滤客户
export function chooseFilterContact(params: { id: number; contactIds: number[] }) {
  return request.post<{ success: boolean }>({
    url: '/roomTagPull/filterContact',
    data: params
  })
}

export function labelContactShow(params: { id: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    avatar: string
    phone: string
    joinTime: string
  }>>({
    url: '/roomTagPull/showContact',
    params
  })
}

// 群信息
export function showRoomApi(params: { roomId: number }) {
  return request.get<{
    id: number
    name: string
    avatar: string
    owner: { id: number; name: string }
    memberCount: number
    createTime: string
  }>({
    url: '/contactMessageBatchSend/showRoom',
    params
  })
}