<template>
  <div class="ship-create">
    <div class="page-header">
      <h2 class="page-title">发布船只</h2>
    </div>

    <div class="card-container">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        @submit.prevent="handleSubmit"
      >
        <el-form-item label="船只名称" prop="name">
          <el-input v-model="form.name" placeholder="请输入船只名称" />
        </el-form-item>

        <el-form-item label="船型" prop="ship_type">
          <el-select v-model="form.ship_type" placeholder="请选择船型" style="width: 100%">
            <el-option label="帆船" value="sailboat" />
            <el-option label="摩托艇" value="motorboat" />
            <el-option label="游艇" value="yacht" />
            <el-option label="渔船" value="fishing" />
            <el-option label="货船" value="cargo" />
          </el-select>
        </el-form-item>

        <el-form-item label="描述" prop="description">
          <el-input
            v-model="form.description"
            type="textarea"
            :rows="4"
            placeholder="请描述您的船只"
          />
        </el-form-item>

        <el-row :gutter="16">
          <el-col :span="8">
            <el-form-item label="载客量" prop="capacity">
              <el-input-number v-model="form.capacity" :min="1" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="客舱数">
              <el-input-number v-model="form.cabin_count" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="卫生间">
              <el-input-number v-model="form.bathroom_count" :min="0" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-row :gutter="16">
          <el-col :span="8">
            <el-form-item label="船长(米)">
              <el-input-number v-model="form.length" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="船宽(米)">
              <el-input-number v-model="form.width" :min="0" :precision="2" style="width: 100%" />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="建造年份">
              <el-input-number v-model="form.year_built" :min="1900" :max="2100" style="width: 100%" />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="母港">
          <el-input v-model="form.home_port" placeholder="请输入母港名称" />
        </el-form-item>

        <el-form-item label="航行区域">
          <el-input v-model="form.sailing_area" placeholder="请输入航行区域" />
        </el-form-item>

        <el-form-item label="设备配置">
          <el-input
            v-model="form.equipment"
            type="textarea"
            :rows="3"
            placeholder="请列出设备配置"
          />
        </el-form-item>

        <el-form-item label="特色功能">
          <el-input
            v-model="form.features"
            type="textarea"
            :rows="3"
            placeholder="请描述特色功能"
          />
        </el-form-item>

        <el-row :gutter="16">
          <el-col :span="8">
            <el-form-item label="日租价格" prop="daily_rate">
              <el-input-number
                v-model="form.daily_rate"
                :min="0"
                :precision="2"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="小时价格" prop="hourly_rate">
              <el-input-number
                v-model="form.hourly_rate"
                :min="0"
                :precision="2"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
          <el-col :span="8">
            <el-form-item label="押金">
              <el-input-number
                v-model="form.deposit_amount"
                :min="0"
                :precision="2"
                style="width: 100%"
              />
            </el-form-item>
          </el-col>
        </el-row>

        <el-form-item label="需要保险">
          <el-switch v-model="form.insurance_required" />
        </el-form-item>

        <el-form-item label="取消政策">
          <el-input
            v-model="form.cancellation_policy"
            type="textarea"
            :rows="3"
            placeholder="请描述取消政策"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" size="large" :loading="loading" @click="handleSubmit">
            发布船只
          </el-button>
          <el-button size="large" @click="$router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { createShipApi } from '@/api/ship'
import type { CreateShipRequest } from '@/types/ship'

const router = useRouter()
const formRef = ref<FormInstance>()
const loading = ref(false)

const form = reactive<CreateShipRequest>({
  name: '',
  ship_type: 'sailboat',
  description: '',
  capacity: 4,
  cabin_count: 0,
  bathroom_count: 0,
  length: 0,
  width: 0,
  year_built: new Date().getFullYear(),
  home_port: '',
  sailing_area: '',
  equipment: '',
  features: '',
  daily_rate: 0,
  hourly_rate: 0,
  deposit_amount: 0,
  insurance_required: true,
  cancellation_policy: ''
})

const rules: FormRules = {
  name: [
    { required: true, message: '请输入船只名称', trigger: 'blur' }
  ],
  ship_type: [
    { required: true, message: '请选择船型', trigger: 'change' }
  ],
  capacity: [
    { required: true, message: '请输入载客量', trigger: 'blur' }
  ],
  daily_rate: [
    { required: true, message: '请输入日租价格', trigger: 'blur' }
  ],
  hourly_rate: [
    { required: true, message: '请输入小时价格', trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await createShipApi(form)
        ElMessage.success('船只发布成功')
        router.push('/my-ships')
      } catch (error: any) {
        ElMessage.error(error.message || '发布失败')
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style lang="scss" scoped>
.ship-create {
  .card-container {
    max-width: 800px;
    margin: 0 auto;
    background: #fff;
    border-radius: 8px;
    padding: 24px;
  }
}
</style>
