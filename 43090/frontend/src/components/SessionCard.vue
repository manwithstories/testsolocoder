<template>
  <el-card class="session-card" shadow="hover" @click="$router.push(`/sessions/${session.id}`)">
    <div class="session-header">
      <h3>{{ session.name }}</h3>
      <el-tag :type="statusType" size="small">{{ statusText }}</el-tag>
    </div>
    <p class="session-desc">{{ session.description }}</p>
    <div class="session-time">
      <div class="time-item">
        <el-icon><Clock /></el-icon>
        <span>开始: {{ formatTime(session.start_time) }}</span>
      </div>
      <div class="time-item">
        <el-icon><Clock /></el-icon>
        <span>结束: {{ formatTime(session.end_time) }}</span>
      </div>
    </div>
    <div class="session-info">
      <span>最小加价: ¥{{ session.min_increment.toFixed(2) }}</span>
      <span>延时: {{ session.extend_time }}秒</span>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { computed } from 'vue'
import dayjs from 'dayjs'
import type { AuctionSession } from '@/types'

const props = defineProps<{
  session: AuctionSession
}>()

const statusText = computed(() => {
  const map: Record<number, string> = {
    0: '未开始',
    1: '进行中',
    2: '已结束',
    3: '已取消',
  }
  return map[props.session.status] || ''
})

const statusType = computed(() => {
  const map: Record<number, string> = {
    0: 'info',
    1: 'success',
    2: 'info',
    3: 'danger',
  }
  return map[props.session.status] || 'info'
})

const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}
</script>

<style scoped>
.session-card {
  cursor: pointer;
  transition: transform 0.3s;
  margin-bottom: 20px;
}

.session-card:hover {
  transform: translateY(-3px);
}

.session-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 10px;
}

.session-header h3 {
  margin: 0;
  font-size: 18px;
  color: #303133;
}

.session-desc {
  color: #606266;
  font-size: 14px;
  margin-bottom: 15px;
  overflow: hidden;
  text-overflow: ellipsis;
  display: -webkit-box;
  -webkit-line-clamp: 2;
  -webkit-box-orient: vertical;
}

.session-time {
  margin-bottom: 10px;
}

.time-item {
  display: flex;
  align-items: center;
  gap: 6px;
  color: #909399;
  font-size: 13px;
  margin-bottom: 5px;
}

.session-info {
  display: flex;
  justify-content: space-between;
  padding-top: 10px;
  border-top: 1px solid #f0f0f0;
  font-size: 13px;
  color: #606266;
}
</style>
