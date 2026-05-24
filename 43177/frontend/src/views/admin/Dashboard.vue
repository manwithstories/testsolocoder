<template>
  <div class="dashboard-page">
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-icon" style="background: #409eff;">
            <el-icon :size="24"><User /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.total_users || 0 }}</div>
            <div class="stat-label">总用户数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-icon" style="background: #67c23a;">
            <el-icon :size="24"><Document /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.total_orders || 0 }}</div>
            <div class="stat-label">总工单数</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-icon" style="background: #e6a23c;">
            <el-icon :size="24"><Money /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">¥{{ stats.total_revenue || 0 }}</div>
            <div class="stat-label">总收入</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-icon" style="background: #f56c6c;">
            <el-icon :size="24"><Warning /></el-icon>
          </div>
          <div class="stat-info">
            <div class="stat-value">{{ stats.pending_technicians || 0 }}</div>
            <div class="stat-label">待审核技师</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" class="stats-row">
      <el-col :span="8">
        <el-card class="stat-card small">
          <div class="stat-info">
            <div class="stat-value">{{ stats.pending_orders || 0 }}</div>
            <div class="stat-label">待处理工单</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="stat-card small">
          <div class="stat-info">
            <div class="stat-value">{{ stats.pending_withdraw || 0 }}</div>
            <div class="stat-label">待处理提现</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card class="stat-card small">
          <div class="stat-info">
            <div class="stat-value">{{ stats.low_stock_parts || 0 }}</div>
            <div class="stat-label">库存预警配件</div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-card class="mt-20">
      <template #header>
        <div class="card-header">快捷操作</div>
      </template>
      <div class="quick-actions">
        <el-button type="primary" @click="router.push('/admin/technicians/verify')">
          技师审核
        </el-button>
        <el-button type="success" @click="router.push('/admin/orders')">
          工单管理
        </el-button>
        <el-button type="warning" @click="router.push('/admin/refunds')">
          退款审核
        </el-button>
        <el-button type="danger" @click="router.push('/admin/reviews')">
          差评处理
        </el-button>
        <el-button type="info" @click="router.push('/admin/withdraws')">
          提现审核
        </el-button>
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { User, Document, Money, Warning } from '@element-plus/icons-vue'
import { adminApi } from '@/api/admin'

const router = useRouter()
const stats = ref<any>({})

onMounted(() => {
  loadStats()
})

async function loadStats() {
  try {
    const res = await adminApi.getDashboard()
    stats.value = res.data || {}
  } catch (error) {
    console.error('Failed to load stats:', error)
  }
}
</script>

<style scoped>
.dashboard-page {
  padding: 0;
}

.stats-row {
  margin-bottom: 20px;
}

.stat-card {
  display: flex;
  align-items: center;
}

.stat-card.small {
  justify-content: center;
}

.stat-icon {
  width: 60px;
  height: 60px;
  border-radius: 8px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
  margin-right: 15px;
}

.stat-info {
  flex: 1;
}

.stat-value {
  font-size: 24px;
  font-weight: 600;
  color: #303133;
}

.stat-label {
  font-size: 14px;
  color: #909399;
  margin-top: 5px;
}

.mt-20 {
  margin-top: 20px;
}

.card-header {
  font-weight: 600;
}

.quick-actions {
  display: flex;
  gap: 15px;
  flex-wrap: wrap;
}
</style>
