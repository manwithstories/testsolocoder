<template>
  <div class="admin-statistics">
    <el-card class="filter-card">
      <el-form :inline="true" :model="filterForm">
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="filterForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchStatistics">查询</el-button>
          <el-button @click="resetFilter">重置</el-button>
          <el-button type="success" @click="handleExportPDF">导出PDF报告</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-row :gutter="24">
      <el-col :span="6">
        <el-card class="stat-card total-sales">
          <div class="stat-icon">💰</div>
          <div class="stat-content">
            <div class="stat-value">¥{{ salesData.totalSales.toFixed(2) }}</div>
            <div class="stat-label">总销售额</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card total-orders">
          <div class="stat-icon">📋</div>
          <div class="stat-content">
            <div class="stat-value">{{ salesData.totalOrders }}</div>
            <div class="stat-label">总订单数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card total-tickets">
          <div class="stat-icon">🎫</div>
          <div class="stat-content">
            <div class="stat-value">{{ salesData.totalTickets }}</div>
            <div class="stat-label">售票总数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card average-price">
          <div class="stat-icon">📊</div>
          <div class="stat-content">
            <div class="stat-value">¥{{ salesData.avgPrice.toFixed(2) }}</div>
            <div class="stat-label">平均票价</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="24" style="margin-top: 24px;">
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <span>日销售趋势</span>
          </template>
          <div ref="dailySalesChart" class="chart-container"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <span>销售区域分布</span>
          </template>
          <div ref="areaSalesChart" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="24" style="margin-top: 24px;">
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <span>观众画像 - 年龄分布</span>
          </template>
          <div ref="ageChart" class="chart-container"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card class="chart-card">
          <template #header>
            <span>观众画像 - 性别分布</span>
          </template>
          <div ref="genderChart" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="24" style="margin-top: 24px;">
      <el-col :span="24">
        <el-card class="chart-card">
          <template #header>
            <div class="card-header">
              <span>热门演出排名</span>
              <el-select v-model="rankingType" style="width: 150px;" @change="fetchShowRanking">
                <el-option label="按销售额" value="sales" />
                <el-option label="按销量" value="tickets" />
              </el-select>
            </div>
          </template>
          <el-table :data="showRanking" style="width: 100%">
            <el-table-column label="排名" width="80">
              <template #default="{ $index }">
                <el-tag :type="$index < 3 ? 'warning' : 'info'">{{ $index + 1 }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="name" label="演出名称" min-width="200" />
            <el-table-column prop="artist" label="艺人" width="150" />
            <el-table-column label="销售额" width="150">
              <template #default="{ row }">¥{{ row.total_sales?.toFixed(2) || '0.00' }}</template>
            </el-table-column>
            <el-table-column label="售票数" width="120">
              <template #default="{ row }">{{ row.ticket_count || 0 }} 张</template>
            </el-table-column>
            <el-table-column label="上座率" width="120">
              <template #default="{ row }">
                <el-progress
                  :percentage="Math.round((row.ticket_count || 0) / (row.total_seats || 1) * 100)"
                  :stroke-width="12"
                />
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick } from 'vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import { statisticsApi } from '@/api'

const dailySalesChart = ref<HTMLElement>()
const areaSalesChart = ref<HTMLElement>()
const ageChart = ref<HTMLElement>()
const genderChart = ref<HTMLElement>()

let dailySalesChartInstance: echarts.ECharts | null = null
let areaSalesChartInstance: echarts.ECharts | null = null
let ageChartInstance: echarts.ECharts | null = null
let genderChartInstance: echarts.ECharts | null = null

const filterForm = reactive({
  dateRange: [] as string[]
})

const salesData = reactive({
  totalSales: 0,
  totalOrders: 0,
  totalTickets: 0,
  avgPrice: 0
})

const rankingType = ref('sales')
const showRanking = ref<any[]>([])

async function fetchStatistics() {
  try {
    const params: any = {}
    if (filterForm.dateRange && filterForm.dateRange.length === 2) {
      params.start_date = filterForm.dateRange[0]
      params.end_date = filterForm.dateRange[1]
    }

    const [salesRes, dailyRes, audienceRes] = await Promise.all([
      statisticsApi.getSales(params),
      statisticsApi.getDailySales(params),
      statisticsApi.getAudienceProfile(params)
    ])

    salesData.totalSales = salesRes?.total_sales || 0
    salesData.totalOrders = salesRes?.total_orders || 0
    salesData.totalTickets = salesRes?.total_tickets || 0
    salesData.avgPrice = salesData.totalTickets > 0 ? salesData.totalSales / salesData.totalTickets : 0

    await nextTick()
    initDailySalesChart(dailyRes || [])
    initAreaSalesChart(salesRes?.area_sales || [])
    initAgeChart(audienceRes?.age_distribution || [])
    initGenderChart(audienceRes?.gender_distribution || [])

    fetchShowRanking()
  } catch (err) {
    console.error(err)
  }
}

async function fetchShowRanking() {
  try {
    const params: any = { type: rankingType.value }
    if (filterForm.dateRange && filterForm.dateRange.length === 2) {
      params.start_date = filterForm.dateRange[0]
      params.end_date = filterForm.dateRange[1]
    }

    const res = await statisticsApi.getSales(params)
    showRanking.value = res?.shows || []
  } catch (err) {
    console.error(err)
  }
}

function resetFilter() {
  filterForm.dateRange = []
  fetchStatistics()
}

async function handleExportPDF() {
  try {
    const params: any = {}
    if (filterForm.dateRange && filterForm.dateRange.length === 2) {
      params.start_date = filterForm.dateRange[0]
      params.end_date = filterForm.dateRange[1]
    }

    const blob = await statisticsApi.exportPDF(params) as Blob
    const url = window.URL.createObjectURL(new Blob([blob]))
    const link = document.createElement('a')
    link.href = url
    link.setAttribute('download', `statistics_report_${new Date().getTime()}.pdf`)
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    window.URL.revokeObjectURL(url)
    ElMessage.success('PDF报告导出成功')
  } catch (err: any) {
    ElMessage.error(err.message || 'PDF导出失败')
  }
}

function initDailySalesChart(data: any[]) {
  if (!dailySalesChart.value) return

  if (!dailySalesChartInstance) {
    dailySalesChartInstance = echarts.init(dailySalesChart.value)
  }

  const option = {
    tooltip: {
      trigger: 'axis',
      formatter: (params: any) => {
        return `${params[0].axisValue}<br/>销售额: ¥${params[0].value.toFixed(2)}`
      }
    },
    xAxis: {
      type: 'category',
      data: data.map((item: any) => item.date),
      axisLabel: { rotate: 45 }
    },
    yAxis: {
      type: 'value',
      axisLabel: {
        formatter: '¥{value}'
      }
    },
    series: [{
      data: data.map((item: any) => item.sales),
      type: 'line',
      smooth: true,
      areaStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: 'rgba(64, 158, 255, 0.3)' },
          { offset: 1, color: 'rgba(64, 158, 255, 0.05)' }
        ])
      },
      lineStyle: {
        color: '#409eff',
        width: 2
      },
      itemStyle: { color: '#409eff' }
    }]
  }

  dailySalesChartInstance.setOption(option)
}

function initAreaSalesChart(data: any[]) {
  if (!areaSalesChart.value) return

  if (!areaSalesChartInstance) {
    areaSalesChartInstance = echarts.init(areaSalesChart.value)
  }

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: ¥{c} ({d}%)'
    },
    legend: {
      orient: 'vertical',
      left: 'left'
    },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      avoidLabelOverlap: false,
      itemStyle: {
        borderRadius: 10,
        borderColor: '#fff',
        borderWidth: 2
      },
      label: {
        show: true,
        formatter: '{b}\n{d}%'
      },
      data: data.map((item: any) => ({
        value: item.sales,
        name: item.area_name
      }))
    }]
  }

  areaSalesChartInstance.setOption(option)
}

function initAgeChart(data: any[]) {
  if (!ageChart.value) return

  if (!ageChartInstance) {
    ageChartInstance = echarts.init(ageChart.value)
  }

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c}人 ({d}%)'
    },
    xAxis: {
      type: 'category',
      data: data.map((item: any) => item.age_group)
    },
    yAxis: {
      type: 'value'
    },
    series: [{
      type: 'bar',
      data: data.map((item: any) => item.count),
      itemStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: '#67c23a' },
          { offset: 1, color: '#95d475' }
        ]),
        borderRadius: [4, 4, 0, 0]
      },
      barWidth: '50%'
    }]
  }

  ageChartInstance.setOption(option)
}

function initGenderChart(data: any[]) {
  if (!genderChart.value) return

  if (!genderChartInstance) {
    genderChartInstance = echarts.init(genderChart.value)
  }

  const option = {
    tooltip: {
      trigger: 'item',
      formatter: '{b}: {c}人 ({d}%)'
    },
    series: [{
      type: 'pie',
      radius: '60%',
      data: data.map((item: any) => ({
        value: item.count,
        name: item.gender === 'male' ? '男' : item.gender === 'female' ? '女' : '未知'
      })),
      color: ['#409eff', '#f56c6c', '#909399'],
      label: {
        formatter: '{b}: {d}%'
      }
    }]
  }

  genderChartInstance.setOption(option)
}

onMounted(() => {
  fetchStatistics()

  window.addEventListener('resize', () => {
    dailySalesChartInstance?.resize()
    areaSalesChartInstance?.resize()
    ageChartInstance?.resize()
    genderChartInstance?.resize()
  })
})
</script>

<style lang="scss" scoped>
.admin-statistics {
  .filter-card {
    margin-bottom: 24px;
  }

  .stat-card {
    display: flex;
    align-items: center;
    gap: 16px;

    .stat-icon {
      font-size: 40px;
      width: 64px;
      height: 64px;
      display: flex;
      align-items: center;
      justify-content: center;
      border-radius: 12px;
    }

    &.total-sales .stat-icon {
      background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    }

    &.total-orders .stat-icon {
      background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
    }

    &.total-tickets .stat-icon {
      background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
    }

    &.average-price .stat-icon {
      background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
    }

    .stat-content {
      .stat-value {
        font-size: 24px;
        font-weight: bold;
        color: #303133;
        margin-bottom: 4px;
      }

      .stat-label {
        color: #909399;
        font-size: 14px;
      }
    }
  }

  .chart-card {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }

    .chart-container {
      height: 350px;
      width: 100%;
    }
  }
}
</style>
