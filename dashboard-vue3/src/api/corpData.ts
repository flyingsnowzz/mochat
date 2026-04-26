import request from '@/utils/http'

// 后台首页-首页统计
export function corpData() {
  return request.get<{
    todayAddContact: number
    todayAddRoom: number
    todayAddIntoRoom: number
    todayLossContact: number
    todayQuitRoom: number
    totalContact: number
    totalRoom: number
  }>({
    url: '/corpData/index'
  })
}

// 后台首页-折线图
export function lineChat(params: { startDate: string; endDate: string }) {
  return request.get<{
    dates: string[]
    addContacts: number[]
    addRooms: number[]
    addIntoRooms: number[]
    lossContacts: number[]
    quitRooms: number[]
  }>({
    url: '/corpData/lineChat',
    params
  })
}