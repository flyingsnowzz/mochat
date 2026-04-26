<template>
  <pro-layout
    :menus="sideMenus"
    :collapsed="sideCollapsed"
    :mediaQuery="query"
    :isMobile="isMobile"
    :handleMediaQuery="handleMediaQuery"
    :handleCollapse="handleCollapse"
    :logo="logoRender"
    v-bind="settings"
  >
    <template v-slot:rightContentRender>
      <global-header :is-mobile="isMobile" :theme="settings.theme"></global-header>
    </template>
    <template v-slot:footerRender>
      <div class="footer">
        Powered by <a class="mochat" href="https://mo.chat/" target="_blank">MoChat</a>
      </div>
    </template>
    <div class="breadcrumb-wrapper">
      <div class="chunk"></div>
      <div class="breadcrumb">
        <a-breadcrumb>
          <a-breadcrumb-item v-for="(item, index) in breadcrumb" :key="index">{{ item }}</a-breadcrumb-item>
        </a-breadcrumb>
      </div>
    </div>
    <router-view />
  </pro-layout>
</template>

<script>
/**
 * 基础布局组件
 * 功能说明：应用的基础布局结构，包含侧边栏、头部、内容区域
 * 主要功能：
 * 1. 使用 pro-layout 组件实现高级布局
 * 2. 侧边栏菜单展示（支持折叠/展开）
 * 3. 响应式布局（自动适配 PC/移动端）
 * 4. 面包屑导航展示
 * 5. 顶部右侧区域（用户信息、通知等）
 * 6. 页脚展示
 *
 * 布局配置：
 * - layout: sidemenu（侧边菜单模式）
 * - theme: dark（暗色主题）
 * - fixSiderbar: true（固定侧边栏）
 * - fixedHeader: true（固定头部）
 *
 * 技术实现：
 * - 使用 pro-layout 组件库实现高级布局
 * - 使用 Vuex mapState 获取菜单和折叠状态
 * - 使用 GlobalHeader 组件渲染顶部区域
 * - 使用 Logo 组件渲染侧边栏顶部 Logo
 */
import { mapState, mapMutations, mapGetters } from 'vuex'
import GlobalHeader from '@/components/GlobalHeader/index'
import Logo from '@/components/Logo'

export default {
  name: 'BasicLayout',
  components: {
    GlobalHeader
  },
  data () {
    return {
      // 侧栏收起状态
      collapsed: false,
      settings: {
        // 布局类型
        layout: 'sidemenu', // 'sidemenu', 'topmenu'
        // CONTENT_WIDTH_TYPE
        contentWidth: 'Fluid',
        // 主题 'dark' | 'light'
        theme: 'dark',
        // 主色调
        primaryColor: '#1890ff',
        fixSiderbar: true,
        fixedHeader: true
      },
      // 媒体查询
      query: {},
      // 是否手机模式
      isMobile: false,
      breadcrumbText: ''
    }
  },
  computed: {
    ...mapState({
      sideMenus: state => state.permission.sideMenus,
      sideCollapsed: state => state.app.sideCollapsed
    }),
    ...mapGetters(['breadcrumb'])
  },
  created () {
    // 处理侧栏收起状态
    // this.$watch('collapsed', () => {
    //   this.sidebarType(this.collapsed)
    // })
    this.$watch('isMobile', () => {
      this.toggleMobileType(this.isMobile)
    })
  },
  mounted () {
    const userAgent = navigator.userAgent
    if (userAgent.indexOf('Edge') > -1) {
      this.$nextTick(() => {
        this.sidebarType(!this.sideCollapsed)
        setTimeout(() => {
          this.sidebarType(!this.sideCollapsed)
        }, 16)
      })
    }
  },
  methods: {
    ...mapMutations({
      sidebarType: 'SIDEBAR_TYPE',
      toggleMobileType: 'TOGGLE_MOBILE_TYPE'
    }),
    handleMediaQuery (val) {
      this.query = val
      if (this.isMobile && !val['screen-xs']) {
        this.isMobile = false
        return
      }
      if (!this.isMobile && val['screen-xs']) {
        this.isMobile = true
        this.sidebarType(false)
        this.settings.contentWidth = 'Fluid'
        // this.settings.fixSiderbar = false
      }
    },
    handleCollapse (val) {
      this.sidebarType(val)
    },
    logoRender () {
      return <Logo></Logo>
    }
  }
}
</script>

<style lang="less">
@import "./BasicLayout.less";
</style>
