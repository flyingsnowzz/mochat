import request from '@/utils/http'

// 应用信息
export function chatTool() {
  return request.get<{
    appId: string
    appSecret: string
    token: string
    encodingAesKey: string
  }>({
    url: '/chatTool/config'
  })
}

// 添加应用
export function addChatTool(params: {
  name: string
  appId: string
  appSecret: string
  token: string
  encodingAesKey: string
}) {
  return request.post<{ success: boolean }>({
    url: '/agent/store',
    data: params
  })
}