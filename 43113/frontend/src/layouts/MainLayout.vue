<template>
  <el-container class="main-layout">
    <el-header class="header">
      <div class="header-left">
        <router-link to="/" class="logo">
          <el-icon><ChatLineRound /></el-icon>
          <span>问答社区</span>
        </router-link>
      </div>

      <div class="header-center">
        <el-input
          v-model="searchKeyword"
          placeholder="搜索问题..."
          class="search-input"
          clearable
          @keyup.enter="handleSearch"
        >
          <template #prefix>
            <el-icon><Search /></el-icon>
          </template>
        </el-input>
      </div>

      <div class="header-right">
        <router-link to="/questions/ask" v-if="userStore.isLoggedIn">
          <el-button type="primary" :icon="Edit">提问</el-button>
        </router-link>

        <el-dropdown v-if="userStore.isLoggedIn" @command="handleCommand">
          <div class="user-info">
            <el-avatar :size="32" :src="userStore.userInfo?.avatar">
              {{ userStore.userInfo?.nickname?.charAt(0) || 'U' }}
            </el-avatar>
            <span class="username">{{ userStore.userInfo?.nickname || userStore.userInfo?.username }}</span>
            <el-badge v-if="unreadCount > 0" :value="unreadCount" class="notification-badge" />
          </div>
          <template #dropdown>
            <el-dropdown-menu>
              <el-dropdown-item command="profile">
                <el-icon><User /></el-icon>个人中心
              </el-dropdown-item>
              <el-dropdown-item command="notifications">
                <el-icon><Bell /></el-icon>
                消息通知
                <el-badge v-if="unreadCount > 0" :value="unreadCount" class="menu-badge" />
              </el-dropdown-item>
              <el-dropdown-item command="points">
                <el-icon><Coin /></el-icon>积分: {{ userStore.userInfo?.points || 0 }}
              </el-dropdown-item>
              <el-dropdown-item command="admin" v-if="userStore.isAdmin">
                <el-icon><Setting /></el-icon>管理后台
              </el-dropdown-item>
              <el-dropdown-item divided command="logout">
                <el-icon><SwitchButton /></el-icon>退出登录
              </el-dropdown-item>
            </el-dropdown-menu>
          </template>
        </el-dropdown>

        <template v-else>
          <router-link to="/login">
            <el-button>登录</el-button>
          </router-link>
          <router-link to="/register">
            <el-button type="primary">注册</el-button>
          </router-link>
        </template>
      </div>
    </el-header>

    <el-container class="main-content">
      <el-aside width="200px" class="sidebar">
        <el-menu
          :default-active="activeMenu"
          router
          class="sidebar-menu"
        >
          <el-menu-item index="/">
            <el-icon><HomeFilled /></el-icon>
            <span>首页</span>
          </el-menu-item>
          <el-menu-item index="/questions">
            <el-icon><Document /></el-icon>
            <span>问题列表</span>
          </el-menu-item>
          <el-menu-item index="/questions/ask" v-if="userStore.isLoggedIn">
            <el-icon><Edit /></el-icon>
            <span>我要提问</span>
          </el-menu-item>
          <el-menu-item index="/user/favorites" v-if="userStore.isLoggedIn">
            <el-icon><Star /></el-icon>
            <span>我的收藏</span>
          </el-menu-item>
          <el-menu-item index="/user/notifications" v-if="userStore.isLoggedIn">
            <el-icon><Bell /></el-icon>
            <span>消息通知</span>
            <el-badge v-if="unreadCount > 0" :value="unreadCount" class="menu-badge" />
          </el-menu-item>
        </el-menu>

        <div class="sidebar-footer" v-if="userStore.isLoggedIn">
          <div class="user-stats">
            <div class="stat-item">
              <span class="stat-label">等级</span>
              <span class="stat-value">Lv.{{ userStore.userInfo?.level || 1 }}</span>
            </div>
            <div class="stat-item">
              <span class="stat-label">积分</span>
              <span class="stat-value">{{ userStore.userInfo?.points || 0 }}</span>
            </div>
          </div>
        </div>
      </el-aside>

      <el-main class="content-area">
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
import { ref, computed, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'
import { notificationApi } from '@/api'

const router = useRouter()
const route = useRoute()
const userStore = useUserStore()

const searchKeyword = ref('')
const unreadCount = ref(0)

const activeMenu = computed(() => route.path)

const handleSearch = () => {
  if (searchKeyword.value.trim()) {
    router.push({ path: '/search', query: { keyword: searchKeyword.value } })
  }
}

const handleCommand = (command: string) => {
  switch (command) {
    case 'profile':
      router.push('/user/profile')
      break
    case 'notifications':
      router.push('/user/notifications')
      break
    case 'points':
      router.push('/user/points')
      break
    case 'admin':
      router.push('/admin')
      break
    case 'logout':
      userStore.logout()
      router.push('/')
      break
  }
}

const fetchUnreadCount = async () => {
  if (userStore.isLoggedIn) {
    try {
      const res = await notificationApi.getUnreadCount()
      unreadCount.value = res.data?.unreadCount || 0
    } catch (e) {
      // ignore
    }
  }
}

onMounted(() => {
  fetchUnreadCount()
  setInterval(fetchUnreadCount, 30000)
})
</script>

<style scoped lang="scss">
.main-layout {
  height: 100vh;
}

.header {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 0 20px;
  background: #fff;
  border-bottom: 1px solid #e4e7ed;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.05);
  z-index: 100;
}

.header-left {
  .logo {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 20px;
    font-weight: bold;
    color: #409eff;

    .el-icon {
      font-size: 24px;
    }
  }
}

.header-center {
  flex: 1;
  max-width: 400px;
  margin: 0 20px;

  .search-input {
    width: 100%;
  }
}

.header-right {
  display: flex;
  align-items: center;
  gap: 12px;

  .user-info {
    display: flex;
    align-items: center;
    gap: 8px;
    cursor: pointer;
    position: relative;

    .username {
      font-size: 14px;
    }

    .notification-badge {
      position: absolute;
      top: -2px;
      right: -8px;
    }
  }
}

.main-content {
  height: calc(100vh - 60px);
}

.sidebar {
  background: #fff;
  border-right: 1px solid #e4e7ed;
  display: flex;
  flex-direction: column;

  .sidebar-menu {
    border-right: none;
    flex: 1;
  }

  .sidebar-footer {
    padding: 16px;
    border-top: 1px solid #e4e7ed;

    .user-stats {
      display: flex;
      justify-content: space-around;

      .stat-item {
        text-align: center;

        .stat-label {
          display: block;
          font-size: 12px;
          color: #909399;
        }

        .stat-value {
          display: block;
          font-size: 18px;
          font-weight: bold;
          color: #409eff;
        }
      }
    }
  }
}

.content-area {
  background: #f5f7fa;
  padding: 20px;
  overflow-y: auto;
}

.menu-badge {
  margin-left: 8px;
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
