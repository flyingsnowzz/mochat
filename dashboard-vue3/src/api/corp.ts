import request from '@/utils/http'

// 企业列表
export function corpList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{ id: number; name: string; wxCorpid: string; status: string }>>({
    url: '/corp/index',
    params
  })
}

// 创建企业
export function corpStore(params: { name: string; wxCorpid: string; status: Api.Common.EnableStatus }) {
  return request.post<{ success: boolean }>({
    url: '/corp/store',
    data: params
  })
}

// 修改企业
export function corpUpdate(params: { id: number; name: string; wxCorpid: string; status: Api.Common.EnableStatus }) {
  return request.put<{ success: boolean }>({
    url: '/corp/update',
    data: params
  })
}

// 删除企业
export function corpDelete(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/corp/destroy',
    data: params
  })
}

// 企业详情
export function corpDetail(params: { id: number }) {
  return request.get<{ id: number; name: string; wxCorpid: string; status: string }>({
    url: '/corp/show',
    params
  })
}

// 企业状态更改
export function corpStatusUpdate(params: { id: number; status: Api.Common.EnableStatus }) {
  return request.put<{ success: boolean }>({
    url: '/corp/statusUpdate',
    data: params
  })
}