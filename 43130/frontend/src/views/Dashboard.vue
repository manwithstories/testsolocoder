<template>
  <div class="dashboard">
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-info">
              <p class="stat-label">婚礼总数</p>
              <p class="stat-value">{{ stats.total_weddings || 0 }}</p>
            </div>
            <el-icon class="stat-icon" color="#409EFF"><House /></el-icon>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-info">
              <p class="stat-label">嘉宾总数</p>
              <p class="stat-value">{{ stats.total_guests || 0 }}</p>
            </div>
            <el-icon class="stat-icon" color="#67C23A"><User /></el-icon>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-info">
              <p class="stat-label">供应商数</p>
              <p class="stat-value">{{ stats.total_vendors || 0 }}</p>
            </div>
            <el-icon class="stat-icon" color="#E6A23C"><OfficeBuilding /></el-icon>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-info">
              <p class="stat-label">预算执行率</p>
              <p class="stat-value">{{ formatRate(stats.budget?.rate) }}%</p>
            </div>
            <el-icon class="stat-icon" color="#F56C6C"><Money /></el-icon>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="chart-row">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <div class="card-header">
              <span>预算分布</span>
              <el-button type="primary" size="small" @click="exportReport">导出报告</el-button>
            </div>
          </template>
          <div ref="budgetChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <span>任务完成进度</span>
          </template>
          <div ref="taskChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="detail-row">
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <span>RSVP 统计</span>
          </template>
          <div class="rsvp-stats">
            <div class="rsvp-item accepted">
              <div class="rsvp-value">{{ stats.rsvp?.Accepted || 0 }}</div>
              <div class="rsvp-label">已接受</div>
            </div>
            <div class="rsvp-item declined">
              <div class="rsvp-value">{{ stats.rsvp?.Declined || 0 }}</div>
              <div class="rsvp-label">已拒绝</div>
            </div>
            <div class="rsvp-item pending">
              <div class="rsvp-value">{{ stats.rsvp?.Pending || 0 }}</div>
              <div class="rsvp-label">待回复</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card shadow="hover">
          <template #header>
            <span>即将到期任务</span>
          </template>
          <el-table :data="upcomingTasks" style="width: 100%" max-height="300">
            <el-table-column prop="title" label="任务名称" />
            <el-table-column prop="category" label="分类" width="100" />
            <el-table-column label="截止日期" width="120">
              <template #default="{ row }">
                {{ formatDate(row.due_date) }}
              </template>
            </el-table-column>
            <el-table-column prop="priority" label="优先级" width="80">
              <template #default="{ row }">
                <el-tag :type="priorityType(row.priority)" size="small">{{ row.priority }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, nextTick } from 'vue'
import { useWeddingStore } from '@/store/wedding'
import { dashboardApi } from '@/api/dashboard'
import { ElMessage } from 'element-plus'
import { House, User, OfficeBuilding, Money } from '@element-plus/icons-vue'
import * as echarts from 'echarts'
import dayjs from 'dayjs'

const weddingStore = useWeddingStore()

const stats = ref<any>({})
const budgetChartRef = ref<HTMLElement>()
const taskChartRef = ref<HTMLElement>()
const upcomingTasks = ref<any[]>([])

let budgetChart: echarts.ECharts | null = null
let taskChart: echarts.ECharts | null = null

function formatRate(rate?: number) {
  return rate?.toFixed(1) || '0'
}

function formatDate(date?: string) {
  return date ? dayjs(date).format('YYYY-MM-DD') : '-'
}

function priorityType(priority: string) {
  const types: Record<string, string> = {
    high: 'danger',
    medium: 'warning',
    low: 'success'
  }
  return types[priority] || 'info'
}

async function fetchStats() {
  try {
    const params: any = {}
    if (weddingStore.currentWeddingId) {
      params.wedding_id = weddingStore.currentWeddingId
    }
    const res = await dashboardApi.getStats(params)
    stats.value = res.data
  } catch (error) {
    console.error('Failed to fetch stats:', error)
  }
}

async function fetchBudgetChart() {
  try {
    const params: any = {}
    if (weddingStore.currentWeddingId) {
      params.wedding_id = weddingStore.currentWeddingId
    }
    const res = await dashboardApi.getBudgetChart(params)
    
    if (budgetChartRef.value) {
      budgetChart = echarts.init(budgetChartRef.value)
      budgetChart.setOption({
        tooltip: { trigger: 'item' },
        legend: { orient: 'vertical', left: 'left' },
        series: [{
          type: 'pie',
          radius: ['40%', '70%'],
          data: res.data.map((item: any) => ({
            name: item.category,
            value: item.actual_cost
          }))
        }]
      })
    }
  } catch (error) {
    console.error('Failed to fetch budget chart:', error)
  }
}

async function fetchTaskChart() {
  try {
    const params: any = {}
    if (weddingStore.currentWeddingId) {
      params.wedding_id = weddingStore.currentWeddingId
    }
    const res = await dashboardApi.getTaskProgress(params)
    
    if (taskChartRef.value) {
      taskChart = echarts.init(taskChartRef.value)
      taskChart.setOption({
        tooltip: { trigger: 'axis' },
        legend: { data: ['总数', '已完成'] },
        xAxis: { type: 'category', data: res.data.map((item: any) => item.category) },
        yAxis: { type: 'value' },
        series: [
          { name: '总数', type: 'bar', data: res.data.map((item: any) => item.total) },
          { name: '已完成', type: 'bar', data: res.data.map((item: any) => item.complete) }
        ]
      })
    }
  } catch (error) {
    console.error('Failed to fetch task chart:', error)
  }
}

async function fetchUpcomingTasks() {
  try {
    const params: any = {}
    if (weddingStore.currentWeddingId) {
      params.wedding_id = weddingStore.currentWeddingId
    }
    const res = await dashboardApi.getUpcomingTasks(params)
    upcomingTasks.value = res.data
  } catch (error) {
    console.error('Failed to fetch upcoming tasks:', error)
  }
}

async function exportReport() {
  try {
    const params: any = {}
    if (weddingStore.currentWeddingId) {
      params.wedding_id = weddingStore.currentWeddingId
    }
    const res = await dashboardApi.exportReport(params)
    
    const blob = new Blob([res as any], {
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = 'wedding_report.xlsx'
    link.click()
    URL.revokeObjectURL(url)
    
    ElMessage.success('报告导出成功')
  } catch (error) {
    ElMessage.error('导出报告失败')
  }
}

function handleResize() {
  budgetChart?.resize()
  taskChart?.resize()
}

onMounted(async () => {
  await fetchStats()
  await nextTick()
  await fetchBudgetChart()
  await fetchTaskChart()
  await fetchUpcomingTasks()
  
  window.addEventListener('resize', handleResize)
})
</script>

<style scoped>
.dashboard {
  padding: 0;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  border-radius: 8px;
}

.stat-content {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.stat-label {
  color: #909399;
  font-size: 14px;
  margin: 0;
}

.stat-value {
  font-size: 28px;
  font-weight: 600;
  color: #303133;
  margin: 8px 0 0;
}

.stat-icon {
  font-size: 48px;
  opacity: 0.8;
}

.chart-row {
  margin-bottom: 20px;
}

.card-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}

.chart-container {
  height: 300px;
}

.rsvp-stats {
  display: flex;
  justify-content: space-around;
  padding: 20px 0;
}

.rsvp-item {
  text-align: center;
  padding: 20px;
  border-radius: 8px;
  min-width: 100px;
}

.rsvp-item.accepted {
  background-color: #f0f9eb;
  color: #67C23A;
}

.rsvp-item.declined {
  background-color: #fef0f0;
  color: #F56C6C;
}

.rsvp-item.pending {
  background-color: #fdf6ec;
  color: #E6A23C;
}

.rsvp-value {
  font-size: 32px;
  font-weight: 600;
}

.rsvp-label {
  font-size: 14px;
  margin-top: 8px;
}
</style>
