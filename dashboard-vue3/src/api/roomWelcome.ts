import request from '@/utils/http'

// 获取列表
export function getList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    content: string
    roomCount: number
    employeeCount: number
    status: string
    createTime: string
  }>>({
    url: '/roomWelcome/index',
    params
  })
}

// 删除
export function del(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/roomWelcome/destroy',
    data: params
  })
}

// 获取详情
export function getDetail(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    content: string
    rooms: { id: number; name: string }[]
    employees: { id: number; name: string }[]
    status: string
    createTime: string
  }>({
    url: '/roomWelcome/show',
    params
  })
}

// 修改
export function update(params: {
  id: number
  name: string
  content: string
  roomIds: number[]
  employeeIds: number[]
}) {
  return request.put<{ success: boolean }>({
    url: '/roomWelcome/update',
    data: params
  })
}

// 新增
export function add(params: {
  name: string
  content: string
  roomIds: number[]
  employeeIds: number[]
}) {
  return request.post<{ success: boolean }>({
    url: '/roomWelcome/store',
    data: params
  })
}