import request from '@/utils/http'

// 客户列表
export function contactList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    name: string
    phone: string
    state: string
    importId: number
    allotStatus: string
    employee: { id: number; name: string }
    createTime: string
  }>>({
    url: '/contactBatchAdd/index',
    params
  })
}

// 分配客户
export function allot(params: { contactIds: number[]; employeeId: number }) {
  return request.post<{ success: boolean }>({
    url: '/contactBatchAdd/allot',
    data: params
  })
}

// 数据统计
export function dataStatistic() {
  return request.get<{
    totalCount: number
    allotCount: number
    pendingCount: number
  }>({
    url: '/contactBatchAdd/dataStatistic'
  })
}

// 删除客户
export function contactDel(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/contactBatchAdd/destroy',
    data: params
  })
}

// 导入记录
export function importData(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<{
    id: number
    fileName: string
    totalCount: number
    successCount: number
    failCount: number
    status: string
    createTime: string
  }>>({
    url: '/contactBatchAdd/importIndex',
    params
  })
}

// 删除导入记录
export function importDestroyApi(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/contactBatchAdd/importDestroy',
    data: params
  })
}

// 导入客户
export function importContact(params: { file: File }) {
  const formData = new FormData()
  formData.append('file', params.file)
  return request.post<{ success: boolean; importId: number }>({
    url: '/contactBatchAdd/importStore',
    data: formData
  })
}

// 获取设置
export function getSetting() {
  return request.get<{
    autoAllot: boolean
    allotRule: string
    duplicateHandle: string
  }>({
    url: '/contactBatchAdd/settingEdit'
  })
}

// 修改设置
export function updateSetting(params: {
  autoAllot: boolean
  allotRule: string
  duplicateHandle: string
}) {
  return request.post<{ success: boolean }>({
    url: '/contactBatchAdd/settingUpdate',
    data: params
  })
}

// 成员下拉框
export function department() {
  return request.get<{ id: number; name: string; parentId: number }[]>({
    url: '/workDepartment/index'
  })
}

// 提醒
export function remindApi(params: { id: number }) {
  return request.get<{ success: boolean; message?: string }>({
    url: '/contactBatchAdd/remind',
    params
  })
}