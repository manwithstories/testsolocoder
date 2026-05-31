<template>
  <AppLayout>
    <div class="page">
      <h2>新建工单</h2>
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px" @submit.prevent="onSubmit">
        <el-form-item label="类型" prop="type">
          <el-radio-group v-model="form.type">
            <el-radio value="complaint">投诉</el-radio>
            <el-radio value="refund">退款</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="订单号">
          <el-input-number v-model="form.order_id" :min="0" />
        </el-form-item>
        <el-form-item label="家政人员ID">
          <el-input-number v-model="form.staff_id" :min="0" />
        </el-form-item>
        <el-form-item label="标题" prop="title"><el-input v-model="form.title" /></el-form-item>
        <el-form-item label="内容" prop="content"><el-input v-model="form.content" type="textarea" :rows="5" /></el-form-item>
        <el-form-item>
          <el-button type="primary" native-type="submit" :loading="loading">提交</el-button>
        </el-form-item>
      </el-form>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import AppLayout from '../../components/AppLayout.vue'
import { createTicket } from '../../api/ticket'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)
const form = reactive({
  type: 'complaint' as 'complaint' | 'refund',
  order_id: 0,
  staff_id: 0,
  title: '',
  content: '',
})
const rules: FormRules = {
  type: [{ required: true, message: '请选择类型', trigger: 'change' }],
  title: [{ required: true, message: '请填写标题', trigger: 'blur' }],
  content: [{ required: true, message: '请填写内容', trigger: 'blur' }],
}

async function onSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  loading.value = true
  try {
    await createTicket({
      type: form.type,
      title: form.title,
      content: form.content,
      order_id: form.order_id > 0 ? form.order_id : undefined,
      staff_id: form.staff_id > 0 ? form.staff_id : undefined,
    })
    ElMessage.success('工单已提交')
    router.push('/tickets')
  } catch (e: any) {
    ElMessage.error(e.message || '失败')
  } finally {
    loading.value = false
  }
}
</script>
