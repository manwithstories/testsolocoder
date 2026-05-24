<template>
  <div class="layout">
    <header class="header">
      <div class="header-content container">
        <div class="logo" @click="goHome">
          <el-icon :size="28"><Music /></el-icon>
          <span>独立音乐人平台</span>
        </div>
        
        <nav class="nav">
          <router-link to="/home" class="nav-item">首页</router-link>
          <router-link to="/ranking" class="nav-item">排行榜</router-link>
          <router-link to="/works" class="nav-item">作品</router-link>
          <router-link to="/artists" class="nav-item">音乐人</router-link>
          <router-link to="/events" class="nav-item">演出</router-link>
          <router-link to="/playlists" class="nav-item">歌单</router-link>
        </nav>
        
        <div class="user-area">
          <template v-if="userStore.isLoggedIn">
            <el-badge :value="unreadCount" :hidden="unreadCount === 0" class="notification-badge">
              <el-button text @click="goToNotifications">
                <el-icon :size="20"><Bell /></el-icon>
              </el-button>
            </el-badge>
            <el-dropdown @command="handleCommand">
              <div class="user-info">
                <el-avatar :size="32" :src="userStore.user?.avatar">
                  {{ userStore.user?.nickname?.charAt(0) || 'U' }}
                </el-avatar>
                <span class="username">{{ userStore.user?.nickname || userStore.user?.username }}</span>
                <el-icon><ArrowDown /></el-icon>
              </div>
              <template #dropdown>
                <el-dropdown-menu>
                  <el-dropdown-item command="profile">个人中心</el-dropdown-item>
                  <el-dropdown-item command="works" v-if="userStore.isArtist">我的作品</el-dropdown-item>
                  <el-dropdown-item command="revenue" v-if="userStore.isArtist">我的收益</el-dropdown-item>
                  <el-dropdown-item command="tickets">我的票</el-dropdown-item>
                  <el-dropdown-item command="playlists">我的歌单</el-dropdown-item>
                  <el-dropdown-item command="follows">关注列表</el-dropdown-item>
                  <el-dropdown-item command="admin" v-if="userStore.isAdmin">管理后台</el-dropdown-item>
                  <el-dropdown-item divided command="logout">退出登录</el-dropdown-item>
                </el-dropdown-menu>
              </template>
            </el-dropdown>
          </template>
          <template v-else>
            <el-button type="primary" @click="goToLogin">登录</el-button>
            <el-button @click="goToRegister">注册</el-button>
          </template>
        </div>
      </div>
    </header>
    
    <main class="main">
      <div class="container">
        <router-view v-slot="{ Component }">
          <transition name="fade" mode="out-in">
            <component :is="Component" />
          </transition>
        </router-view>
      </div>
    </main>
    
    <footer class="footer">
      <div class="container">
        <p>© 2024 独立音乐人平台 版权所有</p>
      </div>
    </footer>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { communityApi } from '@/api/community'

const router = useRouter()
const userStore = useUserStore()
const unreadCount = ref(0)

onMounted(() => {
  if (userStore.isLoggedIn) {
    loadUnreadCount()
  }
})

async function loadUnreadCount() {
  try {
    const res = await communityApi.getUnreadNotificationCount()
    unreadCount.value = res.count
  } catch (e) {
    console.error(e)
  }
}

function goHome() {
  router.push('/home')
}

function goToLogin() {
  router.push('/login')
}

function goToRegister() {
  router.push('/register')
}

function goToNotifications() {
  router.push('/user/notifications')
}

function handleCommand(command: string) {
  switch (command) {
    case 'profile':
      router.push('/user/profile')
      break
    case 'works':
      router.push('/user/works')
      break
    case 'revenue':
      router.push('/user/revenue')
      break
    case 'tickets':
      router.push('/user/tickets')
      break
    case 'playlists':
      router.push('/user/playlists')
      break
    case 'follows':
      router.push('/user/follows')
      break
    case 'admin':
      router.push('/admin/dashboard')
      break
    case 'logout':
      userStore.logout()
      router.push('/login')
      break
  }
}
</script>

<style scoped lang="scss">
.layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background: #fff;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  position: sticky;
  top: 0;
  z-index: 100;
  
  .header-content {
    display: flex;
    align-items: center;
    justify-content: space-between;
    height: 64px;
  }
  
  .logo {
    display: flex;
    align-items: center;
    gap: 10px;
    font-size: 18px;
    font-weight: 600;
    cursor: pointer;
    color: var(--primary-color);
  }
  
  .nav {
    display: flex;
    gap: 8px;
    
    .nav-item {
      padding: 8px 16px;
      border-radius: 4px;
      color: var(--text-light);
      transition: all 0.3s;
      
      &:hover, &.router-link-active {
        color: var(--primary-color);
        background: rgba(64, 158, 255, 0.1);
      }
    }
  }
  
  .user-area {
    display: flex;
    align-items: center;
    gap: 16px;
    
    .notification-badge {
      margin-right: 8px;
    }
    
    .user-info {
      display: flex;
      align-items: center;
      gap: 8px;
      cursor: pointer;
      
      .username {
        font-size: 14px;
      }
    }
  }
}

.main {
  flex: 1;
  padding: 24px 0;
}

.footer {
  background: #fff;
  padding: 20px 0;
  border-top: 1px solid var(--border-color);
  text-align: center;
  color: var(--text-light);
  font-size: 14px;
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
