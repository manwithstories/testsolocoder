<template>
  <div class="part-requests-page">
    <div class="container">
      <div class="page-header">
        <h2 class="page-title">配件申请</h2>
        <el-button type="primary" @click="showCreateDialog = true">
          新建申请
        </el-button>
      </div>

      <el-tabs v-model="activeTab">
        <el-tab-pane label="全部" name="all">
          <el-table :data="filteredRequests" style="width: 100%">
            <el-table-column prop="request_no" label="申请编号" width="200" />
            <el-table-column prop="total_amount" label="总金额" width="120">
              <template #default="{ row }">¥{{ row.total_amount }}</template>
            </el-table-column>
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="getRequestStatusType(row.status)" size="small">
                  {{ getRequestStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
            <el-table-column prop="created_at" label="申请时间" width="180">
              <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
            </el-table-column>
            <el-table-column label="操作" width="150">
              <template #default="{ row }">
                <el-button size="small" @click="viewRequest(row)">查看</el-button>
                <el-button
                  v-if="row.status === 'shipped'"
                  type="success"
                  size="small"
                  @click="receiveRequest(row)"
                >
                  确认收货
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="待审批" name="pending">
          <el-table :data="filteredRequests" style="width: 100%">
            <el-table-column prop="request_no" label="申请编号" width="200" />
            <el-table-column prop="total_amount" label="总金额" width="120">
              <template #default="{ row }">¥{{ row.total_amount }}</template>
            </el-table-column>
            <el-table-column prop="created_at" label="申请时间" width="180">
              <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
            </el-table-column>
          </el-table>
        </el-tab-pane>

        <el-tab-pane label="已通过" name="approved">
          <el-table :data="filteredRequests" style="width: 100%">
            <el-table-column prop="request_no" label="申请编号" width="200" />
            <el-table-column prop="total_amount" label="总金额" width="120">
              <template #default="{ row }">¥{{ row.total_amount }}</template>
            </el-table-column>
            <el-table-column prop="created_at" label="申请时间" width="180">
              <template #default="{ row }">{{ formatTime(row.created_at) }}</template>
            </el-table-column>
          </el-table>
        </el-tab-pane>
      </el-tabs>

      <div v-if="filteredRequests.length === 0" class="empty-state">
        <el-empty description="暂无申请" />
      </div>
    </div>

    <el-dialog v-model="showCreateDialog" title="新建配件申请" width="600px">
      <el-form :model="createForm" label-width="80px">
        <el-form-item label="配件列表">
          <div v-for="(item, index) in createForm.items" :key="index" class="item-row">
            <el-select v-model="item.part_id" placeholder="选择配件" style="width: 200px;">
              <el-option
                v-for="part in parts"
                :key="part.id"
                :label="`${part.name} (¥${part.price})`"
                :value="part.id"
              />
            </el-select>
            <el-input-number v-model="item.quantity" :min="1" style="width: 120px; margin-left: 10px;" />
            <el-button type="danger" size="small" @click="removeItem(index)" style="margin-left: 10px;">
              删除
            </el-button>
          </div>
          <el-button type="primary" size="small" @click="addItem" style="margin-top: 10px;">
            添加配件
          </el-button>
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="createForm.remark" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="submitRequest" :loading="submitting">
          提交申请
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showDetailDialog" title="申请详情" width="600px">
      <el-descriptions :column="2" border v-if="selectedRequest">
        <el-descriptions-item label="申请编号">{{ selectedRequest.request_no }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getRequestStatusType(selectedRequest.status)">
            {{ getRequestStatusText(selectedRequest.status) }}
          </el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="总金额">¥{{ selectedRequest.total_amount }}</el-descriptions-item>
        <el-descriptions-item label="申请时间">{{ formatTime(selectedRequest.created_at) }}</el-descriptions-item>
        <el-descriptions-item label="备注" :span="2">{{ selectedRequest.remark || '-' }}</el-descriptions-item>
      </el-descriptions>
      <el-table :data="selectedItems" style="width: 100%; margin-top: 20px;">
        <el-table-column prop="part.name" label="配件名称" />
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
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { partApi } from '@/api/part'
import type { PartRequest, PartRequestItem, Part } from '@/types'
import dayjs from 'dayjs'

const activeTab = ref('all')
const requests = ref<PartRequest[]>([])
const parts = ref<Part[]>([])
const showCreateDialog = ref(false)
const showDetailDialog = ref(false)
const submitting = ref(false)
const selectedRequest = ref<PartRequest | null>(null)
const selectedItems = ref<PartRequestItem[]>([])

const createForm = reactive({
  items: [{ part_id: null as number | null, quantity: 1 }],
  remark: ''
})

const filteredRequests = computed(() => {
  if (activeTab.value === 'all') return requests.value
  return requests.value.filter(r => r.status === activeTab.value)
})

onMounted(async () => {
  await loadRequests()
  await loadParts()
})

async function loadRequests() {
  try {
    const res = await partApi.getPartRequests({ page: 1, page_size: 50 })
    requests.value = res.data?.list || []
  } catch (error) {
    console.error('Failed to load requests:', error)
  }
}

async function loadParts() {
  try {
    const res = await partApi.getParts({ page: 1, page_size: 100 })
    parts.value = res.data?.list || []
  } catch (error) {
    console.error('Failed to load parts:', error)
  }
}

function getRequestStatusType(status: string): string {
  const statusMap: Record<string, string> = {
    pending: 'warning',
    approved: 'primary',
    rejected: 'danger',
    shipped: 'success',
    received: 'info'
  }
  return statusMap[status] || 'info'
}

function getRequestStatusText(status: string): string {
  const statusMap: Record<string, string> = {
    pending: '待审批',
    approved: '已通过',
    rejected: '已拒绝',
    shipped: '已发货',
    received: '已收货'
  }
  return statusMap[status] || status
}

function formatTime(time: string): string {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

function addItem() {
  createForm.items.push({ part_id: null, quantity: 1 })
}

function removeItem(index: number) {
  if (createForm.items.length > 1) {
    createForm.items.splice(index, 1)
  }
}

async function submitRequest() {
  const validItems = createForm.items.filter(item => item.part_id && item.quantity > 0)
  if (validItems.length === 0) {
    ElMessage.error('请添加配件')
    return
  }

  submitting.value = true
  try {
    await partApi.createPartRequest({
      items: validItems.map(item => ({ part_id: item.part_id!, quantity: item.quantity })),
      remark: createForm.remark
    })
    ElMessage.success('申请已提交')
    showCreateDialog.value = false
    createForm.items = [{ part_id: null, quantity: 1 }]
    createForm.remark = ''
    loadRequests()
  } catch (error) {
    console.error('Failed to submit request:', error)
  } finally {
    submitting.value = false
  }
}

async function viewRequest(request: PartRequest) {
  selectedRequest.value = request
  try {
    const res = await partApi.getPartRequestDetail(request.id)
    selectedItems.value = res.data?.items || []
    showDetailDialog.value = true
  } catch (error) {
    console.error('Failed to load request detail:', error)
  }
}

async function receiveRequest(request: PartRequest) {
  try {
    await partApi.receivePartRequest(request.id)
    ElMessage.success('已确认收货')
    loadRequests()
  } catch (error) {
    console.error('Failed to receive request:', error)
  }
}
</script>

<style scoped>
.part-requests-page {
  min-height: 100vh;
  background-color: #f5f7fa;
}

.item-row {
  display: flex;
  align-items: center;
  margin-bottom: 10px;
}
</style>
