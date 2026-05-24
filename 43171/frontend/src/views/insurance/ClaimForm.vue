<template>
  <el-card>
    <template #header>
      <span>申请理赔</span>
    </template>
    <el-form :model="form" :rules="rules" ref="formRef" label-width="100px" style="max-width: 600px">
      <el-form-item label="订单ID" prop="order_id">
        <el-input-number v-model="form.order_id" :min="1" />
      </el-form-item>
      <el-form-item label="损坏描述" prop="damage_desc">
        <el-input v-model="form.damage_desc" type="textarea" :rows="4" />
      </el-form-item>
      <el-form-item label="预估费用" prop="estimated_cost">
        <el-input-number v-model="form.estimated_cost" :min="0" :precision="2" />
      </el-form-item>
      <el-form-item label="损坏照片">
        <el-upload
          :action="'/api/upload/damage'"
          :headers="{ Authorization: 'Bearer ' + token }"
          :show-file-list="false"
          :on-success="handleUpload"
          list-type="picture-card"
          multiple
        >
          <el-icon><Plus /></el-icon>
        </el-upload>
      </el-form-item>
      <el-form-item>
        <el-button type="primary" @click="handleSubmit" :loading="loading">提交</el-button>
        <el-button @click="$router.back()">取消</el-button>
      </el-form-item>
    </el-form>
  </el-card>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { useUserStore } from '@/stores/user'
import request from '@/utils/request'

const router = useRouter()
const userStore = useUserStore()
const token = computed(() => userStore.token)

const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive({
  order_id: 0,
  damage_desc: '',
  estimated_cost: 0,
  damage_images: ''
})

const rules: FormRules = {
  order_id: [{ required: true, message: '请输入订单ID', trigger: 'blur' }],
  damage_desc: [{ required: true, message: '请输入损坏描述', trigger: 'blur' }],
  estimated_cost: [{ required: true, message: '请输入预估费用', trigger: 'blur' }]
}

function handleUpload(res: any) {
  form.damage_images = res.data.url
}

async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await request.post('/insurance/claims', form)
        ElMessage.success('申请成功，请等待审核')
        router.push('/insurance')
      } catch (e: any) {
        ElMessage.error(e.message || '申请失败')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>
