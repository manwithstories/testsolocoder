<template>
  <el-container class="layout-container">
    <el-header class="header">
      <div class="header-left">
        <router-link to="/" class="logo">
          <span class="logo-icon">🚴</span>
          <span class="logo-text">跑腿服务</span>
        </router-link>
      </div>
      <div class="header-center">
        <el-menu
          :default-active="activeMenu"
          mode="horizontal"
          router
          class="nav-menu"
        >
          <el-menu-item index="/tasks">任务大厅</el-menu-item>
          <el-menu-item v-if="userStore.isPublisher || userStore.isAdmin" index="/tasks/create">
            发布任务
          </el-menu-item>
          <el-menu-item index="/my-tasks">我的任务</el-menu-item>
          <el-menu-item index="/orders">我的订单</el-menu-item>
          <el-menu-item v-if="userStore.isAdmin" index="/admin">管理后台</el-menu-item>
        </el-menu>
      </div>
      <div class="header-right">
        <template v-if="userStore.isLoggedIn">
          <el-dropdown>
            <div class="user-info">
              <el-avatar :size="36" :src="userStore.userInfo?.avatar">
                {{ userStore.userInfo?.nickname?.charAt(0) || 'U' }}
              </el-avatar>
              <span class="user-name">{{ userStore.userInfo?.nickname || '用户' }}</span>
            </div>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item @click="goTo('/profile')">个人中心</el-dropdown-item>
                <el-dropdown-item @click="goTo('/wallet')">我的钱包</el-dropdown-item>
                <el-dropdown-item @click="goTo('/reviews')">评价中心</el-dropdown-item>
                <el-dropdown-item divided @click="handleLogout">退出登录</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </template>
        <template v-else>
          <el-button type="primary" @click="goTo('/login')">登录</el-button>
          <el-button @click="goTo('/register')">注册</el-button>
        </template>
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
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)

const goTo = (path: string) => {
  router.push(path)
}

const handleLogout = () => {
  userStore.logout()
  ElMessage.success('已退出登录')
  router.push('/login')
}
</script>

<style lang="scss" scoped>
.layout-container {
  min-height: 100vh;
}

.header {
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.08);
  display: flex;
  align-items: center;
  padding: 0 24px;
  height: 64px;
}

.header-left {
  .logo {
    display: flex;
    align-items: center;
    text-decoration: none;
    color: #303133;
  }

  .logo-icon {
    font-size: 28px;
    margin-right: 8px;
  }

  .logo-text {
    font-size: 20px;
    font-weight: bold;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }
}

.header-center {
  flex: 1;
  margin: 0 40px;
}

.nav-menu {
  border-bottom: none;
}

.header-right {
  display: flex;
  align-items: center;
  gap: 16px;

  .user-info {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
  }

  .user-name {
    color: #606266;
  }
}

.main-content {
  padding: 20px;
  background: #f5f7fa;
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
