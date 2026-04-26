import request from '@/utils/http'

// 角色列表
export function roleList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.SystemManage.RoleList>({
    url: '/role/index',
    params
  })
}

// 角色状态更改
export function statusUpdate(params: { id: number; status: Api.Common.EnableStatus }) {
  return request.put<{ success: boolean }>({
    url: '/role/statusUpdate',
    data: params
  })
}

// 添加角色
export function roleStore(params: Api.SystemManage.RoleCreateParams) {
  return request.post<{ success: boolean }>({
    url: '/role/store',
    data: params
  })
}

// 修改角色
export function roleUpdate(params: Api.SystemManage.RoleCreateParams) {
  return request.put<{ success: boolean }>({
    url: '/role/update',
    data: params
  })
}

// 删除角色
export function roleDelete(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/role/destroy',
    data: params
  })
}

// 角色详情
export function roleDetail(params: { id: number }) {
  return request.get<Api.SystemManage.RoleListItem>({
    url: '/role/show',
    params
  })
}

// 查看人员
export function roleShowEmployee(params: { id: number }) {
  return request.get<{ records: any[]; total: number }>({
    url: '/role/showEmployee',
    params
  })
}

// 查看权限
export function rolePermission(params: { id: number }) {
  return request.get<{ permissions: string[] }>({
    url: '/role/permissionShow',
    params
  })
}

// 添加角色权限
export function rolePermissionStore(params: Api.SystemManage.RolePermissionParams) {
  return request.post<{ success: boolean }>({
    url: '/role/permissionStore',
    data: params
  })
}