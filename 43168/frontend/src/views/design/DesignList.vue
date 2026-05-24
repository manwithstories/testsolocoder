<template>
  <div class="design-list">
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchForm" @submit.prevent>
        <el-form-item label="方案名称">
          <el-input v-model="searchForm.name" placeholder="请输入方案名称" clearable style="width: 200px" />
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
          <span>方案列表</span>
          <el-button type="primary" :icon="Plus" @click="handleCreate">新建方案</el-button>
        </div>
      </template>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column label="封面" width="120">
          <template #default="{ row }">
            <el-image
              v-if="row.coverImage"
              :src="row.coverImage"
              fit="cover"
              style="width: 80px; height: 60px; border-radius: 4px"
              :preview-src-list="[row.coverImage]"
            />
            <div v-else class="empty-cover">
              <el-icon :size="24"><Picture /></el-icon>
            </div>
          </template>
        </el-table-column>
        <el-table-column prop="name" label="方案名称" min-width="160" show-overflow-tooltip />
        <el-table-column prop="ownerName" label="业主" width="120" show-overflow-tooltip />
        <el-table-column prop="houseType" label="房型" width="100" show-overflow-tooltip />
        <el-table-column label="面积" width="100">
          <template #default="{ row }">{{ row.area ? `${row.area} ㎡` : '-' }}</template>
        </el-table-column>
        <el-table-column label="预算" width="120">
          <template #default="{ row }">{{ row.budget ? `¥${row.budget.toLocaleString()}` : '-' }}</template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="170" />
        <el-table-column label="操作" width="260" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleDetail(row)">详情</el-button>
            <el-button link type="primary" @click="handleEdit(row)">编辑</el-button>
            <el-popconfirm title="确认删除该方案？" @confirm="handleDelete(row)">
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
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Search, Refresh, Plus, Picture } from '@element-plus/icons-vue'
import { listDesigns, deleteDesign, type DesignProject, type ListDesignParams } from '@/api/design'

const router = useRouter()
const loading = ref(false)
const list = ref<DesignProject[]>([])

const searchForm = reactive<ListDesignParams>({
  page: 1,
  pageSize: 10,
  name: '',
  status: undefined as unknown as number
})

const pagination = reactive({ page: 1, pageSize: 10, total: 0 })

const statusOptions = [
  { label: '草稿', value: 0 },
  { label: '进行中', value: 1 },
  { label: '待审核', value: 2 },
  { label: '已通过', value: 3 },
  { label: '已驳回', value: 4 },
  { label: '已归档', value: 5 }
]

function getStatusLabel(status: number) {
  return statusOptions.find((i) => i.value === status)?.label ?? '未知'
}
function getStatusType(status: number) {
  const map: Record<number, string> = {
    0: 'info',
    1: 'primary',
    2: 'warning',
    3: 'success',
    4: 'danger',
    5: ''
  }
  return map[status] ?? ''
}

async function fetchList() {
  loading.value = true
  try {
    const res = await listDesigns({
      name: searchForm.name || undefined,
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
  searchForm.name = ''
  searchForm.status = undefined as unknown as number
  pagination.page = 1
  fetchList()
}

function handleCreate() {
  router.push('/designs/new')
}

function handleEdit(row: DesignProject) {
  router.push(`/designs/${row.id}/edit`)
}

function handleDetail(row: DesignProject) {
  router.push(`/designs/${row.id}`)
}

async function handleDelete(row: DesignProject) {
  await deleteDesign(row.id)
  ElMessage.success('删除成功')
  fetchList()
}

onMounted(fetchList)
</script>

<style lang="scss" scoped>
.design-list {
  .search-card {
    margin-bottom: 16px;
  }
  .card-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .empty-cover {
    width: 80px;
    height: 60px;
    background: #f5f7fa;
    display: flex;
    align-items: center;
    justify-content: center;
    color: #c0c4cc;
    border-radius: 4px;
  }
  .pagination-wrapper {
    display: flex;
    justify-content: flex-end;
    padding-top: 16px;
  }
}
</style>
