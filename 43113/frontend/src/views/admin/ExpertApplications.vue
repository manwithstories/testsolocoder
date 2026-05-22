<template>
  <div class="expert-applications-page">
    <el-card>
      <template #header>
        <div class="card-header">
          <span>专家审核</span>
          <el-select v-model="statusFilter" placeholder="选择状态" clearable style="width: 150px" @change="fetchApplications">
            <el-option label="待审核" value="pending" />
            <el-option label="已通过" value="approved" />
            <el-option label="已拒绝" value="rejected" />
          </el-select>
        </div>
      </template>

      <el-table :data="applications" style="width: 100%">
        <el-table-column label="ID" prop="id" width="80" />
        <el-table-column label="申请人" width="200">
          <template #default="{ row }">
            <div class="user-cell">
              <el-avatar :size="32" :src="row.user?.avatar">
                {{ row.user?.nickname?.charAt(0) || 'U' }}
              </el-avatar>
              <div class="user-info">
                <span class="username">{{ row.user?.nickname || row.user?.username }}</span>
                <span class="level">Lv.{{ row.user?.level }}</span>
              </div>
            </div>
          </template>
        </el-table-column>
        <el-table-column label="擅长领域" prop="field" />
        <el-table-column label="个人描述" prop="description" show-overflow-tooltip />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusName(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="申请时间" width="180">
          <template #default="{ row }">
            {{ formatTime(row.createdAt) }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200">
          <template #default="{ row }">
            <el-button
              v-if="row.status === 'pending'"
              type="success"
              size="small"
              @click="reviewApplication(row.id, 'approved', '')"
            >
              通过
            </el-button>
            <el-button
              v-if="row.status === 'pending'"
              type="danger"
              size="small"
              @click="reviewApplication(row.id, 'rejected', '不符合要求')"
            >
              拒绝
            </el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pagination"
        v-model:current-page="page"
        v-model:page-size="pageSize"
        :total="total"
        layout="prev, pager, next, total"
        @current-change="fetchApplications"
      />
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { userApi } from '@/api/user'
import type { ExpertApplication } from '@/types'
import dayjs from 'dayjs'

const applications = ref<ExpertApplication[]>([])
const page = ref(1)
const pageSize = ref(10)
const total = ref(0)
const statusFilter = ref('')

const fetchApplications = async () => {
  try {
    const res = await userApi.getExpertApplications({
      page: page.value,
      pageSize: pageSize.value,
      status: statusFilter.value || undefined
    })
    applications.value = res.data?.list || []
    total.value = res.data?.total || 0
  } catch (e) {
    console.error(e)
  }
}

const reviewApplication = async (id: number, status: string, remark: string) => {
  try {
    await userApi.reviewExpertApplication(id, { status, remark })
    fetchApplications()
  } catch (e) {
    console.error(e)
  }
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    pending: 'warning',
    approved: 'success',
    rejected: 'danger'
  }
  return map[status] || 'info'
}

const getStatusName = (status: string) => {
  const map: Record<string, string> = {
    pending: '待审核',
    approved: '已通过',
    rejected: '已拒绝'
  }
  return map[status] || status
}

const formatTime = (time: string) => {
  return dayjs(time).format('YYYY-MM-DD HH:mm:ss')
}

onMounted(() => {
  fetchApplications()
})
</script>

<style scoped lang="scss">
.expert-applications-page {
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }

  .user-cell {
    display: flex;
    align-items: center;
    gap: 12px;

    .user-info {
      display: flex;
      flex-direction: column;
      gap: 4px;

      .username {
        font-weight: 500;
      }

      .level {
        font-size: 12px;
        color: #909399;
      }
    }
  }

  .pagination {
    justify-content: center;
    margin-top: 20px;
  }
}
</style>
