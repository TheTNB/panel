<script setup lang="ts">
import VersionModal from '@/views/app/VersionModal.vue'

import { NButton, NDataTable, NPopconfirm, NSwitch } from 'naive-ui'
import { useI18n } from 'vue-i18n'

import TheIcon from '@/components/custom/TheIcon.vue'
import { router } from '@/router'
import { renderIcon } from '@/utils'
import type { App } from '@/views/app/types'
import app from '../../api/panel/app'

const { t } = useI18n()

const versionModalShow = ref(false)
const versionModalOperation = ref('安装')
const versionModalInfo = ref<App>({} as App)

const columns: any = [
  {
    key: 'icon',
    fixed: 'left',
    width: 80,
    align: 'center',
    render(row: any) {
      return h(TheIcon, {
        icon: row.icon,
        size: 26
      })
    }
  },
  {
    title: t('appIndex.columns.name'),
    key: 'name',
    width: 300,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: t('appIndex.columns.description'),
    key: 'description',
    minWidth: 300,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: t('appIndex.columns.installedVersion'),
    key: 'installed_version',
    width: 100,
    ellipsis: { tooltip: true }
  },
  {
    title: t('appIndex.columns.show'),
    key: 'show',
    width: 100,
    align: 'center',
    render(row: any) {
      return h(NSwitch, {
        size: 'small',
        rubberBand: false,
        value: row.show,
        onUpdateValue: () => handleShowChange(row)
      })
    }
  },
  {
    title: t('appIndex.columns.actions'),
    key: 'actions',
    width: 280,
    align: 'center',
    hideInExcel: true,
    render(row: any) {
      return [
        row.installed && row.update_exist
          ? h(
              NPopconfirm,
              {
                onPositiveClick: () => handleUpdate(row.slug)
              },
              {
                default: () => {
                  return t('appIndex.confirm.update', { app: row.name })
                },
                trigger: () => {
                  return h(
                    NButton,
                    {
                      size: 'small',
                      type: 'warning'
                    },
                    {
                      default: () => t('appIndex.buttons.update'),
                      icon: renderIcon('material-symbols:arrow-circle-up-outline-rounded', {
                        size: 14
                      })
                    }
                  )
                }
              }
            )
          : null,
        row.installed
          ? h(
              NButton,
              {
                size: 'small',
                type: 'success',
                onClick: () => handleManage(row.slug)
              },
              {
                default: () => t('appIndex.buttons.manage'),
                icon: renderIcon('material-symbols:settings-outline', { size: 14 })
              }
            )
          : null,
        row.installed
          ? h(
              NPopconfirm,
              {
                onPositiveClick: () => handleUninstall(row.slug)
              },
              {
                default: () => {
                  return t('appIndex.confirm.uninstall', { app: row.name })
                },
                trigger: () => {
                  return h(
                    NButton,
                    {
                      size: 'small',
                      type: 'error',
                      style: 'margin-left: 15px;'
                    },
                    {
                      default: () => t('appIndex.buttons.uninstall'),
                      icon: renderIcon('material-symbols:delete-outline', { size: 14 })
                    }
                  )
                }
              }
            )
          : null,
        !row.installed
          ? h(
              NButton,
              {
                size: 'small',
                type: 'info',
                onClick: () => {
                  versionModalShow.value = true
                  versionModalOperation.value = '安装'
                  versionModalInfo.value = row
                }
              },
              {
                default: () => t('appIndex.buttons.install'),
                icon: renderIcon('material-symbols:download-rounded', { size: 14 })
              }
            )
          : null
      ]
    }
  }
]

const apps = ref<App[]>([] as App[])

const selectedRowKeys = ref<any>([])

const pagination = reactive({
  page: 1,
  pageCount: 1,
  pageSize: 20,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [20, 50, 100, 200]
})

const handleShowChange = (row: any) => {
  app.updateShow(row.slug, !row.show).then(() => {
    window.$message.success(t('appIndex.alerts.setup'))
    row.show = !row.show
  })
}

const handleUpdate = (slug: string) => {
  app.update(slug).then(() => {
    window.$message.success(t('appIndex.alerts.update'))
  })
}

const handleUninstall = (slug: string) => {
  app.uninstall(slug).then(() => {
    window.$message.success(t('appIndex.alerts.uninstall'))
  })
}

const handleManage = (slug: string) => {
  router.push({ name: 'apps-' + slug + '-index' })
}

const handleUpdateCache = () => {
  app.updateCache().then(() => {
    window.$message.success(t('appIndex.alerts.cache'))
    onPageChange(1)
  })
}

const getAppList = async (page: number, limit: number) => {
  const { data } = await app.list(page, limit)
  return data
}

const onChecked = (rowKeys: any) => {
  selectedRowKeys.value = rowKeys
}

const onPageChange = (page: number) => {
  pagination.page = page
  getAppList(page, pagination.pageSize).then((res) => {
    apps.value = res.items
    pagination.itemCount = res.total
    pagination.pageCount = res.total / pagination.pageSize + 1
  })
}

const onPageSizeChange = (pageSize: number) => {
  pagination.pageSize = pageSize
  onPageChange(1)
}

onMounted(() => {
  onPageChange(pagination.page)
})
</script>

<template>
  <common-page show-footer>
    <template #action>
      <div flex items-center>
        <n-button class="ml-16" type="primary" @click="handleUpdateCache">
          <TheIcon :size="18" icon="material-symbols:refresh" />
          {{ $t('appIndex.buttons.updateCache') }}
        </n-button>
      </div>
    </template>
    <n-flex vertical>
      <n-alert type="warning">{{ $t('appIndex.alerts.warning') }}</n-alert>
      <n-data-table
        striped
        remote
        :scroll-x="1000"
        :loading="false"
        :columns="columns"
        :data="apps"
        :row-key="(row: any) => row.slug"
        :pagination="pagination"
        @update:checked-row-keys="onChecked"
        @update:page="onPageChange"
        @update:page-size="onPageSizeChange"
      />
      <version-modal
        v-model:show="versionModalShow"
        v-model:operation="versionModalOperation"
        v-model:info="versionModalInfo"
      />
    </n-flex>
  </common-page>
</template>
