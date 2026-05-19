<template>
  <div class="stats-page">
    <el-row :gutter="20">
      <el-col :span="24">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>年度阅读趋势</span>
              <el-select v-model="trendYear" style="width: 120px" @change="loadYearlyTrend">
                <el-option v-for="y in years" :key="y" :label="y" :value="y" />
              </el-select>
            </div>
          </template>
          <div ref="trendChartRef" class="chart" style="height: 350px"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>阅读热力图</span>
              <el-select v-model="heatmapYear" style="width: 120px" @change="loadHeatmap">
                <el-option v-for="y in years" :key="y" :label="y" :value="y" />
              </el-select>
            </div>
          </template>
          <div ref="heatmapChartRef" class="chart" style="height: 200px"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <span>阅读时长分布</span>
          </template>
          <div ref="durationChartRef" class="chart" style="height: 200px"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <span>分类统计</span>
          </template>
          <div ref="categoryChartRef" class="chart" style="height: 250px"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="never">
          <template #header>
            <span>标签统计</span>
          </template>
          <div ref="tagChartRef" class="chart" style="height: 250px"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import * as echarts from 'echarts'
import { getYearlyTrend, getReadingHeatmap, getDurationDistribution, getCategoryStats, getTagStats } from '@/api/stats'
import type { MonthlyStats, HeatmapData, DurationStats, CategoryStats, TagStats } from '@/types'

const trendYear = ref(new Date().getFullYear())
const heatmapYear = ref(new Date().getFullYear())
const years = [new Date().getFullYear(), new Date().getFullYear() - 1, new Date().getFullYear() - 2]

const trendChartRef = ref<HTMLElement>()
const heatmapChartRef = ref<HTMLElement>()
const durationChartRef = ref<HTMLElement>()
const categoryChartRef = ref<HTMLElement>()
const tagChartRef = ref<HTMLElement>()

let trendChart: echarts.ECharts | null = null
let heatmapChart: echarts.ECharts | null = null
let durationChart: echarts.ECharts | null = null
let categoryChart: echarts.ECharts | null = null
let tagChart: echarts.ECharts | null = null

const loadYearlyTrend = async () => {
  try {
    const data = await getYearlyTrend(trendYear.value)
    renderTrendChart(data.monthly)
  } catch (e) {}
}

const loadHeatmap = async () => {
  try {
    const data = await getReadingHeatmap(heatmapYear.value)
    renderHeatmapChart(data.data)
  } catch (e) {}
}

const loadDuration = async () => {
  try {
    const data = await getDurationDistribution()
    renderDurationChart(data)
  } catch (e) {}
}

const loadCategoryStats = async () => {
  try {
    const data = await getCategoryStats()
    renderCategoryChart(data)
  } catch (e) {}
}

const loadTagStats = async () => {
  try {
    const data = await getTagStats()
    renderTagChart(data)
  } catch (e) {}
}

const renderTrendChart = (data: MonthlyStats[]) => {
  if (!trendChartRef.value) return
  if (!trendChart) trendChart = echarts.init(trendChartRef.value)
  trendChart.setOption({
    tooltip: { trigger: 'axis' },
    legend: { data: ['完成数量', '阅读页数'] },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: data.map(d => d.month) },
    yAxis: [
      { type: 'value', name: '数量' },
      { type: 'value', name: '页数' }
    ],
    series: [
      { name: '完成数量', type: 'bar', data: data.map(d => d.completed), itemStyle: { color: '#409eff' } },
      { name: '阅读页数', type: 'line', yAxisIndex: 1, data: data.map(d => d.pages_read), itemStyle: { color: '#67c23a' }, smooth: true }
    ]
  })
}

const renderHeatmapChart = (data: HeatmapData[]) => {
  if (!heatmapChartRef.value) return
  if (!heatmapChart) heatmapChart = echarts.init(heatmapChartRef.value)
  const dateMap: Record<string, number> = {}
  data.forEach(d => { dateMap[d.date] = d.count })
  const calendarData: [string, number][] = []
  const startDate = new Date(heatmapYear.value, 0, 1)
  const endDate = new Date(heatmapYear.value, 11, 31)
  for (let d = new Date(startDate); d <= endDate; d.setDate(d.getDate() + 1)) {
    const dateStr = d.toISOString().split('T')[0]
    calendarData.push([dateStr, dateMap[dateStr] || 0])
  }
  heatmapChart.setOption({
    tooltip: { position: 'top' },
    visualMap: {
      min: 0, max: Math.max(...data.map(d => d.count), 1),
      calculable: true, orient: 'horizontal', left: 'center', bottom: '0',
      inRange: { color: ['#ebedf0', '#c6e48b', '#7bc96f', '#239a3b', '#196127'] }
    },
    calendar: {
      top: 20, left: 50, right: 50, cellSize: ['auto', 12],
      range: heatmapYear.value, itemStyle: { borderWidth: 0.5, borderColor: '#fff' },
      yearLabel: { show: false }, monthLabel: { nameMap: 'ZH' },
      dayLabel: { firstDay: 1, nameMap: ['日', '一', '二', '三', '四', '五', '六'] }
    },
    series: [{ type: 'heatmap', coordinateSystem: 'calendar', data: calendarData }]
  })
}

const renderDurationChart = (data: DurationStats[]) => {
  if (!durationChartRef.value) return
  if (!durationChart) durationChart = echarts.init(durationChartRef.value)
  durationChart.setOption({
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'category', data: data.map(d => d.range) },
    yAxis: { type: 'value' },
    series: [{
      type: 'bar', data: data.map(d => d.count),
      itemStyle: {
        color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
          { offset: 0, color: '#667eea' },
          { offset: 1, color: '#764ba2' }
        ])
      }
    }]
  })
}

const renderCategoryChart = (data: CategoryStats[]) => {
  if (!categoryChartRef.value) return
  if (!categoryChart) categoryChart = echarts.init(categoryChartRef.value)
  categoryChart.setOption({
    tooltip: { trigger: 'item' },
    series: [{
      type: 'pie', radius: ['40%', '70%'],
      avoidLabelOverlap: false,
      itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
      label: { show: true, formatter: '{b}: {c}' },
      data: data.map(d => ({ name: d.name, value: d.count }))
    }]
  })
}

const renderTagChart = (data: TagStats[]) => {
  if (!tagChartRef.value) return
  if (!tagChart) tagChart = echarts.init(tagChartRef.value)
  tagChart.setOption({
    tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: { type: 'value' },
    yAxis: { type: 'category', data: data.map(d => d.name) },
    series: [{
      type: 'bar',
      data: data.map((d, i) => ({
        value: d.count,
        itemStyle: { color: d.color }
      }))
    }]
  })
}

onMounted(async () => {
  await nextTick()
  loadYearlyTrend()
  loadHeatmap()
  loadDuration()
  loadCategoryStats()
  loadTagStats()
})
</script>

<style scoped lang="scss">
.stats-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .chart {
    width: 100%;
  }
}
</style>
