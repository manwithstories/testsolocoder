<template>
  <div class="application-list">
    <div class="page-header flex-between">
      <h2 class="page-title">注册申请</h2>
      <el-button type="primary" @click="goToNewApplication" v-if="isEntrepreneur">
        <el-icon><Plus /></el-icon>
        新建申请
      </el-button>
    </div>

    <div class="filter-bar">
      <el-form :inline="true" :model="filterForm" @submit.prevent>
        <el-form-item label="状态">
          <el-select v-model="filterForm.status" placeholder="全部" clearable style="width: 160px">
            <el-option label="全部" value="" />
            <el-option label="草稿" value="draft" />
            <el-option label="待审核" value="pending_review" />
            <el-option label="审核中" value="reviewing" />
            <el-option label="处理中" value="processing" />
            <el-option label="已完成" value="completed" />
            <el-option label="已驳回" value="rejected" />
            <el-option label="已取消" value="cancelled" />
            <el-option label="待支付" value="payment_pending" />
          </el-select>
        </el-form-item>
        <el-form-item label="关键词">
          <el-input
            v-model="filterForm.keyword"
            placeholder="公司名称/申请编号"
            clearable
            style="width: 200px"
            @keyup.enter="fetchData"
          />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="fetchData">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
          <el-button @click="resetFilter">重置</el-button>
        </el-form-item>
      </el-form>
    </div>

    <div class="table-wrapper">
      <el-table :data="applications" style="width: 100%" v-loading="loading">
        <el-table-column prop="applicationNo" label="申请编号" width="160" />
        <el-table-column prop="companyName" label="公司名称" />
        <el-table-column prop="companyType" label="公司类型" width="120">
          <template #default="{ row }">
            {{ getCompanyTypeText(row.companyType) }}
          </template>
        </el-table-column>
        <el-table-column prop="registeredCapital" label="注册资本" width="120">
          <template #default="{ row }">
            {{ formatMoney(row.registeredCapital) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <span :class="['status-tag', getStatusClass(row.status)]">
              {{ getStatusText(row.status) }}
            </span>
          </template>
        </el-table-column>
        <el-table-column prop="progressPercent" label="进度" width="150">
          <template #default="{ row }">
            <el-progress :percentage="row.progressPercent" :status="getProgressStatus(row.status)" />
          </template>
        </el-table-column>
        <el-table-column label="创业者" width="120" v-if="!isEntrepreneur">
          <template #default="{ row }">
            {{ row.entrepreneur?.realName || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="代办专员" width="120" v-if="!isAgent">
          <template #default="{ row }">
            {{ row.agent?.realName || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="160">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="viewDetail(row.id)">详情</el-button>
            <el-button
              v-if="row.status === 'draft' && isEntrepreneur"
              type="success"
              link
              @click="submitApplication(row.id)"
            >
              提交
            </el-button>
            <el-button
              v-if="row.status === 'payment_pending' && isEntrepreneur"
              type="warning"
              link
              @click="goToPay(row.id)"
            >
              支付
            </el-button>
            <el-button
              v-if="row.status === 'draft' && isEntrepreneur"
              type="danger"
              link
              @click="cancelApplication(row.id)"
            >
              取消
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper mt-16">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :total="pagination.total"
          :page-sizes="[10, 20, 50, 100]"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchData"
          @current-change="fetchData"
        />
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus, Search } from '@element-plus/icons-vue'
import { useUserStore } from '@/store/user'
import { applicationApi } from '@/api/application'
import { Application, ApplicationStatus, CompanyType } from '@/types'
import dayjs from 'dayjs'

const router = useRouter()
const userStore = useUserStore()
const loading = ref(false)
const applications = ref<Application[]>([])

const isEntrepreneur = computed(() => userStore.userRole === 'entrepreneur')
const isAgent = computed(() => userStore.userRole === 'agent')

const filterForm = reactive({
  status: '',
  keyword: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const fetchData = async () => {
  loading.value = true
  try {
    const res = await applicationApi.getList({
      page: pagination.page,
      pageSize: pagination.pageSize,
      status: filterForm.status || undefined,
      keyword: filterForm.keyword || undefined
    })
    applications.value = res?.list || []
    pagination.total = res?.total || 0
  } catch (error) {
    console.error('获取申请列表失败:', error)
  } finally {
    loading.value = false
  }
}

const resetFilter = () => {
  filterForm.status = ''
  filterForm.keyword = ''
  pagination.page = 1
  fetchData()
}

const getStatusText = (status: ApplicationStatus) => {
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

const getStatusClass = (status: ApplicationStatus) => {
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

const getProgressStatus = (status: ApplicationStatus) => {
  if (status === 'completed') return 'success'
  if (status === 'rejected') return 'exception'
  return ''
}

const getCompanyTypeText = (type: CompanyType) => {
  const map: Record<CompanyType, string> = {
    llc: '有限责任公司',
    joint_stock: '股份有限公司',
    sole: '个人独资',
    partnership: '合伙企业'
  }
  return map[type] || type
}

const formatMoney = (amount: number) => {
  return `¥${amount.toLocaleString()}`
}

const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm')
}

const goToNewApplication = () => {
  router.push('/applications/new')
}

const viewDetail = (id: number) => {
  router.push(`/applications/${id}`)
}

const goToPay = (id: number) => {
  router.push(`/applications/${id}?tab=fee`)
}

const submitApplication = async (id: number) => {
  try {
    await ElMessageBox.confirm('确认提交该申请？提交后将进入审核流程。', '提示', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await applicationApi.submit(id)
    ElMessage.success('提交成功')
    fetchData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '提交失败')
    }
  }
}

const cancelApplication = async (id: number) => {
  try {
    await ElMessageBox.confirm('确认取消该申请？', '提示', {
      confirmButtonText: '确认',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await applicationApi.cancel(id)
    ElMessage.success('取消成功')
    fetchData()
  } catch (error: any) {
    if (error !== 'cancel') {
      ElMessage.error(error.message || '取消失败')
    }
  }
}

onMounted(fetchData)
</script>

<style scoped>
.pagination-wrapper {
  display: flex;
  justify-content: flex-end;
}
</style>
