<script>
/**
 * 路由视图组件
 * 功能说明：支持页面缓存的路由视图组件
 * 主要功能：
 * 1. 根据 keepAlive 属性决定是否缓存页面
 * 2. 支持路由 meta.keepAlive 配置
 * 3. 使用 keep-alive 实现页面缓存
 *
 * 业务场景：
 * - 需要缓存的页面（如列表页返回后保持滚动位置）
 * - 不需要缓存的页面（如表单页每次进入需要刷新）
 *
 * 技术实现：
 * - 使用 render 函数动态渲染组件
 * - 使用 keep-alive 包裹 router-view 实现缓存
 * - 支持组件级和路由级两种缓存配置
 */
export default {
  name: 'RouteView',
  props: {
    keepAlive: {
      type: Boolean,
      default: true
    }
  },
  data () {
    return {}
  },
  render () {
    const { $route: { meta } } = this
    const inKeep = (
      <keep-alive>
        <router-view />
      </keep-alive>
    )
    const notKeep = (
      <router-view />
    )
    return this.keepAlive || meta.keepAlive ? inKeep : notKeep
  }
}
</script>
