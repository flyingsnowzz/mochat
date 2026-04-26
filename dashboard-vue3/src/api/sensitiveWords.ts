import request from '@/utils/http'

// 敏感词词库列表
export function sensitiveWordsList(params: { groupId?: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    word: string
    groupId: number
    groupName: string
    status: Api.Common.EnableStatus
    createTime: string
  }>>({
    url: '/sensitiveWord/index',
    params
  })
}

// 删除敏感词
export function delSensitiveWords(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/sensitiveWord/destroy',
    data: params
  })
}

// 敏感词状态修改
export function sensitiveWordsStatus(params: { id: number; status: Api.Common.EnableStatus }) {
  return request.put<{ success: boolean }>({
    url: '/sensitiveWord/statusUpdate',
    data: params
  })
}

// 敏感词移动
export function sensitiveWordsMove(params: { id: number; targetGroupId: number }) {
  return request.put<{ success: boolean }>({
    url: '/sensitiveWord/move',
    data: params
  })
}

// 新增敏感词
export function sensitiveWordsAdd(params: { groupId: number; word: string }) {
  return request.post<{ success: boolean }>({
    url: '/sensitiveWord/store',
    data: params
  })
}

// 敏感词修改分组
export function sensitiveWordsGroupUp(params: { id: number; name: string }) {
  return request.put<{ success: boolean }>({
    url: '/sensitiveWordGroup/update',
    data: params
  })
}

// 敏感词分组列表
export function sensitiveWordsGroupList() {
  return request.get<{
    id: number
    name: string
    wordCount: number
  }[]>({
    url: '/sensitiveWordGroup/select'
  })
}

// 分组添加
export function sensitiveWordsGroupAdd(params: { name: string }) {
  return request.post<{ success: boolean }>({
    url: '/sensitiveWordGroup/store',
    data: params
  })
}