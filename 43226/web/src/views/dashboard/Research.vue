<template>
  <div class="dashboard-research">
    <div class="card-shadow p-20">
      <div class="flex-between mb-20">
        <h2 class="page-title">学术申请审批</h2>
        <div class="action-buttons">
          <el-button type="primary" @click="exportData">
            <el-icon><Download /></el-icon> 导出数据
          </el-button>
        </div>
      </div>

      <el-tabs v-model="activeTab" type="border-card" @tab-change="handleTabChange">
        <el-tab-pane label="待审批" name="pending">
          <template #label>
            <span class="tab-label">
              <el-icon><Clock /></el-icon>
              待审批
              <el-badge v-if="pendingCount > 0" :value="pendingCount" class="ml-10" />
            </span>
          </template>
          <div class="tab-content">
            <el-form :inline="true" :model="query" class="mb-20">
              <el-form-item>
                <el-input
                  v-model="query.keyword"
                  placeholder="搜索申请人、机构"
                  clearable
                  @keyup.enter="fetchList"
                />
              </el-form-item>
              <el-form-item label="申请时间">
                <el-date-picker
                  v-model="query.start_date"
                  type="date"
                  placeholder="开始日期"
                  value-format="YYYY-MM-DD"
                />
                <span class="mx-10">至</span>
                <el-date-picker
                  v-model="query.end_date"
                  type="date"
                  placeholder="结束日期"
                  value-format="YYYY-MM-DD"
                />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="fetchList">查询</el-button>
                <el-button @click="resetQuery">重置</el-button>
              </el-form-item>
            </el-form>

            <el-table :data="list" v-loading="loading" border>
              <el-table-column prop="id" label="ID" width="60" />
              <el-table-column label="申请人" width="120">
                <template #default="{ row }">
                  <div class="flex-center gap-10">
                    <el-avatar :size="32" :src="row.user?.avatar">
                      {{ row.user?.nickname?.charAt(0) }}
                    </el-avatar>
                    <span>{{ row.user?.nickname }}</span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="申请藏品" width="150">
                <template #default="{ row }">
                  <div class="flex-center gap-10">
                    <img
                      v-if="row.collection?.image_url"
                      :src="row.collection.image_url"
                      style="width: 40px; height: 40px; object-fit: cover; border-radius: 4px;"
                    />
                    <div>
                      <div>{{ row.collection?.name }}</div>
                      <div class="text-xs text-gray">{{ row.collection?.code }}</div>
                    </div>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="institution" label="申请机构" width="150" show-overflow-tooltip />
              <el-table-column prop="purpose" label="申请用途" min-width="150" show-overflow-tooltip />
              <el-table-column prop="created_at" label="申请时间" width="160">
                <template #default="{ row }">{{ formatDateTime(row.created_at) }}</template>
              </el-table-column>
              <el-table-column label="操作" width="200" fixed="right">
                <template #default="{ row }">
                  <el-button type="primary" size="small" link @click="handleViewDetail(row)">查看</el-button>
                  <el-button type="success" size="small" link @click="handleApprove(row)">通过</el-button>
                  <el-button type="danger" size="small" link @click="handleReject(row)">拒绝</el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>

        <el-tab-pane label="已通过" name="approved">
          <template #label>
            <span class="tab-label">
              <el-icon><CircleCheck /></el-icon>
              已通过
            </span>
          </template>
          <div class="tab-content">
            <el-form :inline="true" :model="query" class="mb-20">
              <el-form-item>
                <el-input
                  v-model="query.keyword"
                  placeholder="搜索申请人、机构"
                  clearable
                  @keyup.enter="fetchList"
                />
              </el-form-item>
              <el-form-item label="审批时间">
                <el-date-picker
                  v-model="query.start_date"
                  type="date"
                  placeholder="开始日期"
                  value-format="YYYY-MM-DD"
                />
                <span class="mx-10">至</span>
                <el-date-picker
                  v-model="query.end_date"
                  type="date"
                  placeholder="结束日期"
                  value-format="YYYY-MM-DD"
                />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="fetchList">查询</el-button>
                <el-button @click="resetQuery">重置</el-button>
              </el-form-item>
            </el-form>

            <el-table :data="list" v-loading="loading" border>
              <el-table-column prop="id" label="ID" width="60" />
              <el-table-column label="申请人" width="120">
                <template #default="{ row }">
                  <div class="flex-center gap-10">
                    <el-avatar :size="32" :src="row.user?.avatar">
                      {{ row.user?.nickname?.charAt(0) }}
                    </el-avatar>
                    <span>{{ row.user?.nickname }}</span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="申请藏品" width="150">
                <template #default="{ row }">
                  <div class="flex-center gap-10">
                    <img
                      v-if="row.collection?.image_url"
                      :src="row.collection.image_url"
                      style="width: 40px; height: 40px; object-fit: cover; border-radius: 4px;"
                    />
                    <div>
                      <div>{{ row.collection?.name }}</div>
                      <div class="text-xs text-gray">{{ row.collection?.code }}</div>
                    </div>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="institution" label="申请机构" width="150" show-overflow-tooltip />
              <el-table-column label="审批人" width="100">
                <template #default="{ row }">{{ row.reviewer?.nickname || '-' }}</template>
              </el-table-column>
              <el-table-column prop="approved_at" label="通过时间" width="160">
                <template #default="{ row }">{{ row.approved_at ? formatDateTime(row.approved_at) : '-' }}</template>
              </el-table-column>
              <el-table-column label="操作" width="100" fixed="right">
                <template #default="{ row }">
                  <el-button type="primary" size="small" link @click="handleViewDetail(row)">查看</el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>

        <el-tab-pane label="已拒绝" name="rejected">
          <template #label>
            <span class="tab-label">
              <el-icon><CircleClose /></el-icon>
              已拒绝
            </span>
          </template>
          <div class="tab-content">
            <el-form :inline="true" :model="query" class="mb-20">
              <el-form-item>
                <el-input
                  v-model="query.keyword"
                  placeholder="搜索申请人、机构"
                  clearable
                  @keyup.enter="fetchList"
                />
              </el-form-item>
              <el-form-item label="审批时间">
                <el-date-picker
                  v-model="query.start_date"
                  type="date"
                  placeholder="开始日期"
                  value-format="YYYY-MM-DD"
                />
                <span class="mx-10">至</span>
                <el-date-picker
                  v-model="query.end_date"
                  type="date"
                  placeholder="结束日期"
                  value-format="YYYY-MM-DD"
                />
              </el-form-item>
              <el-form-item>
                <el-button type="primary" @click="fetchList">查询</el-button>
                <el-button @click="resetQuery">重置</el-button>
              </el-form-item>
            </el-form>

            <el-table :data="list" v-loading="loading" border>
              <el-table-column prop="id" label="ID" width="60" />
              <el-table-column label="申请人" width="120">
                <template #default="{ row }">
                  <div class="flex-center gap-10">
                    <el-avatar :size="32" :src="row.user?.avatar">
                      {{ row.user?.nickname?.charAt(0) }}
                    </el-avatar>
                    <span>{{ row.user?.nickname }}</span>
                  </div>
                </template>
              </el-table-column>
              <el-table-column label="申请藏品" width="150">
                <template #default="{ row }">
                  <div class="flex-center gap-10">
                    <img
                      v-if="row.collection?.image_url"
                      :src="row.collection.image_url"
                      style="width: 40px; height: 40px; object-fit: cover; border-radius: 4px;"
                    />
                    <div>
                      <div>{{ row.collection?.name }}</div>
                      <div class="text-xs text-gray">{{ row.collection?.code }}</div>
                    </div>
                  </div>
                </template>
              </el-table-column>
              <el-table-column prop="institution" label="申请机构" width="150" show-overflow-tooltip />
              <el-table-column label="审批人" width="100">
                <template #default="{ row }">{{ row.reviewer?.nickname || '-' }}</template>
              </el-table-column>
              <el-table-column prop="reviewed_at" label="拒绝时间" width="160">
                <template #default="{ row }">{{ row.reviewed_at ? formatDateTime(row.reviewed_at) : '-' }}</template>
              </el-table-column>
              <el-table-column label="操作" width="100" fixed="right">
                <template #default="{ row }">
                  <el-button type="primary" size="small" link @click="handleViewDetail(row)">查看</el-button>
                </template>
              </el-table-column>
            </el-table>
          </div>
        </el-tab-pane>
      </el-tabs>

      <div class="pagination mt-20">
        <el-pagination
          v-model:current-page="query.page"
          v-model:page-size="query.page_size"
          :total="total"
          layout="total, sizes, prev, pager, next, jumper"
          :page-sizes="[10, 20, 50, 100]"
          @size-change="fetchList"
          @current-change="fetchList"
        />
      </div>
    </div>

    <el-dialog v-model="showDetail" title="申请详情" width="600px">
      <div v-if="currentApplication" class="application-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="申请编号">
            <el-tag type="primary">#{{ currentApplication.id }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="申请状态">
            <el-tag :type="statusType(currentApplication.status)">
              {{ statusText(currentApplication.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="申请人">
            <div class="flex-center gap-10">
              <el-avatar :size="28" :src="currentApplication.user?.avatar">
                {{ currentApplication.user?.nickname?.charAt(0) }}
              </el-avatar>
              <span>{{ currentApplication.user?.nickname }}</span>
            </div>
          </el-descriptions-item>
          <el-descriptions-item label="申请机构">
            {{ currentApplication.institution }}
          </el-descriptions-item>
          <el-descriptions-item label="申请藏品" :span="2">
            <div class="flex-center gap-10">
              <img
                v-if="currentApplication.collection?.image_url"
                :src="currentApplication.collection.image_url"
                style="width: 50px; height: 50px; object-fit: cover; border-radius: 4px;"
              />
              <div>
                <div class="font-medium">{{ currentApplication.collection?.name }}</div>
                <div class="text-sm text-gray">{{ currentApplication.collection?.code }}</div>
              </div>
            </div>
          </el-descriptions-item>
          <el-descriptions-item label="申请用途" :span="2">
            <div style="white-space: pre-wrap;">{{ currentApplication.purpose }}</div>
          </el-descriptions-item>
          <el-descriptions-item label="申请时间">
            {{ formatDateTime(currentApplication.created_at) }}
          </el-descriptions-item>
          <el-descriptions-item label="审批时间">
            {{ currentApplication.reviewed_at ? formatDateTime(currentApplication.reviewed_at) : '-' }}
          </el-descriptions-item>
          <el-descriptions-item v-if="currentApplication.reviewer" label="审批人">
            <div class="flex-center gap-10">
              <el-avatar :size="24" :src="currentApplication.reviewer?.avatar">
                {{ currentApplication.reviewer?.nickname?.charAt(0) }}
              </el-avatar>
              <span>{{ currentApplication.reviewer?.nickname }}</span>
            </div>
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="currentApplication.review_comment" class="mt-20">
          <h4 class="mb-10">审批意见</h4>
          <div class="review-comment">
            {{ currentApplication.review_comment }}
          </div>
        </div>
      </div>
      <template #footer>
        <el-button @click="showDetail = false">关闭</el-button>
        <el-button v-if="currentApplication?.status === 'pending'" type="success" @click="handleApprove(currentApplication)">
          通过申请
        </el-button>
        <el-button v-if="currentApplication?.status === 'pending'" type="danger" @click="handleReject(currentApplication)">
          拒绝申请
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showApproveDialog" title="通过申请" width="500px">
      <el-alert
        title="确定通过该学术申请吗？"
        type="success"
        show-icon
        :closable="false"
        class="mb-20"
      />
      <el-form :model="reviewForm" :rules="reviewRules" ref="reviewFormRef" label-width="100px">
        <el-form-item label="审批意见" prop="review_comment">
          <el-input
            v-model="reviewForm.review_comment"
            type="textarea"
            :rows="4"
            placeholder="请输入审批意见（可选）"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showApproveDialog = false">取消</el-button>
        <el-button type="success" :loading="submitting" @click="submitApprove">确认通过</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showRejectDialog" title="拒绝申请" width="500px">
      <el-alert
        title="请填写拒绝原因"
        type="warning"
        show-icon
        :closable="false"
        class="mb-20"
      />
      <el-form :model="reviewForm" :rules="rejectRules" ref="reviewFormRef" label-width="100px">
        <el-form-item label="拒绝原因" prop="review_comment">
          <el-input
            v-model="reviewForm.review_comment"
            type="textarea"
            :rows="4"
            placeholder="请输入拒绝原因"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showRejectDialog = false">取消</el-button>
        <el-button type="danger" :loading="submitting" @click="submitReject">确认拒绝</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, computed } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Download, Clock, CircleCheck, CircleClose } from '@element-plus/icons-vue'
import * as guideApi from '@/api/guide'
import type { ResearchApplication } from '@/types'
import dayjs from 'dayjs'

const loading = ref(false)
const submitting = ref(false)
const list = ref<ResearchApplication[]>([])
const total = ref(0)
const pendingCount = ref(0)
const activeTab = ref<'pending' | 'approved' | 'rejected'>('pending')
const showDetail = ref(false)
const showApproveDialog = ref(false)
const showRejectDialog = ref(false)
const currentApplication = ref<ResearchApplication | null>(null)
const reviewFormRef = ref<FormInstance>()

const query = reactive({
  page: 1,
  page_size: 10,
  keyword: '',
  start_date: '',
  end_date: ''
})

const reviewForm = reactive({
  review_comment: ''
})

const reviewRules: FormRules = {
  review_comment: [{ max: 500, message: '审批意见不能超过500字', trigger: 'blur' }]
}

const rejectRules: FormRules = {
  review_comment: [
    { required: true, message: '请输入拒绝原因', trigger: 'blur' },
    { max: 500, message: '拒绝原因不能超过500字', trigger: 'blur' }
  ]
}

const statusType = (status: string) => {
  if (status === 'approved') return 'success'
  if (status === 'rejected') return 'danger'
  return 'warning'
}

const statusText = (status: string) => {
  const map: Record<string, string> = {
    pending: '待审批',
    approved: '已通过',
    rejected: '已拒绝'
  }
  return map[status] || status
}

const formatDateTime = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
}

const handleTabChange = () => {
  query.page = 1
  fetchList()
}

const fetchList = async () => {
  try {
    loading.value = true
    const params: any = {
      page: query.page,
      page_size: query.page_size,
      status: activeTab.value
    }
    if (query.start_date) params.start_date = query.start_date
    if (query.end_date) params.end_date = query.end_date
    const res = await guideApi.listResearchApplications(params)
    list.value = res.data.list
    total.value = res.data.total
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const fetchPendingCount = async () => {
  try {
    const res = await guideApi.listResearchApplications({ page: 1, page_size: 1, status: 'pending' })
    pendingCount.value = res.data.total
  } catch (e) {
    console.error(e)
  }
}

const resetQuery = () => {
  query.page = 1
  query.keyword = ''
  query.start_date = ''
  query.end_date = ''
  fetchList()
}

const handleViewDetail = async (row: ResearchApplication) => {
  try {
    loading.value = true
    const res = await guideApi.getResearchApplication(row.id)
    currentApplication.value = res.data
    showDetail.value = true
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

const handleApprove = (row: ResearchApplication) => {
  currentApplication.value = row
  reviewForm.review_comment = ''
  showDetail.value = false
  showApproveDialog.value = true
}

const handleReject = (row: ResearchApplication) => {
  currentApplication.value = row
  reviewForm.review_comment = ''
  showDetail.value = false
  showRejectDialog.value = true
}

const submitApprove = async () => {
  if (!currentApplication.value) return
  if (!reviewFormRef.value) return
  await reviewFormRef.value.validate()
  try {
    submitting.value = true
    await guideApi.reviewResearchApplication(currentApplication.value.id, {
      status: 'approved',
      review_comment: reviewForm.review_comment
    })
    ElMessage.success('审批通过')
    showApproveDialog.value = false
    fetchList()
    fetchPendingCount()
  } catch (e) {
    console.error(e)
  } finally {
    submitting.value = false
  }
}

const submitReject = async () => {
  if (!currentApplication.value) return
  if (!reviewFormRef.value) return
  await reviewFormRef.value.validate()
  try {
    submitting.value = true
    await guideApi.reviewResearchApplication(currentApplication.value.id, {
      status: 'rejected',
      review_comment: reviewForm.review_comment
    })
    ElMessage.success('已拒绝')
    showRejectDialog.value = false
    fetchList()
    fetchPendingCount()
  } catch (e) {
    console.error(e)
  } finally {
    submitting.value = false
  }
}

const exportData = () => {
  ElMessage.info('导出功能开发中...')
}

onMounted(() => {
  fetchList()
  fetchPendingCount()
})
</script>

<style scoped lang="scss">
.dashboard-research {
  .page-title {
    margin: 0;
    font-size: 20px;
  }

  .tab-label {
    display: inline-flex;
    align-items: center;

    .el-icon {
      margin-right: 5px;
    }
  }

  .tab-content {
    padding-top: 20px;
  }

  .text-xs {
    font-size: 12px;
  }

  .text-gray {
    color: #909399;
  }

  .text-sm {
    font-size: 13px;
  }

  .font-medium {
    font-weight: 500;
  }

  .application-detail {
    .review-comment {
      padding: 15px;
      background: #f5f7fa;
      border-radius: 6px;
      line-height: 1.6;
      white-space: pre-wrap;
    }

    h4 {
      margin: 0;
      font-size: 14px;
      color: #303133;
    }
  }

  .ml-10 {
    margin-left: 10px;
  }

  .mt-20 {
    margin-top: 20px;
  }

  .mb-10 {
    margin-bottom: 10px;
  }

  .mb-20 {
    margin-bottom: 20px;
  }
}
</style>
