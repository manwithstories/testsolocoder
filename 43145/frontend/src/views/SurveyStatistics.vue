<template>
  <div class="page-container">
    <div class="page-header">
      <h1 class="page-title">数据统计 - {{ survey?.title }}</h1>
      <div>
        <el-button :icon="Download" @click="handleExportExcel">导出Excel</el-button>
        <el-button :icon="Printer" @click="handleExportPDF">导出报告</el-button>
      </div>
    </div>

    <div class="statistics-summary">
      <div class="stat-card">
        <div class="stat-card-value">{{ stats?.total_responses || 0 }}</div>
        <div class="stat-card-label">总答卷数</div>
      </div>
      <div class="stat-card">
        <div class="stat-card-value">{{ stats?.completed_count || 0 }}</div>
        <div class="stat-card-label">已完成</div>
      </div>
      <div class="stat-card">
        <div class="stat-card-value">{{ stats?.completion_rate?.toFixed(2) || 0 }}%</div>
        <div class="stat-card-label">完成率</div>
      </div>
      <div class="stat-card">
        <div class="stat-card-value">{{ stats?.avg_duration?.toFixed(0) || 0 }}s</div>
        <div class="stat-card-label">平均用时</div>
      </div>
    </div>

    <div class="card">
      <div class="card-header" style="font-weight: 600;">时间趋势</div>
      <div class="card-body" style="height: 300px;">
        <v-chart class="chart" :option="timeChartOption" autoresize />
      </div>
    </div>

    <div v-for="qs in stats?.question_stats" :key="qs.question_id" class="card">
      <div class="card-header" style="font-weight: 600;">
        {{ qs.question_title }}
        <el-tag size="small" style="margin-left: 8px;">{{ getQuestionTypeText(qs.question_type) }}</el-tag>
        <span style="color: #909399; margin-left: 8px; font-weight: normal;">
          {{ qs.response_count }} 人回答
        </span>
      </div>
      <div class="card-body">
        <div v-if="qs.option_stats && qs.option_stats.length > 0" style="height: 300px;">
          <v-chart class="chart" :option="getOptionStatChart(qs)" autoresize />
        </div>
        <div v-else-if="qs.rating_stats" style="height: 300px;">
          <v-chart class="chart" :option="getRatingChart(qs)" autoresize />
        </div>
        <div v-else-if="qs.text_stats">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="平均长度">{{ qs.text_stats.avg_length.toFixed(0) }} 字符</el-descriptions-item>
            <el-descriptions-item label="最长">{{ qs.text_stats.max_length }} 字符</el-descriptions-item>
            <el-descriptions-item label="最短">{{ qs.text_stats.min_length }} 字符</el-descriptions-item>
          </el-descriptions>
        </div>
        <div v-if="qs.word_cloud && qs.word_cloud.length > 0" style="margin-top: 16px;">
          <h4 style="margin-bottom: 12px;">词云</h4>
          <div class="word-cloud">
            <span
              v-for="(item, index) in qs.word_cloud.slice(0, 30)"
              :key="index"
              :style="{ fontSize: Math.max(12, 24 - index * 0.5) + 'px', color: getWordColor(index) }"
            >
              {{ item.word }}
            </span>
          </div>
        </div>
      </div>
    </div>

    <el-empty v-if="!stats" description="暂无统计数据" />
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Download, Printer } from '@element-plus/icons-vue'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { BarChart, LineChart, PieChart } from 'echarts/charts'
import {
  TitleComponent, TooltipComponent, LegendComponent, GridComponent
} from 'echarts/components'
import VChart from 'vue-echarts'
import { statisticsApi, exportApi } from '@/api/statistics'
import { surveyApi } from '@/api/survey'
import type { Survey, Statistics, QuestionStat } from '@/types'
import { download } from '@/utils/request'

use([
  CanvasRenderer,
  BarChart,
  LineChart,
  PieChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
])

const route = useRoute()
const surveyId = computed(() => Number(route.params.id))

const survey = ref<Survey>()
const stats = ref<Statistics>()

const getQuestionTypeText = (type: string) => {
  const map: Record<string, string> = {
    single_choice: '单选题',
    multi_choice: '多选题',
    fill_in: '填空题',
    rating: '评分题',
    ranking: '排序题',
    matrix: '矩阵题'
  }
  return map[type] || type
}

const getWordColor = (index: number) => {
  const colors = ['#409eff', '#67c23a', '#e6a23c', '#f56c6c', '#909399', '#9b59b6']
  return colors[index % colors.length]
}

const timeChartOption = computed(() => ({
  tooltip: { trigger: 'axis' },
  xAxis: {
    type: 'category',
    data: stats.value?.time_distribution?.map(d => d.date) || []
  },
  yAxis: {
    type: 'value'
  },
  series: [{
    type: 'line',
    data: stats.value?.time_distribution?.map(d => d.count) || [],
    smooth: true,
    areaStyle: {}
  }]
}))

const getOptionStatChart = (qs: QuestionStat) => ({
  tooltip: { trigger: 'axis', axisPointer: { type: 'shadow' } },
  xAxis: {
    type: 'category',
    data: qs.option_stats?.map(o => o.text) || [],
    axisLabel: { interval: 0, rotate: 30 }
  },
  yAxis: {
    type: 'value'
  },
  series: [{
    type: 'bar',
    data: qs.option_stats?.map(o => o.count) || [],
    label: {
      show: true,
      formatter: (params: any) => `${params.value} (${qs.option_stats?.[params.dataIndex]?.percent?.toFixed(1)}%)`
    }
  }]
})

const getRatingChart = (qs: QuestionStat) => ({
  tooltip: { trigger: 'axis' },
  xAxis: {
    type: 'category',
    data: qs.rating_stats ? Object.keys(qs.rating_stats.distribution).map(Number).sort((a, b) => a - b).map(String) : []
  },
  yAxis: {
    type: 'value'
  },
  series: [{
    type: 'bar',
    data: qs.rating_stats ? Object.keys(qs.rating_stats.distribution).map(Number).sort((a, b) => a - b).map(k => qs.rating_stats!.distribution[k]) : []
  }]
})

const loadData = async () => {
  try {
    const [surveyRes, statsRes] = await Promise.all([
      surveyApi.getById(surveyId.value),
      statisticsApi.getStatistics({ survey_id: surveyId.value })
    ])
    survey.value = surveyRes
    stats.value = statsRes
  } catch (error) {
    console.error('Failed to load statistics')
  }
}

const handleExportExcel = () => {
  download(`/export/${surveyId.value}/excel`, `survey_${surveyId.value}.xlsx`)
}

const handleExportPDF = async () => {
  try {
    const report = await exportApi.exportPDF(surveyId.value)
    ElMessage.success('报告已生成')
    console.log('PDF Report:', report)
  } catch (error: any) {
    ElMessage.error(error.message || '导出失败')
  }
}

onMounted(loadData)
</script>

<style scoped>
.chart {
  width: 100%;
  height: 100%;
}

.word-cloud {
  display: flex;
  flex-wrap: wrap;
  gap: 12px;
  align-items: center;
  padding: 16px;
  background: #f5f7fa;
  border-radius: 8px;
}
</style>
