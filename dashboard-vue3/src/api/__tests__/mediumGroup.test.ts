import { describe, it, expect, vi, beforeEach } from 'vitest'

// 模拟所有可能的依赖
vi.mock('@/utils/http', () => ({
  default: {
    get: vi.fn(),
    post: vi.fn(),
    put: vi.fn(),
    del: vi.fn(),
    request: vi.fn()
  }
}))

// 模拟其他可能的依赖
vi.mock('@/config', () => ({
  default: {
    API_BASE_URL: 'http://localhost:8000'
  }
}))

// 导入 API 函数
import { materialLibraryList, getMaterialLibrary, addMaterialLibrary, delMaterialLibrary, moveGroup, editMaterialLibrary, mediumGroup, addMediumGroup, editMediumGroup, delMediumGroup, upLoad } from '../mediumGroup'

// 导入模拟的 http 模块
import http from '@/utils/http'

describe('MediumGroup API', () => {
  beforeEach(() => {
    // 清除所有模拟调用
    vi.clearAllMocks()
  })

  it('should export materialLibraryList function', () => {
    expect(typeof materialLibraryList).toBe('function')
  })

  it('should call get method with correct parameters when materialLibraryList', async () => {
    // 模拟响应
    const mockResponse = { list: [], total: 0, page: 1, perPage: 10 }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { page: 1, perPage: 10 }
    const result = await materialLibraryList(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/medium/index',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export getMaterialLibrary function', () => {
    expect(typeof getMaterialLibrary).toBe('function')
  })

  it('should call get method with correct parameters when getMaterialLibrary', async () => {
    // 模拟响应
    const mockResponse = { id: 1, name: 'Test Material', url: 'test.jpg' }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1 }
    const result = await getMaterialLibrary(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/medium/show',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export addMaterialLibrary function', () => {
    expect(typeof addMaterialLibrary).toBe('function')
  })

  it('should call post method with correct parameters when addMaterialLibrary', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.post.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { name: 'Test Material', url: 'test.jpg', groupId: 1 }
    const result = await addMaterialLibrary(params)

    // 验证调用
    expect(http.post).toHaveBeenCalledWith({
      url: '/medium/store',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export delMaterialLibrary function', () => {
    expect(typeof delMaterialLibrary).toBe('function')
  })

  it('should call del method with correct parameters when delMaterialLibrary', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.del.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1 }
    const result = await delMaterialLibrary(params)

    // 验证调用
    expect(http.del).toHaveBeenCalledWith({
      url: '/medium/destroy',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export moveGroup function', () => {
    expect(typeof moveGroup).toBe('function')
  })

  it('should call put method with correct parameters when moveGroup', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.put.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { ids: [1, 2], groupId: 3 }
    const result = await moveGroup(params)

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: '/medium/groupUpdate',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export editMaterialLibrary function', () => {
    expect(typeof editMaterialLibrary).toBe('function')
  })

  it('should call put method with correct parameters when editMaterialLibrary', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.put.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1, name: 'Test Material', url: 'test.jpg', groupId: 1 }
    const result = await editMaterialLibrary(params)

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: '/medium/update',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export mediumGroup function', () => {
    expect(typeof mediumGroup).toBe('function')
  })

  it('should call get method with correct parameters when mediumGroup', async () => {
    // 模拟响应
    const mockResponse = { list: [], total: 0, page: 1, perPage: 10 }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { page: 1, perPage: 10 }
    const result = await mediumGroup(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/mediumGroup/index',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export addMediumGroup function', () => {
    expect(typeof addMediumGroup).toBe('function')
  })

  it('should call post method with correct parameters when addMediumGroup', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.post.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { name: 'Test Group' }
    const result = await addMediumGroup(params)

    // 验证调用
    expect(http.post).toHaveBeenCalledWith({
      url: '/mediumGroup/store',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export editMediumGroup function', () => {
    expect(typeof editMediumGroup).toBe('function')
  })

  it('should call put method with correct parameters when editMediumGroup', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.put.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1, name: 'Test Group' }
    const result = await editMediumGroup(params)

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: '/mediumGroup/update',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export delMediumGroup function', () => {
    expect(typeof delMediumGroup).toBe('function')
  })

  it('should call del method with correct parameters when delMediumGroup', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.del.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1 }
    const result = await delMediumGroup(params)

    // 验证调用
    expect(http.del).toHaveBeenCalledWith({
      url: '/mediumGroup/destroy',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export upLoad function', () => {
    expect(typeof upLoad).toBe('function')
  })

  it('should call post method with correct parameters when upLoad', async () => {
    // 模拟响应
    const mockResponse = { success: true, url: 'test.jpg' }
    http.post.mockResolvedValue(mockResponse)

    // 测试参数
    const params = new FormData()
    params.append('file', new Blob(['test'], { type: 'text/plain' }))
    const result = await upLoad(params)

    // 验证调用
    expect(http.post).toHaveBeenCalledWith({
      url: '/common/upload',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })
})
