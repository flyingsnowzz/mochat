<!--
/**
 * 定时打标签规则列表页面
 *
 * 功能说明：
 * 1. 显示所有定时打标签规则，支持查看、删除、开关控制
 * 2. 支持按客户标签筛选规则
 * 3. 支持按规则名称搜索
 * 4. 定时规则用于在特定时间（如客户生日）自动为客户添加标签
 *
 * 业务场景：
 * - 客户生日自动打标签
 * - 重要日期提醒打标签
 * - 节假日自动打标签
 *
 * 使用技术：
 * - Vue 2.x
 * - Ant Design Vue 组件库
 * - Vuex 状态管理
 * - Vue Router 路由管理
 */
-->
<template>
  <!-- 卡片容器 -->
  <a-card>
    <div class="tag-box">
      <!-- 页面顶部：标题、操作按钮和搜索条件 -->
      <div class="page_top">
        <!-- 添加规则按钮 - 跳转到创建页面 -->
        <a-button type="primary" @click="$router.push('/autoTag/dayPartCreate')">添加规则</a-button>

        <!-- 搜索条件区域 -->
        <div class="search_term">
          <!-- 客户标签筛选 -->
          <div class="choice_tags">
            <span>客户标签：</span>
            <!-- 标签选择弹窗触发区域 -->
            <div class="showTagsBox" @click="showModel">
              <!-- 无标签时显示提示文字 -->
              <span class="showTagTips" v-if="selectTags.length==0">请选择标签</span>
              <!-- 有标签时显示已选标签 -->
              <div v-else>
                <a-tag v-for="(item,index) in selectTags" :key="index">{{ item.name }}</a-tag>
              </div>
            </div>
          </div>

          <!-- 标签选择组件 -->
          <selectTags @onChange="acceptData" :controlPopup="showPopup" ref="popupRef" />

          <!-- 规则名称搜索 -->
          <div class="">
            <span>搜索规则：</span>
            <a-input-search
              placeholder="请输入规则名称搜索"
              style="width: 200px"
              allow-clear
              v-model="searchRuleData"
              @search="searchRuleName"
              @change="emptyRule"
            />
          </div>
        </div>
      </div>

      <!-- 规则列表表格 -->
      <div class="rule-list">
        <a-table :columns="table.col" :data-source="table.data">
          <!-- 创建人列 - 显示创建人昵称 -->
          <div slot="nickname" slot-scope="text">
            <a-tag><a-icon type="user" />{{ text }}</a-tag>
          </div>

          <!-- 标签列 - 显示已添加的标签列表 -->
          <div slot="tags" slot-scope="text">
            <a-tag v-for="(item,index) in text" :key="index">{{ item }}</a-tag>
          </div>

          <!-- 规则状态列 - 显示开关状态 -->
          <div slot="on_off" slot-scope="text,record">
            <!-- 开关切换组件 -->
            <a-switch
              size="small"
              :defaultChecked="text==1"
              @change="openSwitch($event,record)"
            />
            <!-- 状态文字提示 -->
            <span v-if="text==1">已开启</span>
            <span v-else>已关闭</span>
          </div>

          <!-- 操作列 - 查看详情和删除 -->
          <div slot="operate" slot-scope="text,record">
            <span>
              <!-- 查看详情链接 -->
              <a @click="$router.push('/autoTag/dayPartShow?idRow='+record.id)">详情</a>
              <a-divider type="vertical" />
              <!-- 删除链接 -->
              <a @click="deltableRow(record)">删除</a>
            </span>
          </div>
        </a-table>
      </div>
    </div>
  </a-card>
</template>

<script>
/**
 * 定时打标签规则列表页脚本部分
 *
 * 主要功能：
 * 1. 获取并展示定时打标签规则列表
 * 2. 按标签和规则名称筛选
 * 3. 开关控制规则启用/禁用
 * 4. 删除规则
 */
import { indexApi, destroyApi, onOffApi } from '@/api/autoTag'
import selectTags from '@/components/addlabel/selectTags'

export default {
  // 注册子组件
  components: { selectTags },

  data () {
    return {
      // 已选择的客户标签列表
      selectTags: [],
      // 选中标签的ID数组，用于请求参数
      tagsParamArr: [],
      // 控制标签选择弹窗显示
      showPopup: false,

      // 表格配置和数据
      table: {
        // 表格列配置
        col: [
          {
            key: 'name',
            dataIndex: 'name',
            title: '规则名称'
          },
          {
            key: 'mark_tag_count',
            dataIndex: 'mark_tag_count',
            title: '已打标签数'
          },
          {
            key: 'nickname',
            dataIndex: 'nickname',
            title: '创建人',
            scopedSlots: { customRender: 'nickname' }
          },
          {
            key: 'created_at',
            dataIndex: 'created_at',
            title: '创建时间'
          },
          {
            key: 'tags',
            dataIndex: 'tags',
            title: '添加的标签',
            scopedSlots: { customRender: 'tags' }
          },
          {
            key: 'on_off',
            dataIndex: 'on_off',
            title: '规则状态',
            scopedSlots: { customRender: 'on_off' }
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
      },

      // 表格类型：3 表示定时打标签类型
      typeTable: 3,
      // 搜索规则名称
      searchRuleData: ''
    }
  },

  // 组件创建时获取表格数据
  created () {
    this.getTableData({
      type: this.typeTable
    })
  },

  methods: {
    /**
     * 更改规则开关状态
     * @param {Boolean} e - 开关状态
     * @param {Object} record - 当前规则记录
     */
    openSwitch (e, record) {
      // 根据开关状态设置规则状态：1-开启，2-关闭
      if (e) {
        record.on_off = 1
      } else {
        record.on_off = 2
      }
      // 调用接口更新规则状态
      onOffApi({
        id: record.id,
        on_off: record.on_off
      }).then((res) => {
        // 可以添加成功/失败提示
      })
    },

    /**
     * 显示标签选择弹窗
     * 调用子组件的show方法
     */
    showModel () {
      this.$refs.popupRef.show(this.selectTags)
    },

    /**
     * 接收子组件（标签选择组件）传来的数据
     * @param {Array} e - 选中的标签数组
     */
    acceptData (e) {
      // 重置标签ID数组
      this.tagsParamArr = []
      // 提取选中标签的ID
      e.forEach((item, index) => {
        this.tagsParamArr[index] = item.id
      })
      // 保存选中的标签对象数组
      this.selectTags = e
      // 重新获取表格数据
      this.getTableData({
        type: this.typeTable,
        name: this.searchRuleData,
        tags: this.tagsParamArr
      })
    },

    /**
     * 搜索规则名称
     * 点击搜索按钮时触发
     */
    searchRuleName () {
      this.getTableData({
        type: this.typeTable,
        name: this.searchRuleData,
        tags: this.tagsParamArr
      })
    },

    /**
     * 清空搜索框时触发
     * 如果搜索框为空，则重新获取所有数据
     * @param {Event} e - 输入事件
     */
    emptyRule (e) {
      if (this.searchRuleData == '') {
        this.getTableData({
          type: this.typeTable
        })
      }
    },

    /**
     * 删除表格数据行
     * @param {Object} record - 要删除的规则记录
     */
    deltableRow (record) {
      const that = this
      // 显示确认对话框
      this.$confirm({
        title: '提示',
        content: '是否删除',
        okText: '删除',
        okType: 'danger',
        cancelText: '取消',
        // 点击确定按钮时执行删除
        onOk () {
          // 找到要删除的行的索引
          const indexRow = that.table.data.indexOf(record)
          // 调用删除接口
          destroyApi({
            id: record.id
          }).then((res) => {
            // 删除成功后显示消息并从列表中移除
            that.$message.success('删除成功')
            that.table.data.splice(indexRow, 1)
          })
        }
      })
    },

    /**
     * 获取表格数据
     * @param {Object} paramData - 请求参数
     */
    getTableData (paramData) {
      indexApi(paramData).then((res) => {
        // 更新表格数据
        this.table.data = res.data.list
      })
    }
  }
}
</script>

<style lang="less" scoped>
/**
 * 样式部分
 */

/* 标签选择框样式 */
.showTagsBox{
  border: 1px solid #D9D9D9;  // 边框颜色
  width: 200px;               // 宽度
  cursor: pointer;            // 鼠标样式
  border-radius: 2px;         // 圆角
  padding-left: 5px;           // 左侧内边距
}

/* 无标签时的提示文字样式 */
.showTagTips{
  color: #BFBFBF;             // 灰色文字
}

/* 页面顶部布局 */
.page_top{
  display: flex;
  justify-content:space-between;  // 两端对齐
}

/* 规则列表容器 */
.rule-list{
  margin-top: 15px;           // 上边距
}

/* 搜索条件区域布局 */
.search_term{
  display: flex;              // 弹性布局
}

/* 标签选择区域 */
.choice_tags{
  display: flex;              // 弹性布局
  line-height: 32px;          // 行高
  margin-right: 15px;          // 右侧外边距
}
</style>
