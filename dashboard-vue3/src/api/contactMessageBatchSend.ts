import request from '@/utils/http'

// 消息列表
export function indexApi(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    content: string
    contactNum: number
    employeeNum: number
    sendStatus: string
    sendTime: string
    createTime: string
  }>>({
    url: '/contactMessageBatchSend/index',
    params
  })
}

// 新建
export function storeApi(params: {
  name: string
  employeeIds: number[]
  contactFilter: {
    tagIds?: number[]
    departmentIds?: number[]
  }
  content: string
  sendTime?: string
}) {
  return request.post<{ success: boolean }>({
    url: '/contactMessageBatchSend/store',
    data: params
  })
}

// 基础信息、数据统计
export function showApi(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    content: string
    contactNum: number
    employeeNum: number
    receiveNum: number
    unreadNum: number
    sendStatus: string
    sendTime: string
  }>({
    url: '/contactMessageBatchSend/show',
    params
  })
}

// 删除
export function destroyApi(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/contactMessageBatchSend/destroy',
    data: params
  })
}

// 预览
export function messageShowApi(params: { id: number }) {
  return request.get<{
    content: string
    mediaType: string
    mediaUrl?: string
  }>({
    url: '/contactMessageBatchSend/messageShow',
    params
  })
}

// 提醒发送
export function remindApi(params: { id: number }) {
  return request.post<{ success: boolean }>({
    url: '/contactMessageBatchSend/remind',
    data: params
  })
}

// 客户详情
export function contactReceiveIndexApi(params: { id: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    contact: { id: number; name: string; avatar: string }
    receiveStatus: string
    receiveTime: string
  }>>({
    url: '/contactMessageBatchSend/contactReceiveIndex',
    params
  })
}

// 成员详情
export function employeeSendIndexApi(params: { id: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    employee: { id: number; name: string; avatar: string }
    sendCount: number
    sendTime: string
  }>>({
    url: '/contactMessageBatchSend/employeeSendIndex',
    params
  })
}