<template>
  <div class="dashboard">
    <div class="page-header">
      <h2>管理后台</h2>
    </div>
    
    <el-row :gutter="20" class="stats-row">
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="icon" style="background: rgba(64, 158, 255, 0.1);">
              <el-icon :size="24" color="#409eff"><User /></el-icon>
            </div>
            <div class="info">
              <div class="value">{{ stats.totalUsers || 0 }}</div>
              <div class="label">用户总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="icon" style="background: rgba(103, 194, 58, 0.1);">
              <el-icon :size="24" color="#67c23a"><Headset /></el-icon>
            </div>
            <div class="info">
              <div class="value">{{ stats.totalWorks || 0 }}</div>
              <div class="label">作品总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="icon" style="background: rgba(230, 162, 60, 0.1);">
              <el-icon :size="24" color="#e6a23c"><Calendar /></el-icon>
            </div>
            <div class="info">
              <div class="value">{{ stats.totalEvents || 0 }}</div>
              <div class="label">演出总数</div>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card>
          <div class="stat-item">
            <div class="icon" style="background: rgba(245, 108, 108, 0.1);">
              <el-icon :size="24" color="#f56c6c"><Wallet /></el-icon>
            </div>
            <div class="info">
              <div class="value">¥{{ stats.pendingWithdraw || 0 }}</div>
              <div class="label">待审核提现</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <el-row :gutter="20" class="content-row">
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>待审核提现</span>
          </template>
          <el-table :data="pendingWithdraw" style="width: 100%">
            <el-table-column prop="amount" label="金额" width="100">
              <template #default="{ row }">
                ¥{{ row.amount.toFixed(2) }}
              </template>
            </el-table-column>
            <el-table-column prop="method" label="方式" width="80" />
            <el-table-column prop="created_at" label="时间" />
          </el-table>
        </el-card>
      </el-col>
      <el-col :span="12">
        <el-card>
          <template #header>
            <span>最新用户</span>
          </template>
          <el-table :data="latestUsers" style="width: 100%">
            <el-table-column prop="username" label="用户名" />
            <el-table-column prop="role" label="角色" width="100">
              <template #default="{ row }">
                <el-tag size="small">{{ getRoleText(row.role) }}</el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="注册时间" />
          </el-table>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { userApi } from '@/api/auth'
import { revenueApi } from '@/api/revenue'

const stats = ref({
  totalUsers: 0,
  totalWorks: 0,
  totalEvents: 0,
  pendingWithdraw: 0
})

const pendingWithdraw = ref<any[]>([])
const latestUsers = ref<any[]>([])

onMounted(() => {
  loadData()
})

async function loadData() {
  try {
    const [users, withdraw] = await Promise.all([
      userApi.list({ page: 1, page_size: 5 }),
      revenueApi.getAdminWithdrawList({ page: 1, page_size: 5, status: 0 })
    ])
    
    latestUsers.value = users.list
    pendingWithdraw.value = withdraw.list
    stats.value.totalUsers = users.total
    stats.value.pendingWithdraw = withdraw.total
  } catch (e) {
    console.error(e)
  }
}

function getRoleText(role: string) {
  const texts: Record<string, string> = {
    admin: '管理员',
    artist: '音乐人',
    label: '厂牌',
    fan: '乐迷'
  }
  return texts[role] || role
}
</script>

<style scoped lang="scss">
.dashboard {
  .stats-row {
    margin-bottom: 24px;
    
    .stat-item {
      display: flex;
      align-items: center;
      gap: 16px;
      
      .icon {
        width: 48px;
        height: 48px;
        border-radius: 8px;
        display: flex;
        align-items: center;
        justify-content: center;
      }
      
      .info {
        .value {
          font-size: 20px;
          font-weight: 600;
        }
        
        .label {
          font-size: 13px;
          color: var(--text-light);
        }
      }
    }
  }
  
  .content-row {
    :deep(.el-card) {
      margin-bottom: 0;
    }
  }
}
</style>
