<template>
  <div class="finance-dashboard">
    <div class="page-header">
      <h2 class="page-title">财务结算</h2>
    </div>

    <el-row :gutter="16" class="stat-cards">
      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card">
          <div class="stat-label">总收入</div>
          <div class="stat-value" style="color: #52c41a">¥{{ summary.total_income }}</div>
        </div>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card">
          <div class="stat-label">总支出</div>
          <div class="stat-value" style="color: #ff4d4f">¥{{ summary.total_expense }}</div>
        </div>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card">
          <div class="stat-label">净余额</div>
          <div class="stat-value" :style="{ color: summary.net_balance >= 0 ? '#52c41a' : '#ff4d4f' }">
            ¥{{ summary.net_balance }}
          </div>
        </div>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card">
          <div class="stat-label">平台服务费</div>
          <div class="stat-value">¥{{ summary.platform_fee }}</div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :xs="24" :md="16">
        <div class="card-container">
          <div class="card-header">
            <span>最近交易</span>
            <el-link type="primary" @click="$router.push('/transactions')">查看全部</el-link>
          </div>
          <el-table :data="transactions" style="width: 100%">
            <el-table-column prop="description" label="描述" show-overflow-tooltip />
            <el-table-column label="金额" width="120">
              <template #default="{ row }">
                <span :style="{ color: row.transaction_type === 'income' ? '#52c41a' : '#ff4d4f' }">
                  {{ row.transaction_type === 'income' ? '+' : '-' }}{{ row.currency }} {{ row.amount }}
                </span>
              </template>
            </el-table-column>
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'completed' ? 'success' : 'warning'" size="small">
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="时间" width="180">
              <template #default="{ row }">
                {{ formatDateTime(row.created_at) }}
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-col>
      <el-col :xs="24" :md="8">
        <div class="card-container">
          <div class="card-header">
            <span>操作</span>
          </div>
          <div class="quick-actions">
            <el-button type="primary" size="large" @click="$router.push('/transactions')">
              <el-icon><Document /></el-icon>
              交易记录
            </el-button>
            <el-button size="large" @click="$router.push('/settlements')">
              <el-icon><Tickets /></el-icon>
              结算单
            </el-button>
            <el-button size="large" @click="exportReport">
              <el-icon><Download /></el-icon>
              导出报表
            </el-button>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getFinancialSummaryApi, getTransactionsApi, exportFinancialReportApi } from '@/api/finance'
import type { Transaction, FinancialSummary } from '@/types/finance'
import dayjs from 'dayjs'

const summary = ref<FinancialSummary>({
  total_income: 0,
  total_expense: 0,
  net_balance: 0,
  platform_fee: 0,
  dock_fee: 0
})

const transactions = ref<Transaction[]>([])

const formatDateTime = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm')

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待处理',
    completed: '已完成',
    failed: '失败',
    refunded: '已退款'
  }
  return map[status] || status
}

const fetchData = async () => {
  try {
    const res: any = await getFinancialSummaryApi()
    summary.value = res.data || summary.value
  } catch (error) {
    console.error('Failed to fetch summary:', error)
  }

  try {
    const res: any = await getTransactionsApi({ page_size: 5 })
    transactions.value = res.data || []
  } catch (error) {
    console.error('Failed to fetch transactions:', error)
  }
}

const exportReport = async () => {
  try {
    const params = {
      start_date: dayjs().subtract(30, 'day').format('YYYY-MM-DD'),
      end_date: dayjs().format('YYYY-MM-DD'),
      format: 'csv' as const
    }
    const res: any = await exportFinancialReportApi(params)
    const blob = new Blob([res], { type: 'text/csv' })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `financial_report_${dayjs().format('YYYYMMDD')}.csv`
    a.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('报表导出成功')
  } catch (error) {
    ElMessage.error('导出失败')
  }
}

onMounted(fetchData)
</script>

<style lang="scss" scoped>
.finance-dashboard {
  .card-container {
    background: #fff;
    border-radius: 8px;
    padding: 20px;
    margin-bottom: 16px;

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;
      font-weight: 600;
    }
  }

  .quick-actions {
    display: flex;
    flex-direction: column;
    gap: 12px;

    .el-button {
      justify-content: flex-start;

      .el-icon {
        margin-right: 8px;
      }
    }
  }
}
</style>
