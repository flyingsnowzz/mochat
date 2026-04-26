/**
 * 用户状态模块
 * 功能说明：管理用户登录状态和用户信息
 * 主要功能：
 * 1. 用户登录（Login）
 * 2. 获取用户信息（GetInfo）
 * 3. 用户登出（Logout）
 * 4. 状态管理：token、roles、userInfo、corpId、corpName
 *
 * 业务场景：
 * - 用户登录系统，获取 Token 和用户信息
 * - 保存用户登录状态到本地存储
 * - 管理当前登录企业的 ID 和名称
 * - 存储用户权限角色
 *
 * 技术实现：
 * - 使用 Promise 封装异步操作
 * - 使用 storage 持久化 Token
 * - Token 格式：Bearer {token}
 * - mutations 同步修改状态
 */
import storage from 'store'
import { login, getInfo } from '@/api/login'

const user = {
  state: {
    corpId: undefined,
    corpName: '',
    token: '',
    roles: [],
    userInfo: null
  },

  mutations: {
    SET_TOKEN: (state, token) => {
      state.token = token
    },
    SET_ROLES: (state, roles) => {
      state.roles = roles
    },
    SET_USER_INFO: (state, userInfo) => {
      state.userInfo = userInfo
    },
    SET_CORP_ID: (state, corpId) => {
      state.corpId = corpId
    },
    SET_CORP_NAME: (state, corpName) => {
      state.corpName = corpName
    }
  },

  actions: {
    // 登录
    Login ({ commit }, userInfo) {
      return new Promise((resolve, reject) => {
        login(userInfo).then(response => {
          const result = response.data
          const token = `Bearer ${result.token}`
          storage.set('ACCESS_TOKEN', token, result.expire * 60 * 1000)
          commit('SET_TOKEN', token)
          resolve()
        }).catch(error => {
          reject(error)
        })
      })
    },

    // 获取用户信息
    GetInfo ({ commit }) {
      return new Promise((resolve, reject) => {
        getInfo().then(response => {
          const data = response.data
          commit('SET_USER_INFO', data)
          resolve(response)
        }).catch(error => {
          reject(error)
        })
      })
    },

    // 登出
    Logout ({ commit, state }) {
      return new Promise((resolve) => {
        commit('SET_TOKEN', '')
        storage.remove('ACCESS_TOKEN')
        resolve()
      })
    }

  }
}

export default user
