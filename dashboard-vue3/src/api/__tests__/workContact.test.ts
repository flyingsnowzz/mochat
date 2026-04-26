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
import { getUserPortrait, editUserPortrait, addTag, getWorkContactInfo, editWorkContactInfo, workContactList, allTag, groupChatList, customersSource, UserPortraitList, synContact, track } from '../workContact'

// 导入模拟的 http 模块
import http from '@/utils/http'

describe('WorkContact API', () => {
  beforeEach(() => {
    // 清除所有模拟调用
    vi.clearAllMocks()
  })

  it('should export getUserPortrait function', () => {
    expect(typeof getUserPortrait).toBe('function')
  })

  it('should call get method with correct parameters when getUserPortrait', async () => {
    // 模拟响应
    const mockResponse = { fields: [] }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { contactId: 1 }
    const result = await getUserPortrait(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/contactFieldPivot/index',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export editUserPortrait function', () => {
    expect(typeof editUserPortrait).toBe('function')
  })

  it('should call put method with correct parameters when editUserPortrait', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.put.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { contactId: 1, fields: [{ id: 1, value: 'test' }] }
    const result = await editUserPortrait(params)

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: '/contactFieldPivot/update',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export addTag function', () => {
    expect(typeof addTag).toBe('function')
  })

  it('should call post method with correct parameters when addTag', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.post.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { contactIds: [1, 2], tagIds: [3, 4] }
    const result = await addTag(params)

    // 验证调用
    expect(http.post).toHaveBeenCalledWith({
      url: '/workContact/batchLabeling',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export getWorkContactInfo function', () => {
    expect(typeof getWorkContactInfo).toBe('function')
  })

  it('should call get method with correct parameters when getWorkContactInfo', async () => {
    // 模拟响应
    const mockResponse = { id: 1, name: 'Test User', phone: '13800138000' }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1 }
    const result = await getWorkContactInfo(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workContact/show',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export editWorkContactInfo function', () => {
    expect(typeof editWorkContactInfo).toBe('function')
  })

  it('should call put method with correct parameters when editWorkContactInfo', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.put.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1, name: 'Test User', phone: '13800138000' }
    const result = await editWorkContactInfo(params)

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: '/workContact/update',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export workContactList function', () => {
    expect(typeof workContactList).toBe('function')
  })

  it('should call get method with correct parameters when workContactList', async () => {
    // 模拟响应
    const mockResponse = { list: [], total: 0, page: 1, perPage: 10 }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { page: 1, perPage: 10, name: 'Test' }
    const result = await workContactList(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workContact/index',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export allTag function', () => {
    expect(typeof allTag).toBe('function')
  })

  it('should call get method with correct parameters when allTag', async () => {
    // 模拟响应
    const mockResponse = [{ id: 1, name: 'Test Tag', color: '#ff0000' }]
    http.get.mockResolvedValue(mockResponse)

    // 调用函数
    const result = await allTag()

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workContactTag/allTag'
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export groupChatList function', () => {
    expect(typeof groupChatList).toBe('function')
  })

  it('should call get method with correct parameters when groupChatList', async () => {
    // 模拟响应
    const mockResponse = [{ id: 1, name: 'Test Group', memberCount: 10 }]
    http.get.mockResolvedValue(mockResponse)

    // 调用函数
    const result = await groupChatList()

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workRoom/roomIndex'
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export customersSource function', () => {
    expect(typeof customersSource).toBe('function')
  })

  it('should call get method with correct parameters when customersSource', async () => {
    // 模拟响应
    const mockResponse = [{ value: 'wechat', label: '微信' }]
    http.get.mockResolvedValue(mockResponse)

    // 调用函数
    const result = await customersSource()

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workContact/source'
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export UserPortraitList function', () => {
    expect(typeof UserPortraitList).toBe('function')
  })

  it('should call get method with correct parameters when UserPortraitList', async () => {
    // 模拟响应
    const mockResponse = { list: [], total: 0, page: 1, perPage: 10 }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { page: 1, perPage: 10, fieldId: 1 }
    const result = await UserPortraitList(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/contactField/portrait',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export synContact function', () => {
    expect(typeof synContact).toBe('function')
  })

  it('should call put method with correct parameters when synContact', async () => {
    // 模拟响应
    const mockResponse = { success: true, message: '同步成功' }
    http.put.mockResolvedValue(mockResponse)

    // 调用函数
    const result = await synContact()

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: '/workContact/synContact'
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export track function', () => {
    expect(typeof track).toBe('function')
  })

  it('should call get method with correct parameters when track', async () => {
    // 模拟响应
    const mockResponse = { records: [] }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { contactId: 1 }
    const result = await track(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workContact/track',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })
})
