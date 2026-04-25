<template>
  <div class="work-employee">
    <a-card>
      <div class="btn-box">
        <a-button v-permission="'/workEmployee/index@search'" type="primary" @click="() => {this.visible = true}">条件筛选</a-button>
        <div>
          <p>最后一次同步时间：{{ syncTime }}</p>
          <a-button v-permission="'/workEmployee/index@sync'" @click="() => { this.syncEmployeeDefine() }">同步企业微信通讯录</a-button>
        </div>
      </div>
      <div class="table-wrapper">
        <a-table
          :columns="columns"
          :data-source="tableData"
          :rowKey="record => record.id"
          :pagination="pagination"
          @change="handleTableChange">
          <div class="msg-box" slot="name" slot-scope="text, record">
            <img v-if="record.thumbAvatar" :src="record.thumbAvatar" :onerror="errorImg" style="width:40px; height:40px"/>
            <a-icon v-if="!record.thumbAvatar" type="user" style="font-size:35px "/>
            <span style="marginRight: 10px">{{ record.name }}</span>
            <a-icon v-if="record.gender === '男'" type="man" />
            <a-icon v-if="record.gender === '女'" type="woman" />
          </div>
        </a-table>
      </div>
      <a-drawer
        title="条件筛选"
        placement="right"
        width="400"
        :closable="false"
        :visible="visible"
        :after-visible-change="afterVisibleChange"
        @close="onClose"
      >
        <a-form>
          <a-form-item label="选择成员">
            <a-input v-model="searchData.name"/>
          </a-form-item>
          <a-form-item label="成功状态">
            <a-select v-model="searchData.status">
              <a-select-option v-for="item in statusList" :key="item.id">
                {{ item.name }}
              </a-select-option>
            </a-select>
          </a-form-item>
          <a-form-item label="外部联系人权限">
            <a-select v-model="searchData.contactAuth">
              <a-select-option v-for="item in contactAuth" :key="item.id">
                {{ item.name }}
              </a-select-option>
            </a-select>
          </a-form-item>
        </a-form>
        <div class="footer">
          <a-button style="marginRight: 5px" @click="() => {this.visible = false; this.searchData = {}}">取消</a-button>
          <a-button type="primary" @click="() => { this.pagination.current = 1; this.getTableData() }">确定</a-button>
        </div>
      </a-drawer>
    </a-card>
  </div>
</template>

<script>
/**
 * 企业成员管理页面
 * 功能说明：管理企业微信的员工成员数据
 * 主要功能：
 * 1. 查看企业成员列表
 * 2. 条件筛选：
 *    - 选择成员（按名称搜索）
 *    - 成功状态
 *    - 外部联系人权限
 * 3. 同步企业微信通讯录
 * 4. 查看最后一次同步时间
 * 5. 展示成员统计数据：
 *    - 发起申请数
 *    - 新增客户数
 *    - 聊天数
 *    - 发送消息数
 *    - 已回复聊天占比
 *    - 平均首次回复时长
 *    - 删除/拉黑客户数
 *
 * 业务场景：
 * - 企业管理员查看和管理企业微信员工
 * - 同步企业微信通讯录数据到系统
 * - 分析员工的工作绩效和客户服务情况
 *
 * 技术实现：
 * - 使用 a-table 展示成员列表
 * - 使用 a-drawer 抽屉组件进行条件筛选
 * - 使用 a-form 构建筛选表单
 * - 使用权限指令 v-permission 控制按钮权限
 */
import { enterMembersList, syncTime, syncEmployee } from '@/api/workEmployee'
export default {
  data () {
    return {
      visible: false,
      columns: [
        {
          title: '企业成员',
          align: 'center',
          dataIndex: 'name',
          width: '200px',
          scopedSlots: { customRender: 'name' }
        },
        {
          title: '状态',
          align: 'center',
          dataIndex: 'statusName'
        },
        {
          title: '外部联系人权限',
          align: 'center',
          dataIndex: 'contactAuthName'
        },
        {
          title: '发起申请数',
          align: 'center',
          dataIndex: 'applyNums'
        },
        {
          title: '新增客户数',
          align: 'center',
          dataIndex: 'addNums'
        },
        {
          title: '聊天数',
          align: 'center',
          dataIndex: 'messageNums'
        },
        {
          title: '发送消息数',
          align: 'center',
          dataIndex: 'sendMessageNums'
        },
        {
          title: '已回复聊天占比',
          align: 'center',
          dataIndex: 'replyMessageRatio'
        },
        {
          title: '平均首次回复时长',
          align: 'center',
          dataIndex: 'averageReply'
        },
        {
          title: '删除/拉黑客户数',
          align: 'center',
          dataIndex: 'invalidContact'
        }
      ],
      tableData: [],
      pagination: {
        total: 0,
        current: 1,
        pageSize: 10,
        showSizeChanger: true
      },
      syncTime: '',
      searchData: {},
      statusList: [],
      contactAuth: [],
      errorImg: 'this.src="' + require('@/assets/avatar.png') + '"'
    }
  },
  created () {
    this.getTableData()
    this.getSyncTime()
  },
  methods: {
    getTableData () {
      const params = {
        status: this.searchData.status,
        name: this.searchData.name,
        contactAuth: this.searchData.contactAuth,
        page: this.pagination.current,
        perPage: this.pagination.pageSize
      }
      enterMembersList(params).then(res => {
        this.tableData = res.data.list
        this.pagination.total = res.data.page.total
      })
    },
    // 同步时间获取-获取搜索条件数据
    getSyncTime () {
      syncTime().then(res => {
        this.syncTime = res.data.syncTime
        this.statusList = res.data.status
        this.contactAuth = res.data.contactAuth
      })
    },
    // 同步企业微信
    syncEmployeeDefine () {
      syncEmployee().then(res => {
        this.getTableData()
        this.$message.info('后台同步中,请稍等')
      })
    },
    onClose () {
      this.visible = false
    },
    afterVisibleChange () {

    },
    handleTableChange ({ current, pageSize }) {
      this.pagination.current = current
      this.pagination.pageSize = pageSize
      this.getTableData()
    }
  }
}
</script>

<style lang="less" scoped>
.work-employee {
  .btn-box {
    width: 100%;
    display: flex;
    justify-content: space-between;
    margin-bottom: 20px;
    .ant-btn {
      width: 80px;
    }
    div {
      width: 500px;
      display: flex;
      align-items: center;
      p{
        margin: 0;
        font-size: 16px;
      }
      .ant-btn {
        margin-left: 20px;
        width: 160px;
      }
    }
  }
  .table-wrapper {
    .msg-box {
      width: 100%;
      display: flex;
      align-items: center;
      justify-content: space-between;
    }
  }
}
</style>
