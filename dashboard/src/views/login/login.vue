<!--
/**
 * 登录页面组件
 *
 * 功能说明：
 * 1. 这是MoChat系统的用户登录页面，是系统的入口点
 * 2. 支持两种登录方式：
 *    - tab1: 手机号 + 密码登录（主要方式）
 *    - tab2: 手机号 + 验证码 + 新密码登录（用于忘记密码后重置）
 * 3. 登录成功后会自动获取用户权限列表并跳转到首页
 * 4. 页面底部显示MoChat的版权信息
 *
 * 使用技术：
 * - Vue 2.x
 * - Ant Design Vue 组件库
 * - Vuex 状态管理
 * - Vue Router 路由管理
 */
-->
<template>
  <!-- 登录页面容器 - 全屏居中显示，带有背景图片 -->
  <div class="login-wrapper">
    <!-- 登录表单容器 -->
    <div class="contant">
      <!-- MoChat Logo 图片 -->
      <div class="img-wrapper">
        <img class="img" :src="require('@/assets/title.png')" alt="MoChat Logo">
      </div>

      <!-- 登录表单 - 使用 Ant Design Vue 的 Form 组件 -->
      <a-form
        id="formLogin"
        class="user-layout-login"
        ref="formLogin"
        :form="form"
        @submit="handleSubmit"
      >
        <!-- tab1: 手机号 + 密码登录表单 -->
        <div v-if="customActiveKey == 'tab1'">
          <!-- 手机号输入框 -->
          <a-form-item>
            <a-input
              class="input"
              size="large"
              type="text"
              placeholder="手机号"
              v-decorator="[
                'phone',
                {rules: [{ required: true, pattern: /^1[345789]\d{9}$/, message: '请输入正确的手机号' }], validateTrigger: 'change'}
              ]"
            >
              <!-- 手机号输入框前缀图标 -->
              <img slot="prefix" :src="require('@/assets/user.png')" alt="用户图标">
            </a-input>
          </a-form-item>

          <!-- 密码输入框 -->
          <a-form-item>
            <a-input-password
              class="input"
              size="large"
              placeholder="密码"
              v-decorator="[
                'password',
                {rules: [{ required: true, message: '请输入密码' }], validateTrigger: 'blur'}
              ]"
            >
              <!-- 密码输入框前缀图标 -->
              <img slot="prefix" :src="require('@/assets/lock.png')" alt="密码图标">
            </a-input-password>
          </a-form-item>

          <!-- 登录错误提示信息 -->
          <a-alert v-if="isLoginError" type="error" show-icon message="登录失败" />
        </div>

        <!-- tab2: 手机号 + 验证码 + 新密码登录表单（忘记密码功能） -->
        <div v-if="customActiveKey == 'tab2'">
          <!-- 手机号输入框 -->
          <a-form-item>
            <a-input
              class="input"
              size="large"
              type="text"
              placeholder="手机号"
              v-decorator="['mobile', {rules: [{ required: true, pattern: /^1[345789]\d{9}$/, message: '请输入正确的手机号' }], validateTrigger: 'change'}]"
            >
              <a-icon slot="prefix" type="mobile" :style="{ color: 'rgba(0,0,0,.25)' }" />
            </a-input>
          </a-form-item>

          <!-- 验证码输入框 + 获取验证码按钮 -->
          <a-row :gutter="16">
            <a-col class="gutter-row" :span="16">
              <a-form-item>
                <a-input
                  class="input"
                  size="large"
                  type="text"
                  placeholder="验证码"
                  v-decorator="['captcha', {rules: [{ required: true, message: '请输入验证码' }], validateTrigger: 'blur'}]"
                >
                  <a-icon slot="prefix" type="mail" :style="{ color: 'rgba(0,0,0,.25)' }" />
                </a-input>
              </a-form-item>
            </a-col>
            <a-col class="gutter-row" :span="8">
              <!-- 获取验证码按钮 - 倒计时功能 -->
              <a-button
                class="getCaptcha"
                tabindex="-1"
                :disabled="state.smsSendBtn"
                @click.stop.prevent="getCaptcha"
                v-text="!state.smsSendBtn && '获取验证码' || (state.time+' s')"
              />
            </a-col>
          </a-row>

          <!-- 新密码输入框 -->
          <a-input
            class="input"
            size="large"
            placeholder="设置一个新的密码"
            v-decorator="[
              'newPassword',
              {rules: [{ required: true, message: '请输入密码' }], validateTrigger: 'blur'}
            ]"
          >
            <a-icon slot="prefix" type="lock" :style="{ color: 'rgba(0,0,0,.25)' }" />
          </a-input>
        </div>

        <!-- 登录/确定按钮 -->
        <a-form-item>
          <!-- tab1 的登录按钮 -->
          <div class="btn-wrapper" v-if="customActiveKey == 'tab1'">
            <a-button
              type="primary"
              html-type="submit"
              class="login-button"
              :loading="state.loginBtn"
              :disabled="state.loginBtn"
            >
              登录
            </a-button>
            <!-- 忘记密码按钮（暂时隐藏） -->
            <a-button
              v-if="false"
              class="login-button"
              @click="handleTabClick"
              :loading="state.loginBtn"
              :disabled="state.loginBtn">
              忘记密码
            </a-button>
          </div>

          <!-- tab2 的确定按钮 -->
          <div class="confirm-wrapper" v-if="customActiveKey == 'tab2'">
            <a-button
              type="primary"
              html-type="submit"
              class="login-button"
              :loading="state.loginBtn"
              :disabled="state.loginBtn"
            >
              确定
            </a-button>
          </div>
        </a-form-item>
      </a-form>
    </div>

    <!-- 页面底部版权信息 -->
    <div class="footer">
      Powered by <a class="mochat" href="https://mo.chat/" target="_blank">MoChat</a>
    </div>
  </div>
</template>

<script>
/**
 * 登录页面脚本部分
 *
 * 主要功能：
 * 1. 处理用户登录请求
 * 2. 表单验证
 * 3. 登录成功后的权限获取和路由跳转
 * 4. 忘记密码功能的验证码倒计时
 */
import { mapActions } from 'vuex'
import store from '@/store'
import { resetRoutes } from '@/utils/menu'

export default {
  data () {
    return {
      // 当前激活的标签页：'tab1' 为手机号+密码登录，'tab2' 为手机号+验证码+新密码
      customActiveKey: 'tab1',
      // 登录按钮loading状态
      loginBtn: false,
      // 登录类型：0-邮箱，1-手机号，2-电话号码
      loginType: 0,
      // 是否显示登录错误提示
      isLoginError: false,
      // 表单实例，用于获取表单值和进行表单验证
      form: this.$form.createForm(this),
      // 按钮和倒计时状态管理
      state: {
        time: 60, // 验证码倒计时时间（秒）
        loginBtn: false, // 登录按钮loading状态
        loginType: 0, // 登录类型
        smsSendBtn: false // 验证码发送按钮状态
      }
    }
  },

  created () {
    // 组件创建时的生命周期钩子
    // 可以在这里进行初始数据的加载
  },

  methods: {
    // 从 Vuex 映射的登录和登出方法
    ...mapActions(['Login', 'Logout']),

    /**
     * 切换到忘记密码标签页
     * 目前此功能被隐藏，tab2 暂时不可用
     */
    handleTabClick () {
      this.customActiveKey = 'tab2'
    },

    /**
     * 处理表单提交事件
     * 根据当前激活的标签页验证不同的字段
     * @param {Event} e - 表单提交事件
     */
    handleSubmit (e) {
      e.preventDefault()
      const {
        form: { validateFields },
        state,
        customActiveKey,
        Login
      } = this

      // 开始登录按钮的loading状态
      state.loginBtn = true

      // 根据标签页确定需要验证的字段
      // tab1: 手机号 + 密码
      // tab2: 手机号 + 验证码 + 新密码
      const validateFieldsKey = customActiveKey === 'tab1' ? ['phone', 'password'] : ['mobile', 'captcha', 'newPassword']

      // 验证表单字段
      validateFields(validateFieldsKey, { force: true }, (err, values) => {
        if (!err) {
          // 表单验证通过，准备登录参数
          const loginParams = { ...values }
          // 调用登录接口
          Login(loginParams)
            .then((res) => {
              this.loginSuccess(res)
            })
            .catch(err => this.requestFailed(err))
            .finally(() => {
              // 无论成功或失败，都重置登录按钮状态
              state.loginBtn = false
            })
        } else {
          // 表单验证失败，600ms后重置登录按钮状态
          setTimeout(() => {
            state.loginBtn = false
          }, 600)
        }
      })
    },

    /**
     * 获取验证码
     * 点击后开始60秒倒计时
     * @param {Event} e - 点击事件
     */
    getCaptcha (e) {
      e.preventDefault()
      const { form: { validateFields }, state } = this

      // 验证手机号字段
      validateFields(['mobile'], { force: true }, (err, values) => {
        if (!err) {
          // 手机号验证通过，调用获取验证码接口
          this.loginSuccess()
          state.smsSendBtn = true

          // 开始倒计时
          const interval = window.setInterval(() => {
            if (state.time-- <= 0) {
              // 倒计时结束，重置状态
              state.time = 60
              state.smsSendBtn = false
              window.clearInterval(interval)
            }
          }, 1000)
        }
      })
    },

    /**
     * 登录成功后的处理
     * 1. 获取用户权限列表
     * 2. 重置路由
     * 3. 跳转到首页
     * @param {Object} res - 登录接口返回的数据
     */
    async loginSuccess (res) {
      // 从服务器获取用户的权限列表
      const data = await store.dispatch('getPermissionList')
      if (data) {
        // 权限获取成功，重置路由并跳转到首页
        resetRoutes()
        this.$router.push({ path: '/' })
        this.isLoginError = false
      } else {
        // 权限获取失败，执行登出操作
        store.dispatch('Logout')
      }
    },

    /**
     * 登录请求失败的处理
     * @param {Object} err - 错误信息
     */
    requestFailed (err) {
      // 显示登录错误提示
      this.isLoginError = true
      // 在控制台输出错误信息，便于调试
      console.log(err)
    }
  }
}
</script>

<style lang="less" scoped>
/**
 * 登录页面样式
 * 采用深色主题，蓝色渐变背景
 */

/* 登录页面容器 - 全屏布局，元素居中 */
.login-wrapper {
  display: flex;
  height: 100%;
  align-items: center;
  justify-content: center;
  flex-direction: column;
  /* 设置背景图片 */
  background: url('../../assets/background.png') no-repeat center;

  /* Logo 容器样式 */
  .img-wrapper {
    width: 142px;
    margin: 0 auto;
    padding: 42px 0 45px 0;
    .img {
      width: 142px;
    }
  }

  /* 登录表单容器 - 深色背景，圆角，阴影效果 */
  .contant {
    width: 370px;
    height: 430px;
    margin-bottom: 25px;
    background: #293152; // 深蓝色背景
    border-radius: 10px; // 圆角
    box-shadow: 0px 11px 24px 0px rgba(18,21,40,0.58) // 阴影效果
  }

  /* 输入框样式 */
  .input {
    height: 60px; // 输入框高度
    font-size: 16px; // 字体大小
  }

  /* 底部版权信息样式 */
  .footer {
    color: #e4eafa;
    font-size: 18px;
    font-weight: 400;
    margin-top: 80px;
  }
}

/* 用户登录表单样式 */
.user-layout-login {
  padding: 0 21px;

  label {
    font-size: 14px;
  }

  /* 获取验证码按钮样式 */
  .getCaptcha {
    display: block;
    width: 100%;
    min-width: 70px;
    height: 40px;
  }

  .forge-password {
    font-size: 14px;
  }

  /* 登录按钮容器 */
  .btn-wrapper{
    width: 100%;
    margin-top: 25px;
    display: flex;
    justify-content: center;
  }

  /* 确定按钮容器 */
  .confirm-wrapper {
    width: 100%;
    margin-top: 40px;
    display: flex;
    justify-content: center;
  }

  /* 登录/确定按钮样式 */
  button.login-button {
    padding: 0 15px;
    font-size: 20px;
    width: 100%;
    height: 60px;
    font-weight: 500;
  }
}
</style>
