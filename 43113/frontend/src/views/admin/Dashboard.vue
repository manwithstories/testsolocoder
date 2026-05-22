<template>
  <div class="dashboard-page">
    <el-row :gutter="20">
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-item">
            <el-icon :size="32" color="#409eff"><Document /></el-icon>
            <div class="stat-info">
              <span class="stat-value">{{ stats?.totalQuestions || 0 }}</span>
              <span class="stat-label">总问题数</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-item">
            <el-icon :size="32" color="#67c23a"><ChatDotRound /></el-icon>
            <div class="stat-info">
              <span class="stat-value">{{ stats?.totalAnswers || 0 }}</span>
              <span class="stat-label">总回答数</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card">
          <div class="stat-item">
            <el-icon :size="32" color="#e6a23c"><User /></el-icon>
            <div class="stat-info">
              <span class="stat-value">{{ stats?.totalUsers || 0 }}</span>
              <span class="stat-label">总用户数</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card class="stat-card warning">
          <div class="stat-item">
            <el-icon :size="32" color="#f56c6c"><Warning /></el-icon>
            <div class="stat-info">
              <span class="stat-value">{{ stats?.pendingAuditCount || 0 }}</span>
              <span class="stat-label">待审核</span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="8">
        <el-card>
          <template #header>今日新增</template>
          <div class="today-stats">
            <div class="today-item">
              <span class="label">新增问题</span>
              <span class="value">{{ stats?.todayNewQuestions || 0 }}</span>
            </div>
            <div class="today-item">
              <span class="label">新增回答</span>
              <span class="value">{{ stats?.todayNewAnswers || 0 }}</span>
            </div>
            <div class="today-item">
              <span class="label">新增用户</span>
              <span class="value">{{ stats?.todayNewUsers || 0 }}</span>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="16">
        <el-card>
          <template #header>快捷操作</template>
          <div class="quick-actions">
            <router-link to="/admin/audit">
              <el-button type="primary" size="large">
                <el-icon><DocumentChecked /></el-icon>
                内容审核
              </el-button>
            </router-link>
            <router-link to="/admin/reports">
              <el-button type="warning" size="large">
                <el-icon><Warning /></el-icon>
                举报处理
              </el-button>
            </router-link>
            <router-link to="/admin/users">
              <el-button type="success" size="large">
                <el-icon><User /></el-icon>
                用户管理
              </el-button>
            </router-link>
            <router-link to="/admin/stats">
              <el-button type="info" size="large">
                <el-icon><TrendCharts /></el-icon>
                统计报表
              </el-button>
            </router-link>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { statsApi } from '@/api'
import type { DashboardStats } from '@/types'

const stats = ref<DashboardStats | null>(null)

const fetchStats = async () => {
  try {
    const res = await statsApi.getDashboardStats()
    stats.value = res.data || null
  } catch (e) {
    console.error(e)
  }
}

onMounted(() => {
  fetchStats()
})
</script>

<style scoped lang="scss">
.dashboard-page {
  .stat-card {
    .stat-item {
      display: flex;
      align-items: center;
      gap: 16px;

      .stat-info {
        display: flex;
        flex-direction: column;

        .stat-value {
          font-size: 28px;
          font-weight: bold;
        }

        .stat-label {
          font-size: 14px;
          color: #909399;
        }
      }
    }

    &.warning {
      .stat-value {
        color: #f56c6c;
      }
    }
  }

  .today-stats {
    .today-item {
      display: flex;
      justify-content: space-between;
      align-items: center;
      padding: 12px 0;
      border-bottom: 1px solid #e4e7ed;

      &:last-child {
        border-bottom: none;
      }

      .label {
        color: #606266;
      }

      .value {
        font-size: 18px;
        font-weight: bold;
        color: #409eff;
      }
    }
  }

  .quick-actions {
    display: flex;
    gap: 16px;
    flex-wrap: wrap;

    .el-button {
      min-width: 120px;
    }
  }
}
</style>
