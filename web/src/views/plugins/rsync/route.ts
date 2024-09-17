import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'rsync',
  path: '/plugins/rsync',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'plugins-rsync-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'Rsync',
        icon: 'mdi:sync',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
