<template>
  <AppLayout>
    <div class="page">
      <h2>写评价</h2>
      <el-form :model="form" :rules="rules" ref="formRef" label-width="80px" @submit.prevent="onSubmit">
        <el-form-item label="订单号"><el-input :model-value="String(orderId)" disabled /></el-form-item>
        <el-form-item label="家政人员ID" prop="staff_id">
          <el-input-number v-model="form.staff_id" :min="1" />
        </el-form-item>
        <el-form-item label="评分" prop="rating">
          <el-rate v-model="form.rating" :max="5" />
        </el-form-item>
        <el-form-item label="评价内容">
          <el-input v-model="form.content" type="textarea" :rows="4" />
        </el-form-item>
        <el-form-item label="图片"><el-input v-model="form.images" placeholder="逗号分隔" /></el-form-item>
        <el-form-item>
          <el-button type="primary" native-type="submit" :loading="loading">提交</el-button>
        </el-form-item>
      </el-form>
    </div>
  </AppLayout>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import AppLayout from '../../components/AppLayout.vue'
import { createReview } from '../../api/review'
import { getOrder } from '../../api/order'

const route = useRoute()
const router = useRouter()
const orderId = Number(route.params.id)
const formRef = ref<FormInstance>()
const loading = ref(false)
const form = reactive({ order_id: orderId, staff_id: 0, rating: 5, content: '', images: '' })
const rules: FormRules = {
  staff_id: [{ required: true, type: 'number', min: 1, message: '请输入家政人员ID', trigger: 'blur' }],
  rating: [{ required: true, type: 'number', message: '请评分', trigger: 'change' }],
}

onMounted(async () => {
  const res = await getOrder(orderId)
  form.staff_id = (res.data as any).data.staff_id
})

async function onSubmit() {
  const valid = await formRef.value?.validate().catch(() => false)
  if (!valid) return
  loading.value = true
  try {
    await createReview(form)
    ElMessage.success('评价已提交')
    router.push('/orders')
  } catch (e: any) {
    ElMessage.error(e.message || '失败')
  } finally {
    loading.value = false
  }
}
</script>
