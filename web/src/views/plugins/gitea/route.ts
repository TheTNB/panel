import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'gitea',
  path: '/plugins/gitea',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'plugins-gitea-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'Gitea',
        icon: 'mdi:git',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
