<template>
  <AppLayout>
    <div class="page">
      <div class="row" style="margin-bottom:16px">
        <h2 style="margin:0">财务结算</h2>
        <el-date-picker v-model="month" type="month" placeholder="选择月份" value-format="YYYY-MM" @change="load" />
        <el-button type="primary" :icon="Download" @click="exportCSV">导出CSV</el-button>
      </div>
      <el-descriptions :column="3" bordered>
        <el-descriptions-item label="订单数">{{ data.order_count || 0 }}</el-descriptions-item>
        <el-descriptions-item label="总收入">¥{{ data.income?.toFixed(2) || '0.00' }}</el-descriptions-item>
        <el-descriptions-item label="支出">¥{{ data.payout?.toFixed(2) || '0.00' }}</el-descriptions-item>
        <el-descriptions-item label="净利润">¥{{ data.net?.toFixed(2) || '0.00' }}</el-descriptions-item>
      </el-descriptions>

      <el-table :data="data.settlements || []" border stripe style="margin-top:16px">
        <el-table-column prop="id" label="编号" width="80" />
        <el-table-column prop="order_id" label="订单号" width="100" />
        <el-table-column prop="total_amount" label="总金额" width="120" />
        <el-table-column prop="company_share" label="公司分成" width="120" />
        <el-table-column prop="staff_share" label="人员分成" width="120" />
        <el-table-column prop="status" label="状态" width="100" />
        <el-table-column prop="created_at" label="时间" width="180" />
      </el-table>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { Download } from '@element-plus/icons-vue'
import AppLayout from '../../components/AppLayout.vue'
import { companyMonthly, exportFinanceCSV } from '../../api/finance'

const month = ref<string>('')
const data = ref<any>({})

async function load() {
  const res = await companyMonthly(month.value || undefined)
  data.value = (res.data as any).data || {}
}

async function exportCSV() {
  const res: any = await exportFinanceCSV()
  const blob = new Blob([res.data], { type: 'text/csv' })
  const url = URL.createObjectURL(blob)
  const a = document.createElement('a')
  a.href = url
  a.download = `finance-${month.value || 'all'}.csv`
  a.click()
  URL.revokeObjectURL(url)
}

onMounted(load)
</script>
