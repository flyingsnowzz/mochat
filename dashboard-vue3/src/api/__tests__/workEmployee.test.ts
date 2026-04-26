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
import { enterMembersList, syncEmployee, syncTime } from '../workEmployee'

// 导入模拟的 http 模块
import http from '@/utils/http'

describe('WorkEmployee API', () => {
  beforeEach(() => {
    // 清除所有模拟调用
    vi.clearAllMocks()
  })

  it('should export enterMembersList function', () => {
    expect(typeof enterMembersList).toBe('function')
  })

  it('should call get method with correct parameters when enterMembersList', async () => {
    // 模拟响应
    const mockResponse = { list: [], total: 0, page: 1, perPage: 10 }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { page: 1, perPage: 10, name: 'Test' }
    const result = await enterMembersList(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workEmployee/index',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export syncEmployee function', () => {
    expect(typeof syncEmployee).toBe('function')
  })

  it('should call put method with correct parameters when syncEmployee', async () => {
    // 模拟响应
    const mockResponse = { success: true, message: '同步成功' }
    http.put.mockResolvedValue(mockResponse)

    // 调用函数
    const result = await syncEmployee()

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: '/workEmployee/synEmployee'
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export syncTime function', () => {
    expect(typeof syncTime).toBe('function')
  })

  it('should call get method with correct parameters when syncTime', async () => {
    // 模拟响应
    const mockResponse = { lastSyncTime: '2023-01-01 00:00:00', syncStatus: 'success' }
    http.get.mockResolvedValue(mockResponse)

    // 调用函数
    const result = await syncTime()

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workEmployee/searchCondition'
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })
})
