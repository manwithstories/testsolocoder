<template>
  <div class="page-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>消息通知</span>
          <el-button type="primary" @click="handleMarkAllRead">全部已读</el-button>
        </div>
      </template>

      <el-table :data="notifications" stripe v-loading="loading">
        <el-table-column label="标题">
          <template #default="{ row }">
            <span :class="{ 'font-bold': !row.is_read }">{{ row.title }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="type" label="类型" width="100">
          <template #default="{ row }">
            <el-tag :type="getTypeTag(row.type)" size="small">{{ getTypeText(row.type) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="content" label="内容" show-overflow-tooltip />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag v-if="!row.is_read" type="warning" size="small">未读</el-tag>
            <el-tag v-else type="success" size="small">已读</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="时间" width="180">
          <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button
              v-if="!row.is_read"
              type="primary"
              link
              @click="handleMarkRead(row)"
            >
              标记已读
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="fetchList"
          @size-change="fetchList"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getNotifications, markAsRead, markAllAsRead } from '@/api/notification'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import type { Notification } from '@/types'

const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const notifications = ref<Notification[]>([])

const getTypeTag = (type: string) => {
  const tags: Record<string, string> = {
    appointment: 'primary',
    appointment_cancel: 'warning',
    payment: 'success',
    daily_schedule: 'info',
    low_stock: 'danger',
    leave: 'info'
  }
  return tags[type] || 'info'
}

const getTypeText = (type: string) => {
  const texts: Record<string, string> = {
    appointment: '预约',
    appointment_cancel: '取消预约',
    payment: '支付',
    daily_schedule: '日程',
    low_stock: '库存预警',
    leave: '请假'
  }
  return texts[type] || type
}

const formatDateTime = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss')

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getNotifications({ page: page.value, page_size: pageSize.value })
    notifications.value = res.data.list
    total.value = res.data.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleMarkRead = async (row: Notification) => {
  try {
    await markAsRead(row.id)
    ElMessage.success('操作成功')
    fetchList()
  } catch (e) {
    console.error(e)
  }
}

const handleMarkAllRead = async () => {
  try {
    await markAllAsRead()
    ElMessage.success('操作成功')
    fetchList()
  } catch (e) {
    console.error(e)
  }
}

onMounted(fetchList)
</script>

<style scoped lang="scss">
.page-container {
  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }

  .font-bold {
    font-weight: bold;
  }
}
</style>
