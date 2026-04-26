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
import { workRoomGroupList, deleteGroup, createGroup, updateGroup, workRoomList, synList, batchUpdate, workContactRoom, workDepartmentList, departmentList, statisticsIndex, statistics, workRoomAutoPullList, autoPullUpdate, autoPullCreate, autoPullMove, autoPullShow, workContactTagGroup, addWorkContactTag, tagList, roomList, optCroup, tagGetList, addGroup, labelShow, delRoomTag, remindRoomTag, chooseContactRoomTag, chooseFilterContact, labelContactShow, showRoomApi } from '../workRoom'

// 导入模拟的 http 模块
import http from '@/utils/http'

describe('WorkRoom API', () => {
  beforeEach(() => {
    // 清除所有模拟调用
    vi.clearAllMocks()
  })

  it('should export workRoomGroupList function', () => {
    expect(typeof workRoomGroupList).toBe('function')
  })

  it('should call get method with correct parameters when workRoomGroupList', async () => {
    // 模拟响应
    const mockResponse = { list: [], total: 0, page: 1, perPage: 10 }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { page: 1, perPage: 10 }
    const result = await workRoomGroupList(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workRoomGroup/index',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export deleteGroup function', () => {
    expect(typeof deleteGroup).toBe('function')
  })

  it('should call del method with correct parameters when deleteGroup', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.del.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1 }
    const result = await deleteGroup(params)

    // 验证调用
    expect(http.del).toHaveBeenCalledWith({
      url: '/workRoomGroup/destroy',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export createGroup function', () => {
    expect(typeof createGroup).toBe('function')
  })

  it('should call post method with correct parameters when createGroup', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.post.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { name: 'Test Group' }
    const result = await createGroup(params)

    // 验证调用
    expect(http.post).toHaveBeenCalledWith({
      url: '/workRoomGroup/store',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export updateGroup function', () => {
    expect(typeof updateGroup).toBe('function')
  })

  it('should call put method with correct parameters when updateGroup', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.put.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1, name: 'Test Group' }
    const result = await updateGroup(params)

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: '/workRoomGroup/update',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export workRoomList function', () => {
    expect(typeof workRoomList).toBe('function')
  })

  it('should call get method with correct parameters when workRoomList', async () => {
    // 模拟响应
    const mockResponse = { list: [], total: 0, page: 1, perPage: 10 }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { page: 1, perPage: 10, name: 'Test Room' }
    const result = await workRoomList(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workRoom/index',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export synList function', () => {
    expect(typeof synList).toBe('function')
  })

  it('should call put method with correct parameters when synList', async () => {
    // 模拟响应
    const mockResponse = { success: true, message: '同步成功' }
    http.put.mockResolvedValue(mockResponse)

    // 调用函数
    const result = await synList()

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: '/workRoom/syn'
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
    const params = { ids: [1, 2], groupId: 3 }
    const result = await batchUpdate(params)

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: '/workRoom/batchUpdate',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export workContactRoom function', () => {
    expect(typeof workContactRoom).toBe('function')
  })

  it('should call get method with correct parameters when workContactRoom', async () => {
    // 模拟响应
    const mockResponse = { list: [], total: 0, page: 1, perPage: 10 }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { roomId: 1, page: 1, perPage: 10 }
    const result = await workContactRoom(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workContactRoom/index',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export workDepartmentList function', () => {
    expect(typeof workDepartmentList).toBe('function')
  })

  it('should call get method with correct parameters when workDepartmentList', async () => {
    // 模拟响应
    const mockResponse = [{ id: 1, name: 'Test Department', avatar: 'avatar.jpg', position: 'Manager' }]
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { departmentId: 1 }
    const result = await workDepartmentList(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workEmployeeDepartment/memberIndex',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export departmentList function', () => {
    expect(typeof departmentList).toBe('function')
  })

  it('should call get method with correct parameters when departmentList', async () => {
    // 模拟响应
    const mockResponse = [{ id: 1, name: 'Test Department', parentId: 0 }]
    http.get.mockResolvedValue(mockResponse)

    // 调用函数
    const result = await departmentList()

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workDepartment/index'
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export statisticsIndex function', () => {
    expect(typeof statisticsIndex).toBe('function')
  })

  it('should call get method with correct parameters when statisticsIndex', async () => {
    // 模拟响应
    const mockResponse = { list: [], total: 0, page: 1, perPage: 10 }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { page: 1, perPage: 10, roomId: 1 }
    const result = await statisticsIndex(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workRoom/statisticsIndex',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export statistics function', () => {
    expect(typeof statistics).toBe('function')
  })

  it('should call get method with correct parameters when statistics', async () => {
    // 模拟响应
    const mockResponse = { dates: [], addMembers: [], quitMembers: [], totalMembers: [] }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { roomId: 1, startDate: '2023-01-01', endDate: '2023-01-31' }
    const result = await statistics(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workRoom/statistics',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export workRoomAutoPullList function', () => {
    expect(typeof workRoomAutoPullList).toBe('function')
  })

  it('should call get method with correct parameters when workRoomAutoPullList', async () => {
    // 模拟响应
    const mockResponse = { list: [], total: 0, page: 1, perPage: 10 }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { page: 1, perPage: 10 }
    const result = await workRoomAutoPullList(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workRoomAutoPull/index',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export autoPullUpdate function', () => {
    expect(typeof autoPullUpdate).toBe('function')
  })

  it('should call put method with correct parameters when autoPullUpdate', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.put.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1, name: 'Test Auto Pull', roomId: 1, tagIds: [1, 2], employeeId: 1 }
    const result = await autoPullUpdate(params)

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: '/workRoomAutoPull/update',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export autoPullCreate function', () => {
    expect(typeof autoPullCreate).toBe('function')
  })

  it('should call post method with correct parameters when autoPullCreate', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.post.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { name: 'Test Auto Pull', roomId: 1, tagIds: [1, 2], employeeId: 1 }
    const result = await autoPullCreate(params)

    // 验证调用
    expect(http.post).toHaveBeenCalledWith({
      url: '/workRoomAutoPull/store',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export autoPullMove function', () => {
    expect(typeof autoPullMove).toBe('function')
  })

  it('should call put method with correct parameters when autoPullMove', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.put.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1, targetRoomId: 2 }
    const result = await autoPullMove(params)

    // 验证调用
    expect(http.put).toHaveBeenCalledWith({
      url: '/workRoomAutoPull/move',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export autoPullShow function', () => {
    expect(typeof autoPullShow).toBe('function')
  })

  it('should call get method with correct parameters when autoPullShow', async () => {
    // 模拟响应
    const mockResponse = { id: 1, name: 'Test Auto Pull', room: { id: 1, name: 'Test Room' }, tags: [], employee: { id: 1, name: 'Test Employee' } }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1 }
    const result = await autoPullShow(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workRoomAutoPull/show',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export workContactTagGroup function', () => {
    expect(typeof workContactTagGroup).toBe('function')
  })

  it('should call get method with correct parameters when workContactTagGroup', async () => {
    // 模拟响应
    const mockResponse = [{ id: 1, name: 'Test Tag Group', tags: [] }]
    http.get.mockResolvedValue(mockResponse)

    // 调用函数
    const result = await workContactTagGroup()

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workContactTagGroup/index'
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export addWorkContactTag function', () => {
    expect(typeof addWorkContactTag).toBe('function')
  })

  it('should call post method with correct parameters when addWorkContactTag', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.post.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { groupId: 1, name: 'Test Tag', color: '#ff0000' }
    const result = await addWorkContactTag(params)

    // 验证调用
    expect(http.post).toHaveBeenCalledWith({
      url: '/workContactTag/store',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export tagList function', () => {
    expect(typeof tagList).toBe('function')
  })

  it('should call get method with correct parameters when tagList', async () => {
    // 模拟响应
    const mockResponse = [{ id: 1, name: 'Test Tag', color: '#ff0000', groupId: 1 }]
    http.get.mockResolvedValue(mockResponse)

    // 调用函数
    const result = await tagList()

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workContactTag/allTag'
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export roomList function', () => {
    expect(typeof roomList).toBe('function')
  })

  it('should call get method with correct parameters when roomList', async () => {
    // 模拟响应
    const mockResponse = [{ id: 1, name: 'Test Room', memberCount: 10 }]
    http.get.mockResolvedValue(mockResponse)

    // 调用函数
    const result = await roomList()

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/workRoom/roomIndex'
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export optCroup function', () => {
    expect(typeof optCroup).toBe('function')
  })

  it('should call get method with correct parameters when optCroup', async () => {
    // 模拟响应
    const mockResponse = [{ id: 1, name: 'Test Room' }]
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { tagIds: [1, 2] }
    const result = await optCroup(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/roomTagPull/roomList',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export tagGetList function', () => {
    expect(typeof tagGetList).toBe('function')
  })

  it('should call get method with correct parameters when tagGetList', async () => {
    // 模拟响应
    const mockResponse = { list: [], total: 0, page: 1, perPage: 10 }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { page: 1, perPage: 10 }
    const result = await tagGetList(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/roomTagPull/index',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export addGroup function', () => {
    expect(typeof addGroup).toBe('function')
  })

  it('should call post method with correct parameters when addGroup', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.post.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { name: 'Test Tag Group', roomId: 1, tagIds: [1, 2], employeeId: 1 }
    const result = await addGroup(params)

    // 验证调用
    expect(http.post).toHaveBeenCalledWith({
      url: '/roomTagPull/store',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export labelShow function', () => {
    expect(typeof labelShow).toBe('function')
  })

  it('should call get method with correct parameters when labelShow', async () => {
    // 模拟响应
    const mockResponse = { id: 1, name: 'Test Tag Group', room: { id: 1, name: 'Test Room' }, tags: [], employee: { id: 1, name: 'Test Employee' }, contactCount: 10, status: 'active' }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1 }
    const result = await labelShow(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/roomTagPull/show',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export delRoomTag function', () => {
    expect(typeof delRoomTag).toBe('function')
  })

  it('should call del method with correct parameters when delRoomTag', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.del.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1 }
    const result = await delRoomTag(params)

    // 验证调用
    expect(http.del).toHaveBeenCalledWith({
      url: '/roomTagPull/destroy',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export remindRoomTag function', () => {
    expect(typeof remindRoomTag).toBe('function')
  })

  it('should call get method with correct parameters when remindRoomTag', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1 }
    const result = await remindRoomTag(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/roomTagPull/remindSend',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export chooseContactRoomTag function', () => {
    expect(typeof chooseContactRoomTag).toBe('function')
  })

  it('should call get method with correct parameters when chooseContactRoomTag', async () => {
    // 模拟响应
    const mockResponse = { list: [], total: 0, page: 1, perPage: 10 }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1, page: 1, perPage: 10 }
    const result = await chooseContactRoomTag(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/roomTagPull/chooseContact',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export chooseFilterContact function', () => {
    expect(typeof chooseFilterContact).toBe('function')
  })

  it('should call post method with correct parameters when chooseFilterContact', async () => {
    // 模拟响应
    const mockResponse = { success: true }
    http.post.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1, contactIds: [1, 2] }
    const result = await chooseFilterContact(params)

    // 验证调用
    expect(http.post).toHaveBeenCalledWith({
      url: '/roomTagPull/filterContact',
      data: params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export labelContactShow function', () => {
    expect(typeof labelContactShow).toBe('function')
  })

  it('should call get method with correct parameters when labelContactShow', async () => {
    // 模拟响应
    const mockResponse = { list: [], total: 0, page: 1, perPage: 10 }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { id: 1, page: 1, perPage: 10 }
    const result = await labelContactShow(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/roomTagPull/showContact',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })

  it('should export showRoomApi function', () => {
    expect(typeof showRoomApi).toBe('function')
  })

  it('should call get method with correct parameters when showRoomApi', async () => {
    // 模拟响应
    const mockResponse = { id: 1, name: 'Test Room', avatar: 'avatar.jpg', owner: { id: 1, name: 'Test Owner' }, memberCount: 10, createTime: '2023-01-01 00:00:00' }
    http.get.mockResolvedValue(mockResponse)

    // 测试参数
    const params = { roomId: 1 }
    const result = await showRoomApi(params)

    // 验证调用
    expect(http.get).toHaveBeenCalledWith({
      url: '/contactMessageBatchSend/showRoom',
      params
    })

    // 验证返回值
    expect(result).toEqual(mockResponse)
  })
})
