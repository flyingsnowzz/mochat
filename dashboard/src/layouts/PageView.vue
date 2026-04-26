<template>
  <page-header-wrapper :title="false">
    <router-view />
  </page-header-wrapper>
</template>

<script>
/**
 * 页面视图组件
 * 功能说明：带页面头部包装的视图组件
 * 主要功能：
 * 1. 使用 page-header-wrapper 组件包装内容
 * 2. 支持页面标题展示
 * 3. 支持多标签页缓存
 *
 * 业务场景：
 * - 需要展示页面标题的标准页面
 * - 需要统一页面头部样式的页面
 *
 * 技术实现：
 * - 使用 page-header-wrapper 组件提供统一的头部样式
 * - 结合 router-view 实现多标签页功能
 */
export default {
  name: 'PageView'
}
</script>
