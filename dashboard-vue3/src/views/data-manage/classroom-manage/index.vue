<template>
  <div class="classroom-manage-page art-full-height">
    <ElCard class="art-table-card">
      <!-- 标题区域 -->
      <div class="flex items-center justify-between mb-4">
        <div class="flex items-center">
          <!-- <ElButton link icon="ri:arrow-left-line" class="mr-2 text-xl" @click="handleBack" /> -->
          <h2 class="text-lg font-medium m-0">教室管理</h2>
        </div>
        <ElDropdown trigger="click">
          <ElButton link icon="ri:more-2-fill" class="text-xl" />
          <template #dropdown>
            <ElDropdownMenu>
              <ElDropdownItem>导出数据</ElDropdownItem>
            </ElDropdownMenu>
          </template>
        </ElDropdown>
      </div>

      <!-- 搜索区域 -->
      <ArtSearchBar
        v-model="formFilters"
        :items="formItems"
        :showExpand="false"
        :noBorder="true"
        @search="handleSearch"
        @reset="handleReset"
        class="mb-4"
      />

      <!-- 按钮行与表格头配置 -->
      <ArtTableHeader v-model:columns="columnChecks">
        <template #left>
          <ElButton v-ripple @click="handleAdd">新增教室</ElButton>
          <ElButton v-ripple @click="handleSyncAccess">同步门禁设备</ElButton>
        </template>
      </ArtTableHeader>

      <!-- 表格区域 -->
      <ArtTable :columns="columns" :data="filteredTableData" />

      <!-- 分页区域 -->
      <div class="mt-4 flex justify-end">
        <ElPagination
          v-model:current-page="pagination.current"
          v-model:page-size="pagination.size"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
        />
      </div>
    </ElCard>
  </div>
</template>

<script setup lang="ts">
  import { useTableColumns } from '@/hooks/core/useTableColumns'
  import { ElTag, ElMessage } from 'element-plus'
  import { reactive, ref, computed, h, resolveComponent } from 'vue'

  defineOptions({ name: 'ClassroomManage' })

  interface ClassroomItem {
    id: number
    name: string
    type: string
    controller: string
    seatCount: number
    examSeatCount: number
    hasProjector: boolean
    hasNetwork: boolean
    remark: string
  }

  const ArtButtonTable = resolveComponent('ArtButtonTable')

  // 搜索条件
  const initialSearchState = {
    name: ''
  }

  const formFilters = reactive({ ...initialSearchState })
  const appliedFilters = reactive({ ...initialSearchState })

  const formItems = computed(() => [
    {
      label: '教室名称',
      key: 'name',
      type: 'input',
      props: {
        placeholder: '请输入教室名称',
        clearable: true
      }
    }
  ])

  // 模拟数据
  const tableData = ref<ClassroomItem[]>([
    {
      id: 1,
      name: '教一 101',
      type: '多媒体教室',
      controller: '智能中控 V1',
      seatCount: 60,
      examSeatCount: 30,
      hasProjector: true,
      hasNetwork: true,
      remark: '常规上课教室'
    },
    {
      id: 2,
      name: '理科楼 203',
      type: '计算机机房',
      controller: '智能中控 V2',
      seatCount: 50,
      examSeatCount: 50,
      hasProjector: true,
      hasNetwork: true,
      remark: '含 50 台学生电脑'
    },
    {
      id: 3,
      name: '主楼 405',
      type: '普通教室',
      controller: '无',
      seatCount: 40,
      examSeatCount: 20,
      hasProjector: false,
      hasNetwork: false,
      remark: '仅做自习使用'
    }
  ])

  // 分页状态
  const pagination = reactive({
    current: 1,
    size: 10,
    total: tableData.value.length
  })

  // 根据搜索条件过滤
  const filteredTableData = computed(() => {
    const keyword = appliedFilters.name.trim().toLowerCase()
    let result = tableData.value

    if (keyword) {
      result = result.filter((item) => item.name.toLowerCase().includes(keyword))
    }

    // 简单前端分页
    const start = (pagination.current - 1) * pagination.size
    const end = start + pagination.size
    return result.slice(start, end)
  })

  // 列配置
  const { columnChecks, columns } = useTableColumns<ClassroomItem>(() => [
    { type: 'index', width: 70, label: '序号' },
    { prop: 'name', label: '教室名称', minWidth: 120 },
    { prop: 'type', label: '类型', width: 120 },
    { prop: 'controller', label: '控制器', width: 120 },
    { prop: 'seatCount', label: '座位数', width: 100 },
    { prop: 'examSeatCount', label: '考试座位数', width: 120 },
    {
      prop: 'hasProjector',
      label: '投影',
      width: 80,
      formatter: (row) =>
        h(ElTag, { type: row.hasProjector ? 'success' : 'info' }, () =>
          row.hasProjector ? '有' : '无'
        )
    },
    {
      prop: 'hasNetwork',
      label: '网络',
      width: 80,
      formatter: (row) =>
        h(ElTag, { type: row.hasNetwork ? 'success' : 'info' }, () =>
          row.hasNetwork ? '有' : '无'
        )
    },
    { prop: 'remark', label: '备注', minWidth: 150 },
    {
      prop: 'operation',
      label: '操作',
      width: 140,
      fixed: 'right',
      formatter: (row) =>
        h('div', [
          h(ArtButtonTable, {
            type: 'edit',
            onClick: () => handleEdit(row)
          }),
          h(ArtButtonTable, {
            type: 'delete',
            onClick: () => handleDelete(row)
          })
        ])
    }
  ])

  const handleSearch = () => {
    Object.assign(appliedFilters, { ...formFilters })
    pagination.current = 1 // 搜索后回到第一页
  }

  const handleReset = () => {
    Object.assign(formFilters, { ...initialSearchState })
    Object.assign(appliedFilters, { ...initialSearchState })
    pagination.current = 1 // 重置后回到第一页
  }

  const handleAdd = () => {
    ElMessage.info('新增教室')
  }

  const handleSyncAccess = () => {
    ElMessage.success('已触发门禁设备同步')
  }

  const handleEdit = (row: ClassroomItem) => {
    ElMessage.info(`编辑：${row.name}`)
  }

  const handleDelete = (row: ClassroomItem) => {
    ElMessage.warning(`删除：${row.name}`)
  }
</script>
