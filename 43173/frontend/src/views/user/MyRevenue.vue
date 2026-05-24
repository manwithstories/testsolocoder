<template>
  <div class="my-revenue">
    <div class="page-header">
      <h2>我的收益</h2>
      <el-button type="primary" @click="goToWithdraw">
        <el-icon><Wallet /></el-icon>
        申请提现
      </el-button>
    </div>
    
    <el-row :gutter="20" class="stats-row">
      <el-col :span="8">
        <el-card>
          <div class="stat-item">
            <div class="label">可提现余额</div>
            <div class="value primary">¥{{ formatAmount(summary.available_balance || 0) }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <div class="stat-item">
            <div class="label">累计收益</div>
            <div class="value">¥{{ formatAmount(summary.total_revenue || 0) }}</div>
          </div>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <div class="stat-item">
            <div class="label">已提现</div>
            <div class="value">¥{{ formatAmount(summary.total_withdraw || 0) }}</div>
          </div>
        </el-card>
      </el-col>
    </el-row>
    
    <div class="section-header">
      <h3>收益明细</h3>
      <div class="filter-bar">
        <el-date-picker
          v-model="dateRange"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          @change="loadRecords"
        />
        <el-button type="primary" @click="exportExcel">
          <el-icon><Download /></el-icon>
          导出
        </el-button>
      </div>
    </div>
    
    <el-table :data="records" v-loading="loading" style="width: 100%">
      <el-table-column prop="type" label="类型" width="100">
        <template #default="{ row }">
          <el-tag :type="getTypeTag(row.type)">{{ getTypeText(row.type) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="amount" label="金额" width="120">
        <template #default="{ row }">
          <span class="amount">¥{{ formatAmount(row.amount) }}</span>
        </template>
      </el-table-column>
      <el-table-column prop="play_count" label="播放次数" width="100" />
      <el-table-column prop="status" label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="row.status === 1 ? 'success' : 'info'">
            {{ row.status === 1 ? '已结算' : '待结算' }}
          </el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="时间" width="180" />
    </el-table>
    
    <div class="pagination">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadRecords"
      />
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { revenueApi } from '@/api/revenue'

const router = useRouter()

const loading = ref(false)
const summary = ref<any>({})
const records = ref<any[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const dateRange = ref<string[]>([])

onMounted(() => {
  loadSummary()
  loadRecords()
})

async function loadSummary() {
  try {
    summary.value = await revenueApi.getSummary()
  } catch (e) {
    console.error(e)
  }
}

async function loadRecords() {
  loading.value = true
  try {
    const params: any = {
      page: page.value,
      page_size: pageSize.value
    }
    
    if (dateRange.value?.length === 2) {
      params.start_date = dateRange.value[0]
      params.end_date = dateRange.value[1]
    }
    
    const res = await revenueApi.getRecords(params)
    records.value = res.list
    total.value = res.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function exportExcel() {
  window.open('/api/v1/export/revenue', '_blank')
}

function goToWithdraw() {
  router.push('/user/withdraw')
}

function formatAmount(amount: number) {
  return amount.toFixed(2)
}

function getTypeTag(type: string): 'primary' | 'success' | 'warning' | 'info' | 'danger' {
  const types: Record<string, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    play: 'info',
    subscription: 'success',
    ticket: 'warning'
  }
  return types[type] || 'info'
}

function getTypeText(type: string) {
  const texts: Record<string, string> = {
    play: '播放',
    subscription: '订阅',
    ticket: '演出'
  }
  return texts[type] || type
}
</script>

<style scoped lang="scss">
.my-revenue {
  .stats-row {
    margin-bottom: 24px;
    
    .stat-item {
      .label {
        font-size: 14px;
        color: var(--text-light);
        margin-bottom: 8px;
      }
      
      .value {
        font-size: 24px;
        font-weight: 600;
        
        &.primary {
          color: var(--primary-color);
        }
      }
    }
  }
  
  .section-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 16px;
    
    h3 {
      margin: 0;
    }
    
    .filter-bar {
      display: flex;
      gap: 12px;
    }
  }
  
  .amount {
    font-weight: 500;
    color: var(--primary-color);
  }
  
  .pagination {
    display: flex;
    justify-content: center;
    margin-top: 20px;
  }
}
</style>
