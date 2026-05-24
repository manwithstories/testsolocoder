<template>
  <el-card>
    <template #header>
      <span>发布航拍需求</span>
    </template>
    <el-form :model="form" :rules="rules" ref="formRef" label-width="100px" style="max-width: 600px">
      <el-form-item label="需求标题" prop="title">
        <el-input v-model="form.title" />
      </el-form-item>
      <el-form-item label="详细描述" prop="description">
        <el-input v-model="form.description" type="textarea" :rows="4" />
      </el-form-item>
      <el-form-item label="区域" prop="region">
        <el-input v-model="form.region" />
      </el-form-item>
      <el-form-item label="详细地址">
        <el-input v-model="form.address" />
      </el-form-item>
      <el-form-item label="服务日期">
        <el-date-picker v-model="form.service_date" type="date" placeholder="选择日期" />
      </el-form-item>
      <el-form-item label="服务时间">
        <el-input v-model="form.service_time" placeholder="如: 09:00-12:00" />
      </el-form-item>
      <el-form-item label="时长(分钟)">
        <el-input-number v-model="form.duration" :min="0" />
      </el-form-item>
      <el-form-item label="预算下限">
        <el-input-number v-model="form.budget_min" :min="0" :precision="2" />
      </el-form-item>
      <el-form-item label="预算上限">
        <el-input-number v-model="form.budget_max" :min="0" :precision="2" />
      </el-form-item>
      <el-form-item label="备注">
        <el-input v-model="form.remark" type="textarea" :rows="2" />
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleSubmit" :loading="loading">发布</el-button>
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
  title: '',
  description: '',
  region: '',
  address: '',
  service_date: '',
  service_time: '',
  duration: 0,
  budget_min: 0,
  budget_max: 0,
  remark: ''
})

const rules: FormRules = {
  title: [{ required: true, message: '请输入需求标题', trigger: 'blur' }],
  description: [{ required: true, message: '请输入详细描述', trigger: 'blur' }],
  region: [{ required: true, message: '请输入区域', trigger: 'blur' }]
}

async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const data = {
          ...form,
          service_date: form.service_date ? dayjs(form.service_date as string).toISOString() : null
        }
        await request.post('/services', data)
        ElMessage.success('发布成功')
        router.push('/services')
      } catch (e: any) {
        ElMessage.error(e.message || '发布失败')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>
