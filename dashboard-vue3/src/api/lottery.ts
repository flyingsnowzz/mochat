import request from '@/utils/http'

// 获取列表
export function getList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    type: string
    status: string
    participateCount: number
    winnerCount: number
    startTime: string
    endTime: string
    createTime: string
  }>>({
    url: '/lottery/index',
    params
  })
}

// 新建活动
export function addActivity(params: {
  name: string
  type: string
  content: string
  startTime: string
  endTime: string
  prizes: {
    name: string
    type: string
    value: string
    count: number
    probability?: number
  }[]
  rules: string
}) {
  return request.post<{ success: boolean }>({
    url: '/lottery/store',
    data: params
  })
}

// 客户数据详情
export function dataDetails(params: { id: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    contact: { id: number; name: string; avatar: string; phone: string }
    prize: { name: string; type: string; value: string }
    participateTime: string
    status: string
  }>>({
    url: '/lottery/showContact',
    params
  })
}

// 详情
export function getDetails(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    type: string
    content: string
    status: string
    participateCount: number
    winnerCount: number
    startTime: string
    endTime: string
    prizes: {
      name: string
      type: string
      value: string
      count: number
      givenCount: number
      probability?: number
    }[]
    rules: string
    createTime: string
  }>({
    url: '/lottery/show',
    params
  })
}

// 删除
export function del(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/lottery/destroy',
    data: params
  })
}

// 分享
export function share(params: { id: number }) {
  return request.get<{
    shareUrl: string
    qrCode: string
  }>({
    url: '/lottery/share',
    params
  })
}

// 修改
export function update(params: {
  id: number
  name: string
  type: string
  content: string
  startTime: string
  endTime: string
  prizes: {
    name: string
    type: string
    value: string
    count: number
    probability?: number
  }[]
  rules: string
}) {
  return request.put<{ success: boolean }>({
    url: '/lottery/update',
    data: params
  })
}

// 修改-详情
export function modify(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    type: string
    content: string
    startTime: string
    endTime: string
    prizes: {
      id: number
      name: string
      type: string
      value: string
      count: number
      probability?: number
    }[]
    rules: string
  }>({
    url: '/lottery/info',
    params
  })
}

// 核销
export function writeOffApi(params: { id: number; contactId: number }) {
  return request.get<{ success: boolean; message?: string }>({
    url: '/lottery/writeOff',
    params
  })
}

// 批量打标签
export function batchContactTagsApi(params: { lotteryId: number; contactIds: number[]; tagIds: number[] }) {
  return request.put<{ success: boolean }>({
    url: '/lottery/batchContactTags',
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