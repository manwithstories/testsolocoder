<template>
  <div class="technician-layout">
    <el-container>
      <el-aside width="220px" class="aside">
        <div class="user-info">
          <el-avatar :size="50" :src="userStore.userInfo?.avatar">
            {{ userStore.userInfo?.nickname?.charAt(0) || 'T' }}
          </el-avatar>
          <div class="user-name">{{ userStore.userInfo?.nickname || userStore.userInfo?.username }}</div>
          <div class="user-role">维修技师中心</div>
        </div>
        <el-menu
          :default-active="activeMenu"
          class="side-menu"
          @select="handleMenuSelect"
        >
          <el-menu-item index="services">
            <el-icon><Tools /></el-icon>
            <span>服务管理</span>
          </el-menu-item>
          <el-menu-item index="orders">
            <el-icon><List /></el-icon>
            <span>维修订单</span>
          </el-menu-item>
          <el-menu-item index="stats">
            <el-icon><DataLine /></el-icon>
            <span>数据统计</span>
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
import { computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useUserStore } from '@/store/user'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const activeMenu = computed(() => {
  const path = route.path
  if (path.includes('services/create')) return 'services'
  if (path.includes('services')) return 'services'
  if (path.includes('orders')) return 'orders'
  if (path.includes('stats')) return 'stats'
  return 'services'
})

function handleMenuSelect(index: string) {
  router.push(`/technician/${index}`)
}
</script>

<style lang="scss" scoped>
.technician-layout {
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
}

.el-main {
  margin-left: 220px;
  padding: 20px;
}
</style>
