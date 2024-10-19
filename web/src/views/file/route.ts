import type { RouteType } from '~/types/router'

const Layout = () => import('@/layout/IndexView.vue')

export default {
  name: 'file',
  path: '/file',
  component: Layout,
  meta: {
    order: 50
  },
  children: [
    {
      name: 'file-index',
      path: '',
      component: () => import('./IndexView.vue'),
      meta: {
        title: 'fileIndex.title',
        icon: 'mdi:folder-open-outline',
        role: ['admin'],
        requireAuth: true
      }
    }
  ]
} as RouteType
