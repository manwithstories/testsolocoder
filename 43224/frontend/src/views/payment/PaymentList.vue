<template>
  <div class="page-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>费用管理</span>
        </div>
      </template>

      <div class="stats-bar">
        <el-statistic title="总收入(元)" :value="stats.total_revenue" :precision="2" />
        <el-statistic title="已支付(元)" :value="stats.paid_amount" :precision="2" />
        <el-statistic title="待支付(元)" :value="stats.pending_amount" :precision="2" />
      </div>

      <el-table :data="payments" stripe v-loading="loading" style="margin-top: 20px">
        <el-table-column prop="project.title" label="项目名称" min-width="150" />
        <el-table-column prop="client.username" label="客户" width="120" />
        <el-table-column prop="translator.username" label="译者" width="120" />
        <el-table-column prop="amount" label="金额(元)" width="120">
          <template #default="{ row }">{{ row.amount?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="row.status === 'paid' ? 'success' : 'warning'" size="small">
              {{ row.status === 'paid' ? '已支付' : '待支付' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="160">
          <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="120" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'pending'"
              type="success"
              link
              @click="handleConfirm(row)"
            >确认支付</el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { listPayments, confirmPayment, getPaymentStatistics } from '@/api/statistics'
import dayjs from 'dayjs'

const payments = ref<any[]>([])
const loading = ref(false)
const stats = reactive({
  total_revenue: 0,
  paid_amount: 0,
  pending_amount: 0
})

async function loadData() {
  loading.value = true
  try {
    const [paymentsRes, statsRes] = await Promise.all([
      listPayments(),
      getPaymentStatistics()
    ])
    if (Array.isArray(paymentsRes)) {
      payments.value = paymentsRes
    } else {
      payments.value = (paymentsRes as any)?.list || []
    }
    Object.assign(stats, statsRes || {})
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function handleConfirm(row: any) {
  try {
    await ElMessageBox.confirm('确认已支付该费用？', '提示', { type: 'warning' })
    await confirmPayment(row.id)
    ElMessage.success('确认成功')
    loadData()
  } catch (_) {}
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

onMounted(loadData)
</script>

<style lang="scss" scoped>
.page-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .stats-bar {
    display: flex;
    gap: 40px;
    padding: 16px;
    background: #f5f7fa;
    border-radius: 4px;
  }
}
</style>
