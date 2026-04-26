import request from '@/utils/http'

// 流失客户列表
export function getLossContactList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    avatar: string
    phone: string
    lossTime: string
    lossReason: string
    employee: { id: number; name: string; avatar: string }
  }>>({
    url: '/workContact/lossContact',
    params
  })
}