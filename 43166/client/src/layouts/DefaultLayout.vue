<template>
  <el-container class="default-layout">
    <el-aside width="220px" class="sidebar">
      <div class="logo">
        <el-icon :size="28"><Building /></el-icon>
        <span>工商注册平台</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        class="sidebar-menu"
        background-color="#304156"
        text-color="#bfcbd9"
        active-text-color="#409EFF"
        router
      >
        <el-menu-item v-for="item in menuItems" :key="item.path" :index="item.path">
          <el-icon v-if="item.icon">
            <component :is="item.icon" />
          </el-icon>
          <template #title>{{ item.title }}</template>
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <div class="header-left">
          <el-breadcrumb separator="/">
            <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
            <el-breadcrumb-item v-if="currentRoute">{{ currentRoute }}</el-breadcrumb-item>
          </el-breadcrumb>
        </div>
        <div class="header-right">
          <el-badge :value="unreadCount" :hidden="unreadCount === 0" class="notification-badge">
            <el-icon :size="20" class="cursor-pointer" @click="goToNotifications">
              <Bell />
            </el-icon>
          </el-badge>
          <el-dropdown trigger="click" @command="handleCommand">
            <div class="user-info">
              <el-avatar :size="32" :src="userStore.userInfo?.avatar">
                {{ userStore.userInfo?.realName?.charAt(0) || 'U' }}
              </el-avatar>
              <span class="username">{{ userStore.userInfo?.realName || '用户' }}</span>
              <el-icon><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>
                  个人中心
                </el-dropdown-item>
                <el-dropdown-item command="logout" divided>
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main class="main-content">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { notificationApi } from '@/api/notification'
import { UserRole } from '@/types'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const unreadCount = ref(0)

const activeMenu = computed(() => route.path)
const currentRoute = computed(() => route.meta.title as string)

const menuItems = computed(() => {
  const allMenus = [
    { path: '/dashboard', title: '工作台', icon: 'HomeFilled', roles: ['admin', 'entrepreneur', 'agent'] },
    { path: '/applications', title: '注册申请', icon: 'Document', roles: ['admin', 'entrepreneur', 'agent'] },
    { path: '/fees', title: '费用管理', icon: 'Money', roles: ['admin', 'entrepreneur'] },
    { path: '/notifications', title: '消息中心', icon: 'Bell', roles: ['admin', 'entrepreneur', 'agent'] },
    { path: '/admin/agents', title: '专员管理', icon: 'UserFilled', roles: ['admin'] },
    { path: '/admin/fee-standards', title: '费用标准', icon: 'Setting', roles: ['admin'] },
    { path: '/admin/discounts', title: '优惠策略', icon: 'Discount', roles: ['admin'] },
    { path: '/admin/notification-templates', title: '通知模板', icon: 'MessageBox', roles: ['admin'] },
    { path: '/admin/statistics', title: '统计分析', icon: 'DataLine', roles: ['admin'] },
    { path: '/admin/exports', title: '数据导出', icon: 'Download', roles: ['admin'] }
  ]

  return allMenus.filter(item => item.roles.includes(userStore.userRole || ''))
})

const fetchUnreadCount = async () => {
  try {
    const res = await notificationApi.getUnreadCount()
    unreadCount.value = res?.count || 0
  } catch (error) {
    console.error('获取未读消息数量失败:', error)
  }
}

const goToNotifications = () => {
  router.push('/notifications')
}

const handleCommand = (command: string) => {
  if (command === 'profile') {
    router.push('/profile')
  } else if (command === 'logout') {
    userStore.logout()
    router.push('/login')
  }
}

onMounted(() => {
  if (!userStore.userInfo) {
    userStore.fetchUserInfo()
  }
  fetchUnreadCount()
})
</script>

<style scoped>
.default-layout {
  height: 100vh;
}

.sidebar {
  background-color: #304156;
  transition: width 0.3s;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: #fff;
  font-size: 16px;
  font-weight: 600;
  border-bottom: 1px solid #1f2d3d;
}

.sidebar-menu {
  border-right: none;
  height: calc(100vh - 60px);
}

.sidebar-menu:not(.el-menu--collapse) {
  width: 220px;
}

.header {
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.header-left {
  flex: 1;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 20px;
}

.notification-badge {
  cursor: pointer;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
  padding: 0 10px;
  border-radius: 4px;
  transition: background 0.3s;
}

.user-info:hover {
  background: #f5f7fa;
}

.username {
  font-size: 14px;
  color: #606266;
}

.main-content {
  background-color: #f0f2f5;
  padding: 20px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
