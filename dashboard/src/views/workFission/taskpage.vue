<template>
  <a-card>
    <a-tabs type="card" v-model="tabPage" @change="switchPage">
      <a-tab-pane :key="1" tab="企微任务宝">
        <workFissionIndex />
      </a-tab-pane>
      <a-tab-pane :key="2" tab="数据详情" force-render>
        <workFissionDataShow />
      </a-tab-pane>
    </a-tabs>
  </a-card>
</template>
<script>
/**
 * 企微任务宝主页面
 * 功能说明：企微任务宝活动的综合管理页面
 * 主要功能：
 * 1. 企微任务宝：管理任务宝活动列表
 * 2. 数据详情：查看活动的数据统计分析
 *
 * 业务场景：
 * - 集中管理企微任务宝活动
 * - 查看活动数据和效果分析
 *
 * 技术实现：
 * - 使用 a-tabs 组件切换不同功能模块
 * - 引入 workFissionIndex 和 workFissionDataShow 子组件
 */
import workFissionIndex from '@/views/workFission/index'
import workFissionDataShow from '@/views/workFission/dataShow'
export default {
  components: {
    workFissionIndex,
    workFissionDataShow
  },
  data () {
    return {
      tabPage: 1
    }
  },
  methods: {
    switchPage () {
    }
  }
}
</script>
