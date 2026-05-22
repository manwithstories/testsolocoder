<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">订单管理</h1>
      <el-button type="primary" :icon="Download" @click="exportOrders">导出Excel</el-button>
    </div>

    <div class="search-bar">
      <el-select v-model="filters.status" placeholder="订单状态" clearable style="width: 140px">
        <el-option label="待处理" value="pending" />
        <el-option label="已支付" value="paid" />
        <el-option label="已完成" value="completed" />
        <el-option label="已退款" value="refunded" />
        <el-option label="已取消" value="cancelled" />
      </el-select>
      <el-date-picker
        v-model="filters.dateRange"
        type="daterange"
        range-separator="至"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
        value-format="YYYY-MM-DD"
        style="width: 260px"
      />
      <el-button type="primary" @click="loadOrders">搜索</el-button>
    </div>

    <el-table :data="orders" v-loading="loading" style="width: 100%">
      <el-table-column prop="order_no" label="订单号" width="180" />
      <el-table-column label="用户" width="120">
        <template #default="{ row }">
          {{ row.user?.username }}
        </template>
      </el-table-column>
      <el-table-column label="车辆" min-width="150">
        <template #default="{ row }">
          {{ row.car?.brand }} {{ row.car?.model }}
        </template>
      </el-table-column>
      <el-table-column label="预订号" width="180">
        <template #default="{ row }">
          {{ row.booking?.booking_no }}
        </template>
      </el-table-column>
      <el-table-column label="金额" width="120" align="right">
        <template #default="{ row }">
          ¥{{ row.final_amount.toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column label="支付状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getPaymentStatusType(row.payment_status)">
            {{ getPaymentStatusText(row.payment_status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="订单状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getOrderStatusType(row.status)">
            {{ getOrderStatusText(row.status) }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column label="创建时间" width="170">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="180" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.payment_status === 'pending'"
            type="success"
            link
            size="small"
            @click="updateStatus(row, 'paid')"
          >
            确认支付
          </el-button>
          <el-button
            v-if="row.payment_status === 'paid'"
            type="warning"
            link
            size="small"
            @click="handleRefund(row)"
          >
            退款
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div style="text-align: center; margin-top: 20px;">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        layout="total, prev, pager, next"
        @current-change="loadOrders"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Download } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { orderApi } from '@/api'
import type { Order } from '@/types'

const orders = ref<Order[]>([])
const loading = ref(false)

const filters = reactive({
  status: '',
  dateRange: [] as string[]
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

onMounted(() => {
  loadOrders()
})

const loadOrders = async () => {
  loading.value = true
  try {
    const params: Record<string, any> = {
      page: pagination.page,
      page_size: pagination.pageSize,
      status: filters.status
    }
    if (filters.dateRange && filters.dateRange.length === 2) {
      params.start_date = filters.dateRange[0]
      params.end_date = filters.dateRange[1]
    }
    const res = await orderApi.getOrders(params)
    orders.value = res.data.items
    pagination.total = res.data.total
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

const updateStatus = async (row: Order, status: string) => {
  try {
    await orderApi.updateOrderStatus(row.id, status)
    ElMessage.success('操作成功')
    loadOrders()
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  }
}

const handleRefund = async (row: Order) => {
  try {
    const { value: reason } = await ElMessageBox.prompt('请输入退款原因', '退款', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputPlaceholder: '请输入退款原因'
    })

    await orderApi.refundOrder(row.id, reason)
    ElMessage.success('退款成功')
    loadOrders()
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error(err.message || '操作失败')
    }
  }
}

const exportOrders = async () => {
  try {
    const params: Record<string, any> = {
      status: filters.status
    }
    if (filters.dateRange && filters.dateRange.length === 2) {
      params.start_date = filters.dateRange[0]
      params.end_date = filters.dateRange[1]
    }
    const res = await orderApi.exportOrders(params)
    window.open(res.data.file_path, '_blank')
  } catch (err: any) {
    ElMessage.error(err.message || '导出失败')
  }
}

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const getPaymentStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    paid: 'success',
    failed: 'danger',
    refunded: 'info'
  }
  return map[status] || 'info'
}

const getPaymentStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待支付',
    paid: '已支付',
    failed: '支付失败',
    refunded: '已退款'
  }
  return map[status] || status
}

const getOrderStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    paid: 'primary',
    completed: 'success',
    refunded: 'info',
    cancelled: 'danger'
  }
  return map[status] || 'info'
}

const getOrderStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待处理',
    paid: '已支付',
    completed: '已完成',
    refunded: '已退款',
    cancelled: '已取消'
  }
  return map[status] || status
}
</script>
