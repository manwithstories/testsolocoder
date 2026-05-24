<template>
  <div class="statistics-page">
    <el-card shadow="never" class="filter-card">
      <el-form :inline="true" :model="query" @submit.prevent>
        <el-form-item label="时间范围">
          <el-date-picker
            v-model="dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            style="width: 280px"
          />
        </el-form-item>
        <el-form-item label="统计粒度">
          <el-radio-group v-model="query.type">
            <el-radio value="daily">日</el-radio>
            <el-radio value="weekly">周</el-radio>
            <el-radio value="monthly">月</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="fetchData">查询</el-button>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
          <el-button :icon="Download" @click="handleExport('excel')">导出 Excel</el-button>
          <el-button :icon="Printer" @click="handleExport('pdf')">导出 PDF</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-row :gutter="16">
      <el-col v-for="item in overviewCards" :key="item.key" :xs="12" :sm="12" :md="8" :lg="6">
        <el-card shadow="hover" class="overview-card">
          <div class="overview-label">{{ item.label }}</div>
          <div class="overview-value" :style="{ color: item.color }">{{ item.value }}</div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :xs="24" :md="16">
        <el-card shadow="never">
          <template #header>
            <span>销售趋势</span>
          </template>
          <div ref="trendChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
      <el-col :xs="24" :md="8">
        <el-card shadow="never">
          <template #header>
            <span>房型分布</span>
          </template>
          <div ref="houseTypeChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :xs="24" :md="12">
        <el-card shadow="never">
          <template #header>
            <span>面积分布</span>
          </template>
          <div ref="areaChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
      <el-col :xs="24" :md="12">
        <el-card shadow="never">
          <template #header>
            <span>预算分布</span>
          </template>
          <div ref="budgetChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, onBeforeUnmount } from 'vue'
import * as echarts from 'echarts'
import { Search, Refresh, Download, Printer } from '@element-plus/icons-vue'
import { ElMessage } from 'element-plus'
import { saveAs } from 'file-saver'
import {
  getOverview,
  getSalesTrend,
  getCustomerProfile,
  exportExcel,
  exportPdf,
  type StatisticsOverview,
  type SalesTrendItem,
  type StatisticsQuery
} from '@/api/statistics'

const query = reactive<StatisticsQuery>({
  type: 'monthly'
})
const dateRange = ref<string[]>([])

const overview = ref<StatisticsOverview>({
  totalOrders: 0,
  totalAmount: 0,
  totalDesigns: 0,
  totalTickets: 0,
  avgOrderAmount: 0,
  avgReviewScore: 0
})

const salesTrend = ref<SalesTrendItem[]>([])
const customerProfile = ref({
  areaDistribution: [] as { label: string; value: number }[],
  houseTypeDistribution: [] as { label: string; value: number }[],
  budgetDistribution: [] as { label: string; value: number }[]
})

const overviewCards = ref([
  { key: 'totalOrders', label: '订单总数', value: 0, color: '#409eff' },
  { key: 'totalAmount', label: '总销售额', value: '¥0', color: '#67c23a' },
  { key: 'totalDesigns', label: '设计方案数', value: 0, color: '#e6a23c' },
  { key: 'totalTickets', label: '工单总数', value: 0, color: '#f56c6c' },
  { key: 'avgOrderAmount', label: '客单价', value: '¥0', color: '#909399' },
  { key: 'avgReviewScore', label: '平均评分', value: 0, color: '#9b59b6' }
])

function updateOverviewCards() {
  overviewCards.value = [
    { key: 'totalOrders', label: '订单总数', value: overview.value.totalOrders, color: '#409eff' },
    {
      key: 'totalAmount',
      label: '总销售额',
      value: `¥${overview.value.totalAmount.toLocaleString()}`,
      color: '#67c23a'
    },
    { key: 'totalDesigns', label: '设计方案数', value: overview.value.totalDesigns, color: '#e6a23c' },
    { key: 'totalTickets', label: '工单总数', value: overview.value.totalTickets, color: '#f56c6c' },
    {
      key: 'avgOrderAmount',
      label: '客单价',
      value: `¥${overview.value.avgOrderAmount.toLocaleString()}`,
      color: '#909399'
    },
    { key: 'avgReviewScore', label: '平均评分', value: overview.value.avgReviewScore, color: '#9b59b6' }
  ]
}

const trendChartRef = ref<HTMLDivElement>()
const houseTypeChartRef = ref<HTMLDivElement>()
const areaChartRef = ref<HTMLDivElement>()
const budgetChartRef = ref<HTMLDivElement>()

let trendChart: echarts.ECharts | null = null
let houseTypeChart: echarts.ECharts | null = null
let areaChart: echarts.ECharts | null = null
let budgetChart: echarts.ECharts | null = null

function renderTrendChart() {
  if (!trendChartRef.value) return
  if (!trendChart) trendChart = echarts.init(trendChartRef.value)
  trendChart.setOption({
    tooltip: { trigger: 'axis' },
    legend: { data: ['销售额', '订单数'] },
    grid: { left: 40, right: 40, top: 40, bottom: 40 },
    xAxis: { type: 'category', data: salesTrend.value.map((i) => i.date) },
    yAxis: [{ type: 'value', name: '销售额' }, { type: 'value', name: '订单数' }],
    series: [
      {
        name: '销售额',
        type: 'line',
        smooth: true,
        data: salesTrend.value.map((i) => i.amount),
        itemStyle: { color: '#409eff' },
        areaStyle: { opacity: 0.1 }
      },
      {
        name: '订单数',
        type: 'bar',
        yAxisIndex: 1,
        data: salesTrend.value.map((i) => i.count),
        itemStyle: { color: '#67c23a' }
      }
    ]
  })
}

function renderPieChart(
  chartRef: HTMLDivElement | undefined,
  chartInstance: echarts.ECharts | null,
  data: { label: string; value: number }[]
) {
  if (!chartRef) return chartInstance
  const instance = chartInstance ?? echarts.init(chartRef)
  instance.setOption({
    tooltip: { trigger: 'item' },
    legend: { orient: 'vertical', left: 'left' },
    series: [
      {
        type: 'pie',
        radius: ['40%', '70%'],
        data: data.map((i) => ({ name: i.label, value: i.value })),
        label: { formatter: '{b}: {d}%' }
      }
    ]
  })
  return instance
}

function renderAllCharts() {
  renderTrendChart()
  houseTypeChart = renderPieChart(
    houseTypeChartRef.value,
    houseTypeChart,
    customerProfile.value.houseTypeDistribution
  ) as echarts.ECharts
  areaChart = renderPieChart(areaChartRef.value, areaChart, customerProfile.value.areaDistribution) as echarts.ECharts
  budgetChart = renderPieChart(
    budgetChartRef.value,
    budgetChart,
    customerProfile.value.budgetDistribution
  ) as echarts.ECharts
}

function handleResize() {
  trendChart?.resize()
  houseTypeChart?.resize()
  areaChart?.resize()
  budgetChart?.resize()
}

async function fetchData() {
  const params: StatisticsQuery = { ...query }
  if (dateRange.value?.length === 2) {
    params.startDate = dateRange.value[0]
    params.endDate = dateRange.value[1]
  }
  try {
    const [ov, trend, profile] = await Promise.all([
      getOverview(params),
      getSalesTrend(params),
      getCustomerProfile()
    ])
    overview.value = ov as any
    salesTrend.value = (trend as any) ?? []
    customerProfile.value = (profile as any) ?? {
      areaDistribution: [],
      houseTypeDistribution: [],
      budgetDistribution: []
    }
    updateOverviewCards()
    renderAllCharts()
  } catch (e) {
    console.error(e)
  }
}

function handleReset() {
  dateRange.value = []
  query.type = 'monthly'
  fetchData()
}

async function handleExport(type: 'excel' | 'pdf') {
  try {
    const params: StatisticsQuery = { ...query }
    if (dateRange.value?.length === 2) {
      params.startDate = dateRange.value[0]
      params.endDate = dateRange.value[1]
    }
    const blob = type === 'excel' ? await exportExcel(params) : await exportPdf(params)
    const filename = type === 'excel' ? 'statistics.xlsx' : 'statistics.pdf'
    saveAs(blob as any, filename)
    ElMessage.success('导出成功')
  } catch (e) {
    ElMessage.error('导出失败')
  }
}

onMounted(() => {
  fetchData()
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  trendChart?.dispose()
  houseTypeChart?.dispose()
  areaChart?.dispose()
  budgetChart?.dispose()
})
</script>

<style lang="scss" scoped>
.statistics-page {
  .filter-card {
    margin-bottom: 16px;
  }
  .overview-card {
    text-align: center;
    .overview-label {
      color: #909399;
      font-size: 14px;
    }
    .overview-value {
      font-size: 26px;
      font-weight: bold;
      margin-top: 8px;
    }
  }
  .chart-container {
    width: 100%;
    height: 360px;
  }
}
</style>
