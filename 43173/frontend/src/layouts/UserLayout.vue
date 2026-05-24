<template>
  <div class="user-layout">
    <div class="sidebar">
      <div class="user-card">
        <el-avatar :size="64" :src="userStore.user?.avatar">
          {{ userStore.user?.nickname?.charAt(0) || 'U' }}
        </el-avatar>
        <div class="user-name">{{ userStore.user?.nickname || userStore.user?.username }}</div>
        <div class="user-role">{{ roleText }}</div>
      </div>
      
      <el-menu
        :default-active="activeMenu"
        router
        background-color="transparent"
        text-color="#606266"
        active-text-color="#409eff"
      >
        <el-menu-item index="/user/profile">
          <el-icon><User /></el-icon>
          <span>个人中心</span>
        </el-menu-item>
        
        <template v-if="userStore.isArtist">
          <el-menu-item index="/user/works">
            <el-icon><Headset /></el-icon>
            <span>我的作品</span>
          </el-menu-item>
          <el-menu-item index="/user/upload">
            <el-icon><Upload /></el-icon>
            <span>上传作品</span>
          </el-menu-item>
          <el-menu-item index="/user/albums">
            <el-icon><Collection /></el-icon>
            <span>我的专辑</span>
          </el-menu-item>
          <el-menu-item index="/user/events">
            <el-icon><Calendar /></el-icon>
            <span>我的演出</span>
          </el-menu-item>
          <el-menu-item index="/user/revenue">
            <el-icon><Money /></el-icon>
            <span>我的收益</span>
          </el-menu-item>
          <el-menu-item index="/user/withdraw">
            <el-icon><Wallet /></el-icon>
            <span>申请提现</span>
          </el-menu-item>
          <el-menu-item index="/user/stats">
            <el-icon><DataAnalysis /></el-icon>
            <span>数据统计</span>
          </el-menu-item>
        </template>
        
        <el-menu-item index="/user/tickets">
          <el-icon><Ticket /></el-icon>
          <span>我的票</span>
        </el-menu-item>
        <el-menu-item index="/user/playlists">
          <el-icon><List /></el-icon>
          <span>我的歌单</span>
        </el-menu-item>
        <el-menu-item index="/user/follows">
          <el-icon><UserFilled /></el-icon>
          <span>关注列表</span>
        </el-menu-item>
        <el-menu-item index="/user/notifications">
          <el-icon><Bell /></el-icon>
          <span>消息通知</span>
        </el-menu-item>
      </el-menu>
    </div>
    
    <div class="content">
      <router-view />
    </div>
  </div>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import { useRoute } from 'vue-router'
import { useUserStore } from '@/stores/user'

const route = useRoute()
const userStore = useUserStore()

const activeMenu = computed(() => route.path)

const roleText = computed(() => {
  switch (userStore.user?.role) {
    case 'artist': return '独立音乐人'
    case 'label': return '厂牌'
    case 'admin': return '管理员'
    default: return '乐迷'
  }
})
</script>

<style scoped lang="scss">
.user-layout {
  display: flex;
  min-height: calc(100vh - 64px - 60px);
  gap: 24px;
}

.sidebar {
  width: 240px;
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
  flex-shrink: 0;
  height: fit-content;
  
  .user-card {
    text-align: center;
    padding-bottom: 20px;
    border-bottom: 1px solid var(--border-color);
    margin-bottom: 16px;
    
    .user-name {
      margin-top: 12px;
      font-size: 16px;
      font-weight: 500;
    }
    
    .user-role {
      margin-top: 4px;
      font-size: 13px;
      color: var(--text-light);
    }
  }
  
  :deep(.el-menu) {
    border: none;
  }
  
  :deep(.el-menu-item) {
    border-radius: 4px;
    margin-bottom: 4px;
    
    &:hover {
      background: rgba(64, 158, 255, 0.1);
    }
  }
}

.content {
  flex: 1;
  background: #fff;
  border-radius: 8px;
  padding: 24px;
  box-shadow: 0 2px 12px rgba(0, 0, 0, 0.05);
}
</style>
