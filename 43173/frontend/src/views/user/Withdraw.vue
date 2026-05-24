<template>
  <div class="withdraw">
    <div class="page-header">
      <h2>申请提现</h2>
    </div>
    
    <el-card class="balance-card">
      <div class="balance-info">
        <div class="label">可提现余额</div>
        <div class="amount">¥{{ formatAmount(balance.available_balance || 0) }}</div>
      </div>
    </el-card>
    
    <el-form 
      ref="formRef"
      :model="form" 
      :rules="rules"
      label-width="100px"
      style="max-width: 500px; margin-top: 24px;"
    >
      <el-form-item label="提现金额" prop="amount">
        <el-input-number 
          v-model="form.amount" 
          :min="10" 
          :max="balance.available_balance || 100000"
          :precision="2"
          :step="100"
        />
        <div class="tip">最低提现金额 ¥10.00</div>
      </el-form-item>
      
      <el-form-item label="提现方式" prop="method">
        <el-radio-group v-model="form.method">
          <el-radio value="alipay">支付宝</el-radio>
          <el-radio value="wechat">微信</el-radio>
          <el-radio value="bank">银行卡</el-radio>
        </el-radio-group>
      </el-form-item>
      
      <el-form-item label="收款账号" prop="account">
        <el-input v-model="form.account" placeholder="请输入收款账号" />
      </el-form-item>
      
      <el-form-item label="账户名称" prop="account_name">
        <el-input v-model="form.account_name" placeholder="请输入账户名称" />
      </el-form-item>
      
      <el-form-item label="开户行" v-if="form.method === 'bank'">
        <el-input v-model="form.bank_name" placeholder="请输入开户行" />
      </el-form-item>
      
      <el-form-item>
        <el-button type="primary" :loading="submitting" @click="submit">
          提交申请
        </el-button>
      </el-form-item>
    </el-form>
    
    <div class="withdraw-history">
      <h3>提现记录</h3>
      <el-table :data="history" style="width: 100%">
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">
            <span>¥{{ formatAmount(row.amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="actual_amount" label="实际到账" width="120">
          <template #default="{ row }">
            <span class="success">¥{{ formatAmount(row.actual_amount) }}</span>
          </template>
        </el-table-column>
        <el-table-column prop="method" label="方式" width="100">
          <template #default="{ row }">
            {{ getMethodText(row.method) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusTag(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="申请时间" width="180" />
      </el-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { revenueApi } from '@/api/revenue'

const formRef = ref<FormInstance>()
const submitting = ref(false)
const balance = ref<any>({})
const history = ref<any[]>([])

const form = reactive({
  amount: 100,
  method: 'alipay',
  account: '',
  account_name: '',
  bank_name: ''
})

const rules: FormRules = {
  amount: [{ required: true, message: '请输入提现金额', trigger: 'blur' }],
  method: [{ required: true, message: '请选择提现方式', trigger: 'change' }],
  account: [{ required: true, message: '请输入收款账号', trigger: 'blur' }],
  account_name: [{ required: true, message: '请输入账户名称', trigger: 'blur' }]
}

onMounted(() => {
  loadBalance()
  loadHistory()
})

async function loadBalance() {
  try {
    balance.value = await revenueApi.getSummary()
  } catch (e) {
    console.error(e)
  }
}

async function loadHistory() {
  try {
    const res = await revenueApi.getWithdrawList({ page: 1, page_size: 20 })
    history.value = res.list
  } catch (e) {
    console.error(e)
  }
}

async function submit() {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    
    submitting.value = true
    try {
      await revenueApi.requestWithdraw(form)
      ElMessage.success('申请提交成功')
      loadBalance()
      loadHistory()
      form.amount = 100
      form.account = ''
      form.account_name = ''
    } catch (e) {
      console.error(e)
    } finally {
      submitting.value = false
    }
  })
}

function formatAmount(amount: number) {
  return amount.toFixed(2)
}

function getMethodText(method: string) {
  const texts: Record<string, string> = {
    alipay: '支付宝',
    wechat: '微信',
    bank: '银行卡'
  }
  return texts[method] || method
}

function getStatusTag(status: number): 'primary' | 'success' | 'warning' | 'info' | 'danger' {
  const tags: Record<number, 'primary' | 'success' | 'warning' | 'info' | 'danger'> = {
    0: 'warning',
    1: 'success',
    2: 'danger',
    3: 'success',
    4: 'info'
  }
  return tags[status] || 'info'
}

function getStatusText(status: number) {
  const texts: Record<number, string> = {
    0: '待审核',
    1: '已通过',
    2: '已拒绝',
    3: '已打款',
    4: '失败'
  }
  return texts[status] || '未知'
}
</script>

<style scoped lang="scss">
.withdraw {
  .balance-card {
    margin-bottom: 24px;
    
    .balance-info {
      text-align: center;
      
      .label {
        font-size: 14px;
        color: var(--text-light);
        margin-bottom: 8px;
      }
      
      .amount {
        font-size: 32px;
        font-weight: 600;
        color: var(--primary-color);
      }
    }
  }
  
  .tip {
    font-size: 12px;
    color: var(--text-light);
    margin-top: 4px;
  }
  
  .withdraw-history {
    margin-top: 40px;
    
    h3 {
      margin: 0 0 16px 0;
    }
  }
  
  .success {
    color: var(--success-color);
  }
}
</style>
