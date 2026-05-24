<template>
  <el-container class="layout-container">
    <el-header class="header">
      <div class="header-inner">
        <div class="logo" @click="$router.push('/')">
          <span class="logo-icon">👜</span>
          <span class="logo-text">奢侈品交易平台</span>
        </div>
        <el-menu
          :default-active="activeMenu"
          mode="horizontal"
          router
          class="main-menu"
        >
          <el-menu-item index="/">首页</el-menu-item>
          <el-menu-item index="/products">商品列表</el-menu-item>
          <template v-if="userStore.isLoggedIn">
            <el-menu-item v-if="isSeller" index="/seller/products">我的商品</el-menu-item>
            <el-menu-item v-if="isAuthenticator" index="/authenticator/tasks">鉴定任务</el-menu-item>
            <el-menu-item index="/my-orders">我的订单</el-menu-item>
            <el-menu-item v-if="isAdmin" index="/dashboard">数据仪表盘</el-menu-item>
            <el-menu-item v-if="isAdmin" index="/admin/users">用户管理</el-menu-item>
            <el-menu-item v-if="isAdmin" index="/admin/authenticators">鉴定师审核</el-menu-item>
          </template>
        </el-menu>
        <div class="user-info">
          <template v-if="userStore.isLoggedIn">
            <el-dropdown @command="handleCommand">
              <span class="user-name">
                <el-avatar :size="32" :src="userStore.user?.avatar">
                  {{ userStore.username?.charAt(0)?.toUpperCase() }}
                </el-avatar>
                {{ userStore.username }}
              </span>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                  <el-dropdown-item command="logout" divided>退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
          <template v-else>
            <el-button type="primary" @click="$router.push('/login')">登录</el-button>
            <el-button @click="$router.push('/register')">注册</el-button>
          </template>
        </div>
      </div>
    </el-header>
    <el-main class="main">
      <router-view v-slot="{ Component }">
        <transition name="fade" mode="out-in">
          <component :is="Component" />
        </transition>
      </router-view>
    </el-main>
    <el-footer class="footer">
      <p>© 2024 奢侈品二手交易与鉴定平台 - 专业鉴定·安全交易</p>
    </el-footer>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)

const isSeller = computed(() => userStore.userRole === 'seller' || userStore.userRole === 'admin')
const isAuthenticator = computed(() => userStore.userRole === 'authenticator' || userStore.userRole === 'admin')
const isAdmin = computed(() => userStore.userRole === 'admin')

const handleCommand = (command: string) => {
  if (command === 'profile') {
    router.push('/profile')
  } else if (command === 'logout') {
    userStore.logout()
    router.push('/login')
  }
}
</script>

<style lang="scss" scoped>
.layout-container {
  min-height: 100vh;
}

.header {
  background: #fff;
  border-bottom: 1px solid var(--border-color);
  padding: 0;
  height: 64px;
  
  .header-inner {
    max-width: 1400px;
    margin: 0 auto;
    height: 100%;
    display: flex;
    align-items: center;
    padding: 0 20px;
  }
  
  .logo {
    display: flex;
    align-items: center;
    cursor: pointer;
    margin-right: 40px;
    
    .logo-icon {
      font-size: 24px;
      margin-right: 8px;
    }
    
    .logo-text {
      font-size: 18px;
      font-weight: 600;
      color: var(--primary-color);
    }
  }
  
  .main-menu {
    border-bottom: none;
    flex: 1;
  }
  
  .user-info {
    display: flex;
    align-items: center;
    gap: 12px;
    
    .user-name {
      display: flex;
      align-items: center;
      gap: 8px;
      cursor: pointer;
    }
  }
}

.main {
  background: #f5f7fa;
  min-height: calc(100vh - 64px - 60px);
}

.footer {
  background: #fff;
  border-top: 1px solid var(--border-color);
  text-align: center;
  color: var(--text-light);
  font-size: 14px;
  height: 60px;
  display: flex;
  align-items: center;
  justify-content: center;
}

.fade-enter-active,
.fade-leave-active {
  transition: opacity 0.2s ease;
}

.fade-enter-from,
.fade-leave-to {
  opacity: 0;
}
</style>
