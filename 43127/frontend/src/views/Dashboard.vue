<template>
  <div class="dashboard">
    <el-row :gutter="20" class="stat-row">
      <el-col :span="6">
        <div class="stat-card">
          <div class="label">总房源数</div>
          <div class="value primary">{{ overview.totalProperties }}</div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card">
          <div class="label">出租率</div>
          <div class="value success">{{ overview.occupancyRate }}%</div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card">
          <div class="label">本月收入</div>
          <div class="value">¥{{ overview.totalIncome }}</div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card">
          <div class="label">待处理工单</div>
          <div class="value warning">{{ overview.pendingRepairs }}</div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="chart-row">
      <el-col :span="12">
        <div class="card">
          <div class="card-header">
            <h3>出租率趋势</h3>
          </div>
          <div ref="occupancyChart" class="chart-container"></div>
        </div>
      </el-col>
      <el-col :span="12">
        <div class="card">
          <div class="card-header">
            <h3>收入趋势</h3>
          </div>
          <div ref="incomeChart" class="chart-container"></div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="chart-row">
      <el-col :span="12">
        <div class="card">
          <div class="card-header">
            <h3>维修工单统计</h3>
          </div>
          <div ref="repairChart" class="chart-container"></div>
        </div>
      </el-col>
      <el-col :span="12">
        <div class="card">
          <div class="card-header">
            <h3>即将到期合同</h3>
          </div>
          <el-table :data="expiringContracts" style="width: 100%">
            <el-table-column prop="contract" label="合同ID" width="100" />
          </el-table>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, shallowRef, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import { getOverview, getOccupancyTrend, getIncomeTrend, getRepairStats } from '@/api/business'
import { getExpiringContracts } from '@/api/tenant'

const overview = ref({
  totalProperties: 0,
  rentedProperties: 0,
  occupancyRate: '0',
  totalIncome: 0,
  pendingRepairs: 0,
  activeContracts: 0
})

const occupancyChart = shallowRef<HTMLElement>()
const incomeChart = shallowRef<HTMLElement>()
const repairChart = shallowRef<HTMLElement>()

let occupancyChartInstance: any = null
let incomeChartInstance: any = null
let repairChartInstance: any = null

const expiringContracts = ref<{ contract: string }[]>([])

onMounted(async () => {
  await loadData()
})

onUnmounted(() => {
  occupancyChartInstance?.dispose()
  incomeChartInstance?.dispose()
  repairChartInstance?.dispose()
})

async function loadData() {
  try {
    const [overviewRes, occupancyRes, incomeRes, repairRes, expiringRes] = await Promise.all([
      getOverview(),
      getOccupancyTrend(),
      getIncomeTrend(),
      getRepairStats(),
      getExpiringContracts()
    ])

    overview.value = overviewRes.data

    expiringContracts.value = expiringRes.data.expiring.map((item: string) => ({ contract: item }))

    initOccupancyChart(occupancyRes.data.months, occupancyRes.data.rates)
    initIncomeChart(incomeRes.data.months, incomeRes.data.incomes)
    initRepairChart(repairRes.data.byCategory, repairRes.data.byStatus)
  } catch (error) {
    console.error('Failed to load dashboard data:', error)
  }
}

function initOccupancyChart(months: string[], rates: number[]) {
  if (!occupancyChart.value) return
  occupancyChartInstance = echarts.init(occupancyChart.value)
  occupancyChartInstance.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: months },
    yAxis: { type: 'value', max: 100 },
    series: [{
      data: rates,
      type: 'line',
      smooth: true,
      areaStyle: {},
      itemStyle: { color: '#67c23a' }
    }]
  })
}

function initIncomeChart(months: string[], incomes: number[]) {
  if (!incomeChart.value) return
  incomeChartInstance = echarts.init(incomeChart.value)
  incomeChartInstance.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: months },
    yAxis: { type: 'value' },
    series: [{
      data: incomes,
      type: 'bar',
      itemStyle: { color: '#409eff' }
    }]
  })
}

function initRepairChart(byCategory: Record<string, number>, byStatus: Record<string, number>) {
  if (!repairChart.value) return
  repairChartInstance = echarts.init(repairChart.value)
  
  const categories = Object.keys(byCategory)
  const values = Object.values(byCategory)

  repairChartInstance.setOption({
    tooltip: { trigger: 'item' },
    legend: { orient: 'vertical', left: 'left' },
    series: [{
      type: 'pie',
      radius: '60%',
      data: categories.map((name, i) => ({ name, value: values[i] }))
    }]
  })
}
</script>

<style scoped>
.dashboard {
  padding: 0;
}

.stat-row {
  margin-bottom: 20px;
}

.stat-card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

.stat-card .label {
  color: #909399;
  font-size: 14px;
  margin-bottom: 8px;
}

.stat-card .value {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
}

.stat-card .value.primary {
  color: #409eff;
}

.stat-card .value.success {
  color: #67c23a;
}

.stat-card .value.warning {
  color: #e6a23c;
}

.chart-row {
  margin-bottom: 20px;
}

.card {
  background: #fff;
  border-radius: 8px;
  padding: 20px;
  box-shadow: 0 2px 12px 0 rgba(0, 0, 0, 0.05);
}

.card-header {
  margin-bottom: 15px;
}

.card-header h3 {
  font-size: 16px;
  font-weight: 600;
}

.chart-container {
  height: 300px;
}
</style>
