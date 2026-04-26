// 群裂变模块
import request from '@/utils/http'

// 新增
export function storeApi(params: {
  name: string
  roomId: number
  employeeId: number
  tagIds: number[]
  content: string
  rewardContent?: string
}) {
  return request.post<{ success: boolean }>({
    url: '/roomFission/store',
    data: params
  })
}

// 详情-修改
export function infoApi(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    roomId: number
    employeeId: number
    tagIds: number[]
    content: string
    rewardContent?: string
    status: string
  }>({
    url: '/roomFission/info',
    params
  })
}

// 修改
export function updateApi(params: {
  id: number
  name: string
  roomId: number
  employeeId: number
  tagIds: number[]
  content: string
  rewardContent?: string
}) {
  return request.put<{ success: boolean }>({
    url: '/roomFission/update',
    data: params
  })
}

// 删除
export function destroyApi(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/roomFission/destroy',
    data: params
  })
}

// 获取列表
export function indexApi(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    room: { id: number; name: string }
    employee: { id: number; name: string }
    contactCount: number
    inviteCount: number
    conversionRate: number
    status: string
    createTime: string
  }>>({
    url: '/roomFission/index',
    params
  })
}

// 邀请
export function inviteApi(params: { id: number; contactIds?: number[] }) {
  return request.post<{ success: boolean; message?: string }>({
    url: '/roomFission/invite',
    data: params
  })
}

// 客服成员
export function department() {
  return request.get<{ id: number; name: string; parentId: number }[]>({
    url: '/workDepartment/index'
  })
}

// 客户标签列表
export function contactTagListApi(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    color: string
    groupId: number
    contactNum: number
  }>>({
    url: '/workContactTag/contactTagList',
    params
  })
}

// 数据总览
export function showApi(params: { id: number }) {
  return request.get<{
    contactCount: number
    inviteCount: number
    acceptCount: number
    writeOffCount: number
    conversionRate: number
  }>({
    url: '/roomFission/show',
    params
  })
}

// 群聊数据
export function showRoomApi(params: { id: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    contact: { id: number; name: string; avatar: string }
    room: { id: number; name: string }
    inviteStatus: string
    acceptStatus: string
    createTime: string
  }>>({
    url: '/roomFission/showRoom',
    params
  })
}

// 客户数据
export function showContactApi(params: { id: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    contact: { id: number; name: string; avatar: string; phone: string }
    inviteStatus: string
    inviteTime: string
    acceptStatus: string
    acceptTime: string
    writeOffStatus: string
    writeOffTime: string
  }>>({
    url: '/roomFission/showContact',
    params
  })
}

// 核销
export function writeOffApi(params: { id: number; contactId: number }) {
  return request.get<{ success: boolean; message?: string }>({
    url: '/roomFission/writeOff',
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
    status: string
  }>>({
    url: '/officialAccount/index',
    params
  })
}