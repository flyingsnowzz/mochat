/**
 * 通用工具函数
 * 功能说明：提供常用的工具函数
 * 主要函数：
 * 1. timeFix - 根据当前时间返回问候语（早上好/上午好/中午好/下午好/晚上好）
 * 2. isIE - 检测当前浏览器是否为 IE 浏览器
 * 3. createValidate - 创建表单验证函数
 * 4. createFunc - 创建验证规则配置对象
 *
 * 技术实现：
 * - 基于原生 JavaScript 实现
 * - 用于表单验证的封装
 */
export function timeFix () {
  const time = new Date()
  const hour = time.getHours()
  return hour < 9 ? '早上好' : hour <= 11 ? '上午好' : hour <= 13 ? '中午好' : hour < 20 ? '下午好' : '晚上好'
}

export function isIE () {
  const bw = window.navigator.userAgent
  const compare = (s) => bw.indexOf(s) >= 0
  const ie11 = (() => 'ActiveXObject' in window)()
  return compare('MSIE') || ie11
}
export const createValidate = (callback, value, message) => {
  if (!value) {
    return callback(new Error(message))
  } else {
    callback()
  }
}
export const createFunc = (func, change) => {
  return {
    validator: func,
    trigger: change || 'blur'
  }
}
