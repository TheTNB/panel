import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'plugin',
  path: '/plugin',
  component: Layout,
  meta: {
    order: 8
  },
  children: [
    {
      name: 'plugin-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'pluginIndex.title',
        icon: 'mdi:puzzle-outline',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
