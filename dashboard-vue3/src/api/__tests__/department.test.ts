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
import { departmentList, showEmployee } from '../department'

// 导入模拟的 http 模块
import http from '@/utils/http'

describe('Department API', () => {
  beforeEach(() => {
    // 清除所有模拟调用
    vi.clearAllMocks()
  })

  it('should export departmentList function', () => {
    expect(typeof departmentList).toBe('function')
  })

  it('should call get method with correct parameters when departmentList', async () => {
    // 模拟响应
    const mockResponse = { list: [], total: 0, page: 1, perPage: 10 }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { page: 1, perPage: 10 }
    const result = await departmentList(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workDepartment/pageIndex',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export showEmployee function', () => {
    expect(typeof showEmployee).toBe('function')
  })

  it('should call get method with correct parameters when showEmployee', async () => {
    // 模拟响应
    const mockResponse = { list: [], total: 0, page: 1, perPage: 10 }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { departmentId: 1, page: 1, perPage: 10 }
    const result = await showEmployee(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workDepartment/showEmployee',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })
})
