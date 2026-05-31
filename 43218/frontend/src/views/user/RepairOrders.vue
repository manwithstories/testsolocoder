<template>
  <div class="repair-orders-page">
    <div class="page-header">
      <h2 class="page-title">维修订单</h2>
      <el-tabs v-model="activeTab" @tab-change="handleTabChange">
        <el-tab-pane label="全部" name="all" />
        <el-tab-pane label="待接单" name="pending" />
        <el-tab-pane label="维修中" name="repairing" />
        <el-tab-pane label="待取件" name="completed" />
        <el-tab-pane label="已完成" name="finished" />
      </el-tabs>
    </div>

    <div class="order-list" v-loading="loading">
      <div
        v-for="order in orders"
        :key="order.id"
        class="order-card card"
      >
        <div class="order-header">
          <span class="order-no">订单号：{{ order.orderNo }}</span>
          <span class="order-time">{{ formatTime(order.createdAt) }}</span>
          <el-tag :type="getStatusType(order.status)">
            {{ getStatusText(order.status) }}
          </el-tag>
        </div>
        <div class="order-content">
          <div class="device-info">
            <h4>{{ order.deviceBrand }} {{ order.deviceModel }}</h4>
            <p>设备类型：{{ order.deviceType }}</p>
            <p>故障描述：{{ order.faultDescription }}</p>
            <p>联系人：{{ order.contactName }} - {{ order.contactPhone }}</p>
            <p v-if="order.address">地址：{{ order.address }}</p>
          </div>
          <div class="order-info">
            <p>服务费用：<span class="price-text">¥{{ order.servicePrice.toFixed(2) }}</span></p>
            <p>最终费用：<span class="price-text">¥{{ order.finalPrice.toFixed(2) }}</span></p>
          </div>
          <div class="order-actions">
            <el-button
              v-if="order.status === 4 && order.buyerId === userStore.userInfo?.id"
              type="success"
              @click="handlePickUp(order)"
            >
              确认取件
            </el-button>
            <el-button
              v-if="order.status === 1 || order.status === 2"
              @click="handleCancel(order)"
            >
              取消订单
            </el-button>
          </div>
        </div>
      </div>

      <div class="empty-state" v-if="!loading && orders.length === 0">
        <el-empty description="暂无维修订单" />
      </div>
    </div>

    <div class="pagination-wrapper" v-if="total > 0">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :page-sizes="[10, 20, 50]"
        :total="total"
        layout="total, sizes, prev, pager, next, jumper"
        @size-change="fetchOrders"
        @current-change="fetchOrders"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { useUserStore } from '@/store/user'
import { repairApi } from '@/api'
import { RepairStatus, RepairStatusText } from '@/types'
import type { RepairOrder } from '@/types'

const userStore = useUserStore()

const loading = ref(false)
const orders = ref<RepairOrder[]>([])
const total = ref(0)
const activeTab = ref('all')

const pagination = reactive({
  page: 1,
  pageSize: 10
})

function getStatusText(status: number): string {
  return RepairStatusText[status] || '未知'
}

function getStatusType(status: number): string {
  const typeMap: Record<number, string> = {
    1: 'warning',
    2: 'primary',
    3: 'primary',
    4: 'success',
    5: 'success',
    6: 'info',
    7: 'warning',
    8: 'danger'
  }
  return typeMap[status] || 'info'
}

function formatTime(time: string): string {
  return new Date(time).toLocaleString()
}

function handleTabChange() {
  pagination.page = 1
  fetchOrders()
}

function getStatusFilter(): number | undefined {
  const statusMap: Record<string, number> = {
    pending: 1,
    repairing: 3,
    completed: 4,
    finished: 5
  }
  return activeTab.value === 'all' ? undefined : statusMap[activeTab.value]
}

async function fetchOrders() {
  loading.value = true
  try {
    const res = await repairApi.getOrderList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      status: getStatusFilter()
    })
    orders.value = res.data
    total.value = res.pagination.total
  } catch (error) {
    console.error('Failed to fetch orders:', error)
  } finally {
    loading.value = false
  }
}

async function handlePickUp(order: RepairOrder) {
  try {
    await ElMessageBox.confirm('确认取件？', '提示', {
      confirmButtonText: '确认',
      cancelButtonText: '取消'
    })
    await repairApi.pickUpDevice(order.id)
    ElMessage.success('已确认取件')
    fetchOrders()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to pick up:', error)
    }
  }
}

async function handleCancel(order: RepairOrder) {
  try {
    await ElMessageBox.confirm('确定取消订单？', '提示', {
      confirmButtonText: '确认',
      cancelButtonText: '取消'
    })
    await repairApi.cancelOrder(order.id)
    ElMessage.success('订单已取消')
    fetchOrders()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to cancel:', error)
    }
  }
}

onMounted(() => {
  fetchOrders()
})
</script>

<style lang="scss" scoped>
.repair-orders-page {
  .order-list {
    .order-card {
      margin-bottom: 16px;

      .order-header {
        display: flex;
        align-items: center;
        gap: 16px;
        padding-bottom: 12px;
        border-bottom: 1px solid #f0f0f0;
        margin-bottom: 12px;

        .order-no {
          font-weight: 500;
        }

        .order-time {
          color: var(--text-lighter-color);
          font-size: 13px;
        }
      }

      .order-content {
        display: flex;
        justify-content: space-between;
        align-items: flex-start;

        .device-info {
          h4 {
            margin-bottom: 8px;
          }

          p {
            margin-bottom: 4px;
            color: var(--text-light-color);
            font-size: 13px;
          }
        }

        .order-info {
          text-align: right;

          p {
            margin-bottom: 8px;
          }
        }

        .order-actions {
          display: flex;
          gap: 8px;
        }
      }
    }
  }

  .pagination-wrapper {
    display: flex;
    justify-content: center;
    padding: 20px 0;
  }
}
</style>
