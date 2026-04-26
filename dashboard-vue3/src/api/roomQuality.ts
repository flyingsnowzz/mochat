// 群质量模块
import request from '@/utils/http'

// 新增
export function storeApi(params: {
  name: string
  type: string
  rules: {
    condition: string
    value: string
    action: string
  }[]
}) {
  return request.post<{ success: boolean }>({
    url: '/roomQuality/store',
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
    roomCount: number
    contactCount: number
    createTime: string
  }>>({
    url: '/roomQuality/index',
    params
  })
}

// 规则状态
export function statusApi(params: { id: number; status: Api.Common.EnableStatus }) {
  return request.put<{ success: boolean }>({
    url: '/roomQuality/status',
    data: params
  })
}

// 详情
export function infoApi(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    type: string
    status: string
    rules: {
      condition: string
      value: string
      action: string
    }[]
    createTime: string
  }>({
    url: '/roomQuality/info',
    params
  })
}

// 修改
export function updateApi(params: {
  id: number
  name: string
  type: string
  rules: {
    condition: string
    value: string
    action: string
  }[]
}) {
  return request.put<{ success: boolean }>({
    url: '/roomQuality/update',
    data: params
  })
}

// 详情-客户数据
export function showContactApi(params: { id: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    room: { id: number; name: string }
    contact: { id: number; name: string; avatar: string }
    status: string
    createTime: string
  }>>({
    url: '/roomQuality/showContact',
    params
  })
}

// 删除
export function destroyApi(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/roomQuality/destroy',
    data: params
  })
}

// 获取群聊消息
export function indexgroupApi(params: { roomId: number } & Api.Common.CommonSearchParams) {
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

// 客户数据-详情
export function contactDetailApi(params: { id: number; contactId: number }) {
  return request.get<{
    id: number
    contact: { id: number; name: string; avatar: string; phone: string }
    room: { id: number; name: string }
    joinTime: string
    messages: {
      content: string
      msgType: string
      createTime: string
    }[]
  }>({
    url: '/roomQuality/contactDetail',
    params
  })
}