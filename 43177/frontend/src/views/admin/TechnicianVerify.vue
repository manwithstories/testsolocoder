<template>
  <div class="technician-verify-page">
    <div class="page-header">
      <h2 class="page-title">技师审核</h2>
    </div>

    <el-card>
      <div class="filter-bar">
        <el-select v-model="statusFilter" placeholder="选择状态" clearable style="width: 150px;" @change="loadList">
          <el-option label="待审核" value="pending" />
          <el-option label="已通过" value="approved" />
          <el-option label="已拒绝" value="rejected" />
        </el-select>
      </div>

      <el-table :data="profiles" style="width: 100%">
        <el-table-column label="用户名" width="150">
          <template #default="{ row }">{{ row.user?.username }}</template>
        </el-table-column>
        <el-table-column label="真实姓名" width="150">
          <template #default="{ row }">{{ row.user?.real_name || '-' }}</template>
        </el-table-column>
        <el-table-column label="手机号" width="150">
          <template #default="{ row }">{{ row.user?.phone }}</template>
        </el-table-column>
        <el-table-column prop="specialty" label="专业" width="150" />
        <el-table-column prop="experience_years" label="从业年限" width="100">
          <template #default="{ row }">{{ row.experience_years }}年</template>
        </el-table-column>
        <el-table-column prop="certificate_no" label="证书编号" width="200" />
        <el-table-column prop="verify_status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.verify_status)">
              {{ getStatusText(row.verify_status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button size="small" @click="viewDetail(row)">查看</el-button>
            <el-button
              v-if="row.verify_status === 'pending'"
              type="success"
              size="small"
              @click="approve(row)"
            >
              通过
            </el-button>
            <el-button
              v-if="row.verify_status === 'pending'"
              type="danger"
              size="small"
              @click="reject(row)"
            >
              拒绝
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <div v-if="total > pageSize" class="pagination">
        <el-pagination
          v-model:current-page="currentPage"
          :page-size="pageSize"
          :total="total"
          layout="prev, pager, next"
          @current-change="loadList"
        />
      </div>
    </el-card>

    <el-dialog v-model="showDetailDialog" title="技师详情" width="600px">
      <el-descriptions :column="2" border v-if="selectedProfile">
        <el-descriptions-item label="用户名">
          {{ (selectedProfile as any).user?.username }}
        </el-descriptions-item>
        <el-descriptions-item label="真实姓名">
          {{ (selectedProfile as any).user?.real_name || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="手机号">
          {{ (selectedProfile as any).user?.phone }}
        </el-descriptions-item>
        <el-descriptions-item label="邮箱">
          {{ (selectedProfile as any).user?.email || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="专业">
          {{ selectedProfile.specialty || '-' }}
        </el-descriptions-item>
        <el-descriptions-item label="从业年限">
          {{ selectedProfile.experience_years }}年
        </el-descriptions-item>
        <el-descriptions-item label="证书编号">
          {{ selectedProfile.certificate_no }}
        </el-descriptions-item>
        <el-descriptions-item label="服务范围">
          {{ selectedProfile.service_radius }}km
        </el-descriptions-item>
      </el-descriptions>

      <div v-if="selectedProfile?.certificate_image" class="certificate-image">
        <img :src="selectedProfile.certificate_image" alt="证书" />
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { adminApi } from '@/api/admin'
import type { TechnicianProfile } from '@/types'

const profiles = ref<TechnicianProfile[]>([])
const currentPage = ref(1)
const pageSize = ref(10)
const total = ref(0)
const statusFilter = ref('pending')
const showDetailDialog = ref(false)
const selectedProfile = ref<TechnicianProfile | null>(null)

onMounted(() => {
  loadList()
})

async function loadList() {
  try {
    const res = await adminApi.getTechnicianVerifyList({
      page: currentPage.value,
      page_size: pageSize.value,
      status: statusFilter.value || undefined
    })
    profiles.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (error) {
    console.error('Failed to load profiles:', error)
  }
}

function getStatusType(status: string): string {
  const typeMap: Record<string, string> = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger'
  }
  return typeMap[status] || 'info'
}

function getStatusText(status: string): string {
  const textMap: Record<string, string> = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝'
  }
  return textMap[status] || status
}

function viewDetail(profile: TechnicianProfile) {
  selectedProfile.value = profile
  showDetailDialog.value = true
}

async function approve(profile: TechnicianProfile) {
  try {
    await ElMessageBox.confirm('确定要通过该技师的认证吗？', '确认', { type: 'success' })
    await adminApi.verifyTechnician(profile.id, { is_verified: true })
    ElMessage.success('审核通过')
    loadList()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to approve:', error)
    }
  }
}

async function reject(profile: TechnicianProfile) {
  try {
    const { value: remark } = await ElMessageBox.prompt('请输入拒绝原因', '拒绝认证', {
      inputValidator: (value) => {
        if (!value) return '请输入拒绝原因'
        return true
      }
    })
    await adminApi.verifyTechnician(profile.id, {
      is_verified: false,
      verify_remark: remark
    })
    ElMessage.success('已拒绝')
    loadList()
  } catch (error: any) {
    if (error !== 'cancel') {
      console.error('Failed to reject:', error)
    }
  }
}
</script>

<style scoped>
.technician-verify-page {
  padding: 0;
}

.filter-bar {
  display: flex;
  gap: 15px;
  margin-bottom: 20px;
}

.pagination {
  display: flex;
  justify-content: center;
  margin-top: 20px;
}

.certificate-image {
  margin-top: 20px;
  text-align: center;
}

.certificate-image img {
  max-width: 100%;
  max-height: 300px;
  border-radius: 8px;
}
</style>
