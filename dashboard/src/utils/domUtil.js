/**
 * DOM 操作工具函数
 * 功能说明：提供 DOM 相关的工具函数
 * 主要函数：
 * 1. setDocumentTitle - 设置页面标题
 *
 * 特殊处理：
 * - 在微信 iOS 浏览器中，页面标题需要通过 iframe 才能动态更新
 * - 这是因为微信 iOS 浏览器对 document.title 的修改有延迟
 *
 * 技术实现：
 * - 创建隐藏的 iframe 来触发标题更新
 * - 加载完成后自动移除 iframe
 */
export const setDocumentTitle = function (title) {
  document.title = title
  const ua = navigator.userAgent
  // eslint-disable-next-line
  const regex = /\bMicroMessenger\/([\d\.]+)/
  if (regex.test(ua) && /ip(hone|od|ad)/i.test(ua)) {
    const i = document.createElement('iframe')
    i.src = '/favicon.ico'
    i.style.display = 'none'
    i.onload = function () {
      setTimeout(function () {
        i.remove()
      }, 9)
    }
    document.body.appendChild(i)
  }
}
