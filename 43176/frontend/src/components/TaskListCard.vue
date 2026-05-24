<template>
  <div class="task-list-card">
    <div v-if="loading" class="loading">
      <el-skeleton :rows="3" animated />
    </div>

    <div v-else-if="tasks.length === 0" class="empty">
      <el-empty description="暂无任务" />
    </div>

    <div v-else class="task-list">
      <el-card
        v-for="task in tasks"
        :key="task.id"
        class="task-item"
        shadow="hover"
        @click="$emit('view', task.id)"
      >
        <div class="task-header">
          <span :class="['task-type-tag', task.type]">
            {{ getTaskTypeLabel(task.type) }}
          </span>
          <span :class="['status-tag', task.status]">
            {{ getStatusLabel(task.status) }}
          </span>
        </div>
        <h3 class="task-title">{{ task.title }}</h3>
        <p class="task-desc">{{ task.description }}</p>
        <div class="task-location">
          <el-icon><Location /></el-icon>
          <span>{{ task.start_addr }} → {{ task.end_addr }}</span>
        </div>
        <div class="task-footer">
          <div class="task-time">
            <el-icon><Clock /></el-icon>
            <span>{{ formatDeadline(task.deadline) }}</span>
          </div>
          <div class="task-reward">
            <span>报酬</span>
            <span class="reward-amount">¥{{ task.reward }}</span>
          </div>
        </div>
        <div class="task-actions">
          <el-button
            v-if="showAccept && task.status === 'pending'"
            type="primary"
            size="small"
            @click.stop="$emit('accept', task.id)"
          >
            接单
          </el-button>
          <el-button
            v-if="task.status === 'accepted' || task.status === 'in_progress'"
            type="primary"
            size="small"
            @click.stop="$emit('track', task.id)"
          >
            追踪
          </el-button>
          <el-button
            v-if="task.status === 'completed'"
            size="small"
            @click.stop="$emit('review', task.id)"
          >
            评价
          </el-button>
          <el-button size="small" @click.stop="$emit('view', task.id)">
            详情
          </el-button>
        </div>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { Location, Clock } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import type { Task, TaskType, TaskStatus } from '@/types'

defineProps<{
  tasks: Task[]
  loading?: boolean
  showAccept?: boolean
}>()

defineEmits<{
  (e: 'view', id: number): void
  (e: 'accept', id: number): void
  (e: 'track', id: number): void
  (e: 'review', id: number): void
}>()

const taskTypeLabels: Record<TaskType, string> = {
  buy: '代购',
  pickup: '代取',
  deliver: '代送',
  queue: '排队代办',
  errand: '其他代办'
}

const statusLabels: Record<TaskStatus, string> = {
  pending: '待接单',
  accepted: '已接单',
  in_progress: '进行中',
  completed: '已完成',
  cancelled: '已取消',
  timeout: '已超时'
}

const getTaskTypeLabel = (type: TaskType) => taskTypeLabels[type] || '代办'
const getStatusLabel = (status: TaskStatus) => statusLabels[status] || status

const formatDeadline = (deadline: string) => {
  const d = dayjs(deadline)
  const now = dayjs()
  const diffHours = d.diff(now, 'hour')
  if (diffHours < 0) return '已过期'
  if (diffHours < 24) return `剩余 ${diffHours} 小时`
  return d.format('MM-DD HH:mm')
}
</script>

<style lang="scss" scoped>
.task-list-card {
  .task-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
    gap: 16px;
  }

  .task-item {
    cursor: pointer;
    transition: transform 0.3s;

    &:hover {
      transform: translateY(-2px);
    }
  }

  .task-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;
  }

  .task-title {
    font-size: 16px;
    font-weight: 600;
    margin-bottom: 6px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .task-desc {
    color: #606266;
    font-size: 14px;
    margin-bottom: 10px;
    display: -webkit-box;
    -webkit-line-clamp: 2;
    -webkit-box-orient: vertical;
    overflow: hidden;
  }

  .task-location {
    display: flex;
    align-items: center;
    color: #909399;
    font-size: 13px;
    margin-bottom: 8px;

    .el-icon {
      margin-right: 4px;
    }
  }

  .task-footer {
    display: flex;
    justify-content: space-between;
    align-items: center;
    padding-top: 10px;
    border-top: 1px solid #ebeef5;
  }

  .task-time {
    display: flex;
    align-items: center;
    color: #909399;
    font-size: 13px;

    .el-icon {
      margin-right: 4px;
    }
  }

  .task-reward {
    display: flex;
    align-items: baseline;
    gap: 4px;

    span {
      color: #909399;
      font-size: 12px;
    }

    .reward-amount {
      color: #ff5722;
      font-size: 18px;
      font-weight: bold;
    }
  }

  .task-actions {
    display: flex;
    gap: 8px;
    margin-top: 12px;
  }
}

@media (max-width: 768px) {
  .task-list {
    grid-template-columns: 1fr !important;
  }
}
</style>
