<template>
  <el-card>
    <template #header>
      <span>添加飞行记录</span>
    </template>
    <el-form :model="form" :rules="rules" ref="formRef" label-width="100px" style="max-width: 600px">
      <el-form-item label="设备ID" prop="drone_id">
        <el-input-number v-model="form.drone_id" :min="1" />
      </el-form-item>
      <el-form-item label="关联订单">
        <el-input-number v-model="form.order_id" :min="0" />
      </el-form-item>
      <el-form-item label="关联服务">
        <el-input-number v-model="form.service_id" :min="0" />
      </el-form-item>
      <el-form-item label="起飞点">
        <el-input v-model="form.start_point" />
      </el-form-item>
      <el-form-item label="降落点">
        <el-input v-model="form.end_point" />
      </el-form-item>
      <el-form-item label="航线">
        <el-input v-model="form.route" type="textarea" :rows="2" />
      </el-form-item>
      <el-form-item label="最高高度(m)">
        <el-input-number v-model="form.altitude_max" :min="0" :precision="1" />
      </el-form-item>
      <el-form-item label="平均高度(m)">
        <el-input-number v-model="form.altitude_avg" :min="0" :precision="1" />
      </el-form-item>
      <el-form-item label="时长(分钟)" prop="duration">
        <el-input-number v-model="form.duration" :min="0" />
      </el-form-item>
      <el-form-item label="距离(km)">
        <el-input-number v-model="form.distance" :min="0" :precision="2" />
      </el-form-item>
      <el-form-item label="飞行日期" prop="flight_date">
        <el-date-picker v-model="form.flight_date" type="date" placeholder="选择日期" />
      </el-form-item>
      <el-form-item label="飞行日志">
        <el-input v-model="form.flight_log" type="textarea" :rows="3" />
      </el-form-item>
      <el-form-item label="备注">
        <el-input v-model="form.remark" type="textarea" :rows="2" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleSubmit" :loading="loading">提交</el-button>
        <el-button @click="$router.back()">取消</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import request from '@/utils/request'
import dayjs from 'dayjs'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  drone_id: 0,
  order_id: undefined as number | undefined,
  service_id: undefined as number | undefined,
  start_point: '',
  end_point: '',
  route: '',
  altitude_max: 0,
  altitude_avg: 0,
  duration: 0,
  distance: 0,
  flight_date: '',
  flight_log: '',
  remark: ''
})

const rules: FormRules = {
  drone_id: [{ required: true, message: '请输入设备ID', trigger: 'blur' }],
  duration: [{ required: true, message: '请输入飞行时长', trigger: 'blur' }],
  flight_date: [{ required: true, message: '请选择飞行日期', trigger: 'change' }]
}

async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await request.post('/flights', {
          ...form,
          flight_date: dayjs(form.flight_date as string).format('YYYY-MM-DD')
        })
        ElMessage.success('添加成功')
        router.push('/flights')
      } catch (e: any) {
        ElMessage.error(e.message || '添加失败')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>
