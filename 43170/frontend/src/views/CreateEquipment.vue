<template>
  <div class="page-container">
    <div class="page-header">
      <el-button @click="goBack">
        <el-icon><ArrowLeft /></el-icon>
        返回
      </el-button>
      <h1 class="page-title">添加设备</h1>
    </div>

    <div class="card">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="100px"
      >
        <el-form-item label="设备名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入设备名称" />
        </el-form-item>
        <el-form-item label="分类" prop="category">
          <el-select v-model="form.category" placeholder="请选择分类" style="width: 100%">
            <el-option label="相机" value="相机" />
            <el-option label="镜头" value="镜头" />
            <el-option label="灯光" value="灯光" />
            <el-option label="稳定器" value="稳定器" />
            <el-option label="无人机" value="无人机" />
            <el-option label="其他" value="其他" />
          </el-select>
        </el-form-item>
        <el-form-item label="品牌" prop="brand">
          <el-input v-model="form.brand" placeholder="请输入品牌" />
        </el-form-item>
        <el-form-item label="型号" prop="model">
          <el-input v-model="form.model" placeholder="请输入型号" />
        </el-form-item>
        <el-form-item label="购买时间" prop="purchaseDate">
          <el-date-picker
            v-model="form.purchaseDate"
            type="date"
            placeholder="选择日期"
            style="width: 100%"
            value-format="YYYY-MM-DD"
          />
        </el-form-item>
        <el-form-item label="日租金" prop="dailyRent">
          <el-input-number
            v-model="form.dailyRent"
            :min="0"
            :precision="2"
            :step="10"
            style="width: 100%"
          />
          <span style="margin-left: 10px; color: #909399">元/天</span>
        </el-form-item>
        <el-form-item label="押金" prop="deposit">
          <el-input-number
            v-model="form.deposit"
            :min="0"
            :precision="2"
            :step="100"
            style="width: 100%"
          />
          <span style="margin-left: 10px; color: #909399">元</span>
        </el-form-item>
        <el-form-item label="设备描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="4"
            placeholder="请输入设备描述"
          />
        </el-form-item>
        <el-form-item label="设备图片">
          <el-upload
            v-model:file-list="fileList"
            :action="uploadAction"
            :headers="uploadHeaders"
            :auto-upload="false"
            :limit="9"
            list-type="picture-card"
            :on-change="handleFileChange"
            :on-exceed="handleExceed"
            :on-preview="handlePreview"
            :on-remove="handleRemove"
            accept="image/jpeg,image/png"
          >
            <el-icon><Plus /></el-icon>
            <template #tip>
              <div class="el-upload__tip">
                最多上传9张图片，支持JPG/PNG格式
              </div>
            </template>
          </el-upload>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :loading="submitting" @click="handleSubmit">
            提交
          </el-button>
          <el-button @click="goBack">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules, type UploadUserFile } from 'element-plus'
import { equipmentApi } from '@/api/equipment'
import { useUserStore } from '@/stores/user'

const router = useRouter()
const userStore = useUserStore()
const formRef = ref<FormInstance>()
const submitting = ref(false)
const fileList = ref<UploadUserFile[]>([])
const newEquipmentId = ref<number | null>(null)

const uploadAction = ref('')
const uploadHeaders = computed(() => ({
  Authorization: `Bearer ${userStore.token}`
}))

const form = reactive({
  name: '',
  category: '',
  brand: '',
  model: '',
  purchaseDate: '',
  dailyRent: 0,
  deposit: 0,
  description: ''
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入设备名称', trigger: 'blur' },
    { max: 100, message: '名称长度不能超过100个字符', trigger: 'blur' }
  ],
  category: [
    { required: true, message: '请选择分类', trigger: 'change' }
  ],
  brand: [
    { required: true, message: '请输入品牌', trigger: 'blur' },
    { max: 50, message: '品牌长度不能超过50个字符', trigger: 'blur' }
  ],
  model: [
    { required: true, message: '请输入型号', trigger: 'blur' },
    { max: 100, message: '型号长度不能超过100个字符', trigger: 'blur' }
  ],
  dailyRent: [
    { required: true, message: '请输入日租金', trigger: 'blur' }
  ],
  deposit: [
    { required: true, message: '请输入押金', trigger: 'blur' }
  ]
}

function handleFileChange(file: UploadUserFile) {
  if (newEquipmentId.value && file.raw) {
    uploadImage(file.raw)
  }
}

async function uploadImage(file: File) {
  try {
    if (newEquipmentId.value) {
      await equipmentApi.uploadImage(newEquipmentId.value, file)
    }
  } catch (error) {
    console.error('Image upload failed:', error)
    ElMessage.error('图片上传失败')
  }
}

function handleExceed() {
  ElMessage.warning('最多只能上传9张图片')
}

function handlePreview(_file: UploadUserFile) {
  // Preview handler
}

function handleRemove(file: UploadUserFile) {
  const responseData = file.response as { data?: { id?: number } } | undefined
  if (responseData?.data?.id && newEquipmentId.value) {
    equipmentApi.deleteImage(newEquipmentId.value, responseData.data.id)
  }
}

async function handleSubmit() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        const response = await equipmentApi.createEquipment(form)
        newEquipmentId.value = response.data.id

        const uploadPromises = fileList.value
          .filter(f => f.raw)
          .map(f => equipmentApi.uploadImage(newEquipmentId.value!, f.raw!))

        await Promise.all(uploadPromises)

        ElMessage.success('设备添加成功')
        router.push('/my-equipments')
      } catch (error) {
        console.error('Failed to create equipment:', error)
      } finally {
        submitting.value = false
      }
    }
  })
}

function goBack() {
  router.back()
}
</script>

<style scoped>
.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.el-upload__tip {
  font-size: 12px;
  color: #909399;
  margin-top: 8px;
}
</style>
