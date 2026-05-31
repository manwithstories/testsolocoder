<template>
  <div class="page-container">
    <el-row :gutter="20">
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>项目统计</span>
            </div>
          </template>
          <div class="stat-grid">
            <div class="stat-item">
              <div class="stat-label">项目总数</div>
              <div class="stat-value">{{ projectStats.total_projects || 0 }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">待审核</div>
              <div class="stat-value warning">{{ projectStats.pending_projects || 0 }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">进行中</div>
              <div class="stat-value primary">{{ projectStats.in_progress_projects || 0 }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">已完成</div>
              <div class="stat-value success">{{ projectStats.completed_projects || 0 }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">完成率</div>
              <div class="stat-value">{{ projectStats.completion_rate }}</div>
            </div>
            <div class="stat-item">
              <div class="stat-label">总字数</div>
              <div class="stat-value">{{ projectStats.total_words || 0 }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>收入趋势</span>
          </template>
          <div ref="chartRef" style="height: 280px"></div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>语言对分布</span>
          </template>
          <el-table :data="languagePairStats" size="small">
            <el-table-column prop="language_pair" label="语言对" />
            <el-table-column prop="project_count" label="项目数" width="100" />
            <el-table-column prop="total_words" label="总字数" width="100" />
            <el-table-column prop="total_revenue" label="收入(元)" width="120">
              <template #default="{ row }">{{ row.total_revenue?.toFixed(2) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <div class="card-header">
              <span>译者效率排行</span>
            </div>
          </template>
          <el-table :data="translatorStats" size="small" max-height="300">
            <el-table-column prop="user.username" label="译者" />
            <el-table-column prop="completed_count" label="完成项目" width="100" />
            <el-table-column prop="total_words" label="总字数" width="100" />
            <el-table-column prop="avg_rating" label="评分" width="80">
              <template #default="{ row }">{{ row.avg_rating?.toFixed(1) }}</template>
            </el-table-column>
          </el-table>
        </el-card>
      </el-col>
    </el-row>

    <el-card style="margin-top: 20px">
      <template #header>
        <div class="card-header">
          <span>导出报表</span>
        </div>
      </template>
      <div class="export-bar">
        <el-button type="primary" @click="handleExportExcel">
          <el-icon><Download /></el-icon>导出 Excel
        </el-button>
        <el-button type="success" @click="handleExportPDF">
          <el-icon><Download /></el-icon>导出 PDF
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, onUnmounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  getProjectStatistics, getTranslatorStatistics,
  getRevenueTrend, getLanguagePairStatistics,
  exportExcel, exportPDF
} from '@/api/statistics'
import * as echarts from 'echarts'
import { saveAs } from 'file-saver'

const chartRef = ref<HTMLElement>()
let chart: echarts.ECharts | null = null

const projectStats = ref<any>({})
const translatorStats = ref<any[]>([])
const languagePairStats = ref<any[]>([])

async function loadData() {
  try {
    const [projStats, transStats, revenueData, lpStats] = await Promise.all([
      getProjectStatistics(),
      getTranslatorStatistics(),
      getRevenueTrend({ months: 6 }),
      getLanguagePairStatistics()
    ])
    projectStats.value = projStats || {}
    translatorStats.value = transStats || []
    languagePairStats.value = lpStats || []

    renderChart(revenueData || [])
  } catch (e) {
    console.error(e)
  }
}

function renderChart(data: any[]) {
  if (!chartRef.value) return
  chart = echarts.init(chartRef.value)
  chart.setOption({
    tooltip: { trigger: 'axis' },
    xAxis: { type: 'category', data: data.map(d => d.month) },
    yAxis: { type: 'value' },
    series: [{
      name: '收入',
      type: 'bar',
      data: data.map(d => d.revenue),
      itemStyle: { color: '#67C23A' }
    }]
  })
}

async function handleExportExcel() {
  try {
    const blob = await exportExcel() as any
    saveAs(blob, `statistics_${Date.now()}.xlsx`)
    ElMessage.success('导出成功')
  } catch (e) {
    ElMessage.error('导出失败')
  }
}

async function handleExportPDF() {
  try {
    const blob = await exportPDF() as any
    saveAs(blob, `statistics_${Date.now()}.pdf`)
    ElMessage.success('导出成功')
  } catch (e) {
    ElMessage.error('导出失败')
  }
}

function handleResize() {
  chart?.resize()
}

onMounted(() => {
  loadData()
  window.addEventListener('resize', handleResize)
})

onUnmounted(() => {
  chart?.dispose()
  window.removeEventListener('resize', handleResize)
})
</script>

<style lang="scss" scoped>
.page-container {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .stat-grid {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 16px;
  }

  .stat-item {
    padding: 16px;
    background: #f5f7fa;
    border-radius: 4px;
    text-align: center;

    .stat-label {
      font-size: 13px;
      color: #909399;
      margin-bottom: 8px;
    }

    .stat-value {
      font-size: 20px;
      font-weight: bold;
      color: #303133;

      &.primary { color: #409EFF; }
      &.success { color: #67C23A; }
      &.warning { color: #E6A23C; }
    }
  }

  .export-bar {
    display: flex;
    gap: 16px;
  }
}
</style>
