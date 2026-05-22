<template>
  <div class="repair-list">
    <div class="page-header">
      <h2 class="page-title">维修工单</h2>
      <div>
        <el-button type="primary" @click="goToCreate">
          <el-icon><Plus /></el-icon>提交报修
        </el-button>
        <el-button type="success" @click="handleExport">
          <el-icon><Download /></el-icon>导出Excel
        </el-button>
      </div>
    </div>

    <div class="search-bar card">
      <el-select v-model="statusFilter" placeholder="状态" clearable style="width: 150px">
        <el-option label="待处理" :value="1" />
        <el-option label="处理中" :value="2" />
        <el-option label="已完成" :value="3" />
        <el-option label="已关闭" :value="4" />
      </el-select>
      <el-select v-model="categoryFilter" placeholder="分类" clearable style="width: 150px">
        <el-option label="水电维修" value="水电维修" />
        <el-option label="家电维修" value="家电维修" />
        <el-option label="家具维修" value="家具维修" />
        <el-option label="门窗维修" value="门窗维修" />
        <el-option label="其他" value="其他" />
      </el-select>
      <el-button type="primary" @click="loadData">搜索</el-button>
    </div>

    <div class="card">
      <el-table :data="orders" v-loading="loading">
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="title" label="标题" min-width="150" />
        <el-table-column prop="category" label="分类" width="100" />
        <el-table-column label="房源" min-width="150">
          <template #default="{ row }">
            {{ row.property?.title }}
          </template>
        </el-table-column>
        <el-table-column label="租户" width="120">
          <template #default="{ row }">
            {{ row.tenant?.name }}
          </template>
        </el-table-column>
        <el-table-column label="处理人" width="120">
          <template #default="{ row }">
            {{ row.handler?.realName || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="priority" label="优先级" width="100">
          <template #default="{ row }">
            <el-tag :type="getPriorityType(row.priority)">{{ getPriorityText(row.priority) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button link type="primary" @click="viewDetail(row.id)">详情</el-button>
            <el-button
              link
              type="success"
              @click="handleAssign(row)"
              v-if="row.status === 1"
            >派单</el-button>
            <el-button
              link
              type="warning"
              @click="handleComplete(row)"
              v-if="row.status === 2"
            >完成</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-pagination
        class="pagination"
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        :page-sizes="[10, 20, 50]"
        layout="total, sizes, prev, pager, next, jumper"
        @current-change="loadData"
        @size-change="loadData"
      />
    </div>

    <el-dialog v-model="assignDialogVisible" title="派单" width="400px">
      <el-form label-width="80px">
        <el-form-item label="处理人">
          <el-select v-model="selectedHandlerId" placeholder="请选择处理人" style="width: 100%">
            <el-option
              v-for="u in handlers"
              :key="u.id"
              :label="u.realName"
              :value="u.id"
            />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="assignDialogVisible = false">取消</el-button>
        <el-button type="primary" @click="confirmAssign">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import type { RepairOrder, User } from '@/types'
import { getRepairs, assignRepair, updateRepairStatus } from '@/api/business'
import { getUsers } from '@/api/auth'
import { exportRepairOrders } from '@/api/business'

const router = useRouter()
const loading = ref(false)
const orders = ref<RepairOrder[]>([])
const statusFilter = ref<number | ''>('')
const categoryFilter = ref('')
const assignDialogVisible = ref(false)
const assignOrderId = ref<number | null>(null)
const selectedHandlerId = ref(0)
const handlers = ref<User[]>([])

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

onMounted(() => {
  loadData()
  loadHandlers()
})

async function loadData() {
  loading.value = true
  try {
    const res = await getRepairs({
      page: pagination.page,
      pageSize: pagination.pageSize,
      status: statusFilter.value || undefined,
      category: categoryFilter.value || undefined
    })
    orders.value = res.data.list
    pagination.total = res.data.total
  } catch (error) {
    console.error('Failed to load repairs:', error)
  } finally {
    loading.value = false
  }
}

async function loadHandlers() {
  try {
    const res = await getUsers({ role: 'maintainer', pageSize: 100 })
    handlers.value = res.data.list
  } catch (error) {
    console.error('Failed to load handlers:', error)
  }
}

function getStatusType(status: number) {
  switch (status) {
    case 1: return 'warning'
    case 2: return 'primary'
    case 3: return 'success'
    default: return 'info'
  }
}

function getStatusText(status: number) {
  switch (status) {
    case 1: return '待处理'
    case 2: return '处理中'
    case 3: return '已完成'
    default: return '已关闭'
  }
}

function getPriorityType(priority: number) {
  switch (priority) {
    case 1: return 'danger'
    case 2: return 'warning'
    default: return 'info'
  }
}

function getPriorityText(priority: number) {
  switch (priority) {
    case 1: return '紧急'
    case 2: return '高'
    default: return '普通'
  }
}

function goToCreate() {
  router.push('/repairs/create')
}

function viewDetail(id: number) {
  router.push(`/repairs/${id}`)
}

function handleAssign(row: RepairOrder) {
  assignOrderId.value = row.id
  selectedHandlerId.value = 0
  assignDialogVisible.value = true
}

async function confirmAssign() {
  if (!assignOrderId.value || !selectedHandlerId.value) {
    ElMessage.warning('请选择处理人')
    return
  }
  try {
    await assignRepair(assignOrderId.value, selectedHandlerId.value)
    ElMessage.success('派单成功')
    assignDialogVisible.value = false
    loadData()
  } catch (error) {
    console.error(error)
  }
}

async function handleComplete(row: RepairOrder) {
  try {
    await ElMessageBox.confirm('确定要标记为已完成吗？', '提示', { type: 'warning' })
    await updateRepairStatus(row.id, { status: 3 })
    ElMessage.success('操作成功')
    loadData()
  } catch (error) {
    if (error !== 'cancel') {
      console.error(error)
    }
  }
}

async function handleExport() {
  try {
    const res = await exportRepairOrders()
    const blob = new Blob([res as any], { type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet' })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = 'repair_orders.xlsx'
    link.click()
    URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    console.error(error)
  }
}
</script>

<style scoped>
.repair-list {
  padding: 0;
}

.pagination {
  margin-top: 20px;
  display: flex;
  justify-content: flex-end;
}
</style>
