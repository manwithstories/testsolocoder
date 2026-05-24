<template>
  <el-card>
    <template #header>
      <div class="header">
        <span>保险理赔</span>
        <el-button type="primary" @click="$router.push('/insurance/create')">
          <el-icon><Plus /></el-icon>
          申请理赔
        </el-button>
      </div>
    </template>

    <el-table :data="claims" v-loading="loading">
      <el-table-column prop="claim_no" label="理赔编号" width="180" />
      <el-table-column label="订单" min-width="150">
        <template #default="{ row }">{{ row.order?.order_no || '-' }}</template>
      </el-table-column>
      <el-table-column prop="damage_desc" label="损坏描述" show-overflow-tooltip />
      <el-table-column label="预估费用" width="120">
        <template #default="{ row }">¥{{ row.estimated_cost }}</template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="statusTagType(row.status)">{{ statusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="申请时间" width="180" />
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button type="primary" link @click="viewDetail(row)">详情</el-button>
          <el-button v-if="row.status === 'pending' && role === 'owner'" type="success" link @click="reviewClaim(row, 'approved')">通过</el-button>
          <el-button v-if="row.status === 'pending' && role === 'owner'" type="danger" link @click="reviewClaim(row, 'rejected')">拒绝</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-pagination
      style="margin-top: 20px; justify-content: flex-end"
      layout="prev, pager, next, total"
      :total="total"
      :page-size="pageSize"
      v-model:current-page="currentPage"
      @current-change="fetchClaims"
    />
  </el-card>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useUserStore } from '@/stores/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'

const userStore = useUserStore()
const role = computed(() => userStore.role)

const loading = ref(false)
const claims = ref<InsuranceClaim[]>([])
const total = ref(0)
const currentPage = ref(1)
const pageSize = ref(10)

onMounted(() => {
  fetchClaims()
})

async function fetchClaims() {
  loading.value = true
  try {
    const res: any = await request.get('/insurance/claims', {
      params: { page: currentPage.value, page_size: pageSize.value }
    })
    claims.value = res.data.list || []
    total.value = res.data.total || 0
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function viewDetail(row: InsuranceClaim) {
  ElMessage.info(`理赔详情: ${row.claim_no}`)
}

async function reviewClaim(row: InsuranceClaim, status: string) {
  const actionText = status === 'approved' ? '通过' : '拒绝'
  ElMessageBox.confirm(`确定${actionText}该理赔申请吗？`, '审核确认', {
    type: 'warning'
  }).then(async () => {
    try {
      await request.put('/insurance/claims/review', {
        claim_id: row.id,
        status,
        actual_cost: row.estimated_cost,
        deducted_amount: status === 'approved' ? row.estimated_cost : 0,
        review_remark: `${actionText}理赔申请`
      })
      ElMessage.success('审核成功')
      fetchClaims()
    } catch (e: any) {
      ElMessage.error(e.message || '审核失败')
    }
  }).catch(() => {})
}

function statusText(status: string) {
  const map: Record<string, string> = {
    pending: '待审核', reviewing: '审核中', approved: '已通过', rejected: '已拒绝'
  }
  return map[status] || status
}

function statusTagType(status: string) {
  const map: Record<string, string> = {
    pending: 'warning', reviewing: 'primary', approved: 'success', rejected: 'danger'
  }
  return map[status] || ''
}
</script>

<style scoped>
.header {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
</style>
