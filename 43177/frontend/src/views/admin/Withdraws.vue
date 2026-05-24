<template>
  <div class="admin-withdraws-page">
    <div class="page-header">
      <h2 class="page-title">提现审核</h2>
    </div>

    <el-card>
      <el-table :data="withdraws" style="width: 100%">
        <el-table-column prop="request_no" label="申请编号" width="200" />
        <el-table-column label="申请人" width="150">
          <template #default="{ row }">{{ row.technician?.username || '-' }}</template>
        </el-table-column>
        <el-table-column prop="amount" label="金额" width="120">
          <template #default="{ row }">¥{{ row.amount }}</template>
        </el-table-column>
        <el-table-column prop="bank_name" label="开户银行" width="150" />
        <el-table-column prop="bank_account" label="银行账号" width="200" />
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="申请时间" width="180">
          <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'pending'"
              type="success"
              size="small"
              @click="approve(row)"
            >
              通过
            </el-button>
            <el-button
              v-if="row.status === 'pending'"
              type="danger"
              size="small"
              @click="reject(row)"
            >
              拒绝
            </el-button>
            <el-button
              v-if="row.status === 'approved'"
              type="primary"
              size="small"
              @click="complete(row)"
            >
              打款
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div v-if="withdraws.length === 0" class="empty-state">
        <el-empty description="暂无提现申请" />
      </div>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { adminFinanceApi } from '@/api/finance'
import type { WithdrawRequest } from '@/types'
import dayjs from 'dayjs'

const withdraws = ref<WithdrawRequest[]>([])

onMounted(() => {
  loadWithdraws()
})

async function loadWithdraws() {
  try {
    const res = await adminFinanceApi.getWithdraws({ page: 1, page_size: 50 })
    withdraws.value = res.data?.list || []
  } catch (error) {
    console.error('Failed to load withdraws:', error)
  }
}

function getStatusType(status: string): string {
  const typeMap: Record<string, string> = {
    pending: 'warning',
    approved: 'primary',
    rejected: 'danger',
    completed: 'success'
  }
  return typeMap[status] || 'info'
}

function getStatusText(status: string): string {
  const textMap: Record<string, string> = {
    pending: '待审批',
    approved: '已通过',
    rejected: '已拒绝',
    completed: '已打款'
  }
  return textMap[status] || status
}

function formatTime(time: string): string {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

async function approve(request: WithdrawRequest) {
  try {
    await ElMessageBox.confirm(
      `确定要通过该提现申请吗？金额：¥${request.amount}`,
      '确认',
      { type: 'warning' }
    )
    await adminFinanceApi.approveWithdraw(request.id)
    ElMessage.success('审批通过')
    loadWithdraws()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to approve:', error)
    }
  }
}

async function reject(request: WithdrawRequest) {
  try {
    const { value: remark } = await ElMessageBox.prompt('请输入拒绝原因', '拒绝提现', {
      inputValidator: (value) => {
        if (!value) return '请输入拒绝原因'
        return true
      }
    })
    await adminFinanceApi.rejectWithdraw(request.id, { remark })
    ElMessage.success('已拒绝')
    loadWithdraws()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to reject:', error)
    }
  }
}

async function complete(request: WithdrawRequest) {
  try {
    await ElMessageBox.confirm('确定要完成打款吗？', '确认')
    await adminFinanceApi.completeWithdraw(request.id)
    ElMessage.success('打款完成')
    loadWithdraws()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to complete:', error)
    }
  }
}
</script>

<style scoped>
.admin-withdraws-page {
  padding: 0;
}

.empty-state {
  padding: 40px;
}
</style>
