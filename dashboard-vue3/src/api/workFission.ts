import request from '@/utils/http'

// 获取列表
export function getList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    fissionType: string
    room: { id: number; name: string }
    employee: { id: number; name: string }
    contactCount: number
    inviteCount: number
    status: string
    createTime: string
  }>>({
    url: '/workFission/index',
    params
  })
}

// 删除
export function del(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/workFission/destroy',
    data: params
  })
}

// 获取详情
export function getDetails(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    fissionType: string
    room: { id: number; name: string }
    employee: { id: number; name: string }
    contactCount: number
    inviteCount: number
    status: string
    createTime: string
  }>({
    url: '/workFission/show',
    params
  })
}

// 获取修改用详情
export function getInfo(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    fissionType: string
    roomId: number
    employeeId: number
    tagIds: number[]
    content: string
  }>({
    url: '/workFission/info',
    params
  })
}

// 修改
export function update(params: {
  id: number
  name: string
  fissionType: string
  roomId: number
  employeeId: number
  tagIds: number[]
  content: string
}) {
  return request.put<{ success: boolean }>({
    url: '/workFission/update',
    data: params
  })
}

// 新增
export function add(params: {
  name: string
  fissionType: string
  roomId: number
  employeeId: number
  tagIds: number[]
  content: string
}) {
  return request.post<{ success: boolean }>({
    url: '/workFission/store',
    data: params
  })
}

// 获取统计信息
export function getStatistics(params: { id: number }) {
  return request.get<{
    contactCount: number
    inviteCount: number
    acceptCount: number
    conversionRate: number
  }>({
    url: '/workFission/statistics',
    params
  })
}

// 获取客户列表
export function getUserList(params: { id: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    contact: { id: number; name: string; avatar: string; phone: string }
    inviteStatus: string
    inviteTime: string
    acceptStatus: string
    acceptTime: string
  }>>({
    url: '/workFission/inviteData',
    params
  })
}

// 获取客户邀请详情
export function getInviteInfo(params: { id: number }) {
  return request.get<{
    id: number
    contact: { id: number; name: string; avatar: string; phone: string }
    inviteStatus: string
    inviteTime: string
    acceptStatus: string
    acceptTime: string
    timeline: {
      event: string
      time: string
    }[]
  }>({
    url: '/workFission/inviteDetail',
    params
  })
}

// 获取预计发送邀请数量
export function chooseContact(params: { id: number }) {
  return request.get<{
    contactCount: number
    filteredCount: number
  }>({
    url: '/workFission/chooseContact',
    params
  })
}

// 发送邀请信息
export function inviteMsg(params: { id: number; contactIds?: number[] }) {
  return request.post<{ success: boolean; message?: string }>({
    url: '/workFission/invite',
    data: params
  })
}