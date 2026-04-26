import request from '@/utils/http'

// 菜单列表
export function menuList() {
  return request.get<Api.SystemManage.MenuListItem[]>({
    url: 'menu/index'
  })
}

// 菜单下拉
export function menuSelect() {
  return request.get<Api.SystemManage.MenuListItem[]>({
    url: 'menu/select'
  })
}

// 菜单详情
export function menuDetail(params: { id: number }) {
  return request.get<Api.SystemManage.MenuListItem>({
    url: 'menu/show',
    params
  })
}

// 菜单禁用启用
export function statusUpdate(params: { id: number; status: Api.Common.EnableStatus }) {
  return request.put<{ success: boolean }>({
    url: 'menu/statusUpdate',
    data: params
  })
}

// 菜单移除
export function destroy(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: 'menu/destroy',
    data: params
  })
}

// 添加菜单
export function menuStore(params: Api.SystemManage.MenuCreateParams) {
  return request.post<{ success: boolean }>({
    url: 'menu/store',
    data: params
  })
}

// 修改菜单
export function menuUpdate(params: Api.SystemManage.MenuCreateParams) {
  return request.put<{ success: boolean }>({
    url: 'menu/update',
    data: params
  })
}

// 已使用图标
export function iconUsed() {
  return request.get<{ icons: string[] }>({
    url: 'menu/iconIndex'
  })
}