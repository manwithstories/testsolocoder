<template>
  <div class="seller-stats-page">
    <div class="page-header">
      <h2 class="page-title">数据统计</h2>
    </div>

    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-icon sales">
            <el-icon :size="28"><Money /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">¥{{ totalSales.toFixed(2) }}</div>
            <div class="stat-label">总销售额</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-icon orders">
            <el-icon :size="28"><List /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ orderCount }}</div>
            <div class="stat-label">订单数量</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-icon products">
            <el-icon :size="28"><Goods /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ productCount }}</div>
            <div class="stat-label">商品数量</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-icon customers">
            <el-icon :size="28"><User /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ customerCount }}</div>
            <div class="stat-label">客户数量</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="16">
        <el-card>
          <template #header>
            <span>销售趋势</span>
            <el-date-picker
              v-model="dateRange"
              type="daterange"
              range-separator="至"
              start-placeholder="开始日期"
              end-placeholder="结束日期"
              style="float: right"
              @change="fetchSalesData"
            />
          </template>
          <v-chart :option="salesChartOption" style="height: 350px" autoresize />
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header>
            <span>商品分类销售占比</span>
          </template>
          <v-chart :option="categoryChartOption" style="height: 350px" autoresize />
        </el-card>
      </el-col>
    </el-row>

    <el-card style="margin-top: 20px">
      <template #header>
        <span>最近订单</span>
        <el-button type="primary" link style="float: right" @click="$router.push('/seller/orders')">
          查看全部
        </el-button>
      </template>
      <el-table :data="recentOrders" stripe>
        <el-table-column prop="orderNo" label="订单号" width="180" />
        <el-table-column prop="productTitle" label="商品名称" min-width="200" show-overflow-tooltip />
        <el-table-column prop="finalPrice" label="金额" width="100">
          <template #default="{ row }">
            <span class="price-text">¥{{ row.finalPrice.toFixed(2) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag size="small">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
      </el-table>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { orderApi } from '@/api'
import { OrderStatusText } from '@/types'
import { use } from 'echarts/core'
import { CanvasRenderer } from 'echarts/renderers'
import { LineChart, PieChart } from 'echarts/charts'
import {
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
} from 'echarts/components'
import VChart from 'vue-echarts'

use([
  CanvasRenderer,
  LineChart,
  PieChart,
  TitleComponent,
  TooltipComponent,
  LegendComponent,
  GridComponent
])

const totalSales = ref(0)
const orderCount = ref(0)
const productCount = ref(0)
const customerCount = ref(0)
const recentOrders = ref<any[]>([])

const dateRange = ref<[Date, Date] | null>(null)

const salesChartOption = computed(() => ({
  tooltip: {
    trigger: 'axis'
  },
  legend: {
    data: ['销售额', '订单数']
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
    data: ['1月', '2月', '3月', '4月', '5月', '6月', '7月']
  },
  yAxis: [
    {
      type: 'value',
      name: '销售额(元)',
      position: 'left'
    },
    {
      type: 'value',
      name: '订单数',
      position: 'right'
    }
  ],
  series: [
    {
      name: '销售额',
      type: 'line',
      smooth: true,
      data: [12000, 19000, 15000, 28000, 22000, 35000, 42000]
    },
    {
      name: '订单数',
      type: 'line',
      yAxisIndex: 1,
      smooth: true,
      data: [24, 38, 30, 56, 44, 70, 84]
    }
  ]
}))

const categoryChartOption = computed(() => ({
  tooltip: {
    trigger: 'item'
  },
  legend: {
    orient: 'vertical',
    left: 'left'
  },
  series: [
    {
      name: '销售占比',
      type: 'pie',
      radius: ['40%', '70%'],
      avoidLabelOverlap: false,
      itemStyle: {
        borderRadius: 10,
        borderColor: '#fff',
        borderWidth: 2
      },
      label: {
        show: false
      },
      emphasis: {
        label: {
          show: true,
          fontSize: 16,
          fontWeight: 'bold'
        }
      },
      data: [
        { value: 1048, name: '手机' },
        { value: 735, name: '电脑' },
        { value: 580, name: '耳机' },
        { value: 484, name: '相机' },
        { value: 300, name: '其他' }
      ]
    }
  ]
}))

function getStatusText(status: number): string {
  return OrderStatusText[status] || '未知'
}

function formatTime(time: string): string {
  return new Date(time).toLocaleString()
}

async function fetchStats() {
  try {
    const res = await orderApi.getSellerOrders({
      page: 1,
      pageSize: 10
    })
    recentOrders.value = res.data
    orderCount.value = res.pagination.total

    let total = 0
    const customers = new Set<number>()
    res.data.forEach((order: any) => {
      if (order.status === 5) {
        total += order.finalPrice
      }
      customers.add(order.buyerId)
    })
    totalSales.value = total
    customerCount.value = customers.size
  } catch (error) {
    console.error('Failed to fetch stats:', error)
  }
}

function fetchSalesData() {
}

onMounted(() => {
  fetchStats()
})
</script>

<style lang="scss" scoped>
.seller-stats-page {
  .stat-card {
    display: flex;
    align-items: center;
    gap: 20px;

    .stat-icon {
      width: 60px;
      height: 60px;
      border-radius: 12px;
      display: flex;
      align-items: center;
      justify-content: center;
      color: #fff;

      &.sales {
        background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
      }

      &.orders {
        background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
      }

      &.products {
        background: linear-gradient(135deg, #4facfe 0%, #00f2fe 100%);
      }

      &.customers {
        background: linear-gradient(135deg, #43e97b 0%, #38f9d7 100%);
      }
    }

    .stat-info {
      flex: 1;

      .stat-value {
        font-size: 24px;
        font-weight: 600;
        margin-bottom: 4px;
      }

      .stat-label {
        color: var(--text-lighter-color);
        font-size: 13px;
      }
    }
  }
}
</style>
