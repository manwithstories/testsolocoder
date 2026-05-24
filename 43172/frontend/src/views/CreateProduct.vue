<template>
  <div class="page-container">
    <div class="card">
      <h2 class="section-title">发布商品</h2>
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
      >
        <el-form-item label="商品标题" prop="title">
          <el-input v-model="form.title" placeholder="请输入商品标题" maxlength="200" />
        </el-form-item>
        <el-form-item label="商品分类" prop="category">
          <el-select v-model="form.category" placeholder="请选择分类" style="width: 100%">
            <el-option
              v-for="cat in CATEGORY_OPTIONS"
              :key="cat.value"
              :label="cat.label"
              :value="cat.value"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="品牌">
          <el-input v-model="form.brand_name" placeholder="请输入品牌名称" />
        </el-form-item>
        <el-form-item label="商品描述" prop="description">
          <el-input
              v-model="form.description"
              type="textarea"
              :rows="6"
              placeholder="请输入详细描述，至少10个字符"
            />
        </el-form-item>
        <el-form-item label="成色">
          <el-input v-model="form.condition" placeholder="请输入成色描述" />
        </el-form-item>
        <el-form-item label="颜色">
          <el-input v-model="form.color" placeholder="请输入颜色" />
        </el-form-item>
        <el-form-item label="尺寸">
          <el-input v-model="form.size" placeholder="请输入尺寸" />
        </el-form-item>
        <el-form-item label="材质">
          <el-input v-model="form.material" placeholder="请输入材质" />
        </el-form-item>
        <el-form-item label="原价">
          <el-input-number
              v-model="form.original_price"
              :min="0"
              :precision="2"
              :step="100"
              style="width: 100%"
            />
        </el-form-item>
        <el-form-item label="售价" prop="price">
          <el-input-number
              v-model="form.price"
              :min="0.01"
              :precision="2"
              :step="100"
              style="width: 100%"
            />
        </el-form-item>
        <el-form-item label="库存数量" prop="stock">
          <el-input-number
              v-model="form.stock"
              :min="1"
              :step="1"
              style="width: 100%"
            />
        </el-form-item>
        <el-form-item label="商品图片">
          <el-upload
            :auto-upload="false"
            :limit="10"
            list-type="picture-card"
            :on-change="handleImageChange"
            :on-remove="handleImageRemove"
            :file-list="fileList"
            accept="image/*"
          >
            <el-icon><Plus /></el-icon>
          </el-upload>
          <div class="upload-tip">支持 jpg、png、gif 格式，最多10张</div>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">
            保存草稿
          </el-button>
          <el-button type="success" :loading="loading" @click="handleSubmitAndPublish">
            保存并上架
          </el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { productApi } from '@/api/product'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules, UploadFile } from 'element-plus'
import { CATEGORY_OPTIONS } from '@/types'
import { Plus } from '@element-plus/icons-vue'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)
const fileList = ref<UploadFile[]>([])
const filesToUpload = ref<File[]>([])

const form = reactive({
  title: '',
  category: '',
  brand_name: '',
  description: '',
  condition: '',
  color: '',
  size: '',
  material: '',
  original_price: 0,
  price: 0,
  stock: 1
})

const rules: FormRules = {
  title: [
    { required: true, message: '请输入商品标题', trigger: 'blur' },
    { min: 5, max: 200, message: '标题长度在 5 到 200 个字符', trigger: 'blur' }
  ],
  category: [
    { required: true, message: '请选择商品分类', trigger: 'change' }
  ],
  description: [
    { required: true, message: '请输入商品描述', trigger: 'blur' },
    { min: 10, message: '描述至少10个字符', trigger: 'blur' }
  ],
  price: [
    { required: true, message: '请输入售价', trigger: 'blur' }
  ],
  stock: [
    { required: true, message: '请输入库存数量', trigger: 'blur' }
  ]
}

const handleImageChange = (file: UploadFile) => {
  if (file.raw) {
    filesToUpload.value.push(file.raw)
  }
}

const handleImageRemove = (file: UploadFile) => {
  const index = filesToUpload.value.indexOf(file.raw!)
  if (index > -1) {
    filesToUpload.value.splice(index, 1)
  }
}

const handleSubmit = async (publish: boolean) => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const res = await productApi.createProduct({
          title: form.title,
          description: form.description,
          category: form.category,
          brand_name: form.brand_name || undefined,
          original_price: form.original_price || undefined,
          price: form.price,
          condition: form.condition || undefined,
          color: form.color || undefined,
          size: form.size || undefined,
          material: form.material || undefined,
          stock: form.stock
        })

        if (res.code === 201 && res.data) {
          const productId = (res.data as any).id
          
          if (filesToUpload.value.length > 0 && productId) {
            const formData = new FormData()
            filesToUpload.value.forEach(file => {
              formData.append('images', file)
            })
            await productApi.uploadImages(productId, formData)
          }

          if (publish) {
            await productApi.updateProductStatus(productId, 'on_sale')
          }

          ElMessage.success(publish ? '商品已上架' : '草稿已保存')
          router.push('/seller/products')
        }
      } catch (error) {
        console.error('Create product error:', error)
      } finally {
        loading.value = false
      }
    }
  })
}

const handleSubmitAndPublish = () => handleSubmit(true)
</script>

<style lang="scss" scoped>
.upload-tip {
  font-size: 12px;
  color: var(--text-light);
  margin-top: 4px;
}
</style>
