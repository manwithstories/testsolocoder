<template>
  <div class="page-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>我的预约</span>
          <el-button type="primary" :icon="Plus" @click="$router.push('/appointments/create')">
            新建预约
          </el-button>
        </div>
      </template>

      <el-table :data="appointments" stripe v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column label="服务">
          <template #default="{ row }">
            {{ row.service?.name }}
          </template>
        </el-table-column>
        <el-table-column label="技师">
          <template #default="{ row }">
            {{ row.technician?.name }}
          </template>
        </el-table-column>
        <el-table-column label="日期" width="120">
          <template #default="{ row }">
            {{ formatDate(row.appointment_date) }}
          </template>
        </el-table-column>
        <el-table-column label="时间" width="140">
          <template #default="{ row }">
            {{ row.start_time }} - {{ row.end_time }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="160" fixed="right">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'confirmed' || row.status === 'pending'"
              type="warning"
              link
              @click="handleCancel(row)"
            >
              取消
            </el-button>
            <el-button
              v-if="row.status === 'confirmed' || row.status === 'pending'"
              type="primary"
              link
              @click="handleReschedule(row)"
            >
              改期
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination">
        <el-pagination
          v-model:current-page="page"
          v-model:page-size="pageSize"
          :total="total"
          :page-sizes="[10, 20, 50]"
          layout="total, sizes, prev, pager, next, jumper"
          @current-change="fetchList"
          @size-change="fetchList"
        />
      </div>
    </el-card>

    <el-dialog v-model="rescheduleVisible" title="改期" width="500px">
      <el-form :model="rescheduleForm" label-width="80px">
        <el-form-item label="日期" required>
          <el-date-picker
            v-model="rescheduleForm.appointment_date"
            type="date"
            placeholder="选择日期"
            :disabled-date="disabledDate"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="时间" required>
          <el-select
            v-model="rescheduleForm.start_time"
            placeholder="请选择时间"
            loading="slotsLoading"
          >
            <el-option
              v-for="slot in availableSlots"
              :key="slot"
              :label="slot"
              :value="slot.split('-')[0]"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rescheduleVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="handleSubmitReschedule">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { getMyAppointments, cancelAppointment, rescheduleAppointment } from '@/api/appointment'
import { getAvailableSlots } from '@/api/appointment'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import type { Appointment } from '@/types'

const loading = ref(false)
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const appointments = ref<Appointment[]>([])

const rescheduleVisible = ref(false)
const reschedulingId = ref<number | null>(null)
const submitting = ref(false)
const slotsLoading = ref(false)
const availableSlots = ref<string[]>([])
const rescheduleForm = ref({
  appointment_date: '',
  start_time: ''
})

const formatDate = (date: string) => dayjs(date).format('YYYY-MM-DD')

const disabledDate = (time: Date) => {
  return time.getTime() < Date.now() - 86400000
}

const getStatusType = (status: string) => {
  const types: Record<string, string> = {
    pending: 'info',
    confirmed: 'primary',
    paid: 'success',
    completed: 'success',
    cancelled: 'danger'
  }
  return types[status] || 'info'
}

const getStatusText = (status: string) => {
  const texts: Record<string, string> = {
    pending: '待确认',
    confirmed: '已确认',
    paid: '已支付',
    completed: '已完成',
    cancelled: '已取消'
  }
  return texts[status] || status
}

const fetchList = async () => {
  loading.value = true
  try {
    const res = await getMyAppointments({ page: page.value, page_size: pageSize.value })
    appointments.value = res.data.list
    total.value = res.data.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleCancel = async (row: Appointment) => {
  try {
    const { value } = await ElMessageBox.prompt('请输入取消原因', '取消预约', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      inputType: 'textarea'
    })
    await cancelAppointment({ id: row.id, cancel_reason: value })
    ElMessage.success('取消成功')
    fetchList()
  } catch (e) {
    if (e !== 'cancel') console.error(e)
  }
}

const handleReschedule = async (row: Appointment) => {
  reschedulingId.value = row.id
  rescheduleForm.value = {
    appointment_date: '',
    start_time: ''
  }
  rescheduleVisible.value = true
}

const handleSubmitReschedule = async () => {
  if (!rescheduleForm.value.appointment_date || !rescheduleForm.value.start_time) {
    ElMessage.warning('请填写必填项')
    return
  }

  submitting.value = true
  try {
    await rescheduleAppointment({
      id: reschedulingId.value!,
      appointment_date: rescheduleForm.value.appointment_date,
      start_time: rescheduleForm.value.start_time
    })
    ElMessage.success('改期成功')
    rescheduleVisible.value = false
    fetchList()
  } catch (e: any) {
    ElMessage.error(e.message || '改期失败')
  } finally {
    submitting.value = false
  }
}

onMounted(fetchList)
</script>

<style scoped lang="scss">
.page-container {
  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }

  .pagination {
    margin-top: 20px;
    display: flex;
    justify-content: flex-end;
  }
}
</style>
