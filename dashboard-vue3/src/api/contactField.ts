import request from '@/utils/http'

// 客户资料卡列表
export function contactFieldList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    type: string
    options: string[]
    sort: number
    status: Api.Common.EnableStatus
    createTime: string
  }>>({
    url: '/contactField/index',
    params
  })
}

// 客户资料卡-新增属性
export function addContactField(params: {
  name: string
  type: string
  options?: string[]
  sort?: number
  status: Api.Common.EnableStatus
}) {
  return request.post<{ success: boolean }>({
    url: '/contactField/store',
    data: params
  })
}

// 客户资料卡-编辑属性
export function editContactField(params: {
  id: number
  name: string
  type: string
  options?: string[]
  sort?: number
  status: Api.Common.EnableStatus
}) {
  return request.put<{ success: boolean }>({
    url: '/contactField/update',
    data: params
  })
}

// 客户资料卡-删除属性
export function delContactField(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/contactField/destroy',
    data: params
  })
}

// 状态修改
export function statusUpdate(params: { id: number; status: Api.Common.EnableStatus }) {
  return request.put<{ success: boolean }>({
    url: '/contactField/statusUpdate',
    data: params
  })
}

// 批量修改
export function batchUpdate(params: { ids: number[]; updates: Record<string, any> }) {
  return request.put<{ success: boolean }>({
    url: '/contactField/batchUpdate',
    data: params
  })
}