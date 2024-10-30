import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'docker',
  path: '/apps/docker',
  component: Layout,
  isHidden: true,
  children: [
    {
      name: 'apps-docker-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'Docker',
        icon: 'logos:docker-icon',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
