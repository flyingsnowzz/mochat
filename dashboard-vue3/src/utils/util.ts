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
export function timeFix (): string {
  const time = new Date()
  const hour = time.getHours()
  return hour < 9 ? '早上好' : hour <= 11 ? '上午好' : hour <= 13 ? '中午好' : hour < 20 ? '下午好' : '晚上好'
}

export function isIE (): boolean {
  const bw = window.navigator.userAgent
  const compare = (s: string): boolean => bw.indexOf(s) >= 0
  const ie11 = (): boolean => 'ActiveXObject' in window
  return compare('MSIE') || ie11()
}

export type ValidateCallback = (error?: Error) => void

export const createValidate = (callback: ValidateCallback, value: any, message: string): void => {
  if (!value) {
    return callback(new Error(message))
  } else {
    callback()
  }
}

export interface ValidatorFunc {
  validator: (rule: any, value: any, callback: ValidateCallback) => void
  trigger: string
}

export const createFunc = (func: (rule: any, value: any, callback: ValidateCallback) => void, change?: string): ValidatorFunc => {
  return {
    validator: func,
    trigger: change || 'blur'
  }
}