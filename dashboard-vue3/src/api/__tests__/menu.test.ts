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
import { menuList, menuSelect, menuDetail, statusUpdate, destroy, menuStore, menuUpdate, iconUsed } from '../menu'

// 导入模拟的 http 模块
import http from '@/utils/http'

describe('Menu API', () => {
  beforeEach(() => {
    // 清除所有模拟调用
    vi.clearAllMocks()
  })

  it('should export menuList function', () => {
    expect(typeof menuList).toBe('function')
  })

  it('should call get method with correct parameters when menuList', async () => {
    // 模拟响应
    const mockResponse = [{ id: 1, name: 'Dashboard', path: '/dashboard' }]
    http.get.mockResolvedValue(mockResponse)

    // 调用函数
    const result = await menuList()

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: 'menu/index'
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export menuSelect function', () => {
    expect(typeof menuSelect).toBe('function')
  })

  it('should call get method with correct parameters when menuSelect', async () => {
    // 模拟响应
    const mockResponse = [{ id: 1, name: 'Dashboard', path: '/dashboard' }]
    http.get.mockResolvedValue(mockResponse)

    // 调用函数
    const result = await menuSelect()

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: 'menu/select'
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export menuDetail function', () => {
    expect(typeof menuDetail).toBe('function')
  })

  it('should call get method with correct parameters when menuDetail', async () => {
    // 模拟响应
    const mockResponse = { id: 1, name: 'Dashboard', path: '/dashboard' }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1 }
    const result = await menuDetail(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: 'menu/show',
      params
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
      url: 'menu/statusUpdate',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export destroy function', () => {
    expect(typeof destroy).toBe('function')
  })

  it('should call del method with correct parameters when destroy', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.del.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1 }
    const result = await destroy(params)

    // 验证调用
    expect(http.del).toHaveBeenCalledWith({
      url: 'menu/destroy',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export menuStore function', () => {
    expect(typeof menuStore).toBe('function')
  })

  it('should call post method with correct parameters when menuStore', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.post.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { name: 'Test Menu', path: '/test', parentId: 0 }
    const result = await menuStore(params)

    // 验证调用
    expect(http.post).toHaveBeenCalledWith({
      url: 'menu/store',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export menuUpdate function', () => {
    expect(typeof menuUpdate).toBe('function')
  })

  it('should call put method with correct parameters when menuUpdate', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.put.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1, name: 'Test Menu', path: '/test', parentId: 0 }
    const result = await menuUpdate(params)

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: 'menu/update',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export iconUsed function', () => {
    expect(typeof iconUsed).toBe('function')
  })

  it('should call get method with correct parameters when iconUsed', async () => {
    // 模拟响应
    const mockResponse = { icons: ['home', 'user', 'settings'] }
    http.get.mockResolvedValue(mockResponse)

    // 调用函数
    const result = await iconUsed()

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: 'menu/iconIndex'
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })
})
