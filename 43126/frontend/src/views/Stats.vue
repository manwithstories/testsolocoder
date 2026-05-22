<template>
  <div class="stats">
    <el-card class="filter-card">
      <el-form :inline="true" :model="filterForm">
        <el-form-item label="开始日期">
          <el-date-picker
            v-model="filterForm.start_date"
            type="date"
            placeholder="选择开始日期"
            value-format="YYYY-MM-DD"
            style="width: 160px"
          />
        </el-form-item>
        <el-form-item label="结束日期">
          <el-date-picker
            v-model="filterForm.end_date"
            type="date"
            placeholder="选择结束日期"
            value-format="YYYY-MM-DD"
            style="width: 160px"
          />
        </el-form-item>
        <el-form-item label="部门">
          <el-input v-model="filterForm.department" placeholder="部门名称（选填）" clearable style="width: 160px" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadStats">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
        </el-form-item>
        <el-form-item>
          <el-button type="success" @click="exportExcel">
            <el-icon><Download /></el-icon>
            导出Excel
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-item">
            <div class="stat-label">总预订次数</div>
            <div class="stat-value">{{ stats?.summary?.total_bookings || 0 }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-item">
            <div class="stat-label">总使用时长(小时)</div>
            <div class="stat-value">{{ stats?.summary?.total_hours?.toFixed(1) || 0 }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-item">
            <div class="stat-label">总收入(元)</div>
            <div class="stat-value">¥{{ stats?.summary?.total_revenue?.toFixed(2) || 0 }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-item">
            <div class="stat-label">整体利用率</div>
            <div class="stat-value">{{ stats?.summary?.util_rate?.toFixed(1) || 0 }}%</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>会议室利用率</span>
          </template>
          <el-table :data="stats?.utilization || []" style="width: 100%">
            <el-table-column prop="room_name" label="会议室" />
            <el-table-column prop="floor" label="楼层" width="80" />
            <el-table-column prop="bookings" label="预订次数" width="100" />
            <el-table-column prop="total_hours" label="使用时长" width="100">
              <template #default="{ row }">
                {{ row.total_hours?.toFixed(1) }}h
              </template>
            </el-table-column>
            <el-table-column label="利用率" width="120">
              <template #default="{ row }">
                <el-progress :percentage="row.util_rate || 0" :stroke-width="8" />
              </template>
            </el-table-column>
            <el-table-column prop="revenue" label="收入" width="120">
              <template #default="{ row }">
                ¥{{ row.revenue?.toFixed(2) }}
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>热门时段</span>
          </template>
          <div ref="hourlyChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row style="margin-top: 20px">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>收入趋势</span>
          </template>
          <div ref="revenueChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick, watch } from 'vue'
import { ElMessage } from 'element-plus'
import { api } from '@/api'
import * as echarts from 'echarts'
import dayjs from 'dayjs'

const stats = ref<any>(null)
const hourlyChartRef = ref<HTMLElement>()
const revenueChartRef = ref<HTMLElement>()
let hourlyChart: echarts.ECharts | null = null
let revenueChart: echarts.ECharts | null = null

const filterForm = reactive({
  start_date: dayjs().subtract(30, 'day').format('YYYY-MM-DD'),
  end_date: dayjs().format('YYYY-MM-DD'),
  department: ''
})

onMounted(() => {
  loadStats()
})

async function loadStats() {
  try {
    const res: any = await api.getStats(filterForm)
    stats.value = res.data
    await nextTick()
    renderCharts()
  } catch (e) {
    console.error(e)
  }
}

function renderCharts() {
  if (hourlyChartRef.value && stats.value?.hourly_popular) {
    if (!hourlyChart) {
      hourlyChart = echarts.init(hourlyChartRef.value)
    }
    const data = stats.value.hourly_popular.sort((a: any, b: any) => a.hour.localeCompare(b.hour))
    hourlyChart.setOption({
      tooltip: { trigger: 'axis' },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: { type: 'category', data: data.map((d: any) => d.hour) },
      yAxis: { type: 'value' },
      series: [{
        type: 'bar',
        data: data.map((d: any) => d.count),
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#667eea' },
            { offset: 1, color: '#764ba2' }
          ])
        },
        barWidth: '60%'
      }]
    })
  }

  if (revenueChartRef.value && stats.value?.revenue_trend) {
    if (!revenueChart) {
      revenueChart = echarts.init(revenueChartRef.value)
    }
    const data = stats.value.revenue_trend.sort((a: any, b: any) => a.date.localeCompare(b.date))
    revenueChart.setOption({
      tooltip: { trigger: 'axis' },
      legend: { data: ['收入', '预订次数'] },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: { type: 'category', data: data.map((d: any) => d.date.slice(5)) },
      yAxis: [
        { type: 'value', name: '收入(元)' },
        { type: 'value', name: '预订次数' }
      ],
      series: [
        {
          name: '收入',
          type: 'line',
          smooth: true,
          data: data.map((d: any) => d.revenue),
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(102, 126, 234, 0.3)' },
              { offset: 1, color: 'rgba(102, 126, 234, 0.05)' }
            ])
          },
          lineStyle: { color: '#667eea' }
        },
        {
          name: '预订次数',
          type: 'bar',
          yAxisIndex: 1,
          data: data.map((d: any) => d.bookings),
          itemStyle: { color: '#67C23A' }
        }
      ]
    })
  }
}

async function exportExcel() {
  try {
    const res: any = await api.exportStats(filterForm)
    const blob = new Blob([res], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const url = window.URL.createObjectURL(blob)
    const a = document.createElement('a')
    a.href = url
    a.download = `统计报表_${filterForm.start_date}_${filterForm.end_date}.xlsx`
    a.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (e) {
    console.error(e)
  }
}
</script>

<style scoped>
.filter-card {
  border-radius: 8px;
}

.stat-card {
  border-radius: 8px;
}

.stat-item {
  text-align: center;
}

.stat-label {
  color: #909399;
  font-size: 14px;
  margin-bottom: 8px;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
}

.chart-container {
  width: 100%;
  height: 300px;
}
</style>
