<template>
  <div class="notification-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>消息通知</span>
          <el-button type="primary" size="small" @click="markAllAsRead">
            全部已读
          </el-button>
        </div>
      </template>

      <el-table :data="notifications" style="width: 100%">
        <el-table-column label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getNotificationTypeTag(row.type)" size="small">
              {{ getNotificationTypeName(row.type) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="title" label="标题" />
        <el-table-column prop="content" label="内容" show-overflow-tooltip />
        <el-table-column label="状态" width="80">
          <template #default="{ row }">
            <el-tag v-if="row.isRead" type="info" size="small">已读</el-tag>
            <el-tag v-else type="warning" size="small">未读</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button
              v-if="!row.isRead"
              type="primary"
              size="small"
              link
              @click="markAsRead(row.id)"
            >
              标记已读
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="notifications.length === 0" description="暂无通知" />

      <el-pagination
        class="pagination"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="prev, pager, next"
        @current-change="fetchNotifications"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { notificationApi } from '@/api'
import type { Notification } from '@/types'
import dayjs from 'dayjs'

const notifications = ref<Notification[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)

const fetchNotifications = async () => {
  try {
    const res = await notificationApi.getNotifications({
      page: page.value,
      pageSize: pageSize.value
    })
    notifications.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    console.error(e)
  }
}

const markAsRead = async (id: number) => {
  try {
    await notificationApi.markAsRead(id)
    fetchNotifications()
  } catch (e) {
    console.error(e)
  }
}

const markAllAsRead = async () => {
  try {
    await notificationApi.markAllAsRead()
    fetchNotifications()
  } catch (e) {
    console.error(e)
  }
}

const getNotificationTypeTag = (type: string) => {
  const map: Record<string, string> = {
    answer: 'primary',
    comment: 'success',
    like: 'warning',
    follow: 'info',
    system: 'danger'
  }
  return map[type] || ''
}

const getNotificationTypeName = (type: string) => {
  const map: Record<string, string> = {
    answer: '回答',
    comment: '评论',
    like: '点赞',
    follow: '关注',
    system: '系统'
  }
  return map[type] || type
}

const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(() => {
  fetchNotifications()
})
</script>

<style scoped lang="scss">
.notification-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .pagination {
    justify-content: center;
    margin-top: 20px;
  }
}
</style>
