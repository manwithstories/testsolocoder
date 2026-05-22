<template>
  <el-container class="user-layout">
    <el-aside width="220px" class="aside">
      <el-menu
        :default-active="activeMenu"
        router
        background-color="#f5f7fa"
        text-color="#606266"
        active-text-color="#409eff"
      >
        <el-menu-item index="/user/profile">
          <el-icon><User /></el-icon>
          <span>个人中心</span>
        </el-menu-item>
        <el-menu-item index="/user/items">
          <el-icon><Goods /></el-icon>
          <span>我的拍卖品</span>
        </el-menu-item>
        <el-menu-item index="/user/bids">
          <el-icon><Cpu /></el-icon>
          <span>我的出价</span>
        </el-menu-item>
        <el-menu-item index="/user/auto-bids">
          <el-icon><SwitchButton /></el-icon>
          <span>自动出价</span>
        </el-menu-item>
        <el-menu-item index="/user/orders">
          <el-icon><List /></el-icon>
          <span>我的订单</span>
        </el-menu-item>
        <el-menu-item index="/user/notifications">
          <el-icon><Bell /></el-icon>
          <span>消息中心</span>
        </el-menu-item>
      </el-menu>
    </el-aside>
    <el-main class="main">
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'

const route = useRoute()
const activeMenu = computed(() => route.path)
</script>

<style scoped>
.user-layout {
  background: white;
  border-radius: 8px;
  min-height: 600px;
}

.aside {
  border-right: 1px solid #f0f0f0;
}

.aside .el-menu {
  border-right: none;
}

.main {
  padding: 20px;
}
</style>
