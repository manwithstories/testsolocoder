<template>
  <div class="application-detail">
    <div class="page-header flex-between">
      <h2 class="page-title">申请详情</h2>
      <el-button @click="goBack">返回</el-button>
    </div>

    <el-tabs v-model="activeTab" v-loading="loading">
      <el-tab-pane label="基本信息" name="info">
        <div class="card">
          <el-descriptions :column="2" border>
            <el-descriptions-item label="申请编号">{{ application?.applicationNo }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <span :class="['status-tag', getStatusClass(application?.status)]">
                {{ getStatusText(application?.status) }}
              </span>
            </el-descriptions-item>
            <el-descriptions-item label="公司名称">{{ application?.companyName }}</el-descriptions-item>
            <el-descriptions-item label="公司类型">{{ getCompanyTypeText(application?.companyType) }}</el-descriptions-item>
            <el-descriptions-item label="注册资本">{{ formatMoney(application?.registeredCapital) }}</el-descriptions-item>
            <el-descriptions-item label="注册地址">{{ application?.registeredAddress }}</el-descriptions-item>
            <el-descriptions-item label="经营范围" :span="2">{{ application?.businessScope }}</el-descriptions-item>
            <el-descriptions-item label="股东信息" :span="2">{{ application?.shareholderInfo }}</el-descriptions-item>
            <el-descriptions-item label="创业者">{{ application?.entrepreneur?.realName }}</el-descriptions-item>
            <el-descriptions-item label="代办专员">{{ application?.agent?.realName || '-' }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">{{ formatTime(application?.createdAt) }}</el-descriptions-item>
            <el-descriptions-item label="完成时间">{{ formatTime(application?.completedAt) || '-' }}</el-descriptions-item>
          </el-descriptions>

          <div class="form-section mt-24" v-if="application?.reviewComments">
            <div class="form-section-title">审核意见</div>
            <div class="review-comments">{{ application.reviewComments }}</div>
          </div>
        </div>
      </el-tab-pane>

      <el-tab-pane label="办理进度" name="process">
        <div class="card">
          <div class="progress-wrapper">
            <el-steps :active="currentStepIndex" finish-status="success" align-center>
              <el-step
                v-for="step in processSteps"
                :key="step.id"
                :title="step.stepName"
                :status="getStepStatus(step.status)"
              />
            </el-steps>
          </div>

          <el-table :data="processSteps" style="width: 100%; margin-top: 24px">
            <el-table-column prop="stepOrder" label="序号" width="80">
              <template #default="{ $index }">
                {{ $index + 1 }}
              </template>
            </el-table-column>
            <el-table-column prop="stepName" label="环节名称" />
            <el-table-column prop="status" label="状态" width="100">
              <template #default="{ row }">
                <span :class="['status-tag', getStepStatusClass(row.status)]">
                  {{ getStepStatusText(row.status) }}
                </span>
              </template>
            </el-table-column>
            <el-table-column prop="handler.realName" label="处理人" width="120">
              <template #default="{ row }">
                {{ row.handler?.realName || '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="startedAt" label="开始时间" width="160">
              <template #default="{ row }">
                {{ formatTime(row.startedAt) || '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="completedAt" label="完成时间" width="160">
              <template #default="{ row }">
                {{ formatTime(row.completedAt) || '-' }}
              </template>
            </el-table-column>
            <el-table-column prop="remark" label="备注" />
            <el-table-column label="操作" width="200" v-if="isAgent">
              <template #default="{ row }">
                <el-button
                  v-if="row.status === 'pending'"
                  type="primary"
                  link
                  @click="startStep(row)"
                >
                  开始处理
                </el-button>
                <el-button
                  v-if="row.status === 'in_progress'"
                  type="success"
                  link
                  @click="completeStep(row)"
                >
                  完成
                </el-button>
                <el-button type="info" link @click="viewStepDetail(row)">
                  详情
                </el-button>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-tab-pane>

      <el-tab-pane label="费用信息" name="fee">
        <div class="card">
          <div v-if="applicationFee">
            <el-descriptions :column="2" border>
              <el-descriptions-item label="总金额">{{ formatMoney(applicationFee.totalAmount) }}</el-descriptions-item>
              <el-descriptions-item label="优惠金额">{{ formatMoney(applicationFee.discountAmount) }}</el-descriptions-item>
              <el-descriptions-item label="实付金额">{{ formatMoney(applicationFee.paidAmount) }}</el-descriptions-item>
              <el-descriptions-item label="支付状态">
                <span :class="['status-tag', applicationFee.status === 'paid' ? 'completed' : 'pending']">
                  {{ applicationFee.status === 'paid' ? '已支付' : '待支付' }}
                </span>
              </el-descriptions-item>
              <el-descriptions-item label="支付方式">{{ applicationFee.paymentMethod || '-' }}</el-descriptions-item>
              <el-descriptions-item label="交易号">{{ applicationFee.transactionNo || '-' }}</el-descriptions-item>
            </el-descriptions>

            <h4 class="mt-24 mb-16">费用明细</h4>
            <el-table :data="applicationFee.feeItems" style="width: 100%">
              <el-table-column prop="itemName" label="项目" />
              <el-table-column prop="amount" label="金额" width="150">
                <template #default="{ row }">
                  {{ formatMoney(row.amount) }}
                </template>
              </el-table-column>
              <el-table-column prop="description" label="说明" />
            </el-table>

            <div v-if="applicationFee.status === 'pending'" class="mt-24">
              <el-button type="primary" size="large" @click="handlePay">
                在线支付
              </el-button>
            </div>
          </div>
          <div v-else class="empty-state">
            <el-icon :size="48"><Wallet /></el-icon>
            <p>暂无费用信息</p>
          </div>
        </div>
      </el-tab-pane>

      <el-tab-pane label="材料文件" name="materials">
        <div class="card">
          <el-row :gutter="20">
            <el-col :span="8" v-if="application?.idCardFront">
              <div class="material-item">
                <div class="material-label">身份证正面</div>
                <el-image :src="getImageUrl(application.idCardFront)" :preview-src-list="[getImageUrl(application.idCardFront)]" fit="cover" class="material-image" />
              </div>
            </el-col>
            <el-col :span="8" v-if="application?.idCardBack">
              <div class="material-item">
                <div class="material-label">身份证反面</div>
                <el-image :src="getImageUrl(application.idCardBack)" :preview-src-list="[getImageUrl(application.idCardBack)]" fit="cover" class="material-image" />
              </div>
            </el-col>
            <el-col :span="8" v-if="application?.licensePreview">
              <div class="material-item">
                <div class="material-label">营业执照预审</div>
                <el-image :src="getImageUrl(application.licensePreview)" :preview-src-list="[getImageUrl(application.licensePreview)]" fit="cover" class="material-image" />
              </div>
            </el-col>
          </el-row>
        </div>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="stepDialogVisible" title="环节详情" width="600px">
      <el-descriptions :column="1" border v-if="currentStep">
        <el-descriptions-item label="环节名称">{{ currentStep.stepName }}</el-descriptions-item>
        <el-descriptions-item label="状态">{{ getStepStatusText(currentStep.status) }}</el-descriptions-item>
        <el-descriptions-item label="描述">{{ currentStep.description || '-' }}</el-descriptions-item>
        <el-descriptions-item label="备注">{{ currentStep.remark || '-' }}</el-descriptions-item>
        <el-descriptions-item label="证明文件">
          <el-link v-if="currentStep.certificateFile" :href="getImageUrl(currentStep.certificateFile)" target="_blank">
            查看文件
          </el-link>
          <span v-else>-</span>
        </el-descriptions-item>
      </el-descriptions>
    </el-dialog>

    <el-dialog v-model="payDialogVisible" title="支付费用" width="500px">
      <el-form :model="payForm" label-width="80px">
        <el-form-item label="支付金额">
          <span class="pay-amount">{{ formatMoney(applicationFee?.paidAmount || 0) }}</span>
        </el-form-item>
        <el-form-item label="支付方式">
          <el-radio-group v-model="payForm.method">
            <el-radio label="alipay">支付宝</el-radio>
            <el-radio label="wechat">微信支付</el-radio>
            <el-radio label="bank">银行卡</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="payDialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="paying" @click="confirmPay">确认支付</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed, watch } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Money } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/user'
import { applicationApi } from '@/api/application'
import { processApi } from '@/api/process'
import { feeApi } from '@/api/fee'
import { Application, ProcessStep, ApplicationFee, ApplicationStatus, ProcessStepStatus, CompanyType } from '@/types'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const application = ref<Application | null>(null)
const processSteps = ref<ProcessStep[]>([])
const applicationFee = ref<ApplicationFee | null>(null)
const activeTab = ref('info')
const stepDialogVisible = ref(false)
const currentStep = ref<ProcessStep | null>(null)
const payDialogVisible = ref(false)
const paying = ref(false)

const isAgent = computed(() => userStore.userRole === 'agent')

const payForm = reactive({
  method: 'alipay'
})

const currentStepIndex = computed(() => {
  const index = processSteps.value.findIndex(s => s.status === 'in_progress')
  return index >= 0 ? index : processSteps.value.filter(s => s.status === 'completed').length
})

const fetchApplication = async () => {
  loading.value = true
  try {
    const id = Number(route.params.id)
    const res = await applicationApi.getById(id)
    application.value = res || null
    processSteps.value = res?.processSteps || []
    applicationFee.value = res?.fee || null
  } catch (error) {
    console.error('获取申请详情失败:', error)
  } finally {
    loading.value = false
  }
}

const getStatusText = (status?: ApplicationStatus) => {
  if (!status) return ''
  const map: Record<ApplicationStatus, string> = {
    draft: '草稿',
    pending_review: '待审核',
    reviewing: '审核中',
    processing: '处理中',
    completed: '已完成',
    rejected: '已驳回',
    cancelled: '已取消',
    payment_pending: '待支付'
  }
  return map[status] || status
}

const getStatusClass = (status?: ApplicationStatus) => {
  if (!status) return ''
  const map: Record<ApplicationStatus, string> = {
    draft: 'draft',
    pending_review: 'pending',
    reviewing: 'pending',
    processing: 'processing',
    completed: 'completed',
    rejected: 'rejected',
    cancelled: 'cancelled',
    payment_pending: 'pending'
  }
  return map[status] || ''
}

const getCompanyTypeText = (type?: CompanyType) => {
  if (!type) return ''
  const map: Record<CompanyType, string> = {
    llc: '有限责任公司',
    joint_stock: '股份有限公司',
    sole: '个人独资',
    partnership: '合伙企业'
  }
  return map[type] || type
}

const getStepStatus = (status: ProcessStepStatus) => {
  const map: Record<ProcessStepStatus, string> = {
    pending: 'wait',
    in_progress: 'process',
    completed: 'finish',
    failed: 'error',
    skipped: 'success'
  }
  return map[status] || ''
}

const getStepStatusText = (status: ProcessStepStatus) => {
  const map: Record<ProcessStepStatus, string> = {
    pending: '待处理',
    in_progress: '处理中',
    completed: '已完成',
    failed: '失败',
    skipped: '已跳过'
  }
  return map[status] || status
}

const getStepStatusClass = (status: ProcessStepStatus) => {
  const map: Record<ProcessStepStatus, string> = {
    pending: 'draft',
    in_progress: 'processing',
    completed: 'completed',
    failed: 'rejected',
    skipped: 'cancelled'
  }
  return map[status] || ''
}

const formatMoney = (amount?: number) => {
  if (!amount) return '¥0'
  return `¥${amount.toLocaleString()}`
}

const formatTime = (time?: string | null) => {
  if (!time) return ''
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

const getImageUrl = (path: string) => {
  if (!path) return ''
  return `/uploads/${path}`
}

const startStep = async (step: ProcessStep) => {
  try {
    await processApi.startStep(application.value!.id, step.id)
    ElMessage.success('已开始处理')
    fetchApplication()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  }
}

const completeStep = async (step: ProcessStep) => {
  try {
    await processApi.completeStep(application.value!.id, step.id)
    ElMessage.success('已完成')
    fetchApplication()
  } catch (error: any) {
    ElMessage.error(error.message || '操作失败')
  }
}

const viewStepDetail = (step: ProcessStep) => {
  currentStep.value = step
  stepDialogVisible.value = true
}

const handlePay = () => {
  payDialogVisible.value = true
}

const confirmPay = async () => {
  if (!application.value) return

  paying.value = true
  try {
    await feeApi.pay({
      applicationId: application.value.id,
      paymentMethod: payForm.method
    })
    ElMessage.success('支付成功')
    payDialogVisible.value = false
    fetchApplication()
  } catch (error: any) {
    ElMessage.error(error.message || '支付失败')
  } finally {
    paying.value = false
  }
}

const goBack = () => {
  router.back()
}

watch(() => route.query.tab, (tab) => {
  if (tab) {
    activeTab.value = tab as string
  }
}, { immediate: true })

onMounted(fetchApplication)
</script>

<style scoped>
.material-item {
  text-align: center;
}

.material-label {
  font-size: 14px;
  color: #606266;
  margin-bottom: 8px;
}

.material-image {
  width: 100%;
  height: 200px;
  border-radius: 8px;
  overflow: hidden;
}

.review-comments {
  padding: 16px;
  background: #f5f7fa;
  border-radius: 4px;
  color: #606266;
}

.pay-amount {
  font-size: 24px;
  font-weight: 600;
  color: #f56c6c;
}
</style>
