<template>
  <el-card>
    <template #header>
      <span>{{ isEdit ? '编辑设备' : '添加设备' }}</span>
    </template>
    <el-form :model="form" :rules="rules" ref="formRef" label-width="120px" style="max-width: 700px">
      <el-form-item label="设备名称" prop="name">
        <el-input v-model="form.name" />
      </el-form-item>
      <el-form-item label="品牌" prop="brand">
        <el-input v-model="form.brand" />
      </el-form-item>
      <el-form-item label="型号" prop="model">
        <el-input v-model="form.model" />
      </el-form-item>
      <el-form-item label="序列号" prop="serial_no">
        <el-input v-model="form.serial_no" />
      </el-form-item>
      <el-form-item label="重量(kg)">
        <el-input-number v-model="form.weight" :min="0" :precision="2" />
      </el-form-item>
      <el-form-item label="续航(分钟)">
        <el-input-number v-model="form.battery_life" :min="0" />
      </el-form-item>
      <el-form-item label="云台规格">
        <el-input v-model="form.gimbal_spec" />
      </el-form-item>
      <el-form-item label="相机规格">
        <el-input v-model="form.camera_spec" />
      </el-form-item>
      <el-form-item label="最大速度(m/s)">
        <el-input-number v-model="form.max_speed" :min="0" :precision="1" />
      </el-form-item>
      <el-form-item label="最大高度(m)">
        <el-input-number v-model="form.max_altitude" :min="0" :precision="1" />
      </el-form-item>
      <el-form-item label="区域" prop="region">
        <el-input v-model="form.region" />
      </el-form-item>
      <el-form-item label="可用时段" prop="available_from">
        <el-date-picker
          v-model="form.available_range"
          type="daterange"
          range-separator="至"
          start-placeholder="开始日期"
          end-placeholder="结束日期"
          value-format="YYYY-MM-DD"
          style="width: 100%"
        />
      </el-form-item>
      <el-form-item label="日租金" prop="price_per_day">
        <el-input-number v-model="form.price_per_day" :min="0" :precision="2" />
      </el-form-item>
      <el-form-item label="押金">
        <el-input-number v-model="form.deposit" :min="0" :precision="2" />
      </el-form-item>
      <el-form-item label="描述">
        <el-input v-model="form.description" type="textarea" :rows="3" />
      </el-form-item>
      <el-form-item label="图片">
        <el-upload
          :action="'/api/upload/image'"
          :headers="{ Authorization: 'Bearer ' + token }"
          :show-file-list="false"
          :on-success="handleUpload"
          list-type="picture-card"
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
import { ref, reactive, computed, onMounted, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { useUserStore } from '@/stores/user'
import request from '@/utils/request'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const token = computed(() => userStore.token)

const formRef = ref<FormInstance>()
const loading = ref(false)
const isEdit = computed(() => !!route.params.id)

const form = reactive({
  name: '',
  brand: '',
  model: '',
  serial_no: '',
  weight: 0,
  battery_life: 0,
  gimbal_spec: '',
  camera_spec: '',
  max_speed: 0,
  max_altitude: 0,
  region: '',
  price_per_day: 0,
  deposit: 0,
  description: '',
  images: '',
  available_from: '',
  available_to: '',
  available_range: [] as string[]
})

const rules: FormRules = {
  name: [{ required: true, message: '请输入设备名称', trigger: 'blur' }],
  serial_no: [{ required: true, message: '请输入序列号', trigger: 'blur' }],
  region: [{ required: true, message: '请输入区域', trigger: 'blur' }],
  price_per_day: [{ required: true, message: '请输入日租金', trigger: 'blur' }]
}

watch(() => form.available_range, (val) => {
  if (val && val.length === 2) {
    form.available_from = val[0]
    form.available_to = val[1]
  } else {
    form.available_from = ''
    form.available_to = ''
  }
})

onMounted(() => {
  if (isEdit.value) {
    fetchDrone()
  }
})

async function fetchDrone() {
  try {
    const res: any = await request.get(`/drones/${route.params.id}`)
    Object.assign(form, res.data)
    if (res.data.available_from || res.data.available_to) {
      form.available_range = [res.data.available_from || '', res.data.available_to || '']
    }
  } catch (e) {
    console.error(e)
  }
}

function handleUpload(res: any) {
  form.images = res.data.url
}

async function handleSubmit() {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const submitData = {
          ...form,
          available_from: form.available_from,
          available_to: form.available_to
        }
        delete (submitData as any).available_range
        
        if (isEdit.value) {
          await request.put(`/drones/${route.params.id}`, submitData)
        } else {
          await request.post('/drones', submitData)
        }
        ElMessage.success(isEdit.value ? '修改成功' : '添加成功')
        router.push('/my-drones')
      } catch (e: any) {
        ElMessage.error(e.message || (isEdit.value ? '修改失败' : '添加失败'))
      } finally {
        loading.value = false
      }
    }
  })
}
</script>
