<template>
  <div>
    <el-row :gutter="20" style="margin-bottom: 20px;">
      <el-col :span="6">
        <div class="stat-card">
          <div class="stat-value">¥{{ stats?.monthly_revenue?.toFixed(2) || '0.00' }}</div>
          <div class="stat-label">本月营收</div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card">
          <div class="stat-value">{{ stats?.monthly_orders || 0 }}</div>
          <div class="stat-label">本月订单</div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card">
          <div class="stat-value">{{ stats?.car_utilization?.toFixed(1) || 0 }}%</div>
          <div class="stat-label">车辆利用率</div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card">
          <div class="stat-value">{{ stats?.total_cars || 0 }}</div>
          <div class="stat-label">车辆总数</div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="16">
        <div class="chart-container">
          <h3 style="margin-bottom: 20px;">营收趋势（近30天）</h3>
          <div ref="revenueChartRef" style="height: 300px;"></div>
        </div>
      </el-col>
      <el-col :span="8">
        <div class="chart-container">
          <h3 style="margin-bottom: 20px;">车辆状态分布</h3>
          <div ref="statusChartRef" style="height: 300px;"></div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px;">
      <el-col :span="24">
        <div class="chart-container">
          <h3 style="margin-bottom: 20px;">热门车型排行</h3>
          <el-table :data="stats?.popular_cars || []" style="width: 100%">
            <el-table-column type="index" label="排名" width="80" />
            <el-table-column label="车型">
              <template #default="{ row }">
                {{ row.brand }} {{ row.model }}
              </template>
            </el-table-column>
            <el-table-column label="评分" width="150">
              <template #default="{ row }">
                <el-rate :model-value="row.rating" disabled size="small" />
              </template>
            </el-table-column>
            <el-table-column prop="review_count" label="评价数" width="100" />
            <el-table-column prop="booking_count" label="预订数" width="100" />
          </el-table>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import { statsApi } from '@/api'
import type { DashboardStats } from '@/types'

const stats = ref<DashboardStats | null>(null)
const revenueChartRef = ref<HTMLElement>()
const statusChartRef = ref<HTMLElement>()
let revenueChart: echarts.ECharts | null = null
let statusChart: echarts.ECharts | null = null

onMounted(async () => {
  await loadStats()
  initCharts()
})

onUnmounted(() => {
  revenueChart?.dispose()
  statusChart?.dispose()
})

const loadStats = async () => {
  try {
    const res = await statsApi.getDashboardStats()
    stats.value = res.data
  } catch {
    // ignore
  }
}

const initCharts = () => {
  if (revenueChartRef.value && stats.value) {
    revenueChart = echarts.init(revenueChartRef.value)
    revenueChart.setOption({
      tooltip: {
        trigger: 'axis'
      },
      grid: {
        left: '3%',
        right: '4%',
        bottom: '3%',
        containLabel: true
      },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: stats.value.revenue_trend.map(item => item.date)
      },
      yAxis: {
        type: 'value'
      },
      series: [{
        name: '营收',
        type: 'line',
        data: stats.value.revenue_trend.map(item => item.revenue),
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
        itemStyle: {
          color: '#409eff'
        }
      }]
    })
  }

  if (statusChartRef.value && stats.value) {
    statusChart = echarts.init(statusChartRef.value)
    const statusMap: Record<string, string> = {
      available: '可用',
      rented: '出租中',
      maintenance: '维护中',
      disabled: '已停用'
    }
    statusChart.setOption({
      tooltip: {
        trigger: 'item'
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
          show: false
        },
        data: stats.value.status_breakdown.map(item => ({
          value: item.count,
          name: statusMap[item.status] || item.status
        }))
      }]
    })
  }
}
</script>
