<template>
  <div class="admin-refunds-page">
    <div class="page-header">
      <h2 class="page-title">退款审核</h2>
    </div>

    <el-card>
      <el-table :data="orders" style="width: 100%">
        <el-table-column prop="order_no" label="工单号" width="200" />
        <el-table-column prop="title" label="服务内容" min-width="200" />
        <el-table-column label="客户" width="150">
          <template #default="{ row }">{{ row.customer?.username || '-' }}</template>
        </el-table-column>
        <el-table-column prop="refund_reason" label="退款原因" min-width="200" />
        <el-table-column prop="refund_amount" label="退款金额" width="120">
          <template #default="{ row }">¥{{ row.refund_amount }}</template>
        </el-table-column>
        <el-table-column prop="final_price" label="原费用" width="120">
          <template #default="{ row }">¥{{ row.final_price }}</template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button type="success" size="small" @click="approve(row)">通过</el-button>
            <el-button type="danger" size="small" @click="reject(row)">拒绝</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div v-if="orders.length === 0" class="empty-state">
        <el-empty description="暂无退款申请" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { adminOrderApi } from '@/api/order'
import type { Order } from '@/types'

const orders = ref<Order[]>([])

onMounted(() => {
  loadRefunds()
})

async function loadRefunds() {
  try {
    const res = await adminOrderApi.getRefundList({ page: 1, page_size: 50 })
    orders.value = res.data?.list || []
  } catch (error) {
    console.error('Failed to load refunds:', error)
  }
}

async function approve(order: Order) {
  try {
    await ElMessageBox.confirm(
      `确定要通过该退款申请吗？退款金额：¥${order.refund_amount}`,
      '确认',
      { type: 'warning' }
    )
    await adminOrderApi.approveRefund(order.id)
    ElMessage.success('退款已通过')
    loadRefunds()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to approve:', error)
    }
  }
}

async function reject(order: Order) {
  try {
    const { value: reason } = await ElMessageBox.prompt('请输入拒绝原因', '拒绝退款', {
      inputValidator: (value) => {
        if (!value) return '请输入拒绝原因'
        return true
      }
    })
    await adminOrderApi.rejectRefund(order.id, { reason })
    ElMessage.success('已拒绝退款申请')
    loadRefunds()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to reject:', error)
    }
  }
}
</script>

<style scoped>
.admin-refunds-page {
  padding: 0;
}

.empty-state {
  padding: 40px;
}
</style>
