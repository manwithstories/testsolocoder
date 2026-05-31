<template>
  <AppLayout>
    <div class="page">
      <h2>预约服务</h2>
      <div class="card">
        <el-form :model="form" :rules="rules" ref="formRef" label-width="100px" @submit.prevent="onSubmit">
          <el-form-item label="选择服务" prop="service_id">
            <el-select v-model="form.service_id" placeholder="请选择服务" style="width:100%" filterable>
              <el-option v-for="s in services" :key="s.id" :value="s.id"
                :label="`${s.name} (¥${s.min_price}-¥${s.max_price}, ${s.duration}分钟)`" />
            </el-select>
          </el-form-item>
          <el-form-item label="开始时间" prop="start_at">
            <el-date-picker v-model="form.start_at" type="datetime" placeholder="开始时间"
              value-format="YYYY-MM-DDTHH:mm:ss" style="width:100%" />
          </el-form-item>
          <el-form-item label="结束时间" prop="end_at">
            <el-date-picker v-model="form.end_at" type="datetime" placeholder="结束时间"
              value-format="YYYY-MM-DDTHH:mm:ss" style="width:100%" />
          </el-form-item>
          <el-form-item label="服务地址" prop="address">
            <el-input v-model="form.address" placeholder="详细地址" />
          </el-form-item>
          <el-form-item label="备注">
            <el-input v-model="form.remark" type="textarea" :rows="3" />
          </el-form-item>
          <el-form-item label="指定人员">
            <el-select v-model="form.staff_id" placeholder="留空由系统匹配" clearable style="width:100%">
              <el-option v-for="s in staffs" :key="s.id" :value="s.id"
                :label="`${s.real_name || s.username} (评分${s.rating})`" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button type="primary" native-type="submit" :loading="loading">提交预约</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import AppLayout from '../components/AppLayout.vue'
import { listServices, type ServiceItem } from '../../api/service'
import { listStaff } from '../../api/user'
import { createBooking } from '../../api/booking'

const route = useRoute()
const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)
const services = ref<ServiceItem[]>([])
const staffs = ref<any[]>([])
const form = reactive({
  service_id: Number(route.query.service_id) || undefined as any,
  start_at: '',
  end_at: '',
  address: '',
  remark: '',
  staff_id: undefined as number | undefined,
})

const rules: FormRules = {
  service_id: [{ required: true, type: 'number', message: '请选择服务', trigger: 'change' }],
  start_at: [{ required: true, message: '请选择开始时间', trigger: 'change' }],
  end_at: [
    { required: true, message: '请选择结束时间', trigger: 'change' },
    {
      validator: (_r, v, cb) => {
        if (!form.start_at) return cb()
        if (new Date(v) <= new Date(form.start_at)) return cb(new Error('结束时间必须晚于开始时间'))
        cb()
      },
      trigger: 'change',
    },
  ],
  address: [{ required: true, message: '请填写服务地址', trigger: 'blur' }],
}

async function load() {
  const [sRes, uRes] = await Promise.all([listServices(), listStaff()])
  services.value = (sRes.data as any).data || []
  staffs.value = (uRes.data as any).data || []
}

async function onSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  loading.value = true
  try {
    await createBooking({
      service_id: form.service_id,
      start_at: form.start_at,
      end_at: form.end_at,
      address: form.address,
      remark: form.remark,
      staff_id: form.staff_id,
    })
    ElMessage.success('预约成功,等待确认')
    router.push('/bookings')
  } catch (e: any) {
    ElMessage.error(e.message || '预约失败')
  } finally {
    loading.value = false
  }
}

onMounted(load)
</script>
