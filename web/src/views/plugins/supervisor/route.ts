import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'supervisor',
  path: '/plugins/supervisor',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'plugins-supervisor-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'SuperVisor',
        icon: 'mdi:monitor-dashboard',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
