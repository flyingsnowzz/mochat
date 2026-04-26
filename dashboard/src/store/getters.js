/**
 * Vuex Getters 计算属性
 * 功能说明：提供从 Vuex Store 获取数据的便捷方式
 * 主要功能：
 * 1. 获取移动端状态（isMobile）
 * 2. 获取企业信息（corpId、corpName）
 * 3. 获取用户信息（token、roles、userInfo）
 * 4. 获取路由权限（addRouters、permissionList、defaultRoutePath）
 * 5. 获取面包屑数据（breadcrumb）
 *
 * 技术实现：
 * - 使用箭头函数从对应的 module 中获取状态
 * - 为组件提供响应式的状态访问方式
 */
const getters = {
  isMobile: state => state.app.isMobile,
  corpId: state => state.user.corpId,
  corpName: state => state.user.corpName,
  token: state => state.user.token,
  roles: state => state.user.roles,
  userInfo: state => state.user.userInfo,
  addRouters: state => state.permission.addRouters,
  breadcrumb: state => state.permission.breadcrumb,
  permissionList: state => state.permission.permissionList,
  defaultRoutePath: state => state.permission.defaultRoutePath
}

export default getters
