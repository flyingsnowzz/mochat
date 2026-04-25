<template>
  <a-result status="404" title="404" sub-title="对不起，MoChat该页面未找到！">
    <template #extra>
      <a-button type="primary" @click="toHome">
        回首页
      </a-button>
    </template>
  </a-result>
</template>

<script>
/**
 * 404错误页面
 * 功能说明：当用户访问不存在的页面时显示
 * 主要功能：
 * 1. 显示404错误提示
 * 2. 提供返回首页按钮
 *
 * 技术实现：
 * - 使用 a-result 组件展示404状态
 * - 点击按钮跳转到首页
 */
export default {
  name: 'Error404',
  methods: {
    toHome () {
      this.$router.push({ path: '/' })
    }
  }
}
</script>
