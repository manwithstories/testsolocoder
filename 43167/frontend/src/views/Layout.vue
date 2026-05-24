<template>
  <el-container class="layout">
    <el-header>
      <div class="brand">手表交易与鉴定平台</div>
      <el-menu mode="horizontal" :default-active="$route.path" router class="menu">
        <el-menu-item index="/watches">手表市场</el-menu-item>
        <el-menu-item v-if="auth.user?.role === 'seller'" index="/my-watches">我的手表</el-menu-item>
        <el-menu-item v-if="auth.user?.role === 'seller'" index="/publish">发布手表</el-menu-item>
        <el-menu-item index="/auth-orders">鉴定申请</el-menu-item>
        <el-menu-item v-if="auth.user?.role === 'appraiser'" index="/auth-review">鉴定审核</el-menu-item>
        <el-menu-item index="/trades">我的交易</el-menu-item>
        <el-menu-item index="/favorites">收藏夹</el-menu-item>
        <el-menu-item index="/messages">消息中心</el-menu-item>
        <el-menu-item index="/stats">数据统计</el-menu-item>
        <el-menu-item index="/profile">个人中心</el-menu-item>
      </el-menu>
      <div class="user">
        <span>{{ auth.user?.username }} ({{ auth.user?.role }})</span>
        <el-button size="small" @click="logout">退出</el-button>
      </div>
    </el-header>
    <el-main>
      <router-view v-slot="{ Component }">
        <component :is="Component" />
      </router-view>
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { useAuthStore } from '@/stores/auth'
import { useRouter } from 'vue-router'
const auth = useAuthStore()
const router = useRouter()
function logout() {
  auth.logout()
  router.push('/login')
}
</script>

<style scoped>
.layout { height: 100vh; }
.el-header { display: flex; align-items: center; background: #fff; border-bottom: 1px solid #eee; }
.brand { font-size: 18px; font-weight: 600; margin-right: 24px; }
.menu { flex: 1; }
.user { margin-left: auto; display: flex; align-items: center; gap: 12px; }
</style>
