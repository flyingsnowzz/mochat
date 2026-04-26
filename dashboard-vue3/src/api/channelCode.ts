import request from '@/utils/http'

// 渠道码分组列表
export function channelCodeGroup(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    createTime: string
  }>>({
    url: '/channelCodeGroup/index',
    params
  })
}

// 渠道码分组修改
export function channelCodeGroupUpdate(params: { id: number; name: string }) {
  return request.put<{ success: boolean }>({
    url: '/channelCodeGroup/update',
    data: params
  })
}

// 渠道码分组新增
export function channelCodeGroupAdd(params: { name: string }) {
  return request.post<{ success: boolean }>({
    url: '/channelCodeGroup/store',
    data: params
  })
}

// 渠道码分组移动
export function channelCodeGroupMove(params: { id: number; targetGroupId: number }) {
  return request.put<{ success: boolean }>({
    url: '/channelCodeGroup/move',
    data: params
  })
}

// 渠道码新增
export function channelCodeAdd(params: {
  group_id: number
  name: string
  type: string
  auto_add_friend: number
  tags: string[]
  drainage_employee: { id: number; name: string }[]
  welcome_message: { content: string; media_id?: string }[]
}) {
  return request.post<{ success: boolean }>({
    url: '/channelCode/store',
    data: params
  })
}

// 渠道码编辑
export function channelCodeUpdate(params: {
  id: number
  group_id: number
  name: string
  type: string
  auto_add_friend: number
  tags: string[]
  drainage_employee: { id: number; name: string }[]
  welcome_message: { content: string; media_id?: string }[]
}) {
  return request.put<{ success: boolean }>({
    url: '/channelCode/update',
    data: params
  })
}

// 渠道码列表
export function channelCodeList(params: { group_id?: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    group_id: number
    group_name: string
    name: string
    type: string
    auto_add_friend: number
    tags: string[]
    contact_num: number
    create_time: string
  }>>({
    url: '/channelCode/index',
    params
  })
}

// 渠道码客户
export function channelCodeContact(params: { channelCodeId: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    phone: string
    avatar: string
    create_time: string
  }>>({
    url: '/channelCode/contact',
    params
  })
}

// 渠道码详情
export function channelCodeDetail(params: { id: number }) {
  return request.get<{
    id: number
    group_id: number
    name: string
    type: string
    auto_add_friend: number
    tags: string[]
    drainage_employee: { id: number; name: string }[]
    welcome_message: { content: string; media_id?: string }[]
  }>({
    url: '/channelCode/show',
    params
  })
}

// 统计分页数据
export function statisticsIndex(params: { channelCodeId?: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    channel_code_id: number
    date: string
    new_contact_num: number
    total_contact_num: number
  }>>({
    url: '/channelCode/statisticsIndex',
    params
  })
}

// 统计折线图
export function statistics(params: { channelCodeId: number; startDate: string; endDate: string }) {
  return request.get<{
    dates: string[]
    newContacts: number[]
    totalContacts: number[]
  }>({
    url: '/channelCode/statistics',
    params
  })
}

// 成员下拉框
export function department() {
  return request.get<{ id: number; name: string; parentId: number }[]>({
    url: '/workDepartment/index'
  })
}

// 部门成员
export function member(params: { departmentId: number }) {
  return request.get<{ id: number; name: string; avatar: string; position: string }[]>({
    url: '/workDepartment/memberIndex',
    params
  })
}