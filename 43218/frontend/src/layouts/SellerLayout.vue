<template>
  <div class="seller-layout">
    <el-container>
      <el-aside width="220px" class="aside">
        <div class="user-info">
          <el-avatar :size="50" :src="userStore.userInfo?.avatar">
            {{ userStore.userInfo?.nickname?.charAt(0) || 'S' }}
          </el-avatar>
          <div class="user-name">{{ userStore.userInfo?.nickname || userStore.userInfo?.username }}</div>
          <div class="user-role">卖家中心</div>
        </div>
        <el-menu
          :default-active="activeMenu"
          class="side-menu"
          @select="handleMenuSelect"
        >
          <el-menu-item index="products">
            <el-icon><Goods /></el-icon>
            <span>商品管理</span>
          </el-menu-item>
          <el-menu-item index="orders">
            <el-icon><List /></el-icon>
            <span>订单管理</span>
          </el-menu-item>
          <el-menu-item index="negotiations">
            <el-icon><ChatDotRound /></el-icon>
            <span>议价管理</span>
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
  if (path.includes('products/create')) return 'products'
  if (path.includes('products')) return 'products'
  if (path.includes('orders')) return 'orders'
  if (path.includes('negotiations')) return 'negotiations'
  if (path.includes('stats')) return 'stats'
  return 'products'
})

function handleMenuSelect(index: string) {
  router.push(`/seller/${index}`)
}
</script>

<style lang="scss" scoped>
.seller-layout {
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
