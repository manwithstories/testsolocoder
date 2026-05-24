<template>
  <div class="create-order-page">
    <div class="container">
      <div class="page-header">
        <h2 class="page-title">创建工单</h2>
        <el-button @click="router.back()">返回</el-button>
      </div>

      <el-card>
        <el-form :model="form" :rules="rules" ref="formRef" label-width="100px">
          <el-form-item label="服务分类" prop="category_id">
            <el-select v-model="form.category_id" placeholder="请选择服务分类" @change="handleCategoryChange">
              <el-option
                v-for="category in categories"
                :key="category.id"
                :label="category.name"
                :value="category.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item label="服务项目" prop="service_item_id">
            <el-select v-model="form.service_item_id" placeholder="请选择服务项目" @change="handleServiceItemChange">
              <el-option
                v-for="item in filteredServiceItems"
                :key="item.id"
                :label="`${item.name} (¥${item.min_price}-¥${item.max_price})`"
                :value="item.id"
              />
            </el-select>
          </el-form-item>

          <el-form-item v-if="selectedServiceItem" label="预计费用">
            <span>¥{{ selectedServiceItem.min_price }} - ¥{{ selectedServiceItem.max_price }}</span>
            <span style="margin-left: 20px; color: #909399;">
              预计时长：{{ selectedServiceItem.estimated_time }}分钟
            </span>
          </el-form-item>

          <el-form-item label="问题标题" prop="title">
            <el-input v-model="form.title" placeholder="请简要描述维修问题" />
          </el-form-item>

          <el-form-item label="详细描述" prop="description">
            <el-input
              v-model="form.description"
              type="textarea"
              :rows="4"
              placeholder="请详细描述您遇到的问题"
            />
          </el-form-item>

          <el-form-item label="服务地址" prop="address">
            <el-input v-model="form.address" placeholder="请输入详细地址" />
          </el-form-item>

          <el-form-item label="联系人" prop="contact_name">
            <el-input v-model="form.contact_name" placeholder="请输入联系人姓名" />
          </el-form-item>

          <el-form-item label="联系电话" prop="contact_phone">
            <el-input v-model="form.contact_phone" placeholder="请输入联系电话" />
          </el-form-item>

          <el-form-item label="预约时间">
            <el-date-picker
              v-model="form.appointment_time"
              type="datetime"
              placeholder="选择预约时间"
              format="YYYY-MM-DD HH:mm"
            />
          </el-form-item>

          <el-form-item label="紧急程度">
            <el-radio-group v-model="form.urgent_level">
              <el-radio :value="0">普通</el-radio>
              <el-radio :value="1">加急</el-radio>
              <el-radio :value="2">特急</el-radio>
            </el-radio-group>
          </el-form-item>

          <el-form-item>
            <el-button type="primary" @click="submitOrder" :loading="loading">
              提交工单
            </el-button>
            <el-button @click="router.back()">取消</el-button>
          </el-form-item>
        </el-form>
      </el-card>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter, useRoute } from 'vue-router'
import { ElMessage, type FormInstance, type FormRules } from 'element-plus'
import { categoryApi, serviceItemApi } from '@/api/category'
import { orderApi } from '@/api/order'
import type { Category, ServiceItem } from '@/types'

const router = useRouter()
const route = useRoute()

const formRef = ref<FormInstance>()
const loading = ref(false)
const categories = ref<Category[]>([])
const serviceItems = ref<ServiceItem[]>([])

const form = reactive({
  category_id: null as number | null,
  service_item_id: null as number | null,
  title: '',
  description: '',
  address: '',
  contact_name: '',
  contact_phone: '',
  appointment_time: null as string | null,
  urgent_level: 0
})

const rules: FormRules = {
  category_id: [{ required: true, message: '请选择服务分类', trigger: 'change' }],
  service_item_id: [{ required: true, message: '请选择服务项目', trigger: 'change' }],
  title: [{ required: true, message: '请输入问题标题', trigger: 'blur' }],
  address: [{ required: true, message: '请输入服务地址', trigger: 'blur' }],
  contact_name: [{ required: true, message: '请输入联系人', trigger: 'blur' }],
  contact_phone: [{ required: true, message: '请输入联系电话', trigger: 'blur' }]
}

const filteredServiceItems = computed(() => {
  if (!form.category_id) return serviceItems.value
  return serviceItems.value.filter(item => item.category_id === form.category_id)
})

const selectedServiceItem = computed(() => {
  if (!form.service_item_id) return null
  return serviceItems.value.find(item => item.id === form.service_item_id)
})

onMounted(async () => {
  await loadCategories()
  await loadServiceItems()

  const categoryId = route.query.category_id
  const serviceItemId = route.query.service_item_id
  if (categoryId) {
    form.category_id = Number(categoryId)
  }
  if (serviceItemId) {
    form.service_item_id = Number(serviceItemId)
    const item = serviceItems.value.find(i => i.id === form.service_item_id)
    if (item) {
      form.category_id = item.category_id
    }
  }
})

async function loadCategories() {
  try {
    const res = await categoryApi.getCategories()
    categories.value = res.data || []
  } catch (error) {
    console.error('Failed to load categories:', error)
  }
}

async function loadServiceItems() {
  try {
    const res = await serviceItemApi.getServiceItems()
    serviceItems.value = res.data || []
  } catch (error) {
    console.error('Failed to load service items:', error)
  }
}

function handleCategoryChange() {
  form.service_item_id = null
}

function handleServiceItemChange() {
  const item = selectedServiceItem.value
  if (item && !form.title) {
    form.title = item.name
  }
}

async function submitOrder() {
  if (!formRef.value) return

  await formRef.value.validate(async (valid) => {
    if (valid) {
      loading.value = true
      try {
        const res = await orderApi.createOrder({
          service_item_id: form.service_item_id!,
          title: form.title,
          description: form.description,
          address: form.address,
          contact_name: form.contact_name,
          contact_phone: form.contact_phone,
          appointment_time: form.appointment_time || undefined,
          urgent_level: form.urgent_level
        })
        ElMessage.success('工单创建成功')
        if (res.data?.order_id) {
          router.push(`/orders/${res.data.order_id}`)
        } else {
          router.push('/orders')
        }
      } catch (error) {
        console.error('Failed to create order:', error)
      } finally {
        loading.value = false
      }
    }
  })
}
</script>

<style scoped>
.create-order-page {
  min-height: 100vh;
  background-color: #f5f7fa;
}
</style>
