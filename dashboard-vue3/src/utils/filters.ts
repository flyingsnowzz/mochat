/**
 * 工具函数 - 替代 Vue 2 过滤器
 * 功能说明：提供常用的格式化函数
 * 主要函数：
 * 1. numberFormat - 数字格式化（添加千分位分隔符）
 * 2. formatDate - 日期格式化
 *
 * 技术实现：
 * - 使用原生 JavaScript 实现
 * - 替代 Vue 2 中的过滤器功能
 *
 * 使用方式：
 * - 在组件中导入并使用：import { numberFormat, formatDate } from '@/utils/filters'
 * - 在模板中使用：{{ numberFormat(value) }}
 */

/**
 * 数字格式化（添加千分位分隔符）
 * @param value 数字值
 * @returns 格式化后的字符串
 */
export function numberFormat(value: number | string): string {
  if (!value) {
    return '0'
  }
  const num = typeof value === 'string' ? parseFloat(value) : value
  if (isNaN(num)) {
    return '0'
  }
  return num.toString().replace(/(\d)(?=(?:\d{3})+$)/g, '$1,')
}

/**
 * 日期格式化
 * @param dateStr 日期字符串或日期对象
 * @param pattern 格式化模式
 * @returns 格式化后的日期字符串
 */
export function formatDate(dateStr: string | Date, pattern: string = 'YYYY-MM-DD HH:mm:ss'): string {
  const date = typeof dateStr === 'string' ? new Date(dateStr) : dateStr
  if (isNaN(date.getTime())) {
    return ''
  }

  const year = date.getFullYear()
  const month = String(date.getMonth() + 1).padStart(2, '0')
  const day = String(date.getDate()).padStart(2, '0')
  const hours = String(date.getHours()).padStart(2, '0')
  const minutes = String(date.getMinutes()).padStart(2, '0')
  const seconds = String(date.getSeconds()).padStart(2, '0')

  return pattern
    .replace('YYYY', String(year))
    .replace('MM', month)
    .replace('DD', day)
    .replace('HH', hours)
    .replace('mm', minutes)
    .replace('ss', seconds)
}

/**
 * 日期格式化（别名，与 formatDate 功能相同）
 * @param dateStr 日期字符串或日期对象
 * @param pattern 格式化模式
 * @returns 格式化后的日期字符串
 */
export function dayjs(dateStr: string | Date, pattern: string = 'YYYY-MM-DD HH:mm:ss'): string {
  return formatDate(dateStr, pattern)
}

/**
 * 日期格式化（别名，与 formatDate 功能相同）
 * @param dateStr 日期字符串或日期对象
 * @param pattern 格式化模式
 * @returns 格式化后的日期字符串
 */
export function moment(dateStr: string | Date, pattern: string = 'YYYY-MM-DD HH:mm:ss'): string {
  return formatDate(dateStr, pattern)
}