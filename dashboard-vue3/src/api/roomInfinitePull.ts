import request from '@/utils/http'

// 获取列表
export function getList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    room: { id: number; name: string }
    employee: { id: number; name: string }
    tagCount: number
    contactCount: number
    status: string
    createTime: string
  }>>({
    url: '/roomInfinitePull/index',
    params
  })
}

// 获取详情
export function getInfo(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    roomId: number
    employeeId: number
    tagIds: number[]
    content: string
    status: string
  }>({
    url: '/roomInfinitePull/info',
    params
  })
}

// 修改
export function update(params: {
  id: number
  name: string
  roomId: number
  employeeId: number
  tagIds: number[]
  content: string
}) {
  return request.put<{ success: boolean }>({
    url: '/roomInfinitePull/update',
    data: params
  })
}

// 删除
export function del(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/roomInfinitePull/destroy',
    data: params
  })
}

// 新增
export function add(params: {
  name: string
  roomId: number
  employeeId: number
  tagIds: number[]
  content: string
}) {
  return request.post<{ success: boolean }>({
    url: '/roomInfinitePull/store',
    data: params
  })
}