<template>
  <div class="review-form">
    <el-page-header @back="router.back()">
      <template #content>
        <span>发表评价</span>
      </template>
    </el-page-header>

    <el-card shadow="never" style="margin-top: 16px">
      <el-form ref="formRef" :model="form" :rules="rules" label-width="100px" style="max-width: 720px">
        <el-form-item label="订单号" prop="orderNo">
          <el-input v-model="form.orderNo" placeholder="请输入订单号" />
        </el-form-item>
        <el-form-item label="产品评分" prop="productScore">
          <el-rate v-model="form.productScore" show-score text-color="#ff9900" :max="5" />
        </el-form-item>
        <el-form-item label="服务评分" prop="serviceScore">
          <el-rate v-model="form.serviceScore" show-score text-color="#ff9900" :max="5" />
        </el-form-item>
        <el-form-item label="评价内容" prop="content">
          <el-input
            v-model="form.content"
            type="textarea"
            :rows="4"
            placeholder="请输入评价内容"
            maxlength="500"
            show-word-limit
          />
        </el-form-item>
        <el-form-item label="上传图片">
          <el-upload
            list-type="picture-card"
            :auto-upload="false"
            :file-list="imageFileList"
            :on-change="handleImageChange"
            :on-remove="handleImageRemove"
            accept="image/*"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
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
import type { FormInstance, FormRules, UploadFile } from 'element-plus'
import { ElMessage } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import { createReview, type ReviewFormData } from '@/api/review'

const router = useRouter()
const formRef = ref<FormInstance>()
const submitting = ref(false)
const imageFileList = ref<UploadFile[]>([])

const form = reactive<ReviewFormData & { orderNo?: string }>({
  productScore: 5,
  serviceScore: 5,
  content: '',
  images: []
})

const rules: FormRules = {
  productScore: [{ required: true, message: '请选择产品评分', trigger: 'change' }],
  serviceScore: [{ required: true, message: '请选择服务评分', trigger: 'change' }],
  content: [
    { required: true, message: '请输入评价内容', trigger: 'blur' },
    { min: 5, message: '评价内容不少于 5 个字', trigger: 'blur' }
  ]
}

function handleImageChange(file: UploadFile) {
  const reader = new FileReader()
  reader.onload = (e) => {
    const url = e.target?.result as string
    if (!form.images) form.images = []
    form.images.push(url)
    imageFileList.value.push({ name: file.name, url, uid: file.uid })
  }
  if (file.raw) reader.readAsDataURL(file.raw)
}

function handleImageRemove(file: UploadFile) {
  const idx = imageFileList.value.findIndex((f) => f.uid === file.uid)
  if (idx > -1) {
    imageFileList.value.splice(idx, 1)
    form.images?.splice(idx, 1)
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (!valid) return
    submitting.value = true
    try {
      await createReview(form)
      ElMessage.success('评价提交成功')
      router.back()
    } finally {
      submitting.value = false
    }
  })
}
</script>

<style lang="scss" scoped>
.review-form {
}
</style>
