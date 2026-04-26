import request from '@/utils/http'

// 企业信息添加
export function addInformation(params: any) {
  return request.post<{ success: boolean }>({
    url: '/workMessageConfig/corpStore',
    data: params
  })
}

// 企业成员查看
export function getEnterMembers(params: { id: number }) {
  return request.get<any>({
    url: '/workMessageConfig/corpShow',
    params
  })
}

// 列表
export function wechatAuthList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<any>>({
    url: '/workMessageConfig/corpIndex',
    params
  })
}
