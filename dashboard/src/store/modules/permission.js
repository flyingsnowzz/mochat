/**
 * 权限状态模块
 * 功能说明：管理路由权限和菜单权限
 * 主要功能：
 * 1. 获取用户权限列表（getPermissionList）
 * 2. 动态生成路由（addRouters）
 * 3. 生成侧边栏菜单（sideMenus）
 * 4. 顶部菜单状态管理（topMenuKey）
 * 5. 面包屑导航管理（breadcrumb）
 * 6. 默认路由路径管理（defaultRoutePath）
 *
 * 业务场景：
 * - 用户登录后，根据权限动态生成可访问的路由
 * - 根据用户权限显示不同的菜单项
 * - 管理页面的面包屑导航
 * - 控制用户对特定页面的访问权限
 *
 * 技术实现：
 * - 使用 dealPermissionData 处理权限数据，生成层级菜单结构
 * - 将权限数据分为顶部菜单（topMenus）和二级菜单（secondMenus）
 * - 结合 errorPage 路由处理权限不足的情况
 * - 使用 async/await 处理异步权限获取
 */
import { dealPermissionData } from '@/router/router.config'
import { errorPage } from '@/router/base/error'
import { permissionByUser } from '@/api/login'

const permission = {
  state: {
    addRouters: [],
    sideMenus: [],
    topMenuKey: '',
    breadcrumb: [],
    permissionList: [],
    defaultRoutePath: ''
  },
  mutations: {
    SET_ROUTERS: (state, routers) => {
      state.addRouters = routers.concat(errorPage)
    },
    CLEAR_ROUTERS: (state) => {
      state.addRouters = []
    },
    SET_SIDE_MENUS: (state, sideMenus) => {
      state.sideMenus = sideMenus
    },
    SET_TOP_MENU_KEY: (state, topMenuKey) => {
      state.topMenuKey = topMenuKey
    },
    SET_BREADCRUMB: (state, breadcrumb) => {
      state.breadcrumb = breadcrumb
    },
    SET_PERMISSION_LIST: (state, permissionList) => {
      state.permissionList = permissionList
    },
    SET_DEFAULT_ROUTER_PATH: (state, defaultRoutePath) => {
      state.defaultRoutePath = defaultRoutePath
    }
  },
  actions: {
    async getPermissionList ({ commit }) {
      try {
        const { data } = await permissionByUser()
        if (data.length > 0) {
          const { topMenus, secondMenus, path } = dealPermissionData(data)
          commit('SET_PERMISSION_LIST', topMenus)
          commit('SET_ROUTERS', secondMenus)
          commit('SET_DEFAULT_ROUTER_PATH', path)
          return secondMenus
        }
      } catch (e) {
        console.log(e)
      }
    }
  }
}

export default permission
