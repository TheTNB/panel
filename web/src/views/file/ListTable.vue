<script setup lang="ts">
import type { DataTableColumns } from 'naive-ui'
import { NButton, NInput, NPopconfirm, NPopselect, NSpace } from 'naive-ui'
import type { RowData } from 'naive-ui/es/data-table/src/interface'
import file from '@/api/panel/file'
import TheIcon from '@/components/custom/TheIcon.vue'
import { checkName, checkPath, getExt, getFilename, getIconByExt, isArchive } from '@/utils/file'
import EditModal from '@/views/file/EditModal.vue'
import EventBus from '@/utils/event'
import type { Marked } from '@/views/file/types'

const loading = ref(false)
const path = defineModel<string>('path', { type: String, required: true })
const selected = defineModel<any[]>('selected', { type: Array, default: () => [] })
const marked = defineModel<Marked[]>('marked', { type: Array, default: () => [] })
const archive = defineModel<boolean>('archive', { type: Boolean, required: true })
const permission = defineModel<boolean>('permission', { type: Boolean, required: true })
const editorModal = ref(false)
const editorFile = ref('')

const renameModal = ref(false)
const renameModel = ref({
  source: '',
  target: ''
})
const unArchiveModal = ref(false)
const unArchiveModel = ref({
  path: '',
  file: ''
})

const columns: DataTableColumns<RowData> = [
  {
    type: 'selection',
    fixed: 'left'
  },
  {
    title: '名称',
    key: 'name',
    width: '180',
    ellipsis: {
      tooltip: true
    },
    render(row) {
      let icon = 'bi:file-earmark'
      if (row.dir) {
        icon = 'bi:folder'
      } else {
        icon = getIconByExt(getExt(row.name))
      }

      return h(
        NSpace,
        {
          onClick: () => {
            if (row.dir) {
              path.value = row.full
            }
          }
        },
        () => [
          h(TheIcon, { icon }),
          h(
            'p',
            {},
            {
              default: () => {
                if (row.symlink) {
                  return row.name + ' -> ' + row.link
                } else {
                  return row.name
                }
              }
            }
          )
        ]
      )
    }
  },
  { title: '权限', key: 'mode', width: '80' },
  { title: '所有者', key: 'owner', width: '80' },
  { title: '组', key: 'group', width: '80' },
  { title: '大小', key: 'size', width: '80' },
  { title: '修改时间', key: 'modify', width: '150' },
  {
    title: '操作',
    key: 'action',
    width: '340',
    render(row) {
      return h(
        NSpace,
        {},
        {
          default: () => [
            h(
              NButton,
              {
                tertiary: true,
                size: 'small',
                onClick: () => {
                  if (!row.dir && !row.symlink) {
                    editorFile.value = row.full
                    editorModal.value = true
                  } else {
                    path.value = row.full
                  }
                }
              },
              {
                default: () => {
                  if (!row.dir && !row.symlink) {
                    return '编辑'
                  } else {
                    return '打开'
                  }
                }
              }
            ),
            h(
              NButton,
              {
                tertiary: true,
                size: 'small',
                onClick: () => {
                  if (row.dir) {
                    selected.value = [row.full]
                    archive.value = true
                  } else {
                    window.open('/api/panel/file/download?path=' + encodeURIComponent(row.full))
                  }
                }
              },
              {
                default: () => {
                  if (row.dir) {
                    return '压缩'
                  } else {
                    return '下载'
                  }
                }
              }
            ),
            h(
              NButton,
              {
                tertiary: true,
                size: 'small',
                onClick: () => {
                  renameModel.value.source = getFilename(row.name)
                  renameModel.value.target = getFilename(row.name)
                  renameModal.value = true
                }
              },
              { default: () => '重命名' }
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
                      tertiary: true,
                      size: 'small',
                      type: 'error'
                    },
                    { default: () => '删除' }
                  )
                }
              }
            ),
            h(
              NPopselect,
              {
                options: [
                  { label: '复制', value: 'copy' },
                  { label: '移动', value: 'move' },
                  { label: '权限', value: 'permission' },
                  { label: '压缩', value: 'archive' },
                  { label: '解压', value: 'unarchive', disabled: !isArchive(row.name) }
                ],
                onUpdateValue: (value) => {
                  switch (value) {
                    case 'copy':
                      marked.value = [
                        {
                          name: row.name,
                          source: row.full,
                          type: 'copy'
                        }
                      ]
                      window.$message.success('标记成功，请前往目标路径粘贴')
                      break
                    case 'move':
                      marked.value = [
                        {
                          name: row.name,
                          source: row.full,
                          type: 'move'
                        }
                      ]
                      window.$message.success('标记成功，请前往目标路径粘贴')
                      break
                    case 'permission':
                      selected.value = [row.full]
                      permission.value = true
                      break
                    case 'archive':
                      selected.value = [row.full]
                      archive.value = true
                      break
                    case 'unarchive':
                      unArchiveModel.value.file = row.full
                      unArchiveModel.value.path = path.value
                      unArchiveModal.value = true
                      break
                  }
                }
              },
              {
                default: () => {
                  return h(
                    NButton,
                    {
                      tertiary: true,
                      size: 'small'
                    },
                    { default: () => '更多' }
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
  pageSize: 20,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [20, 50, 100, 200, 500, 1000]
})

const handlePageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  handlePageChange(1)
}

const handlePageChange = async (page: number) => {
  loading.value = true
  await getList(path.value, page, pagination.pageSize!).finally(() => {
    loading.value = false
  })
}

const handleRefresh = async () => {
  loading.value = true
  await getList(path.value, pagination.page, pagination.pageSize!).finally(() => {
    loading.value = false
  })
}

const getList = async (path: string, page: number, limit: number) => {
  await file.list(path, page, limit).then((res) => {
    data.value = res.data.items
    pagination.page = page
    pagination.itemCount = res.data.total
    pagination.pageCount = res.data.total / pagination.pageSize! + 1
  })
}

const handleRename = () => {
  const source = path.value + '/' + renameModel.value.source
  const target = path.value + '/' + renameModel.value.target
  if (!checkName(renameModel.value.source) || !checkName(renameModel.value.target)) {
    window.$message.error('名称不合法')
    console.log(source, target)
    return
  }

  file.move(source, target).then(() => {
    window.$message.success('重命名成功')
    renameModal.value = false
    EventBus.emit('file:refresh')
  })
}

const handleUnArchive = () => {
  // 移除首位的 / 去检测
  if (
    !unArchiveModel.value.path.startsWith('/') ||
    !checkPath(unArchiveModel.value.path.slice(1))
  ) {
    window.$message.error('路径不合法')
    return
  }
  const message = window.$message.loading('正在解压中...', {
    duration: 0
  })
  file
    .unArchive(unArchiveModel.value.file, unArchiveModel.value.path)
    .then(() => {
      message?.destroy()
      window.$message.success('解压成功')
      unArchiveModal.value = false
      EventBus.emit('file:refresh')
    })
    .catch(() => {
      message?.destroy()
      window.$message.error('解压失败')
    })
}

const onChecked = (rowKeys: any) => {
  selected.value = rowKeys
}

onMounted(() => {
  watch(
    path,
    () => {
      handlePageChange(1)
      EventBus.emit('push-history', path.value)
    },
    { immediate: true }
  )
  EventBus.on('file:refresh', handleRefresh)
})

onUnmounted(() => {
  EventBus.off('file:refresh', handleRefresh)
})
</script>

<template>
  <n-data-table
    size="small"
    :columns="columns"
    :data="data"
    :loading="loading"
    :pagination="pagination"
    :row-key="(row: any) => row.full"
    :checked-row-keys="selected"
    max-height="60vh"
    remote
    striped
    virtual-scroll
    @update:page="handlePageChange"
    @update:page-size="handlePageSizeChange"
    @update:checked-row-keys="onChecked"
  />

  <edit-modal v-model:show="editorModal" v-model:file="editorFile" />
  <n-modal
    v-model:show="renameModal"
    preset="card"
    :title="'重命名 - ' + renameModel.source"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-flex vertical>
      <n-form>
        <n-form-item label="新名称">
          <n-input v-model:value="renameModel.target" />
        </n-form-item>
      </n-form>
      <n-button type="primary" @click="handleRename">保存</n-button>
    </n-flex>
  </n-modal>
  <n-modal
    v-model:show="unArchiveModal"
    preset="card"
    :title="'解压缩 - ' + unArchiveModel.file"
    style="width: 60vw"
    size="huge"
    :bordered="false"
    :segmented="false"
  >
    <n-flex vertical>
      <n-form>
        <n-form-item label="解压到">
          <n-input v-model:value="unArchiveModel.path" />
        </n-form-item>
      </n-form>
      <n-button type="primary" @click="handleUnArchive">解压</n-button>
    </n-flex>
  </n-modal>
</template>

<style scoped lang="scss"></style>
