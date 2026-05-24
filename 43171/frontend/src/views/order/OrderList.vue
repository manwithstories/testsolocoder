<template>
  <el-card>
    <template #header>
      <div class="header">
        <span>租赁订单</span>
        <el-select v-model="filterStatus" placeholder="全部状态" style="width: 150px" @change="fetchOrders">
          <el-option value="" label="全部" />
          <el-option value="pending" label="待支付" />
          <el-option value="paid" label="已支付" />
          <el-option value="picked" label="已取机" />
          <el-option value="returned" label="已归还" />
          <el-option value="completed" label="已完成" />
          <el-option value="cancelled" label="已取消" />
        </el-select>
      </div>
    </template>

    <el-table :data="orders" v-loading="loading">
      <el-table-column prop="order_no" label="订单号" width="180" />
      <el-table-column label="设备" min-width="150">
        <template #default="{ row }">{{ row.drone?.name || '-' }}</template>
      </el-table-column>
      <el-table-column label="租赁日期">
        <template #default="{ row }">{{ row.start_date }} ~ {{ row.end_date }}</template>
      </el-table-column>
      <el-table-column prop="total_amount" label="金额" width="120">
        <template #default="{ row }">¥{{ row.total_amount }}</template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.status)">{{ statusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="创建时间" width="180" />
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="$router.push(`/order/${row.id}`)">详情</el-button>
          <el-button v-if="row.status === 'pending'" type="success" link @click="payOrder(row)">支付</el-button>
          <el-button v-if="row.status === 'paid'" type="warning" link @click="pickup(row)">取机</el-button>
          <el-button v-if="row.status === 'picked'" type="success" link @click="confirmReturn(row)">归还</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      style="margin-top: 20px; justify-content: flex-end"
      layout="prev, pager, next, total"
      :total="total"
      :page-size="pageSize"
      v-model:current-page="currentPage"
      @current-change="fetchOrders"
    />
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const loading = ref(false)
const orders = ref<RentalOrder[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)
const filterStatus = ref('')

onMounted(() => {
  fetchOrders()
})

async function fetchOrders() {
  loading.value = true
  try {
    const res: any = await request.get('/my-orders', {
      params: {
        page: currentPage.value,
        page_size: pageSize.value,
        status: filterStatus.value
      }
    })
    orders.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function payOrder(row: RentalOrder) {
  ElMessageBox.confirm(`确定支付订单"${row.order_no}"金额 ¥${row.total_amount}？`, '支付确认', {
    type: 'warning'
  }).then(async () => {
    try {
      await request.post('/orders/pay', { order_id: row.id, pay_type: 'balance' })
      ElMessage.success('支付成功')
      fetchOrders()
    } catch (e: any) {
      ElMessage.error(e.message || '支付失败')
    }
  }).catch(() => {})
}

async function pickup(row: RentalOrder) {
  try {
    await request.post(`/orders/${row.id}/pickup`)
    ElMessage.success('取机成功')
    fetchOrders()
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  }
}

async function confirmReturn(row: RentalOrder) {
  ElMessageBox.confirm('确定归还设备吗？', '归还确认', {
    type: 'warning'
  }).then(async () => {
    try {
      await request.post('/orders/confirm-return', { order_id: row.id })
      ElMessage.success('归还成功')
      fetchOrders()
    } catch (e: any) {
      ElMessage.error(e.message || '操作失败')
    }
  }).catch(() => {})
}

function statusText(status: string) {
  const map: Record<string, string> = {
    pending: '待支付', paid: '已支付', picked: '已取机',
    returned: '已归还', completed: '已完成', cancelled: '已取消'
  }
  return map[status] || status
}

function statusTagType(status: string) {
  const map: Record<string, string> = {
    pending: 'warning', paid: 'primary', picked: 'success',
    returned: 'info', completed: 'success', cancelled: 'info'
  }
  return map[status] || ''
}
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
