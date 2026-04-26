// 雷达模块
import request from '@/utils/http'

// 新增消息
export function storeApi(params: {
  name: string
  content: string
  articleUrl?: string
  mediaId?: string
  tagIds?: number[]
  employeeIds?: number[]
}) {
  return request.post<{ success: boolean }>({
    url: '/radar/store',
    data: params
  })
}

// 修改消息
export function updateApi(params: {
  id: number
  name: string
  content: string
  articleUrl?: string
  mediaId?: string
  tagIds?: number[]
  employeeIds?: number[]
}) {
  return request.put<{ success: boolean }>({
    url: '/radar/update',
    data: params
  })
}

// 获取列表
export function indexApi(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    content: string
    articleUrl: string
    tagCount: number
    clickCount: number
    contactCount: number
    createTime: string
  }>>({
    url: '/radar/index',
    params
  })
}

// 删除
export function destroyApi(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/radar/destroy',
    data: params
  })
}

// 详情-修改
export function infoApi(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    content: string
    articleUrl: string
    mediaId: string
    tags: { id: number; name: string }[]
    employees: { id: number; name: string }[]
  }>({
    url: '/radar/info',
    params
  })
}

// 添加渠道
export function storeChannelApi(params: {
  radarId: number
  name: string
  channelType: string
}) {
  return request.post<{ success: boolean }>({
    url: '/radar/storeChannel',
    data: params
  })
}

// 生成渠道链接
export function storeChannelLinkApi(params: {
  channelId: number
  linkType: string
}) {
  return request.post<{
    shortLink: string
    qrCode: string
  }>({
    url: '/radar/storeChannelLink',
    data: params
  })
}

// 渠道列表
export function indexChannelApi(params: { radarId: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    channelType: string
    linkCount: number
    clickCount: number
    contactCount: number
    createTime: string
  }>>({
    url: '/radar/indexChannel',
    params
  })
}

// 渠道链接列表
export function indexChannelLinkApi(params: { channelId: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    shortLink: string
    linkType: string
    clickCount: number
    contactCount: number
    createTime: string
  }>>({
    url: '/radar/indexChannelLink',
    params
  })
}

// 详情
export function showApi(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    content: string
    articleUrl: string
    clickCount: number
    contactCount: number
    tags: { id: number; name: string }[]
    employees: { id: number; name: string }[]
    createTime: string
  }>({
    url: '/radar/show',
    params
  })
}

// 详情-客户数据
export function showContactApi(params: { id: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    contact: { id: number; name: string; avatar: string; phone: string }
    clickTime: string
    source: string
  }>>({
    url: '/radar/showContact',
    params
  })
}

// 详情-渠道数据
export function showChannelApi(params: { id: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    channel: { id: number; name: string }
    clickCount: number
    contactCount: number
  }>>({
    url: '/radar/showChannel',
    params
  })
}

// 成员下拉框
export function department() {
  return request.get<{ id: number; name: string; parentId: number }[]>({
    url: '/workDepartment/index'
  })
}

// 生成雷达文章
export function radarArticleApi(params: { url: string }) {
  return request.get<{
    title: string
    cover: string
    description: string
  }>({
    url: '/radar/radarArticle',
    params
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