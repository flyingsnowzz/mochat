import request from '@/utils/http'

// 企业成员列表
export function enterMembersList(params: {
  name?: string
  phone?: string
  departmentId?: number
  status?: Api.Common.EnableStatus
} & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    avatar: string
    phone: string
    email: string
    department: { id: number; name: string }
    position: string
    status: Api.Common.EnableStatus
    createTime: string
  }>>({
    url: '/workEmployee/index',
    params
  })
}

// 同步企业成员
export function syncEmployee() {
  return request.put<{ success: boolean; message?: string }>({
    url: '/workEmployee/synEmployee'
  })
}

// 同步时间-成员列表搜索条件数据
export function syncTime() {
  return request.get<{
    lastSyncTime: string
    syncStatus: 'pending' | 'success' | 'failed'
  }>({
    url: '/workEmployee/searchCondition'
  })
}