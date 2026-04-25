<template>
  <div class="step3">
    <div class="block qr-code">
      <div class="title">
        企微群聊二维码设置
      </div>
      <div class="content">
        <a-alert
          class="mb16"
          type="warning"
          show-icon>
          <p slot="description">
            1、注意群聊二维码有效期为7天，需要及时更换哦，以防客户进群失效哦～<br>
            2、上传群二维码请 仔细核对 ，群聊二维码和群聊 对应错误 将导致 客户进群失效哦~<br>
            3、建议同一群聊不要应用于多个群裂变活动，可能导致群上限控制精度降低，出现客户无法正常加入群聊，或数据计算异常的情况"
          </p>
        </a-alert>
        <div class="item">
          <span class="label required">群聊二维码：</span>
          <div class="input">
            <a-alert class="mb16" message="若群内添加了小助理机器人，请将【人数上限减1】，可能会导致客户进群失败" type="info" show-icon/>
            <div class="groups">
              <div class="g-item" v-for="(item,index) in rooms" :key="index">
                <div class="g-title">群聊二维码{{ index+1 }}</div>
                <div class="group">
                  <!--                  群聊二维码-->
                  <div class="upload mb16">
                    <!--    选择二维码-->
                    <!--                    v-model="item.room_qrcode"  -->
                    <m-upload :def="false" text="请上传二维码" :ref="`qrCode${index}`" @change="(e)=>{item.room_qrcode=e}" />
                  </div>
                  <div class="select">
                    <div class="select-group" @click="electionGroup(index)" v-if="JSON.stringify(item.room)=='{}'">选择群聊</div>
                    <div class="group_style" v-else>
                      <img src="../../../../../assets/avatar-room-default.svg">
                      <div class="group_News">
                        <div class="name">{{ item.room.name }}</div>
                        <div class="num">{{ item.room.contact_num }}/{{ item.room.roomMax }}</div>
                      </div>
                      <a-icon type="close-circle" @click="delgroup(index)" />
                    </div>
                    <div class="count">群人数上限：
                      <a-input-number v-model="item.room_max" :min="1" style="width: 72px"></a-input-number>
                    </div>
                  </div>
                  <div class="box">
                    <div class="icon" @click="delGroupNews(index)">
                      <a-icon type="rest"/>
                    </div>
                  </div>
                </div>
              </div>
              <div class="add-group" @click="addGroupNews">
                <a-icon type="plus-circle"/>
                <div>
                  添加群聊
                </div>
              </div>
            </div>
          </div>
        </div>
        <!--        群聊分配规则-->
        <div class="item">
          <span class="label required">群聊分配规则：</span>
          <div class="input">
            <a-radio >顺序分配</a-radio>
            <span>当前面的群人数达到上限后，自动发送后面的群二维码</span>
          </div>
        </div>
      </div>
    </div>
    <div class="block">
      <div class="title">入群欢迎语素材设置</div>
      <div class="welcome-box">
        <div class="content welcome-text">
          <a-alert
            class="mb16"
            type="info"
            show-icon>
            <p slot="description">
              1、 设置入群欢迎语素材目的是为了邀请客户参与活动，是整个裂变活动中重要的步骤<br>
              2、活动创建成功后，欢迎语链接将自动填入裂变活动链接<br>
              3、群欢迎语素材设置完成后，请提示对应的群主/群管理员前往移动端活动群聊设置入群欢迎语，客户进群才会收到参与活动链接哦～ 如何设置
            </p>
          </a-alert>
          <div class="item">
            <span class="label required">欢迎语1：</span>
            <div class="content text-1">
              <div class="input">
                <div class="insert-btn-group">
                  <span>[插入客户名称]</span>
                </div>
                <m-enter-text @change="receiveText" v-if="!banModify"/>
                <div class="textarea" v-else>
                  <textarea v-model="text" :disabled="banModify"></textarea>
                </div>
              </div>
            </div>
          </div>
          <div class="item">
            <span class="label required">欢迎语2：</span>
            <div class="content text-2">
              <div class="link-card">
                <div class="title">{{ link_title }}</div>
                <div class="desc">
                  <span>{{ link_desc }}</span>
                  <img src="../../../../../assets/default-cover.png" v-if="link_pic==''">
                  <img :src="link_pic" v-else>
                </div>
              </div>
              <div class="link-form">
                <div class="item">
                  <span class="label required">链接标题：</span>
                  <div class="input">
                    <a-input v-model="link_title" :disabled="banModify"/>
                  </div>
                </div>
                <div class="item">
                  <span class="label required">链接摘要：</span>
                  <div class="input">
                    <a-input v-model="link_desc" :disabled="banModify"/>
                  </div>
                </div>
                <div class="item">
                  <span>链接封面：</span>
                  <div class="input">
                    <m-upload :def="false" text="请上传封面" v-model="link_pic" v-show="!banModify" ref="coverImg"/>
                    <img :src="link_pic" alt="" v-if="banModify" style="width: 130px;height: 130px;">
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>
        <div class="preview">
          <div class="tips">
            群内欢迎语预览
          </div>
          <m-preview ref="preview"/>
        </div>
      </div>
    </div>
    <!--    选择群聊-->
    <selectGroup ref="selectGroup" @change="acceptData"/>
  </div>
</template>

<script>
/**
 * 群裂变创建向导 - 步骤3：企微群聊二维码和入群欢迎语设置
 * 功能说明：配置裂变活动使用的企业微信群二维码和欢迎语
 * 主要功能：
 * 1. 上传企业微信群聊二维码
 * 2. 选择群聊并设置群人数上限
 * 3. 添加多个群聊（按顺序分配客户）
 * 4. 设置入群欢迎语素材
 *
 * 业务场景：
 * - 客户扫码后会被分配到企业微信群
 * - 多个群时按顺序分配，当前群满员后自动切换到下一个群
 * - 入群欢迎语用于客户入群后的自动消息推送
 *
 * 技术实现：
 * - 使用 selectGroup 组件选择群聊
 * - 使用 qrCode 组件上传群二维码
 * - 使用 m-upload 组件上传欢迎语素材
 */
import selectGroup from '@/components/Select/selectGroup'
import qrCode from '@/components/Upload/qrCode'

export default {
  components: { selectGroup, qrCode },
  data () {
    return {
      // 禁止修改
      banModify: false,
      text: '',
      link_title: '点击这里，完成任务领取奖品吧👇',
      link_desc: '快来参加活动吧',
      link_pic: '',
      // 群聊索引
      indexGroup: 0,
      imgLink: 0,
      // 是否新建群聊信息
      newbulidNews: false,
      // 群聊
      rooms: [
        {
          // 二维码
          room_qrcode: '',
          // 群上限
          room_max: 1,
          // 群聊
          room: {}
        }
      ],
      //  欢迎语
      welcome: {
      }
    }
  },
  mounted () {
    this.$refs.preview.setLink(this.link_title, this.link_desc, '')
  },
  watch: {
    link_title: function (val) {
      if (this.link_pic != '') {
        this.$refs.preview.setLink(this.link_title, this.link_desc, this.link_pic)
      } else {
        this.$refs.preview.setLink(this.link_title, this.link_desc, '')
      }
    },
    link_desc: function (val) {
      if (this.link_pic != '') {
        this.$refs.preview.setLink(this.link_title, this.link_desc, this.link_pic)
      } else {
        this.$refs.preview.setLink(this.link_title, this.link_desc, '')
      }
    },
    link_pic: function (val) {
      this.$refs.preview.setLink(this.link_title, this.link_desc, this.link_pic)
    }
  },
  methods: {
    // 父组件的传值
    parentNews (data) {
      if (this.banModify == false) {
        this.banModify = true
        this.rooms = []
        data.rooms.forEach((item, index) => {
          const roomsData = {
            // 二维码
            room_qrcode: '',
            // 群上限
            room_max: 1,
            // 群聊
            room: {}
          }
          roomsData.room_max = item.roomMax
          roomsData.room = item.room
          this.rooms.push(roomsData)
          this.$nextTick(() => {
            this.$refs[`qrCode${index}`][0].setUrl(item.roomQrcode)
          })
        })
        this.text = data.welcome.text
        this.link_title = data.welcome.linkTitle
        this.link_desc = data.welcome.linkDesc
        this.$refs.coverImg.setUrl(data.welcome.linkPic)
      }
    },
    // 向父组件抛出
    outputStep3 () {
      // 验证群聊二维码
      for (let i = 0; i < this.rooms.length; i++) {
        if (this.rooms[i].room_qrcode == '') {
          this.$message.error('群聊二维码' + (i + 1) + '二维码不能为空')
          return false
        }
        if (this.rooms[i].room_max == 1) {
          this.$message.error('群聊二维码' + (i + 1) + '群上限不能为1')
          return false
        }
        if (JSON.stringify(this.rooms[i].room) == '{}') {
          this.$message.error('请为群聊二维码' + (i + 1) + '选择群聊')
          return false
        }
      }
      // 欢迎语1
      if (this.text == '') {
        this.$message.error('欢迎语1不能为空')
        return false
      }
      this.welcome.text = this.text
      // 欢迎语二
      if (this.link_title == '') {
        this.$message.error('链接标题不能为空')
        return false
      }
      this.welcome.link_title = this.link_title
      if (this.link_desc == '') {
        this.$message.error('链接摘要不能为空')
        return false
      }
      this.welcome.link_desc = this.link_desc

      this.welcome.link_pic = this.link_pic
      const setp3News = {
        rooms: this.rooms,
        welcome: this.welcome
      }
      return setp3News
    },
    // 欢迎语1
    receiveText  (e) {
      this.text = e
      this.$refs.preview.setText(this.text)
    },
    // 添加群组
    addGroupNews () {
      this.$refs.selectGroup.show()
      this.newbulidNews = true
    },
    // 删除群组
    delgroup (index) {
      this.rooms[index].room = {}
    },
    // 删除群组信息
    delGroupNews (index) {
      if (this.rooms.length == 1) {
        this.$message.warning('最后一个不能删除')
      } else {
        this.rooms.splice(index, 1)
      }
    },
    // 选择群聊
    electionGroup (index) {
      this.$refs.selectGroup.show()
      this.indexGroup = index
    },
    // 接收组件传值
    acceptData (e) {
      for (let i = 0; i < this.rooms.length; i++) {
        if (this.rooms[i].room.id == e[0].id) {
          this.$message.warning('不能选中相同的群聊')
          return false
        }
      }
      if (this.newbulidNews) {
      //  添加群组
        const groupNew = {
          // 二维码
          room_qrcode: '',
          // 群上限
          room_max: 1,
          // 群聊
          room: {}
        }
        this.indexGroup = this.rooms.length
        this.rooms.push(groupNew)
        this.newbulidNews = false
      }
      this.rooms[this.indexGroup].room = e[0]
    }

  }
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

.qr-code {
  width: 850px;
}

.groups {
  display: -webkit-box;
  display: flex;
  flex-wrap: wrap;
  width: 766px;
  height: 100%;
  background: #fbfbfb;
  border: 1px solid hsla(0, 0%, 59.2%, .08);
  padding: 16px 8px;
  max-height: 666px;
  overflow-y: scroll;

  .g-item {
    margin-bottom: 50px;
    margin-right: 10px;
  }

  .g-title {
    margin-bottom: 10px;
  }

  .group {
    width: 350px;
    background: #fff;
    padding: 12px;
    border: 1px solid hsla(0, 0%, 59.2%, .08);
    margin-top: 8px;
    display: flex;
    flex-wrap: wrap;

    .move {
      width: 100%;
      border-top: 1px solid #ededed;
      padding-top: 8px;
      display: flex;
      align-items: center;
      justify-content: center;
    }
    .box {
      flex: 1;
      display: flex;
      align-items: center;
      justify-content: center;
      .icon {
        width: 30px;
        height: 30px;
        border: 1px solid #ccc;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        cursor: pointer;
      }
    }
    .select-group {
      display: flex;
      justify-content: center;
      align-items: center;
      width: 144px;
      height: 48px;
      background: #fff;
      border-radius: 2px;
      border: 1px solid #e8e8eb;
      font-size: 14px;
      color: rgba(0, 0, 0, .65);
      cursor: pointer;

      .icon img {
        width: 32px;
        height: 32px;
        margin-right: 8px;
      }

      .info {
        .name {
          color: #222;
          font-size: 14px;
          width: 100%;
          overflow: hidden;
          text-overflow: ellipsis;
          white-space: nowrap;
        }

        .count {
          color: #999;
          opacity: .85;
          font-size: 13px;
          margin-top: 0;
        }
      }
    }

    .count {
      display: flex;
      align-items: center;
      margin-top: 20px;

      input {
        width: 50px;
      }
    }
  }
}
.group_style{
  max-width: 144px;
  max-height: 48px;
  height: 60px;
  background: #fff;
  border: 1px solid #e6e6e6;
  border-radius: 1px;
  display: flex;
  position: relative;
  padding: 8px;

  img {
    width: 30px;
    height: 30px;
  }

  i{
    position: absolute;
    right: -8px;
    top: -8px;
    font-size: 16px;
    color: #B5B5B5;
    cursor: pointer;
  }
}
.group_News{
  margin-left: 7px;
  .name{
    width: 95px;
    color: #222;
    font-size: 13px;
    overflow: hidden;
    text-overflow:ellipsis;
    white-space: nowrap;
    margin-top: -5px;
  }
  .num{
    color: #999;
    opacity: .85;
    font-size: 12px;
  }
}
.welcome-text {
  width: 760px;

  .label {
    width: 70px;
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

.link-card {
  width: 250px;
  height: 120px;
  background: #fff;
  border: 1px solid #f0f0f0;
  padding: 13px;
  margin-left: 16px;
  margin-top: 16px;

  .title {
    width: 100%;
    font-size: 14px;
    text-align: left;
    max-height: 38px;
    overflow: hidden;
    padding-bottom: 0;
    border: none;
    margin: 0 0 25px;
  }

  .desc {
    max-height: 54px;
    display: flex;

    span {
      flex: 1;
      overflow: hidden;
      margin-right: 5px;
      word-break: break-all;
      color: rgba(0, 0, 0, .45);
    }

    img {
      width: 40px;
      height: 40px;
    }
  }
}

.add-group {
  display: flex;
  justify-content: center;
  align-items: center;
  width: 274px;
  height: 140px;
  background: #fff;
  border-radius: 2px;
  border: 1px dashed #ddd;
  color: #000;
  margin-top: 31px;
  flex-direction: column;
  cursor: pointer;
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
.textarea {
  overflow-y: auto;
  overflow-x: hidden;
  white-space: pre-wrap;
  word-break: break-all;
  textarea {
    width: 100%;
    height: 175px;
    padding: 6px 13px;
    border: none;
    background: #fbfbfb;
    outline: none;
    resize: none;
  }
}
</style>
