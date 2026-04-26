import request from '@/utils/http'

// 获取列表
export function getList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    type: string
    status: string
    employeeCount: number
    contactCount: number
    createTime: string
  }>>({
    url: '/contactSop/index',
    params
  })
}

// 创建
export function add(params: {
  name: string
  type: string
  employeeIds: number[]
  rules: {
    action: string
    content: string
    delay: number
  }[]
}) {
  return request.post<{ success: boolean }>({
    url: '/contactSop/store',
    data: params
  })
}

// 添加客服
export function setEmployee(params: { id: number; employeeIds: number[] }) {
  return request.put<{ success: boolean }>({
    url: '/contactSop/setEmployee',
    data: params
  })
}

// 开关规则
export function switchState(params: { id: number; status: Api.Common.EnableStatus }) {
  return request.put<{ success: boolean }>({
    url: '/contactSop/state',
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
    employees: { id: number; name: string; avatar: string }[]
    rules: {
      action: string
      content: string
      delay: number
    }[]
    contactCount: number
    createTime: string
  }>({
    url: '/contactSop/info',
    params
  })
}

// 删除规则
export function del(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/contactSop/destroy',
    data: params
  })
}

// 编辑规则
export function edit(params: {
  id: number
  name: string
  type: string
  employeeIds: number[]
  rules: {
    action: string
    content: string
    delay: number
  }[]
}) {
  return request.put<{ success: boolean }>({
    url: '/contactSop/update',
    data: params
  })
}