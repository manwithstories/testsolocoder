<template>
  <div class="notification-list">
    <div class="page-header flex-between">
      <h2 class="page-title">消息中心</h2>
      <el-button @click="markAllAsRead" :disabled="unreadCount === 0">
        全部标为已读
      </el-button>
    </div>

    <div class="filter-bar">
      <el-radio-group v-model="filterType" @change="fetchData">
        <el-radio-button label="">全部</el-radio-button>
        <el-radio-button label="unread">未读</el-radio-button>
        <el-radio-button label="read">已读</el-radio-button>
      </el-radio-group>
    </div>

    <div class="table-wrapper">
      <div
        v-for="notification in notifications"
        :key="notification.id"
        class="notification-item"
        :class="{ unread: !notification.isRead }"
        @click="handleNotificationClick(notification)"
      >
        <div class="notification-header">
          <el-icon class="notification-icon" :size="20">
            <Bell />
          </el-icon>
          <span class="notification-title">{{ notification.title }}</span>
          <span class="notification-time">{{ formatTime(notification.createdAt) }}</span>
          <el-tag v-if="!notification.isRead" type="danger" size="small">未读</el-tag>
        </div>
        <div class="notification-content">{{ notification.content }}</div>
        <div class="notification-footer" v-if="notification.application">
          <el-button type="primary" link @click.stop="goToApplication(notification.applicationId!)">
            查看申请详情
          </el-button>
        </div>
      </div>

      <div v-if="notifications.length === 0 && !loading" class="empty-state">
        <el-icon :size="48"><Bell /></el-icon>
        <p>暂无通知</p>
      </div>

      <div v-if="loading" class="loading-wrapper">
        <el-icon class="is-loading" :size="32"><Loading /></el-icon>
      </div>
    </div>

    <div class="pagination-wrapper mt-16">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50, 100]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="fetchData"
        @current-change="fetchData"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Bell, Loading } from '@element-plus/icons-vue'
import { notificationApi } from '@/api/notification'
import { Notification } from '@/types'
import dayjs from 'dayjs'

const router = useRouter()
const loading = ref(false)
const notifications = ref<Notification[]>([])
const filterType = ref('')
const unreadCount = ref(0)

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const fetchData = async () => {
  loading.value = true
  try {
    const isRead = filterType.value === 'read' ? true : filterType.value === 'unread' ? false : undefined
    const res = await notificationApi.getList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      isRead
    })
    notifications.value = res?.list || []
    pagination.total = res?.total || 0
  } catch (error) {
    console.error('获取通知列表失败:', error)
  } finally {
    loading.value = false
  }
}

const fetchUnreadCount = async () => {
  try {
    const res = await notificationApi.getUnreadCount()
    unreadCount.value = res?.count || 0
  } catch (error) {
    console.error('获取未读数量失败:', error)
  }
}

const handleNotificationClick = async (notification: Notification) => {
  if (!notification.isRead) {
    try {
      await notificationApi.markAsRead(notification.id)
      notification.isRead = true
      unreadCount.value--
    } catch (error) {
      console.error('标记已读失败:', error)
    }
  }
}

const markAllAsRead = async () => {
  try {
    await notificationApi.markAllAsRead()
    ElMessage.success('已全部标记为已读')
    notifications.value.forEach(n => n.isRead = true)
    unreadCount.value = 0
  } catch (error) {
    ElMessage.error('操作失败')
  }
}

const goToApplication = (applicationId: number) => {
  router.push(`/applications/${applicationId}`)
}

const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

onMounted(() => {
  fetchData()
  fetchUnreadCount()
})
</script>

<style scoped>
.notification-item {
  padding: 16px;
  border-radius: 8px;
  margin-bottom: 12px;
  background: #fff;
  cursor: pointer;
  transition: all 0.2s;
  border: 1px solid #ebeef5;
}

.notification-item:hover {
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.1);
}

.notification-item.unread {
  background: #ecf5ff;
  border-color: #d9ecff;
}

.notification-header {
  display: flex;
  align-items: center;
  gap: 8px;
  margin-bottom: 8px;
}

.notification-icon {
  color: #409eff;
}

.notification-title {
  font-size: 15px;
  font-weight: 600;
  color: #303133;
  flex: 1;
}

.notification-time {
  font-size: 12px;
  color: #909399;
}

.notification-content {
  font-size: 14px;
  color: #606266;
  line-height: 1.6;
  padding-left: 28px;
}

.notification-footer {
  padding-left: 28px;
  margin-top: 8px;
}

.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
}
</style>
