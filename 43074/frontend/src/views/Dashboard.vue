<template>
  <div class="dashboard">
    <el-row :gutter="20" class="stats-cards">
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card card-hover">
          <div class="stat-content">
            <div class="stat-icon blue">
              <el-icon :size="32"><Reading /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ overview?.total_books || 0 }}</div>
              <div class="stat-label">全部图书</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card card-hover">
          <div class="stat-content">
            <div class="stat-icon green">
              <el-icon :size="32"><Clock /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ overview?.reading_books || 0 }}</div>
              <div class="stat-label">正在阅读</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card card-hover">
          <div class="stat-content">
            <div class="stat-icon purple">
              <el-icon :size="32"><CircleCheck /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ overview?.completed_books || 0 }}</div>
              <div class="stat-label">已读完</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card shadow="hover" class="stat-card card-hover">
          <div class="stat-content">
            <div class="stat-icon orange">
              <el-icon :size="32"><Share /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ overview?.currently_borrowed || 0 }}</div>
              <div class="stat-label">借出中</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="content-row">
      <el-col :span="16">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>年度阅读趋势</span>
              <el-select v-model="trendYear" style="width: 100px" @change="loadYearlyTrend">
                <el-option v-for="y in years" :key="y" :label="y" :value="y" />
              </el-select>
            </div>
          </template>
          <div ref="trendChartRef" class="chart-container" style="height: 300px"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card shadow="hover">
          <template #header>
            <span>正在阅读</span>
          </template>
          <div v-loading="loadingReading" class="reading-list">
            <div v-if="!currentlyReading.length" class="empty">
              <el-empty description="暂无在读书籍" :image-size="80" />
            </div>
            <div v-for="book in currentlyReading" :key="book.id" class="reading-item card-hover" @click="goToBook(book.id)">
              <img v-if="book.cover_image" :src="book.cover_image" class="mini-cover" />
              <div v-else class="mini-cover placeholder">
                <el-icon><Reading /></el-icon>
              </div>
              <div class="reading-info">
                <div class="book-title">{{ book.title }}</div>
                <div class="book-progress">
                  <el-progress :percentage="Math.round(book.reading_progress)" :stroke-width="8" />
                </div>
              </div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="content-row">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>阅读热力图</span>
              <el-select v-model="heatmapYear" style="width: 100px" @change="loadHeatmap">
                <el-option v-for="y in years" :key="y" :label="y" :value="y" />
              </el-select>
            </div>
          </template>
          <div ref="heatmapChartRef" class="chart-container" style="height: 200px"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <span>分类统计</span>
          </template>
          <div ref="categoryChartRef" class="chart-container" style="height: 200px"></div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import { getOverview, getYearlyTrend, getReadingHeatmap, getCategoryStats } from '@/api/stats'
import { getCurrentlyReading } from '@/api/book'
import type { OverviewStats, Book, MonthlyStats, HeatmapData, CategoryStats } from '@/types'

const router = useRouter()

const overview = ref<OverviewStats | null>(null)
const currentlyReading = ref<Book[]>([])
const loadingReading = ref(false)
const trendYear = ref(new Date().getFullYear())
const heatmapYear = ref(new Date().getFullYear())
const years = [new Date().getFullYear(), new Date().getFullYear() - 1, new Date().getFullYear() - 2]

const trendChartRef = ref<HTMLElement>()
const heatmapChartRef = ref<HTMLElement>()
const categoryChartRef = ref<HTMLElement>()
let trendChart: echarts.ECharts | null = null
let heatmapChart: echarts.ECharts | null = null
let categoryChart: echarts.ECharts | null = null

const loadOverview = async () => {
  try {
    overview.value = await getOverview()
  } catch (e) {}
}

const loadCurrentlyReading = async () => {
  loadingReading.value = true
  try {
    currentlyReading.value = await getCurrentlyReading()
  } catch (e) {
  } finally {
    loadingReading.value = false
  }
}

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

const loadCategoryStats = async () => {
  try {
    const data = await getCategoryStats()
    renderCategoryChart(data)
  } catch (e) {}
}

const renderTrendChart = (data: MonthlyStats[]) => {
  if (!trendChartRef.value) return
  if (!trendChart) {
    trendChart = echarts.init(trendChartRef.value)
  }
  trendChart.setOption({
    tooltip: { trigger: 'axis' },
    legend: { data: ['完成数量', '阅读页数'] },
    grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
    xAxis: {
      type: 'category',
      data: data.map(d => d.month)
    },
    yAxis: [
      { type: 'value', name: '数量' },
      { type: 'value', name: '页数' }
    ],
    series: [
      {
        name: '完成数量',
        type: 'bar',
        data: data.map(d => d.completed),
        itemStyle: { color: '#409eff' }
      },
      {
        name: '阅读页数',
        type: 'line',
        yAxisIndex: 1,
        data: data.map(d => d.pages_read),
        itemStyle: { color: '#67c23a' },
        smooth: true
      }
    ]
  })
}

const renderHeatmapChart = (data: HeatmapData[]) => {
  if (!heatmapChartRef.value) return
  if (!heatmapChart) {
    heatmapChart = echarts.init(heatmapChartRef.value)
  }
  const dateMap: Record<string, number> = {}
  data.forEach(d => { dateMap[d.date] = d.count })
  const startDate = new Date(heatmapYear.value, 0, 1)
  const endDate = new Date(heatmapYear.value, 11, 31)
  const calendarData: [string, number][] = []
  for (let d = new Date(startDate); d <= endDate; d.setDate(d.getDate() + 1)) {
    const dateStr = d.toISOString().split('T')[0]
    calendarData.push([dateStr, dateMap[dateStr] || 0])
  }
  heatmapChart.setOption({
    tooltip: { position: 'top' },
    visualMap: {
      min: 0,
      max: Math.max(...data.map(d => d.count), 1),
      calculable: true,
      orient: 'horizontal',
      left: 'center',
      bottom: '0',
      inRange: { color: ['#ebedf0', '#c6e48b', '#7bc96f', '#239a3b', '#196127'] }
    },
    calendar: {
      top: 20,
      left: 50,
      right: 50,
      cellSize: ['auto', 13],
      range: heatmapYear.value,
      itemStyle: { borderWidth: 0.5, borderColor: '#fff' },
      yearLabel: { show: false },
      monthLabel: { nameMap: 'ZH' },
      dayLabel: { firstDay: 1, nameMap: ['日', '一', '二', '三', '四', '五', '六'] }
    },
    series: [{
      type: 'heatmap',
      coordinateSystem: 'calendar',
      data: calendarData
    }]
  })
}

const renderCategoryChart = (data: CategoryStats[]) => {
  if (!categoryChartRef.value) return
  if (!categoryChart) {
    categoryChart = echarts.init(categoryChartRef.value)
  }
  categoryChart.setOption({
    tooltip: { trigger: 'item' },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      avoidLabelOverlap: false,
      itemStyle: { borderRadius: 6, borderColor: '#fff', borderWidth: 2 },
      label: { show: true, formatter: '{b}: {c}' },
      data: data.map(d => ({ name: d.name, value: d.count }))
    }]
  })
}

const goToBook = (id: number) => {
  router.push(`/books/${id}`)
}

onMounted(async () => {
  loadOverview()
  loadCurrentlyReading()
  await nextTick()
  loadYearlyTrend()
  loadHeatmap()
  loadCategoryStats()
})
</script>

<style scoped lang="scss">
.dashboard {
  .stats-cards {
    margin-bottom: 20px;
  }

  .stat-card {
    .stat-content {
      display: flex;
      align-items: center;
      gap: 16px;
    }

    .stat-icon {
      width: 56px;
      height: 56px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;

      &.blue { background: linear-gradient(135deg, #667eea 0%, #764ba2 100%); }
      &.green { background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%); }
      &.purple { background: linear-gradient(135deg, #8e2de2 0%, #4a00e0 100%); }
      &.orange { background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%); }
    }

    .stat-info {
      .stat-value {
        font-size: 28px;
        font-weight: 600;
        color: #303133;
        line-height: 1.2;
      }
      .stat-label {
        font-size: 14px;
        color: #909399;
        margin-top: 4px;
      }
    }
  }

  .content-row {
    margin-bottom: 20px;
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .reading-list {
    max-height: 300px;
    overflow-y: auto;

    .empty {
      padding: 20px 0;
    }
  }

  .reading-item {
    display: flex;
    gap: 12px;
    padding: 12px;
    border-radius: 8px;
    cursor: pointer;
    margin-bottom: 8px;

    &:hover {
      background-color: #f5f7fa;
    }

    .mini-cover {
      width: 48px;
      height: 64px;
      object-fit: cover;
      border-radius: 4px;

      &.placeholder {
        background-color: #ebeef5;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #c0c4cc;
      }
    }

    .reading-info {
      flex: 1;
      min-width: 0;

      .book-title {
        font-weight: 500;
        color: #303133;
        margin-bottom: 8px;
        overflow: hidden;
        text-overflow: ellipsis;
        white-space: nowrap;
      }
    }
  }
}
</style>
