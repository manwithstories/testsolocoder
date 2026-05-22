<template>
  <div class="notifications">
    <div class="header">
      <h3>消息中心</h3>
      <el-button @click="markAllAsRead" :disabled="unreadCount === 0">全部标为已读</el-button>
    </div>
    <el-table :data="notifications" v-loading="loading" @selection-change="handleSelectionChange">
      <el-table-column type="selection" width="55" />
      <el-table-column prop="title" label="标题" min-width="200" />
      <el-table-column prop="content" label="内容" min-width="300" />
      <el-table-column prop="type" label="类型" width="120">
        <template #default="{ row }">
          <el-tag size="small" :type="getTypeTag(row.type)">{{ getTypeText(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="is_read" label="状态" width="80">
        <template #default="{ row }">
          <span v-if="row.is_read === 0" class="unread-dot" />
          <span v-else class="read-text">已读</span>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="时间" width="160">
        <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
      </el-table-column>
    </el-table>
    <div class="table-footer">
      <div v-if="selectedIds.length > 0" class="selected-info">
        已选择 {{ selectedIds.length }} 条
        <el-button size="small" type="primary" @click="markSelectedAsRead">标为已读</el-button>
      </div>
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="fetchNotifications"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import type { Notification } from '@/types'
import { notificationApi } from '@/api'

const notifications = ref<Notification[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(20)
const selectedIds = ref<number[]>([])

const unreadCount = computed(() => {
  return notifications.value.filter(n => n.is_read === 0).length
})

const fetchNotifications = async () => {
  loading.value = true
  try {
    const res = await notificationApi.getMyNotifications({ page: page.value, page_size: pageSize.value })
    notifications.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

const handleSelectionChange = (selection: Notification[]) => {
  selectedIds.value = selection.map(n => n.id)
}

const markSelectedAsRead = async () => {
  if (selectedIds.value.length === 0) return
  try {
    await notificationApi.markAsRead({ notification_ids: selectedIds.value })
    ElMessage.success('已标为已读')
    fetchNotifications()
  } catch (e) {}
}

const markAllAsRead = async () => {
  try {
    await notificationApi.markAllAsRead()
    ElMessage.success('已全部标为已读')
    fetchNotifications()
  } catch (e) {}
}

const getTypeTag = (type: string) => {
  const map: Record<string, string> = {
    bid_outbid: 'warning',
    auction_ending: 'info',
    bid_success: 'success',
    payment: 'primary',
    system: 'info',
  }
  return map[type] || 'info'
}

const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    bid_outbid: '出价被超越',
    auction_ending: '拍卖即将结束',
    bid_success: '竞拍成功',
    bid_failed: '竞拍失败',
    payment: '支付通知',
    system: '系统消息',
  }
  return map[type] || type
}

const formatTime = (time: string) => dayjs(time).format('YYYY-MM-DD HH:mm')

onMounted(() => {
  fetchNotifications()
})
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.header h3 {
  margin: 0;
}

.read-text {
  color: #909399;
  font-size: 13px;
}

.unread-dot {
  display: inline-block;
  width: 8px;
  height: 8px;
  background: #f56c6c;
  border-radius: 50%;
}

.table-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-top: 20px;
}

.selected-info {
  display: flex;
  align-items: center;
  gap: 10px;
}
</style>
