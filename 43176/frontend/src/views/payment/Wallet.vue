<template>
  <div class="wallet-container">
    <el-card class="balance-card">
      <div class="balance-info">
        <div class="balance-label">账户余额</div>
        <div class="balance-amount">¥{{ balance.toFixed(2) }}</div>
      </div>
      <div class="balance-actions">
        <el-button type="primary" size="large" @click="showDepositDialog = true">
          充值
        </el-button>
        <el-button size="large" @click="showWithdrawDialog = true">
          提现
        </el-button>
      </div>
    </el-card>

    <el-card class="history-card">
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
          <div class="tx-info">
            <div class="tx-type">{{ getTypeLabel(tx.type) }}</div>
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
    </el-card>

    <el-dialog v-model="showDepositDialog" title="充值" width="400px">
      <el-form :model="depositForm" label-width="80px">
        <el-form-item label="金额">
          <el-input-number v-model="depositForm.amount" :min="1" :max="10000" />
        </el-form-item>
        <el-form-item label="支付方式">
          <el-radio-group v-model="depositForm.payment_method">
            <el-radio value="alipay">支付宝</el-radio>
            <el-radio value="wechat">微信</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showDepositDialog = false">取消</el-button>
        <el-button type="primary" @click="handleDeposit" :loading="depositLoading">
          确认充值
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showWithdrawDialog" title="提现" width="400px">
      <el-form :model="withdrawForm" label-width="80px">
        <el-form-item label="金额">
          <el-input-number v-model="withdrawForm.amount" :min="1" :max="balance" />
        </el-form-item>
        <el-form-item label="提现方式">
          <el-radio-group v-model="withdrawForm.account_type">
            <el-radio value="alipay">支付宝</el-radio>
            <el-radio value="wechat">微信</el-radio>
            <el-radio value="bank">银行卡</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="账号">
          <el-input v-model="withdrawForm.account_no" placeholder="请输入账号" />
        </el-form-item>
        <el-form-item label="姓名">
          <el-input v-model="withdrawForm.account_name" placeholder="请输入姓名" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showWithdrawDialog = false">取消</el-button>
        <el-button type="primary" @click="handleWithdraw" :loading="withdrawLoading">
          确认提现
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { List } from '@element-plus/icons-vue'
import { paymentApi } from '@/api'
import { useUserStore } from '@/stores/user'
import dayjs from 'dayjs'
import type { Transaction, TransactionType, TransactionStatus } from '@/types'

const userStore = useUserStore()
const loading = ref(false)
const depositLoading = ref(false)
const withdrawLoading = ref(false)
const showDepositDialog = ref(false)
const showWithdrawDialog = ref(false)
const activeTab = ref('all')
const transactions = ref<Transaction[]>([])

const balance = computed(() => userStore.userInfo?.balance || 0)

const depositForm = reactive({
  amount: 100,
  payment_method: 'alipay'
})

const withdrawForm = reactive({
  amount: 100,
  account_type: 'alipay',
  account_no: '',
  account_name: ''
})

const filteredTransactions = computed(() => {
  if (activeTab.value === 'all') return transactions.value
  if (activeTab.value === 'deposit') return transactions.value.filter(t => t.type === 'deposit')
  if (activeTab.value === 'withdraw') return transactions.value.filter(t => t.type === 'withdraw')
  if (activeTab.value === 'income') return transactions.value.filter(t => t.type === 'settlement')
  return transactions.value
})

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
    const res = await paymentApi.history({ page_size: 50 })
    if (res.code === 200) {
      transactions.value = res.data.items
    }
  } catch (error) {
    console.error('Failed to fetch transactions:', error)
  } finally {
    loading.value = false
  }
}

const handleDeposit = async () => {
  depositLoading.value = true
  try {
    const res = await paymentApi.deposit(depositForm)
    if (res.code === 200) {
      ElMessage.success('充值成功')
      userStore.updateUserInfo({ balance: balance.value + depositForm.amount })
      showDepositDialog.value = false
      fetchTransactions()
    }
  } catch (error) {
    console.error('Deposit failed:', error)
  } finally {
    depositLoading.value = false
  }
}

const handleWithdraw = async () => {
  if (!withdrawForm.account_no || !withdrawForm.account_name) {
    ElMessage.warning('请填写完整的提现信息')
    return
  }
  withdrawLoading.value = true
  try {
    const res = await paymentApi.withdraw(withdrawForm)
    if (res.code === 200) {
      ElMessage.success('提现申请已提交，等待审核')
      userStore.updateUserInfo({ balance: balance.value - withdrawForm.amount })
      showWithdrawDialog.value = false
      fetchTransactions()
    }
  } catch (error) {
    console.error('Withdraw failed:', error)
  } finally {
    withdrawLoading.value = false
  }
}

onMounted(() => {
  fetchTransactions()
})
</script>

<style lang="scss" scoped>
.wallet-container {
  max-width: 800px;
  margin: 0 auto;
  padding: 20px;
}

.balance-card {
  margin-bottom: 20px;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  border: none;

  :deep(.el-card__body) {
    padding: 30px;
  }

  .balance-info {
    text-align: center;
    margin-bottom: 20px;

    .balance-label {
      color: rgba(255, 255, 255, 0.8);
      font-size: 14px;
      margin-bottom: 8px;
    }

    .balance-amount {
      color: #fff;
      font-size: 48px;
      font-weight: bold;
    }
  }

  .balance-actions {
    display: flex;
    justify-content: center;
    gap: 16px;

    .el-button {
      min-width: 120px;
    }
  }
}

.history-card {
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

      .tx-info {
        flex: 1;

        .tx-type {
          font-weight: 500;
          margin-bottom: 4px;
        }

        .tx-time {
          color: #909399;
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
}
</style>
