<template>
  <div class="product-form">
    <el-page-header :icon="ArrowLeft" @back="goBack">
      <template #content>
        <span>{{ isEdit ? '编辑产品' : '新增产品' }}</span>
      </template>
    </el-page-header>

    <el-card v-loading="loading" class="form-card" shadow="never">
      <el-form
        ref="formRef"
        :model="formData"
        :rules="rules"
        label-width="100px"
        label-position="right"
      >
        <el-row :gutter="20">
          <el-col :xs="24" :md="12">
            <el-form-item label="产品名称" prop="name">
              <el-input v-model="formData.name" placeholder="请输入产品名称" maxlength="64" show-word-limit />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="12">
            <el-form-item label="SKU" prop="sku">
              <el-input v-model="formData.sku" placeholder="请输入SKU" maxlength="32" />
            </el-form-item>
          </el-col>

          <el-col :xs="24" :md="12">
            <el-form-item label="分类" prop="category">
              <el-select v-model="formData.category" placeholder="请选择分类" style="width: 100%">
                <el-option
                  v-for="cat in categories"
                  :key="cat.id"
                  :label="cat.name"
                  :value="cat.name"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="12">
            <el-form-item label="价格" prop="price">
              <el-input-number
                v-model="formData.price"
                :min="0"
                :precision="2"
                :step="1"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>

          <el-col :xs="24" :md="12">
            <el-form-item label="库存" prop="stock">
              <el-input-number
                v-model="formData.stock"
                :min="0"
                :step="1"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :xs="24" :md="12">
            <el-form-item label="状态" prop="status">
              <el-radio-group v-model="formData.status">
                <el-radio :value="1">上架</el-radio>
                <el-radio :value="0">下架</el-radio>
              </el-radio-group>
            </el-form-item>
          </el-col>

          <el-col :span="24">
            <el-form-item label="热门推荐" prop="isHot">
              <el-switch
                v-model="formData.isHot"
                active-text="是"
                inactive-text="否"
              />
            </el-form-item>
          </el-col>

          <el-col :span="24">
            <el-form-item label="描述" prop="description">
              <el-input
                v-model="formData.description"
                type="textarea"
                :rows="4"
                placeholder="请输入产品描述"
                maxlength="1000"
                show-word-limit
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-divider>产品图片</el-divider>
        <el-form-item label="图片">
          <el-upload
            list-type="picture-card"
            :file-list="fileList"
            :http-request="handleUpload"
            :on-success="handleUploadSuccess"
            :on-remove="handleRemoveImage"
            :before-upload="beforeUpload"
            :limit="8"
            multiple
          >
            <el-icon><Plus /></el-icon>
            <template #tip>
              <div class="upload-tip">
                支持 jpg / png，单张不超过 10MB，最多 8 张
              </div>
            </template>
          </el-upload>
        </el-form-item>

        <el-divider>自定义选项</el-divider>

        <div class="options-wrapper">
          <div
            v-for="(opt, oIdx) in formData.options"
            :key="oIdx"
            class="option-block"
          >
            <div class="option-header">
              <el-select v-model="opt.type" placeholder="选项类型" style="width: 140px">
                <el-option label="尺寸" value="size" />
                <el-option label="材质" value="material" />
                <el-option label="颜色" value="color" />
              </el-select>
              <el-input
                v-model="opt.name"
                placeholder="选项名称（如：沙发尺寸）"
                style="flex: 1; margin: 0 8px"
              />
              <el-checkbox v-model="opt.required">必选</el-checkbox>
              <el-button
                link
                type="danger"
                :icon="Delete"
                @click="removeOption(oIdx)"
              >
                删除
              </el-button>
            </div>
            <div class="option-values">
              <div
                v-for="(val, vIdx) in opt.values"
                :key="vIdx"
                class="option-value-row"
              >
                <el-input
                  v-model="val.value"
                  placeholder="选项值（如：1.8m / 米白）"
                  style="flex: 1"
                />
                <el-input-number
                  v-model="val.priceAdjustment"
                  :step="1"
                  :precision="2"
                  placeholder="价格调整"
                  style="width: 160px; margin-left: 8px"
                />
                <el-button
                  link
                  type="danger"
                  :icon="Delete"
                  @click="removeOptionValue(oIdx, vIdx)"
                />
              </div>
              <el-button
                link
                type="primary"
                :icon="Plus"
                @click="addOptionValue(oIdx)"
              >
                添加选项值
              </el-button>
            </div>
          </div>

          <el-button
            v-if="formData.options && formData.options.length < 3"
            type="primary"
            plain
            :icon="Plus"
            @click="addOption"
          >
            添加自定义选项
          </el-button>
        </div>

        <el-form-item style="margin-top: 24px; text-align: center">
          <el-button type="primary" :loading="submitting" @click="handleSubmit">
            {{ isEdit ? '保存' : '提交' }}
          </el-button>
          <el-button @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { onMounted, reactive, ref, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules, type UploadUserFile } from 'element-plus'
import { ArrowLeft, Plus, Delete } from '@element-plus/icons-vue'
import { createProduct, updateProduct, getProduct, uploadProductImage, listProductCategories } from '@/api/product'
import type { ProductCategory, ProductFormData, ProductOption, ProductOptionValue, ProductImage } from '@/types'

const route = useRoute()
const router = useRouter()

const formRef = ref<FormInstance>()
const loading = ref(false)
const submitting = ref(false)
const categories = ref<ProductCategory[]>([])
const fileList = ref<UploadUserFile[]>([])

const isEdit = computed(() => !!route.params.id)

const formData = reactive<ProductFormData>({
  name: '',
  sku: '',
  category: '',
  description: '',
  price: 0,
  stock: 0,
  isHot: false,
  status: 1,
  images: [],
  options: []
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入产品名称', trigger: 'blur' },
    { min: 2, max: 64, message: '长度在 2 到 64 个字符', trigger: 'blur' }
  ],
  category: [{ required: true, message: '请选择分类', trigger: 'change' }],
  price: [
    { required: true, message: '请输入价格', trigger: 'blur' },
    {
      validator: (_r, value, cb) => {
        if (value === undefined || value === null || value < 0) {
          cb(new Error('价格不能小于 0'))
        } else {
          cb()
        }
      },
      trigger: 'blur'
    }
  ],
  stock: [
    { required: true, message: '请输入库存', trigger: 'blur' }
  ]
}

async function loadCategories() {
  try {
    categories.value = await listProductCategories()
  } catch {
    categories.value = [
      { id: 1, name: '沙发' },
      { id: 2, name: '椅子' },
      { id: 3, name: '桌子' },
      { id: 4, name: '床' },
      { id: 5, name: '柜' }
    ]
  }
}

async function loadProduct() {
  if (!isEdit.value) return
  loading.value = true
  try {
    const data = await getProduct(route.params.id as string)
    Object.assign(formData, {
      name: data.name,
      sku: data.sku,
      category: data.category,
      description: data.description,
      price: data.price,
      stock: data.stock,
      isHot: data.isHot,
      status: data.status,
      images: data.images || [],
      options: data.options || []
    })
    fileList.value = (data.images || []).map((img) => ({
      name: `img-${img.id}`,
      url: img.url,
      uid: img.id as unknown as number
    }))
  } catch (err) {
    console.error(err)
  } finally {
    loading.value = false
  }
}

function addOption() {
  const opt: ProductOption = {
    type: 'size',
    name: '',
    required: false,
    values: [{ value: '', priceAdjustment: 0 }]
  }
  if (!formData.options) formData.options = []
  formData.options.push(opt)
}

function removeOption(idx: number) {
  formData.options!.splice(idx, 1)
}

function addOptionValue(idx: number) {
  const newVal: ProductOptionValue = { value: '', priceAdjustment: 0 }
  formData.options![idx].values.push(newVal)
}

function removeOptionValue(oIdx: number, vIdx: number) {
  formData.options![oIdx].values.splice(vIdx, 1)
}

function beforeUpload(file: File) {
  const isImage = file.type.startsWith('image/')
  if (!isImage) {
    ElMessage.error('仅支持上传图片')
    return false
  }
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error('图片大小不能超过 10MB')
    return false
  }
  return true
}

async function handleUpload(options: { file: File; onSuccess: (res: unknown) => void; onError: (err: unknown) => void }) {
  try {
    const fd = new FormData()
    fd.append('file', options.file)
    const data = await uploadProductImage(fd)
    options.onSuccess(data)
  } catch (err) {
    options.onError(err)
  }
}

function handleUploadSuccess(response: { url: string }, uploadFile: UploadUserFile) {
  uploadFile.url = response.url
  if (!formData.images) formData.images = []
  formData.images.push({ url: response.url, productId: 0, id: Date.now() } as ProductImage)
}

function handleRemoveImage(uploadFile: UploadUserFile) {
  formData.images = (formData.images || []).filter((img) => img.url !== uploadFile.url)
}

async function handleSubmit() {
  if (!formRef.value) return
  try {
    await formRef.value.validate()
  } catch {
    ElMessage.warning('请检查表单填写')
    return
  }
  submitting.value = true
  try {
    if (isEdit.value) {
      await updateProduct(route.params.id as string, formData)
      ElMessage.success('更新成功')
    } else {
      await createProduct(formData)
      ElMessage.success('创建成功')
    }
    router.push('/products')
  } catch (err) {
    console.error(err)
  } finally {
    submitting.value = false
  }
}

function goBack() {
  router.back()
}

onMounted(() => {
  loadCategories()
  loadProduct()
})
</script>

<style lang="scss" scoped>
.product-form {
  .form-card {
    margin-top: 16px;
  }

  .upload-tip {
    color: #909399;
    font-size: 12px;
    margin-top: 4px;
  }

  .options-wrapper {
    padding: 0 0 0 100px;

    .option-block {
      border: 1px solid #ebeef5;
      border-radius: 4px;
      padding: 12px 16px;
      margin-bottom: 12px;
      background-color: #fafbfc;
    }

    .option-header {
      display: flex;
      align-items: center;
      margin-bottom: 8px;
    }

    .option-values {
      .option-value-row {
        display: flex;
        margin-bottom: 8px;
      }
    }
  }
}
</style>
