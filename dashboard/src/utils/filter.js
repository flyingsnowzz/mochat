/**
 * Vue 过滤器
 * 功能说明：提供全局使用的过滤器
 * 主要过滤器：
 * 1. NumberFormat - 数字格式化（添加千分位分隔符）
 * 2. dayjs - 日期格式化（使用 moment.js）
 * 3. moment - 日期格式化（使用 moment.js）
 *
 * 使用方式：
 * - 在模板中使用：{{ value | NumberFormat }}
 * - 在 JS 中使用：this.$options.filters.NumberFormat(value)
 *
 * 技术实现：
 * - 使用 Vue.filter 注册全局过滤器
 * - 使用 moment.js 处理日期格式化
 */
import Vue from 'vue'
import moment from 'moment'
import 'moment/locale/zh-cn'
moment.locale('zh-cn')

Vue.filter('NumberFormat', function (value) {
  if (!value) {
    return '0'
  }
  const intPartFormat = value.toString().replace(/(\d)(?=(?:\d{3})+$)/g, '$1,') // 将整数部分逢三一断
  return intPartFormat
})

Vue.filter('dayjs', function (dataStr, pattern = 'YYYY-MM-DD HH:mm:ss') {
  return moment(dataStr).format(pattern)
})

Vue.filter('moment', function (dataStr, pattern = 'YYYY-MM-DD HH:mm:ss') {
  return moment(dataStr).format(pattern)
})
