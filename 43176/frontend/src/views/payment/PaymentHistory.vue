<template>
  <div class="payment-history-container">
    <el-card>
      <template #header>
        <div class="card-header">
          <el-icon><List /></el-icon>
          <span>交易记录</span>
        </div>
      </template>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="全部" name="all" />
        <el-tab-pane label="充值" name="deposit" />
        <el-tab-pane label="提现" name="withdraw" />
        <el-tab-pane label="收入" name="income" />
      </el-tabs>

      <div v-if="loading" class="loading">
        <el-skeleton :rows="4" animated />
      </div>

      <div v-else class="history-list">
        <div
          v-for="tx in filteredTransactions"
          :key="tx.id"
          class="history-item"
        >
          <div class="tx-icon">
            <el-icon :size="24"><component :is="getTypeIcon(tx.type)" /></el-icon>
          </div>
          <div class="tx-info">
            <div class="tx-type">{{ getTypeLabel(tx.type) }}</div>
            <div class="tx-desc">{{ tx.description }}</div>
            <div class="tx-time">{{ formatTime(tx.created_at) }}</div>
          </div>
          <div class="tx-amount" :class="getAmountClass(tx.type)">
            {{ getAmountPrefix(tx.type) }}¥{{ tx.amount.toFixed(2) }}
          </div>
          <div class="tx-status">
            <el-tag :type="getStatusTagType(tx.status)" size="small">
              {{ getStatusLabel(tx.status) }}
            </el-tag>
          </div>
        </div>
        <el-empty v-if="filteredTransactions.length === 0" description="暂无记录" />
      </div>

      <div class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="20"
          :total="total"
          layout="total, prev, pager, next"
          @current-change="fetchTransactions"
        />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { List, Money, Wallet, Top, Bottom } from '@element-plus/icons-vue'
import { paymentApi } from '@/api'
import dayjs from 'dayjs'
import type { Transaction, TransactionType, TransactionStatus } from '@/types'

const loading = ref(false)
const activeTab = ref('all')
const transactions = ref<Transaction[]>([])
const total = ref(0)
const currentPage = ref(1)

const filteredTransactions = computed(() => {
  if (activeTab.value === 'all') return transactions.value
  if (activeTab.value === 'deposit') return transactions.value.filter(t => t.type === 'deposit')
  if (activeTab.value === 'withdraw') return transactions.value.filter(t => t.type === 'withdraw')
  if (activeTab.value === 'income') return transactions.value.filter(t => t.type === 'settlement')
  return transactions.value
})

const getTypeIcon = (type: TransactionType) => {
  const icons: Record<TransactionType, any> = {
    deposit: Top,
    withdraw: Bottom,
    payment: Wallet,
    refund: Money,
    settlement: Money,
    service_fee: Wallet
  }
  return icons[type] || Money
}

const getTypeLabel = (type: TransactionType) => {
  const labels: Record<TransactionType, string> = {
    deposit: '充值',
    withdraw: '提现',
    payment: '支付',
    refund: '退款',
    settlement: '收入',
    service_fee: '服务费'
  }
  return labels[type] || type
}

const getAmountPrefix = (type: TransactionType) => {
  if (type === 'deposit' || type === 'settlement' || type === 'refund') return '+'
  if (type === 'withdraw' || type === 'payment' || type === 'service_fee') return '-'
  return ''
}

const getAmountClass = (type: TransactionType) => {
  if (type === 'deposit' || type === 'settlement' || type === 'refund') return 'amount-positive'
  return 'amount-negative'
}

const getStatusLabel = (status: TransactionStatus) => {
  const labels: Record<TransactionStatus, string> = {
    pending: '处理中',
    completed: '已完成',
    failed: '失败',
    cancelled: '已取消'
  }
  return labels[status] || status
}

const getStatusTagType = (status: TransactionStatus): 'primary' | 'success' | 'warning' | 'danger' | 'info' => {
  const types: Record<TransactionStatus, 'primary' | 'success' | 'warning' | 'danger' | 'info'> = {
    pending: 'warning',
    completed: 'success',
    failed: 'danger',
    cancelled: 'info'
  }
  return types[status] || 'info'
}

const formatTime = (time: string) => dayjs(time).format('YYYY-MM-DD HH:mm')

const fetchTransactions = async () => {
  loading.value = true
  try {
    const res = await paymentApi.history({
      page: currentPage.value,
      page_size: 20
    })
    if (res.code === 200) {
      transactions.value = res.data.items
      total.value = res.data.total
    }
  } catch (error) {
    console.error('Failed to fetch transactions:', error)
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchTransactions()
})
</script>

<style lang="scss" scoped>
.payment-history-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;

  .card-header {
    display: flex;
    align-items: center;
    gap: 8px;
    font-weight: 600;
  }

  .history-list {
    .history-item {
      display: flex;
      align-items: center;
      padding: 16px 0;
      border-bottom: 1px solid #ebeef5;

      &:last-child {
        border-bottom: none;
      }

      .tx-icon {
        width: 48px;
        height: 48px;
        background: #f5f7fa;
        border-radius: 50%;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #667eea;
        margin-right: 16px;
      }

      .tx-info {
        flex: 1;

        .tx-type {
          font-weight: 500;
          margin-bottom: 4px;
        }

        .tx-desc {
          color: #909399;
          font-size: 13px;
          margin-bottom: 4px;
        }

        .tx-time {
          color: #c0c4cc;
          font-size: 12px;
        }
      }

      .tx-amount {
        font-size: 18px;
        font-weight: bold;
        margin-right: 16px;

        &.amount-positive {
          color: #67c23a;
        }

        &.amount-negative {
          color: #f56c6c;
        }
      }
    }
  }

  .pagination {
    display: flex;
    justify-content: center;
    margin-top: 20px;
  }
}
</style>
