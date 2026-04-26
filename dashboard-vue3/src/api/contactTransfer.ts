// 在职转接模块
import request from '@/utils/http'

// 在职转接列表
export function infoApi(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    fromEmployee: { id: number; name: string; avatar: string }
    toEmployee: { id: number; name: string; avatar: string }
    contact: { id: number; name: string; avatar: string }
    status: string
    createTime: string
  }>>({
    url: '/contactTransfer/info',
    params
  })
}

// 同步离职待分配客户列表
export function saveUnassignedListApi() {
  return request.get<{ success: boolean; message?: string }>({
    url: '/contactTransfer/saveUnassignedList'
  })
}

// 离职待分配客户列表
export function unassignedListApi(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    contact: { id: number; name: string; avatar: string; phone: string }
   离职员工: { id: number; name: string }
    createTime: string
  }>>({
    url: '/contactTransfer/unassignedList',
    params
  })
}

// 分配 离职/在职客户
export function indexApi(params: {
  contactIds: number[]
  type: 'transfer' | 'dimission'
  toEmployeeId: number
}) {
  return request.post<{ success: boolean }>({
    url: '/contactTransfer/index',
    data: params
  })
}

// 离职待分配群列表
export function roomApi(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    room: { id: number; name: string; avatar: string }
   离职员工: { id: number; name: string }
    memberCount: number
    createTime: string
  }>>({
    url: '/contactTransfer/room',
    params
  })
}

// 分配离职客服群
export function postRoomApi(params: { roomIds: number[]; toEmployeeId: number }) {
  return request.post<{ success: boolean }>({
    url: '/contactTransfer/room',
    data: params
  })
}

// 分配记录
export function logApi(params: { contactId?: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    contact: { id: number; name: string; avatar: string }
    fromEmployee: { id: number; name: string }
    toEmployee: { id: number; name: string }
    type: string
    createTime: string
  }>>({
    url: '/contactTransfer/log',
    params
  })
}

// 会话内容
export function messageList(params: { contactId: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    content: string
    sender: { id: number; name: string; avatar: string }
    createTime: string
  }>>({
    url: '/workMessage/index',
    params
  })
}

// 成员下拉框
export function department() {
  return request.get<{ id: number; name: string; parentId: number }[]>({
    url: '/workDepartment/index'
  })
}