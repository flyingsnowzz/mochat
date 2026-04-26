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
import { roleList, statusUpdate, roleStore, roleUpdate, roleDelete, roleDetail, roleShowEmployee, rolePermission, rolePermissionStore } from '../role'

// 导入模拟的 http 模块
import http from '@/utils/http'

describe('Role API', () => {
  beforeEach(() => {
    // 清除所有模拟调用
    vi.clearAllMocks()
  })

  it('should export roleList function', () => {
    expect(typeof roleList).toBe('function')
  })

  it('should call get method with correct parameters when roleList', async () => {
    // 模拟响应
    const mockResponse = { list: [], total: 0, page: 1, perPage: 10 }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { page: 1, perPage: 10 }
    const result = await roleList(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/role/index',
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
      url: '/role/statusUpdate',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export roleStore function', () => {
    expect(typeof roleStore).toBe('function')
  })

  it('should call post method with correct parameters when roleStore', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.post.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { name: 'Test Role', status: 1 }
    const result = await roleStore(params)

    // 验证调用
    expect(http.post).toHaveBeenCalledWith({
      url: '/role/store',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export roleUpdate function', () => {
    expect(typeof roleUpdate).toBe('function')
  })

  it('should call put method with correct parameters when roleUpdate', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.put.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1, name: 'Test Role', status: 1 }
    const result = await roleUpdate(params)

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: '/role/update',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export roleDelete function', () => {
    expect(typeof roleDelete).toBe('function')
  })

  it('should call del method with correct parameters when roleDelete', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.del.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1 }
    const result = await roleDelete(params)

    // 验证调用
    expect(http.del).toHaveBeenCalledWith({
      url: '/role/destroy',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export roleDetail function', () => {
    expect(typeof roleDetail).toBe('function')
  })

  it('should call get method with correct parameters when roleDetail', async () => {
    // 模拟响应
    const mockResponse = { id: 1, name: 'Test Role', status: 1 }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1 }
    const result = await roleDetail(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/role/show',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export roleShowEmployee function', () => {
    expect(typeof roleShowEmployee).toBe('function')
  })

  it('should call get method with correct parameters when roleShowEmployee', async () => {
    // 模拟响应
    const mockResponse = { records: [], total: 0 }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1 }
    const result = await roleShowEmployee(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/role/showEmployee',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export rolePermission function', () => {
    expect(typeof rolePermission).toBe('function')
  })

  it('should call get method with correct parameters when rolePermission', async () => {
    // 模拟响应
    const mockResponse = { permissions: ['dashboard', 'user'] }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1 }
    const result = await rolePermission(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/role/permissionShow',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export rolePermissionStore function', () => {
    expect(typeof rolePermissionStore).toBe('function')
  })

  it('should call post method with correct parameters when rolePermissionStore', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.post.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { roleId: 1, permissions: ['dashboard', 'user'] }
    const result = await rolePermissionStore(params)

    // 验证调用
    expect(http.post).toHaveBeenCalledWith({
      url: '/role/permissionStore',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })
})
