<template>
  <div class="dashboard">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-item">
            <el-icon :size="32" color="#409EFF"><Folder /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ projectStats.total_projects || 0 }}</div>
              <div class="stat-label">项目总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-item">
            <el-icon :size="32" color="#67C23A"><CircleCheck /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ projectStats.completed_projects || 0 }}</div>
              <div class="stat-label">已完成项目</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-item">
            <el-icon :size="32" color="#E6A23C"><Document /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ projectStats.total_words || 0 }}</div>
              <div class="stat-label">翻译字数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-item">
            <el-icon :size="32" color="#F56C6C"><Money /></el-icon>
            <div class="stat-info">
              <div class="stat-value">{{ paymentStats.total_revenue?.toFixed(2) || '0.00' }}</div>
              <div class="stat-label">总收入(元)</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="16">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>收入趋势</span>
            </div>
          </template>
          <div ref="chartRef" style="height: 300px"></div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>项目状态分布</span>
            </div>
          </template>
          <div ref="pieChartRef" style="height: 300px"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="24">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>最近项目</span>
              <el-button type="primary" link @click="$router.push('/projects')">查看全部</el-button>
            </div>
          </template>
          <el-table :data="recentProjects" stripe>
            <el-table-column prop="title" label="项目名称" />
            <el-table-column prop="source_lang" label="源语言" width="100" />
            <el-table-column prop="target_lang" label="目标语言" width="100" />
            <el-table-column prop="word_count" label="字数" width="100" />
            <el-table-column prop="total_amount" label="金额(元)" width="120">
              <template #default="{ row }">{{ row.total_amount?.toFixed(2) }}</template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="120">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="创建时间" width="180">
              <template #default="{ row }">{{ formatDate(row.created_at) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { getProjectStatistics } from '@/api/statistics'
import { getPaymentStatistics, getRevenueTrend } from '@/api/statistics'
import { listProjects } from '@/api/project'
import * as echarts from 'echarts'
import dayjs from 'dayjs'

const chartRef = ref<HTMLElement>()
const pieChartRef = ref<HTMLElement>()
let lineChart: echarts.ECharts | null = null
let pieChart: echarts.ECharts | null = null

const projectStats = ref<any>({})
const paymentStats = ref<any>({})
const recentProjects = ref<any[]>([])

async function loadData() {
  try {
    const [projStats, payStats, revenueData, projects] = await Promise.all([
      getProjectStatistics(),
      getPaymentStatistics(),
      getRevenueTrend({ months: 6 }),
      listProjects({ page: 1, page_size: 5 })
    ])
    projectStats.value = projStats
    paymentStats.value = payStats

    if (Array.isArray(projects)) {
      recentProjects.value = projects
    } else if (projects?.list) {
      recentProjects.value = projects.list
    }

    renderRevenueChart(revenueData || [])
    renderProjectPieChart(projStats)
  } catch (e) {
    console.error('加载数据失败', e)
  }
}

function renderRevenueChart(data: any[]) {
  if (!chartRef.value) return
  lineChart = echarts.init(chartRef.value)
  lineChart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: {
      type: 'category',
      data: data.map(d => d.month)
    },
    yAxis: { type: 'value' },
    series: [{
      name: '收入',
      type: 'line',
      smooth: true,
      data: data.map(d => d.revenue),
      areaStyle: { opacity: 0.3 },
      itemStyle: { color: '#409EFF' }
    }]
  })
}

function renderProjectPieChart(stats: any) {
  if (!pieChartRef.value) return
  pieChart = echarts.init(pieChartRef.value)
  const data = [
    { name: '待审核', value: stats.pending_projects || 0 },
    { name: '进行中', value: stats.in_progress_projects || 0 },
    { name: '已完成', value: stats.completed_projects || 0 }
  ]
  pieChart.setOption({
    tooltip: { trigger: 'item' },
    series: [{
      type: 'pie',
      radius: ['40%', '70%'],
      data,
      label: { show: true, formatter: '{b}: {c}' }
    }]
  })
}

function getStatusType(status: string) {
  const map: Record<string, string> = {
    pending: 'warning',
    approved: 'primary',
    assigned: 'info',
    in_progress: '',
    review: 'warning',
    completed: 'success',
    cancelled: 'danger'
  }
  return map[status] || ''
}

function getStatusText(status: string) {
  const map: Record<string, string> = {
    pending: '待审核',
    approved: '已通过',
    assigned: '已分配',
    in_progress: '进行中',
    review: '审核中',
    completed: '已完成',
    cancelled: '已取消'
  }
  return map[status] || status
}

function formatDate(date: string) {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

function handleResize() {
  lineChart?.resize()
  pieChart?.resize()
}

onMounted(() => {
  loadData()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  lineChart?.dispose()
  pieChart?.dispose()
  window.removeEventListener('resize', handleResize)
})
</script>

<style lang="scss" scoped>
.dashboard {
  .stat-card {
    .stat-item {
      display: flex;
      align-items: center;
      gap: 16px;
    }

    .stat-info {
      .stat-value {
        font-size: 24px;
        font-weight: bold;
        color: #303133;
      }

      .stat-label {
        font-size: 14px;
        color: #909399;
        margin-top: 4px;
      }
    }
  }

  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
}
</style>
