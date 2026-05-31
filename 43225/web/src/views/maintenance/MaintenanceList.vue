<template>
  <div class="maintenance-list">
    <div class="page-header">
      <h2 class="page-title">维修保养</h2>
      <el-button type="primary" @click="showCreateDialog = true">
        <el-icon><Plus /></el-icon>
        添加记录
      </el-button>
    </div>

    <div class="filter-bar">
      <el-select v-model="filters.status" placeholder="状态" clearable style="width: 140px">
        <el-option label="已计划" value="scheduled" />
        <el-option label="进行中" value="in_progress" />
        <el-option label="已完成" value="completed" />
        <el-option label="已取消" value="cancelled" />
      </el-select>
      <el-button type="primary" @click="fetchMaintenances">
        <el-icon><Search /></el-icon>
        搜索
      </el-button>
    </div>

    <div class="card-container">
      <el-table :data="maintenances" v-loading="loading" style="width: 100%">
        <el-table-column prop="title" label="标题" min-width="150" />
        <el-table-column label="船只">
          <template #default="{ row }">
            {{ row.ship?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column prop="maintenance_type" label="类型" width="100">
          <template #default="{ row }">
            {{ getTypeText(row.maintenance_type) }}
          </template>
        </el-table-column>
        <el-table-column prop="planned_date" label="计划日期" width="120">
          <template #default="{ row }">
            {{ formatDate(row.planned_date) }}
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)" size="small">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="费用" width="120">
          <template #default="{ row }">
            {{ row.currency }} {{ row.cost }}
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link @click="handleUpdateStatus(row.id, 'in_progress')" v-if="row.status === 'scheduled'">
              开始
            </el-button>
            <el-button type="success" link @click="handleUpdateStatus(row.id, 'completed')" v-if="row.status === 'in_progress'">
              完成
            </el-button>
            <el-button type="primary" link @click="$router.push(`/maintenance/${row.id}`)">详情</el-button>
          </template>
        </el-table-column>
      </el-table>

      <el-empty v-if="maintenances.length === 0 && !loading" description="暂无保养记录" />
    </div>

    <el-dialog v-model="showCreateDialog" title="添加保养记录" width="600px">
      <el-form :model="form" label-width="100px">
        <el-form-item label="船只">
          <el-select v-model="form.ship_id" placeholder="请选择船只" style="width: 100%">
            <el-option v-for="ship in ships" :key="ship.id" :label="ship.name" :value="ship.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="标题">
          <el-input v-model="form.title" />
        </el-form-item>
        <el-form-item label="类型">
          <el-select v-model="form.maintenance_type" style="width: 100%">
            <el-option label="常规保养" value="routine" />
            <el-option label="维修" value="repair" />
            <el-option label="检查" value="inspection" />
            <el-option label="大修" value="overhaul" />
          </el-select>
        </el-form-item>
        <el-form-item label="计划日期">
          <el-date-picker v-model="form.planned_date" type="date" style="width: 100%" />
        </el-form-item>
        <el-form-item label="费用">
          <el-input-number v-model="form.cost" :min="0" :precision="2" style="width: 100%" />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="form.description" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" @click="handleCreate">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import { getMaintenancesApi, createMaintenanceApi, updateMaintenanceApi } from '@/api/maintenance'
import { getMyShipsApi } from '@/api/ship'
import type { MaintenanceRecord, CreateMaintenanceRequest } from '@/types/maintenance'
import type { Ship } from '@/types/ship'
import dayjs from 'dayjs'

const loading = ref(false)
const maintenances = ref<MaintenanceRecord[]>([])
const ships = ref<Ship[]>([])
const showCreateDialog = ref(false)

const filters = reactive({
  status: ''
})

const form = reactive<CreateMaintenanceRequest>({
  ship_id: '',
  title: '',
  maintenance_type: 'routine',
  planned_date: '',
  description: '',
  cost: 0
})

const formatDate = (date: string) => dayjs(date).format('YYYY-MM-DD')

const getTypeText = (type: string) => {
  const map: Record<string, string> = {
    routine: '常规保养',
    repair: '维修',
    inspection: '检查',
    overhaul: '大修'
  }
  return map[type] || type
}

const getStatusType = (status: string) => {
  const map: Record<string, string> = {
    scheduled: 'warning',
    in_progress: 'primary',
    completed: 'success',
    cancelled: 'info'
  }
  return map[status] || 'info'
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    scheduled: '已计划',
    in_progress: '进行中',
    completed: '已完成',
    cancelled: '已取消'
  }
  return map[status] || status
}

const fetchMaintenances = async () => {
  loading.value = true
  try {
    const params: any = {}
    if (filters.status) params.status = filters.status
    const res: any = await getMaintenancesApi(params)
    maintenances.value = res.data || []
  } catch (error) {
    console.error('Failed to fetch maintenances:', error)
  } finally {
    loading.value = false
  }
}

const fetchShips = async () => {
  try {
    const res: any = await getMyShipsApi()
    ships.value = res.data || []
  } catch (error) {
    console.error('Failed to fetch ships:', error)
  }
}

const handleCreate = async () => {
  try {
    await createMaintenanceApi(form)
    ElMessage.success('添加成功')
    showCreateDialog.value = false
    fetchMaintenances()
  } catch (error) {
    ElMessage.error('添加失败')
  }
}

const handleUpdateStatus = async (id: string, status: string) => {
  try {
    await updateMaintenanceApi(id, { status })
    ElMessage.success('状态更新成功')
    fetchMaintenances()
  } catch (error) {
    ElMessage.error('更新失败')
  }
}

onMounted(() => {
  fetchMaintenances()
  fetchShips()
})
</script>

<style lang="scss" scoped>
.maintenance-list {
  .card-container {
    background: #fff;
    border-radius: 8px;
    padding: 20px;
  }
}
</style>
