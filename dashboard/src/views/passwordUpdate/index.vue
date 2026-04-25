<template>
  <div class="pass-word">
    <a-card>
      <a-form-model
        class="pass-model"
        ref="ruleForm"
        :model="form"
        :label-col="{ span: 5 }"
        :wrapper-col="{ span: 12 }">
        <a-form-model-item label="手机号码：">
          <h3>{{ userPhone }}</h3>
        </a-form-model-item>
        <a-form-model-item label="旧密码：">
          <a-input v-model="form.oldPassword" />
        </a-form-model-item>
        <a-form-model-item label="新密码：">
          <a-input v-model="form.newPassword" />
        </a-form-model-item>
        <a-form-model-item label="确认新密码：">
          <a-input v-model="form.againNewPassword" />
        </a-form-model-item>
        <div class="footer">
          <a-button v-permission="'/passwordUpdate/index@save'" type="primary" @click="updatePassWord">保存</a-button>
        </div>
      </a-form-model>
    </a-card>
  </div>
</template>
<script>
/**
 * 密码修改页面
 * 功能说明：用于用户修改自己的登录密码
 * 主要功能：
 * 1. 显示当前用户的手机号码
 * 2. 输入旧密码
 * 3. 输入新密码
 * 4. 确认新密码
 * 5. 保存密码修改
 *
 * 业务场景：
 * - 用户登录后可以修改自己的登录密码
 * - 修改成功后会自动退出登录，需要重新登录
 *
 * 技术实现：
 * - 使用 a-form-model 表单组件
 * - 从 vuex 中获取用户信息（手机号码）
 * - 调用 passWordUpdate API 修改密码
 * - 修改成功后调用 Logout 退出登录
 */
import { passWordUpdate } from '@/api/passWordUpdate'
import { mapGetters } from 'vuex'
export default {
  computed: {
    ...mapGetters(['userInfo'])
  },
  data () {
    return {
      userPhone: '',
      form: {}
    }
  },
  created () {
    const time = this.userInfo ? 0 : 2000
    setTimeout(() => {
      this.userPhone = this.userInfo.userPhone
    }, time)
  },
  methods: {
    updatePassWord () {
      passWordUpdate(this.form).then(res => {
        this.$store.dispatch('Logout').then(() => {
          this.$router.push({ name: 'login' })
        })
      })
    }
  }
}
</script>
<style lang='less' scoped>
.pass-model {
  width:500px;
  margin: 0 auto;
  .footer {
    width: 23%;
    margin: 0 auto;
    .ant-btn {
      width: 100px;
    }
  }
}

</style>
