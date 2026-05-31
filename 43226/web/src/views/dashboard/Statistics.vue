<template>
  <div class="statistics-page">
    <div class="filter-card card-shadow p-20 mb-20">
      <el-form :inline="true" :model="filterForm" class="filter-form">
        <el-form-item label="日期范围">
          <el-date-picker
            v-model="filterForm.dateRange"
            type="daterange"
            range-separator="至"
            start-placeholder="开始日期"
            end-placeholder="结束日期"
            value-format="YYYY-MM-DD"
            @change="handleDateChange"
          />
        </el-form-item>
        <el-form-item label="快速筛选">
          <el-radio-group v-model="filterForm.quickDate" @change="handleQuickDateChange">
            <el-radio-button value="today">今日</el-radio-button>
            <el-radio-button value="week">本周</el-radio-button>
            <el-radio-button value="month">本月</el-radio-button>
            <el-radio-button value="quarter">本季度</el-radio-button>
            <el-radio-button value="year">本年</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="handleSearch">
            <el-icon><Search /></el-icon>
            查询
          </el-button>
          <el-button @click="handleReset">重置</el-button>
        </el-form-item>
        <el-form-item class="export-buttons">
          <el-button type="success" @click="exportExcel">
            <el-icon><Download /></el-icon>
            导出Excel
          </el-button>
          <el-button type="warning" @click="exportPDF">
            <el-icon><Document /></el-icon>
            导出PDF
          </el-button>
        </el-form-item>
      </el-form>
    </div>

    <el-row :gutter="20">
      <el-col :span="6">
        <div class="stat-card card-shadow p-20">
          <div class="stat-icon icon-blue">
            <el-icon size="28"><User /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.totalVisitors }}</div>
            <div class="stat-label">总参观人数</div>
            <div class="stat-trend" :class="stats.visitorTrend >= 0 ? 'trend-up' : 'trend-down'">
              <el-icon v-if="stats.visitorTrend >= 0"><Top /></el-icon>
              <el-icon v-else><Bottom /></el-icon>
              {{ Math.abs(stats.visitorTrend) }}%
            </div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card card-shadow p-20">
          <div class="stat-icon icon-green">
            <el-icon size="28"><CircleCheck /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.conversionRate }}%</div>
            <div class="stat-label">预约转化率</div>
            <div class="stat-trend" :class="stats.conversionTrend >= 0 ? 'trend-up' : 'trend-down'">
              <el-icon v-if="stats.conversionTrend >= 0"><Top /></el-icon>
              <el-icon v-else><Bottom /></el-icon>
              {{ Math.abs(stats.conversionTrend) }}%
            </div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card card-shadow p-20">
          <div class="stat-icon icon-orange">
            <el-icon size="28"><Tickets /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.totalReservations }}</div>
            <div class="stat-label">预约总数</div>
            <div class="stat-trend" :class="stats.reservationTrend >= 0 ? 'trend-up' : 'trend-down'">
              <el-icon v-if="stats.reservationTrend >= 0"><Top /></el-icon>
              <el-icon v-else><Bottom /></el-icon>
              {{ Math.abs(stats.reservationTrend) }}%
            </div>
          </div>
        </div>
      </el-col>
      <el-col :span="6">
        <div class="stat-card card-shadow p-20">
          <div class="stat-icon icon-red">
            <el-icon size="28"><Star /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.avgRating }}</div>
            <div class="stat-label">平均评分</div>
            <div class="stat-trend" :class="stats.ratingTrend >= 0 ? 'trend-up' : 'trend-down'">
              <el-icon v-if="stats.ratingTrend >= 0"><Top /></el-icon>
              <el-icon v-else><Bottom /></el-icon>
              {{ Math.abs(stats.ratingTrend) }}%
            </div>
          </div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
      <el-col :span="12">
        <div class="ranking-card card-shadow p-20">
          <h3 class="card-title">
            <el-icon class="title-icon"><Trophy /></el-icon>
            热门展览 TOP3
          </h3>
          <div class="ranking-list">
            <div
              v-for="(item, index) in topExhibitions"
              :key="item.id"
              class="ranking-item"
            >
              <div class="rank-number" :class="`rank-${index + 1}`">{{ index + 1 }}</div>
              <div class="rank-info">
                <div class="rank-name">{{ item.title }}</div>
                <div class="rank-sub">
                  <el-tag size="small" type="info">{{ item.location }}</el-tag>
                  <span class="view-count">
                    <el-icon><View /></el-icon>
                    {{ item.view_count }} 次浏览
                  </span>
                </div>
              </div>
              <div class="rank-progress">
                <el-progress
                  :percentage="Math.round((item.view_count / topExhibitions[0].view_count) * 100)"
                  :color="getRankColor(index)"
                  :stroke-width="8"
                  :show-text="false"
                />
                <span class="progress-text">{{ item.view_count }}</span>
              </div>
            </div>
          </div>
        </div>
      </el-col>
      <el-col :span="12">
        <div class="ranking-card card-shadow p-20">
          <h3 class="card-title">
            <el-icon class="title-icon"><Gem /></el-icon>
            热门藏品 TOP3
          </h3>
          <div class="ranking-list">
            <div
              v-for="(item, index) in topCollections"
              :key="item.id"
              class="ranking-item"
            >
              <div class="rank-number" :class="`rank-${index + 1}`">{{ index + 1 }}</div>
              <div class="rank-info">
                <div class="rank-name">{{ item.name }}</div>
                <div class="rank-sub">
                  <el-tag size="small" type="info">{{ item.category?.name }}</el-tag>
                  <span class="view-count">
                    <el-icon><View /></el-icon>
                    {{ item.view_count }} 次浏览
                  </span>
                </div>
              </div>
              <div class="rank-progress">
                <el-progress
                  :percentage="Math.round((item.view_count / topCollections[0].view_count) * 100)"
                  :color="getRankColor(index)"
                  :stroke-width="8"
                  :show-text="false"
                />
                <span class="progress-text">{{ item.view_count }}</span>
              </div>
            </div>
          </div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
      <el-col :span="16">
        <div class="chart-card card-shadow p-20">
          <h3 class="card-title">参观趋势分析</h3>
          <div ref="visitTrendRef" class="chart-container"></div>
        </div>
      </el-col>
      <el-col :span="8">
        <div class="chart-card card-shadow p-20">
          <h3 class="card-title">预约状态分布</h3>
          <div ref="reservationPieRef" class="chart-container pie-chart"></div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="mt-20">
      <el-col :span="12">
        <div class="chart-card card-shadow p-20">
          <h3 class="card-title">展览参观人数统计</h3>
          <div ref="exhibitionBarRef" class="chart-container"></div>
        </div>
      </el-col>
      <el-col :span="12">
        <div class="chart-card card-shadow p-20">
          <h3 class="card-title">访客时段分布</h3>
          <div ref="timeSlotRef" class="chart-container"></div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, nextTick, onUnmounted } from 'vue'
import * as echarts from 'echarts'
import type { Exhibition, Collection } from '@/types'
import dayjs from 'dayjs'
import { ElMessage } from 'element-plus'

const loading = ref(false)

const filterForm = reactive({
  dateRange: [] as string[],
  quickDate: 'week' as string
})

const stats = ref({
  totalVisitors: 15680,
  conversionRate: 78.5,
  totalReservations: 3256,
  avgRating: 4.8,
  visitorTrend: 12.5,
  conversionTrend: 5.2,
  reservationTrend: 8.3,
  ratingTrend: 2.1
})

const topExhibitions = ref<Exhibition[]>([
  {
    id: 1,
    title: '古代书画艺术展',
    location: '一号展厅',
    view_count: 8520,
    status: 'published'
  } as Exhibition,
  {
    id: 2,
    title: '瓷器珍品特展',
    location: '二号展厅',
    view_count: 6340,
    status: 'published'
  } as Exhibition,
  {
    id: 3,
    title: '青铜器文化展',
    location: '三号展厅',
    view_count: 5120,
    status: 'published'
  } as Exhibition
])

const topCollections = ref<Collection[]>([
  {
    id: 1,
    name: '青花缠枝莲纹瓶',
    view_count: 12580,
    category: { name: '瓷器' } as any
  } as Collection,
  {
    id: 2,
    name: '清明上河图局部',
    view_count: 9870,
    category: { name: '书画' } as any
  } as Collection,
  {
    id: 3,
    name: '司母戊鼎仿制品',
    view_count: 8450,
    category: { name: '青铜器' } as any
  } as Collection
])

const visitTrendRef = ref<HTMLElement>()
const reservationPieRef = ref<HTMLElement>()
const exhibitionBarRef = ref<HTMLElement>()
const timeSlotRef = ref<HTMLElement>()

let visitTrendChart: echarts.ECharts | null = null
let reservationPieChart: echarts.ECharts | null = null
let exhibitionBarChart: echarts.ECharts | null = null
let timeSlotChart: echarts.ECharts | null = null

const getRankColor = (index: number) => {
  const colors = ['#f7ba2a', '#c0c4cc', '#e6a23c']
  return colors[index] || '#409eff'
}

const handleQuickDateChange = (value: string) => {
  const today = dayjs()
  let start: dayjs.Dayjs
  let end: dayjs.Dayjs = today

  switch (value) {
    case 'today':
      start = today
      end = today
      break
    case 'week':
      start = today.startOf('week' as dayjs.OpUnitType)
      break
    case 'month':
      start = today.startOf('month' as dayjs.OpUnitType)
      break
    case 'quarter':
      start = today.startOf('quarter' as dayjs.OpUnitType)
      break
    case 'year':
      start = today.startOf('year' as dayjs.OpUnitType)
      break
    default:
      start = today.startOf('week' as dayjs.OpUnitType)
  }

  filterForm.dateRange = [
    start.format('YYYY-MM-DD'),
    end.format('YYYY-MM-DD')
  ]
}

const handleDateChange = () => {
  filterForm.quickDate = ''
}

const handleSearch = () => {
  loading.value = true
  setTimeout(() => {
    loading.value = false
    initCharts()
    ElMessage.success('数据刷新成功')
  }, 500)
}

const handleReset = () => {
  filterForm.dateRange = []
  filterForm.quickDate = 'week'
  handleQuickDateChange('week')
  handleSearch()
}

const generateTrendData = (days: number) => {
  const dates: string[] = []
  const visitors: number[] = []
  const reservations: number[] = []
  const today = dayjs()

  for (let i = days - 1; i >= 0; i--) {
    const date = today.subtract(i, 'day')
    dates.push(date.format('MM-DD'))
    visitors.push(Math.floor(Math.random() * 200) + 100)
    reservations.push(Math.floor(Math.random() * 100) + 50)
  }

  return { dates, visitors, reservations }
}

const initCharts = () => {
  const trendData = generateTrendData(7)

  if (visitTrendRef.value) {
    visitTrendChart = echarts.init(visitTrendRef.value)
    visitTrendChart.setOption({
      tooltip: { trigger: 'axis' },
      legend: { data: ['实际参观', '预约人数'] },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: {
        type: 'category',
        boundaryGap: false,
        data: trendData.dates
      },
      yAxis: { type: 'value' },
      series: [
        {
          name: '实际参观',
          type: 'line',
          smooth: true,
          data: trendData.visitors,
          lineStyle: { width: 3 },
          itemStyle: { color: '#409eff' },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(64, 158, 255, 0.4)' },
              { offset: 1, color: 'rgba(64, 158, 255, 0.05)' }
            ])
          }
        },
        {
          name: '预约人数',
          type: 'line',
          smooth: true,
          data: trendData.reservations,
          lineStyle: { width: 3 },
          itemStyle: { color: '#67c23a' }
        }
      ]
    })
  }

  if (reservationPieRef.value) {
    reservationPieChart = echarts.init(reservationPieRef.value)
    reservationPieChart.setOption({
      tooltip: { trigger: 'item', formatter: '{b}: {c} ({d}%)' },
      legend: { orient: 'vertical', left: 'left' },
      series: [
        {
          type: 'pie',
          radius: ['40%', '70%'],
          center: ['60%', '50%'],
          avoidLabelOverlap: false,
          itemStyle: { borderRadius: 8, borderColor: '#fff', borderWidth: 2 },
          label: { show: false },
          emphasis: {
            label: { show: true, fontSize: 16, fontWeight: 'bold' }
          },
          labelLine: { show: false },
          data: [
            { value: 1856, name: '已完成', itemStyle: { color: '#67c23a' } },
            { value: 892, name: '已确认', itemStyle: { color: '#409eff' } },
            { value: 328, name: '待确认', itemStyle: { color: '#e6a23c' } },
            { value: 180, name: '已取消', itemStyle: { color: '#909399' } }
          ]
        }
      ]
    })
  }

  if (exhibitionBarRef.value) {
    exhibitionBarChart = echarts.init(exhibitionBarRef.value)
    exhibitionBarChart.setOption({
      tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
      legend: { data: ['参观人数', '预约人数'] },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: {
        type: 'category',
        data: ['古代书画', '瓷器珍品', '青铜器', '玉器', '近现代艺术', '钱币']
      },
      yAxis: { type: 'value' },
      series: [
        {
          name: '参观人数',
          type: 'bar',
          barWidth: '30%',
          data: [8520, 6340, 5120, 4280, 3650, 2890],
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#667eea' },
              { offset: 1, color: '#764ba2' }
            ]),
            borderRadius: [4, 4, 0, 0]
          }
        },
        {
          name: '预约人数',
          type: 'bar',
          barWidth: '30%',
          data: [6800, 5200, 4100, 3400, 2900, 2200],
          itemStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: '#43e97b' },
              { offset: 1, color: '#38f9d7' }
            ]),
            borderRadius: [4, 4, 0, 0]
          }
        }
      ]
    })
  }

  if (timeSlotRef.value) {
    timeSlotChart = echarts.init(timeSlotRef.value)
    timeSlotChart.setOption({
      tooltip: { trigger: 'axis' },
      legend: { data: ['访客数量'] },
      grid: { left: '3%', right: '4%', bottom: '3%', containLabel: true },
      xAxis: {
        type: 'category',
        data: ['09:00', '10:00', '11:00', '12:00', '13:00', '14:00', '15:00', '16:00', '17:00']
      },
      yAxis: { type: 'value' },
      series: [
        {
          name: '访客数量',
          type: 'line',
          smooth: true,
          data: [120, 280, 350, 180, 150, 320, 420, 380, 200],
          lineStyle: { width: 3 },
          itemStyle: { color: '#f093fb' },
          areaStyle: {
            color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
              { offset: 0, color: 'rgba(240, 147, 251, 0.4)' },
              { offset: 1, color: 'rgba(240, 147, 251, 0.05)' }
            ])
          }
        }
      ]
    })
  }
}

const handleResize = () => {
  visitTrendChart?.resize()
  reservationPieChart?.resize()
  exhibitionBarChart?.resize()
  timeSlotChart?.resize()
}

const exportExcel = () => {
  const headers = ['日期', '参观人数', '预约人数', '转化率', '热门展览', '热门藏品']
  const rows = [
    ['2024-01-01', '230', '180', '78%', '古代书画艺术展', '青花缠枝莲纹瓶'],
    ['2024-01-02', '210', '165', '78%', '瓷器珍品特展', '清明上河图局部'],
    ['2024-01-03', '195', '150', '77%', '青铜器文化展', '司母戊鼎仿制品'],
    ['2024-01-04', '220', '175', '79%', '古代书画艺术展', '青花缠枝莲纹瓶'],
    ['2024-01-05', '240', '190', '79%', '瓷器珍品特展', '清明上河图局部'],
    ['2024-01-06', '310', '250', '81%', '青铜器文化展', '司母戊鼎仿制品'],
    ['2024-01-07', '290', '230', '79%', '古代书画艺术展', '青花缠枝莲纹瓶']
  ]

  let csvContent = '\uFEFF'
  csvContent += headers.join(',') + '\n'
  rows.forEach(row => {
    csvContent += row.join(',') + '\n'
  })

  const blob = new Blob([csvContent], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  const url = URL.createObjectURL(blob)
  link.setAttribute('href', url)
  link.setAttribute('download', `统计报表_${dayjs().format('YYYYMMDD')}.csv`)
  link.style.visibility = 'hidden'
  document.body.appendChild(link)
  link.click()
  document.body.removeChild(link)
  ElMessage.success('Excel导出成功')
}

const exportPDF = () => {
  const printContent = document.querySelector('.statistics-page')
  if (printContent) {
    const printWindow = window.open('', '_blank')
    if (printWindow) {
      printWindow.document.write(`
        <!DOCTYPE html>
        <html>
        <head>
          <title>博物馆统计报表</title>
          <style>
            body { font-family: Arial, sans-serif; padding: 20px; }
            h1 { text-align: center; color: #333; }
            .header { margin-bottom: 20px; }
            .stats { display: flex; justify-content: space-around; margin: 20px 0; }
            .stat-item { text-align: center; }
            .stat-value { font-size: 24px; font-weight: bold; color: #409eff; }
            .stat-label { color: #666; }
            table { width: 100%; border-collapse: collapse; margin: 20px 0; }
            th, td { border: 1px solid #ddd; padding: 12px; text-align: left; }
            th { background: #f5f7fa; }
            .rank-list { margin: 10px 0; }
            .rank-item { padding: 8px; border-bottom: 1px solid #eee; }
            .date-range { text-align: right; color: #666; }
          </style>
        </head>
        <body>
          <div class="header">
            <h1>博物馆统计分析报表</h1>
            <div class="date-range">生成日期：${dayjs().format('YYYY年MM月DD日')}</div>
          </div>
          <div class="stats">
            <div class="stat-item">
              <div class="stat-value">${stats.value.totalVisitors}</div>
              <div class="stat-label">总参观人数</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">${stats.value.conversionRate}%</div>
              <div class="stat-label">预约转化率</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">${stats.value.totalReservations}</div>
              <div class="stat-label">预约总数</div>
            </div>
            <div class="stat-item">
              <div class="stat-value">${stats.value.avgRating}</div>
              <div class="stat-label">平均评分</div>
            </div>
          </div>
          <h3>热门展览 TOP3</h3>
          <div class="rank-list">
            ${topExhibitions.value.map((item, i) => `
              <div class="rank-item">${i + 1}. ${item.title} - ${item.view_count}次浏览</div>
            `).join('')}
          </div>
          <h3>热门藏品 TOP3</h3>
          <div class="rank-list">
            ${topCollections.value.map((item, i) => `
              <div class="rank-item">${i + 1}. ${item.name} - ${item.view_count}次浏览</div>
            `).join('')}
          </div>
          <h3>参观数据明细</h3>
          <table>
            <thead>
              <tr>
                <th>日期</th>
                <th>参观人数</th>
                <th>预约人数</th>
                <th>转化率</th>
              </tr>
            </thead>
            <tbody>
              ${generateTrendData(7).dates.map((date, i) => `
                <tr>
                  <td>${date}</td>
                  <td>${generateTrendData(7).visitors[i]}</td>
                  <td>${generateTrendData(7).reservations[i]}</td>
                  <td>${Math.round(generateTrendData(7).reservations[i] / generateTrendData(7).visitors[i] * 100)}%</td>
                </tr>
              `).join('')}
            </tbody>
          </table>
        </body>
        </html>
      `)
      printWindow.document.close()
      printWindow.focus()
      setTimeout(() => {
        printWindow.print()
      }, 500)
    }
  }
  ElMessage.success('PDF导出功能已打开，请在新窗口中选择打印为PDF')
}

onMounted(() => {
  handleQuickDateChange(filterForm.quickDate)
  nextTick(() => {
    initCharts()
  })
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  window.removeEventListener('resize', handleResize)
  visitTrendChart?.dispose()
  reservationPieChart?.dispose()
  exhibitionBarChart?.dispose()
  timeSlotChart?.dispose()
})
</script>

<style scoped lang="scss">
.statistics-page {
  .filter-card {
    .filter-form {
      display: flex;
      flex-wrap: wrap;
      align-items: center;
      gap: 10px;

      .export-buttons {
        margin-left: auto;
      }
    }
  }

  .stat-card {
    display: flex;
    align-items: center;
    gap: 20px;
    border-radius: 8px;
    position: relative;
    overflow: hidden;

    .stat-icon {
      width: 60px;
      height: 60px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;
      flex-shrink: 0;

      &.icon-blue { background: linear-gradient(135deg, #667eea, #764ba2); }
      &.icon-green { background: linear-gradient(135deg, #43e97b, #38f9d7); }
      &.icon-orange { background: linear-gradient(135deg, #fa709a, #fee140); }
      &.icon-red { background: linear-gradient(135deg, #f093fb, #f5576c); }
    }

    .stat-info {
      flex: 1;

      .stat-value {
        font-size: 28px;
        font-weight: 700;
        margin-bottom: 4px;
      }

      .stat-label {
        color: #909399;
        font-size: 14px;
        margin-bottom: 4px;
      }

      .stat-trend {
        font-size: 12px;
        display: flex;
        align-items: center;
        gap: 2px;

        &.trend-up {
          color: #67c23a;
        }

        &.trend-down {
          color: #f56c6c;
        }
      }
    }
  }

  .ranking-card {
    border-radius: 8px;

    .card-title {
      font-size: 18px;
      margin-bottom: 20px;
      display: flex;
      align-items: center;
      gap: 8px;

      .title-icon {
        color: #e6a23c;
      }
    }

    .ranking-list {
      .ranking-item {
        display: flex;
        align-items: center;
        gap: 15px;
        padding: 15px 0;
        border-bottom: 1px solid #f0f2f5;

        &:last-child {
          border-bottom: none;
        }

        .rank-number {
          width: 32px;
          height: 32px;
          border-radius: 50%;
          display: flex;
          align-items: center;
          justify-content: center;
          font-weight: bold;
          color: #fff;
          font-size: 16px;
          flex-shrink: 0;

          &.rank-1 { background: linear-gradient(135deg, #ffd700, #ffb700); }
          &.rank-2 { background: linear-gradient(135deg, #c0c4cc, #909399); }
          &.rank-3 { background: linear-gradient(135deg, #cd7f32, #b8860b); }
        }

        .rank-info {
          flex: 1;
          min-width: 0;

          .rank-name {
            font-weight: 500;
            margin-bottom: 6px;
            overflow: hidden;
            text-overflow: ellipsis;
            white-space: nowrap;
          }

          .rank-sub {
            display: flex;
            align-items: center;
            gap: 10px;

            .view-count {
              display: flex;
              align-items: center;
              gap: 4px;
              color: #909399;
              font-size: 12px;
            }
          }
        }

        .rank-progress {
          width: 120px;
          display: flex;
          align-items: center;
          gap: 10px;

          .progress-text {
            font-size: 14px;
            font-weight: 600;
            color: #606266;
            min-width: 50px;
            text-align: right;
          }
        }
      }
    }
  }

  .chart-card {
    border-radius: 8px;

    .card-title {
      font-size: 18px;
      margin-bottom: 20px;
    }

    .chart-container {
      height: 320px;

      &.pie-chart {
        height: 320px;
      }
    }
  }
}
</style>
