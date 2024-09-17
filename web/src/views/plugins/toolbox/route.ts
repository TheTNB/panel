import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'toolbox',
  path: '/plugins/toolbox',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'plugins-toolbox-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: '系统工具箱',
        icon: 'mdi:tools',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
