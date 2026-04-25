<template>
  <div class="work-fission-invite">
    <a-card>
      <div class="block qr-code">
        <div class="title">
          选择客户
        </div>
        <div class="content">
          <div class="option">
            <div class="item">
              <span class="label required">所属员工：</span>
              <a-button @click="selectMemberShow">选择员工</a-button>
              <div class="ml16">
                <a-tag v-for="v in filterSendForm.services" :key="v.id">
                  {{ v.name }}
                </a-tag>
              </div>
            </div>
            <div class="item">
              <span class="label required">选择客服：</span>
              <a-radio-group
                :options="serviceOptions"
                default-value="0"
                v-model="filterValue"
                @change="resetSendCount"
              />
            </div>
            <div class="item">
              <div class="filter" v-if="filterValue === '1'">
                <div class="client-info">
                  <span class="client-title">性别：</span>
                  <a-radio-group
                    :options="genderOptions"
                    default-value="-1"
                    v-model="genderValue"
                    @change="resetSendCount"
                  />
                </div>
                <div class="client-info">
                  <span class="client-title">添加时间：</span>
                  <a-range-picker v-model="filterSendForm.time"/>
                </div>
              </div>
            </div>
            <div class="select-all">
              将群发消息给全部账号的
              <a v-if="sendCount === -1" @click="getSendCount">获取</a>
              <span v-else>
                {{ sendCount }}
              </span>
              个客户
            </div>
          </div>
        </div>
      </div>
      <div class="block">
        <div class="title">
          邀请信息
        </div>
        <div class="welcome-box">
          <div class="content welcome-text">
            <div class="item">
              <span class="label required">邀请文案：</span>
              <div class="content text-1">
                <div class="input">
                  <m-enter-text ref="text" v-model="form.msg.text"/>
                </div>
              </div>
            </div>
            <div class="item">
              <span class="label required">邀请链接：</span>
              <div class="content text-2">
                <div class="link-form">
                  <div class="item">
                    <span class="label required">链接标题：</span>
                    <div class="input">
                      <a-input v-model="form.msg.link.title"/>
                    </div>
                  </div>
                  <div class="item">
                    <span class="label required">链接摘要：</span>
                    <div class="input">
                      <a-input v-model="form.msg.link.desc"/>
                    </div>
                  </div>
                  <div class="item">
                    <span>链接封面：</span>
                    <div class="input">
                      <m-upload v-model="form.msg.link.image"/>
                    </div>
                  </div>
                </div>
              </div>
            </div>
          </div>
          <div class="preview">
            <div class="tips">
              企业微信邀请客户参与预览
            </div>
            <m-preview ref="preview"/>
          </div>
        </div>
      </div>

      <a-button type="primary" size="large" @click="sendInvite">
        通知成员发送
      </a-button>
    </a-card>

    <selectMember ref="selectMember" @change="selectMemberChange"/>
  </div>
</template>

<script>
/**
 * 邀请客户参与活动页面
 * 功能说明：邀请客户参与企微任务宝活动
 * 主要功能：
 * 1. 选择客户（按员工、客服、性别、添加时间筛选）
 * 2. 设置邀请文案
 * 3. 设置邀请链接（标题、描述、封面）
 * 4. 发送邀请消息
 * 5. 获取发送数量
 *
 * 业务场景：
 * - 向符合条件的客户发送活动邀请
 * - 客户收到邀请后可以参与任务宝活动
 *
 * 技术实现：
 * - 使用 selectMember 组件选择员工
 * - 使用 a-radio-group 进行选项选择
 * - 使用 a-range-picker 选择时间范围
 * - 使用 m-enter-text 编辑邀请文案
 * - 调用 chooseContact 和 inviteMsg API
 */
import selectMember from '@/components/Select/member'
// eslint-disable-next-line no-unused-vars
import { chooseContact, inviteMsg } from '@/api/workFission'

export default {
  data () {
    return {
      form: {
        msg: {
          text: '',
          link: {
            title: '',
            desc: '',
            image: ''
          }
        }
      },
      serviceOptions: [
        { label: '全部客户', value: '0' },
        { label: '筛选客户', value: '1' }
      ],
      genderOptions: [
        { label: '全部性别', value: '-1' },
        { label: '仅男性粉丝', value: '1' },
        { label: '仅女性粉丝', value: '2' },
        { label: '未知性别', value: '0' }
      ],
      filterValue: '0',
      genderValue: '-1',
      filterSendForm: {
        time: [],
        services: []
      },
      sendCount: -1
    }
  },
  mounted () {
    const msg = '你好，我们正在进行xxx活动，只要邀请x位好友添加我的微信就可以获得奖品\n' +
      '\n' +
      '参与流程：\n' +
      '①点击下面链接，生成专属海报\n' +
      '②进入链接后长按保存海报，将海报发给好友或朋友圈\n' +
      '③邀请x位好友扫码添加，即可成功获得奖品\n' +
      '④进入链接点击查看进度，完成任务后点击「领取奖励」即可领取哦\n' +
      '\n' +
      '注意事项：请不要直接转发活动链接给好友，是无法成功记录数据的哦~'

    this.$refs.text.addUserName(msg)
    this.form.msg.text = msg

    this.form.msg.link.title = '点击这里，完成任务领取奖品吧👇'
    this.form.msg.link.desc = '快来参加活动吧'
  },
  methods: {
    selectMemberShow () {
      this.$refs.selectMember.setSelect(this.filterSendForm.services)
    },

    sendInvite () {
      const msg = this.getVerify()

      if (msg) {
        this.$message.error(msg)

        return false
      }

      const loading = this.$message.loading('发送中', 1.5)

      inviteMsg(this.getInviteMsgFormData()).then(res => {
        if (res.code === 200) {
          loading.then(() => {
            this.$message.success('发送成功', 2.5)
          }).then(() => {
            this.$router.push('/workFission/taskpage')
          })
        }
      })
    },

    getVerify () {
      if (!this.form.msg.text) return '邀请文案未填写'
      if (!this.form.msg.link.title) return '邀请链接标题未填写'
      if (!this.form.msg.link.desc) return '邀请链接描述未填写'
      if (!this.form.msg.link.image) return '邀请链接封面未上传'
      if (!this.filterSendForm.time.length && this.filterValue !== '0') return '请选择时间'
      if (!this.filterSendForm.services.length) return '请选择员工'
    },

    getInviteMsgFormData () {
      return {
        fission_id: this.$route.query.id,
        text: this.form.msg.text,
        link_title: this.form.msg.link.title,
        link_desc: this.form.msg.link.desc,
        link_pic: this.form.msg.link.image,
        filter: this.getSendFilterParams()
      }
    },

    selectMemberChange (e) {
      this.filterSendForm.services = e
    },

    resetSendCount () {
      this.sendCount = -1
    },

    getSendFilterParams () {
      const params = {}

      params.employee_ids = JSON.stringify(this.filterSendForm.services.map(v => {
        return v.id
      }))

      params.is_all = this.filterValue

      if (this.filterValue === '0') {
        params.start_time = ''
        params.end_time = ''
      } else {
        params.start_time = this.filterSendForm.time[0].format('YYYY-MM-DD hh:ss')
        params.end_time = this.filterSendForm.time[1].format('YYYY-MM-DD hh:ss')
      }

      if (this.genderValue !== '-1') {
        params.gender = this.genderValue
      }

      return params
    },

    getSendCount () {
      if (!this.filterSendForm.time.length && this.filterValue !== '0') {
        this.$message.error('请选择时间')

        return false
      }

      if (!this.filterSendForm.services.length) {
        this.$message.error('请选择员工')

        return false
      }

      const params = this.getSendFilterParams()

      chooseContact(params).then(res => {
        this.sendCount = res.data
      })
    }
  },
  watch: {
    'form.msg': {
      handler () {
        this.$refs.preview.setText(this.form.msg.text)
        this.$refs.preview.setLink(this.form.msg.link.title, this.form.msg.link.desc, this.form.msg.link.image)
      },
      deep: true
    }
  },
  components: { selectMember }
}
</script>

<style lang="less" scoped>
.block {
  margin-bottom: 60px;

  .title {
    font-size: 15px;
    line-height: 21px;
    color: rgba(0, 0, 0, .85);
    border-bottom: 1px solid #e9ebf3;
    padding-bottom: 16px;
    margin-bottom: 16px;
    position: relative;

    span {
      font-size: 13px;
      margin-left: 11px;
      color: rgba(0, 0, 0, .45);
      font-weight: 400;
    }
  }

  .required:after {
    content: "*";
    display: inline-block;
    margin-right: 4px;
    color: #f5222d;
    font-size: 14px;
    line-height: 1;
    position: absolute;
    left: -10px;
    top: 6px;
  }

  .item {
    display: flex;
    align-items: center;
    margin-bottom: 16px;

    .label {
      color: rgba(0, 0, 0, .85);
      position: relative;
    }
  }
}

.welcome-box {
  display: flex;
  align-items: center;

  .preview {
    margin-left: 30px;

    .tips {
      text-align: center;
      margin-bottom: 16px;
    }
  }
}

.welcome-text {
  width: 760px;
  margin-left: 50px;

  .label {
    width: 80px;
  }
}

.account {
  height: 38px;
  border-radius: 4px;
  border: 1px solid #e5e5e5;
  padding-left: 17px;
  padding-right: 11px;
  font-size: 14px;
  color: rgba(0, 0, 0, .85);
  cursor: pointer;
  display: flex;
  align-items: center;
  transition: .3s all;

  .info {
    display: flex;
    align-items: center;
    flex: 1;

    img {
      width: 23px;
      height: 23px;
      margin-right: 4px;
      font-size: 14px;
      border-radius: 2px;
    }
  }

  .icon {
    margin-left: 10px;
  }

  &:hover {
    border: 1px solid #1890ff;
  }
}

.select-all {
  display: inline-block;
  font-size: 13px;
  height: 40px;
  color: rgba(0, 0, 0, .65);
  line-height: 20px;
  padding: 10px 8px 8px;
  background: #effaff;
}

.option {
  margin-left: 40px;
}

.filter {
  width: 850px;
  background: #fbfbfb;
  border-radius: 2px;
  border: 1px solid #ebebeb;
  padding-top: 13px;

  .client-info {
    display: flex;
    margin-bottom: 25px;

    .client-title {
      width: 90px;
      text-align: right;
      margin-right: 15px;
    }
  }

  .tags {
    display: flex;
    align-items: center;
    cursor: pointer;

    .tag {
      height: 27px;
      padding: 0 14px;
      background: #e7f7ff;
      border: 1px solid #1890ff;
      color: #1890ff;
      border-radius: 4px;
      margin-right: 10px;
      margin-bottom: 10px;
    }
  }
}

.text-1 {
  width: 100%;
  border: 1px solid #eee;
  background: #fbfbfb;
  border-radius: 2px;

  .insert-btn-group {
    width: 100%;
    flex: 1;
    border-bottom: 1px dashed #e9e9e9;
    padding: 6px 15px;
    color: #e8971d;
    cursor: pointer;
  }

  .textarea {
    overflow-y: auto;
    overflow-x: hidden;
    white-space: pre-wrap;
    word-break: break-all;

    textarea {
      width: 100%;
      height: 110px;
      padding: 6px 13px;
      border: none;
      background: #fbfbfb;
      outline: none;
      resize: none;
    }

    .word-count {
      font-size: 13px;
      color: rgba(0, 0, 0, .25);
      margin-left: 10px;
    }
  }
}

.text-2 {
  width: 100%;
  border: 1px solid #eee;
  background: #fbfbfb;
  border-radius: 2px;
  padding-left: 20px;
}

.link-form {
  margin-top: 16px;

  .item {
    display: flex;
    align-items: center;
    margin-bottom: 14px;

    .ant-input {
      width: 348px;
    }
  }
}

.instructions-img {
  width: 758px;
}

.mb16 {
  margin-bottom: 16px;
}

/deep/ .ant-alert-description {
  font-size: 13px;
}

/deep/ .ant-alert-with-description {
  padding: 9px 31px 3px 64px !important;
}
</style>
