// 自动标签模块
import request from '@/utils/http'

// 新增消息
export function storeApi(params: {
  name: string
  type: 'keyword' | 'room' | 'time'
  keywordRule?: {
    keywords: string[]
    tagIds: number[]
  }
  roomRule?: {
    roomIds: number[]
    tagIds: number[]
  }
  timeRule?: {
    timePeriods: { startTime: string; endTime: string; tagIds: number[] }[]
  }
}) {
  return request.post<{ success: boolean }>({
    url: '/autoTag/store',
    data: params
  })
}

// 获取列表
export function indexApi(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    type: string
    status: string
    tagCount: number
    contactCount: number
    createTime: string
  }>>({
    url: '/autoTag/index',
    params
  })
}

// 删除
export function destroyApi(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/autoTag/destroy',
    data: params
  })
}

// 规则状态
export function onOffApi(params: { id: number; status: Api.Common.EnableStatus }) {
  return request.put<{ success: boolean }>({
    url: '/autoTag/onOff',
    data: params
  })
}

// 详情
export function showApi(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    type: string
    status: string
    rules: any
    tagCount: number
    contactCount: number
    createTime: string
  }>({
    url: '/autoTag/show',
    params
  })
}

// 详情-关键字 客户数据
export function showContactKeyWordApi(params: { id: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    contact: { id: number; name: string; avatar: string; phone: string }
    keyword: string
    tagNames: string[]
    createTime: string
  }>>({
    url: '/autoTag/showContactKeyWord',
    params
  })
}

// 关键字打标签-定时任务
export function KeyWordTagApi(params: { id: number }) {
  return request.get<{ success: boolean; message?: string }>({
    url: '/Task/AutoTag/KeyWordTag',
    params
  })
}

// 详情-客户入群-打标签客户数据
export function showContactRoomApi(params: { id: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    contact: { id: number; name: string; avatar: string; phone: string }
    room: { id: number; name: string }
    tagNames: string[]
    createTime: string
  }>>({
    url: '/autoTag/showContactRoom',
    params
  })
}

// 详情 - 分时段打标签客户数据
export function showContactTimeApi(params: { id: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    contact: { id: number; name: string; avatar: string; phone: string }
    timePeriod: string
    tagNames: string[]
    createTime: string
  }>>({
    url: '/autoTag/showContactTime',
    params
  })
}

// 会话对象列表
export function toUsersList(params: {
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