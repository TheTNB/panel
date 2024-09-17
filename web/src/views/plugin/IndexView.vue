<script setup lang="ts">
import type { Plugin } from '@/views/plugin/types'
import { NButton, NDataTable, NPopconfirm, NSwitch } from 'naive-ui'
import plugin from '@/api/panel/plugin'
import { renderIcon } from '@/utils'
import { router } from '@/router'
import { useI18n } from 'vue-i18n'

const { t } = useI18n()

const columns: any = [
  { type: 'selection', fixed: 'left' },
  {
    title: t('pluginIndex.columns.name'),
    key: 'name',
    width: 150,
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: t('pluginIndex.columns.description'),
    key: 'description',
    resizable: true,
    ellipsis: { tooltip: true }
  },
  {
    title: t('pluginIndex.columns.installedVersion'),
    key: 'installed_version',
    width: 100,
    ellipsis: { tooltip: true }
  },
  {
    title: t('pluginIndex.columns.version'),
    key: 'version',
    width: 100,
    ellipsis: { tooltip: true }
  },
  {
    title: t('pluginIndex.columns.show'),
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
    title: t('pluginIndex.columns.actions'),
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
                  return t('pluginIndex.confirm.update', { plugin: row.name })
                },
                trigger: () => {
                  return h(
                    NButton,
                    {
                      size: 'small',
                      type: 'warning'
                    },
                    {
                      default: () => t('pluginIndex.buttons.update'),
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
                default: () => t('pluginIndex.buttons.manage'),
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
                  return t('pluginIndex.confirm.uninstall', { plugin: row.name })
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
                      default: () => t('pluginIndex.buttons.uninstall'),
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
                  return t('pluginIndex.confirm.install', { plugin: row.name })
                },
                trigger: () => {
                  return h(
                    NButton,
                    {
                      size: 'small',
                      type: 'info'
                    },
                    {
                      default: () => t('pluginIndex.buttons.install'),
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

const plugins = ref<Plugin[]>([] as Plugin[])

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
  plugin.updateShow(row.slug, !row.show).then(() => {
    window.$message.success(t('pluginIndex.alerts.setup'))
    row.show = !row.show
  })
}

const handleInstall = (slug: string) => {
  plugin.install(slug).then(() => {
    window.$message.success(t('pluginIndex.alerts.install'))
  })
}

const handleUpdate = (slug: string) => {
  plugin.update(slug).then(() => {
    window.$message.success(t('pluginIndex.alerts.update'))
  })
}

const handleUninstall = (slug: string) => {
  plugin.uninstall(slug).then(() => {
    window.$message.success(t('pluginIndex.alerts.uninstall'))
  })
}

const handleManage = (slug: string) => {
  router.push({ name: 'plugins-' + slug + '-index' })
}

const getPluginList = async (page: number, limit: number) => {
  const { data } = await plugin.list(page, limit)
  return data
}

const onChecked = (rowKeys: any) => {
  selectedRowKeys.value = rowKeys
}

const onPageChange = (page: number) => {
  pagination.page = page
  getPluginList(page, pagination.pageSize).then((res) => {
    plugins.value = res.items
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
      <n-alert type="info">{{ $t('pluginIndex.alerts.info') }}</n-alert>
      <n-alert type="warning">{{ $t('pluginIndex.alerts.warning') }}</n-alert>
      <n-data-table
        striped
        remote
        :loading="false"
        :columns="columns"
        :data="plugins"
        :row-key="(row: any) => row.slug"
        :pagination="pagination"
        @update:checked-row-keys="onChecked"
        @update:page="onPageChange"
        @update:page-size="onPageSizeChange"
      />
    </n-space>
  </common-page>
</template>
