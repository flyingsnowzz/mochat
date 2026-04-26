import request from '@/utils/http'

// 新增
export function addActivity(params: {
  name: string
  content: string
  startTime: string
  endTime: string
  roomIds: number[]
  employeeIds: number[]
  checkInRules: {
    type: string
    value: string
  }[]
}) {
  return request.post<{ success: boolean }>({
    url: '/roomClockIn/store',
    data: params
  })
}

// 修改
export function updateActivity(params: {
  id: number
  name: string
  content: string
  startTime: string
  endTime: string
  roomIds: number[]
  employeeIds: number[]
  checkInRules: {
    type: string
    value: string
  }[]
}) {
  return request.put<{ success: boolean }>({
    url: '/roomClockIn/update',
    data: params
  })
}

// 删除
export function delList(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/roomClockIn/destroy',
    data: params
  })
}

// 获取列表
export function getList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    content: string
    startTime: string
    endTime: string
    roomCount: number
    employeeCount: number
    checkInCount: number
    status: string
    createTime: string
  }>>({
    url: '/roomClockIn/index',
    params
  })
}

// 详情
export function detailsList(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    content: string
    startTime: string
    endTime: string
    rooms: { id: number; name: string }[]
    employees: { id: number; name: string; avatar: string }[]
    checkInRules: {
      type: string
      value: string
    }[]
    checkInCount: number
    status: string
    createTime: string
  }>({
    url: '/roomClockIn/show',
    params
  })
}

// 详情-客户数据详情
export function clientDetails(params: { id: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    contact: { id: number; name: string; avatar: string; phone: string }
    checkInDays: number
    checkInTimes: number
    lastCheckInTime: string
  }>>({
    url: '/roomClockIn/showContact',
    params
  })
}

// 批量打标签
export function batchLabel(params: { clockInId: number; contactIds: number[]; tagIds: number[] }) {
  return request.put<{ success: boolean }>({
    url: '/roomClockIn/batchContactTags',
    data: params
  })
}

// 详情-修改
export function modifyDetails(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    content: string
    startTime: string
    endTime: string
    roomIds: number[]
    employeeIds: number[]
    checkInRules: {
      type: string
      value: string
    }[]
  }>({
    url: '/roomClockIn/info',
    params
  })
}

// 创建客户标签
export function setclientTags(params: { groupId: number; name: string; color?: string }) {
  return request.post<{ success: boolean }>({
    url: '/workContactTag/store',
    data: params
  })
}

// 客户标签列表
export function clientTagsReceive(params: { groupId?: number } & Api.Common.CommonSearchParams) {
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

// 详细打卡天数
export function dayDetail(params: { clockInId: number; contactId: number }) {
  return request.get<{
    records: {
      date: string
      checkInTime: string
      status: string
    }[]
  }>({
    url: '/roomClockIn/dayDetail',
    params
  })
}

// 新建标签组
export function newTagGroup(params: { name: string }) {
  return request.post<{ success: boolean }>({
    url: '/workContactTagGroup/store',
    data: params
  })
}

// 批量打标签
export function batchContactTagsApi(params: { clockInId: number; contactIds: number[]; tagIds: number[] }) {
  return request.put<{ success: boolean }>({
    url: '/roomClockIn/batchContactTags',
    data: params
  })
}

// 公众号列表
export function publicIndexApi(params: Api.Common.CommonSearchParams) {
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