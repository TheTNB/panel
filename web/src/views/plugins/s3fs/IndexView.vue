<script setup lang="ts">
import { NButton, NDataTable, NInput, NPopconfirm } from 'naive-ui'
import s3fs from '@/api/plugins/s3fs'
import { renderIcon } from '@/utils'
import type { S3fs } from '@/views/plugins/s3fs/types'

const addMountModal = ref(false)

const addMountModel = ref({
  ak: '',
  sk: '',
  bucket: '',
  url: '',
  path: ''
})

const pagination = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 15,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [15, 30, 50, 100]
})

const columns: any = [
  {
    title: '挂载路径',
    key: 'path',
    fixed: 'left',
    width: 250,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  { title: 'Bucket', key: 'bucket', resizable: true, ellipsis: { tooltip: true } },
  {
    title: '操作',
    key: 'actions',
    width: 240,
    align: 'center',
    fixed: 'right',
    hideInExcel: true,
    render(row: any) {
      return [
        h(
          NPopconfirm,
          {
            onPositiveClick: () => handleDeleteMount(row.id)
          },
          {
            default: () => {
              return '确定删除挂载' + row.path + '吗？'
            },
            trigger: () => {
              return h(
                NButton,
                {
                  size: 'small',
                  type: 'error'
                },
                {
                  default: () => '卸载',
                  icon: renderIcon('material-symbols:delete-outline', { size: 14 })
                }
              )
            }
          }
        )
      ]
    }
  }
]

const mounts = ref<S3fs[]>([] as S3fs[])

const getMounts = async (page: number, limit: number) => {
  const { data } = await s3fs.list(page, limit)
  return data
}

const onPageChange = (page: number) => {
  pagination.page = page
  getMounts(page, pagination.pageSize).then((res) => {
    mounts.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
  })
}

const onPageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  onPageChange(1)
}

const handleAddMount = async () => {
  await s3fs.add(addMountModel.value)
  window.$message.success('添加成功')
  onPageChange(1)
  addMountModal.value = false
}

const handleDeleteMount = async (id: number) => {
  await s3fs.delete(id)
  window.$message.success('删除成功')
  onPageChange(1)
}

onMounted(() => {
  onPageChange(1)
})
</script>

<template>
  <common-page show-footer>
    <template #action>
      <n-button class="ml-16" type="primary" @click="addMountModal = true">
        <TheIcon :size="18" class="mr-5" icon="material-symbols:add" />
        添加挂载
      </n-button>
    </template>
    <n-card title="挂载列表" :segmented="true" rounded-10>
      <n-data-table
        striped
        remote
        :loading="false"
        :columns="columns"
        :data="mounts"
        :row-key="(row: any) => row.id"
        @update:page="onPageChange"
        @update:page-size="onPageSizeChange"
      />
    </n-card>
  </common-page>
  <n-modal v-model:show="addMountModal" title="添加挂载">
    <n-card closable @close="() => (addMountModal = false)" title="添加挂载" style="width: 60vw">
      <n-form :model="addMountModel">
        <n-form-item path="bucket" label="Bucket（腾讯云COS为: xxxx-用户ID）">
          <n-input
            v-model:value="addMountModel.bucket"
            type="text"
            @keydown.enter.prevent
            placeholder="输入Bucket名字"
          />
        </n-form-item>
        <n-form-item path="ak" label="AK">
          <n-input
            v-model:value="addMountModel.ak"
            type="text"
            @keydown.enter.prevent
            placeholder="输入AK密钥"
          />
        </n-form-item>
        <n-form-item path="sk" label="SK">
          <n-input
            v-model:value="addMountModel.sk"
            type="text"
            @keydown.enter.prevent
            placeholder="输入SK密钥"
          />
        </n-form-item>
        <n-form-item path="url" label="地域节点">
          <n-input
            v-model:value="addMountModel.url"
            type="text"
            @keydown.enter.prevent
            placeholder="输入地域节点的 URL"
          />
        </n-form-item>
        <n-form-item path="path" label="挂载目录">
          <n-input
            v-model:value="addMountModel.path"
            type="text"
            @keydown.enter.prevent
            placeholder="输入挂载目录"
          />
        </n-form-item>
      </n-form>
      <n-row :gutter="[0, 24]">
        <n-col :span="24">
          <n-button type="info" block @click="handleAddMount">提交</n-button>
        </n-col>
      </n-row>
    </n-card>
  </n-modal>
</template>
