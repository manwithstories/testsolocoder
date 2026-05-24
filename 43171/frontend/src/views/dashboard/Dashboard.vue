<template>
  <el-card>
    <template #header>
      <div class="header">
        <span>数据统计</span>
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          @change="fetchData"
          style="width: 280px"
        />
      </div>
    </template>

    <el-row :gutter="20">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>收入趋势</span>
            <el-button type="primary" size="small" @click="exportExcel('revenue')">导出</el-button>
          </template>
          <div ref="revenueChartRef" style="height: 300px"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>区域订单热力图</span>
            <el-button type="primary" size="small" @click="exportExcel('region')">导出</el-button>
          </template>
          <div ref="regionChartRef" style="height: 300px"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="24">
        <el-card>
          <template #header>
            <span>设备利用率</span>
            <el-button type="primary" size="small" @click="exportExcel('drone')">导出</el-button>
          </template>
          <el-table :data="droneStats">
            <el-table-column prop="drone_name" label="设备名称" />
            <el-table-column prop="total_days" label="租赁天数" width="120" />
            <el-table-column label="利用率" width="150">
              <template #default="{ row }">
                <el-progress :percentage="Math.round(row.utilization * 100)" />
              </template>
            </el-table-column>
            <el-table-column prop="income" label="收入" width="150">
              <template #default="{ row }">¥{{ row.income.toFixed(2) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted, watch } from 'vue'
import { ElMessage } from 'element-plus'
import * as echarts from 'echarts'
import request from '@/utils/request'
import dayjs from 'dayjs'

const dateRange = ref<[Date, Date]>([
  dayjs().subtract(30, 'day').toDate(),
  dayjs().toDate()
])

const revenueChartRef = ref<HTMLElement>()
const regionChartRef = ref<HTMLElement>()
let revenueChart: echarts.ECharts | null = null
let regionChart: echarts.ECharts | null = null

const droneStats = ref<DroneStats[]>([])

onMounted(() => {
  fetchData()
})

onUnmounted(() => {
  revenueChart?.dispose()
  regionChart?.dispose()
})

watch(dateRange, () => {
  fetchData()
})

async function fetchData() {
  const startDate = dayjs(dateRange.value[0]).format('YYYY-MM-DD')
  const endDate = dayjs(dateRange.value[1]).format('YYYY-MM-DD')

  try {
    const [revenueRes, regionRes, droneRes]: any[] = await Promise.all([
      request.get('/stats/revenue', { params: { start_date: startDate, end_date: endDate } }),
      request.get('/stats/region', { params: { start_date: startDate, end_date: endDate } }),
      request.get('/stats/drone', { params: { start_date: startDate, end_date: endDate } })
    ])

    initRevenueChart(revenueRes.data)
    initRegionChart(regionRes.data)
    droneStats.value = droneRes.data
  } catch (e) {
    console.error(e)
  }
}

function initRevenueChart(data: RevenueStats[]) {
  if (revenueChartRef.value) {
    revenueChart = echarts.init(revenueChartRef.value)
    revenueChart.setOption({
      tooltip: { trigger: 'axis' },
      xAxis: {
        type: 'category',
        data: data.map(d => d.date)
      },
      yAxis: { type: 'value' },
      series: [{
        name: '收入',
        type: 'line',
        smooth: true,
        data: data.map(d => d.amount),
        areaStyle: {}
      }]
    })
  }
}

function initRegionChart(data: RegionStats[]) {
  if (regionChartRef.value) {
    regionChart = echarts.init(regionChartRef.value)
    regionChart.setOption({
      tooltip: { trigger: 'item' },
      series: [{
        name: '区域订单',
        type: 'pie',
        radius: ['40%', '70%'],
        data: data.map(d => ({ name: d.region, value: d.count }))
      }]
    })
  }
}

async function exportExcel(type: string) {
  const startDate = dayjs(dateRange.value[0]).format('YYYY-MM-DD')
  const endDate = dayjs(dateRange.value[1]).format('YYYY-MM-DD')
  const urls: Record<string, string> = {
    revenue: '/export/revenue',
    region: '/export/region',
    drone: '/export/drone'
  }
  
  try {
    const response = await request.get(urls[type], {
      params: { start_date: startDate, end_date: endDate },
      responseType: 'blob'
    }) as unknown as Blob
    
    const blob = new Blob([response], {
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    const filenameMap: Record<string, string> = {
      revenue: '收入统计',
      region: '区域订单统计',
      drone: '设备利用率统计'
    }
    link.download = `${filenameMap[type]}_${startDate}_${endDate}.xlsx`
    document.body.appendChild(link)
    link.click()
    document.body.removeChild(link)
    URL.revokeObjectURL(url)
    
    ElMessage.success('导出成功')
  } catch (e) {
    console.error(e)
    ElMessage.error('导出失败')
  }
}
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
