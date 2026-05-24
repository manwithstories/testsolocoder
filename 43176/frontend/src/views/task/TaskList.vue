<template>
  <div class="task-list-container">
    <div class="filter-bar">
      <el-tabs v-model="activeTab" class="filter-tabs">
        <el-tab-pane label="全部任务" name="all" />
        <el-tab-pane label="代购" name="buy" />
        <el-tab-pane label="代取" name="pickup" />
        <el-tab-pane label="代送" name="deliver" />
        <el-tab-pane label="排队代办" name="queue" />
      </el-tabs>
      <div class="filter-actions">
        <el-select v-model="sortBy" placeholder="排序方式" style="width: 140px">
          <el-option label="最新发布" value="newest" />
          <el-option label="报酬最高" value="reward" />
          <el-option label="距离最近" value="distance" />
        </el-select>
      </div>
    </div>

    <div v-if="loading" class="loading-container">
      <el-skeleton :rows="5" animated />
    </div>

    <div v-else class="task-grid">
      <el-card
        v-for="task in tasks"
        :key="task.id"
        class="task-card"
        shadow="hover"
        @click="goToDetail(task.id)"
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
        <div v-if="task.distance !== undefined" class="task-distance">
          距离约 {{ task.distance.toFixed(1) }} 公里
        </div>
        <div class="task-footer">
          <div class="task-time">
            <el-icon><Clock /></el-icon>
            <span>{{ formatDeadline(task.deadline) }}</span>
          </div>
          <div class="task-reward">
            <span class="reward-label">报酬</span>
            <span class="reward-amount">¥{{ task.reward }}</span>
          </div>
        </div>
        <div v-if="task.images && task.images.length > 0" class="task-images">
          <el-image
            v-for="(img, index) in task.images.slice(0, 3)"
            :key="index"
            :src="img.image_url"
            :preview-src-list="task.images.map(i => i.image_url)"
            fit="cover"
            class="task-image"
          />
          <div v-if="task.images.length > 3" class="image-more">
            +{{ task.images.length - 3 }}
          </div>
        </div>
      </el-card>
    </div>

    <div v-if="!loading && tasks.length === 0" class="empty-state">
      <el-empty description="暂无任务" />
    </div>

    <div class="pagination-container">
      <el-pagination
        v-model:current-page="currentPage"
        v-model:page-size="pageSize"
        :total="total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="fetchTasks"
        @current-change="fetchTasks"
      />
    </div>

    <el-button
      v-if="userStore.isPublisher || userStore.isAdmin"
      class="fab-btn"
      type="primary"
      circle
      @click="goToCreate"
    >
      <el-icon><Plus /></el-icon>
    </el-button>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { Location, Clock, Plus } from '@element-plus/icons-vue'
import { taskApi } from '@/api'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'
import type { Task, TaskType, TaskStatus } from '@/types'

const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const tasks = ref<Task[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const activeTab = ref('all')
const sortBy = ref('newest')

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

const fetchTasks = async () => {
  loading.value = true
  try {
    const params: any = {
      status: 'pending',
      page: currentPage.value,
      page_size: pageSize.value
    }
    if (activeTab.value !== 'all') {
      params.type = activeTab.value
    }
    const res = await taskApi.list(params)
    if (res.code === 200) {
      tasks.value = res.data.items
      total.value = res.data.total
    }
  } catch (error) {
    console.error('Failed to fetch tasks:', error)
  } finally {
    loading.value = false
  }
}

const goToDetail = (id: number) => {
  router.push(`/tasks/${id}`)
}

const goToCreate = () => {
  router.push('/tasks/create')
}

onMounted(() => {
  fetchTasks()
})
</script>

<style lang="scss" scoped>
.task-list-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.filter-bar {
  display: flex;
  justify-content: space-between;
  align-items: center;
  background: #fff;
  padding: 16px 20px;
  border-radius: 8px;
  margin-bottom: 20px;
  box-shadow: 0 2px 8px rgba(0, 0, 0, 0.06);

  .filter-tabs {
    flex: 1;
  }
}

.loading-container {
  padding: 20px;
}

.task-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(320px, 1fr));
  gap: 20px;
}

.task-card {
  cursor: pointer;
  transition: transform 0.3s, box-shadow 0.3s;

  &:hover {
    transform: translateY(-4px);
  }
}

.task-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 12px;
}

.task-title {
  font-size: 16px;
  font-weight: 600;
  color: #303133;
  margin-bottom: 8px;
  overflow: hidden;
  text-overflow: ellipsis;
  white-space: nowrap;
}

.task-desc {
  color: #606266;
  font-size: 14px;
  line-height: 1.5;
  margin-bottom: 12px;
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

.task-distance {
  color: #67c23a;
  font-size: 13px;
  margin-bottom: 8px;
}

.task-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
  padding-top: 12px;
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

  .reward-label {
    color: #909399;
    font-size: 12px;
    margin-right: 4px;
  }

  .reward-amount {
    color: #ff5722;
    font-size: 20px;
    font-weight: bold;
  }
}

.task-images {
  display: flex;
  gap: 8px;
  margin-top: 12px;
  position: relative;

  .task-image {
    width: 60px;
    height: 60px;
    border-radius: 4px;
  }

  .image-more {
    position: absolute;
    right: 0;
    bottom: 0;
    width: 60px;
    height: 60px;
    background: rgba(0, 0, 0, 0.5);
    color: #fff;
    display: flex;
    align-items: center;
    justify-content: center;
    border-radius: 4px;
    font-size: 14px;
  }
}

.empty-state {
  padding: 60px 20px;
  text-align: center;
}

.pagination-container {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}

.fab-btn {
  position: fixed;
  right: 30px;
  bottom: 30px;
  width: 56px;
  height: 56px;
  font-size: 24px;
  box-shadow: 0 4px 12px rgba(102, 126, 234, 0.4);
}

@media (max-width: 768px) {
  .task-list-container {
    padding: 10px;
  }

  .task-grid {
    grid-template-columns: 1fr;
  }

  .filter-bar {
    flex-direction: column;
    gap: 12px;
  }
}
</style>
