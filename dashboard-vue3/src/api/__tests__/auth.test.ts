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
import { login, getInfo, corpSelect, corpBind, logout, permissionByUser } from '../auth'

// 导入模拟的 http 模块
import http from '@/utils/http'

describe('Auth API', () => {
  beforeEach(() => {
    // 清除所有模拟调用
    vi.clearAllMocks()
  })

  it('should export login function', () => {
    expect(typeof login).toBe('function')
  })

  it('should call post method with correct parameters when login', async () => {
    // 模拟响应
    const mockResponse = { token: 'test-token', refreshToken: 'test-refresh-token' }
    http.post.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { userName: 'admin', password: '123456' }
    const result = await login(params)

    // 验证调用
    expect(http.post).toHaveBeenCalledWith({
      url: '/user/auth',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export getInfo function', () => {
    expect(typeof getInfo).toBe('function')
  })

  it('should call get method with correct parameters when getInfo', async () => {
    // 模拟响应
    const mockResponse = { id: 1, name: 'Admin', phone: '13800138000' }
    http.get.mockResolvedValue(mockResponse)

    // 调用函数
    const result = await getInfo()

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/user/loginShow'
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export corpSelect function', () => {
    expect(typeof corpSelect).toBe('function')
  })

  it('should call get method with correct parameters when corpSelect', async () => {
    // 模拟响应
    const mockResponse = [{ id: 1, name: 'Test Corp', wxCorpid: 'wx123456' }]
    http.get.mockResolvedValue(mockResponse)

    // 调用函数
    const result = await corpSelect()

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/corp/select'
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export corpBind function', () => {
    expect(typeof corpBind).toBe('function')
  })

  it('should call post method with correct parameters when corpBind', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.post.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { corpId: 1 }
    const result = await corpBind(params)

    // 验证调用
    expect(http.post).toHaveBeenCalledWith({
      url: '/corp/bind',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export logout function', () => {
    expect(typeof logout).toBe('function')
  })

  it('should call put method with correct parameters when logout', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.put.mockResolvedValue(mockResponse)

    // 调用函数
    const result = await logout()

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: '/user/logout'
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export permissionByUser function', () => {
    expect(typeof permissionByUser).toBe('function')
  })

  it('should call get method with correct parameters when permissionByUser', async () => {
    // 模拟响应
    const mockResponse = { permissions: ['dashboard', 'user'] }
    http.get.mockResolvedValue(mockResponse)

    // 调用函数
    const result = await permissionByUser()

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/role/permissionByUser'
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })
})
