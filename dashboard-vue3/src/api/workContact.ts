import request from '@/utils/http'

// 客户-客户详情-用户画像
export function getUserPortrait(params: { contactId: number }) {
  return request.get<{
    fields: {
      id: number
      name: string
      value: string
      type: string
    }[]
  }>({
    url: '/contactFieldPivot/index',
    params
  })
}

// 客户-客户详情-编辑用户画像
export function editUserPortrait(params: { contactId: number; fields: { id: number; value: string }[] }) {
  return request.put<{ success: boolean }>({
    url: '/contactFieldPivot/update',
    data: params
  })
}

// 客户-客户详情-批量打标签
export function addTag(params: { contactIds: number[]; tagIds: number[] }) {
  return request.post<{ success: boolean }>({
    url: '/workContact/batchLabeling',
    data: params
  })
}

// 客户-客户详情-查看客户详情基本信息
export function getWorkContactInfo(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    avatar: string
    phone: string
    email: string
    gender: string
    source: string
    tags: { id: number; name: string }[]
    employees: { id: number; name: string; avatar: string }[]
    createTime: string
    updateTime: string
  }>({
    url: '/workContact/show',
    params
  })
}

// 客户-客户详情-编辑客户详情基本信息
export function editWorkContactInfo(params: {
  id: number
  name?: string
  phone?: string
  email?: string
  gender?: string
}) {
  return request.put<{ success: boolean }>({
    url: '/workContact/update',
    data: params
  })
}

// 客户-客户列表
export function workContactList(params: {
  name?: string
  phone?: string
  employeeId?: number
  tagIds?: number[]
  startTime?: string
  endTime?: string
} & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    avatar: string
    phone: string
    gender: string
    source: string
    tags: string[]
    employee: { id: number; name: string; avatar: string }
    createTime: string
  }>>({
    url: '/workContact/index',
    params
  })
}

// 所有标签
export function allTag() {
  return request.get<{ id: number; name: string; color: string }[]>({
    url: '/workContactTag/allTag'
  })
}

// 群聊列表下拉框
export function groupChatList() {
  return request.get<{ id: number; name: string; memberCount: number }[]>({
    url: '/workRoom/roomIndex'
  })
}

// 客户来源下拉框
export function customersSource() {
  return request.get<{ value: string; label: string }[]>({
    url: '/workContact/source'
  })
}

// 客户 - 客户列表筛选 -- 用户画像
export function UserPortraitList(params: { fieldId?: number; fieldValue?: string } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    avatar: string
    phone: string
  }>>({
    url: '/contactField/portrait',
    params
  })
}

// 客户 - 同步客户
export function synContact() {
  return request.put<{ success: boolean; message?: string }>({
    url: '/workContact/synContact'
  })
}

// 客户 - 互动轨迹
export function track(params: { contactId: number }) {
  return request.get<{
    records: {
      id: number
      type: string
      content: string
      time: string
      employee: { name: string; avatar: string }
    }[]
  }>({
    url: '/workContact/track',
    params
  })
}