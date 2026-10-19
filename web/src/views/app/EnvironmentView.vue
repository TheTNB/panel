<script setup lang="ts">
import { NButton, NDataTable, NFlex, NTag } from 'naive-ui'
import { useGettext } from 'vue3-gettext'

import environment from '@/api/panel/environment'
import { useConfirm } from '@/components/system/composables/useConfirm'
import { useTablePagination } from '@/components/system/composables/useTablePagination'
import { router } from '@/router'
import { renderLocalIcon } from '@/utils'
import CustomForm from '@/views/app/CustomForm.vue'

const { $gettext } = useGettext()
const { confirmDelete, confirmAction } = useConfirm()

const selectedType = ref<string>('')
const searchQuery = ref<string>('')

// 编译型环境的安装弹窗（含自定义编译参数输入）
const installModalShow = ref(false)
const installModalRow = ref<any>(null)
const installSubmitting = ref(false)
const customFormRef = ref<{ save: () => Promise<any> } | null>(null)

const { data: types } = useRequest(environment.types, {
  initialData: [],
})

const columns: any = [
  {
    key: 'icon',
    fixed: 'left',
    width: 80,
    align: 'center',
    render(row: any) {
      return renderLocalIcon('environment', row.type, { size: 26 })()
    },
  },
  {
    title: $gettext('Name'),
    key: 'name',
    width: 200,
    ellipsis: { tooltip: true },
  },
  {
    title: $gettext('Description'),
    key: 'description',
    minWidth: 300,
    ellipsis: { tooltip: true },
  },
  {
    title: $gettext('Latest Version'),
    key: 'version',
    width: 160,
    ellipsis: { tooltip: true },
  },
  {
    title: $gettext('Installed Version'),
    key: 'installed_version',
    width: 160,
    ellipsis: { tooltip: true },
  },
  {
    title: $gettext('Actions'),
    key: 'actions',
    width: 320,
    hideInExcel: true,
    render(row: any) {
      return h(NFlex, null, {
        default: () => {
          const items: any[] = []
          if (row.installed && row.has_update) {
            items.push(
              h(
                NButton,
                {
                  size: 'small',
                  type: 'warning',
                  onClick: async () => {
                    const ok = await confirmAction({
                      type: 'warning',
                      title: $gettext('Confirm Update'),
                      content: $gettext('Are you sure to update environment %{ environment }?', {
                        environment: row.name,
                      }),
                    })
                    if (ok) handleUpdate(row.type, row.slug)
                  },
                },
                { default: () => $gettext('Update') },
              ),
            )
          }
          if (row.installed) {
            items.push(
              h(
                NButton,
                {
                  size: 'small',
                  type: 'info',
                  onClick: () => handleManage(row.type, row.slug),
                },
                { default: () => $gettext('Manage') },
              ),
              h(
                NButton,
                {
                  size: 'small',
                  type: 'error',
                  onClick: async () => {
                    const ok = await confirmDelete({
                      title: $gettext('Confirm Uninstall'),
                      content: $gettext('Are you sure to uninstall environment %{ environment }?', {
                        environment: row.name,
                      }),
                      positiveText: $gettext('Uninstall'),
                      countdown: 5,
                    })
                    if (ok) handleUninstall(row.type, row.slug)
                  },
                },
                { default: () => $gettext('Uninstall') },
              ),
            )
          } else {
            items.push(
              h(
                NButton,
                {
                  size: 'small',
                  type: 'success',
                  onClick: async () => {
                    // 编译型环境走带编译参数输入的安装弹窗
                    if (row.custom_supported) {
                      installModalRow.value = row
                      installModalShow.value = true
                      return
                    }
                    const ok = await confirmAction({
                      type: 'info',
                      title: $gettext('Confirm Install'),
                      content: $gettext('Are you sure to install environment %{ environment }?', {
                        environment: row.name,
                      }),
                    })
                    if (ok) handleInstall(row.type, row.slug)
                  },
                },
                { default: () => $gettext('Install') },
              ),
            )
          }
          return items
        },
      })
    },
  },
]

const { loading, data, page, total, pageSize, refresh } = usePagination(
  (page, pageSize) =>
    environment.list(
      page,
      pageSize,
      selectedType.value || undefined,
      searchQuery.value || undefined,
    ),
  {
    initialData: { total: 0, list: [] },
    initialPageSize: 20,
    total: (res: any) => res.total,
    data: (res: any) => res.items,
    watchingStates: [selectedType, searchQuery],
  },
)
const pagination = useTablePagination(page, pageSize, total)

// 处理类型切换
const handleTypeChange = (type: string) => {
  selectedType.value = type
  page.value = 1
}

const handleInstall = (type: string, slug: string) => {
  useRequest(environment.install(type, slug)).onSuccess(() => {
    window.$message.success(
      $gettext('Task submitted, please check the progress in background tasks'),
    )
  })
}

// 编译型环境安装：先保存自定义参数再触发安装
const handleInstallWithCustom = async () => {
  const row = installModalRow.value
  if (!row) return
  installSubmitting.value = true
  try {
    if (customFormRef.value) {
      await customFormRef.value.save()
    }
  } catch {
    installSubmitting.value = false
    return
  }
  useRequest(environment.install(row.type, row.slug))
    .onSuccess(() => {
      installModalShow.value = false
      window.$message.success(
        $gettext('Task submitted, please check the progress in background tasks'),
      )
    })
    .onComplete(() => {
      installSubmitting.value = false
    })
}

const handleUpdate = (type: string, slug: string) => {
  useRequest(environment.update(type, slug)).onSuccess(() => {
    window.$message.success(
      $gettext('Task submitted, please check the progress in background tasks'),
    )
  })
}

const handleUninstall = (type: string, slug: string) => {
  useRequest(environment.uninstall(type, slug)).onSuccess(() => {
    window.$message.success(
      $gettext('Task submitted, please check the progress in background tasks'),
    )
  })
}

const handleManage = (type: string, slug: string) => {
  router.push({ name: 'environment-' + type, params: { slug } })
}

onMounted(() => {
  refresh()
})
</script>

<template>
  <n-flex vertical>
    <n-flex justify="space-between">
      <n-flex>
        <n-tag
          :type="selectedType === '' ? 'primary' : 'default'"
          :bordered="selectedType !== ''"
          class="cursor-pointer"
          @click="handleTypeChange('')"
        >
          {{ $gettext('All') }}
        </n-tag>
        <n-tag
          v-for="type in types"
          :key="type.value"
          :type="selectedType === type.value ? 'primary' : 'default'"
          :bordered="selectedType !== type.value"
          class="cursor-pointer"
          @click="handleTypeChange(type.value)"
        >
          {{ type.label }}
        </n-tag>
      </n-flex>
      <n-input
        v-model:value="searchQuery"
        :placeholder="$gettext('Search')"
        clearable
        class="md:!w-60"
      />
    </n-flex>
    <n-data-table
      v-model:page="page"
      v-model:pageSize="pageSize"
      striped
      remote
      :scroll-x="1280"
      :loading="loading"
      :columns="columns"
      :data="data"
      :row-key="(row: any) => row.slug"
      :pagination="pagination"
    />
  </n-flex>
  <app-modal
    v-model:show="installModalShow"
    :title="$gettext('Install') + ' ' + (installModalRow?.name ?? '')"
  >
    <n-flex vertical>
      <!-- 环境的自定义参数目录名为 type+slug(如 php83),与安装脚本约定一致 -->
      <custom-form
        v-if="installModalRow"
        ref="customFormRef"
        :slug="installModalRow.type + installModalRow.slug"
      />
      <n-button
        type="info"
        block
        :loading="installSubmitting"
        :disabled="installSubmitting"
        @click="handleInstallWithCustom"
      >
        {{ $gettext('Install') }}
      </n-button>
    </n-flex>
  </app-modal>
</template>

<style scoped lang="scss"></style>
