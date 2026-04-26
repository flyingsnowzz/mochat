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
import { addInformation, getEnterMembers, wechatAuthList } from '../workMessageConfig'

// 导入模拟的 http 模块
import http from '@/utils/http'

describe('WorkMessageConfig API', () => {
  beforeEach(() => {
    // 清除所有模拟调用
    vi.clearAllMocks()
  })

  it('should export addInformation function', () => {
    expect(typeof addInformation).toBe('function')
  })

  it('should call post method with correct parameters when addInformation', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.post.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { name: 'Test Corp', corpid: 'wx123456', agentid: '1000001' }
    const result = await addInformation(params)

    // 验证调用
    expect(http.post).toHaveBeenCalledWith({
      url: '/workMessageConfig/corpStore',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export getEnterMembers function', () => {
    expect(typeof getEnterMembers).toBe('function')
  })

  it('should call get method with correct parameters when getEnterMembers', async () => {
    // 模拟响应
    const mockResponse = { id: 1, name: 'Test Corp', members: [] }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1 }
    const result = await getEnterMembers(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workMessageConfig/corpShow',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export wechatAuthList function', () => {
    expect(typeof wechatAuthList).toBe('function')
  })

  it('should call get method with correct parameters when wechatAuthList', async () => {
    // 模拟响应
    const mockResponse = { list: [], total: 0, page: 1, perPage: 10 }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { page: 1, perPage: 10 }
    const result = await wechatAuthList(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workMessageConfig/corpIndex',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })
})
