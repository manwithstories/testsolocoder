<template>
  <el-card>
    <template #header>发布手表</template>
    <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
      <el-form-item label="品牌" prop="brand"><el-input v-model="form.brand" /></el-form-item>
      <el-form-item label="型号" prop="model"><el-input v-model="form.model" /></el-form-item>
      <el-form-item label="参考号"><el-input v-model="form.reference_no" /></el-form-item>
      <el-form-item label="年份"><el-input-number v-model="form.year" :min="1900" :max="2100" /></el-form-item>
      <el-form-item label="机芯"><el-input v-model="form.movement" /></el-form-item>
      <el-form-item label="表径(mm)"><el-input-number v-model="form.case_size_mm" :precision="1" /></el-form-item>
      <el-form-item label="表壳"><el-input v-model="form.case_material" /></el-form-item>
      <el-form-item label="表盘"><el-input v-model="form.dial_color" /></el-form-item>
      <el-form-item label="表带"><el-input v-model="form.bracelet" /></el-form-item>
      <el-form-item label="成色"><el-input v-model="form.condition" /></el-form-item>
      <el-form-item label="描述"><el-input v-model="form.description" type="textarea" :rows="3" /></el-form-item>
      <el-form-item label="价格" prop="price"><el-input-number v-model="form.price" :min="0" :precision="2" /></el-form-item>
      <el-form-item>
        <el-button type="primary" :loading="loading" @click="submit">发布</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import request from '@/utils/request'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)
const form = reactive({
  brand: '', model: '', reference_no: '', year: 2024, movement: '', case_size_mm: 40,
  case_material: '', dial_color: '', bracelet: '', condition: '', description: '', price: 0
})
const rules: FormRules = {
  brand: [{ required: true }],
  model: [{ required: true }],
  price: [{ required: true }]
}
async function submit() {
  await formRef.value?.validate()
  loading.value = true
  try {
    const res: any = await request.post('/watches', form)
    ElMessage.success('发布成功')
    router.push(`/watches/${res.id}`)
  } finally {
    loading.value = false
  }
}
</script>
