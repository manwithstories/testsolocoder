<template>
  <div class="admin-part-requests-page">
    <div class="page-header">
      <h2 class="page-title">配件申请审核</h2>
    </div>

    <el-card>
      <el-table :data="requests" style="width: 100%">
        <el-table-column prop="request_no" label="申请编号" width="200" />
        <el-table-column label="申请人" width="150">
          <template #default="{ row }">{{ row.technician?.username || '-' }}</template>
        </el-table-column>
        <el-table-column prop="total_amount" label="总金额" width="120">
          <template #default="{ row }">¥{{ row.total_amount }}</template>
        </el-table-column>
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
            <el-button size="small" @click="viewDetail(row)">查看</el-button>
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
              @click="ship(row)"
            >
              发货
            </el-button>
          </template>
        </el-table-column>
      </el-table>
    </el-card>

    <el-dialog v-model="showDetailDialog" title="申请详情" width="600px">
      <el-descriptions :column="2" border v-if="selectedRequest">
        <el-descriptions-item label="申请编号">
          {{ selectedRequest.request_no }}
        </el-descriptions-item>
        <el-descriptions-item label="申请人">
          {{ selectedRequest.technician?.username }}
        </el-descriptions-item>
        <el-descriptions-item label="总金额">
          ¥{{ selectedRequest.total_amount }}
        </el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(selectedRequest.status)">
            {{ getStatusText(selectedRequest.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">
          {{ selectedRequest.remark || '-' }}
        </el-descriptions-item>
      </el-descriptions>

      <el-table :data="selectedItems" style="width: 100%; margin-top: 20px;">
        <el-table-column label="配件名称">
          <template #default="{ row }">{{ row.part?.name }}</template>
        </el-table-column>
        <el-table-column prop="quantity" label="数量" width="100" />
        <el-table-column prop="price" label="单价" width="100">
          <template #default="{ row }">¥{{ row.price }}</template>
        </el-table-column>
        <el-table-column label="小计" width="120">
          <template #default="{ row }">¥{{ row.price * row.quantity }}</template>
        </el-table-column>
      </el-table>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { partApi } from '@/api/part'
import type { PartRequest, PartRequestItem } from '@/types'
import dayjs from 'dayjs'

const requests = ref<PartRequest[]>([])
const showDetailDialog = ref(false)
const selectedRequest = ref<PartRequest | null>(null)
const selectedItems = ref<PartRequestItem[]>([])

onMounted(() => {
  loadRequests()
})

async function loadRequests() {
  try {
    const res = await partApi.getPartRequests({ page: 1, page_size: 50 })
    requests.value = res.data?.list || []
  } catch (error) {
    console.error('Failed to load requests:', error)
  }
}

function getStatusType(status: string): string {
  const typeMap: Record<string, string> = {
    pending: 'warning',
    approved: 'primary',
    rejected: 'danger',
    shipped: 'success',
    received: 'info'
  }
  return typeMap[status] || 'info'
}

function getStatusText(status: string): string {
  const textMap: Record<string, string> = {
    pending: '待审批',
    approved: '已通过',
    rejected: '已拒绝',
    shipped: '已发货',
    received: '已收货'
  }
  return textMap[status] || status
}

function formatTime(time: string): string {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

async function viewDetail(request: PartRequest) {
  selectedRequest.value = request
  try {
    const res = await partApi.getPartRequestDetail(request.id)
    selectedItems.value = res.data?.items || []
    showDetailDialog.value = true
  } catch (error) {
    console.error('Failed to load detail:', error)
  }
}

async function approve(request: PartRequest) {
  try {
    await ElMessageBox.confirm('确定要通过该申请吗？', '确认', { type: 'success' })
    await partApi.approvePartRequest(request.id)
    ElMessage.success('审批通过')
    loadRequests()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to approve:', error)
    }
  }
}

async function reject(request: PartRequest) {
  try {
    const { value: remark } = await ElMessageBox.prompt('请输入拒绝原因', '拒绝申请', {
      inputValidator: (value) => {
        if (!value) return '请输入拒绝原因'
        return true
      }
    })
    await partApi.rejectPartRequest(request.id, { remark })
    ElMessage.success('已拒绝')
    loadRequests()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to reject:', error)
    }
  }
}

async function ship(request: PartRequest) {
  try {
    await ElMessageBox.confirm('确定要发货吗？', '确认')
    await partApi.shipPartRequest(request.id)
    ElMessage.success('已发货')
    loadRequests()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to ship:', error)
    }
  }
}
</script>

<style scoped>
.admin-part-requests-page {
  padding: 0;
}
</style>
