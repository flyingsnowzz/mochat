import request from '@/utils/http'

// 素材库-素材库列表
export function materialLibraryList(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<any>>({
    url: '/medium/index',
    params
  })
}

// 素材库-素材库查看
export function getMaterialLibrary(params: { id: number }) {
  return request.get<any>({
    url: '/medium/show',
    params
  })
}

// 素材库-素材库添加
export function addMaterialLibrary(params: any) {
  return request.post<{ success: boolean }>({
    url: '/medium/store',
    data: params
  })
}

// 素材库-素材库删除
export function delMaterialLibrary(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/medium/destroy',
    data: params
  })
}

// 素材库-移动分组
export function moveGroup(params: { ids: number[]; groupId: number }) {
  return request.put<{ success: boolean }>({
    url: '/medium/groupUpdate',
    data: params
  })
}

// 素材库-修改素材库
export function editMaterialLibrary(params: any) {
  return request.put<{ success: boolean }>({
    url: '/medium/update',
    data: params
  })
}

// 素材库分组-素材库分组列表
export function mediumGroup(params: Api.Common.CommonSearchParams) {
  return request.get<Api.Common.PaginatedResponse<any>>({
    url: '/mediumGroup/index',
    params
  })
}

// 素材库分组-素材库分组添加
export function addMediumGroup(params: { name: string }) {
  return request.post<{ success: boolean }>({
    url: '/mediumGroup/store',
    data: params
  })
}

// 素材库分组-素材库分组修改
export function editMediumGroup(params: { id: number; name: string }) {
  return request.put<{ success: boolean }>({
    url: '/mediumGroup/update',
    data: params
  })
}

// 素材库分组-素材库分组删除
export function delMediumGroup(params: { id: number }) {
  return request.del<{ success: boolean }>({
    url: '/mediumGroup/destroy',
    data: params
  })
}

// 上传
export function upLoad(params: FormData) {
  return request.post<any>({
    url: '/common/upload',
    data: params
  })
}
