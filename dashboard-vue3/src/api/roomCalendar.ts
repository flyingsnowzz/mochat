import request from '@/utils/http'

// 获取列表
export function getList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    title: string
    content: string
    startTime: string
    endTime: string
    roomCount: number
    employeeCount: number
    createTime: string
  }>>({
    url: '/roomCalendar/index',
    params
  })
}

// 添加群聊
export function addRoom(params: { calendarId: number; roomIds: number[] }) {
  return request.post<{ success: boolean }>({
    url: '/roomCalendar/addRoom',
    data: params
  })
}

// 删除群聊
export function delRoom(params: { calendarId: number; roomId: number }) {
  return request.del<{ success: boolean }>({
    url: '/roomCalendar/destroyRoom',
    data: params
  })
}

// 新增
export function add(params: {
  title: string
  content: string
  startTime: string
  endTime: string
  roomIds: number[]
  employeeIds: number[]
}) {
  return request.post<{ success: boolean }>({
    url: '/roomCalendar/store',
    data: params
  })
}

// 删除
export function del(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/roomCalendar/destroy',
    data: params
  })
}

// 详情
export function getInfo(params: { id: number }) {
  return request.get<{
    id: number
    title: string
    content: string
    startTime: string
    endTime: string
    rooms: { id: number; name: string }[]
    employees: { id: number; name: string; avatar: string }[]
    createTime: string
  }>({
    url: '/roomCalendar/show',
    params
  })
}

// 修改
export function updateApi(params: {
  id: number
  title: string
  content: string
  startTime: string
  endTime: string
  roomIds: number[]
  employeeIds: number[]
}) {
  return request.put<{ success: boolean }>({
    url: '/roomCalendar/update',
    data: params
  })
}