<template>
  <div class="room-fission-index">
    <div class="add-row">
      <div class="btn">
        <router-link :to="{ path: '/workFission/create' }">
          <a-button type="primary">创建活动</a-button>
        </router-link>
      </div>
      <div class="search">
        <a-input-search placeholder="请输入要搜索的活动" @search="search" />
      </div>
    </div>
    <a-card>
      <a-table :columns="table.columns" :data-source="table.data">
        <div class="btn-group" slot="operating" slot-scope="item">
          <a @click="goInvite(item)">邀请客户参与</a>
          <a-divider type="vertical" />
          <a @click="$refs.details.show(item.id)">详情</a>
          <a-divider type="vertical" />
          <a-popover trigger="click" placement="bottomRight">
            <template slot="content">
              <div class="divider-btn" @click="goUpdate(item)">修改</div>
              <div class="divider-btn" @click="del(item)">删除</div>
            </template>
            <a>编辑</a>
            <a-icon type="caret-down" :style="{ color: '#1890ff' }" />
          </a-popover>
        </div>
        <div slot="member_use" slot-scope="item">
          <a-tag v-for="user in item.serviceEmployees" :key="user.id">
            <a-icon type="user" two-tone-color="#7da3d1" :style="{ color: '#7da3d1' }" />
            {{ user.name }}
          </a-tag>
        </div>
        <div slot="client_tag" slot-scope="item">
          <a-tag v-for="tag in item.contactTags" :key="tag.id">
            {{ tag.name }}
          </a-tag>
        </div>
      </a-table>
    </a-card>

    <Details ref="details" />
  </div>
</template>

<script>
/**
 * 工作裂变活动管理页面
 * 功能说明：管理企业微信的工作裂变活动
 * 主要功能：
 * 1. 查看活动列表（活动名称、使用成员、扫码添加人数、创建时间、活动状态）
 * 2. 创建新活动
 * 3. 搜索活动
 * 4. 邀请客户参与活动
 * 5. 查看活动详情
 * 6. 修改活动
 * 7. 删除活动
 *
 * 业务场景：
 * - 企业创建裂变活动，鼓励员工邀请客户参与
 * - 追踪活动效果和参与情况
 * - 管理活动的生命周期
 *
 * 技术实现：
 * - 使用 a-table 展示活动列表
 * - 使用 a-popover 实现编辑下拉菜单
 * - 使用 a-confirm 弹窗确认删除操作
 * - 使用 Details 组件查看活动详情
 */
import Details from '@/views/workFission/components/details'

import { getList, del } from '@/api/workFission'

export default {
  data() {
    return {
      table: {
        columns: [
          {
            title: '活动名称',
            dataIndex: 'activeName',
          },
          {
            title: '使用成员',
            scopedSlots: { customRender: 'member_use' },
          },
          {
            title: '扫码添加人数',
            dataIndex: 'newFriend',
          },
          {
            title: '创建时间',
            dataIndex: 'createdAt',
          },
          {
            title: '活动状态',
            dataIndex: 'status',
          },
          {
            title: '操作',
            scopedSlots: { customRender: 'operating' },
          },
        ],
        data: [],
      },
    }
  },
  mounted() {
    this.getData()
  },
  methods: {
    goInvite(data) {
      this.$router.push({
        path: '/workFission/invite',
        query: {
          id: data.id,
        },
      })
    },

    goUpdate(data) {
      this.$router.push({
        path: '/workFission/edit',
        query: {
          id: data.id,
        },
      })
    },

    search(e) {
      this.getData(e)
    },

    del(data) {
      const that = this
      this.$confirm({
        title: '提示',
        content: '是否删除',
        okText: '删除',
        okType: 'danger',
        cancelText: '取消',
        onOk() {
          del({ id: data.id }).then((res) => {
            if (res.code === 200) {
              that.getData()
              that.$message.success('删除成功')
            }
          })
        },
      })
    },
    getData(key = '') {
      getList({
        active_name: key,
      }).then((res) => {
        for (const v of res.data.list) {
          if (v.service_employees) v.service_employees = JSON.parse(v.service_employees)
          if (v.contact_tags) v.contact_tags = JSON.parse(v.contact_tags)
        }

        this.table.data = res.data.list
      })
    },
  },
  components: { Details },
}
</script>

<style lang="less" scoped>
.add-row {
  display: flex;
  margin-bottom: 16px;

  .btn {
    flex: 1;

    .ant-btn {
      margin-right: 16px;
    }
  }

  span {
    cursor: pointer;
  }
}

.btn-group {
  font-size: 13px;
}

.divider-btn {
  height: 33px;
  display: flex;
  align-items: center;
  justify-content: center;
  cursor: pointer;
}

/deep/ .ant-card-body {
  padding: 0;
}

/deep/ .ant-table-pagination.ant-pagination {
  margin-right: 20px;
}

/deep/ .ant-popover-inner-content {
  padding: 0 !important;
}
</style>
