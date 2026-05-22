<template>
  <div class="repair-edit">
    <div class="page-header">
      <h2 class="page-title">提交报修</h2>
      <el-button @click="goBack">返回</el-button>
    </div>

    <div class="card">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="选择租户" prop="tenantId">
              <el-select
                v-model="form.tenantId"
                placeholder="请选择租户"
                filterable
                style="width: 100%"
              >
                <el-option
                  v-for="t in tenants"
                  :key="t.id"
                  :label="`${t.name} (${t.phone})`"
                  :value="t.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="选择房源" prop="propertyId">
              <el-select
                v-model="form.propertyId"
                placeholder="请选择房源"
                filterable
                style="width: 100%"
              >
                <el-option
                  v-for="p in properties"
                  :key="p.id"
                  :label="p.title"
                  :value="p.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="报修标题" prop="title">
              <el-input v-model="form.title" placeholder="请输入报修标题" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="维修分类" prop="category">
              <el-select v-model="form.category" placeholder="请选择分类" style="width: 100%">
                <el-option label="水电维修" value="水电维修" />
                <el-option label="家电维修" value="家电维修" />
                <el-option label="家具维修" value="家具维修" />
                <el-option label="门窗维修" value="门窗维修" />
                <el-option label="其他" value="其他" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="优先级" prop="priority">
          <el-radio-group v-model="form.priority">
            <el-radio :value="3">普通</el-radio>
            <el-radio :value="2">高</el-radio>
            <el-radio :value="1">紧急</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="问题描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="4"
            placeholder="请详细描述问题"
          />
        </el-form-item>

        <el-form-item label="现场照片">
          <div class="image-upload">
            <div
              v-for="(url, index) in form.images"
              :key="index"
              class="image-item"
            >
              <img :src="url" alt="" />
              <button class="delete-btn" @click="removeImage(index)">×</button>
            </div>
            <el-upload
              :show-file-list="false"
              :before-upload="beforeUpload"
              :http-request="handleUpload"
              accept="image/*"
            >
              <div class="el-upload--picture-card">
                <el-icon><Plus /></el-icon>
              </div>
            </el-upload>
          </div>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">
            提交报修
          </el-button>
          <el-button @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules, type UploadProps } from 'element-plus'
import type { Tenant, Property } from '@/types'
import { getTenants } from '@/api/tenant'
import { getProperties } from '@/api/property'
import { createRepair } from '@/api/business'
import { uploadImage } from '@/api/property'

const router = useRouter()
const loading = ref(false)
const formRef = ref<FormInstance>()
const tenants = ref<Tenant[]>([])
const properties = ref<Property[]>([])

const form = reactive({
  tenantId: 0,
  propertyId: 0,
  title: '',
  category: '其他',
  priority: 3,
  description: '',
  images: [] as string[]
})

const rules: FormRules = {
  tenantId: [{ required: true, message: '请选择租户', trigger: 'change' }],
  propertyId: [{ required: true, message: '请选择房源', trigger: 'change' }],
  title: [{ required: true, message: '请输入报修标题', trigger: 'blur' }],
  description: [{ required: true, message: '请输入问题描述', trigger: 'blur' }]
}

onMounted(async () => {
  await loadTenants()
  await loadProperties()
})

async function loadTenants() {
  try {
    const res = await getTenants({ pageSize: 100 })
    tenants.value = res.data.list
  } catch (error) {
    console.error('Failed to load tenants:', error)
  }
}

async function loadProperties() {
  try {
    const res = await getProperties({ pageSize: 100, status: 2 })
    properties.value = res.data.list
  } catch (error) {
    console.error('Failed to load properties:', error)
  }
}

const beforeUpload: UploadProps['beforeUpload'] = (file) => {
  const isImage = file.type.startsWith('image/')
  if (!isImage) {
    ElMessage.error('只能上传图片文件')
    return false
  }
  return true
}

async function handleUpload(options: any) {
  try {
    const res = await uploadImage(options.file)
    form.images.push(res.data.url)
  } catch (error) {
    console.error('Upload failed:', error)
  }
}

function removeImage(index: number) {
  form.images.splice(index, 1)
}

async function handleSubmit() {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await createRepair({
          ...form,
          images: form.images.join(',')
        })
        ElMessage.success('报修提交成功')
        router.push('/repairs')
      } catch (error) {
        console.error('Submit failed:', error)
      } finally {
        loading.value = false
      }
    }
  })
}

function goBack() {
  router.back()
}
</script>

<style scoped>
.repair-edit {
  padding: 0;
}

.image-upload {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.image-upload .image-item {
  width: 100px;
  height: 100px;
  border-radius: 4px;
  overflow: hidden;
  position: relative;
  border: 1px solid #dcdfe6;
}

.image-upload .image-item img {
  width: 100%;
  height: 100%;
  object-fit: cover;
}

.image-upload .image-item .delete-btn {
  position: absolute;
  top: 4px;
  right: 4px;
  width: 20px;
  height: 20px;
  background: rgba(0, 0, 0, 0.5);
  color: #fff;
  border: none;
  border-radius: 50%;
  cursor: pointer;
  display: flex;
  align-items: center;
  justify-content: center;
  font-size: 12px;
}
</style>
