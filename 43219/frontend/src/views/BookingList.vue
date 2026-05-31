<template>
  <AppLayout>
    <div class="page">
      <h2>预约记录</h2>
      <div class="row" style="margin-bottom:16px">
        <el-select v-model="status" placeholder="全部状态" clearable style="width:160px" @change="load">
          <el-option value="pending" label="待确认" />
          <el-option value="confirmed" label="已确认" />
          <el-option value="canceled" label="已取消" />
          <el-option value="completed" label="已完成" />
        </el-select>
        <el-button v-if="userStore.role==='customer'" type="primary" @click="$router.push('/booking/new')">新建预约</el-button>
      </div>
      <el-table :data="list" border stripe>
        <el-table-column prop="id" label="编号" width="80" />
        <el-table-column label="服务" width="180">
          <template #default="{ row }">{{ row.service?.name || '-' }}</template>
        </el-table-column>
        <el-table-column label="家政人员" width="140">
          <template #default="{ row }">{{ row.staff?.real_name || row.staff?.username || '-' }}</template>
        </el-table-column>
        <el-table-column prop="start_at" label="开始时间" width="170" />
        <el-table-column prop="end_at" label="结束时间" width="170" />
        <el-table-column prop="price" label="金额" width="100" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="tagType(row.status)">{{ statusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="320">
          <template #default="{ row }">
            <el-button size="small" @click="$router.push(`/staff/${row.staff_id}`)" v-if="row.staff_id">看师傅</el-button>
            <el-button size="small" type="primary" v-if="row.status==='pending'" @click="confirm(row)">确认</el-button>
            <el-button size="small" v-if="row.status==='pending' || row.status==='confirmed'" @click="openReschedule(row)">改期</el-button>
            <el-button size="small" type="danger" v-if="row.status==='pending' || row.status==='confirmed'" @click="cancel(row)">取消</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>

    <el-dialog v-model="resVisible" title="改期">
      <el-form :model="resForm" label-width="80px">
        <el-form-item label="新开始">
          <el-date-picker v-model="resForm.start_at" type="datetime" placeholder="开始时间" value-format="YYYY-MM-DDTHH:mm:ss" />
        </el-form-item>
        <el-form-item label="新结束">
          <el-date-picker v-model="resForm.end_at" type="datetime" placeholder="结束时间" value-format="YYYY-MM-DDTHH:mm:ss" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="resVisible=false">取消</el-button>
        <el-button type="primary" @click="submitReschedule">提交</el-button>
      </template>
    </el-dialog>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, onMounted, reactive } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import AppLayout from '../components/AppLayout.vue'
import { listBookings, confirmBooking, rescheduleBooking, cancelBooking, type Booking } from '../api/booking'
import { useUserStore } from '../stores/user'

const userStore = useUserStore()
const list = ref<Booking[]>([])
const status = ref('')
const resVisible = ref(false)
const currentId = ref<number | null>(null)
const resForm = reactive({ start_at: '', end_at: '' })

async function load() {
  const params: Record<string, string> = {}
  if (status.value) params.status = status.value
  const res = await listBookings(params)
  list.value = (res.data as any).data || []
}

function statusLabel(s: string) {
  return { pending: '待确认', confirmed: '已确认', canceled: '已取消', rejected: '已拒绝', completed: '已完成' } as Record<string, string>[s] || s
}
function tagType(s: string) {
  return { pending: 'warning', confirmed: 'success', canceled: 'info', rejected: 'danger', completed: 'primary' } as Record<string, string>[s] || ''
}

async function confirm(row: Booking) {
  await confirmBooking(row.id)
  ElMessage.success('已确认')
  load()
}

function openReschedule(row: Booking) {
  currentId.value = row.id
  resForm.start_at = ''
  resForm.end_at = ''
  resVisible.value = true
}

async function submitReschedule() {
  if (!resForm.start_at || !resForm.end_at) {
    ElMessage.warning('请填写完整时间')
    return
  }
  await rescheduleBooking(currentId.value!, resForm)
  ElMessage.success('改期请求已提交')
  resVisible.value = false
  load()
}

async function cancel(row: Booking) {
  await ElMessageBox.confirm('确认取消预约?', '提示', { type: 'warning' })
  await cancelBooking(row.id)
  ElMessage.success('已取消')
  load()
}

onMounted(load)
</script>
