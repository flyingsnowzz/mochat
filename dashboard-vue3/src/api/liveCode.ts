import request from '@/utils/http'

export function update(params: {
  id: number
  name: string
  content: string
  roomIds?: number[]
  employeeIds?: number[]
}) {
  return request.put<{ success: boolean }>({
    url: '/roomWelcome/update',
    data: params
  })
}

export function getList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    type: string
    status: string
    contactCount: number
    createTime: string
  }>>({
    url: '/clockIn/index',
    params
  })
}