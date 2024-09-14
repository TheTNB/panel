import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'task',
  path: '/task',
  component: Layout,
  meta: {
    order: 9
  },
  children: [
    {
      name: 'task-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'taskIndex.title',
        icon: 'mdi:archive-sync-outline',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
