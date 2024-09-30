<script setup lang="ts">
import { NButton, NDataTable, NPopconfirm, NSwitch } from 'naive-ui'
import { useI18n } from 'vue-i18n'

import { router } from '@/router'
import { renderIcon } from '@/utils'
import type { App } from '@/views/app/types'
import app from '../../api/panel/app'

const { t } = useI18n()

const columns: any = [
  { type: 'selection', fixed: 'left' },
  {
    title: t('appIndex.columns.name'),
    key: 'name',
    width: 200,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: t('appIndex.columns.description'),
    key: 'description',
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
    title: t('appIndex.columns.version'),
    key: 'version',
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
    fixed: 'right',
    hideInExcel: true,
    render(row: any) {
      return [
        row.installed && row.installed_version != row.version
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
        row.installed && row.installed_version == row.version
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
        row.installed && row.installed_version == row.version
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
              NPopconfirm,
              {
                onPositiveClick: () => handleInstall(row.slug)
              },
              {
                default: () => {
                  return t('appIndex.confirm.install', { app: row.name })
                },
                trigger: () => {
                  return h(
                    NButton,
                    {
                      size: 'small',
                      type: 'info'
                    },
                    {
                      default: () => t('appIndex.buttons.install'),
                      icon: renderIcon('material-symbols:download-rounded', { size: 14 })
                    }
                  )
                }
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
  pageSize: 15,
  itemCount: 0,
  showQuickJumper: true,
  showSizePicker: true,
  pageSizes: [15, 30, 50, 100]
})

const handleShowChange = (row: any) => {
  app.updateShow(row.slug, !row.show).then(() => {
    window.$message.success(t('appIndex.alerts.setup'))
    row.show = !row.show
  })
}

const handleInstall = (slug: string) => {
  app.install(slug).then(() => {
    window.$message.success(t('appIndex.alerts.install'))
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
    <n-space vertical>
      <n-alert type="info">{{ $t('appIndex.alerts.info') }}</n-alert>
      <n-alert type="warning">{{ $t('appIndex.alerts.warning') }}</n-alert>
      <n-data-table
        striped
        remote
        :loading="false"
        :columns="columns"
        :data="apps"
        :row-key="(row: any) => row.slug"
        :pagination="pagination"
        @update:checked-row-keys="onChecked"
        @update:page="onPageChange"
        @update:page-size="onPageSizeChange"
      />
    </n-space>
  </common-page>
</template>
