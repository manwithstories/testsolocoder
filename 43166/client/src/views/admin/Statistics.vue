<template>
  <div class="statistics">
    <div class="page-header">
      <h2 class="page-title">统计分析</h2>
    </div>

    <div class="filter-bar">
      <el-date-picker
        v-model="dateRange"
        type="daterange"
        range-separator="至"
        start-placeholder="开始日期"
        end-placeholder="结束日期"
        value-format="YYYY-MM-DD"
        @change="fetchData"
      />
    </div>

    <el-row :gutter="20" class="mb-24">
      <el-col :span="4">
        <div class="stat-card primary">
          <div class="stat-label">申请总量</div>
          <div class="stat-value">{{ overviewStats?.totalApplications || 0 }}</div>
        </div>
      </el-col>
      <el-col :span="4">
        <div class="stat-card warning">
          <div class="stat-label">待处理</div>
          <div class="stat-value">{{ overviewStats?.pendingApplications || 0 }}</div>
        </div>
      </el-col>
      <el-col :span="4">
        <div class="stat-card info">
          <div class="stat-label">处理中</div>
          <div class="stat-value">{{ overviewStats?.processingApps || 0 }}</div>
        </div>
      </el-col>
      <el-col :span="4">
        <div class="stat-card success">
          <div class="stat-label">已完成</div>
          <div class="stat-value">{{ overviewStats?.completedApps || 0 }}</div>
        </div>
      </el-col>
      <el-col :span="4">
        <div class="stat-card danger">
          <div class="stat-label">已驳回</div>
          <div class="stat-value">{{ overviewStats?.rejectedApps || 0 }}</div>
        </div>
      </el-col>
      <el-col :span="4">
        <div class="stat-card primary">
          <div class="stat-label">总收入</div>
          <div class="stat-value">{{ formatMoney(overviewStats?.totalRevenue || 0) }}</div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="12">
        <div class="card">
          <h3 class="font-weight-600 mb-16">申请状态分布</h3>
          <div ref="statusChartRef" style="height: 300px"></div>
        </div>
      </el-col>
      <el-col :span="12">
        <div class="card">
          <h3 class="font-weight-600 mb-16">公司类型分布</h3>
          <div ref="companyTypeChartRef" style="height: 300px"></div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-24">
      <el-col :span="24">
        <div class="card">
          <h3 class="font-weight-600 mb-16">申请趋势</h3>
          <div ref="trendChartRef" style="height: 350px"></div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-24">
      <el-col :span="24">
        <div class="card">
          <h3 class="font-weight-600 mb-16">代办专员绩效</h3>
          <el-table :data="agentPerformance" style="width: 100%">
            <el-table-column prop="agentName" label="姓名" />
            <el-table-column prop="employeeNo" label="工号" />
            <el-table-column prop="totalHandled" label="处理总数" width="120" />
            <el-table-column prop="completedCount" label="完成数" width="120" />
            <el-table-column prop="inProgressCount" label="处理中" width="120" />
            <el-table-column prop="totalRevenue" label="创收金额" width="150">
              <template #default="{ row }">
                {{ formatMoney(row.totalRevenue) }}
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import * as echarts from 'echarts'
import type { ECharts } from 'echarts'
import { statisticsApi } from '@/api/statistics'
import type { OverviewStats, StatusDistribution, CompanyTypeDistribution, AgentPerformance, TimeSeriesData } from '@/types'

const dateRange = ref<string[]>([])
const overviewStats = ref<OverviewStats | null>(null)
const statusDistribution = ref<StatusDistribution[]>([])
const companyTypeDistribution = ref<CompanyTypeDistribution[]>([])
const agentPerformance = ref<AgentPerformance[]>([])
const timeSeriesData = ref<TimeSeriesData[]>([])

const statusChartRef = ref<HTMLElement | null>(null)
const companyTypeChartRef = ref<HTMLElement | null>(null)
const trendChartRef = ref<HTMLElement | null>(null)

const statusChartInstance = ref<ECharts | null>(null)
const companyTypeChartInstance = ref<ECharts | null>(null)
const trendChartInstance = ref<ECharts | null>(null)

const fetchData = async () => {
  try {
    const params: { startDate?: string; endDate?: string } = {}
    if (dateRange.value.length === 2) {
      params.startDate = dateRange.value[0]
      params.endDate = dateRange.value[1]
    }

    const [overviewRes, statusRes, companyTypeRes, agentRes, trendRes] = await Promise.all([
      statisticsApi.getOverview(params),
      statisticsApi.getStatusDistribution(),
      statisticsApi.getCompanyTypeDistribution(),
      statisticsApi.getAgentPerformance(params),
      statisticsApi.getApplicationTimeSeries({
        startDate: params.startDate || '2024-01-01',
        endDate: params.endDate || new Date().toISOString().split('T')[0],
        interval: 'day'
      })
    ])

    overviewStats.value = overviewRes || null
    statusDistribution.value = statusRes || []
    companyTypeDistribution.value = companyTypeRes || []
    agentPerformance.value = agentRes || []
    timeSeriesData.value = trendRes || []

    await nextTick()
    renderCharts()
  } catch (error) {
    console.error('获取统计数据失败:', error)
  }
}

const renderCharts = () => {
  if (statusChartRef.value) {
    if (statusChartInstance.value) {
      statusChartInstance.value.dispose()
    }
    statusChartInstance.value = echarts.init(statusChartRef.value)
    statusChartInstance.value.setOption({
      tooltip: { trigger: 'item' },
      series: [{
        type: 'pie',
        radius: '60%',
        data: statusDistribution.value.map(item => ({
          name: item.status,
          value: item.count
        }))
      }]
    })
  }

  if (companyTypeChartRef.value) {
    if (companyTypeChartInstance.value) {
      companyTypeChartInstance.value.dispose()
    }
    companyTypeChartInstance.value = echarts.init(companyTypeChartRef.value)
    companyTypeChartInstance.value.setOption({
      tooltip: { trigger: 'item' },
      series: [{
        type: 'pie',
        radius: '60%',
        data: companyTypeDistribution.value.map(item => ({
          name: item.companyType,
          value: item.count
        }))
      }]
    })
  }

  if (trendChartRef.value) {
    if (trendChartInstance.value) {
      trendChartInstance.value.dispose()
    }
    trendChartInstance.value = echarts.init(trendChartRef.value)
    trendChartInstance.value.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: {
        type: 'category',
        data: timeSeriesData.value.map(item => item.date)
      },
      yAxis: {
        type: 'value'
      },
      series: [{
        type: 'line',
        data: timeSeriesData.value.map(item => item.count),
        smooth: true
      }]
    })
  }
}

const formatMoney = (amount: number) => {
  return `¥${amount.toLocaleString()}`
}

onMounted(fetchData)
</script>
