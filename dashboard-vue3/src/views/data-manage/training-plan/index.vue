<template>
  <div class="training-plan-page art-full-height">
    <ElCard class="art-table-card">
      <!-- 标题区域 -->
      <div class="flex items-center justify-between mb-4">
        <div class="flex items-center">
          <h2 class="text-lg font-medium m-0">培养计划</h2>
        </div>
        <ElDropdown trigger="click">
          <ElButton link icon="ri:more-2-fill" class="text-xl" />
          <template #dropdown>
            <ElDropdownMenu>
              <ElDropdownItem>批量导出</ElDropdownItem>
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
      <ArtTableHeader :showZebra="false" v-model:columns="columnChecks">
        <template #left>
          <ElButton v-ripple @click="handleAdd">新增培养计划</ElButton>
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

  defineOptions({ name: 'TrainingPlan' })

  interface TrainingPlanItem {
    id: number
    planName: string
    totalCredits: number
    studentCount: number
    status: 'enabled' | 'disabled'
  }

  const initialSearchState = {
    planName: ''
  }
  const ArtButtonTable = resolveComponent('ArtButtonTable')

  const formFilters = reactive({ ...initialSearchState })
  const appliedFilters = reactive({ ...initialSearchState })

  const formItems = computed(() => [
    {
      label: '培养计划',
      key: 'planName',
      type: 'input',
      props: {
        placeholder: '请输入培养计划名称',
        clearable: true
      }
    }
  ])

  const tableData = ref<TrainingPlanItem[]>([
    {
      id: 1,
      planName: '2026级人工智能本科培养计划',
      totalCredits: 160,
      studentCount: 120,
      status: 'enabled'
    },
    {
      id: 2,
      planName: '2026级软件工程本科培养计划',
      totalCredits: 158,
      studentCount: 180,
      status: 'enabled'
    },
    {
      id: 3,
      planName: '2026级数字媒体艺术培养计划',
      totalCredits: 152,
      studentCount: 90,
      status: 'disabled'
    },
    {
      id: 4,
      planName: '2026级产品设计培养计划',
      totalCredits: 150,
      studentCount: 76,
      status: 'enabled'
    }
  ])

  // 分页状态
  const pagination = reactive({
    current: 1,
    size: 10,
    total: tableData.value.length
  })

  const filteredTableData = computed(() => {
    const keyword = appliedFilters.planName.trim().toLowerCase()
    let result = tableData.value

    if (keyword) {
      result = result.filter((item) => item.planName.toLowerCase().includes(keyword))
    }

    // 简单前端分页
    const start = (pagination.current - 1) * pagination.size
    const end = start + pagination.size
    return result.slice(start, end)
  })

  const { columnChecks, columns } = useTableColumns<TrainingPlanItem>(() => [
    { type: 'index', width: 70, label: '序号' },
    { prop: 'planName', label: '培养计划名称', minWidth: 240 },
    { prop: 'totalCredits', label: '总学分', width: 120 },
    { prop: 'studentCount', label: '学生人数', width: 120 },
    {
      prop: 'status',
      label: '状态',
      width: 120,
      formatter: (row) =>
        h(ElTag, { type: row.status === 'enabled' ? 'success' : 'info' }, () =>
          row.status === 'enabled' ? '启用' : '停用'
        )
    },
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
    ElMessage.info('新增培养计划功能待接入')
  }

  const handleEdit = (row: TrainingPlanItem) => {
    ElMessage.info(`编辑：${row.planName}`)
  }

  const handleDelete = (row: TrainingPlanItem) => {
    ElMessage.warning(`删除：${row.planName}`)
  }
</script>
