/**
 * 菜单工具函数
 * 功能说明：提供菜单和路由相关的操作函数
 * 主要函数：
 * 1. exChangeMenu - 切换菜单路由，更新顶部菜单和侧边栏
 * 2. setBreadcrumb - 设置面包屑导航
 * 3. resetRoutes - 重置路由配置
 *
 * 业务场景：
 * - 页面跳转时更新菜单选中状态
 * - 动态更新面包屑导航
 * - 退出登录时清空动态路由
 *
 * 技术实现：
 * - 操作 Vuex store 管理菜单状态
 * - 直接操作 VueRouter 实例
 * - 通过 mutations 修改状态
 */
import store from '@/store'
import router, { newRouter } from '@/router'

// 更改菜单路由
export function exChangeMenu (path) {
  const firstMenu = store.getters.permissionList.filter(item => {
    return item.routes.includes(path)
  })
  if (firstMenu.length !== 0) {
    const title = store.state.permission.topMenuKey.title
    if (firstMenu[0].title !== title) {
      store.commit('SET_TOP_MENU_KEY', firstMenu[0])
    }
    const { children } = firstMenu[0]
    if (children.length > 0) {
      store.commit('SET_SIDE_MENUS', children)
    }
  }
}
export function setBreadcrumb (path) {
  const sideMenus = store.state.permission.sideMenus
  let firstTitle = ''
  let secondTitle = ''
  sideMenus.find(item => {
    const title = item.meta.title
    item.children.find(inner => {
      if (inner.path == path) {
        firstTitle = title
        secondTitle = inner.meta.title
        return true
      }
    })
  })
  store.commit('SET_BREADCRUMB', [firstTitle, secondTitle])
}
export function resetRoutes () {
  router.matcher = newRouter().matcher
  store.commit('CLEAR_ROUTERS')
}
