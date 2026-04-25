<!--
/**
 * 定时打标签规则详情页面
 *
 * 功能说明：
 * 1. 显示定时打标签规则的详细信息，包括基本信息、规则设置、数据统计和客户列表
 * 2. 支持规则状态的开关控制
 * 3. 支持按客户名称、所属客服、规则、添加好友时间等条件筛选客户
 * 4. 显示打标签的统计数据，包括总数和今日数
 * 5. 支持查看客户详情
 *
 * 业务场景：
 * - 查看定时打标签规则的详细配置
 * - 监控规则的执行情况和打标签效果
 * - 管理规则的启用/禁用状态
 * - 查看规则影响的客户列表
 *
 * 使用技术：
 * - Vue 2.x
 * - Ant Design Vue 组件库
 * - Vuex 状态管理
 * - Vue Router 路由管理
 */
-->
<template>
  <div class="rules-details">
    <!-- 基本信息卡片 -->
    <a-card>
      <div class="title">基本信息</div>
      <div class="information">
        <!-- 左侧基本信息 -->
        <div class="ifm-left">
          <!-- 规则名称 -->
          <div class="name">
            <span>规则名称：</span>
            <span>{{ auto_tag.name }}</span>
          </div>
          <!-- 打标签方式 -->
          <div class="member mt16">
            <span>打标签方式：</span>
            <span>根据客户成为企业微信客户时间段</span>
          </div>
          <!-- 创建者 -->
          <div class="founder mt16">
            <span>创建者：</span>
            <a-tag><a-icon type="user" />{{ auto_tag.nickname }}</a-tag>
          </div>
          <!-- 生效成员 -->
          <div class="add-tag mt16">
            <span>生效成员：</span>
            <span>
              <a-tag v-for="(item,index) in auto_tag.employees" :key="index">
                <a-icon type="user" />{{ item.name }}
              </a-tag>
            </span>
          </div>
          <!-- 创建时间 -->
          <div class="create-time mt16">
            <span>创建时间：</span>
            <span>{{ auto_tag.createdAt }}</span>
          </div>
          <!-- 自动添加标签 -->
          <div class="auto mt16">
            <span>自动添加标签：</span>
            <a-tag v-for="(item,index) in auto_tag.tags" :key="index">{{ item }}</a-tag>
          </div>
          <!-- 规则状态 -->
          <div class="state mt16">
            <span>规则状态：</span>
            <a-switch size="small" :checked="auto_tag.onOff==1" class="mr4"/>
            <span v-if="auto_tag.onOff">已开启</span>
            <span v-else>已关闭</span>
          </div>
        </div>

        <!-- 右侧规则设置 -->
        <div class="ifm-right">
          <div class="title-set">
            <span>规则设置：</span>
            <span class="black">共{{ tagRule.length }}条规则</span>
          </div>
          <!-- 规则详情列表 -->
          <div class="set-content mt18 ml84" v-for="(item,index) in tagRule" :key="index">
            <span class="title-small black">规则{{ index+1 }}</span>
            <p class="keys ml4">
              客户在 每
              <span v-if="item.time_type==1">天</span>
              <span v-if="item.time_type==2">周</span>
              <span v-if="item.time_type==3">月</span> 的
              <!-- 周规则 -->
              <template v-if="item.time_type==2">
                <span v-for="obj in item.schedule" :key="obj">{{ weekData[obj].value }}</span>
              </template>
              <!-- 月规则 -->
              <template v-if="item.time_type==3">
                <span v-for="obj in item.schedule" :key="obj">{{ monthData[obj].value }}</span>
              </template>
              【{{ item.start_time }}-{{ item.end_time }}】内添加生效成员时，
              将会被自动打上
              <a-tag v-for="(obj,idx) in item.tags" :key="idx">{{ obj.tagname }}</a-tag>
              标签
            </p>
          </div>
        </div>
      </div>
    </a-card>

    <!-- 数据统计卡片 -->
    <a-card class="mt16">
      <div class="title">数据统计</div>
      <div class="data-num">
        <div class="data-box mb20">
          <div class="data">
            <!-- 打标签总数 -->
            <div class="item">
              <div class="count">{{ statistics.total_count }}</div>
              <div class="desc">打标签总数</div>
            </div>
            <!-- 今日打标签数 -->
            <div class="item">
              <div class="count">{{ statistics.today_count }}</div>
              <div class="desc">今日打标签数</div>
            </div>
          </div>
        </div>
      </div>
    </a-card>

    <!-- 客户列表卡片 -->
    <a-card class="mt16">
      <!-- 搜索条件 -->
      <div class="search">
        <div class="search-box">
          <!-- 搜索用户 -->
          <div class="customer">
            <span>搜索用户：</span>
            <a-input-search
              placeholder="请输入要搜索的客户"
              style="width: 200px"
              v-model="paramsTable.contact_name"
              @search="searchUser"
              allow-clear
              @change="emptyInput"
            />
          </div>
          <!-- 所属客服 -->
          <div class="customer-service ml20">
            <span>所属客服：</span>
            <a-select style="width: 160px" default-value="请选择客服">
              <a-select-option value="Home">刘波</a-select-option>
              <a-select-option value="Company">小子</a-select-option>
            </a-select>
          </div>
          <!-- 规则筛选 -->
          <div class="join-mode">
            <span class="ml20">规则筛选：</span>
            <a-select style="width: 180px" default-value="请选择规则">
              <a-select-option value="1">规则一</a-select-option>
              <a-select-option value="2">规则二</a-select-option>
            </a-select>
          </div>
          <!-- 添加好友时间 -->
          <div class="add-time ml20">
            <span>添加好友时间：</span>
            <a-range-picker style="width: 220px" @change="searchTime" :allowClear="true" v-model="selectDate"/>
          </div>
        </div>
        <!-- 重置按钮 -->
        <div class="reset"><a-button @click="resetTable">重置</a-button></div>
      </div>

      <!-- 客户列表 -->
      <div class="table-box mt36">
        <div class="store-box mb20">
          <span class="customers-title">共{{ table.data.length }}个客户</span>
          <a-divider type="vertical" />
          <span style="cursor: pointer;" @click="updateTable"><a-icon type="redo" />更新数据</span>
        </div>
        <div class="table">
          <a-table :columns="table.col" :data-source="table.data">
            <!-- 所属客服列 -->
            <div slot="employeeName" slot-scope="text">
              <a-tag><a-icon type="user" />{{ text }}</a-tag>
            </div>
            <!-- 生效规则列 -->
            <div slot="tagRuleId" slot-scope="text">
              <a-tag>规则{{ text }}</a-tag>
            </div>
            <!-- 操作列 -->
            <div slot="operate" slot-scope="text,record">
              <div>
                <a @click="clientDetails(record)">客户详情</a>
              </div>
            </div>
          </a-table>
        </div>
      </div>
    </a-card>
  </div>
</template>

<script>
/**
 * 定时打标签规则详情页脚本部分
 *
 * 主要功能：
 * 1. 获取并展示规则的详细信息
 * 2. 获取并展示打标签的统计数据
 * 3. 获取并展示受规则影响的客户列表
 * 4. 支持按多种条件筛选客户
 * 5. 支持重置筛选条件和更新数据
 */
import { showApi, showContactTimeApi } from '@/api/autoTag'

export default {
  data () {
    return {
      // 周数据映射：用于显示星期几
      weekData: [
        { key: '0', value: '周日' },
        { key: '1', value: '周一' },
        { key: '2', value: '周二' },
        { key: '3', value: '周三' },
        { key: '4', value: '周四' },
        { key: '5', value: '周五' },
        { key: '6', value: '周六' }
      ],
      // 月数据映射：用于显示日期
      monthData: [
        { key: '1', value: '1号' },
        { key: '2', value: '2号' },
        { key: '3', value: '3号' },
        { key: '4', value: '4号' },
        { key: '5', value: '5号' },
        { key: '6', value: '6号' },
        { key: '7', value: '7号' },
        { key: '8', value: '8号' },
        { key: '9', value: '9号' },
        { key: '10', value: '10号' },
        { key: '11', value: '11号' },
        { key: '12', value: '12号' },
        { key: '13', value: '13号' },
        { key: '14', value: '14号' },
        { key: '15', value: '15号' },
        { key: '16', value: '16号' },
        { key: '17', value: '17号' },
        { key: '18', value: '18号' },
        { key: '19', value: '19号' },
        { key: '20', value: '20号' },
        { key: '21', value: '21号' },
        { key: '22', value: '22号' },
        { key: '23', value: '23号' },
        { key: '24', value: '24号' },
        { key: '25', value: '25号' },
        { key: '26', value: '26号' },
        { key: '27', value: '27号' },
        { key: '28', value: '28号' },
        { key: '29', value: '29号' },
        { key: '30', value: '30号' },
        { key: '31', value: '31号' },
        { key: '32', value: '32号' }
      ],
      // 规则基本信息
      auto_tag: [],
      // 规则详情列表
      tagRule: [],
      // 统计数据
      statistics: [],
      // 选择的日期范围
      selectDate: [],
      // 表格请求参数
      paramsTable: {
        contact_name: '', // 搜索客户名称
        employee: '', // 客服
        start_time: '', // 开始时间
        end_time: '', // 结束时间
        page: 1, // 页码
        perPage: 15 // 每页条数
      },
      // 表格配置和数据
      table: {
        // 表格列配置
        col: [
          {
            key: 'contactName',
            dataIndex: 'contactName',
            title: '客户',
            scopedSlots: { customRender: 'contactName' }
          },
          {
            key: 'employeeName',
            dataIndex: 'employeeName',
            title: '所属客服',
            scopedSlots: { customRender: 'employeeName' }
          },
          {
            key: 'tagRuleId',
            dataIndex: 'tagRuleId',
            title: '生效规则',
            scopedSlots: { customRender: 'tagRuleId' }
          },
          {
            key: 'createdAt',
            dataIndex: 'createdAt',
            title: '添加好友时间'
          },
          {
            key: 'operate',
            dataIndex: 'operate',
            title: '操作',
            scopedSlots: { customRender: 'operate' }
          }
        ],
        // 表格数据
        data: []
      }
    }
  },

  // 组件创建时获取数据
  created () {
    // 获取路由参数中的规则ID
    this.idRow = this.$route.query.idRow
    this.paramsTable.id = this.idRow
    // 获取规则详情数据
    this.getDetailsData(this.idRow)
    // 获取客户列表数据
    this.getTableData(this.paramsTable)
  },

  methods: {
    /**
     * 搜索时间范围
     * @param {Array} date - 日期对象数组
     * @param {Array} dateString - 日期字符串数组
     */
    searchTime (date, dateString) {
      this.paramsTable.start_time = dateString[0]
      this.paramsTable.end_time = dateString[1]
      this.getTableData(this.paramsTable)
    },

    /**
     * 重置表格筛选条件
     */
    resetTable () {
      // 重置请求参数
      var resetParams = {
        id: this.idRow,
        page: 1,
        perPage: 15
      }
      // 清空日期选择
      this.selectDate = []
      // 更新参数并重新获取数据
      this.paramsTable = resetParams
      this.getTableData(this.paramsTable)
    },

    /**
     * 清空输入框时触发
     * 如果搜索框为空，则重新获取数据
     */
    emptyInput () {
      if (this.paramsTable.contact_name == '') {
        this.getTableData(this.paramsTable)
      }
    },

    /**
     * 搜索用户
     * 点击搜索按钮时触发
     */
    searchUser () {
      this.getTableData(this.paramsTable)
    },

    /**
     * 更新表格数据
     * 点击更新数据按钮时触发
     */
    updateTable () {
      // 清空现有数据
      this.table.data = []
      // 重新获取数据
      this.getTableData(this.paramsTable)
    },

    /**
     * 查看客户详情
     * @param {Object} record - 客户记录
     */
    clientDetails (record) {
      // 可以跳转到客户详情页面
    },

    /**
     * 获取客户列表数据
     * @param {Object} data - 请求参数
     */
    getTableData (data) {
      showContactTimeApi(data).then((res) => {
        this.table.data = res.data.list
      })
    },

    /**
     * 获取规则详情数据
     * @param {Number} id - 规则ID
     */
    getDetailsData (id) {
      showApi({ id }).then((res) => {
        // 更新规则基本信息
        this.auto_tag = res.data.auto_tag
        this.tagRule = this.auto_tag.tagRule
        // 更新统计数据
        this.statistics = res.data.statistics
      })
    }
  }
}
</script>

<style lang="less" scoped>
/**
 * 样式部分
 */

/* 标题样式 */
.title {
  font-size: 15px;              // 字体大小
  line-height: 21px;            // 行高
  color: rgba(0, 0, 0, .85);     // 字体颜色
  border-bottom: 1px solid #e9ebf3; // 底部边框
  padding-bottom: 16px;          // 底部内边距
  margin-bottom: 16px;           // 底部外边距
  position: relative;             // 相对定位
}

/* 信息区域布局 */
.information{
  display: flex;                  // 弹性布局

  /* 左侧信息区域 */
  .ifm-left{
    flex: 0.4;                   // 占40%宽度
  }

  /* 右侧规则设置区域 */
  .ifm-right{
    flex: 1.6;                   // 占60%宽度
    background-color: #f6f6f6;   // 背景色
    border: 1px solid #e7e7e7;    // 边框
    min-height: 260px;           // 最小高度
  }

  /* 规则设置标题 */
  .title-set{
    padding-top: 14px;            // 顶部内边距
    padding-left: 12px;           // 左侧内边距
  }
}

/* 黑色文字样式 */
.black{
  color: #000;                    // 黑色文字
}

/* 规则标题样式 */
.title-small{
  display: block;                 // 块级元素
  padding-left: 4px;             // 左侧内边距
  border-left: 3px solid #8d8d8d; // 左侧边框
  margin-bottom: 4px;            // 底部外边距
}

/* 规则详情文字样式 */
.keys{
  span{
    font-weight: bold;            // 粗体
    margin-left: 1px;             // 左侧外边距
    margin-right: 1px;            // 右侧外边距
  }
}

/* 数据统计盒子样式 */
.data-box {
  display: flex;                  // 弹性布局
  justify-content: center;        // 水平居中
  flex-direction: column;         // 垂直方向
  margin-top: 25px;              // 顶部外边距
  width: 500px;                  // 宽度
  height: 125px;                 // 高度

  /* 数据内容区域 */
  .data {
    flex: 1;                     // 占满剩余空间
    height: 120px;               // 高度
    background: #fbfdff;          // 背景色
    border: 1px solid #daedff;    // 边框
    display: flex;                // 弹性布局
    align-items: center;          // 垂直居中

    /* 数据项 */
    .item {
      flex: 1;                   // 平均分配空间
      border-right: 1px solid #e9e9e9; // 右侧边框

      /* 数据值 */
      .count {
        font-size: 24px;          // 字体大小
        font-weight: 500;         // 字重
        text-align: center;        // 文本居中
      }

      /* 数据描述 */
      .desc {
        font-size: 13px;          // 字体大小
        text-align: center;        // 文本居中
      }

      /* 最后一个数据项，去除右侧边框 */
      &:last-child {
        border-right: 0;
      }
    }

    /* 最后一个数据盒子，去除右侧外边距 */
    &:last-child {
      margin-right: 0;
    }
  }
}

/* 搜索区域布局 */
.search{
  display: flex;                  // 弹性布局

  /* 搜索框容器 */
  .search-box{
    display: flex;                // 弹性布局
    flex: 1;                     // 占满剩余空间
  }
}

/* 状态区域布局 */
.state{
  display: flex;                  // 弹性布局
  align-items: center;            // 垂直居中
}
</style>
