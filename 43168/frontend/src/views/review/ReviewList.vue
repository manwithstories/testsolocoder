<template>
  <div class="review-list">
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchForm" @submit.prevent>
        <el-form-item label="订单号">
          <el-input v-model="searchForm.orderNo" placeholder="请输入订单号" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable style="width: 200px">
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item>
          <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-card shadow="never">
      <template #header>
        <div class="card-header">
          <span>评价列表</span>
          <el-button type="primary" :icon="Plus" @click="handleCreate">新增评价</el-button>
        </div>
      </template>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="orderNo" label="订单号" width="180" />
        <el-table-column prop="productName" label="产品" min-width="160" show-overflow-tooltip />
        <el-table-column label="产品评分" width="120">
          <template #default="{ row }">
            <el-rate v-model="row.productScore" disabled show-score text-color="#ff9900" />
          </template>
        </el-table-column>
        <el-table-column label="服务评分" width="120">
          <template #default="{ row }">
            <el-rate v-model="row.serviceScore" disabled show-score text-color="#ff9900" />
          </template>
        </el-table-column>
        <el-table-column prop="content" label="评价内容" min-width="200" show-overflow-tooltip />
        <el-table-column prop="ownerName" label="评价人" width="100" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="评价时间" width="170" />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleReply(row)">回复</el-button>
            <el-popconfirm title="确认删除？" @confirm="handleDelete(row)">
              <template #reference>
                <el-button link type="danger">删除</el-button>
              </template>
            </el-popconfirm>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-wrapper">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          background
          @size-change="fetchList"
          @current-change="fetchList"
        />
      </div>
    </el-card>

    <el-dialog v-model="replyVisible" title="回复评价" width="520px">
      <el-input v-model="replyContent" type="textarea" :rows="4" placeholder="请输入回复内容" maxlength="200" show-word-limit />
      <template #footer>
        <el-button @click="replyVisible = false">取消</el-button>
        <el-button type="primary" :loading="replying" @click="submitReply">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Search, Refresh, Plus } from '@element-plus/icons-vue'
import { listReviews, deleteReview, replyReview, type Review, type ListReviewParams } from '@/api/review'

const router = useRouter()
const loading = ref(false)
const list = ref<Review[]>([])
const replyVisible = ref(false)
const replying = ref(false)
const currentReview = ref<Review | null>(null)
const replyContent = ref('')

const searchForm = reactive<ListReviewParams>({
  page: 1,
  pageSize: 10,
  orderNo: '',
  status: undefined as unknown as number
})

const pagination = reactive({ page: 1, pageSize: 10, total: 0 })

const statusOptions = [
  { label: '待审核', value: 0 },
  { label: '已发布', value: 1 },
  { label: '已隐藏', value: 2 }
]

function getStatusLabel(status: number) {
  return statusOptions.find((i) => i.value === status)?.label ?? '未知'
}
function getStatusType(status: number) {
  const map: Record<number, string> = { 0: 'warning', 1: 'success', 2: 'info' }
  return map[status] ?? ''
}

async function fetchList() {
  loading.value = true
  try {
    const res = await listReviews({
      orderNo: searchForm.orderNo || undefined,
      status: searchForm.status ?? undefined,
      page: pagination.page,
      pageSize: pagination.pageSize
    })
    list.value = (res as any)?.list ?? []
    pagination.total = (res as any)?.total ?? 0
  } finally {
    loading.value = false
  }
}

function handleSearch() {
  pagination.page = 1
  fetchList()
}

function handleReset() {
  searchForm.orderNo = ''
  searchForm.status = undefined as unknown as number
  pagination.page = 1
  fetchList()
}

function handleCreate() {
  router.push('/reviews/new')
}

async function handleDelete(row: Review) {
  await deleteReview(row.id)
  ElMessage.success('删除成功')
  fetchList()
}

function handleReply(row: Review) {
  currentReview.value = row
  replyContent.value = ''
  replyVisible.value = true
}

async function submitReply() {
  if (!replyContent.value.trim()) {
    ElMessage.warning('请输入回复内容')
    return
  }
  if (!currentReview.value) return
  replying.value = true
  try {
    await replyReview(currentReview.value.id, replyContent.value)
    ElMessage.success('回复成功')
    replyVisible.value = false
    fetchList()
  } finally {
    replying.value = false
  }
}

onMounted(fetchList)
</script>

<style lang="scss" scoped>
.review-list {
  .search-card {
    margin-bottom: 16px;
  }
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    padding-top: 16px;
  }
}
</style>
