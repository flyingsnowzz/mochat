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
import { subManagementList, addSubManagement, editSubManagement, editStatus, getSubManagement, changeStatus, selectByPhone, selectRole, passwordResetApi } from '../user'

// 导入模拟的 http 模块
import http from '@/utils/http'

describe('User API', () => {
  beforeEach(() => {
    // 清除所有模拟调用
    vi.clearAllMocks()
  })

  it('should export subManagementList function', () => {
    expect(typeof subManagementList).toBe('function')
  })

  it('should call get method with correct parameters when subManagementList', async () => {
    // 模拟响应
    const mockResponse = { list: [], total: 0, page: 1, perPage: 10 }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { page: 1, perPage: 10 }
    const result = await subManagementList(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/user/index',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export addSubManagement function', () => {
    expect(typeof addSubManagement).toBe('function')
  })

  it('should call post method with correct parameters when addSubManagement', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.post.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { name: 'Test User', phone: '13800138000', roleId: 1 }
    const result = await addSubManagement(params)

    // 验证调用
    expect(http.post).toHaveBeenCalledWith({
      url: '/user/store',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export editSubManagement function', () => {
    expect(typeof editSubManagement).toBe('function')
  })

  it('should call put method with correct parameters when editSubManagement', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.put.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1, name: 'Test User', phone: '13800138000', roleId: 1 }
    const result = await editSubManagement(params)

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: '/user/update',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export editStatus function', () => {
    expect(typeof editStatus).toBe('function')
  })

  it('should call put method with correct parameters when editStatus', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.put.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1, status: 1 }
    const result = await editStatus(params)

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: '/user/statusUpdate',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export getSubManagement function', () => {
    expect(typeof getSubManagement).toBe('function')
  })

  it('should call get method with correct parameters when getSubManagement', async () => {
    // 模拟响应
    const mockResponse = { id: 1, name: 'Test User', phone: '13800138000' }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1 }
    const result = await getSubManagement(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/user/show',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export changeStatus function', () => {
    expect(typeof changeStatus).toBe('function')
  })

  it('should call put method with correct parameters when changeStatus', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.put.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1, status: 1 }
    const result = await changeStatus(params)

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: '/user/statusUpdate',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export selectByPhone function', () => {
    expect(typeof selectByPhone).toBe('function')
  })

  it('should call get method with correct parameters when selectByPhone', async () => {
    // 模拟响应
    const mockResponse = { id: 1, name: 'Test Department' }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { phone: '13800138000' }
    const result = await selectByPhone(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workDepartment/selectByPhone',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export selectRole function', () => {
    expect(typeof selectRole).toBe('function')
  })

  it('should call get method with correct parameters when selectRole', async () => {
    // 模拟响应
    const mockResponse = [{ id: 1, name: 'Admin' }, { id: 2, name: 'User' }]
    http.get.mockResolvedValue(mockResponse)

    // 调用函数
    const result = await selectRole()

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/role/select'
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export passwordResetApi function', () => {
    expect(typeof passwordResetApi).toBe('function')
  })

  it('should call put method with correct parameters when passwordResetApi', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.put.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1 }
    const result = await passwordResetApi(params)

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: '/user/passwordReset',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })
})
