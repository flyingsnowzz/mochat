import request from '@/utils/http'

// 新增消息
export function addMessage(params: {
  name: string
  roomIds: number[]
  employeeIds: number[]
  content: string
  sendTime?: string
}) {
  return request.post<{ success: boolean }>({
    url: '/roomMessageBatchSend/store',
    data: params
  })
}

// 查询列表
export function index(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    content: string
    roomNum: number
    employeeNum: number
    sendStatus: string
    sendTime: string
    createTime: string
  }>>({
    url: '/roomMessageBatchSend/index',
    params
  })
}

// 消息提醒
export function remind(params: { id: number }) {
  return request.get<{ success: boolean }>({
    url: '/roomMessageBatchSend/remind',
    params
  })
}

// 消息详情-基础信息、数据统计
export function show(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    content: string
    roomNum: number
    employeeNum: number
    receiveNum: number
    unreadNum: number
    sendStatus: string
    sendTime: string
  }>({
    url: '/roomMessageBatchSend/show',
    params
  })
}

// 删除
export function destroyApi(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/roomMessageBatchSend/destroy',
    data: params
  })
}

// 消息详情-客户群接收详情
export function roomOwnerSendIndex(params: { id: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    room: { id: number; name: string }
    sendStatus: string
    sendTime: string
  }>>({
    url: '/roomMessageBatchSend/roomOwnerSendIndex',
    params
  })
}

// 消息详情-客户群接收详情
export function roomReceiveIndex(params: { id: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    room: { id: number; name: string }
    receiveNum: number
    unreadNum: number
  }>>({
    url: '/roomMessageBatchSend/roomReceiveIndex',
    params
  })
}

// 成员下拉框
export function department() {
  return request.get<{ id: number; name: string; parentId: number }[]>({
    url: '/workDepartment/index'
  })
}