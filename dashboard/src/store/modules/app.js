/**
 * 应用状态模块
 * 功能说明：管理应用全局状态
 * 主要功能：
 * 1. 侧边栏折叠状态管理（sideCollapsed）
 * 2. 移动端适配状态管理（isMobile）
 * 3. 多标签页模式管理（multiTab）
 *
 * 业务场景：
 * - 控制后台管理系统的侧边栏展开/收起
 * - 检测用户设备类型，适配移动端/PC端
 * - 支持多标签页同时打开多个页面
 *
 * 技术实现：
 * - 使用 storage 持久化侧边栏状态
 * - mutations 同步修改状态
 * - 状态变更自动保存到本地存储
 */
import storage from 'store'

const app = {
  state: {
    sideCollapsed: false,
    isMobile: false,
    multiTab: true
  },
  mutations: {
    SIDEBAR_TYPE: (state, type) => {
      state.sideCollapsed = type
      storage.set('SIDEBAR_TYPE', type)
    },
    TOGGLE_MOBILE_TYPE: (state, isMobile) => {
      state.isMobile = isMobile
    }
  },
  actions: {

  }
}

export default app
