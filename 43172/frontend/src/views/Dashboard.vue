<template>
  <div class="page-container">
    <div class="card">
      <h2 class="section-title">数据仪表盘</h2>
      <div class="stats-overview">
        <el-row :gutter="16">
          <el-col :span="6">
            <div class="stat-card">
              <div class="stat-icon orders">📦</div>
              <div class="stat-info">
                <p class="stat-label">总订单数</p>
                <p class="stat-value">{{ stats?.total_orders || 0 }}</p>
              </div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-card">
              <div class="stat-icon amount">💰</div>
              <div class="stat-info">
                <p class="stat-label">总交易额</p>
                <p class="stat-value">¥{{ stats?.total_amount?.toFixed(2) || '0.00' }}</p>
              </div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-card">
              <div class="stat-icon users">👥</div>
              <div class="stat-info">
                <p class="stat-label">用户总数</p>
                <p class="stat-value">{{ stats?.total_users || 0 }}</p>
              </div>
            </div>
          </el-col>
          <el-col :span="6">
            <div class="stat-card">
              <div class="stat-icon products">🛍️</div>
              <div class="stat-info">
                <p class="stat-label">在售商品</p>
                <p class="stat-value">{{ stats?.total_products || 0 }}</p>
              </div>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>

    <el-row :gutter="16">
      <el-col :span="14">
        <div class="card">
          <h3 class="section-title">交易趋势</h3>
          <div v-if="trendData.length > 0" ref="trendChartRef" class="chart-container">
            <v-chart class="chart" :option="trendChartOption" autoresize />
          </div>
          <div v-else class="empty-state">
            <p>暂无数据</p>
          </div>
        </div>
      </el-col>
      <el-col :span="10">
        <div class="card">
          <h3 class="section-title">热门品牌排行</h3>
          <div v-if="stats?.brand_rankings?.length" class="brand-ranking">
            <div
              v-for="(brand, index) in stats.brand_rankings"
              :key="brand.brand_name"
              class="brand-item"
            >
              <span class="rank" :class="{ 'top': index < 3 }">{{ index + 1 }}</span>
              <span class="brand-name">{{ brand.brand_name }}</span>
              <span class="brand-count">{{ brand.order_count }} 单</span>
            </div>
          </div>
          <div v-else class="empty-state">
            <p>暂无数据</p>
          </div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="16">
      <el-col :span="12">
        <div class="card">
          <h3 class="section-title">鉴定统计</h3>
          <div class="auth-stats">
            <div class="auth-stat-item">
              <p class="label">总鉴定数</p>
              <p class="value">{{ stats?.auth_stats?.total || 0 }}</p>
            </div>
            <div class="auth-stat-item">
              <p class="label">已完成</p>
              <p class="value">{{ stats?.auth_stats?.completed || 0 }}</p>
            </div>
            <div class="auth-stat-item">
              <p class="label">通过数</p>
              <p class="value success">{{ stats?.auth_stats?.passed || 0 }}</p>
            </div>
            <div class="auth-stat-item">
              <p class="label">通过率</p>
              <p class="value">{{ stats?.auth_stats?.pass_rate?.toFixed(1) || '0.0' }}%</p>
            </div>
          </div>
        </div>
      </el-col>
      <el-col :span="12">
        <div class="card">
          <h3 class="section-title">最近订单</h3>
          <el-table
            v-if="stats?.recent_orders?.length"
            :data="stats.recent_orders.slice(0, 5)"
            size="small"
          >
            <el-table-column label="订单号" prop="order_number" width="150" />
            <el-table-column label="金额" width="100">
              <template #default="{ row }">
                <span class="price-text">¥{{ row.price.toFixed(2) }}</span>
              </template>
            </el-table-column>
            <el-table-column label="状态">
              <template #default="{ row }">
                <el-tag size="small">{{ row.status }}</el-tag>
              </template>
            </el-table-column>
          </el-table>
          <div v-else class="empty-state">
            <p>暂无数据</p>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted, watch } from 'vue'
import { statisticApi } from '@/api/statistic'
import type { DashboardStats } from '@/types'
import VChart from 'vue-echarts'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent
} from 'echarts/components'

use([
  CanvasRenderer,
  LineChart,
  TitleComponent,
  TooltipComponent,
  GridComponent,
  LegendComponent
])

const stats = ref<DashboardStats | null>(null)
const loading = ref(false)

const trendData = computed(() => stats.value?.transaction_trend || [])

const trendChartOption = computed(() => ({
  tooltip: {
    trigger: 'axis'
  },
  legend: {
    data: ['订单数', '交易额']
  },
  grid: {
    left: '3%',
    right: '4%',
    bottom: '3%',
    containLabel: true
  },
  xAxis: {
    type: 'category',
    boundaryGap: false,
    data: trendData.value.map(item => item.date)
  },
  yAxis: [
    {
      type: 'value',
      name: '订单数',
      position: 'left'
    },
    {
      type: 'value',
      name: '交易额',
      position: 'right',
      axisLabel: {
        formatter: '¥{value}'
      }
    }
  ],
  series: [
    {
      name: '订单数',
      type: 'line',
      data: trendData.value.map(item => item.order_count),
      smooth: true,
      areaStyle: {
        opacity: 0.3
      }
    },
    {
      name: '交易额',
      type: 'line',
      yAxisIndex: 1,
      data: trendData.value.map(item => item.total_amount),
      smooth: true,
      areaStyle: {
        opacity: 0.3
      }
    }
  ]
}))

const loadStats = async () => {
  loading.value = true
  try {
    const res = await statisticApi.getDashboardStats(30)
    if (res.code === 200 && res.data) {
      stats.value = res.data
    }
  } catch (error) {
    console.error('Load stats error:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  loadStats()
})
</script>

<style lang="scss" scoped>
.stats-overview {
  margin-bottom: 20px;
}

.stat-card {
  display: flex;
  align-items: center;
  gap: 16px;
  padding: 20px;
  background: #f8f9fa;
  border-radius: 8px;
  
  .stat-icon {
    width: 48px;
    height: 48px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-size: 24px;
    border-radius: 12px;
    
    &.orders { background: #e3f2fd; }
    &.amount { background: #fff3e0; }
    &.users { background: #e8f5e9; }
    &.products { background: #fce4ec; }
  }
  
  .stat-info {
    .stat-label {
      font-size: 14px;
      color: var(--text-secondary);
      margin-bottom: 4px;
    }
    
    .stat-value {
      font-size: 24px;
      font-weight: 600;
      color: var(--text-primary);
    }
  }
}

.chart-container {
  height: 300px;
}

.brand-ranking {
  .brand-item {
    display: flex;
    align-items: center;
    padding: 12px 0;
    border-bottom: 1px solid var(--border-color);
    
    &:last-child {
      border-bottom: none;
    }
    
    .rank {
      width: 24px;
      height: 24px;
      display: flex;
      align-items: center;
      justify-content: center;
      background: #f0f0f0;
      border-radius: 4px;
      font-size: 12px;
      font-weight: 600;
      margin-right: 12px;
      
      &.top {
        background: var(--secondary-color);
        color: #fff;
      }
    }
    
    .brand-name {
      flex: 1;
      font-size: 14px;
    }
    
    .brand-count {
      font-size: 14px;
      color: var(--text-secondary);
    }
  }
}

.auth-stats {
  display: grid;
  grid-template-columns: repeat(2, 1fr);
  gap: 16px;
  
  .auth-stat-item {
    padding: 16px;
    background: #f8f9fa;
    border-radius: 8px;
    text-align: center;
    
    .label {
      font-size: 14px;
      color: var(--text-secondary);
      margin-bottom: 8px;
    }
    
    .value {
      font-size: 24px;
      font-weight: 600;
      color: var(--text-primary);
      
      &.success {
        color: var(--success-color);
      }
    }
  }
}
</style>
