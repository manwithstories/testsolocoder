<template>
  <div class="admin-reports-page">
    <div class="page-header">
      <h2 class="page-title">财务报表</h2>
      <div class="header-actions">
        <el-date-picker
          v-model="selectedMonth"
          type="month"
          placeholder="选择月份"
          format="YYYY-MM"
          value-format="YYYY-MM"
          @change="loadReport"
        />
        <el-button type="primary" @click="loadReport">
          查询
        </el-button>
        <el-button type="success" @click="settleIncome">
          结算收入
        </el-button>
      </div>
    </div>

    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-label">总订单数</div>
          <div class="stat-value">{{ report?.total_orders || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-label">总收入</div>
          <div class="stat-value">¥{{ report?.total_revenue || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-label">平台收入</div>
          <div class="stat-value">¥{{ report?.platform_income || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-label">技师收入</div>
          <div class="stat-value">¥{{ report?.technician_pay || 0 }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="stats-row">
      <el-col :span="8">
        <el-card class="stat-card">
          <div class="stat-label">总提现金额</div>
          <div class="stat-value">¥{{ report?.total_withdraw || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="stat-card">
          <div class="stat-label">新增技师</div>
          <div class="stat-value">{{ report?.new_technicians || 0 }}</div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="stat-card">
          <div class="stat-label">新增客户</div>
          <div class="stat-value">{{ report?.new_customers || 0 }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="mt-20">
      <template #header>
        <div class="card-header">技师绩效排名</div>
      </template>
      <el-table :data="performances" style="width: 100%">
        <el-table-column type="index" label="排名" width="80" />
        <el-table-column prop="username" label="技师" width="150" />
        <el-table-column prop="real_name" label="真实姓名" width="150" />
        <el-table-column prop="total_orders" label="完成订单" width="120" />
        <el-table-column prop="total_revenue" label="创收金额" width="150">
          <template #default="{ row }">¥{{ row.total_revenue?.toFixed(2) }}</template>
        </el-table-column>
        <el-table-column prop="avg_rating" label="平均评分" width="120">
          <template #default="{ row }">
            <el-rate :model-value="row.avg_rating" disabled size="small" />
          </template>
        </el-table-column>
        <el-table-column prop="completion_rate" label="完成率" width="120">
          <template #default="{ row }">{{ (row.completion_rate * 100).toFixed(1) }}%</template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { adminFinanceApi } from '@/api/finance'
import type { MonthlyReport } from '@/types'
import dayjs from 'dayjs'

const selectedMonth = ref(dayjs().format('YYYY-MM'))
const report = ref<MonthlyReport | null>(null)
const performances = ref<any[]>([])

onMounted(() => {
  loadReport()
  loadPerformances()
})

async function loadReport() {
  try {
    const res = await adminFinanceApi.getMonthlyReport({ month: selectedMonth.value })
    report.value = res.data || null
  } catch (error) {
    console.error('Failed to load report:', error)
  }
}

async function loadPerformances() {
  try {
    const res = await adminFinanceApi.getTechnicianPerformance({ month: selectedMonth.value })
    performances.value = res.data || []
  } catch (error) {
    console.error('Failed to load performances:', error)
  }
}

async function settleIncome() {
  try {
    const res = await adminFinanceApi.settleIncome()
    ElMessage.success(res.data?.message || '结算完成')
    loadReport()
  } catch (error) {
    console.error('Failed to settle:', error)
  }
}
</script>

<style scoped>
.admin-reports-page {
  padding: 0;
}

.header-actions {
  display: flex;
  gap: 10px;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  text-align: center;
}

.stat-label {
  color: #909399;
  font-size: 14px;
  margin-bottom: 10px;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #409eff;
}

.mt-20 {
  margin-top: 20px;
}

.card-header {
  font-weight: 600;
}
</style>
