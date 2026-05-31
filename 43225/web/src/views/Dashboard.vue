<template>
  <div class="dashboard">
    <el-row :gutter="16" class="stat-cards">
      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card">
          <div class="stat-label">可用船只</div>
          <div class="stat-value">{{ stats.totalShips }}</div>
        </div>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card">
          <div class="stat-label">进行中订单</div>
          <div class="stat-value">{{ stats.activeRentals }}</div>
        </div>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card">
          <div class="stat-label">本月收入</div>
          <div class="stat-value">¥{{ stats.monthlyIncome }}</div>
        </div>
      </el-col>
      <el-col :xs="24" :sm="12" :md="6">
        <div class="stat-card">
          <div class="stat-label">待处理</div>
          <div class="stat-value">{{ stats.pendingTasks }}</div>
        </div>
      </el-col>
    </el-row>

    <el-row :gutter="16" style="margin-top: 16px">
      <el-col :xs="24" :md="16">
        <div class="card-container">
          <div class="card-header">
            <span>最近订单</span>
            <el-link type="primary" @click="$router.push('/rentals')">查看全部</el-link>
          </div>
          <el-table :data="recentRentals" style="width: 100%">
            <el-table-column prop="id" label="订单号" width="180" show-overflow-tooltip />
            <el-table-column prop="ship.name" label="船只" />
            <el-table-column prop="start_date" label="开始日期" width="120">
              <template #default="{ row }">
                {{ formatDate(row.start_date) }}
              </template>
            </el-table-column>
            <el-table-column prop="total_amount" label="金额" width="120">
              <template #default="{ row }">
                {{ row.currency }} {{ row.total_amount }}
              </template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getStatusType(row.status)" size="small">
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-col>
      <el-col :xs="24" :md="8">
        <div class="card-container">
          <div class="card-header">
            <span>快捷操作</span>
          </div>
          <div class="quick-actions">
            <el-button type="primary" size="large" @click="$router.push('/ship-create')">
              <el-icon><Plus /></el-icon>
              发布船只
            </el-button>
            <el-button size="large" @click="$router.push('/rental-create')">
              <el-icon><Tickets /></el-icon>
              创建租赁
            </el-button>
            <el-button size="large" @click="$router.push('/berth-reservations')">
              <el-icon><Calendar /></el-icon>
              预约泊位
            </el-button>
            <el-button size="large" @click="$router.push('/finance')">
              <el-icon><Money /></el-icon>
              财务结算
            </el-button>
          </div>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import dayjs from 'dayjs'
import { getMyRentalsApi } from '@/api/rental'
import type { Rental } from '@/types/rental'

const stats = ref({
  totalShips: 0,
  activeRentals: 0,
  monthlyIncome: 0,
  pendingTasks: 0
})

const recentRentals = ref<Rental[]>([])

const formatDate = (date: string) => dayjs(date).format('YYYY-MM-DD')

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    confirmed: 'primary',
    active: 'success',
    completed: 'info',
    cancelled: 'danger',
    refunded: 'info'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待确认',
    confirmed: '已确认',
    active: '进行中',
    completed: '已完成',
    cancelled: '已取消',
    refunded: '已退款'
  }
  return map[status] || status
}

onMounted(async () => {
  try {
    const res: any = await getMyRentalsApi()
    recentRentals.value = (res.data || []).slice(0, 5)
    stats.value.activeRentals = recentRentals.value.filter(
      (r) => r.status === 'active' || r.status === 'confirmed'
    ).length
  } catch (error) {
    console.error('Failed to fetch rentals:', error)
  }
})
</script>

<style lang="scss" scoped>
.dashboard {
  .card-container {
    background: #fff;
    border-radius: 8px;
    padding: 20px;
    margin-bottom: 16px;

    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;
      font-weight: 600;
    }
  }

  .quick-actions {
    display: flex;
    flex-direction: column;
    gap: 12px;

    .el-button {
      justify-content: flex-start;

      .el-icon {
        margin-right: 8px;
      }
    }
  }
}
</style>
