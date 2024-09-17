import { getSession } from '@/utils'

export const activeTab = getSession('activeTab')
export const tabs = getSession('tabs')

export const WITHOUT_TAB_PATHS = ['/404', '/login']
