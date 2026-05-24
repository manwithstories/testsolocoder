<template>
  <div class="admin-orders-page">
    <div class="page-header">
      <h2 class="page-title">工单管理</h2>
    </div>

    <el-card>
      <div class="filter-bar">
        <el-select v-model="statusFilter" placeholder="选择状态" clearable style="width: 150px;" @change="loadOrders">
          <el-option label="待分配" value="pending" />
          <el-option label="待接单" value="assigned" />
          <el-option label="已接单" value="accepted" />
          <el-option label="已到达" value="on_site" />
          <el-option label="维修中" value="repairing" />
          <el-option label="已完成" value="completed" />
          <el-option label="已取消" value="cancelled" />
          <el-option label="退款中" value="refunding" />
          <el-option label="已退款" value="refunded" />
        </el-select>
      </div>

      <el-table :data="orders" style="width: 100%">
        <el-table-column prop="order_no" label="工单号" width="200" />
        <el-table-column prop="title" label="服务内容" min-width="200" />
        <el-table-column label="客户" width="150">
          <template #default="{ row }">{{ row.customer?.username || '-' }}</template>
        </el-table-column>
        <el-table-column label="技师" width="150">
          <template #default="{ row }">{{ row.technician?.username || '-' }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :class="`status-${row.status}`" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="final_price" label="费用" width="100">
          <template #default="{ row }">
            {{ row.final_price ? '¥' + row.final_price : '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="100">
          <template #default="{ row }">
            <el-button size="small" @click="viewOrder(row)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div v-if="total > pageSize" class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="loadOrders"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { adminOrderApi } from '@/api/order'
import type { Order } from '@/types'
import dayjs from 'dayjs'

const router = useRouter()

const orders = ref<Order[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const statusFilter = ref('')

onMounted(() => {
  loadOrders()
})

async function loadOrders() {
  try {
    const res = await adminOrderApi.getAllOrders({
      page: currentPage.value,
      page_size: pageSize.value,
      status: statusFilter.value || undefined
    })
    orders.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (error) {
    console.error('Failed to load orders:', error)
  }
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
</script>

<style scoped>
.admin-orders-page {
  padding: 0;
}

.filter-bar {
  display: flex;
  gap: 15px;
  margin-bottom: 20px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}
</style>
