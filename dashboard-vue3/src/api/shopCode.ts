import request from '@/utils/http'

// 经纬度查询
export function location(params: { address: string }) {
  return request.get<{
    longitude: number
    latitude: number
  }>({
    url: '/shopCode/location',
    params
  })
}

// 关键词列表
export function addressKeyWordList(params: { keyword: string }) {
  return request.get<{
    suggestions: {
      title: string
      address: string
      longitude: number
      latitude: number
    }[]
  }>({
    url: '/shopCode/addressKeyWordList',
    params
  })
}

// 新增
export function storeApi(params: {
  name: string
  type: string
  shopName: string
  shopAddress: string
  longitude: number
  latitude: number
  employeeId: number
  tagIds?: number[]
  welcomeMsg?: string
}) {
  return request.post<{ success: boolean }>({
    url: '/shopCode/store',
    data: params
  })
}

// 修改
export function updateApi(params: {
  id: number
  name: string
  type: string
  shopName: string
  shopAddress: string
  longitude: number
  latitude: number
  employeeId: number
  tagIds?: number[]
  welcomeMsg?: string
}) {
  return request.put<{ success: boolean }>({
    url: '/shopCode/update',
    data: params
  })
}

// 删除
export function destroyApi(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/shopCode/destroy',
    data: params
  })
}

// 详情
export function infoApi(params: { id: number }) {
  return request.get<{
    id: number
    name: string
    type: string
    shopName: string
    shopAddress: string
    longitude: number
    latitude: number
    employee: { id: number; name: string }
    tags: { id: number; name: string }[]
    welcomeMsg: string
    status: string
  }>({
    url: '/shopCode/info',
    params
  })
}

// 开启状态
export function statusApi(params: { id: number; status: Api.Common.EnableStatus }) {
  return request.put<{ success: boolean }>({
    url: '/shopCode/status',
    data: params
  })
}

// 获取列表
export function indexApi(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    type: string
    shopName: string
    shopAddress: string
    contactCount: number
    status: string
    createTime: string
  }>>({
    url: '/shopCode/index',
    params
  })
}

// 拉群活码
export function workRoomIndexApi(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    room: { id: number; name: string }
    employee: { id: number; name: string }
    contactCount: number
    status: string
    createTime: string
  }>>({
    url: '/workRoomAutoPull/index',
    params
  })
}

// 列表查询-城市
export function searchCityApi(params: { keyword: string }) {
  return request.get<{
    cities: { id: string; name: string }[]
  }>({
    url: '/shopCode/searchCity',
    params
  })
}

// 门店地址-关键词列表
export function addressKeyWordListApi(params: { keyword: string }) {
  return request.get<{
    suggestions: {
      title: string
      address: string
      longitude: number
      latitude: number
    }[]
  }>({
    url: '/shopCode/addressKeyWordList',
    params
  })
}

// 门店地址-经纬度查询
export function locationApi(params: { address: string }) {
  return request.get<{
    longitude: number
    latitude: number
  }>({
    url: '/shopCode/location',
    params
  })
}

// 分享
export function shareApi(params: { id: number }) {
  return request.get<{
    shareUrl: string
    qrCode: string
  }>({
    url: '/shopCode/share',
    params
  })
}

// 设置详情
export function pageInfoApi() {
  return request.get<{
    logo: string
    shopName: string
    description: string
    backgroundImage: string
  }>({
    url: '/shopCode/pageInfo'
  })
}

// 设置页面
export function pageSetApi(params: {
  logo: string
  shopName: string
  description: string
  backgroundImage: string
}) {
  return request.post<{ success: boolean }>({
    url: '/shopCode/pageSet',
    data: params
  })
}

// 数据分析-数据总览
export function showApi(params: { id: number }) {
  return request.get<{
    contactCount: number
    todayContactCount: number
    weekContactCount: number
  }>({
    url: '/shopCode/show',
    params
  })
}

// 数据分析-客户数据
export function showContactApi(params: { id: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    contact: { id: number; name: string; avatar: string; phone: string }
    employee: { id: number; name: string }
    createTime: string
  }>>({
    url: '/shopCode/showContact',
    params
  })
}

// 数据分析-门店数据
export function showShopApi(params: { id: number } & Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    shopName: string
    shopAddress: string
    contactCount: number
  }>>({
    url: '/shopCode/showShop',
    params
  })
}

// 成员下拉框
export function department() {
  return request.get<{ id: number; name: string; parentId: number }[]>({
    url: '/workDepartment/index'
  })
}

export function updateEmployeeApi(params: { id: number; employeeId: number }) {
  return request.post<{ success: boolean }>({
    url: '/shopCode/updateEmployee',
    data: params
  })
}

// 修改门店活码
export function updateQrcodeApi(params: {
  id: number
  name: string
  employeeId: number
  tagIds?: number[]
}) {
  return request.post<{ success: boolean }>({
    url: '/shopCode/updateQrcode',
    data: params
  })
}

// 批量打标签
export function batchContactTagsApi(params: { shopCodeId: number; contactIds: number[]; tagIds: number[] }) {
  return request.put<{ success: boolean }>({
    url: '/shopCode/batchContactTags',
    data: params
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