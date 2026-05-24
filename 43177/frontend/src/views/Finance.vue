<template>
  <div class="finance-page">
    <div class="container">
      <div class="page-header">
        <h2 class="page-title">财务管理</h2>
      </div>

      <el-row :gutter="20" class="balance-row">
        <el-col :span="8">
          <el-card class="balance-card">
            <div class="balance-label">账户余额</div>
            <div class="balance-value">¥{{ balance?.balance || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card class="balance-card">
            <div class="balance-label">待提现</div>
            <div class="balance-value">¥{{ balance?.pending_withdraw || 0 }}</div>
          </el-card>
        </el-col>
        <el-col :span="8">
          <el-card class="balance-card">
            <div class="balance-label">可提现</div>
            <div class="balance-value">¥{{ balance?.available_balance || 0 }}</div>
          </el-card>
        </el-col>
      </el-row>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="提现申请" name="withdraw">
          <div class="action-bar">
            <el-button
              type="primary"
              :disabled="!balance || balance.available_balance <= 0"
              @click="showWithdrawDialog = true"
            >
              申请提现
            </el-button>
          </div>

          <el-table :data="withdraws" style="width: 100%">
            <el-table-column prop="request_no" label="申请编号" width="200" />
            <el-table-column prop="amount" label="金额" width="120">
              <template #default="{ row }">¥{{ row.amount }}</template>
            </el-table-column>
            <el-table-column prop="bank_name" label="银行" width="150" />
            <el-table-column prop="bank_account" label="账号" width="200" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getWithdrawStatusType(row.status)" size="small">
                  {{ getWithdrawStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="申请时间" width="180">
              <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="交易记录" name="transactions">
          <el-table :data="transactions" style="width: 100%">
            <el-table-column prop="transaction_no" label="交易号" width="200" />
            <el-table-column prop="type" label="类型" width="100">
              <template #default="{ row }">
                <el-tag :type="getTransactionType(row.type)" size="small">
                  {{ getTransactionTypeText(row.type) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="amount" label="金额" width="120">
              <template #default="{ row }">
                <span :class="row.amount > 0 ? 'income' : 'expense'">
                  {{ row.amount > 0 ? '+' : '' }}¥{{ row.amount }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="balance_after" label="余额" width="120">
              <template #default="{ row }">¥{{ row.balance_after }}</template>
            </el-table-column>
            <el-table-column prop="description" label="描述" />
            <el-table-column prop="created_at" label="时间" width="180">
              <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>
    </div>

    <el-dialog v-model="showWithdrawDialog" title="申请提现" width="500px">
      <el-form :model="withdrawForm" label-width="100px">
        <el-form-item label="提现金额">
          <el-input-number
            v-model="withdrawForm.amount"
            :min="1"
            :max="balance?.available_balance || 0"
            :precision="2"
          />
        </el-form-item>
        <el-form-item label="开户银行">
          <el-input v-model="withdrawForm.bank_name" placeholder="请输入开户银行" />
        </el-form-item>
        <el-form-item label="银行账号">
          <el-input v-model="withdrawForm.bank_account" placeholder="请输入银行账号" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showWithdrawDialog = false">取消</el-button>
        <el-button type="primary" @click="submitWithdraw" :loading="submitting">
          提交申请
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { financeApi } from '@/api/finance'
import type { WithdrawRequest, Transaction } from '@/types'
import dayjs from 'dayjs'

const activeTab = ref('withdraw')
const balance = ref<{ balance: number; pending_withdraw: number; available_balance: number } | null>(null)
const withdraws = ref<WithdrawRequest[]>([])
const transactions = ref<Transaction[]>([])
const showWithdrawDialog = ref(false)
const submitting = ref(false)

const withdrawForm = reactive({
  amount: 0,
  bank_name: '',
  bank_account: ''
})

onMounted(async () => {
  await loadBalance()
  await loadWithdraws()
  await loadTransactions()
})

async function loadBalance() {
  try {
    const res = await financeApi.getBalance()
    balance.value = res.data as any
  } catch (error) {
    console.error('Failed to load balance:', error)
  }
}

async function loadWithdraws() {
  try {
    const res = await financeApi.getWithdraws({ page: 1, page_size: 50 })
    withdraws.value = res.data?.list || []
  } catch (error) {
    console.error('Failed to load withdraws:', error)
  }
}

async function loadTransactions() {
  try {
    const res = await financeApi.getTransactions({ page: 1, page_size: 50 })
    transactions.value = res.data?.list || []
  } catch (error) {
    console.error('Failed to load transactions:', error)
  }
}

function getWithdrawStatusType(status: string): string {
  const statusMap: Record<string, string> = {
    pending: 'warning',
    approved: 'primary',
    rejected: 'danger',
    completed: 'success'
  }
  return statusMap[status] || 'info'
}

function getWithdrawStatusText(status: string): string {
  const statusMap: Record<string, string> = {
    pending: '待审批',
    approved: '已通过',
    rejected: '已拒绝',
    completed: '已打款'
  }
  return statusMap[status] || status
}

function getTransactionType(type: string): string {
  const typeMap: Record<string, string> = {
    income: 'success',
    withdraw: 'warning',
    refund: 'danger'
  }
  return typeMap[type] || 'info'
}

function getTransactionTypeText(type: string): string {
  const typeMap: Record<string, string> = {
    income: '收入',
    withdraw: '提现',
    refund: '退款'
  }
  return typeMap[type] || type
}

function formatTime(time: string): string {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

async function submitWithdraw() {
  if (!withdrawForm.amount || withdrawForm.amount <= 0) {
    ElMessage.error('请输入提现金额')
    return
  }
  if (!withdrawForm.bank_name) {
    ElMessage.error('请输入开户银行')
    return
  }
  if (!withdrawForm.bank_account) {
    ElMessage.error('请输入银行账号')
    return
  }

  submitting.value = true
  try {
    await financeApi.createWithdraw(withdrawForm)
    ElMessage.success('提现申请已提交')
    showWithdrawDialog.value = false
    withdrawForm.amount = 0
    withdrawForm.bank_name = ''
    withdrawForm.bank_account = ''
    await loadBalance()
    await loadWithdraws()
  } catch (error) {
    console.error('Failed to submit withdraw:', error)
  } finally {
    submitting.value = false
  }
}
</script>

<style scoped>
.finance-page {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.balance-row {
  margin-bottom: 20px;
}

.balance-card {
  text-align: center;
}

.balance-label {
  color: #909399;
  font-size: 14px;
  margin-bottom: 10px;
}

.balance-value {
  font-size: 28px;
  font-weight: 600;
  color: #409eff;
}

.action-bar {
  margin-bottom: 20px;
}

.income {
  color: #67c23a;
}

.expense {
  color: #f56c6c;
}
</style>
