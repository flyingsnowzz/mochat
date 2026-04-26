import request from '@/utils/http'

// 企业信息添加
export function corpStore(params: {
  corpId: number
  corpName: string
  wxCorpId: string
  agentId: string
  corpSecret: string
}) {
  return request.post<{ success: boolean }>({
    url: '/workMessageConfig/corpStore',
    data: params
  })
}

// 微信后台配置
export function stepUpdate(params: {
  corpId: number
 WxCustomerApi: string
  wxCustomerSecret: string
  wxCustomerToken: string
  wxCustomerAesKey: string
}) {
  return request.put<{ success: boolean }>({
    url: '/workMessageConfig/stepUpdate',
    data: params
  })
}

// 微信后台配置查看
export function stepCreate(params: { corpId: number }) {
  return request.get<{
    corpId: number
    corpName: string
    wxCorpId: string
    agentId: string
    WxCustomerApi: string
    wxCustomerSecret: string
    wxCustomerToken: string
    wxCustomerAesKey: string
    createTime: string
  }>({
    url: '/workMessageConfig/stepCreate',
    params
  })
}

// 企业成员
export function workEmployee(params: { corpId: number }) {
  return request.get<{ id: number; name: string; avatar: string }[]>({
    url: '/workMessage/fromUsers',
    params
  })
}

// 会话对象列表
export function toUsersList(params: {
  corpId: number
  fromUserId?: number
  keyword?: string
} & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    type: string
    name: string
    avatar: string
    lastMessage: string
    lastMessageTime: string
  }>>({
    url: '/workMessage/toUsers',
    params
  })
}

// 会话内容
export function messageList(params: {
  corpId: number
  fromUserId: number
  toOpenId: string
} & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    content: string
    msgType: string
    fromUser: { id: number; name: string; avatar: string }
    createTime: string
  }>>({
    url: '/workMessage/index',
    params
  })
}