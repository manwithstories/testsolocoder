<template>
  <div class="page-container">
    <div class="page-header">
      <h2>能耗监控</h2>
      <div>
        <el-select v-model="period" style="width: 140px; margin-right: 12px;" @change="loadStatistics">
          <el-option label="今日" value="day" />
          <el-option label="本周" value="week" />
          <el-option label="本月" value="month" />
        </el-select>
        <el-select v-model="exportFormat" style="width: 120px; margin-right: 12px;">
          <el-option label="Excel" value="excel" />
          <el-option label="PDF" value="pdf" />
        </el-select>
        <el-button type="primary" :icon="Download" @click="exportReport">导出报表</el-button>
      </div>
    </div>

    <el-row :gutter="20" class="mb-20">
      <el-col :span="6">
        <div class="stat-card">
          <div class="stat-value">{{ statistics.totalEnergy?.toFixed(2) || '0.00' }} kWh</div>
          <div class="stat-label">总能耗</div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card">
          <div class="stat-value">{{ statistics.peakHour || 0 }}:00</div>
          <div class="stat-label">用电高峰</div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card">
          <div class="stat-value">{{ statistics.avgPower?.toFixed(1) || '0.0' }} W</div>
          <div class="stat-label">平均功率</div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card">
          <div class="stat-value">{{ deviceCount }}</div>
          <div class="stat-label">统计设备</div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20">
      <el-col :span="16">
        <div class="chart-container">
          <h3 class="mb-20">能耗分布 - 按设备</h3>
          <div ref="deviceChartRef" style="height: 350px;"></div>
        </div>
      </el-col>
      <el-col :span="8">
        <div class="chart-container">
          <h3 class="mb-20">每小时能耗</h3>
          <div ref="hourChartRef" style="height: 350px;"></div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
      <el-col :span="12">
        <div class="chart-container">
          <h3 class="mb-20">能耗明细</h3>
          <el-table :data="statistics.byDevice || []" style="width: 100%">
            <el-table-column prop="deviceName" label="设备名称" />
            <el-table-column prop="totalEnergy" label="能耗(kWh)" width="120">
              <template #default="{ row }">{{ row.totalEnergy?.toFixed(2) }}</template>
            </el-table-column>
            <el-table-column prop="percentage" label="占比" width="120">
              <template #default="{ row }">
                <el-progress :percentage="Math.round(row.percentage)" :stroke-width="10" />
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-col>
      <el-col :span="12">
        <div class="chart-container">
          <div class="flex-between mb-20">
            <h3>能耗告警</h3>
            <el-tag type="danger" size="small">{{ alerts.length }} 条未处理</el-tag>
          </div>
          <el-empty v-if="alerts.length === 0" description="暂无告警" />
          <el-table v-else :data="alerts" style="width: 100%">
            <el-table-column prop="alertType" label="类型" width="100">
              <template #default="{ row }">
                <el-tag :type="row.level === 'warning' ? 'warning' : 'danger'" size="small">
                  {{ row.alertType === 'high_consumption' ? '高能耗' : row.alertType }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="message" label="消息" show-overflow-tooltip />
            <el-table-column prop="createdAt" label="时间" width="160">
              <template #default="{ row }">{{ formatTime(row.createdAt) }}</template>
            </el-table-column>
          </el-table>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import { Download } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import dayjs from 'dayjs'
import {
  getEnergyStatistics, getEnergyTrend, listEnergyAlerts, exportEnergyReport
} from '@/api/energy'

const period = ref('day')
const exportFormat = ref('excel')
const statistics = ref<any>({})
const alerts = ref<any[]>([])
const deviceCount = ref(0)

const deviceChartRef = ref<HTMLElement>()
const hourChartRef = ref<HTMLElement>()
let deviceChart: echarts.ECharts | null = null
let hourChart: echarts.ECharts | null = null

onMounted(() => {
  deviceChart = echarts.init(deviceChartRef.value!)
  hourChart = echarts.init(hourChartRef.value!)
  loadStatistics()
  loadAlerts()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  deviceChart?.dispose()
  hourChart?.dispose()
})

function handleResize() {
  deviceChart?.resize()
  hourChart?.resize()
}

async function loadStatistics() {
  try {
    const res = await getEnergyStatistics({ period: period.value })
    statistics.value = res
    deviceCount.value = (res.byDevice || []).length

    const byDevice = res.byDevice || []
    deviceChart?.setOption({
      tooltip: { trigger: 'item', formatter: '{b}: {c} kWh ({d}%)' },
      series: [{
        type: 'pie',
        radius: ['40%', '70%'],
        data: byDevice.map((d: any) => ({
          name: d.deviceName,
          value: Number(d.totalEnergy?.toFixed(2))
        }))
      }]
    })

    const byHour = res.byHour || []
    hourChart?.setOption({
      tooltip: { trigger: 'axis' },
      grid: { left: 40, right: 10, top: 20, bottom: 30 },
      xAxis: {
        type: 'category',
        data: byHour.map((h: any) => `${h.hour}:00`)
      },
      yAxis: { type: 'value' },
      series: [{
        type: 'bar',
        data: byHour.map((h: any) => Number(h.energy?.toFixed(2))),
        itemStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: '#83bff6' },
            { offset: 1, color: '#409eff' }
          ])
        }
      }]
    })
  } catch (e) {
    console.error(e)
  }
}

async function loadAlerts() {
  try {
    const res = await listEnergyAlerts({ resolved: 'false' })
    alerts.value = res
  } catch (e) {
    console.error(e)
  }
}

async function exportReport() {
  try {
    const format = exportFormat.value
    const blob = await exportEnergyReport({ format, period: period.value })
    const url = window.URL.createObjectURL(new Blob([blob]))
    const link = document.createElement('a')
    const ext = format === 'pdf' ? 'pdf' : 'xlsx'
    const mime = format === 'pdf' ? 'application/pdf' : 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    link.href = url
    link.download = `energy_report_${period.value}_${dayjs().format('YYYYMMDD')}.${ext}`
    link.click()
    window.URL.revokeObjectURL(url)
    ElMessage.success('报表导出成功')
  } catch (e) {
    console.error(e)
  }
}

function formatTime(t: string) {
  return dayjs(t).format('YYYY-MM-DD HH:mm')
}
</script>
