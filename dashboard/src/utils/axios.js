/**
 * VueAxios 插件
 * 功能说明：Axios 的 Vue 插件封装
 * 主要功能：
 * 1. 将 axios 实例挂载到 Vue 原型上
 * 2. 提供 Vue.axios 和 this.$http 访问方式
 * 3. 确保 axios 只被安装一次
 *
 * 技术实现：
 * - 创建一个 Vue 插件对象
 * - 在 install 方法中初始化 axios 实例
 * - 使用 Object.defineProperties 定义实例属性
 * - 通过插件机制确保全局只安装一次
 */
const VueAxios = {
  vm: {},
  // eslint-disable-next-line no-unused-vars
  install (Vue, instance) {
    if (this.installed) {
      return
    }
    this.installed = true

    if (!instance) {
      // eslint-disable-next-line no-console
      console.error('You have to install axios')
      return
    }

    Vue.axios = instance

    Object.defineProperties(Vue.prototype, {
      axios: {
        get: function get () {
          return instance
        }
      },
      $http: {
        get: function get () {
          return instance
        }
      }
    })
  }
}

export {
  VueAxios
}
