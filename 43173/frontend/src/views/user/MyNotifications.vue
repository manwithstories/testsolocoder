<template>
  <div class="notifications">
    <div class="page-header">
      <h2>消息通知</h2>
      <el-button type="primary" text @click="markAllRead">
        全部标为已读
      </el-button>
    </div>
    
    <div class="notification-list" v-loading="loading">
      <div 
        v-for="notification in notifications" 
        :key="notification.id"
        class="notification-item"
        :class="{ unread: !notification.is_read }"
        @click="handleNotification(notification)"
      >
        <div class="icon">
          <el-icon :size="20"><Bell /></el-icon>
        </div>
        <div class="content">
          <div class="title">{{ notification.title }}</div>
          <div class="message">{{ notification.content }}</div>
          <div class="time">{{ formatTime(notification.created_at) }}</div>
        </div>
        <div class="unread-dot" v-if="!notification.is_read"></div>
      </div>
      
      <el-empty v-if="notifications.length === 0 && !loading" description="暂无消息" />
    </div>
    
    <div class="pagination" v-if="total > 0">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadNotifications"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { communityApi } from '@/api/community'
import dayjs from 'dayjs'

const loading = ref(false)
const notifications = ref<any[]>([])
const page = ref(1)
const pageSize = ref(20)
const total = ref(0)

onMounted(() => {
  loadNotifications()
})

async function loadNotifications() {
  loading.value = true
  try {
    const res = await communityApi.getNotifications({
      page: page.value,
      page_size: pageSize.value
    })
    notifications.value = res.list
    total.value = res.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function handleNotification(notification: any) {
  if (!notification.is_read) {
    try {
      await communityApi.markNotificationAsRead(notification.id)
      notification.is_read = true
    } catch (e) {
      console.error(e)
    }
  }
}

async function markAllRead() {
  try {
    await communityApi.markAllNotificationsAsRead()
    ElMessage.success('已全部标为已读')
    loadNotifications()
  } catch (e) {
    console.error(e)
  }
}

function formatTime(time: string) {
  const diff = dayjs().diff(dayjs(time), 'minute')
  if (diff < 60) {
    return `${diff}分钟前`
  } else if (diff < 1440) {
    return `${Math.floor(diff / 60)}小时前`
  } else {
    return dayjs(time).format('YYYY-MM-DD HH:mm')
  }
}
</script>

<style scoped lang="scss">
.notifications {
  .notification-list {
    .notification-item {
      display: flex;
      align-items: flex-start;
      gap: 12px;
      padding: 16px;
      border-radius: 8px;
      cursor: pointer;
      transition: background 0.3s;
      position: relative;
      
      &:hover {
        background: rgba(64, 158, 255, 0.05);
      }
      
      &.unread {
        background: rgba(64, 158, 255, 0.03);
      }
      
      .icon {
        width: 40px;
        height: 40px;
        background: rgba(64, 158, 255, 0.1);
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        color: var(--primary-color);
        flex-shrink: 0;
      }
      
      .content {
        flex: 1;
        
        .title {
          font-weight: 500;
          margin-bottom: 4px;
        }
        
        .message {
          font-size: 13px;
          color: var(--text-light);
          margin-bottom: 4px;
        }
        
        .time {
          font-size: 12px;
          color: var(--text-light);
        }
      }
      
      .unread-dot {
        width: 8px;
        height: 8px;
        background: var(--primary-color);
        border-radius: 50%;
        position: absolute;
        top: 20px;
        right: 16px;
      }
    }
  }
  
  .pagination {
    display: flex;
    justify-content: center;
    margin-top: 20px;
  }
}
</style>
