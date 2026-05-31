<template>
  <AppLayout>
    <div class="page">
      <div class="row" style="margin-bottom:16px">
        <h2 style="margin:0">我的收益</h2>
        <el-button type="primary" @click="wdVisible=true">申请提现</el-button>
      </div>
      <el-descriptions :column="3" bordered>
        <el-descriptions-item label="累计收入">¥{{ d.total_earned?.toFixed(2) || '0.00' }}</el-descriptions-item>
        <el-descriptions-item label="已提现">¥{{ d.withdrawn?.toFixed(2) || '0.00' }}</el-descriptions-item>
        <el-descriptions-item label="当前余额">¥{{ d.balance?.toFixed(2) || '0.00' }}</el-descriptions-item>
      </el-descriptions>

      <h3 style="margin-top:24px">分成明细</h3>
      <el-table :data="d.settlements || []" border stripe>
        <el-table-column prop="id" label="编号" width="80" />
        <el-table-column prop="order_id" label="订单号" width="100" />
        <el-table-column prop="total_amount" label="订单金额" width="120" />
        <el-table-column prop="staff_share" label="我的分成" width="120" />
        <el-table-column prop="status" label="状态" width="100" />
        <el-table-column prop="created_at" label="时间" width="180" />
      </el-table>

      <h3 style="margin-top:24px">提现记录</h3>
      <el-table :data="d.withdrawals || []" border stripe>
        <el-table-column prop="id" label="编号" width="80" />
        <el-table-column prop="amount" label="金额" width="120" />
        <el-table-column prop="account" label="账号" />
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag>{{ { pending: '处理中', paid: '已到账', rejected: '已拒绝' }[row.status] || row.status }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="申请时间" width="180" />
      </el-table>
    </div>

    <el-dialog v-model="wdVisible" title="申请提现">
      <el-form :model="wdForm" label-width="80px">
        <el-form-item label="金额"><el-input-number v-model="wdForm.amount" :min="1" :precision="2" /></el-form-item>
        <el-form-item label="账号"><el-input v-model="wdForm.account" placeholder="支付宝/银行卡号" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="wdVisible=false">取消</el-button>
        <el-button type="primary" @click="submitWithdraw">提交</el-button>
      </template>
    </el-dialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import AppLayout from '../../components/AppLayout.vue'
import { myEarnings, requestWithdrawal } from '../../api/finance'

const d = ref<any>({})
const wdVisible = ref(false)
const wdForm = reactive({ amount: 100, account: '' })

async function load() {
  const res = await myEarnings()
  d.value = (res.data as any).data || {}
}

async function submitWithdraw() {
  if (!wdForm.account) { ElMessage.warning('请填写账号'); return }
  await requestWithdrawal(wdForm)
  ElMessage.success('申请已提交')
  wdVisible.value = false
  load()
}

onMounted(load)
</script>
