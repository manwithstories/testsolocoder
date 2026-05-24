<template>
  <div class="ticket-form">
    <el-page-header @back="router.back()">
      <template #content>
        <span>创建工单</span>
      </template>
    </el-page-header>

    <el-card shadow="never" style="margin-top: 16px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px" style="max-width: 720px">
        <el-form-item label="工单标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入工单标题" maxlength="50" show-word-limit />
        </el-form-item>
        <el-form-item label="工单类型" prop="type">
          <el-radio-group v-model="form.type">
            <el-radio v-for="item in typeOptions" :key="item.value" :value="item.value">{{ item.label }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="优先级" prop="priority">
          <el-radio-group v-model="form.priority">
            <el-radio v-for="item in priorityOptions" :key="item.value" :value="item.value">{{ item.label }}</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="关联订单号">
          <el-input v-model="form.relatedOrderNo" placeholder="请输入关联订单号（可选）" />
        </el-form-item>
        <el-form-item label="问题描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="5"
            placeholder="请详细描述问题"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">提交</el-button>
          <el-button @click="router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import type { FormInstance, FormRules } from 'element-plus'
import { ElMessage } from 'element-plus'
import { createTicket, type TicketFormData } from '@/api/ticket'

const router = useRouter()
const formRef = ref<FormInstance>()
const submitting = ref(false)

const typeOptions = [
  { label: '产品问题', value: 1 },
  { label: '安装问题', value: 2 },
  { label: '售后服务', value: 3 },
  { label: '投诉建议', value: 4 },
  { label: '其他', value: 99 }
]

const priorityOptions = [
  { label: '低', value: 1 },
  { label: '中', value: 2 },
  { label: '高', value: 3 },
  { label: '紧急', value: 4 }
]

const form = reactive<TicketFormData & { relatedOrderNo?: string }>({
  title: '',
  description: '',
  type: 1,
  priority: 2
})

const rules: FormRules = {
  title: [{ required: true, message: '请输入工单标题', trigger: 'blur' }],
  type: [{ required: true, message: '请选择工单类型', trigger: 'change' }],
  priority: [{ required: true, message: '请选择优先级', trigger: 'change' }],
  description: [
    { required: true, message: '请输入问题描述', trigger: 'blur' },
    { min: 10, message: '问题描述不少于 10 个字', trigger: 'blur' }
  ]
}

async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      await createTicket(form)
      ElMessage.success('工单创建成功')
      router.back()
    } finally {
      submitting.value = false
    }
  })
}
</script>

<style lang="scss" scoped>
.ticket-form {
}
</style>
