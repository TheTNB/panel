import { router } from '@/router'
import { defineStore } from 'pinia'

export const WITHOUT_TAB_PATHS = ['/404', '/login']

export interface Tab {
  active: string
  tabs: Array<TabItem>
}

export interface TabItem {
  name: string
  path: string
  title: string
}

export const useTabStore = defineStore('tab', {
  state: (): Tab => {
    return {
      active: '',
      tabs: []
    }
  },
  actions: {
    setActiveTab(path: string) {
      this.active = path
    },
    setTabs(tabs: Array<TabItem>) {
      this.tabs = tabs
    },
    addTab(tab: TabItem) {
      this.setActiveTab(tab.path)
      if (WITHOUT_TAB_PATHS.includes(tab.path) || this.tabs.some((item) => item.path === tab.path))
        return
      this.setTabs([...this.tabs, tab])
    },
    removeTab(path: string) {
      if (path === this.active) {
        const activeIndex = this.tabs.findIndex((item) => item.path === path)
        if (activeIndex > 0) router.push(this.tabs[activeIndex - 1].path)
        else router.push(this.tabs[activeIndex + 1].path)
      }
      this.setTabs(this.tabs.filter((tab) => tab.path !== path))
    },
    removeOther(curPath: string) {
      this.setTabs(this.tabs.filter((tab) => tab.path === curPath))
      if (curPath !== this.active) router.push(this.tabs[this.tabs.length - 1].path)
    },
    removeLeft(curPath: string) {
      const curIndex = this.tabs.findIndex((item) => item.path === curPath)
      const filterTabs = this.tabs.filter((item, index) => index >= curIndex)
      this.setTabs(filterTabs)
      if (!filterTabs.find((item) => item.path === this.active))
        router.push(filterTabs[filterTabs.length - 1].path)
    },
    removeRight(curPath: string) {
      const curIndex = this.tabs.findIndex((item) => item.path === curPath)
      const filterTabs = this.tabs.filter((item, index) => index <= curIndex)
      this.setTabs(filterTabs)
      if (!filterTabs.find((item) => item.path === this.active))
        router.push(filterTabs[filterTabs.length - 1].path)
    },
    resetTabs() {
      this.setTabs([])
      this.setActiveTab('')
    }
  },
  persist: true
})
