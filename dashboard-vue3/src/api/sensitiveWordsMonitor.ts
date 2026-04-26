import request from '@/utils/http'

// 敏感词监控列表
export function sensitiveWordsMonitor(params: {
  employeeId?: number
  startTime?: string
  endTime?: string
} & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    employee: { id: number; name: string; avatar: string }
    contact: { id: number; name: string; avatar: string }
    content: string
    sensitiveWords: string[]
    createTime: string
  }>>({
    url: '/sensitiveWordsMonitor/index',
    params
  })
}

// 敏感词监控对话详情
export function dialogueDetails(params: { id: number }) {
  return request.get<{
    id: number
    employee: { id: number; name: string; avatar: string }
    contact: { id: number; name: string; avatar: string }
    messages: {
      id: number
      content: string
      msgType: string
      sender: { id: number; name: string; avatar: string }
      isSensitive: boolean
      sensitiveWords?: string[]
      createTime: string
    }[]
    createTime: string
  }>({
    url: '/sensitiveWordsMonitor/show',
    params
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