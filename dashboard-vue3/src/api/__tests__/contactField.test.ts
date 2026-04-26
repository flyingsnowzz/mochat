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
import { contactFieldList, addContactField, editContactField, delContactField, statusUpdate, batchUpdate } from '../contactField'

// 导入模拟的 http 模块
import http from '@/utils/http'

describe('ContactField API', () => {
  beforeEach(() => {
    // 清除所有模拟调用
    vi.clearAllMocks()
  })

  it('should export contactFieldList function', () => {
    expect(typeof contactFieldList).toBe('function')
  })

  it('should call get method with correct parameters when contactFieldList', async () => {
    // 模拟响应
    const mockResponse = { list: [], total: 0, page: 1, perPage: 10 }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { page: 1, perPage: 10 }
    const result = await contactFieldList(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/contactField/index',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export addContactField function', () => {
    expect(typeof addContactField).toBe('function')
  })

  it('should call post method with correct parameters when addContactField', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.post.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { name: 'Test Field', type: 'text', status: 1 }
    const result = await addContactField(params)

    // 验证调用
    expect(http.post).toHaveBeenCalledWith({
      url: '/contactField/store',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export editContactField function', () => {
    expect(typeof editContactField).toBe('function')
  })

  it('should call put method with correct parameters when editContactField', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.put.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1, name: 'Test Field', type: 'text', status: 1 }
    const result = await editContactField(params)

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: '/contactField/update',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export delContactField function', () => {
    expect(typeof delContactField).toBe('function')
  })

  it('should call del method with correct parameters when delContactField', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.del.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1 }
    const result = await delContactField(params)

    // 验证调用
    expect(http.del).toHaveBeenCalledWith({
      url: '/contactField/destroy',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export statusUpdate function', () => {
    expect(typeof statusUpdate).toBe('function')
  })

  it('should call put method with correct parameters when statusUpdate', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.put.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1, status: 1 }
    const result = await statusUpdate(params)

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: '/contactField/statusUpdate',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export batchUpdate function', () => {
    expect(typeof batchUpdate).toBe('function')
  })

  it('should call put method with correct parameters when batchUpdate', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.put.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { ids: [1, 2], updates: { status: 1 } }
    const result = await batchUpdate(params)

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: '/contactField/batchUpdate',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })
})
