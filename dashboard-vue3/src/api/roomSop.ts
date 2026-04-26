import request from '@/utils/http'

// 获取列表
export function getList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    type: string
    status: string
    roomCount: number
    employeeCount: number
    createTime: string
  }>>({
    url: '/roomSop/index',
    params
  })
}

// 创建
export function add(params: {
  name: string
  type: string
  roomIds: number[]
  rules: {
    action: string
    content: string
    delay: number
  }[]
}) {
  return request.post<{ success: boolean }>({
    url: '/roomSop/store',
    data: params
  })
}

// 添加群
export function setRoom(params: { id: number; roomIds: number[] }) {
  return request.put<{ success: boolean }>({
    url: '/roomSop/setRoom',
    data: params
  })
}

// 开关规则
export function switchState(params: { id: number; status: Api.Common.EnableStatus }) {
  return request.put<{ success: boolean }>({
    url: '/roomSop/state',
    data: params
  })
}

// 获取详情
export function getInfo(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    type: string
    status: string
    rooms: { id: number; name: string }[]
    rules: {
      action: string
      content: string
      delay: number
    }[]
    createTime: string
  }>({
    url: '/roomSop/info',
    params
  })
}

// 删除规则
export function del(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/roomSop/destroy',
    data: params
  })
}

// 编辑规则
export function edit(params: {
  id: number
  name: string
  type: string
  roomIds: number[]
  rules: {
    action: string
    content: string
    delay: number
  }[]
}) {
  return request.put<{ success: boolean }>({
    url: '/roomSop/update',
    data: params
  })
}