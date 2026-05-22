<template>
  <div class="dashboard">
    <el-row :gutter="20" class="stats-cards">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-item">
            <el-icon size="40" color="#409eff"><Goods /></el-icon>
            <div class="stat-info">
              <span class="label">拍卖品总数</span>
              <span class="value">{{ stats?.total_auctions || 0 }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-item">
            <el-icon size="40" color="#67c23a"><Cpu /></el-icon>
            <div class="stat-info">
              <span class="label">出价总数</span>
              <span class="value">{{ stats?.total_bids || 0 }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-item">
            <el-icon size="40" color="#e6a23c"><List /></el-icon>
            <div class="stat-info">
              <span class="label">订单总数</span>
              <span class="value">{{ stats?.total_orders || 0 }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-item">
            <el-icon size="40" color="#f56c6c"><Money /></el-icon>
            <div class="stat-info">
              <span class="label">成交总额</span>
              <span class="value">¥{{ (stats?.total_amount || 0).toFixed(2) }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="stats-cards">
      <el-col :span="8">
        <el-card>
          <template #header>用户统计</template>
          <div class="user-stats">
            <div class="stat-row">
              <span>活跃用户</span>
              <span class="highlight">{{ stats?.active_users || 0 }}</span>
            </div>
            <div class="stat-row">
              <span>新增用户</span>
              <span class="highlight">{{ stats?.new_users || 0 }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header>成交统计</template>
          <div class="user-stats">
            <div class="stat-row">
              <span>成交率</span>
              <span class="highlight">{{ (stats?.success_rate || 0).toFixed(2) }}%</span>
            </div>
            <div class="stat-row">
              <span>平均出价</span>
              <span class="highlight">¥{{ (stats?.average_bid_amount || 0).toFixed(2) }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header>数据导出</template>
          <div class="export-buttons">
            <el-button type="primary" @click="exportOrders">导出订单数据</el-button>
            <el-button type="success" @click="exportBids">导出出价数据</el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import type { Statistics } from '@/types'
import { adminApi } from '@/api'

const stats = ref<Statistics | null>(null)

const fetchStats = async () => {
  try {
    const res = await adminApi.getStatistics({})
    stats.value = res
  } catch (e) {}
}

const exportOrders = async () => {
  try {
    const blob = await adminApi.exportOrders({})
    const url = window.URL.createObjectURL(new Blob([blob]))
    const link = document.createElement('a')
    link.href = url
    link.download = 'orders.csv'
    link.click()
    window.URL.revokeObjectURL(url)
  } catch (e) {}
}

const exportBids = async () => {
  try {
    const blob = await adminApi.exportBids({})
    const url = window.URL.createObjectURL(new Blob([blob]))
    const link = document.createElement('a')
    link.href = url
    link.download = 'bids.csv'
    link.click()
    window.URL.revokeObjectURL(url)
  } catch (e) {}
}

onMounted(() => {
  fetchStats()
})
</script>

<style scoped>
.stats-cards {
  margin-bottom: 20px;
}

.stat-card {
  margin-bottom: 20px;
}

.stat-item {
  display: flex;
  align-items: center;
  gap: 15px;
}

.stat-info {
  display: flex;
  flex-direction: column;
  gap: 5px;
}

.stat-info .label {
  color: #909399;
  font-size: 14px;
}

.stat-info .value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}

.user-stats {
  display: flex;
  flex-direction: column;
  gap: 15px;
}

.stat-row {
  display: flex;
  justify-content: space-between;
  font-size: 14px;
  color: #606266;
}

.highlight {
  font-weight: bold;
  color: #409eff;
}

.export-buttons {
  display: flex;
  flex-direction: column;
  gap: 10px;
}
</style>
