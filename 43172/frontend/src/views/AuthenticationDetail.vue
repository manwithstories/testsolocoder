<template>
  <div class="page-container">
    <div v-if="loading" class="loading-wrapper">
      <el-icon class="is-loading" :size="32"><Loading /></el-icon>
    </div>

    <template v-else-if="authentication">
      <div class="card">
        <h2 class="section-title">鉴定详情</h2>
        <el-descriptions :column="2" border>
          <el-descriptions-item label="鉴定ID">{{ authentication.id }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getAuthStatusType(authentication.status)">
              {{ getAuthStatusLabel(authentication.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="鉴定结果">
            <el-tag v-if="authentication.result" :type="getAuthResultType(authentication.result)">
              {{ getAuthResultLabel(authentication.result) }}
            </el-tag>
            <span v-else>待鉴定</span>
          </el-descriptions-item>
          <el-descriptions-item label="申请时间">
            {{ formatDate(authentication.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item v-if="authentication.accepted_at" label="接单时间">
            {{ formatDate(authentication.accepted_at) }}
          </el-descriptions-item>
          <el-descriptions-item v-if="authentication.completed_at" label="完成时间">
            {{ formatDate(authentication.completed_at) }}
          </el-descriptions-item>
        </el-descriptions>
      </div>

      <div class="card">
        <h3 class="section-title">商品信息</h3>
        <div v-if="authentication.product" class="product-info">
          <img v-if="authentication.product.images?.[0]" :src="authentication.product.images[0].image_url" />
          <div class="product-details">
            <h4>{{ authentication.product.title }}</h4>
            <p>品牌: {{ authentication.product.brand_name || '其他品牌' }}</p>
            <p class="price-text">¥{{ authentication.product.price.toFixed(2) }}</p>
          </div>
        </div>
      </div>

      <div class="card">
        <h3 class="section-title">买家信息</h3>
        <p v-if="authentication.buyer">买家: {{ authentication.buyer.username }}</p>
      </div>

      <div v-if="authentication.authenticator" class="card">
        <h3 class="section-title">鉴定师信息</h3>
        <p>鉴定师: {{ authentication.authenticator.username }}</p>
        <p v-if="authentication.authenticator.authenticator_profile">
          资质编号: {{ authentication.authenticator.authenticator_profile.license_number }}
        </p>
      </div>

      <div v-if="authentication.report_content" class="card">
        <h3 class="section-title">鉴定报告</h3>
        <p class="report-content">{{ authentication.report_content }}</p>
        <div v-if="authentication.report_file" class="report-download">
          <el-button type="primary" :icon="Download" @click="downloadReport">
            下载PDF报告
          </el-button>
        </div>
      </div>

      <div v-if="authentication.authenticator_notes" class="card">
        <h3 class="section-title">鉴定师备注</h3>
        <p>{{ authentication.authenticator_notes }}</p>
      </div>

      <div class="card">
        <div class="action-buttons">
          <el-button
            v-if="authentication.status === 'pending' && userStore.userRole === 'authenticator'"
            type="primary"
            @click="handleAccept"
          >
            接单
          </el-button>
          <el-button
            v-if="authentication.status === 'accepted' && isMyTask"
            type="success"
            @click="completeDialogVisible = true"
          >
            完成鉴定
          </el-button>
        </div>
      </div>
    </template>

    <div v-else class="empty-state">
      <el-icon :size="64"><Box /></el-icon>
      <p>鉴定记录不存在</p>
    </div>

    <el-dialog v-model="completeDialogVisible" title="完成鉴定" width="600px">
      <el-form :model="completeForm" label-width="100px">
        <el-form-item label="鉴定结果" required>
          <el-radio-group v-model="completeForm.result">
            <el-radio-button label="genuine">正品</el-radio-button>
            <el-radio-button label="counterfeit">赝品</el-radio-button>
            <el-radio-button label="inconclusive">无法鉴定</el-radio-button>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="鉴定报告">
          <el-input
            v-model="completeForm.report_content"
            type="textarea"
            :rows="6"
            placeholder="请输入详细的鉴定报告"
          />
        </el-form-item>
        <el-form-item label="备注">
          <el-input
            v-model="completeForm.authenticator_notes"
            type="textarea"
            :rows="3"
            placeholder="可选，输入备注信息"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="completeDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="completing" @click="confirmComplete">
          提交鉴定结果
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute } from 'vue-router'
import { authServiceApi } from '@/api/authentication'
import { useUserStore } from '@/stores/user'
import { ElMessage } from 'element-plus'
import type { Authentication } from '@/types'
import { AUTH_STATUS_OPTIONS, AUTH_RESULT_OPTIONS } from '@/types'
import dayjs from 'dayjs'
import { Loading, Box, Download } from '@element-plus/icons-vue'

const route = useRoute()
const userStore = useUserStore()

const authentication = ref<Authentication | null>(null)
const loading = ref(false)
const completeDialogVisible = ref(false)
const completing = ref(false)

const authId = Number(route.params.id)

const isMyTask = computed(() => 
  authentication.value?.authenticator_id === userStore.userID
)

const completeForm = ref({
  result: 'genuine',
  report_content: '',
  authenticator_notes: ''
})

const loadAuthentication = async () => {
  loading.value = true
  try {
    const res = await authServiceApi.getAuthentication(authId)
    if (res.code === 200 && res.data) {
      authentication.value = res.data
    }
  } catch (error) {
    console.error('Load authentication error:', error)
  } finally {
    loading.value = false
  }
}

const getAuthStatusType = (status: string) => {
  const opt = AUTH_STATUS_OPTIONS.find(o => o.value === status)
  return opt?.type || 'info'
}

const getAuthStatusLabel = (status: string) => {
  const opt = AUTH_STATUS_OPTIONS.find(o => o.value === status)
  return opt?.label || status
}

const getAuthResultType = (result: string) => {
  const opt = AUTH_RESULT_OPTIONS.find(o => o.value === result)
  return opt?.type || 'info'
}

const getAuthResultLabel = (result: string) => {
  const opt = AUTH_RESULT_OPTIONS.find(o => o.value === result)
  return opt?.label || result
}

const formatDate = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss')

const handleAccept = async () => {
  try {
    const res = await authServiceApi.acceptAuthentication(authId)
    if (res.code === 200) {
      ElMessage.success('接单成功')
      loadAuthentication()
    }
  } catch (error) {
    console.error('Accept error:', error)
  }
}

const confirmComplete = async () => {
  completing.value = true
  try {
    const res = await authServiceApi.completeAuthentication(authId, completeForm.value)
    if (res.code === 200) {
      ElMessage.success('鉴定完成')
      completeDialogVisible.value = false
      loadAuthentication()
    }
  } catch (error) {
    console.error('Complete error:', error)
  } finally {
    completing.value = false
  }
}

const downloadReport = () => {
  if (authentication.value?.report_file) {
    window.open(authServiceApi.downloadReport(authId))
  }
}

onMounted(() => {
  loadAuthentication()
})
</script>

<style lang="scss" scoped>
.product-info {
  display: flex;
  gap: 16px;
  align-items: center;
  
  img {
    width: 120px;
    height: 120px;
    object-fit: cover;
    border-radius: 8px;
  }
  
  .product-details {
    h4 {
      font-size: 16px;
      font-weight: 500;
      margin-bottom: 8px;
    }
    
    p {
      margin-bottom: 4px;
      color: var(--text-secondary);
    }
    
    .price-text {
      font-size: 18px;
      font-weight: 600;
      color: var(--danger-color);
    }
  }
}

.report-content {
  white-space: pre-wrap;
  line-height: 1.8;
  color: var(--text-secondary);
}

.report-download {
  margin-top: 16px;
  text-align: right;
}

.action-buttons {
  display: flex;
  gap: 12px;
  justify-content: flex-end;
}
</style>
