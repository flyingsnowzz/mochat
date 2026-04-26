import request from '@/utils/http'

// 欢迎语列表
export function greetingList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    content: string
    roomCount: number
    employeeCount: number
    status: string
    createTime: string
  }>>({
    url: '/greeting/index',
    params
  })
}

// 欢迎语删除
export function delGreeting(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/greeting/destroy',
    data: params
  })
}

// 创建欢迎语
export function greetingStore(params: {
  name: string
  content: string
  roomIds?: number[]
  employeeIds?: number[]
  attachments?: {
    type: string
    url: string
  }[]
}) {
  return request.post<{ success: boolean }>({
    url: '/greeting/store',
    data: params
  })
}

// 修改欢迎语
export function upDateGreeting(params: {
  id: number
  name: string
  content: string
  roomIds?: number[]
  employeeIds?: number[]
  attachments?: {
    type: string
    url: string
  }[]
}) {
  return request.put<{ success: boolean }>({
    url: '/greeting/update',
    data: params
  })
}

// 欢迎语详情
export function greetingDetail(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    content: string
    rooms: { id: number; name: string }[]
    employees: { id: number; name: string }[]
    attachments: {
      type: string
      url: string
    }[]
    status: string
    createTime: string
  }>({
    url: '/greeting/show',
    params
  })
}