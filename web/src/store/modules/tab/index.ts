import { defineStore } from 'pinia'
import { activeTab, tabs, WITHOUT_TAB_PATHS } from './helpers'
import { router } from '@/router'
import { setSession } from '@/utils'

export interface TabItem {
  name: string
  path: string
  title?: string
}

export const useTabStore = defineStore('tab', {
  state() {
    return {
      tabs: <Array<TabItem>>tabs || [],
      activeTab: <string>activeTab || ''
    }
  },
  actions: {
    setActiveTab(path: string) {
      this.activeTab = path
      setSession('activeTab', path)
    },
    setTabs(tabs: Array<TabItem>) {
      this.tabs = tabs
      setSession('tabs', tabs)
    },
    addTab(tab: TabItem) {
      this.setActiveTab(tab.path)
      if (WITHOUT_TAB_PATHS.includes(tab.path) || this.tabs.some((item) => item.path === tab.path))
        return
      this.setTabs([...this.tabs, tab])
    },
    removeTab(path: string) {
      if (path === this.activeTab) {
        const activeIndex = this.tabs.findIndex((item) => item.path === path)
        if (activeIndex > 0) router.push(this.tabs[activeIndex - 1].path)
        else router.push(this.tabs[activeIndex + 1].path)
      }
      this.setTabs(this.tabs.filter((tab) => tab.path !== path))
    },
    removeOther(curPath: string) {
      this.setTabs(this.tabs.filter((tab) => tab.path === curPath))
      if (curPath !== this.activeTab) router.push(this.tabs[this.tabs.length - 1].path)
    },
    removeLeft(curPath: string) {
      const curIndex = this.tabs.findIndex((item) => item.path === curPath)
      const filterTabs = this.tabs.filter((item, index) => index >= curIndex)
      this.setTabs(filterTabs)
      if (!filterTabs.find((item) => item.path === this.activeTab))
        router.push(filterTabs[filterTabs.length - 1].path)
    },
    removeRight(curPath: string) {
      const curIndex = this.tabs.findIndex((item) => item.path === curPath)
      const filterTabs = this.tabs.filter((item, index) => index <= curIndex)
      this.setTabs(filterTabs)
      if (!filterTabs.find((item) => item.path === this.activeTab))
        router.push(filterTabs[filterTabs.length - 1].path)
    },
    resetTabs() {
      this.setTabs([])
      this.setActiveTab('')
    }
  }
})
