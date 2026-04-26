import request from '@/utils/http'

// 登录
export function login(params: Api.Auth.LoginParams) {
  return request.post<Api.Auth.LoginResponse>({
    url: '/user/auth',
    data: params
  })
}

// 登录用户详情
export function getInfo() {
  return request.get<Api.Auth.LoginUserInfo>({
    url: '/user/loginShow'
  })
}

// 企业下拉列表
export function corpSelect() {
  return request.get<Api.Auth.CorpSelectItem[]>({
    url: '/corp/select'
  })
}

// 登录用户绑定企业
export function corpBind(params: { corpId: number }) {
  return request.post<{ success: boolean }>({
    url: '/corp/bind',
    data: params
  })
}

// 登出
export function logout() {
  return request.put<{ success: boolean }>({
    url: '/user/logout'
  })
}

// 用户角色
export function permissionByUser() {
  return request.get<{ permissions: string[] }>({
    url: '/role/permissionByUser'
  })
}

// 兼容性别名
export const fetchLogin = login
export const fetchGetUserInfo = getInfo