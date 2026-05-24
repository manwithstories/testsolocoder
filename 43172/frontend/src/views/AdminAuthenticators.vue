<template>
  <div class="page-container">
    <div class="card">
      <div class="page-header">
        <h2>鉴定师审核</h2>
        <div class="filter-bar">
          <el-radio-group v-model="status" size="default" @change="loadAuthenticators">
            <el-radio-button label="">全部</el-radio-button>
            <el-radio-button label="pending">待审核</el-radio-button>
            <el-radio-button label="approved">已通过</el-radio-button>
            <el-radio-button label="rejected">已拒绝</el-radio-button>
          </el-radio-group>
        </div>
      </div>
    </div>

    <div v-if="loading" class="loading-wrapper">
      <el-icon class="is-loading" :size="32"><Loading /></el-icon>
    </div>

    <el-table
      v-else-if="authenticators.length > 0"
      :data="authenticators"
      style="width: 100%"
      stripe
    >
      <el-table-column label="ID" prop="id" width="80" />
      <el-table-column label="用户" min-width="150">
        <template #default="{ row }">
          <div class="user-cell">
            <el-avatar :size="32" :src="row.user?.avatar">
              {{ row.user?.username?.charAt(0)?.toUpperCase() }}
            </el-avatar>
            <span>{{ row.user?.username }}</span>
          </div>
        </template>
      </el-table-column>
      <el-table-column label="资质编号" prop="license_number" width="150" />
      <el-table-column label="专业领域" prop="specialties" min-width="200" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">{{ getStatusLabel(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="评分" width="80">
        <template #default="{ row }">
          {{ row.rating.toFixed(1) }}
        </template>
      </el-table-column>
      <el-table-column label="累计鉴定" prop="completed_count" width="100" />
      <el-table-column label="申请时间" width="180">
        <template #default="{ row }">
          {{ formatDate(row.created_at) }}
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.status === 'pending'"
            type="success"
            size="small"
            @click="handleApprove(row)"
          >
            通过
          </el-button>
          <el-button
            v-if="row.status === 'pending'"
            type="danger"
            size="small"
            @click="handleReject(row)"
          >
            拒绝
          </el-button>
          <el-button
            v-if="row.license_file"
            size="small"
            @click="viewLicense(row)"
          >
            查看资质
          </el-button>
        </template>
      </el-table-column>
    </el-table>

    <div v-else class="empty-state">
      <el-icon :size="64"><Box /></el-icon>
      <p>暂无鉴定师申请</p>
    </div>

    <div v-if="!loading && total > 0" class="pagination-wrapper">
      <el-pagination
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="total, prev, pager, next"
        @current-change="loadAuthenticators"
        @size-change="loadAuthenticators"
      />
    </div>

    <el-dialog v-model="rejectDialogVisible" title="拒绝原因" width="400px">
      <el-form :model="rejectForm">
        <el-form-item label="拒绝原因">
          <el-input
            v-model="rejectForm.reason"
            type="textarea"
            :rows="3"
            placeholder="请输入拒绝原因"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="rejectDialogVisible = false">取消</el-button>
        <el-button type="danger" :loading="rejecting" @click="confirmReject">
          确认拒绝
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { userApi } from '@/api/user'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { AuthenticatorProfile } from '@/types'
import dayjs from 'dayjs'
import { Loading, Box } from '@element-plus/icons-vue'

const authenticators = ref<AuthenticatorProfile[]>([])
const loading = ref(false)
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const status = ref('')

const rejectDialogVisible = ref(false)
const currentAuthenticator = ref<AuthenticatorProfile | null>(null)
const rejecting = ref(false)

const rejectForm = ref({ reason: '' })

const loadAuthenticators = async () => {
  loading.value = true
  try {
    const res = await userApi.listAuthenticators({
      page: page.value,
      page_size: pageSize.value,
      status: status.value || undefined
    })
    if (res.code === 200) {
      const data = res.data as any
      authenticators.value = data?.list || []
      total.value = data?.total || 0
    }
  } catch (error) {
    console.error('Load authenticators error:', error)
  } finally {
    loading.value = false
  }
}

const getStatusType = (status: string) => {
  const types: Record<string, string> = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger'
  }
  return types[status] || 'info'
}

const getStatusLabel = (status: string) => {
  const labels: Record<string, string> = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝'
  }
  return labels[status] || status
}

const formatDate = (date: string) => dayjs(date).format('YYYY-MM-DD HH:mm:ss')

const handleApprove = async (auth: AuthenticatorProfile) => {
  try {
    await ElMessageBox.confirm(
      `确认通过鉴定师 "${auth.user?.username}" 的申请吗？`,
      '通过确认',
      {
        confirmButtonText: '确认',
        cancelButtonText: '取消',
        type: 'warning'
      }
    )

    const res = await userApi.approveAuthenticator(auth.user_id)
    if (res.code === 200) {
      ElMessage.success('已通过')
      loadAuthenticators()
    }
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Approve error:', error)
    }
  }
}

const handleReject = (auth: AuthenticatorProfile) => {
  currentAuthenticator.value = auth
  rejectForm.value.reason = ''
  rejectDialogVisible.value = true
}

const confirmReject = async () => {
  if (!currentAuthenticator.value) return

  if (!rejectForm.value.reason.trim()) {
    ElMessage.warning('请输入拒绝原因')
    return
  }

  rejecting.value = true
  try {
    const res = await userApi.rejectAuthenticator(
      currentAuthenticator.value.user_id,
      rejectForm.value
    )
    if (res.code === 200) {
      ElMessage.success('已拒绝')
      rejectDialogVisible.value = false
      loadAuthenticators()
    }
  } catch (error) {
    console.error('Reject error:', error)
  } finally {
    rejecting.value = false
  }
}

const viewLicense = (auth: AuthenticatorProfile) => {
  if (auth.license_file) {
    window.open(auth.license_file)
  }
}

onMounted(() => {
  loadAuthenticators()
})
</script>

<style lang="scss" scoped>
.page-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  
  h2 {
    font-size: 20px;
    font-weight: 600;
  }
}

.filter-bar {
  display: flex;
  gap: 12px;
}

.user-cell {
  display: flex;
  align-items: center;
  gap: 8px;
}

.pagination-wrapper {
  display: flex;
  justify-content: center;
  margin-top: 24px;
}
</style>
