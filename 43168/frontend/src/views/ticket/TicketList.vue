<template>
  <div class="ticket-list">
    <el-card shadow="never" class="search-card">
      <el-form :inline="true" :model="searchForm" @submit.prevent>
        <el-form-item label="工单标题">
          <el-input v-model="searchForm.title" placeholder="请输入工单标题" clearable style="width: 200px" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="searchForm.type" placeholder="请选择类型" clearable style="width: 200px">
            <el-option v-for="item in typeOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="状态">
          <el-select v-model="searchForm.status" placeholder="请选择状态" clearable style="width: 200px">
            <el-option v-for="item in statusOptions" :key="item.value" :label="item.label" :value="item.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="优先级">
          <el-select v-model="searchForm.priority" placeholder="请选择优先级" clearable style="width: 200px">
            <el-option v-for="item in priorityOptions" :key="item.value" :label="item.label" :value="item.value" />
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
          <span>工单列表</span>
          <el-button type="primary" :icon="Plus" @click="handleCreate">创建工单</el-button>
        </div>
      </template>

      <el-table :data="list" v-loading="loading" stripe>
        <el-table-column prop="id" label="工单号" width="100" />
        <el-table-column prop="title" label="标题" min-width="160" show-overflow-tooltip />
        <el-table-column label="类型" width="100">
          <template #default="{ row }">{{ getTypeLabel(row.type) }}</template>
        </el-table-column>
        <el-table-column label="优先级" width="100">
          <template #default="{ row }">
            <el-tag :type="getPriorityType(row.priority)">{{ getPriorityLabel(row.priority) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="ownerName" label="提交人" width="100" />
        <el-table-column prop="handlerName" label="处理人" width="100" />
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="createdAt" label="创建时间" width="170" />
        <el-table-column label="操作" width="240" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="handleView(row)">查看</el-button>
            <el-button link type="primary" @click="handleClose(row)" :disabled="row.status === 3">关闭</el-button>
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

    <el-dialog v-model="detailVisible" title="工单详情" width="600px">
      <el-descriptions v-if="currentTicket" :column="2" border>
        <el-descriptions-item label="工单号">{{ currentTicket.id }}</el-descriptions-item>
        <el-descriptions-item label="标题">{{ currentTicket.title }}</el-descriptions-item>
        <el-descriptions-item label="类型">{{ getTypeLabel(currentTicket.type) }}</el-descriptions-item>
        <el-descriptions-item label="优先级">{{ getPriorityLabel(currentTicket.priority) }}</el-descriptions-item>
        <el-descriptions-item label="提交人">{{ currentTicket.ownerName || '-' }}</el-descriptions-item>
        <el-descriptions-item label="处理人">{{ currentTicket.handlerName || '-' }}</el-descriptions-item>
        <el-descriptions-item label="状态">
          <el-tag :type="getStatusType(currentTicket.status)">{{ getStatusLabel(currentTicket.status) }}</el-tag>
        </el-descriptions-item>
        <el-descriptions-item label="创建时间">{{ currentTicket.createdAt }}</el-descriptions-item>
        <el-descriptions-item label="描述" :span="2">{{ currentTicket.description }}</el-descriptions-item>
        <el-descriptions-item v-if="currentTicket.reply" label="处理回复" :span="2">{{ currentTicket.reply }}</el-descriptions-item>
      </el-descriptions>
      <template #footer>
        <el-button @click="detailVisible = false">关闭</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="closeVisible" title="关闭工单" width="520px">
      <el-input v-model="closeReply" type="textarea" :rows="4" placeholder="请输入处理回复" maxlength="200" show-word-limit />
      <template #footer>
        <el-button @click="closeVisible = false">取消</el-button>
        <el-button type="primary" :loading="closing" @click="submitClose">提交</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { reactive, ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage } from 'element-plus'
import { Search, Refresh, Plus } from '@element-plus/icons-vue'
import { listTickets, deleteTicket, closeTicket, type Ticket, type ListTicketParams } from '@/api/ticket'

const router = useRouter()
const loading = ref(false)
const list = ref<Ticket[]>([])
const detailVisible = ref(false)
const closeVisible = ref(false)
const closing = ref(false)
const currentTicket = ref<Ticket | null>(null)
const closeReply = ref('')

const searchForm = reactive<ListTicketParams>({
  page: 1,
  pageSize: 10,
  title: '',
  type: undefined as unknown as number,
  status: undefined as unknown as number,
  priority: undefined as unknown as number
})

const pagination = reactive({ page: 1, pageSize: 10, total: 0 })

const typeOptions = [
  { label: '产品问题', value: 1 },
  { label: '安装问题', value: 2 },
  { label: '售后服务', value: 3 },
  { label: '投诉建议', value: 4 },
  { label: '其他', value: 99 }
]

const priorityOptions = [
  { label: '低', value: 1 },
  { label: '中', value: 2 },
  { label: '高', value: 3 },
  { label: '紧急', value: 4 }
]

const statusOptions = [
  { label: '待处理', value: 0 },
  { label: '处理中', value: 1 },
  { label: '已解决', value: 2 },
  { label: '已关闭', value: 3 }
]

function getTypeLabel(type: number) {
  return typeOptions.find((i) => i.value === type)?.label ?? '未知'
}
function getPriorityLabel(priority: number) {
  return priorityOptions.find((i) => i.value === priority)?.label ?? '未知'
}
function getPriorityType(priority: number) {
  const map: Record<number, string> = { 1: 'info', 2: '', 3: 'warning', 4: 'danger' }
  return map[priority] ?? ''
}
function getStatusLabel(status: number) {
  return statusOptions.find((i) => i.value === status)?.label ?? '未知'
}
function getStatusType(status: number) {
  const map: Record<number, string> = { 0: 'warning', 1: 'primary', 2: 'success', 3: 'info' }
  return map[status] ?? ''
}

async function fetchList() {
  loading.value = true
  try {
    const res = await listTickets({
      title: searchForm.title || undefined,
      type: searchForm.type ?? undefined,
      status: searchForm.status ?? undefined,
      priority: searchForm.priority ?? undefined,
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
  searchForm.title = ''
  searchForm.type = undefined as unknown as number
  searchForm.status = undefined as unknown as number
  searchForm.priority = undefined as unknown as number
  pagination.page = 1
  fetchList()
}

function handleCreate() {
  router.push('/tickets/new')
}

function handleView(row: Ticket) {
  currentTicket.value = row
  detailVisible.value = true
}

function handleClose(row: Ticket) {
  currentTicket.value = row
  closeReply.value = ''
  closeVisible.value = true
}

async function submitClose() {
  if (!currentTicket.value) return
  closing.value = true
  try {
    await closeTicket(currentTicket.value.id, closeReply.value)
    ElMessage.success('关闭成功')
    closeVisible.value = false
    fetchList()
  } finally {
    closing.value = false
  }
}

async function handleDelete(row: Ticket) {
  await deleteTicket(row.id)
  ElMessage.success('删除成功')
  fetchList()
}

onMounted(fetchList)
</script>

<style lang="scss" scoped>
.ticket-list {
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
