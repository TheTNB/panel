<script setup lang="ts">
import file from '@/api/panel/file'
import EventBus from '@/utils/event'
import { NButton, NPopconfirm, NSpace, NTag } from 'naive-ui'

import type { DataTableColumns } from 'naive-ui'
import type { RowData } from 'naive-ui/es/data-table/src/interface'

const show = defineModel<boolean>('show', { type: Boolean, required: true })
const path = defineModel<string>('path', { type: String, required: true })
const keyword = defineModel<string>('keyword', { type: String, required: true })
const sub = defineModel<boolean>('sub', { type: Boolean, required: true })

const loading = ref(false)

const columns: DataTableColumns<RowData> = [
  {
    title: '名称',
    key: 'full',
    minWidth: 300,
    ellipsis: {
      tooltip: true
    }
  },
  {
    title: '大小',
    key: 'size',
    width: 80,
    render(row: any): any {
      return h(NTag, { type: 'info', size: 'small', bordered: false }, { default: () => row.size })
    }
  },
  {
    title: '修改时间',
    key: 'modify',
    width: 200,
    render(row: any): any {
      return h(
        NTag,
        { type: 'warning', size: 'small', bordered: false },
        { default: () => row.modify }
      )
    }
  },
  {
    title: '操作',
    key: 'action',
    width: 200,
    render(row) {
      return h(
        NSpace,
        {},
        {
          default: () => [
            h(
              NButton,
              {
                size: 'small',
                type: 'success',
                tertiary: true,
                onClick: () => {
                  navigator.clipboard.writeText(row.full)
                  window.$message.success('复制成功')
                }
              },
              {
                default: () => {
                  return '复制路径'
                }
              }
            ),
            h(
              NPopconfirm,
              {
                onPositiveClick: () => {
                  file.delete(row.full).then(() => {
                    window.$message.success('删除成功')
                    EventBus.emit('file:refresh')
                  })
                },
                onNegativeClick: () => {}
              },
              {
                default: () => {
                  return '确定删除吗？'
                },
                trigger: () => {
                  return h(
                    NButton,
                    {
                      size: 'small',
                      type: 'error',
                      tertiary: true
                    },
                    { default: () => '删除' }
                  )
                }
              }
            )
          ]
        }
      )
    }
  }
]

const data = ref<RowData[]>([])

const pagination = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 100,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [100, 200, 500, 1000, 1500, 2000, 5000]
})

const handlePageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  handlePageChange(1)
}

const handlePageChange = (page: number) => {
  search(page)
}

const search = async (page: number) => {
  loading.value = true
  await file
    .search(path.value, keyword.value, sub.value, page, pagination.pageSize!)
    .then((res) => {
      data.value = res.data.items
      pagination.itemCount = res.data.total
      pagination.pageCount = res.data.total / pagination.pageSize! + 1
    })
    .catch(() => {
      window.$message.error('搜索失败')
    })
  loading.value = false
}

watch(show, (value) => {
  if (value) {
    search(1)
  }
})
</script>

<template>
  <n-modal
    v-model:show="show"
    preset="card"
    :title="keyword + ' - 搜索结果'"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-data-table
      remote
      striped
      virtual-scroll
      size="small"
      :scroll-x="800"
      :columns="columns"
      :data="data"
      :loading="loading"
      :pagination="pagination"
      :row-key="(row: any) => row.full"
      max-height="60vh"
      @update:page="handlePageChange"
      @update:page-size="handlePageSizeChange"
    />
  </n-modal>
</template>

<style scoped lang="scss"></style>
