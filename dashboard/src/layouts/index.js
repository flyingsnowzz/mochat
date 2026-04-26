/**
 * 布局组件导出模块
 * 功能说明：统一导出所有布局组件
 * 主要组件：
 * 1. BasicLayout - 基础布局（包含侧边栏、头部、内容区域）
 * 2. BlankLayout - 空布局（无侧边栏和头部）
 * 3. RouteView - 路由视图（支持多标签页）
 * 4. PageView - 页面视图（支持多标签页和缓存）
 *
 * 技术实现：
 * - 根据路由 meta.requiresAuth 判断使用哪种布局
 * - 使用 Vue Router 的嵌套路由实现布局嵌套
 */
import BlankLayout from './BlankLayout'
import BasicLayout from './BasicLayout'
import RouteView from './RouteView'
import PageView from './PageView'

export { BasicLayout, BlankLayout, RouteView, PageView }
