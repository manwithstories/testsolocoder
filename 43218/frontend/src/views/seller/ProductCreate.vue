<template>
  <div class="product-create-page">
    <div class="page-header">
      <h2 class="page-title">发布商品</h2>
    </div>

    <el-card>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        style="max-width: 800px"
      >
        <el-form-item label="商品标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入商品标题" maxlength="200" show-word-limit />
        </el-form-item>

        <el-form-item label="商品分类" prop="category">
          <el-select v-model="form.category" placeholder="请选择分类">
            <el-option
              v-for="cat in categories"
              :key="cat"
              :label="cat"
              :value="cat"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="品牌" prop="brand">
          <el-input v-model="form.brand" placeholder="请输入品牌" maxlength="100" />
        </el-form-item>

        <el-form-item label="型号" prop="model">
          <el-input v-model="form.model" placeholder="请输入型号" maxlength="100" />
        </el-form-item>

        <el-form-item label="成色等级" prop="condition">
          <el-radio-group v-model="form.condition">
            <el-radio
              v-for="cond in conditions"
              :key="cond"
              :label="cond"
            >
              {{ cond }}
            </el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="售价" prop="price">
          <el-input-number
            v-model="form.price"
            :min="0"
            :precision="2"
            :step="10"
            placeholder="请输入售价"
          />
          <span class="form-hint">元</span>
        </el-form-item>

        <el-form-item label="原价">
          <el-input-number
            v-model="form.originalPrice"
            :min="0"
            :precision="2"
            :step="10"
            placeholder="请输入原价（选填）"
          />
          <span class="form-hint">元</span>
        </el-form-item>

        <el-form-item label="质保期限" prop="warrantyDays">
          <el-input-number
            v-model="form.warrantyDays"
            :min="0"
            :step="7"
            placeholder="请输入质保天数"
          />
          <span class="form-hint">天（0表示无质保）</span>
        </el-form-item>

        <el-form-item label="商品描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="6"
            placeholder="请详细描述商品情况，包括使用情况、配件、维修历史等"
            maxlength="2000"
            show-word-limit
          />
        </el-form-item>

        <el-form-item label="商品图片">
          <el-upload
            :action="uploadUrl"
            :headers="uploadHeaders"
            list-type="picture-card"
            :limit="6"
            :file-list="fileList"
            :on-preview="handlePreview"
            :on-remove="handleRemove"
            :on-success="handleSuccess"
            :on-error="handleError"
            accept="image/*"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
          <div class="form-hint">最多上传6张图片，建议尺寸800x800px</div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" size="large" :loading="submitting" @click="handleSubmit">
            提交审核
          </el-button>
          <el-button size="large" @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-dialog v-model="previewVisible" title="图片预览" width="800px">
      <img :src="previewImage" style="width: 100%" />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules, type UploadFile, type UploadUserFile } from 'element-plus'
import { useUserStore } from '@/store/user'
import { productApi } from '@/api'
import { ProductCategories, ProductConditions } from '@/types'

const router = useRouter()
const userStore = useUserStore()

const formRef = ref<FormInstance>()
const submitting = ref(false)
const previewVisible = ref(false)
const previewImage = ref('')
const fileList = ref<UploadUserFile[]>([])

const uploadUrl = '/api/v1/upload/multiple'
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${userStore.token}`
}))

const categories = ProductCategories
const conditions = ProductConditions

const form = reactive({
  title: '',
  category: '',
  brand: '',
  model: '',
  condition: '',
  price: 0,
  originalPrice: 0,
  warrantyDays: 30,
  description: '',
  images: ''
})

const rules: FormRules = {
  title: [
    { required: true, message: '请输入商品标题', trigger: 'blur' },
    { min: 2, max: 200, message: '标题长度在2-200个字符', trigger: 'blur' }
  ],
  category: [
    { required: true, message: '请选择商品分类', trigger: 'change' }
  ],
  brand: [
    { required: true, message: '请输入品牌', trigger: 'blur' }
  ],
  model: [
    { required: true, message: '请输入型号', trigger: 'blur' }
  ],
  condition: [
    { required: true, message: '请选择成色等级', trigger: 'change' }
  ],
  price: [
    { required: true, message: '请输入售价', trigger: 'blur' },
    { type: 'number', min: 0, message: '售价不能为负数', trigger: 'blur' }
  ],
  description: [
    { required: true, message: '请输入商品描述', trigger: 'blur' }
  ]
}

function handlePreview(file: UploadFile) {
  previewImage.value = file.url!
  previewVisible.value = true
}

function handleRemove(file: UploadFile) {
  const index = fileList.value.findIndex(item => item.uid === file.uid)
  if (index > -1) {
    fileList.value.splice(index, 1)
  }
}

function handleSuccess(response: any, _file: UploadFile, _fileList: UploadUserFile[]) {
  if (response.code === 0) {
    ElMessage.success('上传成功')
  } else {
    ElMessage.error(response.message || '上传失败')
  }
}

function handleError() {
  ElMessage.error('上传失败')
}

async function handleSubmit() {
  if (!formRef.value) return

  const images: string[] = []
  fileList.value.forEach(file => {
    if (file.response?.data?.url) {
      images.push(file.response.data.url)
    } else if (file.url) {
      images.push(file.url)
    }
  })

  form.images = JSON.stringify(images)

  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        await productApi.create({
          title: form.title,
          category: form.category,
          brand: form.brand,
          model: form.model,
          condition: form.condition,
          price: form.price,
          originalPrice: form.originalPrice || undefined,
          description: form.description,
          warrantyDays: form.warrantyDays,
          images: form.images
        })
        ElMessage.success('商品已提交审核，审核通过后将上架展示')
        router.push('/seller/products')
      } catch (error: any) {
        ElMessage.error(error.message || '发布失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

function handleReset() {
  formRef.value?.resetFields()
  fileList.value = []
}

import { computed } from 'vue'
</script>

<style lang="scss" scoped>
.product-create-page {
  :deep(.el-upload--picture-card) {
    width: 100px;
    height: 100px;
    line-height: 100px;
  }

  .form-hint {
    color: var(--text-lighter-color);
    font-size: 13px;
    margin-left: 8px;
  }
}
</style>
