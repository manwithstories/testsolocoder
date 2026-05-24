<template>
  <div class="admin-panel-container">
    <el-card>
      <template #header>
        <div class="card-header">
        <el-icon><DataBoard /></el-icon>
        <span>管理后台</span>
      </div>
    </template>

    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
        <div class="stat-icon" style="background: #667eea;">
          <el-icon :size="32"><User /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.totalUsers }}</div>
          <div class="stat-label">总用户数</div>
        </div>
      </el-card>
    </el-col>
      <el-col :span="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-icon" style="background: #67c23a;">
          <el-icon :size="32"><ShoppingCart /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.totalTasks }}</div>
          <div class="stat-label">总任务数</div>
        </div>
      </el-card>
    </el-col>
    <el-col :span="6">
      <el-card class="stat-card" shadow="hover">
        <div class="stat-icon" style="background: #e6a23c;">
        <el-icon :size="32"><Wallet /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">¥{{ stats.totalAmount }}</div>
          <div class="stat-label">交易总额</div>
        </div>
      </el-card>
    </el-col>
    <el-col :span="6">
      <el-card class="stat-card" shadow="hover">
        <div class="stat-icon" style="background: #f56c6c;">
          <el-icon :size="32"><Clock /></el-icon>
        </div>
        <div class="stat-info">
          <div class="stat-value">{{ stats.pendingCouriers }}</div>
          <div class="stat-label">待审核跑腿员</div>
        </div>
      </el-card>
    </el-col>
    </el-row>

    <el-divider />

    <el-row :gutter="20">
      <el-col :span="24">
        <h3 class="section-title">快捷操作</h3>
        <div class="quick-actions">
          <router-link to="/admin/users">
            <el-button type="primary" size="large">
              <el-icon><User /></el-icon>
              用户管理
            </el-button>
          </router-link>
          <router-link to="/admin/couriers">
            <el-button type="success" size="large">
              <el-icon><Postcard /></el-icon>
              跑腿员审核
            </el-button>
          </router-link>
        </div>
      </el-col>
    </el-row>
  </el-card>
</div>
</template>

<script setup lang="ts">
import { reactive, onMounted } from 'vue'
import { DataBoard, User, ShoppingCart, Wallet, Clock, Postcard } from '@element-plus/icons-vue'

const stats = reactive({
  totalUsers: 0,
  totalTasks: 0,
  totalAmount: '0.00',
  pendingCouriers: 0
})

const fetchStats = async () => {
  try {
    stats.totalUsers = 0
    stats.totalTasks = 0
    stats.totalAmount = '0.00'
    stats.pendingCouriers = 0
  } catch (error) {
    console.error('Failed to fetch stats:', error)
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<style lang="scss" scoped>
.admin-panel-container {
  max-width: 1200px;
  margin: 0 auto;
  padding: 20px;
}

.card-header {
  display: flex;
  align-items: center;
  gap: 8px;
  font-weight: 600;
  font-size: 18px;
}

.stat-card {
  :deep(.el-card__body) {
    display: flex;
    align-items: center;
    gap: 16px;
    padding: 20px;
  }

  .stat-icon {
    width: 64px;
    height: 64px;
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
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

.section-title {
  font-size: 16px;
  margin-bottom: 16px;
}

.quick-actions {
  display: flex;
  gap: 16px;

  a {
    text-decoration: none;
  }
}
</style>
