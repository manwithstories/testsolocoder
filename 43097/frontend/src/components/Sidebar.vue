<template>
  <el-menu
    :default-active="activeMenu"
    :collapse="appStore.sidebarCollapsed"
    :collapse-transition="false"
    background-color="#304156"
    text-color="#bfcbd9"
    active-text-color="#409EFF"
    class="sidebar-menu"
    router
  >
    <template v-for="menu in filteredMenus" :key="menu.path">
      <el-sub-menu v-if="menu.children && menu.children.length > 0" :index="menu.path">
        <template #title>
          <el-icon><component :is="menu.icon" /></el-icon>
          <span>{{ menu.title }}</span>
        </template>
        <el-menu-item
          v-for="child in menu.children"
          :key="child.path"
          :index="child.path"
        >
          <el-icon><component :is="child.icon" /></el-icon>
          <span>{{ child.title }}</span>
        </el-menu-item>
      </el-sub-menu>
      <el-menu-item v-else :index="menu.path">
        <el-icon><component :is="menu.icon" /></el-icon>
        <template #title>{{ menu.title }}</template>
      </el-menu-item>
    </template>
  </el-menu>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useAppStore } from '@/store/app'
import { useUserStore } from '@/store/user'
import { UserRole } from '@/types'

const route = useRoute()
const appStore = useAppStore()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)

const allMenus = [
  {
    path: '/dashboard',
    title: '仪表盘',
    icon: 'Odometer',
    roles: [UserRole.ADMIN, UserRole.MANAGER, UserRole.RECEPTIONIST, UserRole.STAFF]
  },
  {
    path: '/booking',
    title: '预订管理',
    icon: 'Booking',
    roles: [UserRole.ADMIN, UserRole.MANAGER, UserRole.RECEPTIONIST, UserRole.STAFF],
    children: [
      {
        path: '/booking/list',
        title: '预订列表',
        icon: 'List'
      },
      {
        path: '/booking/create',
        title: '创建预订',
        icon: 'Plus'
      }
    ]
  },
  {
    path: '/checkin',
    title: '入住管理',
    icon: 'HomeFilled',
    roles: [UserRole.ADMIN, UserRole.MANAGER, UserRole.RECEPTIONIST],
    children: [
      {
        path: '/checkin/list',
        title: '入住列表',
        icon: 'List'
      },
      {
        path: '/checkin/create',
        title: '办理入住',
        icon: 'Plus'
      },
      {
        path: '/checkin/today',
        title: '今日入住',
        icon: 'Calendar'
      }
    ]
  },
  {
    path: '/room',
    title: '房间管理',
    icon: 'OfficeBuilding',
    roles: [UserRole.ADMIN, UserRole.MANAGER, UserRole.RECEPTIONIST],
    children: [
      {
        path: '/room/list',
        title: '房间列表',
        icon: 'List'
      },
      {
        path: '/room/type',
        title: '房型管理',
        icon: 'Menu'
      },
      {
        path: '/room/status',
        title: '房间状态',
        icon: 'View'
      }
    ]
  },
  {
    path: '/member',
    title: '会员管理',
    icon: 'UserFilled',
    roles: [UserRole.ADMIN, UserRole.MANAGER, UserRole.RECEPTIONIST],
    children: [
      {
        path: '/member/list',
        title: '会员列表',
        icon: 'List'
      },
      {
        path: '/member/create',
        title: '添加会员',
        icon: 'Plus'
      }
    ]
  },
  {
    path: '/payment',
    title: '支付管理',
    icon: 'Money',
    roles: [UserRole.ADMIN, UserRole.MANAGER, UserRole.RECEPTIONIST],
    children: [
      {
        path: '/payment/list',
        title: '支付记录',
        icon: 'List'
      }
    ]
  },
  {
    path: '/report',
    title: '报表统计',
    icon: 'DataLine',
    roles: [UserRole.ADMIN, UserRole.MANAGER],
    children: [
      {
        path: '/report/daily',
        title: '日报表',
        icon: 'Document'
      },
      {
        path: '/report/room',
        title: '房间报表',
        icon: 'PieChart'
      },
      {
        path: '/report/revenue',
        title: '营收报表',
        icon: 'TrendCharts'
      }
    ]
  },
  {
    path: '/system',
    title: '系统管理',
    icon: 'Setting',
    roles: [UserRole.ADMIN],
    children: [
      {
        path: '/system/user',
        title: '用户管理',
        icon: 'User'
      },
      {
        path: '/system/role',
        title: '角色管理',
        icon: 'Avatar'
      },
      {
        path: '/system/log',
        title: '操作日志',
        icon: 'Document'
      }
    ]
  },
  {
    path: '/profile',
    title: '个人中心',
    icon: 'User',
    roles: [UserRole.ADMIN, UserRole.MANAGER, UserRole.RECEPTIONIST, UserRole.STAFF]
  }
]

const filteredMenus = computed(() => {
  const role = userStore.userRole
  if (!role) return []

  return allMenus.filter(menu => {
    return menu.roles.includes(role)
  })
})
</script>

<style scoped lang="scss">
.sidebar-menu {
  border-right: none;
  height: 100%;
}

:deep(.el-menu) {
  border-right: none;
}

:deep(.el-sub-menu__title:hover),
:deep(.el-menu-item:hover) {
  background-color: #263445 !important;
}

:deep(.el-menu-item.is-active) {
  background-color: #409EFF !important;
  color: #fff !important;
}
</style>
