<template>
  <el-card v-if="service">
    <el-descriptions :column="2" border title="服务需求">
      <el-descriptions-item label="服务编号">{{ service.service_no }}</el-descriptions-item>
      <el-descriptions-item label="状态">
        <el-tag :type="statusTagType(service.status)">{{ statusText(service.status) }}</el-tag>
      </el-descriptions-item>
      <el-descriptions-item label="标题">{{ service.title }}</el-descriptions-item>
      <el-descriptions-item label="区域">{{ service.region }}</el-descriptions-item>
      <el-descriptions-item label="地址">{{ service.address }}</el-descriptions-item>
      <el-descriptions-item label="服务日期">{{ service.service_date || '-' }}</el-descriptions-item>
      <el-descriptions-item label="服务时间">{{ service.service_time || '-' }}</el-descriptions-item>
      <el-descriptions-item label="时长">{{ service.duration }}分钟</el-descriptions-item>
      <el-descriptions-item label="预算">¥{{ service.budget_min }} ~ ¥{{ service.budget_max }}</el-descriptions-item>
      <el-descriptions-item label="成交价" v-if="service.final_price">¥{{ service.final_price }}</el-descriptions-item>
      <el-descriptions-item label="发布人">{{ service.user?.nickname || service.user?.username }}</el-descriptions-item>
      <el-descriptions-item label="飞手">{{ service.pilot?.nickname || service.pilot?.username || '-' }}</el-descriptions-item>
      <el-descriptions-item label="描述" :span="2">{{ service.description }}</el-descriptions-item>
      <el-descriptions-item label="备注" :span="2">{{ service.remark || '-' }}</el-descriptions-item>
    </el-descriptions>

    <el-divider />

    <div class="actions">
      <template v-if="role === 'pilot' && service.status === 'open'">
        <el-button type="primary" @click="showBid = true">竞标报价</el-button>
      </template>
      <template v-if="role === 'client' && service.status === 'assigned'">
        <el-button type="success" @click="updateStatus('progress')">开始服务</el-button>
      </template>
      <template v-if="service.status === 'progress'">
        <el-button type="success" @click="updateStatus('completed')">完成服务</el-button>
      </template>
      <template v-if="service.status === 'open' && role === 'client'">
        <el-button type="danger" @click="updateStatus('cancelled')">取消</el-button>
      </template>
    </div>

    <el-divider />

    <h3>竞标列表</h3>
    <el-table :data="bids">
      <el-table-column label="飞手">
        <template #default="{ row }">{{ row.pilot?.nickname || row.pilot?.username }}</template>
      </el-table-column>
      <el-table-column prop="price" label="报价" width="120">
        <template #default="{ row }">¥{{ row.price }}</template>
      </el-table-column>
      <el-table-column prop="message" label="留言" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="bidStatusTagType(row.status)">{{ bidStatusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="时间" width="180" />
      <el-table-column label="操作" width="120" v-if="role === 'client' && service.status === 'open'">
        <template #default="{ row }">
          <el-button v-if="row.status === 'pending'" type="success" link @click="acceptBid(row)">接受</el-button>
        </template>
      </el-table-column>
    </el-table>

    <el-dialog v-model="showBid" title="竞标报价" width="400px">
      <el-form :model="bidForm" :rules="bidRules" ref="bidFormRef" label-width="80px">
        <el-form-item label="报价" prop="price">
          <el-input-number v-model="bidForm.price" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="留言">
          <el-input v-model="bidForm.message" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showBid = false">取消</el-button>
        <el-button type="primary" @click="submitBid" :loading="submitting">提交</el-button>
      </template>
    </el-dialog>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { useUserStore } from '@/stores/user'
import request from '@/utils/request'

const route = useRoute()
const userStore = useUserStore()
const role = computed(() => userStore.role)
const serviceId = route.params.id as string

const service = ref<AerialService | null>(null)
const bids = ref<ServiceBid[]>([])

const showBid = ref(false)
const submitting = ref(false)
const bidFormRef = ref<FormInstance>()

const bidForm = reactive({
  price: 0,
  message: ''
})

const bidRules: FormRules = {
  price: [{ required: true, message: '请输入报价', trigger: 'blur' }]
}

onMounted(() => {
  fetchService()
  fetchBids()
})

async function fetchService() {
  try {
    const res: any = await request.get(`/services/${serviceId}`)
    service.value = res.data
  } catch (e) {
    console.error(e)
  }
}

async function fetchBids() {
  try {
    const res: any = await request.get(`/services/${serviceId}/bids`)
    bids.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

async function submitBid() {
  if (!bidFormRef.value) return
  await bidFormRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        await request.post('/services/bid', { service_id: service.value?.id, ...bidForm })
        ElMessage.success('竞标成功')
        showBid.value = false
        fetchBids()
      } catch (e: any) {
        ElMessage.error(e.message || '竞标失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

async function acceptBid(bid: ServiceBid) {
  try {
    await request.post('/services/accept-bid', { bid_id: bid.id })
    ElMessage.success('已接受竞标')
    fetchService()
    fetchBids()
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  }
}

async function updateStatus(status: string) {
  try {
    await request.put('/services/status', { service_id: service.value?.id, status })
    ElMessage.success('操作成功')
    fetchService()
  } catch (e: any) {
    ElMessage.error(e.message || '操作失败')
  }
}

function statusText(status: string) {
  const map: Record<string, string> = {
    open: '招募中', assigned: '已分配', progress: '进行中',
    completed: '已完成', cancelled: '已取消'
  }
  return map[status] || status
}

function statusTagType(status: string) {
  const map: Record<string, string> = {
    open: 'success', assigned: 'warning', progress: 'primary',
    completed: 'info', cancelled: 'info'
  }
  return map[status] || ''
}

function bidStatusText(status: string) {
  const map: Record<string, string> = { pending: '待接受', accepted: '已接受', rejected: '已拒绝' }
  return map[status] || status
}

function bidStatusTagType(status: string) {
  const map: Record<string, string> = { pending: 'warning', accepted: 'success', rejected: 'info' }
  return map[status] || ''
}
</script>

<style scoped>
.actions {
  display: flex;
  gap: 10px;
}
</style>
