<template>
  <div class="rental-create">
    <div class="page-header">
      <h2 class="page-title">创建租赁订单</h2>
    </div>

    <div class="card-container">
      <el-form
        ref="formRef"
        :model="form"
        :rules="rules"
        label-width="120px"
        @submit.prevent="handleSubmit"
      >
        <el-form-item label="选择船只" prop="ship_id">
          <el-select
            v-model="form.ship_id"
            placeholder="请选择要租赁的船只"
            filterable
            style="width: 100%"
          >
            <el-option
              v-for="ship in ships"
              :key="ship.id"
              :label="`${ship.name} - ¥${ship.daily_rate}/天`"
              :value="ship.id"
            />
          </el-select>
        </el-form-item>

        <el-form-item label="租赁类型" prop="rental_type">
          <el-radio-group v-model="form.rental_type">
            <el-radio value="daily">按天租赁</el-radio>
            <el-radio value="hourly">按小时租赁</el-radio>
            <el-radio value="voyage">航程租赁</el-radio>
          </el-radio-group>
        </el-form-item>

        <el-form-item label="开始时间" prop="start_date">
          <el-date-picker
            v-model="form.start_date"
            type="datetime"
            placeholder="选择开始时间"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="结束时间" prop="end_date">
          <el-date-picker
            v-model="form.end_date"
            type="datetime"
            placeholder="选择结束时间"
            style="width: 100%"
          />
        </el-form-item>

        <el-form-item label="起始地点">
          <el-input v-model="form.start_location" placeholder="请输入起始地点" />
        </el-form-item>

        <el-form-item label="目的地">
          <el-input v-model="form.end_location" placeholder="请输入目的地" />
        </el-form-item>

        <el-form-item label="保险" prop="insurance_type">
          <el-select v-model="form.insurance_type" style="width: 100%">
            <el-option label="无保险" value="none" />
            <el-option label="基础保险 (5%)" value="basic" />
            <el-option label="高级保险 (10%)" value="premium" />
          </el-select>
        </el-form-item>

        <el-form-item label="紧急联系人" prop="emergency_contact">
          <el-input v-model="form.emergency_contact" placeholder="请输入紧急联系人姓名" />
        </el-form-item>

        <el-form-item label="紧急电话" prop="emergency_phone">
          <el-input v-model="form.emergency_phone" placeholder="请输入紧急联系电话" />
        </el-form-item>

        <el-form-item label="乘客人数">
          <el-input-number v-model="form.passenger_count" :min="0" />
        </el-form-item>

        <el-form-item label="船员人数">
          <el-input-number v-model="form.crew_count" :min="0" />
        </el-form-item>

        <el-form-item label="备注">
          <el-input
            v-model="form.notes"
            type="textarea"
            :rows="3"
            placeholder="其他备注信息"
          />
        </el-form-item>

        <el-form-item>
          <el-button type="primary" size="large" :loading="loading" @click="handleSubmit">
            创建订单
          </el-button>
          <el-button size="large" @click="$router.back()">取消</el-button>
        </el-form-item>
      </el-form>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { getShipsApi } from '@/api/ship'
import { createRentalApi } from '@/api/rental'
import type { Ship } from '@/types/ship'
import type { CreateRentalRequest } from '@/types/rental'

const router = useRouter()
const route = useRoute()
const formRef = ref<FormInstance>()
const loading = ref(false)
const ships = ref<Ship[]>([])

const form = reactive<CreateRentalRequest>({
  ship_id: route.query.ship_id as string || '',
  rental_type: (route.query.rental_type as any) || 'daily',
  start_date: route.query.start_date as string || '',
  end_date: route.query.end_date as string || '',
  start_location: '',
  end_location: '',
  insurance_type: (route.query.insurance_type as any) || 'none',
  emergency_contact: '',
  emergency_phone: '',
  notes: '',
  passenger_count: 0,
  crew_count: 0
})

const rules: FormRules = {
  ship_id: [
    { required: true, message: '请选择船只', trigger: 'change' }
  ],
  rental_type: [
    { required: true, message: '请选择租赁类型', trigger: 'change' }
  ],
  start_date: [
    { required: true, message: '请选择开始时间', trigger: 'change' }
  ],
  end_date: [
    { required: true, message: '请选择结束时间', trigger: 'change' }
  ],
  insurance_type: [
    { required: true, message: '请选择保险类型', trigger: 'change' }
  ],
  emergency_contact: [
    { required: true, message: '请输入紧急联系人', trigger: 'blur' }
  ],
  emergency_phone: [
    { required: true, message: '请输入紧急联系电话', trigger: 'blur' }
  ]
}

const handleSubmit = async () => {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        await createRentalApi(form)
        ElMessage.success('订单创建成功')
        router.push('/my-rentals')
      } catch (error: any) {
        ElMessage.error(error.message || '创建失败')
      } finally {
        loading.value = false
      }
    }
  })
}

onMounted(async () => {
  try {
    const res: any = await getShipsApi({ page_size: 100 })
    ships.value = res.data || []
  } catch (error) {
    console.error('Failed to fetch ships:', error)
  }
})
</script>

<style lang="scss" scoped>
.rental-create {
  .card-container {
    max-width: 800px;
    margin: 0 auto;
    background: #fff;
    border-radius: 8px;
    padding: 24px;
  }
}
</style>
