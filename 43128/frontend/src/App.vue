<template>
  <el-container class="layout-container">
    <el-header v-if="userStore.isLogin">
      <div class="logo">赛事平台</div>
      <el-menu mode="horizontal" :router="true" :default-active="$route.path" class="menu">
        <el-menu-item index="/">首页</el-menu-item>
        <el-menu-item index="/events">赛事</el-menu-item>
        <el-menu-item index="/my/registrations">我的报名</el-menu-item>
        <el-menu-item index="/my/scores">我的成绩</el-menu-item>
        <el-menu-item index="/my/certificates">我的证书</el-menu-item>
        <el-menu-item index="/messages">消息<el-badge v-if="noticeStore.unreadCount>0" :value="noticeStore.unreadCount" class="badge" /></el-menu-item>
        <el-menu-item v-if="userStore.isAdmin" index="/admin/events">后台</el-menu-item>
      </el-menu>
      <div class="user-area">
        <span>{{ userStore.userInfo?.username || '未登录' }}</span>
        <el-button link @click="logout">退出</el-button>
      </div>
    </el-header>
    <el-main>
      <router-view />
    </el-main>
  </el-container>
</template>

<script setup lang="ts">
import { onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { useNoticeStore } from '@/stores/notice'
import { msgApi } from '@/api'

const userStore = useUserStore()
const noticeStore = useNoticeStore()

onMounted(async () => {
  if (userStore.isLogin) {
    await userStore.fetchProfile()
    try {
      const res = await msgApi.unreadCount()
      noticeStore.unreadCount = (res.data as any)?.count || 0
    } catch (_) { /* ignore */ }
  }
})

function logout() {
  userStore.logout()
  location.href = '/login'
}
</script>

<style scoped>
.layout-container { min-height: 100vh; }
.el-header { display: flex; align-items: center; background: #fff; border-bottom: 1px solid #ebeef5; }
.logo { font-weight: 600; font-size: 18px; margin-right: 24px; color: #409eff; }
.menu { flex: 1; }
.user-area { display: flex; align-items: center; gap: 12px; }
.badge { margin-left: 4px; }
</style>
