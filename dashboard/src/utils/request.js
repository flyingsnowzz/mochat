/**
 * Axios 请求封装
 * 功能说明：配置和管理 HTTP 请求
 * 主要功能：
 * 1. 创建 axios 实例（request 和 newRequest）
 * 2. 配置请求拦截器（添加 Token）
 * 3. 配置响应拦截器（统一错误处理）
 * 4. 导出 axios 实例供组件使用
 *
 * 请求实例说明：
 * - request: 基础请求实例，baseURL 为环境变量配置
 * - newRequest: 固定 baseURL 的请求实例
 *
 * 错误处理：
 * - 401 错误：清除登录状态，跳转到登录页
 * - 其他错误：显示错误信息
 *
 * 技术实现：
 * - 使用 axios.create 创建实例
 * - 使用 interceptors 配置拦截器
 * - Token 存储在 localStorage 中
 */
import axios from 'axios'
import store from '@/store'
import storage from 'store'
import router from '@/router'
import message from 'ant-design-vue/es/message'
import { VueAxios } from './axios'

// 创建 axios 实例
const request = axios.create({
  // API 请求的默认前缀
  baseURL: process.env.VUE_APP_API_BASE_URL + '/dashboard',
  timeout: 15000 // 请求超时时间
})

const newRequest = axios.create({
  // API 请求的默认前缀
  baseURL: '//api.mo.chat',
  timeout: 15000 // 请求超时时间
})

// 异常拦截处理器
const errorHandler = (error) => {
  if (error.response) {
    const data = error.response.data
    const status = error.response.status
    if (status === 401) {
      message.error(data.msg)
      store.dispatch('Logout').then(() => {
        router.push({ path: '/login' })
      })
    } else {
      message.error(`${status || ''}  ${data.msg || 'error'}`)
    }
  } else {
    message.error(error.message || '请求出错，请稍后重试！')
  }
  return Promise.reject(error)
}
const requestInterceptor = (config) => {
  const token = storage.get('ACCESS_TOKEN')
  // 如果 token 存在
  // 让每个请求携带自定义 token 请根据实际情况自行修改
  if (token) {
    config.headers['Accept'] = `application/json`
    config.headers.Authorization = token
  }
  return config
}
// request interceptor
request.interceptors.request.use(requestInterceptor, errorHandler)

// response interceptor
request.interceptors.response.use((response) => {
  return response.data
}, errorHandler)

newRequest.interceptors.request.use(requestInterceptor, errorHandler)
// response interceptor
newRequest.interceptors.response.use((response) => {
  return response.data
}, errorHandler)
const installer = {
  vm: {},
  install (Vue) {
    Vue.use(VueAxios, request)
  }
}

export default request

export {
  installer as VueAxios,
  request as axios,
  newRequest
}
