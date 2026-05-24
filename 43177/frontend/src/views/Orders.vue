<template>
  <div class="orders-page">
    <div class="container">
      <div class="page-header">
        <h2 class="page-title">{{ userStore.isTechnician ? '接单中心' : '我的工单' }}</h2>
        <el-button
          v-if="userStore.isCustomer"
          type="primary"
          @click="router.push('/orders/create')"
        >
          创建工单
        </el-button>
      </div>

      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="全部" name="all">
          <div class="order-list">
            <el-card v-for="order in orders" :key="order.id" class="order-card" shadow="hover">
              <div class="order-header">
                <div class="order-title">
                  <h3>{{ order.title }}</h3>
                  <el-tag :class="`status-${order.status}`">
                    {{ getStatusText(order.status) }}
                  </el-tag>
                </div>
                <div class="order-no">{{ order.order_no }}</div>
              </div>
              <div class="order-info">
                <div class="info-item">
                  <span class="label">服务项目：</span>
                  <span>{{ order.service_item?.name }}</span>
                </div>
                <div class="info-item">
                  <span class="label">联系人：</span>
                  <span>{{ order.contact_name }} {{ order.contact_phone }}</span>
                </div>
                <div class="info-item">
                  <span class="label">地址：</span>
                  <span>{{ order.address }}</span>
                </div>
                <div class="info-item" v-if="order.final_price > 0">
                  <span class="label">费用：</span>
                  <span class="price">¥{{ order.final_price }}</span>
                </div>
              </div>
              <div class="order-footer">
                <span class="order-time">
                  {{ formatTime(order.created_at) }}
                </span>
                <div class="order-actions">
                  <el-button size="small" @click="viewOrder(order)">查看详情</el-button>
                  <el-button
                    v-if="order.status === 'assigned' && userStore.isTechnician"
                    type="primary"
                    size="small"
                    @click="acceptOrder(order)"
                  >
                    接单
                  </el-button>
                  <el-button
                    v-if="order.status === 'accepted' && userStore.isTechnician"
                    type="success"
                    size="small"
                    @click="arriveAtSite(order)"
                  >
                    已到达
                  </el-button>
                  <el-button
                    v-if="order.status === 'on_site' && userStore.isTechnician"
                    type="warning"
                    size="small"
                    @click="startRepair(order)"
                  >
                    开始维修
                  </el-button>
                  <el-button
                    v-if="order.status === 'repairing' && userStore.isTechnician"
                    type="success"
                    size="small"
                    @click="completeOrder(order)"
                  >
                    完工
                  </el-button>
                  <el-button
                    v-if="(order.status === 'pending' || order.status === 'assigned') && userStore.isCustomer"
                    type="danger"
                    size="small"
                    @click="cancelOrder(order)"
                  >
                    取消
                  </el-button>
                </div>
              </div>
            </el-card>
          </div>
        </el-tab-pane>

        <el-tab-pane label="待处理" name="pending">
          <div class="order-list">
            <el-card v-for="order in filteredOrders" :key="order.id" class="order-card" shadow="hover">
              <div class="order-header">
                <div class="order-title">
                  <h3>{{ order.title }}</h3>
                  <el-tag :class="`status-${order.status}`">
                    {{ getStatusText(order.status) }}
                  </el-tag>
                </div>
              </div>
              <div class="order-info">
                <div class="info-item">
                  <span class="label">服务项目：</span>
                  <span>{{ order.service_item?.name }}</span>
                </div>
                <div class="info-item">
                  <span class="label">地址：</span>
                  <span>{{ order.address }}</span>
                </div>
              </div>
              <div class="order-footer">
                <el-button size="small" @click="viewOrder(order)">查看详情</el-button>
              </div>
            </el-card>
          </div>
        </el-tab-pane>

        <el-tab-pane label="进行中" name="processing">
          <div class="order-list">
            <el-card v-for="order in filteredOrders" :key="order.id" class="order-card" shadow="hover">
              <div class="order-header">
                <div class="order-title">
                  <h3>{{ order.title }}</h3>
                  <el-tag :class="`status-${order.status}`">
                    {{ getStatusText(order.status) }}
                  </el-tag>
                </div>
              </div>
              <div class="order-info">
                <div class="info-item">
                  <span class="label">服务项目：</span>
                  <span>{{ order.service_item?.name }}</span>
                </div>
              </div>
              <div class="order-footer">
                <el-button size="small" @click="viewOrder(order)">查看详情</el-button>
              </div>
            </el-card>
          </div>
        </el-tab-pane>

        <el-tab-pane label="已完成" name="completed">
          <div class="order-list">
            <el-card v-for="order in filteredOrders" :key="order.id" class="order-card" shadow="hover">
              <div class="order-header">
                <div class="order-title">
                  <h3>{{ order.title }}</h3>
                  <el-tag class="status-completed">已完成</el-tag>
                </div>
              </div>
              <div class="order-info">
                <div class="info-item">
                  <span class="label">费用：</span>
                  <span class="price">¥{{ order.final_price }}</span>
                </div>
              </div>
              <div class="order-footer">
                <el-button size="small" @click="viewOrder(order)">查看详情</el-button>
                <el-button
                  v-if="order.status === 'completed' && userStore.isCustomer"
                  type="primary"
                  size="small"
                  @click="reviewOrder(order)"
                >
                  评价
                </el-button>
              </div>
            </el-card>
          </div>
        </el-tab-pane>
      </el-tabs>

      <div v-if="orders.length === 0 && loading === false" class="empty-state">
        <el-empty description="暂无工单" />
      </div>

      <div v-if="total > pageSize" class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="loadOrders"
        />
      </div>
    </div>

    <el-dialog v-model="reviewDialogVisible" title="评价工单" width="500px">
      <el-form :model="reviewForm" label-width="80px">
        <el-form-item label="评分">
          <el-rate v-model="reviewForm.rating" />
        </el-form-item>
        <el-form-item label="评价内容">
          <el-input v-model="reviewForm.content" type="textarea" :rows="4" placeholder="请输入评价内容" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="reviewDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="submitReview" :loading="submitting">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/store/user'
import { orderApi } from '@/api/order'
import { userApi } from '@/api/user'
import type { Order } from '@/types'
import dayjs from 'dayjs'

const router = useRouter()
const userStore = useUserStore()

const activeTab = ref('all')
const orders = ref<Order[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const loading = ref(false)
const reviewDialogVisible = ref(false)
const submitting = ref(false)
const currentOrder = ref<Order | null>(null)

const reviewForm = ref({
  rating: 5,
  content: ''
})

const filteredOrders = computed(() => {
  const status = activeTab.value
  if (status === 'pending') {
    return orders.value.filter(o => ['pending', 'assigned'].includes(o.status))
  } else if (status === 'processing') {
    return orders.value.filter(o => ['accepted', 'on_site', 'repairing'].includes(o.status))
  } else if (status === 'completed') {
    return orders.value.filter(o => ['completed', 'cancelled', 'refunded'].includes(o.status))
  }
  return orders.value
})

onMounted(() => {
  loadOrders()
})

async function loadOrders() {
  loading.value = true
  try {
    const res = await orderApi.getOrders({
      page: currentPage.value,
      page_size: pageSize.value
    })
    orders.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (error) {
    console.error('Failed to load orders:', error)
  } finally {
    loading.value = false
  }
}

function handleTabChange(tab: string) {
  activeTab.value = tab
}

function getStatusText(status: string): string {
  const statusMap: Record<string, string> = {
    pending: '待分配',
    assigned: '待接单',
    accepted: '已接单',
    on_site: '已到达',
    repairing: '维修中',
    completed: '已完成',
    cancelled: '已取消',
    refunding: '退款中',
    refunded: '已退款'
  }
  return statusMap[status] || status
}

function formatTime(time: string): string {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

function viewOrder(order: Order) {
  router.push(`/orders/${order.id}`)
}

async function acceptOrder(order: Order) {
  try {
    await ElMessageBox.confirm('确定要接这个工单吗？', '确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await orderApi.acceptOrder(order.id)
    ElMessage.success('接单成功')
    loadOrders()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to accept order:', error)
    }
  }
}

async function arriveAtSite(order: Order) {
  try {
    await orderApi.arriveAtSite(order.id)
    ElMessage.success('已标记到达现场')
    loadOrders()
  } catch (error) {
    console.error('Failed to arrive:', error)
  }
}

async function startRepair(order: Order) {
  try {
    await orderApi.startRepair(order.id)
    ElMessage.success('已开始维修')
    loadOrders()
  } catch (error) {
    console.error('Failed to start repair:', error)
  }
}

async function completeOrder(order: Order) {
  try {
    const { value: price } = await ElMessageBox.prompt('请输入最终费用', '完工确认', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputValidator: (value) => {
        if (!value || isNaN(Number(value)) || Number(value) <= 0) {
          return '请输入有效的金额'
        }
        return true
      }
    })
    await orderApi.completeOrder(order.id, { final_price: Number(price) })
    ElMessage.success('工单已完工')
    loadOrders()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to complete order:', error)
    }
  }
}

async function cancelOrder(order: Order) {
  try {
    const { value: reason } = await ElMessageBox.prompt('请输入取消原因', '取消工单', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputValidator: (value) => {
        if (!value) {
          return '请输入取消原因'
        }
        return true
      }
    })
    await orderApi.cancelOrder(order.id, { cancel_reason: reason })
    ElMessage.success('工单已取消')
    loadOrders()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to cancel order:', error)
    }
  }
}

function reviewOrder(order: Order) {
  currentOrder.value = order
  reviewForm.value = { rating: 5, content: '' }
  reviewDialogVisible.value = true
}

async function submitReview() {
  if (!currentOrder.value) return

  submitting.value = true
  try {
    await userApi.createReview({
      order_id: currentOrder.value.id,
      rating: reviewForm.value.rating,
      content: reviewForm.value.content
    })
    ElMessage.success('评价成功')
    reviewDialogVisible.value = false
    loadOrders()
  } catch (error) {
    console.error('Failed to submit review:', error)
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.orders-page {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.order-list {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.order-card {
  transition: all 0.3s;
}

.order-card:hover {
  transform: translateY(-2px);
}

.order-header {
  display: flex;
  justify-content: space-between;
  align-items: flex-start;
  margin-bottom: 15px;
}

.order-title {
  display: flex;
  align-items: center;
  gap: 10px;
}

.order-title h3 {
  margin: 0;
  font-size: 16px;
  color: #303133;
}

.order-no {
  font-size: 12px;
  color: #909399;
}

.order-info {
  margin-bottom: 15px;
}

.info-item {
  margin-bottom: 8px;
  font-size: 14px;
  color: #606266;
}

.label {
  color: #909399;
}

.price {
  color: #f56c6c;
  font-weight: 600;
}

.order-footer {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.order-time {
  font-size: 12px;
  color: #909399;
}

.order-actions {
  display: flex;
  gap: 10px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 30px;
}
</style>
