import request from '@/utils/http'

// 构建PC端授权链接
export function componentloginpageApi(params: { redirectUri: string }) {
  return request.get<{
    authUrl: string
  }>({
    url: '/officialAccount/getPreAuthUrl',
    params
  })
}

// 公众号列表
export function publicApi(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    appId: string
    avatar: string
    status: Api.Common.EnableStatus
    createTime: string
  }>>({
    url: '/officialAccount/index',
    params
  })
}