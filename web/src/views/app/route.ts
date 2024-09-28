import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'app',
  path: '/app',
  component: Layout,
  meta: {
    order: 8
  },
  children: [
    {
      name: 'app-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'appIndex.title',
        icon: 'mdi:puzzle-outline',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
