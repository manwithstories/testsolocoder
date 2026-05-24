<template>
  <div class="dashboard">
    <el-row :gutter="16">
      <el-col v-for="card in cards" :key="card.key" :xs="12" :sm="12" :md="6">
        <el-card shadow="hover" class="stat-card" :body-style="{ padding: '20px' }">
          <div class="stat-content">
            <div class="stat-info">
              <div class="stat-label">{{ card.label }}</div>
              <div class="stat-value" :style="{ color: card.color }">{{ card.value }}</div>
            </div>
            <el-icon class="stat-icon" :size="48" :style="{ color: card.color }">
              <component :is="card.icon" />
            </el-icon>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :xs="24" :md="16">
        <el-card shadow="never">
          <template #header>
            <div class="card-header">
              <span>{{ trendTitle }}</span>
              <el-radio-group v-model="trendRange" size="small" @change="renderTrendChart">
                <el-radio-button value="week">近7天</el-radio-button>
                <el-radio-button value="month">近30天</el-radio-button>
              </el-radio-group>
            </div>
          </template>
          <div ref="trendChartRef" class="chart-container"></div>
        </el-card>
      </el-col>
      <el-col :xs="24" :md="8">
        <el-card shadow="never">
          <template #header>
            <span>快捷入口</span>
          </template>
          <div class="quick-links">
            <div v-for="link in quickLinks" :key="link.path" class="quick-link" @click="goTo(link.path)">
              <el-icon :size="24"><component :is="link.icon" /></el-icon>
              <span>{{ link.title }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :xs="24" :md="12">
        <el-card shadow="never">
          <template #header>
            <span>{{ recentTitle }}</span>
          </template>
          <el-table :data="recentList" size="small">
            <el-table-column prop="id" label="编号" width="80" />
            <el-table-column prop="title" label="标题" min-width="140" show-overflow-tooltip />
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.statusType">{{ row.statusLabel }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="createdAt" label="时间" width="160" />
          </el-table>
        </el-card>
      </el-col>
      <el-col :xs="24" :md="12">
        <el-card shadow="never">
          <template #header>
            <span>待办提醒</span>
          </template>
          <el-timeline>
            <el-timeline-item
              v-for="(item, idx) in todoList"
              :key="idx"
              :timestamp="item.time"
              :type="item.type"
              :hollow="item.hollow"
            >
              {{ item.content }}
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, watch, onMounted, onBeforeUnmount, markRaw } from 'vue'
import { useRouter } from 'vue-router'
import * as echarts from 'echarts'
import {
  Goods,
  List,
  EditPen,
  Tickets,
  DataAnalysis,
  Finished,
  ChatDotRound,
  User
} from '@element-plus/icons-vue'
import { useUserStore } from '@/stores/user'
import type { Role } from '@/types'

const router = useRouter()
const userStore = useUserStore()
const role = computed<Role>(() => userStore.userRole || 'owner')

const trendRange = ref<'week' | 'month'>('week')
const trendChartRef = ref<HTMLDivElement>()
let trendChart: echarts.ECharts | null = null

interface DashboardCard {
  key: string
  label: string
  value: string | number
  color: string
  icon: any
}

const baseCards: Record<Role, DashboardCard[]> = {
  manufacturer: [
    { key: 'products', label: '产品总数', value: 128, color: '#409eff', icon: markRaw(Goods) },
    { key: 'orders', label: '待处理订单', value: 23, color: '#e6a23c', icon: markRaw(List) },
    { key: 'deliveries', label: '待交付', value: 12, color: '#67c23a', icon: markRaw(Finished) },
    { key: 'tickets', label: '待处理工单', value: 5, color: '#f56c6c', icon: markRaw(Tickets) }
  ],
  designer: [
    { key: 'designs', label: '进行中项目', value: 8, color: '#409eff', icon: markRaw(EditPen) },
    { key: 'pending', label: '待审核方案', value: 3, color: '#e6a23c', icon: markRaw(Finished) },
    { key: 'reviews', label: '评审记录', value: 15, color: '#67c23a', icon: markRaw(ChatDotRound) },
    { key: 'tickets', label: '待处理工单', value: 2, color: '#f56c6c', icon: markRaw(Tickets) }
  ],
  owner: [
    { key: 'users', label: '用户总数', value: 256, color: '#409eff', icon: markRaw(User) },
    { key: 'orders', label: '订单总数', value: 892, color: '#e6a23c', icon: markRaw(List) },
    { key: 'designs', label: '设计方案', value: 67, color: '#67c23a', icon: markRaw(EditPen) },
    { key: 'amount', label: '本月销售额', value: '¥586,420', color: '#9b59b6', icon: markRaw(DataAnalysis) }
  ]
}

const cards = computed<DashboardCard[]>(() => baseCards[role.value])

const trendTitle = computed(() => {
  const map: Record<Role, string> = {
    manufacturer: '订单趋势',
    designer: '设计项目趋势',
    owner: '销售趋势'
  }
  return map[role.value]
})

const recentTitle = computed(() => {
  const map: Record<Role, string> = {
    manufacturer: '最近订单',
    designer: '最近评审',
    owner: '最近动态'
  }
  return map[role.value]
})

const recentList = ref<{ id: number; title: string; statusLabel: string; statusType: string; createdAt: string }[]>([
  { id: 1001, title: '订单 #A20260524001', statusLabel: '进行中', statusType: 'primary', createdAt: '2026-05-24 10:32' },
  { id: 1002, title: '订单 #A20260523002', statusLabel: '待处理', statusType: 'warning', createdAt: '2026-05-23 15:20' },
  { id: 1003, title: '订单 #A20260522003', statusLabel: '已完成', statusType: 'success', createdAt: '2026-05-22 09:10' },
  { id: 1004, title: '订单 #A20260521004', statusLabel: '已完成', statusType: 'success', createdAt: '2026-05-21 16:45' }
])

const todoList = ref<{ content: string; time: string; type: string; hollow?: boolean }[]>([
  { content: '有 3 个新订单需要处理', time: '今天 10:30', type: 'primary' },
  { content: '张先生的方案已提交审核', time: '今天 09:15', type: 'success' },
  { content: '有 2 条用户评价需要回复', time: '昨天 18:20', type: 'warning' },
  { content: '客户 #1024 的工单已解决', time: '昨天 14:00', type: 'info', hollow: true }
])

const quickLinksMap: Record<Role, { path: string; title: string; icon: any }[]> = {
  manufacturer: [
    { path: '/products', title: '产品管理', icon: markRaw(Goods) },
    { path: '/orders', title: '订单管理', icon: markRaw(List) },
    { path: '/deliveries', title: '交付管理', icon: markRaw(Finished) }
  ],
  designer: [
    { path: '/designs', title: '设计项目', icon: markRaw(EditPen) },
    { path: '/deliveries', title: '交付管理', icon: markRaw(Finished) },
    { path: '/reviews', title: '评审记录', icon: markRaw(ChatDotRound) }
  ],
  owner: [
    { path: '/users', title: '用户管理', icon: markRaw(User) },
    { path: '/orders', title: '订单管理', icon: markRaw(List) },
    { path: '/statistics', title: '数据统计', icon: markRaw(DataAnalysis) }
  ]
}

const quickLinks = computed(() => quickLinksMap[role.value])

function goTo(path: string) {
  router.push(path)
}

function mockTrendData() {
  const days = trendRange.value === 'week' ? 7 : 30
  const dates: string[] = []
  const values: number[] = []
  const now = new Date()
  for (let i = days - 1; i >= 0; i--) {
    const d = new Date(now)
    d.setDate(d.getDate() - i)
    dates.push(`${d.getMonth() + 1}/${d.getDate()}`)
    values.push(Math.floor(Math.random() * 100) + 20)
  }
  return { dates, values }
}

function renderTrendChart() {
  if (!trendChartRef.value) return
  if (!trendChart) trendChart = echarts.init(trendChartRef.value)
  const { dates, values } = mockTrendData()
  trendChart.setOption({
    tooltip: { trigger: 'axis' },
    grid: { left: 40, right: 20, top: 20, bottom: 30 },
    xAxis: { type: 'category', data: dates, boundaryGap: false },
    yAxis: { type: 'value' },
    series: [
      {
        type: 'line',
        smooth: true,
        data: values,
        areaStyle: {
          color: new echarts.graphic.LinearGradient(0, 0, 0, 1, [
            { offset: 0, color: 'rgba(64,158,255,0.4)' },
            { offset: 1, color: 'rgba(64,158,255,0.05)' }
          ])
        },
        itemStyle: { color: '#409eff' }
      }
    ]
  })
}

function handleResize() {
  trendChart?.resize()
}

watch(
  () => role.value,
  () => {
    renderTrendChart()
  }
)

onMounted(() => {
  renderTrendChart()
  window.addEventListener('resize', handleResize)
})

onBeforeUnmount(() => {
  window.removeEventListener('resize', handleResize)
  trendChart?.dispose()
})
</script>

<style lang="scss" scoped>
.dashboard {
  .stat-card {
    .stat-content {
      display: flex;
      justify-content: space-between;
      align-items: center;
    }
    .stat-label {
      color: #909399;
      font-size: 14px;
    }
    .stat-value {
      font-size: 28px;
      font-weight: bold;
      margin-top: 8px;
    }
    .stat-icon {
      opacity: 0.85;
    }
  }
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .chart-container {
    width: 100%;
    height: 320px;
  }
  .quick-links {
    display: grid;
    grid-template-columns: repeat(3, 1fr);
    gap: 12px;
    .quick-link {
      display: flex;
      flex-direction: column;
      align-items: center;
      justify-content: center;
      padding: 20px 0;
      background: #f5f7fa;
      border-radius: 6px;
      color: #606266;
      cursor: pointer;
      transition: all 0.2s;
      span {
        margin-top: 8px;
        font-size: 13px;
      }
      &:hover {
        background: #409eff;
        color: #fff;
      }
    }
  }
}
</style>
