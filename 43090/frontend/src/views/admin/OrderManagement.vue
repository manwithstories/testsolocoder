<template>
  <div class="order-management">
    <h3>订单管理</h3>
    <el-card class="filter-card">
      <el-input
        v-model="keyword"
        placeholder="搜索订单号/商品名"
        style="width: 300px"
        clearable
        @keyup.enter="fetchOrders"
      />
    </el-card>
    <el-table :data="orders" v-loading="loading">
      <el-table-column prop="order_no" label="订单号" width="180" />
      <el-table-column label="拍卖品" min-width="200">
        <template #default="{ row }">
          <span>{{ row.auction_item?.title }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="price" label="价格" width="120">
        <template #default="{ row }">¥{{ row.price.toFixed(2) }}</template>
      </el-table-column>
      <el-table-column label="买家" width="120">
        <template #default="{ row }">{{ row.buyer?.username }}</template>
      </el-table-column>
      <el-table-column label="卖家" width="120">
        <template #default="{ row }">{{ row.seller?.username }}</template>
      </el-table-column>
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusType(row.status)" size="small">{{ statusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="160">
        <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
      </el-table-column>
    </el-table>
    <el-pagination
      class="pagination"
      v-model:current-page="page"
      v-model:page-size="pageSize"
      :total="total"
      layout="total, prev, pager, next"
      @current-change="fetchOrders"
    />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import dayjs from 'dayjs'
import type { Order } from '@/types'
import { adminApi } from '@/api'

const orders = ref<Order[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const keyword = ref('')

const fetchOrders = async () => {
  loading.value = true
  try {
    const res = await adminApi.getAllOrders({ page: page.value, page_size: pageSize.value, keyword: keyword.value })
    orders.value = res.list
    total.value = res.total
  } finally {
    loading.value = false
  }
}

const statusType = (status: number) => {
  const map: Record<number, string> = {
    0: 'warning', 1: 'primary', 2: 'info', 3: 'success', 4: 'success', 5: 'danger',
  }
  return map[status] || 'info'
}

const statusText = (status: number) => {
  const map: Record<number, string> = {
    0: '待支付', 1: '已支付', 2: '已发货', 3: '已送达', 4: '已完成', 5: '已取消',
  }
  return map[status] || ''
}

const formatTime = (time: string) => dayjs(time).format('YYYY-MM-DD HH:mm')

onMounted(() => {
  fetchOrders()
})
</script>

<style scoped>
h3 {
  margin: 0 0 20px;
}

.filter-card {
  margin-bottom: 20px;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: center;
}
</style>
