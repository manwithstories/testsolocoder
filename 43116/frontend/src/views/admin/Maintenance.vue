<template>
  <div>
    <div class="page-header">
      <h1 class="page-title">维护管理</h1>
      <el-button type="primary" :icon="Plus" @click="showAddDialog = true">添加维护计划</el-button>
    </div>

    <div class="search-bar">
      <el-select v-model="filters.status" placeholder="状态" clearable style="width: 140px">
        <el-option label="已计划" value="scheduled" />
        <el-option label="进行中" value="in_progress" />
        <el-option label="已完成" value="completed" />
        <el-option label="已取消" value="cancelled" />
      </el-select>
      <el-button type="primary" @click="loadMaintenance">搜索</el-button>
    </div>

    <el-table :data="maintenanceList" v-loading="loading" style="width: 100%">
      <el-table-column prop="title" label="标题" min-width="200" />
      <el-table-column label="车辆" min-width="150">
        <template #default="{ row }">
          {{ row.car?.brand }} {{ row.car?.model }}
        </template>
      </el-table-column>
      <el-table-column label="开始时间" width="170">
        <template #default="{ row }">
          {{ formatDate(row.start_date) }}
        </template>
      </el-table-column>
      <el-table-column label="结束时间" width="170">
        <template #default="{ row }">
          {{ formatDate(row.end_date) }}
        </template>
      </el-table-column>
      <el-table-column label="费用" width="120">
        <template #default="{ row }">
          ¥{{ row.cost.toFixed(2) }}
        </template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <el-tag :type="getStatusType(row.status)">{{ getStatusText(row.status) }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="200" fixed="right">
        <template #default="{ row }">
          <el-button
            v-if="row.status === 'scheduled'"
            type="primary"
            link
            size="small"
            @click="startMaintenance(row)"
          >
            开始
          </el-button>
          <el-button
            v-if="row.status === 'in_progress'"
            type="success"
            link
            size="small"
            @click="completeMaintenance(row)"
          >
            完成
          </el-button>
          <el-button
            v-if="row.status === 'scheduled' || row.status === 'in_progress'"
            type="warning"
            link
            size="small"
            @click="cancelMaintenance(row)"
          >
            取消
          </el-button>
          <el-button type="danger" link size="small" @click="deleteMaintenance(row)">删除</el-button>
        </template>
      </el-table-column>
    </el-table>

    <div style="text-align: center; margin-top: 20px;">
      <el-pagination
        v-model:current-page="pagination.page"
        v-model:page-size="pagination.pageSize"
        :total="pagination.total"
        layout="total, prev, pager, next"
        @current-change="loadMaintenance"
      />
    </div>

    <el-dialog v-model="showAddDialog" title="添加维护计划" width="500px">
      <el-form :model="maintenanceForm" label-width="100px">
        <el-form-item label="标题">
          <el-input v-model="maintenanceForm.title" />
        </el-form-item>
        <el-form-item label="车辆">
          <el-select v-model="maintenanceForm.car_id" filterable>
            <el-option
              v-for="car in cars"
              :key="car.id"
              :label="`${car.brand} ${car.model} (${car.license_plate})`"
              :value="car.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="开始时间">
          <el-date-picker
            v-model="maintenanceForm.start_date"
            type="datetime"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="结束时间">
          <el-date-picker
            v-model="maintenanceForm.end_date"
            type="datetime"
            style="width: 100%"
          />
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="maintenanceForm.description" type="textarea" :rows="3" />
        </el-form-item>
        <el-form-item label="费用">
          <el-input-number v-model="maintenanceForm.cost" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="maintenanceForm.notes" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showAddDialog = false">取消</el-button>
        <el-button type="primary" @click="saveMaintenance">保存</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import { maintenanceApi, carApi } from '@/api'
import type { MaintenancePlan, Car } from '@/types'

const maintenanceList = ref<MaintenancePlan[]>([])
const cars = ref<Car[]>([])
const loading = ref(false)
const showAddDialog = ref(false)

const filters = reactive({
  status: ''
})

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const maintenanceForm = reactive({
  title: '',
  car_id: undefined as number | undefined,
  start_date: undefined as Date | undefined,
  end_date: undefined as Date | undefined,
  description: '',
  cost: 0,
  notes: ''
})

onMounted(() => {
  loadMaintenance()
  loadCars()
})

const loadMaintenance = async () => {
  loading.value = true
  try {
    const res = await maintenanceApi.getMaintenance({
      page: pagination.page,
      page_size: pagination.pageSize,
      status: filters.status
    })
    maintenanceList.value = res.data.items
    pagination.total = res.data.total
  } catch {
    // ignore
  } finally {
    loading.value = false
  }
}

const loadCars = async () => {
  try {
    const res = await carApi.getCars({ page: 1, page_size: 100 })
    cars.value = res.data.items
  } catch {
    // ignore
  }
}

const saveMaintenance = async () => {
  try {
    await maintenanceApi.createMaintenance({
      ...maintenanceForm,
      start_date: maintenanceForm.start_date?.toISOString(),
      end_date: maintenanceForm.end_date?.toISOString()
    })
    ElMessage.success('添加成功')
    showAddDialog.value = false
    loadMaintenance()
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  }
}

const startMaintenance = async (row: MaintenancePlan) => {
  try {
    await maintenanceApi.startMaintenance(row.id)
    ElMessage.success('已开始维护')
    loadMaintenance()
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  }
}

const completeMaintenance = async (row: MaintenancePlan) => {
  try {
    await maintenanceApi.completeMaintenance(row.id)
    ElMessage.success('维护已完成')
    loadMaintenance()
  } catch (err: any) {
    ElMessage.error(err.message || '操作失败')
  }
}

const cancelMaintenance = async (row: MaintenancePlan) => {
  try {
    await ElMessageBox.confirm('确定要取消该维护计划吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await maintenanceApi.cancelMaintenance(row.id)
    ElMessage.success('取消成功')
    loadMaintenance()
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error(err.message || '操作失败')
    }
  }
}

const deleteMaintenance = async (row: MaintenancePlan) => {
  try {
    await ElMessageBox.confirm('确定要删除该维护计划吗？', '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })

    await maintenanceApi.deleteMaintenance(row.id)
    ElMessage.success('删除成功')
    loadMaintenance()
  } catch (err: any) {
    if (err !== 'cancel') {
      ElMessage.error(err.message || '删除失败')
    }
  }
}

const formatDate = (date: string) => {
  return dayjs(date).format('YYYY-MM-DD HH:mm')
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
</script>
