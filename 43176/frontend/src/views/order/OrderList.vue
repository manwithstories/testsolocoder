<template>
  <div class="order-list-container">
    <el-tabs v-model="activeTab" class="tabs">
      <el-tab-pane label="全部订单" name="all">
        <div class="order-grid">
          <el-card
            v-for="order in orders"
            :key="order.id"
            class="order-card"
            shadow="hover"
            @click="goToDetail(order.id)"
          >
            <div class="order-header">
              <span :class="['status-tag', order.status]">
                {{ getStatusLabel(order.status) }}
              </span>
              <span class="order-time">{{ formatTime(order.created_at) }}</span>
            </div>
            <h3 class="order-title">{{ order.task?.title }}</h3>
            <div class="order-info">
              <span class="task-type">{{ getTaskTypeLabel(order.task?.type) }}</span>
              <span class="reward">¥{{ order.reward }}</span>
            </div>
            <div class="order-location">
              <el-icon><Location /></el-icon>
              <span>{{ order.task?.start_addr }} → {{ order.task?.end_addr }}</span>
            </div>
          </el-card>
        </div>
      </el-tab-pane>
      <el-tab-pane label="进行中" name="in_progress">
        <div class="order-grid">
          <el-card
            v-for="order in inProgressOrders"
            :key="order.id"
            class="order-card"
            shadow="hover"
            @click="goToDetail(order.id)"
          >
            <div class="order-header">
              <span :class="['status-tag', order.status]">
                {{ getStatusLabel(order.status) }}
              </span>
            </div>
            <h3 class="order-title">{{ order.task?.title }}</h3>
          </el-card>
        </div>
      </el-tab-pane>
      <el-tab-pane label="已完成" name="completed">
        <div class="order-grid">
          <el-card
            v-for="order in completedOrders"
            :key="order.id"
            class="order-card"
            shadow="hover"
            @click="goToDetail(order.id)"
          >
            <div class="order-header">
              <span class="status-tag completed">已完成</span>
            </div>
            <h3 class="order-title">{{ order.task?.title }}</h3>
          </el-card>
        </div>
      </el-tab-pane>
    </el-tabs>

    <div v-if="orders.length === 0 && !loading" class="empty-state">
      <el-empty description="暂无订单" />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { useRouter } from 'vue-router'
import { Location } from '@element-plus/icons-vue'
import { orderApi } from '@/api'
import dayjs from 'dayjs'
import type { Order, OrderStatus, TaskType } from '@/types'

const router = useRouter()
const loading = ref(false)
const activeTab = ref('all')
const orders = ref<Order[]>([])

const statusLabels: Record<OrderStatus, string> = {
  pending: '待处理',
  accepted: '已接单',
  in_progress: '进行中',
  delivered: '已送达',
  completed: '已完成',
  cancelled: '已取消'
}

const taskTypeLabels: Record<TaskType, string> = {
  buy: '代购',
  pickup: '代取',
  deliver: '代送',
  queue: '排队代办',
  errand: '其他代办'
}

const getStatusLabel = (status: OrderStatus) => statusLabels[status] || status
const getTaskTypeLabel = (type?: TaskType) => type ? taskTypeLabels[type] : '代办'

const formatTime = (time: string) => dayjs(time).format('MM-DD HH:mm')

const inProgressOrders = computed(() =>
  orders.value.filter(o => ['accepted', 'in_progress'].includes(o.status))
)

const completedOrders = computed(() =>
  orders.value.filter(o => o.status === 'completed')
)

const fetchOrders = async () => {
  loading.value = true
  try {
    const res = await orderApi.list({ page_size: 50 })
    if (res.code === 200) {
      orders.value = res.data.items
    }
  } catch (error) {
    console.error('Failed to fetch orders:', error)
  } finally {
    loading.value = false
  }
}

const goToDetail = (id: number) => {
  router.push(`/orders/${id}`)
}

watch(activeTab, fetchOrders)

onMounted(() => {
  fetchOrders()
})
</script>

<style lang="scss" scoped>
.order-list-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;

  .tabs {
    background: #fff;
    border-radius: 8px;
    padding: 16px;
  }

  .order-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(300px, 1fr));
    gap: 16px;
  }

  .order-card {
    cursor: pointer;
    transition: transform 0.3s;

    &:hover {
      transform: translateY(-2px);
    }
  }

  .order-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 12px;

    .order-time {
      color: #909399;
      font-size: 12px;
    }
  }

  .order-title {
    font-size: 16px;
    font-weight: 600;
    margin-bottom: 10px;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
  }

  .order-info {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 8px;

    .task-type {
      color: #667eea;
      font-size: 13px;
    }

    .reward {
      color: #ff5722;
      font-size: 18px;
      font-weight: bold;
    }
  }

  .order-location {
    display: flex;
    align-items: center;
    color: #909399;
    font-size: 13px;

    .el-icon {
      margin-right: 4px;
    }
  }
}

.empty-state {
  padding: 40px;
}

@media (max-width: 768px) {
  .order-grid {
    grid-template-columns: 1fr !important;
  }
}
</style>
