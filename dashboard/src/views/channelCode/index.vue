<!--
/**
 * 渠道码管理列表页面
 *
 * 功能说明：
 * 1. 显示所有渠道码列表，支持分页、搜索和筛选
 * 2. 支持按名称、活码类型、分组进行筛选
 * 3. 支持渠道码的编辑、查看客户、下载、移动分组、查看统计等操作
 * 4. 支持新建渠道码（跳转到创建页面）
 * 5. 支持新建分组和修改分组
 * 6. 支持移动渠道码到不同分组
 *
 * 业务场景：
 * - 管理企业的渠道活码
 * - 跟踪不同渠道的客户来源
 * - 批量管理渠道码分组
 * - 查看渠道码的客户统计数据
 *
 * 使用技术：
 * - Vue 2.x
 * - Ant Design Vue 组件库
 * - Vuex 状态管理
 * - Vue Router 路由管理
 */
-->
<template>
  <div class="wrapper">
    <!-- 功能说明横幅 -->
    <div :split="false" class="lists">
      <a-alert :show-icon="false" message="1、可以生成带参数的二维码名片，支持活码功能，即随机选取设置的活码成员推给用户。加企业微信为好友后，可以给微信联系人自动回复相应欢迎消息和打标签。" banner />
      <a-alert :show-icon="false" message="2、每创建一个渠道活码，该码则自动进入【内容引擎】--【图片类型】--分组【渠道码-企业微信】，以素材的方式通过聊天侧边栏快速发送给客户。" banner />
      <a-alert :show-icon="false" message="3、受限于官方，单人类型的渠道码创建后尽量不要再修改成员，否则会造成列表中，该二维码中间的头像与配置的成员头像不一致，但是并不影响功能使用。" banner />
      <a-alert :show-icon="false" message="4、如果企业在企业微信后台为相关成员配置了可用的欢迎语，使用第三方系统配置欢迎语，则均不起效，推送的还是企业微信官方的。" banner />
    </div>

    <!-- 主要内容卡片 -->
    <a-card>
      <!-- 搜索筛选表单 -->
      <a-form :label-col="{ span: 7 }" :wrapper-col="{ span: 14 }">
        <a-row :gutter="16">
          <!-- 名称搜索 -->
          <a-col :lg="8">
            <a-form-item label="名称：">
              <a-input v-model="screentData.name" placeholder="搜索活码名称"></a-input>
            </a-form-item>
          </a-col>
          <!-- 活码类型筛选 -->
          <a-col :lg="8">
            <a-form-item label="活码类型：">
              <a-select v-model="screentData.type">
                <a-select-option v-for="item in typeList" :key="item.value">
                  {{ item.label }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
          <!-- 分组筛选 -->
          <a-col :lg="8">
            <a-form-item label="分组：">
              <a-select v-model="screentData.groupId">
                <a-select-option v-for="item in groupList" :key="item.groupId">
                  {{ item.name }}
                </a-select-option>
              </a-select>
            </a-form-item>
          </a-col>
        </a-row>
      </a-form>

      <!-- 搜索按钮 -->
      <div class="search">
        <a-button v-permission="'/channelCode/index@search'" type="primary" style="marginRight: 10px" @click="getTableData">查询</a-button>
        <a-button @click="() => {this.screentData = {}}">重置</a-button>
      </div>

      <!-- 操作按钮 -->
      <div class="btn-box">
        <a-button v-permission="'/channelCode/index@editGroup'" type="primary" @click="() => {this.editGroupDis = true}">修改分组</a-button>
        <a-button v-permission="'/channelCode/index@add'" type="primary" @click="() => {this.addGroupDis = true}">新建分组</a-button>
        <router-link :to="{path: '/channelCode/store'}">
          <a-button v-permission="'/channelCode/store'" type="primary">+新建</a-button>
        </router-link>
      </div>

      <!-- 渠道码列表表格 -->
      <a-table
        style="marginTop: 20px"
        bordered
        :columns="columns"
        :data-source="tableData"
        :rowKey="record => record.channelCodeId"
        :pagination="pagination"
        @change="handleTableChange">
        <!-- 二维码列 -->
        <div slot="image" slot-scope="text, record" class="img-box">
          <img style="width:90px; height:auto;" :src="record.qrcodeUrl" alt="">
          <span>{{ record.type }}</span>
        </div>
        <!-- 标签列 -->
        <div slot="tags" slot-scope="text, record">
          <div v-if="record.tags.length !== 0" v-for="(item, index) in record.tags" :key="index + 'lseg'">{{ item }}</div>
        </div>
        <!-- 操作列 -->
        <div slot="action" slot-scope="text, record">
          <template>
            <!-- 编辑按钮 -->
            <router-link :to="{path:`/channelCode/store?channelCodeId=${record.channelCodeId}`}">
              <a-button v-permission="'/channelCode/index@edit'" type="link">编辑</a-button>
            </router-link>
            <!-- 客户按钮 -->
            <a-button v-permission="'/channelCode/index@customer'" type="link" @click="getChannelCodeContact(record.channelCodeId)">客户</a-button>
            <!-- 下载按钮 -->
            <a-button v-permission="'/channelCode/index@download'" type="link">
              <a :href="record.qrcodeUrl" download>下载</a>
            </a-button>
            <!-- 移动按钮 -->
            <a-button v-permission="'/channelCode/index@move'" type="link" @click="moveGroup(record.channelCodeId, record.groupId)">移动</a-button>
            <!-- 统计按钮 -->
            <router-link :to="{path: `/channelCode/statistics?channelCodeId=${record.channelCodeId}`}">
              <a-button v-permission="'/channelCode/statistics'" type="link">统计</a-button>
            </router-link>
          </template>
        </div>
      </a-table>
    </a-card>

    <!-- 客户列表模态框 -->
    <a-modal
      width="800px"
      title="扫码客户"
      :visible="contactModal"
      @cancel="() => {this.contactModal = false; this.contactPagination.current = 1;}">
      <a-table
        style="marginTop: 20px;"
        bordered
        :columns="contactColumns"
        :data-source="contactTableData"
        :rowKey="record => record.contactId"
        :pagination="contactPagination"
        @change="handleTableChangeContact">
      </a-table>
      <template slot="footer">
        <a-button type="primary" @click="() => {this.contactModal = false; this.contactPagination.current = 1;}">确定</a-button>
      </template>
    </a-modal>

    <!-- 新建分组模态框 -->
    <a-modal
      title="新建分组"
      :visible="addGroupDis"
      @cancel="() => {this.addGroupDis = false; this.groupName = ''}">
      <a-form-model :label-col="{ span: 6 }" :wrapper-col="{ span: 14}">
        <a-form-model-item label="分组名称：">
          <p>每个分组名称最多15个字。同时新建多个分组时，请用“空格”隔开</p>
          <a-input v-model="groupName"/>
        </a-form-model-item>
      </a-form-model>
      <template slot="footer">
        <a-button @click="() => {this.addGroupDis = false; this.groupName = ''}">取消</a-button>
        <a-button type="primary" :loading="btnLoading" @click="addGroup">确定</a-button>
      </template>
    </a-modal>

    <!-- 修改分组模态框 -->
    <a-modal
      title="修改分组"
      :visible="editGroupDis"
      @cancel="() => {this.editGroupDis = false; this.editGroupData = {}}">
      <a-form-model :label-col="{ span: 6 }" :wrapper-col="{ span:14}">
        <a-form-model-item label="选择分组：">
          <a-select v-model="editGroupData.groupId">
            <a-select-option v-for="item in groupList" :key="item.groupId">
              {{ item.name }}
            </a-select-option>
          </a-select>
        </a-form-model-item>
        <a-form-model-item label="修改分组名称：">
          <a-input v-model="editGroupData.name"/>
        </a-form-model-item>
      </a-form-model>
      <template slot="footer">
        <a-button @click="() => {this.editGroupDis = false; this.editGroupData = {}}">取消</a-button>
        <a-button type="primary" :loading="btnLoading" @click="editGroup">确定</a-button>
      </template>
    </a-modal>

    <!-- 移动分组容器 -->
    <div class="mbox" ref="mbox"></div>
    <!-- 移动分组模态框 -->
    <a-modal
      :getContainer="() => $refs.mbox"
      title="移动分组"
      class="move-box"
      :visible="moveGroupDis"
      @cancel="() => {this.moveGroupDis = false}">
      <div class="group-box">
        <div :class="moveGroupId == item.groupId ? 'active' : ''" v-for="item in groupList" :key="item.groupId" @click="changeGroup(item.groupId)">
          {{ item.name }}
        </div>
      </div>
      <template slot="footer">
        <a-button @click="() => {this.moveGroupDis = false; this.moveGroupData = {}}">取消</a-button>
        <a-button type="primary" @click="moveGroupDefined">确定</a-button>
      </template>
    </a-modal>
  </div>
</template>

<script>
/**
 * 渠道码管理列表页脚本部分
 *
 * 主要功能：
 * 1. 获取并展示渠道码列表
 * 2. 按名称、类型、分组筛选渠道码
 * 3. 分页管理渠道码列表
 * 4. 管理渠道码分组（新建、修改、移动）
 * 5. 查看渠道码的客户列表
 */
import { channelCodeList, channelCodeGroup, channelCodeGroupAdd, channelCodeGroupUpdate, channelCodeContact, channelCodeGroupMove } from '@/api/channelCode'

export default {
  data () {
    return {
      // 渠道码ID
      channelCodeId: '',
      // 修改分组模态框显示状态
      editGroupDis: false,
      // 新建分组模态框显示状态
      addGroupDis: false,
      // 通用模态框显示状态
      visible: false,
      // 按钮加载状态
      btnLoading: false,
      // 客户列表模态框显示状态
      contactModal: false,
      // 移动分组模态框显示状态
      moveGroupDis: false,
      // 表格列配置
      columns: [
        {
          align: 'center',
          title: '二维码',
          dataIndex: 'image',
          scopedSlots: { customRender: 'image' }
        },
        {
          align: 'center',
          title: '名称',
          dataIndex: 'name'
        },
        {
          align: 'center',
          title: '分组',
          dataIndex: 'groupName'
        },
        {
          align: 'center',
          title: '客户数',
          dataIndex: 'contactNum'
        },
        {
          align: 'center',
          title: '标签',
          dataIndex: 'tags',
          scopedSlots: { customRender: 'tags' }
        },
        {
          align: 'center',
          title: '自动添加好友',
          dataIndex: 'autoAddFriend'
        },
        {
          align: 'center',
          title: '操作',
          width: '250px',
          dataIndex: 'action',
          scopedSlots: { customRender: 'action' }
        }
      ],
      // 表格数据
      tableData: [],
      // 分页配置
      pagination: {
        total: 0,
        current: 1,
        pageSize: 10,
        showSizeChanger: true
      },
      // 筛选条件
      screentData: {},
      // 分组列表
      groupList: [],
      // 类型列表
      typeList: [
        {
          label: '全部',
          value: 0
        },
        {
          label: '单人',
          value: 1
        },
        {
          label: '多人',
          value: 2
        }
      ],
      // 客户数据
      contactTableData: [],
      // 客户表格列配置
      contactColumns: [
        {
          align: 'center',
          title: '客户名称',
          dataIndex: 'name'
        },
        {
          align: 'center',
          title: '归属成员',
          dataIndex: 'employees'
        },
        {
          align: 'center',
          title: '添加时间',
          dataIndex: 'createTime'
        }
      ],
      // 客户列表分页配置
      contactPagination: {
        total: 0,
        current: 1,
        pageSize: 10,
        showSizeChanger: true
      },
      // 分组名称
      groupName: '',
      // 修改分组数据
      editGroupData: {},
      // 移动分组ID
      moveGroupId: ''
    }
  },

  // 组件创建时获取数据
  created () {
    // 获取渠道码列表
    this.getTableData()
    // 获取分组列表
    this.getGroupList()
  },

  methods: {
    /**
     * 获取渠道码列表数据
     */
    getTableData () {
      const params = {
        name: this.screentData.name,
        groupId: this.screentData.groupId,
        type: this.screentData.type,
        page: this.pagination.current,
        perPage: this.pagination.pageSize
      }
      channelCodeList(params).then(res => {
        this.tableData = res.data.list
        this.pagination.total = res.data.page.total
      })
    },

    /**
     * 表格分页变化处理
     * @param {Object} pagination - 分页对象
     */
    handleTableChange ({ current, pageSize }) {
      this.pagination.current = current
      this.pagination.pageSize = pageSize
      this.getTableData()
    },

    /**
     * 客户列表分页变化处理
     * @param {Object} pagination - 分页对象
     */
    handleTableChangeContact ({ current, pageSize }) {
      this.contactPagination.current = current
      this.contactPagination.pageSize = pageSize
      this.getChannelCodeContact(this.channelCodeId)
    },

    /**
     * 获取分组列表
     */
    getGroupList () {
      channelCodeGroup().then(res => {
        this.groupList = res.data
      })
    },

    /**
     * 获取渠道码的客户列表
     * @param {Number} id - 渠道码ID
     */
    getChannelCodeContact (id) {
      this.channelCodeId = id
      channelCodeContact({
        channelCodeId: this.channelCodeId,
        page: this.contactPagination.current,
        perPage: this.contactPagination.pageSize
      }).then(res => {
        this.contactModal = true
        this.contactTableData = res.data.list
        this.contactPagination.total = res.data.page.total
      })
    },

    /**
     * 新增分组
     */
    addGroup () {
      // 处理分组名称，支持多个分组（用空格分隔）
      const reg = /\s+/g
      const name = this.groupName.replace(reg, ' ').split(' ')
      let falg = false
      let flag = false
      let fg = false

      // 检查是否有重复分组名
      if (name.length > 1) {
        const lengths = Array.from(new Set(name)).length
        if (lengths < name.length) {
          flag = true
        }
      }

      // 检查分组名是否合法
      name.map((item, index) => {
        if (item == '') {
          fg = true
        }
        if (item.length > 15) {
          this.$message.error('每个敏感词最多15个字')
          falg = true
        }
      })

      // 验证失败处理
      if (falg) {
        return
      }
      if (fg) {
        this.$message.error('请输入正确的分组名称')
        return
      }
      if (flag) {
        this.$message.error('分组名称重复')
        return
      }

      // 调用接口创建分组
      this.btnLoading = true
      channelCodeGroupAdd({
        name: name
      }).then(res => {
        // 刷新分组列表
        this.getGroupList()
        // 关闭模态框
        this.addGroupDis = false
        // 清空分组名称
        this.groupName = ''
        // 重置按钮状态
        this.btnLoading = false
      }).catch(res => {
        this.btnLoading = false
      })
    },

    /**
     * 修改分组
     */
    editGroup () {
      const reg = /\s+/g
      const name = this.editGroupData.name

      // 验证分组选择
      if (this.editGroupData.groupId == undefined) {
        this.$message.error('请选择分组')
        return
      }

      // 不能修改未分组
      if (this.editGroupData.groupId == 0) {
        this.$message.error('不能修改【未分组】')
        return
      }

      // 验证分组名称
      if (reg.test(name)) {
        this.$message.error('请输入正确的分组名称')
        return
      }

      // 调用接口修改分组
      this.btnLoading = true
      channelCodeGroupUpdate({
        groupId: this.editGroupData.groupId,
        name: this.editGroupData.name
      }).then(res => {
        // 清空修改分组数据
        this.editGroupData = {}
        // 刷新分组列表
        this.getGroupList()
        // 关闭模态框
        this.editGroupDis = false
        // 重置按钮状态
        this.btnLoading = false
      }).catch(res => {
        this.btnLoading = false
      })
    },

    /**
     * 移动渠道码到其他分组
     * @param {Number} id - 渠道码ID
     * @param {Number} groupId - 当前分组ID
     */
    moveGroup (id, groupId) {
      this.moveGroupDis = true
      this.channelCodeId = id
      this.moveGroupId = groupId
    },

    /**
     * 选择要移动到的分组
     * @param {Number} id - 目标分组ID
     */
    changeGroup (id) {
      this.moveGroupId = id
    },

    /**
     * 确定移动分组
     */
    moveGroupDefined () {
      channelCodeGroupMove({
        channelCodeId: this.channelCodeId,
        groupId: this.moveGroupId
      }).then(res => {
        // 关闭模态框
        this.moveGroupDis = false
        // 清空移动分组数据
        this.moveGroupId = ''
        this.channelCodeId = ''
        // 刷新渠道码列表
        this.getTableData()
      })
    }
  }
}
</script>

<style lang="less" scoped>
/**
 * 样式部分
 */

.wrapper{
  /* 搜索按钮区域 */
  .search {
    display: flex;
    justify-content: flex-end;
  }

  /* 操作按钮区域 */
  .btn-box {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
    .ant-btn {
      margin-left: 10px;
    }
  }

  /* 二维码图片容器 */
  .img-box {
    display: flex;
    flex-direction: column;
    align-items: center
  }

  /* 移动分组模态框 */
  .mbox {
    .move-box {
      .group-box {
        display: flex;
        flex-wrap: wrap;
        div {
          padding: 0 10px;
          height: 40px;
          line-height: 40px;
          text-align: center;
          margin: 5px 10px;
          border: 1px solid #ccc;
          border-radius: 5px;
          cursor: pointer;
        }
        /* 选中分组样式 */
        .active {
          padding: 0 10px;
          height: 40px;
          line-height: 40px;
          text-align: center;
          margin: 5px 10px;
          border: 2px solid #1890FF;
          border-radius: 5px;
          cursor: pointer;
        }
      }
    }
  }
}
</style>
