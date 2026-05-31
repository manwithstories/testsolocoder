<template>
  <div class="user-layout">
    <el-container>
      <el-aside width="220px" class="aside">
        <div class="user-info">
          <el-avatar :size="60" :src="userStore.userInfo?.avatar">
            {{ userStore.userInfo?.nickname?.charAt(0) || 'U' }}
          </el-avatar>
          <div class="user-name">{{ userStore.userInfo?.nickname || userStore.userInfo?.username }}</div>
          <div class="user-role">{{ getRoleText(userStore.userInfo?.role) }}</div>
        </div>
        <el-menu
          :default-active="activeMenu"
          class="side-menu"
          @select="handleMenuSelect"
        >
          <el-menu-item index="profile">
            <el-icon><User /></el-icon>
            <span>个人中心</span>
          </el-menu-item>
          <el-menu-item index="orders">
            <el-icon><List /></el-icon>
            <span>我的订单</span>
          </el-menu-item>
          <el-menu-item index="repair-orders">
            <el-icon><Tools /></el-icon>
            <span>维修订单</span>
          </el-menu-item>
          <el-menu-item index="favorites">
            <el-icon><Star /></el-icon>
            <span>我的收藏</span>
          </el-menu-item>
          <el-menu-item index="wallet">
            <el-icon><Wallet /></el-icon>
            <span>我的钱包</span>
          </el-menu-item>
          <el-menu-item index="messages">
            <el-icon><Message /></el-icon>
            <span>消息中心</span>
            <el-badge v-if="unreadCount > 0" :value="unreadCount" class="badge" />
          </el-menu-item>
          <el-menu-item index="reviews">
            <el-icon><ChatDotRound /></el-icon>
            <span>评价管理</span>
          </el-menu-item>
        </el-menu>
      </el-aside>
      <el-main>
        <router-view />
      </el-main>
    </el-container>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'
import { notificationApi, messageApi } from '@/api'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const unreadCount = ref(0)

const activeMenu = computed(() => {
  const path = route.path
  if (path.includes('orders') && !path.includes('repair')) return 'orders'
  if (path.includes('repair-orders')) return 'repair-orders'
  if (path.includes('favorites')) return 'favorites'
  if (path.includes('wallet')) return 'wallet'
  if (path.includes('messages')) return 'messages'
  if (path.includes('reviews')) return 'reviews'
  return 'profile'
})

function getRoleText(role?: string): string {
  const roleMap: Record<string, string> = {
    admin: '管理员',
    seller: '卖家',
    buyer: '买家',
    technician: '维修技师'
  }
  return roleMap[role || ''] || '用户'
}

function handleMenuSelect(index: string) {
  router.push(`/user/${index}`)
}

async function fetchUnreadCount() {
  try {
    const [notifRes, msgRes] = await Promise.all([
      notificationApi.getUnreadCount(),
      messageApi.getUnreadCount()
    ])
    unreadCount.value = notifRes.data.unreadCount + msgRes.data.unreadCount
  } catch (error) {
    console.error('Failed to fetch unread count:', error)
  }
}

onMounted(() => {
  fetchUnreadCount()
})
</script>

<style lang="scss" scoped>
.user-layout {
  min-height: 100vh;
  background: #f5f7fa;
}

.aside {
  background: #fff;
  border-right: 1px solid #e4e7ed;
  height: 100vh;
  position: fixed;
  left: 0;
  top: 0;
}

.user-info {
  text-align: center;
  padding: 30px 20px;
  border-bottom: 1px solid #f0f0f0;

  .user-name {
    margin-top: 12px;
    font-weight: 500;
  }

  .user-role {
    margin-top: 4px;
    font-size: 12px;
    color: var(--text-lighter-color);
  }
}

.side-menu {
  border: none;

  .el-menu-item {
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;

    .badge {
      position: absolute;
      right: 20px;
    }
  }
}

.el-main {
  margin-left: 220px;
  padding: 20px;
}
</style>
