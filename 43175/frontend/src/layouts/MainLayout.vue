<template>
  <el-container class="main-layout">
    <el-aside width="220px" class="aside">
      <div class="logo">
        <span class="logo-icon">⚡</span>
        <span class="logo-text">智能能源</span>
      </div>
      <div class="family-selector" v-if="familyStore.families.length > 0">
        <el-select
          v-model="familyStore.currentFamilyId"
          placeholder="选择家庭"
          style="width: 180px; margin: 12px 20px;"
          @change="handleFamilyChange"
        >
          <el-option
            v-for="f in familyStore.families"
            :key="f.id"
            :label="f.name"
            :value="f.id"
          />
        </el-select>
      </div>
      <el-menu
        :default-active="activeMenu"
        router
        background-color="#001529"
        text-color="#ffffffa6"
        active-text-color="#ffffff"
        class="sidebar-menu"
      >
        <el-menu-item index="/dashboard">
          <el-icon><Odometer /></el-icon>
          <span>首页概览</span>
        </el-menu-item>
        <el-menu-item index="/devices">
          <el-icon><Monitor /></el-icon>
          <span>设备管理</span>
        </el-menu-item>
        <el-menu-item index="/groups">
          <el-icon><Grid /></el-icon>
          <span>设备分组</span>
        </el-menu-item>
        <el-menu-item index="/energy">
          <el-icon><TrendCharts /></el-icon>
          <span>能耗监控</span>
        </el-menu-item>
        <el-menu-item index="/scenes">
          <el-icon><MagicStick /></el-icon>
          <span>场景联动</span>
        </el-menu-item>
        <el-menu-item index="/schedules">
          <el-icon><Clock /></el-icon>
          <span>定时任务</span>
        </el-menu-item>
        <el-menu-item index="/families">
          <el-icon><House /></el-icon>
          <span>家庭管理</span>
        </el-menu-item>
        <el-menu-item index="/notifications">
          <el-icon><Bell /></el-icon>
          <span>通知消息</span>
          <el-badge :value="unreadCount" :hidden="unreadCount === 0" class="notification-badge" />
        </el-menu-item>
      </el-menu>
    </el-aside>

    <el-container>
      <el-header class="header">
        <div class="header-left">
          <span class="page-title">{{ pageTitle }}</span>
        </div>
        <div class="header-right">
          <el-dropdown @command="handleCommand">
            <div class="user-info">
              <el-avatar :size="32" :src="userStore.user?.avatar">
                {{ userStore.user?.username?.charAt(0) || 'U' }}
              </el-avatar>
              <span class="username">{{ userStore.user?.username || '用户' }}</span>
              <el-icon><ArrowDown /></el-icon>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="profile">
                  <el-icon><User /></el-icon>
                  个人中心
                </el-dropdown-item>
                <el-dropdown-item divided command="logout">
                  <el-icon><SwitchButton /></el-icon>
                  退出登录
                </el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
      </el-header>

      <el-main class="main-content">
        <router-view />
      </el-main>
    </el-container>
  </el-container>
</template>

<script setup lang="ts">
import { computed, ref, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { useFamilyStore } from '@/stores/family'
import { listNotifications } from '@/api/notification'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()
const familyStore = useFamilyStore()
const unreadCount = ref(0)

const activeMenu = computed(() => route.path)
const pageTitle = computed(() => (route.meta.title as string) || '智能能源管理')

onMounted(async () => {
  if (userStore.token && !userStore.user) {
    await userStore.fetchProfile()
  }
  await familyStore.loadFamilies()
  loadUnreadCount()
})

async function loadUnreadCount() {
  try {
    const res = await listNotifications({ isRead: 'false' })
    unreadCount.value = res.unreadCount || 0
  } catch (e) {
    console.error(e)
  }
}

function handleFamilyChange() {
  localStorage.setItem('currentFamilyId', String(familyStore.currentFamilyId))
}

function handleCommand(command: string) {
  if (command === 'profile') {
    router.push('/profile')
  } else if (command === 'logout') {
    userStore.logout()
    router.push('/login')
  }
}
</script>

<style lang="scss" scoped>
.main-layout {
  height: 100vh;
}

.aside {
  background-color: #001529;
  display: flex;
  flex-direction: column;

  .logo {
    height: 60px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: white;
    font-size: 18px;
    font-weight: bold;
    border-bottom: 1px solid #1f3a57;

    .logo-icon {
      font-size: 24px;
      margin-right: 8px;
    }

    .logo-text {
      background: linear-gradient(135deg, #667eea, #764ba2);
      -webkit-background-clip: text;
      -webkit-text-fill-color: transparent;
    }
  }
}

.sidebar-menu {
  border-right: none;
  flex: 1;

  :deep(.el-menu-item) {
    height: 50px;
    line-height: 50px;
    margin: 4px 8px;
    border-radius: 4px;

    &:hover {
      background-color: #1f3a57 !important;
    }

    &.is-active {
      background-color: #409eff !important;
    }
  }
}

.notification-badge {
  margin-left: 8px;
}

.header {
  background: white;
  box-shadow: 0 1px 4px rgba(0, 21, 41, 0.08);
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 24px;

  .header-left {
    .page-title {
      font-size: 18px;
      font-weight: 600;
      color: #303133;
    }
  }

  .header-right {
    .user-info {
      display: flex;
      align-items: center;
      cursor: pointer;
      padding: 8px 12px;
      border-radius: 4px;

      &:hover {
        background: #f5f7fa;
      }

      .username {
        margin: 0 8px;
        color: #606266;
      }
    }
  }
}

.main-content {
  background-color: #f5f7fa;
  padding: 20px;
  overflow-y: auto;
}
</style>
