<template>
  <div class="page-container">
    <el-row :gutter="20" class="mb-20">
      <el-col :span="6">
        <div class="stat-card">
          <div class="flex-between">
            <div>
              <div class="stat-value">{{ deviceCount }}</div>
              <div class="stat-label">设备总数</div>
            </div>
            <el-icon :size="40" color="#409eff"><Monitor /></el-icon>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card">
          <div class="flex-between">
            <div>
              <div class="stat-value">{{ onlineCount }}</div>
              <div class="stat-label">在线设备</div>
            </div>
            <el-icon :size="40" color="#67c23a"><CircleCheck /></el-icon>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card">
          <div class="flex-between">
            <div>
              <div class="stat-value">{{ totalPower.toFixed(1) }}W</div>
              <div class="stat-label">当前功率</div>
            </div>
            <el-icon :size="40" color="#e6a23c"><Lightning /></el-icon>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card">
          <div class="flex-between">
            <div>
              <div class="stat-value">{{ todayEnergy.toFixed(2) }}kWh</div>
              <div class="stat-label">今日能耗</div>
            </div>
            <el-icon :size="40" color="#f56c6c"><TrendCharts /></el-icon>
          </div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="16">
        <div class="chart-container">
          <div class="flex-between mb-20">
            <h3>能耗趋势</h3>
            <el-radio-group v-model="trendDays" size="small" @change="loadTrend">
              <el-radio-button :value="7">近7天</el-radio-button>
              <el-radio-button :value="30">近30天</el-radio-button>
            </el-radio-group>
          </div>
          <div ref="trendChartRef" style="height: 300px;"></div>
        </div>
      </el-col>
      <el-col :span="8">
        <div class="chart-container">
          <h3 class="mb-20">设备状态分布</h3>
          <div ref="statusChartRef" style="height: 300px;"></div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
      <el-col :span="12">
        <div class="chart-container">
          <div class="flex-between mb-20">
            <h3>高能耗设备 TOP5</h3>
            <el-link type="primary" @click="$router.push('/energy')">查看详情</el-link>
          </div>
          <el-table :data="topDevices" style="width: 100%">
            <el-table-column prop="deviceName" label="设备名称" />
            <el-table-column prop="totalEnergy" label="能耗(kWh)" width="120">
              <template #default="{ row }">{{ row.totalEnergy.toFixed(2) }}</template>
            </el-table-column>
            <el-table-column prop="percentage" label="占比" width="120">
              <template #default="{ row }">{{ row.percentage.toFixed(1) }}%</template>
            </el-table-column>
          </el-table>
        </div>
      </el-col>
      <el-col :span="12">
        <div class="chart-container">
          <div class="flex-between mb-20">
            <h3>最近告警</h3>
            <el-link type="primary" @click="$router.push('/notifications')">查看全部</el-link>
          </div>
          <el-empty v-if="recentAlerts.length === 0" description="暂无告警" :image-size="60" />
          <el-timeline v-else>
            <el-timeline-item
              v-for="alert in recentAlerts"
              :key="alert.id"
              :type="alert.level === 'warning' ? 'warning' : 'danger'"
              :timestamp="alert.createdAt"
            >
              <el-tag :type="alert.level === 'warning' ? 'warning' : 'danger'" size="small">
                {{ alert.alertType === 'high_consumption' ? '高能耗' : alert.alertType }}
              </el-tag>
              <span style="margin-left: 8px;">{{ alert.message }}</span>
            </el-timeline-item>
          </el-timeline>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import { getRealtimeEnergy, getEnergyTrend, listEnergyAlerts } from '@/api/energy'
import { listDevices } from '@/api/device'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const deviceCount = ref(0)
const onlineCount = ref(0)
const totalPower = ref(0)
const todayEnergy = ref(0)
const trendDays = ref(7)
const topDevices = ref<any[]>([])
const recentAlerts = ref<any[]>([])

const trendChartRef = ref<HTMLElement>()
const statusChartRef = ref<HTMLElement>()
let trendChart: echarts.ECharts | null = null
let statusChart: echarts.ECharts | null = null

onMounted(async () => {
  await loadDashboardData()
  initTrendChart()
  initStatusChart()
  loadTrend()

  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  trendChart?.dispose()
  statusChart?.dispose()
})

function handleResize() {
  trendChart?.resize()
  statusChart?.resize()
}

async function loadDashboardData() {
  try {
    const [realtimeRes, alertsRes, devicesRes] = await Promise.all([
      getRealtimeEnergy(),
      listEnergyAlerts({ resolved: 'false' }),
      listDevices()
    ])

    deviceCount.value = realtimeRes.deviceCount || devicesRes.length
    totalPower.value = realtimeRes.totalPower || 0
    onlineCount.value = devicesRes.filter((d: any) => d.status === 'online' || d.status === 'on').length

    if (devicesRes.length > 0) {
      const stats: Record<string, number> = { online: 0, offline: 0, on: 0, off: 0 }
      devicesRes.forEach((d: any) => {
        stats[d.status] = (stats[d.status] || 0) + 1
      })
      updateStatusChart(stats)
    }

    recentAlerts.value = alertsRes.slice(0, 5)

    if (realtimeRes.devices) {
      const sorted = [...realtimeRes.devices].sort((a: any, b: any) => b.hourlyEnergy - a.hourlyEnergy)
      const total = sorted.reduce((sum: number, d: any) => sum + d.hourlyEnergy, 0)
      topDevices.value = sorted.slice(0, 5).map((d: any) => ({
        deviceName: d.deviceName,
        totalEnergy: d.hourlyEnergy,
        percentage: total > 0 ? (d.hourlyEnergy / total) * 100 : 0
      }))
    }
  } catch (e) {
    console.error(e)
  }
}

function initTrendChart() {
  if (!trendChartRef.value) return
  trendChart = echarts.init(trendChartRef.value)
}

function initStatusChart() {
  if (!statusChartRef.value) return
  statusChart = echarts.init(statusChartRef.value)
}

async function loadTrend() {
  try {
    const res = await getEnergyTrend({ days: trendDays.value })
    trendChart?.setOption({
      tooltip: { trigger: 'axis' },
      grid: { left: 40, right: 20, top: 20, bottom: 40 },
      xAxis: {
        type: 'category',
        data: res.trend?.map((t: any) => t.date) || []
      },
      yAxis: { type: 'value', name: 'kWh' },
      series: [{
        data: res.trend?.map((t: any) => t.energy) || [],
        type: 'line',
        smooth: true,
        areaStyle: { opacity: 0.3 },
        lineStyle: { color: '#409eff', width: 2 },
        itemStyle: { color: '#409eff' }
      }]
    })
    todayEnergy.value = res.total || 0
  } catch (e) {
    console.error(e)
  }
}

function updateStatusChart(stats: Record<string, number>) {
  if (!statusChart) return
  statusChart.setOption({
    tooltip: { trigger: 'item' },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      avoidLabelOverlap: false,
      label: { show: true, formatter: '{b}: {c}' },
      data: [
        { value: stats.online || 0, name: '在线', itemStyle: { color: '#67c23a' } },
        { value: stats.offline || 0, name: '离线', itemStyle: { color: '#909399' } },
        { value: stats.on || 0, name: '开启', itemStyle: { color: '#409eff' } },
        { value: stats.off || 0, name: '关闭', itemStyle: { color: '#e6a23c' } }
      ]
    }]
  })
}
</script>
