import request from '@/utils/http'

// 组织架构列表
export function departmentList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    parentId: number
    level: number
    employeeCount: number
  }>>({
    url: '/workDepartment/pageIndex',
    params
  })
}

// 组织架构成员列表
export function showEmployee(params: { departmentId: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    phone: string
    avatar: string
    department: string
    position: string
  }>>({
    url: '/workDepartment/showEmployee',
    params
  })
}