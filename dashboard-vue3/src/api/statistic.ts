import request from '@/utils/http'

// 客户统计
export function contactInfo(params: { startDate?: string; endDate?: string }) {
  return request.get<{
    totalContact: number
    todayAddContact: number
    todayLossContact: number
    weekAddContact: number
    weekLossContact: number
    monthAddContact: number
    monthLossContact: number
  }>({
    url: '/statistic/index',
    params
  })
}

// 联系客户数据
export function employeesInfo(params: { employeeId?: number; startDate?: string; endDate?: string }) {
  return request.get<{
    totalContact: number
    todayContact: number
    weekContact: number
    monthContact: number
  }>({
    url: '/statistic/employees',
    params
  })
}

// 趋势图/列表数据
export function employeesTrendInfo(params: {
  employeeId?: number
  startDate: string
  endDate: string
}) {
  return request.get<{
    dates: string[]
    values: {
      date: string
      contactCount: number
      chatCount: number
    }[]
  }>({
    url: '/statistic/employeesTrend',
    params
  })
}

// 排行榜前十数据
export function topList(params: {
  type: 'contact' | 'chat'
  startDate: string
  endDate: string
  limit?: number
}) {
  return request.get<{
    rankings: {
      rank: number
      employee: { id: number; name: string; avatar: string }
      count: number
    }[]
  }>({
    url: '/statistic/topList',
    params
  })
}

// 按员工查看
export function getEmployeeCounts(params: {
  departmentId?: number
  startDate?: string
  endDate?: string
} & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    employee: { id: number; name: string; avatar: string }
    contactCount: number
    chatCount: number
    newContactCount: number
    lostContactCount: number
  }>>({
    url: '/statistic/employeeCounts',
    params
  })
}

// 成员下拉框
export function department() {
  return request.get<{ id: number; name: string; parentId: number }[]>({
    url: '/workDepartment/index'
  })
}