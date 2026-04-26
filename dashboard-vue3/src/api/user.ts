import request from '@/utils/http'

// 子账户管理列表
export function subManagementList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.SystemManage.UserList>({
    url: '/user/index',
    params
  })
}

// 创建子账户
export function addSubManagement(params: Api.SystemManage.UserCreateParams) {
  return request.post<{ success: boolean }>({
    url: '/user/store',
    data: params
  })
}

// 修改子账户
export function editSubManagement(params: Api.SystemManage.UserCreateParams) {
  return request.put<{ success: boolean }>({
    url: '/user/update',
    data: params
  })
}

// 更改子账户状态
export function editStatus(params: { id: number; status: Api.Common.EnableStatus }) {
  return request.put<{ success: boolean }>({
    url: '/user/statusUpdate',
    data: params
  })
}

// 获取子账户详情
export function getSubManagement(params: { id: number }) {
  return request.get<Api.SystemManage.UserListItem>({
    url: '/user/show',
    params
  })
}

// 子账户启用
export function changeStatus(params: { id: number; status: Api.Common.EnableStatus }) {
  return request.put<{ success: boolean }>({
    url: '/user/statusUpdate',
    data: params
  })
}

// 根据手机号匹配成员部门
export function selectByPhone(params: { phone: string }) {
  return request.get<{ id: number; name: string }>({
    url: '/workDepartment/selectByPhone',
    params
  })
}

// 角色列表
export function selectRole() {
  return request.get<Api.SystemManage.RoleListItem[]>({
    url: '/role/select'
  })
}

// 重置密码
export function passwordResetApi(params: { id: number }) {
  return request.put<{ success: boolean }>({
    url: '/user/passwordReset',
    data: params
  })
}