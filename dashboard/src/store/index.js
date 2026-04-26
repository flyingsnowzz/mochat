/**
 * Vuex Store 入口文件
 * 功能说明：配置和初始化 Vuex 状态管理
 * 主要功能：
 * 1. 引入并注册 Vuex 模块（app、user、permission）
 * 2. 配置模块化状态管理
 * 3. 暴露 getters 给组件使用
 *
 * 模块说明：
 * - app: 应用全局状态（侧边栏折叠、颜色主题、语言等）
 * - user: 用户状态（用户信息、登录状态、权限等）
 * - permission: 路由权限状态（动态路由、菜单权限等）
 *
 * 技术实现：
 * - 使用 Vuex 模块化组织状态
 * - 使用 getters 计算属性获取状态
 */
import Vue from 'vue'
import Vuex from 'vuex'

import app from './modules/app'
import user from './modules/user'

// default router permission control
import permission from './modules/permission'

import getters from './getters'

Vue.use(Vuex)

export default new Vuex.Store({
  modules: {
    app,
    user,
    permission
  },
  state: {

  },
  mutations: {

  },
  actions: {

  },
  getters
})
