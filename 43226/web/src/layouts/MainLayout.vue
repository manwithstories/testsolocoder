<template>
  <div class="main-layout">
    <el-header class="header">
      <div class="header-content">
        <div class="logo" @click="$router.push('/')">
          <el-icon size="32" color="#409EFF"><Museum /></el-icon>
          <span class="logo-text">在线博物馆</span>
        </div>
        <el-menu mode="horizontal" :default-active="$route.path" class="nav-menu" router>
          <el-menu-item index="/">首页</el-menu-item>
          <el-menu-item index="/exhibitions">展览</el-menu-item>
          <el-menu-item index="/collections">藏品</el-menu-item>
        </el-menu>
        <div class="header-right">
          <template v-if="userStore.isLoggedIn">
            <el-button type="primary" text @click="$router.push('/dashboard')">
              管理中心
            </el-button>
            <el-dropdown @command="handleCommand">
              <span class="user-info">
                <el-avatar :size="32" :src="userStore.user?.avatar">
                  {{ userStore.user?.nickname?.charAt(0) }}
                </el-avatar>
                <span class="username">{{ userStore.user?.nickname }}</span>
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
            <el-button text @click="$router.push('/login')">登录</el-button>
            <el-button type="primary" @click="$router.push('/register')">注册</el-button>
          </template>
        </div>
      </div>
    </el-header>
    <main class="main-content">
      <router-view />
    </main>
    <el-footer class="footer">
      <p>© 2024 在线博物馆展览预约与藏品管理平台 | 基于 Go + Vue3 + Redis 技术栈</p>
    </el-footer>
  </div>
</template>

<script setup lang="ts">
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()

const handleCommand = (command: string) => {
  if (command === 'profile') {
    router.push('/dashboard/profile')
  } else if (command === 'logout') {
    userStore.logout()
    ElMessage.success('退出成功')
    router.push('/login')
  }
}
</script>

<style scoped lang="scss">
.main-layout {
  min-height: 100vh;
  display: flex;
  flex-direction: column;
}

.header {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  padding: 0;
  height: 64px;

  .header-content {
    max-width: 1400px;
    margin: 0 auto;
    padding: 0 24px;
    height: 100%;
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .logo {
    display: flex;
    align-items: center;
    gap: 10px;
    cursor: pointer;

    .logo-text {
      color: #fff;
      font-size: 20px;
      font-weight: 600;
    }
  }

  .nav-menu {
    flex: 1;
    background: transparent;
    border: none;
    margin-left: 40px;

    :deep(.el-menu-item) {
      color: rgba(255, 255, 255, 0.85);

      &:hover, &.is-active {
        color: #fff;
        background: rgba(255, 255, 255, 0.1);
      }
    }

    :deep(.el-menu-item.is-active) {
      border-bottom: 2px solid #fff;
    }
  }

  .header-right {
    display: flex;
    align-items: center;
    gap: 16px;

    .el-button {
      color: #fff;
    }

    .user-info {
      display: flex;
      align-items: center;
      gap: 8px;
      cursor: pointer;

      .username {
        color: #fff;
      }
    }
  }
}

.main-content {
  flex: 1;
  background: #f5f7fa;
}

.footer {
  background: #1f2937;
  color: #9ca3af;
  text-align: center;
  padding: 24px;
  height: auto;
}
</style>
