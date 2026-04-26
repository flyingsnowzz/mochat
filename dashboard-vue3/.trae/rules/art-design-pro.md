# Art Design Pro 管理系统页面布局 Rules

用于统一后台管理页面的布局、交互与视觉规范，避免页面实现偏差。后续新页面默认按本规则落地。

## 1. 页面整体布局

- 页面容器统一使用 `art-full-height`。
- 标准结构固定为一个 `ElCard` 内部，按顺序包含：**标题**、**搜索区域**、**按钮行**、**表格**、**分页**。
- **标题部分**：
  - 左侧：**按需包含**一个返回按钮（如 `ri:arrow-left-line`，仅在下钻页面或子页面需要返回时存在），以及主标题、副标题。
  - 右侧：可包含使用 `...`（更多操作）触发的下拉菜单。
- 搜索区与表格区垂直间距保持 `16~20px`，不随意压缩。
- 表格头左侧第一主按钮统一为“新增Xxx”。
- 禁止将搜索项散落到表格头，所有查询条件必须在搜索模块内。

## 2. 搜索模块 Rules

- 默认仅放高频条件，低频条件可放“展开更多”。
- 每个搜索项都要有明确 `label`，不能只靠 placeholder。
- 按钮顺序固定：`查询` 在前，`重置` 在后。
- `重置` 必须同时重置：
  - 输入态（如 `formFilters`）
  - 应用态（如 `appliedFilters`）
  - 分页页码（重置到第一页）
- 文本输入默认开启 `clearable`，查询前建议 `trim()`。

## 3. 表格模块 Rules

- 表格列推荐顺序：`序号` -> 主业务字段 -> 辅助字段 -> `状态` -> `操作`。
- `序号` 使用索引列，宽度建议 `60~70`。
- `状态` 列统一使用标签：启用 `success`，停用 `info/danger`。
- `操作` 列固定在右侧，常规包含：编辑、删除；超过 3 个操作时折叠到“更多”。
- 表格必须具备：加载态、空态、错误反馈能力。
- 列宽优先保障主信息可读，避免布局抖动。

## 4. 分页 Rules（强制）

- **所有列表页表格底部都必须有分页组件**，不得省略。
- 推荐使用 `ArtTable` 内置分页能力（`pagination`）。
- 查询、重置后默认回到第一页。
- 删除当前页最后一条数据后，应自动回退到有效页码，避免空白页。
- 分页参数统一包含：`current`（当前页）、`size`（每页条数）、`total`（总数）。

## 5. 交互与行为 Rules

- 列表查询建议维护双状态：输入态 + 应用态，避免“边输入边抖动刷新”。
- 新增/编辑优先使用弹窗或抽屉，不在列表页直接堆长表单。
- 删除/停用等危险操作必须二次确认。
- 所有关键操作都要有成功/失败反馈（Message）。

## 6. 视觉一致性 Rules

- 卡片圆角统一 `8px`，轻阴影，不用重边框。
- 按钮高度统一 `40px`，避免同页多尺寸混用。
- 表格行高建议 `48px`。
- 页面文案默认中文，避免中英混排。
- 颜色保持低饱和，不使用强刺激色和夸张动效。

## 7. 推荐页面骨架（Vue）

```vue
<template>
  <div class="xxx-page art-full-height">
    <ElCard class="art-table-card">
      <!-- 标题区域 -->
      <div class="flex items-center justify-between mb-4">
        <div class="flex items-center">
          <!-- 返回按钮：仅在需要返回上级页面时存在 -->
          <ElButton
            v-if="showBack"
            link
            icon="ri:arrow-left-line"
            class="mr-2 text-xl"
            @click="handleBack"
          />
          <h2 class="text-lg font-medium m-0">主标题</h2>
          <span class="text-secondary text-sm ml-2">副标题</span>
        </div>
        <ElDropdown trigger="click">
          <ElButton link icon="ri:more-2-fill" class="text-xl" />
          <template #dropdown>
            <ElDropdownMenu>
              <ElDropdownItem>操作一</ElDropdownItem>
              <ElDropdownItem>操作二</ElDropdownItem>
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
          <ElButton v-ripple @click="handleAdd">新增Xxx</ElButton>
        </template>
      </ArtTableHeader>

      <!-- 表格区域 -->
      <ArtTable :columns="columns" :data="tableData" />

      <!-- 分页区域（如果 ArtTable 已内置分页，此处可根据实际情况使用组件内置功能） -->
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
```

## 8. 执行约定

- 新增后台列表页时，先套用本文件规则，再补业务字段。
- 代码评审以本文件为检查清单，不符合项需说明原因。
