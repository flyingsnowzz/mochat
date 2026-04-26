/**
 * 设备检测 Mixin
 * 功能说明：提供设备类型检测的混入对象
 * 主要功能：
 * 1. 将 Vuex 中的 isMobile 状态映射到组件的计算属性
 * 2. 方便组件判断当前是否为移动端设备
 *
 * 业务场景：
 * - 在组件中判断设备类型，渲染不同的 UI
 * - 移动端采用更简洁的布局
 * - PC 端展示完整功能
 *
 * 技术实现：
 * - 使用 Vuex mapState 将状态映射为计算属性
 * - 通过 mixin 方式复用给多个组件
 */
import { mapState } from 'vuex'

const deviceMixin = {
  computed: {
    ...mapState({
      isMobile: state => state.app.isMobile
    })
  }
}

export { deviceMixin }
