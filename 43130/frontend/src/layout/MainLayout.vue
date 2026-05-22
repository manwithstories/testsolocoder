<template>
  <el-container class="main-layout">
    <el-aside width="220px" class="sidebar">
      <div class="logo">
        <el-icon :size="28"><House /></el-icon>
        <span>婚礼策划系统</span>
      </div>
      <el-menu
        :default-active="activeMenu"
        router
        background-color="#001529"
        text-color="#fff"
        active-text-color="#ffd04b"
      >
        <template v-for="route in menuRoutes" :key="route.path">
          <el-menu-item :index="route.fullPath">
            <el-icon><component :is="route.meta.icon" /></el-icon>
            <span>{{ route.meta.title }}</span>
          </el-menu-item>
        </template>
      </el-menu>
    </el-aside>
    
    <el-container>
      <el-header class="header">
        <div class="header-left">
          <el-select
            v-if="weddings.length > 0"
            v-model="currentWeddingId"
            placeholder="选择婚礼"
            style="width: 200px"
            @change="handleWeddingChange"
          >
            <el-option
              v-for="wedding in weddings"
              :key="wedding.id"
              :label="wedding.title"
              :value="wedding.id"
            />
          </el-select>
        </div>
        <div class="header-right">
          <el-badge :value="unreadCount" :hidden="unreadCount === 0" class="notification-badge">
            <el-button :icon="Bell" circle @click="showNotifications = true" />
          </el-badge>
          <el-dropdown @command="handleCommand">
            <span class="user-info">
              <el-avatar :size="32" :src="userInfo?.avatar">
                {{ userInfo?.full_name?.charAt(0) || 'U' }}
              </el-avatar>
              <span class="username">{{ userInfo?.full_name }}</span>
              <el-icon><ArrowDown /></el-icon>
            </span>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
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

    <el-drawer v-model="showNotifications" title="通知消息" direction="rtl" size="360px">
      <div class="notification-header">
        <span>通知列表</span>
        <el-link type="primary" @click="markAllAsRead">全部已读</el-link>
      </div>
      <el-empty v-if="notifications.length === 0" description="暂无通知" />
      <div v-else class="notification-list">
        <div
          v-for="notification in notifications"
          :key="notification.id"
          class="notification-item"
          :class="{ unread: !notification.is_read }"
          @click="handleNotificationClick(notification)"
        >
          <div class="notification-title">{{ notification.title }}</div>
          <div class="notification-content">{{ notification.content }}</div>
          <div class="notification-time">{{ formatTime(notification.created_at) }}</div>
        </div>
      </div>
    </el-drawer>
  </el-container>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/store/user'
import { useWeddingStore } from '@/store/wedding'
import { notificationApi } from '@/api/dashboard'
import { weddingApi } from '@/api/wedding'
import { House, Bell, ArrowDown } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const weddingStore = useWeddingStore()

const userInfo = computed(() => userStore.userInfo)
const weddings = ref<any[]>([])
const currentWeddingId = ref<number | null>(null)
const showNotifications = ref(false)
const notifications = ref<any[]>([])
const unreadCount = ref(0)

const menuRoutes = computed(() => {
  return router.options.routes
    .find(r => r.path === '/')?.children
    ?.filter(r => !r.meta?.hidden && (!r.meta?.roles || userStore.hasRole(r.meta.roles as string[])))
    || []
})

const activeMenu = computed(() => route.path)

async function fetchWeddings() {
  try {
    const res = await weddingApi.getList({ page_size: 100 })
    weddings.value = res.data.list
    if (weddings.value.length > 0 && !currentWeddingId.value) {
      currentWeddingId.value = weddings.value[0].id
      weddingStore.setCurrentWedding(weddings.value[0])
    }
  } catch (error) {
    console.error('Failed to fetch weddings:', error)
  }
}

function handleWeddingChange(id: number) {
  const wedding = weddings.value.find(w => w.id === id)
  if (wedding) {
    weddingStore.setCurrentWedding(wedding)
  }
}

async function fetchNotifications() {
  try {
    const res = await notificationApi.getList()
    notifications.value = res.data.notifications
    unreadCount.value = res.data.unread_count
  } catch (error) {
    console.error('Failed to fetch notifications:', error)
  }
}

async function markAllAsRead() {
  try {
    await notificationApi.markAllAsRead()
    notifications.value.forEach(n => n.is_read = true)
    unreadCount.value = 0
    ElMessage.success('已全部标记为已读')
  } catch (error) {
    console.error('Failed to mark all as read:', error)
  }
}

async function handleNotificationClick(notification: any) {
  if (!notification.is_read) {
    try {
      await notificationApi.markAsRead(notification.id)
      notification.is_read = true
      unreadCount.value--
    } catch (error) {
      console.error('Failed to mark as read:', error)
    }
  }
}

function handleCommand(command: string) {
  if (command === 'profile') {
    router.push('/profile')
  } else if (command === 'logout') {
    userStore.logout()
    router.push('/login')
  }
}

function formatTime(time: string) {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

onMounted(() => {
  fetchWeddings()
  fetchNotifications()
})
</script>

<style scoped>
.main-layout {
  height: 100vh;
}

.sidebar {
  background-color: #001529;
  overflow-x: hidden;
}

.logo {
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
  gap: 8px;
  color: #fff;
  font-size: 16px;
  font-weight: bold;
  border-bottom: 1px solid #1f3548;
}

.sidebar .el-menu {
  border-right: none;
}

.header {
  background: #fff;
  border-bottom: 1px solid #e6e6e6;
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;
}

.user-info {
  display: flex;
  align-items: center;
  gap: 8px;
  cursor: pointer;
}

.username {
  font-size: 14px;
}

.main-content {
  background-color: #f5f7fa;
  padding: 20px;
}

.notification-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-bottom: 16px;
  border-bottom: 1px solid #e6e6e6;
  margin-bottom: 16px;
}

.notification-list {
  max-height: 400px;
  overflow-y: auto;
}

.notification-item {
  padding: 12px;
  border-radius: 4px;
  margin-bottom: 8px;
  cursor: pointer;
  transition: background-color 0.3s;
}

.notification-item:hover {
  background-color: #f5f7fa;
}

.notification-item.unread {
  background-color: #ecf5ff;
}

.notification-title {
  font-weight: 600;
  margin-bottom: 4px;
}

.notification-content {
  color: #666;
  font-size: 13px;
  margin-bottom: 4px;
}

.notification-time {
  color: #999;
  font-size: 12px;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.3s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
