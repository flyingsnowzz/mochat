<template>
  <div class="work-fission-create">
    <div class="steps mb20">
      <a-card>
        <a-steps :current="step">
          <a-step title="活动基本信息"/>
          <a-step title="设置欢迎语和海报"/>
          <a-step title="设置推送消息"/>
          <a-step title="邀请客户参与"/>
        </a-steps>
      </a-card>
    </div>
    <a-card>
      <div class="mb20">
        <step1 ref="step1" v-show="step === 0"></step1>
        <step2 ref="step2" v-show="step === 1"></step2>
        <step3 ref="step3" v-show="step === 2"></step3>
        <step4 ref="step4" v-show="step === 3"></step4>
      </div>
      <a-button
        class="mr"
        size="large"
        v-if="step !== 0"
        @click="step--"
      >
        上一步
      </a-button>
      <a-button
        type="primary"
        size="large"
        v-if="step !== 3"
        @click="next">
        下一步
      </a-button>
      <a-button
        type="primary"
        size="large"
        v-if="step === 3"
        @click="add"
        :loading="loading"
      >
        完成
      </a-button>
    </a-card>
  </div>
</template>

<script>
/**
 * 工作裂变活动创建页面
 * 功能说明：创建企业微信的工作裂变活动，通过多步骤向导完成活动设置
 * 主要功能：
 * 1. 多步骤向导：
 *    - 步骤1：活动基本信息（活动名称、使用成员、任务设置、奖励设置等）
 *    - 步骤2：设置欢迎语和海报（分享文案、海报设计、欢迎语内容）
 *    - 步骤3：设置推送消息（推送对象、消息内容）
 *    - 步骤4：邀请客户参与（邀请消息设置）
 * 2. 表单验证
 * 3. 活动数据提交
 *
 * 业务场景：
 * - 企业创建裂变活动，鼓励员工邀请客户参与
 * - 通过分享海报和消息推送扩大活动影响力
 * - 设置奖励机制提高参与积极性
 *
 * 技术实现：
 * - 使用 a-steps 组件实现多步骤向导
 * - 使用多个 step 子组件处理不同步骤的表单
 * - 分步收集活动数据并提交
 */
import step1 from './components/create/step/step1'
import step2 from './components/create/step/step2'
import step3 from './components/create/step/step3'
import step4 from './components/create/step/step4'

import { add, inviteMsg } from '@/api/workFission'

export default {
  data () {
    return {
      step: 0,
      loading: false
    }
  },
  methods: {
    getFissionData () {
      const data = this.$refs.step1.getFormData()

      const params = {
        active_name: data.name,
        service_employees: data.services,
        auto_pass: data.autoPass,
        auto_add_tag: data.autoAddTag,
        end_time: data.endTime,
        qr_code_invalid: data.expireDay,
        tasks: data.tasks,
        new_friend: data.newFriend,
        delete_invalid: data.deleteInvalid,
        receive_prize: data.receivePrizeType,
        receive_prize_employees: data.receivePrizeServices,
        receive_links: data.receivePrizeLinks,
        contact_tags: []
      }

      return params
    },

    getPosterData () {
      const data = this.$refs.step2.getFormData()

      const params = {
        poster_type: data.posterType,
        foward_text: data.shareText,
        avatar_show: data.poster.avatar,
        nickname_show: data.poster.nickname,
        nickname_color: data.poster.nicknameColor,
        card_corp_image_name: data.card.nickname,
        card_corp_name: data.card.name,
        card_corp_logo: data.card.logoUrl,
        cover_pic: data.poster.imageUrl,
        qrcode_w: data.poster.imageW,
        qrcode_h: data.poster.imageH,
        qrcode_x: data.poster.imageX,
        qrcode_y: data.poster.imageY
      }

      return params
    },

    getWelcomeData () {
      const data = this.$refs.step2.getFormData()

      const params = {
        msg_text: data.welcome.text,
        link_title: data.welcome.link.title,
        link_desc: data.welcome.link.desc,
        link_cover_url: data.welcome.link.imageUrl
      }

      return params
    },

    getPushData () {
      const data = this.$refs.step3.getFormData()

      const params = {
        push_employee: data.pushEmployee,
        push_contact: data.pushContact,
        msg_text: data.msgText,
        msg_complex: {
          ...data.msgComplex,
          msg_complex_type: data.msgComplexType
        }
      }

      return params
    },

    getInviteData () {
      const data = this.$refs.step4.getFormData()

      return {
        text: data.msg.text,
        link_title: data.msg.link.title,
        link_desc: data.msg.link.desc,
        link_pic: data.msg.link.image
      }
    },

    add () {
      const err = this.$refs.step4.getVerify()

      if (err) {
        this.$message.error(err)

        return false
      }
      const params = {
        fission: this.getFissionData(),
        poster: this.getPosterData(),
        welcome: this.getWelcomeData(),
        push: this.getPushData(),
        invite: this.getInviteData()
      }
      this.loading = true
      add(params).then(res => {
        if (res.code === 200) {
          const inviteParams = this.$refs.step4.getInviteMsgFormData()
          if (inviteParams) {
            inviteParams.fission_id = res.data[0]
            inviteMsg(inviteParams).then(res => {
              if (res.code === 200) {
                this.$message.success('添加成功')
                this.loading = false
                this.$router.push('/workFission/taskpage')
              }
            })
            return false
          }
          this.$message.success('添加成功')
          this.loading = false
          this.$router.push('/workFission/taskpage')
        }
      })
    },

    next () {
      let err = ''
      if (this.step === 0) err = this.$refs.step1.getVerify()
      if (this.step === 1) err = this.$refs.step2.getVerify()
      if (this.step === 2) err = this.$refs.step3.getVerify()
      if (err) {
        this.$message.error(err)
        err = ''

        return false
      }

      this.step++
    }
  },
  components: { step1, step2, step3, step4 }
}
</script>

<style lang="less" scoped>
.mr {
  margin-right: 16px;
}
</style>
