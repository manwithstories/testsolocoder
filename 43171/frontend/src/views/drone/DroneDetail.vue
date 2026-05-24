<template>
  <el-card v-if="drone">
    <el-row :gutter="20">
      <el-col :span="10">
        <el-image :src="drone.images?.split(',')[0]" fit="cover" style="width: 100%; height: 300px; border-radius: 8px">
          <template #error>
            <div style="display: flex; align-items: center; justify-content: center; height: 300px; background: #f0f2f5">
              <el-icon :size="64"><Box /></el-icon>
            </div>
          </template>
        </el-image>
      </el-col>
      <el-col :span="14">
        <h2>{{ drone.name }}</h2>
        <p class="brand">{{ drone.brand }} · {{ drone.model }}</p>
        <div class="price">¥{{ drone.price_per_day }}<span class="unit">/天</span></div>
        <div class="rating">
          <el-rate :model-value="drone.rating" disabled />
          <span class="rating-text">{{ drone.rating.toFixed(1) }} ({{ drone.rating_count }}条评价)</span>
        </div>
        <el-descriptions :column="2" border style="margin-top: 20px">
          <el-descriptions-item label="区域">{{ drone.region }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="statusTagType(drone.status)">{{ statusText(drone.status) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="续航">{{ drone.battery_life }}分钟</el-descriptions-item>
          <el-descriptions-item label="最大速度">{{ drone.max_speed }}m/s</el-descriptions-item>
          <el-descriptions-item label="最大高度">{{ drone.max_altitude }}m</el-descriptions-item>
          <el-descriptions-item label="重量">{{ drone.weight }}kg</el-descriptions-item>
          <el-descriptions-item label="云台">{{ drone.gimbal_spec }}</el-descriptions-item>
          <el-descriptions-item label="相机">{{ drone.camera_spec }}</el-descriptions-item>
        </el-descriptions>
        <p class="desc" v-if="drone.description">{{ drone.description }}</p>
        <div style="margin-top: 20px">
          <el-button type="primary" size="large" @click="showOrder = true" :disabled="drone.status !== 'online'">
            立即租赁
          </el-button>
          <el-button size="large" @click="$router.back()">返回</el-button>
        </div>
      </el-col>
    </el-row>

    <el-divider />

    <h3>用户评价</h3>
    <el-table :data="reviews">
      <el-table-column prop="content" label="评价内容" />
      <el-table-column label="评分" width="150">
        <template #default="{ row }">
          <el-rate :model-value="row.rating" disabled size="small" />
        </template>
      </el-table-column>
      <el-table-column label="评价人" width="150">
        <template #default="{ row }">
          {{ row.reviewer?.nickname || row.reviewer?.username }}
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="时间" width="180" />
    </el-table>
  </el-card>

  <el-dialog v-model="showOrder" title="租赁订单" width="500px">
    <el-form :model="orderForm" :rules="orderRules" ref="orderFormRef" label-width="100px">
      <el-form-item label="开始日期" prop="start_date">
        <el-date-picker v-model="orderForm.start_date" type="date" placeholder="选择开始日期" />
      </el-form-item>
      <el-form-item label="结束日期" prop="end_date">
        <el-date-picker v-model="orderForm.end_date" type="date" placeholder="选择结束日期" />
      </el-form-item>
      <el-form-item label="区域" prop="region">
        <el-input v-model="orderForm.region" />
      </el-form-item>
      <el-form-item label="联系地址">
        <el-input v-model="orderForm.address" />
      </el-form-item>
      <el-form-item label="联系人" prop="contact_name">
        <el-input v-model="orderForm.contact_name" />
      </el-form-item>
      <el-form-item label="联系电话" prop="contact_phone">
        <el-input v-model="orderForm.contact_phone" />
      </el-form-item>
      <el-form-item label="备注">
        <el-input v-model="orderForm.remark" type="textarea" />
      </el-form-item>
    </el-form>
    <template #footer>
      <el-button @click="showOrder = false">取消</el-button>
      <el-button type="primary" @click="submitOrder" :loading="submitting">确认下单</el-button>
    </template>
  </el-dialog>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import request from '@/utils/request'
import dayjs from 'dayjs'

const route = useRoute()
const droneId = route.params.id as string

const drone = ref<Drone | null>(null)
const reviews = ref<Review[]>([])

const showOrder = ref(false)
const submitting = ref(false)
const orderFormRef = ref<FormInstance>()

const orderForm = reactive({
  start_date: '',
  end_date: '',
  region: '',
  address: '',
  contact_name: '',
  contact_phone: '',
  remark: ''
})

const orderRules: FormRules = {
  start_date: [{ required: true, message: '请选择开始日期', trigger: 'change' }],
  end_date: [{ required: true, message: '请选择结束日期', trigger: 'change' }],
  region: [{ required: true, message: '请输入区域', trigger: 'blur' }],
  contact_name: [{ required: true, message: '请输入联系人', trigger: 'blur' }],
  contact_phone: [{ required: true, message: '请输入联系电话', trigger: 'blur' }]
}

onMounted(() => {
  fetchDrone()
  fetchReviews()
})

async function fetchDrone() {
  try {
    const res: any = await request.get(`/drones/${droneId}`)
    drone.value = res.data
  } catch (e) {
    console.error(e)
  }
}

async function fetchReviews() {
  try {
    const res: any = await request.get('/reviews', { params: { drone_id: droneId, type: 'rental' } })
    reviews.value = res.data.list || []
  } catch (e) {
    console.error(e)
  }
}

async function submitOrder() {
  if (!orderFormRef.value) return
  await orderFormRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        await request.post('/orders', {
          drone_id: drone.value?.id,
          start_date: dayjs(orderForm.start_date).format('YYYY-MM-DD'),
          end_date: dayjs(orderForm.end_date).format('YYYY-MM-DD'),
          region: orderForm.region,
          address: orderForm.address,
          contact_name: orderForm.contact_name,
          contact_phone: orderForm.contact_phone,
          remark: orderForm.remark
        })
        ElMessage.success('下单成功，请前往支付')
        showOrder.value = false
      } catch (e: any) {
        ElMessage.error(e.message || '下单失败')
      } finally {
        submitting.value = false
      }
    }
  })
}

function statusText(status: string) {
  const map: Record<string, string> = {
    offline: '已下架', online: '可租', rented: '已租', maintenance: '维护中'
  }
  return map[status] || status
}

function statusTagType(status: string) {
  const map: Record<string, string> = {
    offline: 'info', online: 'success', rented: 'warning', maintenance: 'danger'
  }
  return map[status] || ''
}
</script>

<style scoped>
.brand {
  color: #909399;
  margin: 8px 0;
}
.price {
  font-size: 32px;
  color: #f56c6c;
  font-weight: bold;
}
.unit {
  font-size: 16px;
  color: #909399;
  font-weight: normal;
}
.rating {
  display: flex;
  align-items: center;
  gap: 10px;
  margin: 10px 0;
}
.rating-text {
  color: #909399;
}
.desc {
  color: #606266;
  margin-top: 16px;
  line-height: 1.6;
}
</style>
