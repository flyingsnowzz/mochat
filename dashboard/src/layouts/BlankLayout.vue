<template>
  <div>
    <router-view />
  </div>
</template>

<script>
/**
 * 空布局组件
 * 功能说明：简单的空布局，不包含侧边栏、头部等装饰元素
 * 主要功能：
 * 1. 仅渲染路由内容
 * 2. 用于不需要完整布局的页面（如登录页、错误页）
 *
 * 业务场景：
 * - 登录页面
 * - 404/403 错误页面
 * - 其他全屏展示的页面
 *
 * 技术实现：
 * - 最简化的布局组件
 * - 仅包含 router-view 占位符
 */
export default {
  name: 'BlankLayout'
}
</script>

<style scoped>

</style>
