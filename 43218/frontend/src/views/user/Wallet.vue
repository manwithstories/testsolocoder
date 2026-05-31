<template>
  <div class="wallet-page">
    <div class="page-header">
      <h2 class="page-title">我的钱包</h2>
    </div>

    <div class="wallet-content">
      <el-card class="balance-card">
        <div class="balance-info">
          <div class="balance-label">账户余额</div>
          <div class="balance-amount">¥{{ balance.toFixed(2) }}</div>
          <div class="balance-actions">
            <el-button type="primary" size="large" @click="showRechargeDialog = true">
              充值
            </el-button>
            <el-button size="large" @click="showWithdrawDialog = true">
              提现
            </el-button>
          </div>
        </div>
      </el-card>

      <el-card class="logs-card">
        <template #header>
          <span>交易记录</span>
        </template>
        <el-table :data="logs" v-loading="loading" style="width: 100%">
          <el-table-column prop="type" label="类型" width="100">
            <template #default="{ row }">
              <el-tag :type="getTypeTag(row.type)">{{ getTypeText(row.type) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column prop="amount" label="金额" width="120">
            <template #default="{ row }">
              <span :class="{ 'text-success': row.amount > 0, 'text-danger': row.amount < 0 }">
                {{ row.amount > 0 ? '+' : '' }}{{ row.amount.toFixed(2) }}
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="balance" label="余额" width="120">
            <template #default="{ row }">
              ¥{{ row.balance.toFixed(2) }}
            </template>
          </el-table-column>
          <el-table-column prop="description" label="描述" />
          <el-table-column prop="orderNo" label="订单号" width="200" />
          <el-table-column prop="createdAt" label="时间" width="180">
            <template #default="{ row }">
              {{ formatTime(row.createdAt) }}
            </template>
          </el-table-column>
        </el-table>

        <div class="pagination-wrapper">
          <el-pagination
            v-model:current-page="pagination.page"
            v-model:page-size="pagination.pageSize"
            :page-sizes="[10, 20, 50]"
            :total="total"
            layout="total, sizes, prev, pager, next, jumper"
            @size-change="fetchLogs"
            @current-change="fetchLogs"
          />
        </div>
      </el-card>
    </div>

    <el-dialog v-model="showRechargeDialog" title="充值" width="400px">
      <el-form label-width="80px">
        <el-form-item label="充值金额">
          <el-input-number v-model="rechargeAmount" :min="0.01" :precision="2" />
        </el-form-item>
        <el-form-item label="支付方式">
          <el-radio-group v-model="rechargeMethod">
            <el-radio value="alipay">支付宝</el-radio>
            <el-radio value="wechat">微信支付</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRechargeDialog = false">取消</el-button>
        <el-button type="primary" :loading="recharging" @click="submitRecharge">
          确认充值
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showWithdrawDialog" title="提现" width="400px">
      <el-form label-width="80px">
        <el-form-item label="提现金额">
          <el-input-number v-model="withdrawAmount" :min="1" :max="balance" :precision="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showWithdrawDialog = false">取消</el-button>
        <el-button type="primary" :loading="withdrawing" @click="submitWithdraw">
          确认提现
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { userApi } from '@/api'
import type { WalletLog } from '@/types'

const loading = ref(false)
const balance = ref(0)
const logs = ref<WalletLog[]>([])
const total = ref(0)

const showRechargeDialog = ref(false)
const showWithdrawDialog = ref(false)
const rechargeAmount = ref(100)
const withdrawAmount = ref(100)
const rechargeMethod = ref('alipay')
const recharging = ref(false)
const withdrawing = ref(false)

const pagination = reactive({
  page: 1,
  pageSize: 10
})

function getTypeText(type: string): string {
  const typeMap: Record<string, string> = {
    recharge: '充值',
    withdraw: '提现',
    payment: '支付',
    income: '收入',
    refund: '退款',
    commission: '佣金'
  }
  return typeMap[type] || type
}

function getTypeTag(type: string): string {
  const tagMap: Record<string, string> = {
    recharge: 'success',
    withdraw: 'warning',
    payment: 'danger',
    income: 'success',
    refund: 'primary',
    commission: 'info'
  }
  return tagMap[type] || 'info'
}

function formatTime(time: string): string {
  return new Date(time).toLocaleString()
}

async function fetchBalance() {
  try {
    const res = await userApi.getWalletBalance()
    balance.value = res.data.balance
  } catch (error) {
    console.error('Failed to fetch balance:', error)
  }
}

async function fetchLogs() {
  loading.value = true
  try {
    const res = await userApi.getWalletLogs({
      page: pagination.page,
      pageSize: pagination.pageSize
    })
    logs.value = res.data
    total.value = res.pagination.total
  } catch (error) {
    console.error('Failed to fetch logs:', error)
  } finally {
    loading.value = false
  }
}

async function submitRecharge() {
  if (rechargeAmount.value <= 0) {
    ElMessage.warning('请输入正确的充值金额')
    return
  }

  recharging.value = true
  try {
    await userApi.recharge({
      amount: rechargeAmount.value,
      paymentMethod: rechargeMethod.value
    })
    ElMessage.success('充值成功')
    showRechargeDialog.value = false
    fetchBalance()
    fetchLogs()
  } catch (error: any) {
    ElMessage.error(error.message || '充值失败')
  } finally {
    recharging.value = false
  }
}

async function submitWithdraw() {
  if (withdrawAmount.value <= 0 || withdrawAmount.value > balance.value) {
    ElMessage.warning('请输入正确的提现金额')
    return
  }

  withdrawing.value = true
  try {
    await userApi.withdraw({ amount: withdrawAmount.value })
    ElMessage.success('提现申请已提交')
    showWithdrawDialog.value = false
    fetchBalance()
    fetchLogs()
  } catch (error: any) {
    ElMessage.error(error.message || '提现失败')
  } finally {
    withdrawing.value = false
  }
}

onMounted(() => {
  fetchBalance()
  fetchLogs()
})
</script>

<style lang="scss" scoped>
.wallet-page {
  .wallet-content {
    display: flex;
    flex-direction: column;
    gap: 20px;
  }

  .balance-card {
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border: none;
    color: #fff;

    :deep(.el-card__body) {
      padding: 30px;
    }

    .balance-info {
      text-align: center;

      .balance-label {
        font-size: 14px;
        opacity: 0.8;
        margin-bottom: 12px;
      }

      .balance-amount {
        font-size: 42px;
        font-weight: 600;
        margin-bottom: 20px;
      }

      .balance-actions {
        display: flex;
        justify-content: center;
        gap: 16px;
      }
    }
  }

  .logs-card {
    .pagination-wrapper {
      display: flex;
      justify-content: center;
      padding: 20px 0 0;
    }
  }

  .text-success {
    color: var(--success-color);
  }

  .text-danger {
    color: var(--danger-color);
  }
}
</style>
