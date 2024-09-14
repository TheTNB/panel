import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'cert',
  path: '/cert',
  component: Layout,
  meta: {
    order: 2
  },
  children: [
    {
      name: 'cert-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'certIndex.title',
        icon: 'mdi:certificate',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
