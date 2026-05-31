<template>
  <div class="service-detail-page">
    <div class="container">
      <el-breadcrumb class="breadcrumb" separator="/">
        <el-breadcrumb-item :to="{ path: '/' }">首页</el-breadcrumb-item>
        <el-breadcrumb-item :to="{ path: '/services' }">维修服务</el-breadcrumb-item>
        <el-breadcrumb-item>{{ service?.title }}</el-breadcrumb-item>
      </el-breadcrumb>

      <div class="service-content card" v-loading="loading">
        <div class="service-header">
          <div class="service-icon">
            <el-icon :size="64"><Tools /></el-icon>
          </div>
          <div class="service-info">
            <h1 class="service-title">{{ service?.title }}</h1>
            <div class="service-tags">
              <el-tag type="primary">{{ service?.serviceType }}</el-tag>
              <el-tag>预计 {{ service?.estimatedDays }} 天完成</el-tag>
              <el-tag type="success" v-if="service?.status === 1">可接单</el-tag>
            </div>
            <div class="service-price">
              <span class="current-price">¥{{ service?.price.toFixed(2) }}</span>
              <span class="price-range" v-if="service?.minPrice">
                价格范围 ¥{{ service?.minPrice }}-¥{{ service?.maxPrice }}
              </span>
            </div>
          </div>
        </div>

        <div class="service-description">
          <h3>服务描述</h3>
          <p>{{ service?.description }}</p>
        </div>

        <div class="technician-info card">
          <h3>维修技师</h3>
          <div class="technician-detail">
            <el-avatar :size="60" :src="service?.technician?.avatar">
              {{ service?.technician?.nickname?.charAt(0) }}
            </el-avatar>
            <div class="technician-info-detail">
              <div class="technician-name">
                {{ service?.technician?.nickname || service?.technician?.username }}
              </div>
              <div class="technician-stats">
                <span>评分：<el-rate :model-value="service?.rating || 0" disabled size="small" /></span>
                <span>已接 {{ service?.orderCount }} 单</span>
                <span>信用分：{{ service?.technician?.creditScore }}</span>
              </div>
            </div>
          </div>
        </div>

        <div class="action-section">
          <el-button
            type="primary"
            size="large"
            :disabled="!canOrder"
            @click="showOrderDialog = true"
          >
            立即预约维修
          </el-button>
        </div>
      </div>
    </div>

    <el-dialog v-model="showOrderDialog" title="预约维修" width="500px">
      <el-form :model="orderForm" label-width="100px">
        <el-form-item label="服务项目">
          <span>{{ service?.title }}</span>
        </el-form-item>
        <el-form-item label="服务价格">
          <span class="price-text">¥{{ service?.price.toFixed(2) }}</span>
        </el-form-item>
        <el-form-item label="设备类型" prop="deviceType">
          <el-select v-model="orderForm.deviceType" placeholder="请选择设备类型">
            <el-option label="手机" value="手机" />
            <el-option label="电脑" value="电脑" />
            <el-option label="相机" value="相机" />
            <el-option label="耳机" value="耳机" />
            <el-option label="平板" value="平板" />
            <el-option label="其他" value="其他" />
          </el-select>
        </el-form-item>
        <el-form-item label="品牌" prop="deviceBrand">
          <el-input v-model="orderForm.deviceBrand" placeholder="请输入品牌" />
        </el-form-item>
        <el-form-item label="型号" prop="deviceModel">
          <el-input v-model="orderForm.deviceModel" placeholder="请输入型号" />
        </el-form-item>
        <el-form-item label="故障描述" prop="faultDescription">
          <el-input
            v-model="orderForm.faultDescription"
            type="textarea"
            :rows="3"
            placeholder="请详细描述故障情况"
          />
        </el-form-item>
        <el-form-item label="联系人" prop="contactName">
          <el-input v-model="orderForm.contactName" placeholder="请输入联系人姓名" />
        </el-form-item>
        <el-form-item label="联系电话" prop="contactPhone">
          <el-input v-model="orderForm.contactPhone" placeholder="请输入联系电话" />
        </el-form-item>
        <el-form-item label="地址">
          <el-input
            v-model="orderForm.address"
            type="textarea"
            :rows="2"
            placeholder="请输入地址（选填）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showOrderDialog = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="submitOrder">
          提交预约
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, reactive, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { useUserStore } from '@/store/user'
import { repairApi } from '@/api'
import type { RepairService } from '@/types'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()

const loading = ref(false)
const service = ref<RepairService | null>(null)
const showOrderDialog = ref(false)
const submitting = ref(false)

const orderForm = reactive({
  deviceType: '',
  deviceBrand: '',
  deviceModel: '',
  faultDescription: '',
  contactName: '',
  contactPhone: '',
  address: ''
})

const canOrder = computed(() => {
  return userStore.isLoggedIn &&
    service.value?.status === 1 &&
    service.value?.technicianId !== userStore.userInfo?.id
})

async function fetchService() {
  const id = parseInt(route.params.id as string)
  if (!id) return

  loading.value = true
  try {
    const res = await repairApi.getServiceById(id)
    service.value = res.data
  } catch (error) {
    console.error('Failed to fetch service:', error)
    ElMessage.error('获取服务详情失败')
  } finally {
    loading.value = false
  }
}

async function submitOrder() {
  if (!orderForm.deviceType || !orderForm.deviceBrand || !orderForm.deviceModel ||
      !orderForm.faultDescription || !orderForm.contactName || !orderForm.contactPhone) {
    ElMessage.warning('请填写完整的预约信息')
    return
  }

  if (!service.value) return

  submitting.value = true
  try {
    await repairApi.createOrder({
      technicianId: service.value.technicianId,
      serviceId: service.value.id,
      deviceType: orderForm.deviceType,
      deviceBrand: orderForm.deviceBrand,
      deviceModel: orderForm.deviceModel,
      faultDescription: orderForm.faultDescription,
      contactName: orderForm.contactName,
      contactPhone: orderForm.contactPhone,
      address: orderForm.address,
      servicePrice: service.value.price
    })
    ElMessage.success('预约成功')
    showOrderDialog.value = false
    router.push('/user/repair-orders')
  } catch (error: any) {
    ElMessage.error(error.message || '预约失败')
  } finally {
    submitting.value = false
  }
}

onMounted(() => {
  fetchService()
})
</script>

<style lang="scss" scoped>
.service-detail-page {
  min-height: 100vh;
  background: #f5f7fa;
  padding: 20px 0;
}

.breadcrumb {
  margin-bottom: 20px;
}

.service-content {
  padding: 30px;
}

.service-header {
  display: flex;
  gap: 30px;
  margin-bottom: 30px;

  .service-icon {
    width: 120px;
    height: 120px;
    background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
    flex-shrink: 0;
  }

  .service-info {
    flex: 1;

    .service-title {
      font-size: 28px;
      margin-bottom: 16px;
    }

    .service-tags {
      display: flex;
      gap: 12px;
      margin-bottom: 20px;
    }

    .service-price {
      display: flex;
      align-items: baseline;
      gap: 16px;

      .current-price {
        font-size: 36px;
        color: var(--danger-color);
        font-weight: 600;
      }

      .price-range {
        font-size: 14px;
        color: var(--text-lighter-color);
      }
    }
  }
}

.service-description {
  margin-bottom: 30px;

  h3 {
    margin-bottom: 16px;
    padding-bottom: 12px;
    border-bottom: 1px solid #f0f0f0;
  }

  p {
    line-height: 1.8;
    color: var(--text-light-color);
    white-space: pre-wrap;
  }
}

.technician-info {
  margin-bottom: 30px;

  h3 {
    margin-bottom: 20px;
    padding-bottom: 12px;
    border-bottom: 1px solid #f0f0f0;
  }

  .technician-detail {
    display: flex;
    gap: 20px;
    align-items: center;

    .technician-info-detail {
      .technician-name {
        font-size: 18px;
        font-weight: 500;
        margin-bottom: 8px;
      }

      .technician-stats {
        display: flex;
        gap: 20px;
        color: var(--text-lighter-color);
        font-size: 14px;

        span {
          display: flex;
          align-items: center;
          gap: 8px;
        }
      }
    }
  }
}

.action-section {
  text-align: center;

  .el-button {
    width: 200px;
  }
}
</style>
