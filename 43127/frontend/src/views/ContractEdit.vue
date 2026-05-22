<template>
  <div class="contract-edit">
    <div class="page-header">
      <h2 class="page-title">创建合同</h2>
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
                @change="handlePropertyChange"
              >
                <el-option
                  v-for="p in properties"
                  :key="p.id"
                  :label="`${p.title} - ¥${p.rent}/月`"
                  :value="p.id"
                />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="开始日期" prop="startDate">
              <el-date-picker
                v-model="form.startDate"
                type="date"
                placeholder="选择开始日期"
                value-format="YYYY-MM-DD"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="结束日期" prop="endDate">
              <el-date-picker
                v-model="form.endDate"
                type="date"
                placeholder="选择结束日期"
                value-format="YYYY-MM-DD"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="20">
          <el-col :span="8">
            <el-form-item label="月租金(元)" prop="rent">
              <el-input-number v-model="form.rent" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="押金(元)" prop="deposit">
              <el-input-number v-model="form.deposit" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="付款方式" prop="paymentType">
              <el-select v-model="form.paymentType" placeholder="请选择" style="width: 100%">
                <el-option label="押一付一" value="押一付一" />
                <el-option label="押一付三" value="押一付三" />
                <el-option label="押一付六" value="押一付六" />
                <el-option label="年付" value="年付" />
              </el-select>
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="合同文件">
          <el-upload
            :show-file-list="false"
            :before-upload="beforeUpload"
            :http-request="handleUpload"
            accept=".pdf,.doc,.docx"
          >
            <el-button>点击上传合同文件</el-button>
          </el-upload>
          <span v-if="form.fileUrl" style="margin-left: 10px; color: #67c23a;">
            已上传: {{ form.fileUrl }}
          </span>
        </el-form-item>

        <el-form-item>
          <el-button type="primary" :loading="loading" @click="handleSubmit">
            创建合同
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
import { getMyProperties } from '@/api/property'
import { createContract } from '@/api/tenant'
import { uploadImage } from '@/api/property'

const router = useRouter()
const loading = ref(false)
const formRef = ref<FormInstance>()
const tenants = ref<Tenant[]>([])
const properties = ref<Property[]>([])

const form = reactive({
  tenantId: 0,
  propertyId: 0,
  startDate: '',
  endDate: '',
  rent: 0,
  deposit: 0,
  paymentType: '押一付三',
  fileUrl: ''
})

const rules: FormRules = {
  tenantId: [{ required: true, message: '请选择租户', trigger: 'change' }],
  propertyId: [{ required: true, message: '请选择房源', trigger: 'change' }],
  startDate: [{ required: true, message: '请选择开始日期', trigger: 'change' }],
  endDate: [{ required: true, message: '请选择结束日期', trigger: 'change' }],
  rent: [{ required: true, message: '请输入月租金', trigger: 'blur' }]
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
    const res = await getMyProperties({ pageSize: 100, status: 1 })
    properties.value = res.data.list
  } catch (error) {
    console.error('Failed to load properties:', error)
  }
}

function handlePropertyChange(id: number) {
  const property = properties.value.find(p => p.id === id)
  if (property) {
    form.rent = property.rent
    form.deposit = property.deposit
  }
}

const beforeUpload: UploadProps['beforeUpload'] = (file) => {
  const isLt10M = file.size / 1024 / 1024 < 10
  if (!isLt10M) {
    ElMessage.error('文件大小不能超过10MB')
    return false
  }
  return true
}

async function handleUpload(options: any) {
  try {
    const res = await uploadImage(options.file)
    form.fileUrl = res.data.url
    ElMessage.success('上传成功')
  } catch (error) {
    console.error('Upload failed:', error)
  }
}

async function handleSubmit() {
  if (!formRef.value) return
  
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await createContract(form)
        ElMessage.success('合同创建成功')
        router.push('/contracts')
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
.contract-edit {
  padding: 0;
}
</style>
