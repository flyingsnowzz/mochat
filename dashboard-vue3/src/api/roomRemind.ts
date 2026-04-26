// 群提醒模块
import request from '@/utils/http'

// 新增
export function storeApi(params: {
  name: string
  roomIds: number[]
  employeeIds: number[]
  remindRules: {
    type: string
    time: string
    content: string
  }[]
}) {
  return request.post<{ success: boolean }>({
    url: '/roomRemind/store',
    data: params
  })
}

// 详情
export function infoApi(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    rooms: { id: number; name: string }[]
    employees: { id: number; name: string }[]
    remindRules: {
      type: string
      time: string
      content: string
    }[]
    status: string
    createTime: string
  }>({
    url: '/roomRemind/info',
    params
  })
}

// 修改
export function updateApi(params: {
  id: number
  name: string
  roomIds: number[]
  employeeIds: number[]
  remindRules: {
    type: string
    time: string
    content: string
  }[]
}) {
  return request.put<{ success: boolean }>({
    url: '/roomRemind/update',
    data: params
  })
}

// 状态
export function statusApi(params: { id: number }) {
  return request.get<{
    status: string
    lastRunTime: string
    nextRunTime: string
  }>({
    url: '/roomRemind/status',
    params
  })
}

// 获取任务列表
export function indexApi(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    roomCount: number
    employeeCount: number
    status: string
    createTime: string
  }>>({
    url: '/roomRemind/index',
    params
  })
}

// 定时任务
export function roomRemindApi(params: { id: number }) {
  return request.get<{ success: boolean; message?: string }>({
    url: '/task/roomRemind',
    params
  })
}

// 删除
export function destroyApi(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/roomRemind/destroy',
    data: params
  })
}