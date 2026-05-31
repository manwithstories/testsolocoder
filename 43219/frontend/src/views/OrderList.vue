<template>
  <AppLayout>
    <div class="page">
      <h2>订单管理</h2>
      <div class="row" style="margin-bottom:16px">
        <el-select v-model="status" placeholder="全部状态" clearable style="width:160px" @change="load">
          <el-option value="created" label="已创建" />
          <el-option value="reported" label="待验收" />
          <el-option value="confirmed" label="已确认" />
          <el-option value="paid" label="已结算" />
          <el-option value="refunding" label="退款中" />
          <el-option value="refunded" label="已退款" />
        </el-select>
      </div>
      <el-table :data="list" border stripe>
        <el-table-column prop="id" label="订单号" width="80" />
        <el-table-column label="服务" width="200">
          <template #default="{ row }">
            {{ row.booking?.service?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="total_amount" label="金额" width="100" />
        <el-table-column label="状态" width="120">
          <template #default="{ row }">
            <el-tag :type="tagType(row.status)">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="created_at" label="创建时间" width="180" />
        <el-table-column label="操作" width="320">
          <template #default="{ row }">
            <el-button v-if="userStore.role==='staff' && (row.status==='created' || row.status==='reported')"
              size="small" type="primary" @click="openReport(row)">提交报告</el-button>
            <el-button v-if="userStore.role==='customer' && (row.status==='reported' || row.status==='created')"
              size="small" type="success" @click="confirm(row)">确认验收</el-button>
            <el-button v-if="userStore.role==='customer' && row.status!=='refunding' && row.status!=='refunded'"
              size="small" type="danger" @click="openRefund(row)">申请退款</el-button>
            <el-button v-if="userStore.role==='customer' && row.status==='paid'"
              size="small" @click="$router.push(`/orders/${row.id}/review')">写评价</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="reportVisible" title="提交服务报告">
      <el-form :model="reportForm" label-width="80px">
        <el-form-item label="描述"><el-input v-model="reportForm.report_text" type="textarea" :rows="4" /></el-form-item>
        <el-form-item label="图片"><el-input v-model="reportForm.report_images" placeholder="多个逗号分隔" /></el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="reportVisible=false">取消</el-button>
        <el-button type="primary" @click="submitReport">提交</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="refundVisible" title="申请退款">
      <el-input v-model="refundReason" type="textarea" :rows="3" placeholder="请说明退款原因" />
      <template #footer>
        <el-button @click="refundVisible=false">取消</el-button>
        <el-button type="danger" @click="submitRefund">提交</el-button>
      </template>
    </el-dialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import AppLayout from '../components/AppLayout.vue'
import { listOrders, submitReport, confirmOrder, requestRefund, type OrderItem } from '../api/order'
import { useUserStore } from '../stores/user'

const userStore = useUserStore()
const list = ref<OrderItem[]>([])
const status = ref('')
const reportVisible = ref(false)
const refundVisible = ref(false)
const currentId = ref<number | null>(null)
const reportForm = reactive({ report_text: '', report_images: '' })
const refundReason = ref('')

async function load() {
  const params: Record<string, string> = {}
  if (status.value) params.status = status.value
  const res = await listOrders(params)
  list.value = (res.data as any).data || []
}

function statusLabel(s: string) {
  return { created: '已创建', reported: '待验收', confirmed: '已确认', paid: '已结算', refunding: '退款中', refunded: '已退款', closed: '已关闭' }[s] || s
}
function tagType(s: string) {
  return { created: 'info', reported: 'warning', confirmed: '', paid: 'success', refunding: 'danger', refunded: 'info', closed: 'info' }[s] || ''
}

function openReport(row: OrderItem) {
  currentId.value = row.id
  reportForm.report_text = ''
  reportForm.report_images = ''
  reportVisible.value = true
}
async function submitReport() {
  if (!reportForm.report_text) { ElMessage.warning('请填写报告内容') }
  await submitReport(currentId.value!, reportForm)
  ElMessage.success('报告已提交')
  reportVisible.value = false
  load()
}
async function confirm(row: OrderItem) {
  await ElMessageBox.confirm('确认验收并结算?', '提示')
  await confirmOrder(row.id)
  ElMessage.success('已确认')
  load()
}
function openRefund(row: OrderItem) {
  currentId.value = row.id
  refundReason.value = ''
  refundVisible.value = true
}
async function submitRefund() {
  await requestRefund(currentId.value!, { reason: refundReason.value })
  ElMessage.success('退款申请已提交')
  refundVisible.value = false
  load()
}

onMounted(load)
</script>
