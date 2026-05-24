<template>
  <el-card v-if="order">
    <el-descriptions :column="2" border title="订单信息">
      <el-descriptions-item label="订单号">{{ order.order_no }}</el-descriptions-item>
      <el-descriptions-item label="状态">
        <el-tag :type="statusTagType(order.status)">{{ statusText(order.status) }}</el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="设备">{{ order.drone?.name }}</el-descriptions-item>
      <el-descriptions-item label="区域">{{ order.region }}</el-descriptions-item>
      <el-descriptions-item label="租赁日期">{{ order.start_date }} ~ {{ order.end_date }}</el-descriptions-item>
      <el-descriptions-item label="归还日期">{{ order.return_date || '未归还' }}</el-descriptions-item>
      <el-descriptions-item label="天数">{{ order.total_days }}天</el-descriptions-item>
      <el-descriptions-item label="日租金">¥{{ order.price_per_day }}</el-descriptions-item>
      <el-descriptions-item label="租金">¥{{ order.rental_fee }}</el-descriptions-item>
      <el-descriptions-item label="押金">¥{{ order.deposit }}</el-descriptions-item>
      <el-descriptions-item label="保险费">¥{{ order.insurance_fee }}</el-descriptions-item>
      <el-descriptions-item label="滞纳金">¥{{ order.late_fee }}</el-descriptions-item>
      <el-descriptions-item label="总金额">¥{{ order.total_amount }}</el-descriptions-item>
      <el-descriptions-item label="已支付">¥{{ order.paid_amount }}</el-descriptions-item>
      <el-descriptions-item label="联系人">{{ order.contact_name }} ({{ order.contact_phone }})</el-descriptions-item>
      <el-descriptions-item label="地址">{{ order.address }}</el-descriptions-item>
      <el-descriptions-item label="备注" :span="2">{{ order.remark || '-' }}</el-descriptions-item>
      <el-descriptions-item label="创建时间">{{ order.created_at }}</el-descriptions-item>
    </el-descriptions>

    <el-divider />

    <div class="actions">
      <el-button v-if="order.status === 'pending'" type="primary" @click="payOrder">支付</el-button>
      <el-button v-if="order.status === 'paid'" type="warning" @click="pickup">取机</el-button>
      <el-button v-if="order.status === 'picked'" type="success" @click="confirmReturn">归还</el-button>
      <el-button v-if="order.status === 'returned'" type="primary" @click="complete">完成</el-button>
      <el-button v-if="order.status === 'pending' || order.status === 'paid'" type="danger" @click="cancel">取消订单</el-button>
      <el-button @click="$router.back()">返回</el-button>
    </div>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const route = useRoute()
const router = useRouter()
const orderId = route.params.id as string

const order = ref<RentalOrder | null>(null)

onMounted(() => {
  fetchOrder()
})

async function fetchOrder() {
  try {
    const res: any = await request.get(`/orders/${orderId}`)
    order.value = res.data
  } catch (e) {
    console.error(e)
  }
}

async function payOrder() {
  ElMessageBox.confirm(`确定支付金额 ¥${order.value?.total_amount}？`, '支付确认', {
    type: 'warning'
  }).then(async () => {
    try {
      await request.post('/orders/pay', { order_id: order.value?.id, pay_type: 'balance' })
      ElMessage.success('支付成功')
      fetchOrder()
    } catch (e: any) {
      ElMessage.error(e.message || '支付失败')
    }
  }).catch(() => {})
}

async function pickup() {
  try {
    await request.post(`/orders/${orderId}/pickup`)
    ElMessage.success('取机成功')
    fetchOrder()
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  }
}

async function confirmReturn() {
  ElMessageBox.confirm('确定归还设备吗？', '归还确认', { type: 'warning' }).then(async () => {
    try {
      await request.post('/orders/confirm-return', { order_id: order.value?.id })
      ElMessage.success('归还成功')
      fetchOrder()
    } catch (e: any) {
      ElMessage.error(e.message || '操作失败')
    }
  }).catch(() => {})
}

async function complete() {
  try {
    await request.post(`/orders/${orderId}/complete`)
    ElMessage.success('订单已完成')
    fetchOrder()
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  }
}

async function cancel() {
  const { value: reason } = await ElMessageBox.prompt('请输入取消原因', '取消订单', {
    confirmButtonText: '确定',
    cancelButtonText: '取消'
  })
  try {
    await request.post('/orders/cancel', { order_id: order.value?.id, cancel_reason: reason })
    ElMessage.success('订单已取消')
    router.push('/orders')
  } catch (e: any) {
    ElMessage.error(e.message || '取消失败')
  }
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
.actions {
  display: flex;
  gap: 10px;
}
</style>
