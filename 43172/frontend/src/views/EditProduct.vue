<template>
  <div class="page-container">
    <div class="card">
      <h2 class="section-title">编辑商品</h2>
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
              placeholder="请输入详细描述"
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
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">
            保存修改
          </el-button>
          <el-button @click="$router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { productApi } from '@/api/product'
import { ElMessage } from 'element-plus'
import type { FormInstance, FormRules, UploadFile } from 'element-plus'
import { CATEGORY_OPTIONS } from '@/types'
import { Plus } from '@element-plus/icons-vue'

const route = useRoute()
const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)
const fileList = ref<UploadFile[]>([])
const filesToUpload = ref<File[]>([])

const productId = Number(route.params.id)

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

const loadProduct = async () => {
  try {
    const res = await productApi.getProduct(productId)
    if (res.code === 200 && res.data) {
      Object.assign(form, {
        title: res.data.title,
        category: res.data.category,
        brand_name: res.data.brand_name || '',
        description: res.data.description,
        condition: res.data.condition || '',
        color: res.data.color || '',
        size: res.data.size || '',
        material: res.data.material || '',
        original_price: res.data.original_price || 0,
        price: res.data.price,
        stock: res.data.stock
      })
    }
  } catch (error) {
    console.error('Load product error:', error)
  }
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

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const res = await productApi.updateProduct(productId, form)
        if (res.code === 200) {
          if (filesToUpload.value.length > 0) {
            const formData = new FormData()
            filesToUpload.value.forEach(file => {
              formData.append('images', file)
            })
            await productApi.uploadImages(productId, formData)
          }

          ElMessage.success('商品已更新')
          router.push('/seller/products')
        }
      } catch (error) {
        console.error('Update product error:', error)
      } finally {
        loading.value = false
      }
    }
  })
}

onMounted(() => {
  loadProduct()
})
</script>

<style lang="scss" scoped>
</style>
