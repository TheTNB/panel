import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'ssh',
  path: '/ssh',
  component: Layout,
  meta: {
    order: 7
  },
  children: [
    {
      name: 'ssh-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'sshIndex.title',
        icon: 'mdi:console',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
