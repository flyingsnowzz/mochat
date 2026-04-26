import request from '@/utils/http'

// 客户标签列表
export function contactTagList(params: { groupId?: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    color: string
    groupId: number
    groupName: string
    contactNum: number
    createTime: string
  }>>({
    url: '/workContactTag/index',
    params
  })
}

// 删除标签
export function delContactTag(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/workContactTag/destroy',
    data: params
  })
}

// 新增标签
export function addContactTag(params: { groupId: number; name: string; color?: string }) {
  return request.post<{ success: boolean }>({
    url: '/workContactTag/store',
    data: params
  })
}

// 标签详情
export function contactTagDetail(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    color: string
    groupId: number
    groupName: string
    contactNum: number
  }>({
    url: '/workContactTag/detail',
    params
  })
}

// 移动标签
export function moveContactTag(params: { id: number; targetGroupId: number }) {
  return request.put<{ success: boolean }>({
    url: '/workContactTag/move',
    data: params
  })
}

// 编辑标签
export function editContactTag(params: { id: number; name: string; color?: string }) {
  return request.put<{ success: boolean }>({
    url: '/workContactTag/update',
    data: params
  })
}

// 分组列表
export function getContactTagGroup() {
  return request.get<{
    id: number
    name: string
    tagNum: number
    contactNum: number
  }[]>({
    url: '/workContactTagGroup/index'
  })
}

// 新增分组
export function addContactTagGroup(params: { name: string }) {
  return request.post<{ success: boolean }>({
    url: '/workContactTagGroup/store',
    data: params
  })
}

// 编辑分组
export function editContactTagGroup(params: { id: number; name: string }) {
  return request.put<{ success: boolean }>({
    url: '/workContactTagGroup/update',
    data: params
  })
}

// 删除分组
export function delContactTagGroup(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/workContactTagGroup/destroy',
    data: params
  })
}

// 分组详情
export function contactTagGroupDetail(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    tags: { id: number; name: string; color: string }[]
  }>({
    url: '/workContactTagGroup/detail',
    params
  })
}

// 同步标签
export function syncTag() {
  return request.put<{ success: boolean; message?: string }>({
    url: '/workContactTag/synContactTag'
  })
}