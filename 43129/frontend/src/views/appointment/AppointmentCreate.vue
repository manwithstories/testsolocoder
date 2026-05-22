<template>
  <div class="page-container">
    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>创建预约</span>
          <el-button @click="$router.back()">返回</el-button>
        </div>
      </template>

      <el-form :model="form" label-width="100px" style="max-width: 600px">
        <el-form-item label="选择技师" required>
          <el-select v-model="form.technician_id" placeholder="请选择技师" @change="onTechnicianChange">
            <el-option
              v-for="tech in technicians"
              :key="tech.id"
              :label="tech.name + ' - ' + tech.title"
              :value="tech.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="选择服务" required>
          <el-select v-model="form.service_id" placeholder="请选择服务" @change="onServiceChange">
            <el-option
              v-for="svc in services"
              :key="svc.id"
              :label="svc.name + ' - ¥' + svc.price"
              :value="svc.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="选择日期" required>
          <el-date-picker
            v-model="form.appointment_date"
            type="date"
            placeholder="选择日期"
            :disabled-date="disabledDate"
            value-format="YYYY-MM-DD"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="选择时间" required>
          <el-select
            v-model="form.start_time"
            placeholder="请先选择技师和日期"
            :disabled="!form.technician_id || !form.appointment_date"
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

        <el-form-item label="备注">
          <el-input v-model="form.remark" type="textarea" />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">
            提交预约
          </el-button>
          <el-button @click="$router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { getAllTechnicians } from '@/api/technician'
import { getAllServices } from '@/api/service'
import { getAvailableSlots, createAppointment } from '@/api/appointment'
import { getMyCustomer } from '@/api/auth'
import { ElMessage } from 'element-plus'
import dayjs from 'dayjs'
import type { Technician, Service, AvailableSlot } from '@/types'

const router = useRouter()
const technicians = ref<Technician[]>([])
const services = ref<Service[]>([])
const availableSlots = ref<string[]>([])
const slotsLoading = ref(false)
const submitting = ref(false)
const customerId = ref<number | null>(null)

const form = reactive({
  technician_id: null as number | null,
  service_id: null as number | null,
  appointment_date: '',
  start_time: '',
  remark: ''
})

const disabledDate = (time: Date) => {
  return time.getTime() < Date.now() - 86400000
}

const fetchTechnicians = async () => {
  try {
    const res = await getAllTechnicians()
    technicians.value = res.data
  } catch (e) {
    console.error(e)
  }
}

const fetchServices = async () => {
  try {
    const res = await getAllServices()
    services.value = res.data.filter(s => !s.is_package)
  } catch (e) {
    console.error(e)
  }
}

const fetchCustomerId = async () => {
  try {
    const res = await getMyCustomer()
    customerId.value = res.data.id
  } catch (e) {
    console.error(e)
  }
}

const fetchAvailableSlots = async () => {
  if (!form.technician_id || !form.appointment_date) return

  slotsLoading.value = true
  try {
    const duration = services.value.find(s => s.id === form.service_id)?.duration || 60
    const res: any = await getAvailableSlots(form.technician_id, {
      date: form.appointment_date,
      duration
    })
    availableSlots.value = res.data.slots || []
  } catch (e) {
    console.error(e)
  } finally {
    slotsLoading.value = false
  }
}

const onTechnicianChange = () => {
  form.start_time = ''
  availableSlots.value = []
  if (form.appointment_date) {
    fetchAvailableSlots()
  }
}

const onServiceChange = () => {
  form.start_time = ''
  if (form.technician_id && form.appointment_date) {
    fetchAvailableSlots()
  }
}

const handleSubmit = async () => {
  if (!form.technician_id || !form.service_id || !form.appointment_date || !form.start_time) {
    ElMessage.warning('请填写必填项')
    return
  }

  submitting.value = true
  try {
    await createAppointment({
      customer_id: customerId.value!,
      technician_id: form.technician_id,
      service_id: form.service_id,
      appointment_date: form.appointment_date,
      start_time: form.start_time,
      remark: form.remark
    })
    ElMessage.success('预约成功')
    router.push('/my-appointments')
  } catch (e: any) {
    ElMessage.error(e.message || '预约失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchTechnicians()
  fetchServices()
  fetchCustomerId()
})
</script>

<style scoped lang="scss">
.page-container {
  .card-header {
    display: flex;
    align-items: center;
    justify-content: space-between;
  }
}
</style>
