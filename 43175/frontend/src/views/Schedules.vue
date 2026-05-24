<template>
  <div class="page-container">
    <div class="page-header">
      <h2>定时任务</h2>
      <el-button type="primary" :icon="Plus" @click="showCreateDialog = true">
        新建任务
      </el-button>
    </div>

    <div class="table-container">
      <el-table :data="schedules" v-loading="loading" stripe>
        <el-table-column prop="name" label="任务名称" min-width="150" />
        <el-table-column prop="cronExpr" label="执行时间" width="180">
          <template #default="{ row }">
            <el-tag size="small">{{ row.cronExpr }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="action" label="动作" width="100">
          <template #default="{ row }">
            {{ getActionLabel(row.action) }}
          </template>
        </el-table-column>
        <el-table-column prop="isEnabled" label="状态" width="100">
          <template #default="{ row }">
            <el-switch v-model="row.isEnabled" @change="toggleSchedule(row)" />
          </template>
        </el-table-column>
        <el-table-column prop="lastRun" label="上次执行" width="180">
          <template #default="{ row }">{{ row.lastRun || '-' }}</template>
        </el-table-column>
        <el-table-column label="操作" width="150" fixed="right">
          <template #default="{ row }">
            <el-button size="small" type="primary" link @click="viewLogs(row)">日志</el-button>
            <el-button size="small" @click="editSchedule(row)">编辑</el-button>
            <el-button size="small" type="danger" link @click="handleDeleteSchedule(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="schedules.length === 0 && !loading" description="暂无定时任务" />
    </div>

    <el-dialog v-model="showCreateDialog" :title="editingSchedule ? '编辑任务' : '新建任务'" width="550px">
      <el-form :model="scheduleForm" ref="scheduleFormRef" label-width="80px">
        <el-form-item label="任务名称">
          <el-input v-model="scheduleForm.name" />
        </el-form-item>
        <el-form-item label="目标设备">
          <el-select v-model="scheduleForm.deviceId" style="width: 100%">
            <el-option v-for="d in devices" :key="d.id" :label="d.name" :value="d.id" />
          </el-select>
        </el-form-item>
        <el-form-item label="执行动作">
          <el-select v-model="scheduleForm.action" style="width: 100%">
            <el-option v-for="opt in scheduleActionOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="执行时间">
          <el-select v-model="scheduleForm.cronExpr" style="width: 100%">
            <el-option v-for="opt in cronPresets" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="scheduleForm.description" type="textarea" :rows="2" />
        </el-form-item>
        <el-form-item label="启用">
          <el-switch v-model="scheduleForm.isEnabled" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveSchedule">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showLogsDialog" title="执行日志" width="700px">
      <el-table :data="scheduleLogs" style="width: 100%" max-height="400">
        <el-table-column prop="executedAt" label="执行时间" width="180" />
        <el-table-column prop="action" label="动作" width="100" />
        <el-table-column prop="success" label="结果" width="80">
          <template #default="{ row }">
            <el-tag :type="row.success ? 'success' : 'danger'" size="small">
              {{ row.success ? '成功' : '失败' }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="message" label="消息" show-overflow-tooltip />
        <el-table-column prop="energyDelta" label="能耗变化(kWh)" width="140">
          <template #default="{ row }">{{ row.energyDelta?.toFixed(2) || '0.00' }}</template>
        </el-table-column>
      </el-table>
      <el-empty v-if="scheduleLogs.length === 0" description="暂无日志" />
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  listSchedules, createSchedule, updateSchedule, deleteSchedule,
  listScheduleLogs, scheduleActionOptions, cronPresets, type Schedule
} from '@/api/schedule'
import { listDevices } from '@/api/device'
import { useFamilyStore } from '@/stores/family'

const familyStore = useFamilyStore()
const loading = ref(false)
const saving = ref(false)
const schedules = ref<Schedule[]>([])
const devices = ref<any[]>([])
const showCreateDialog = ref(false)
const showLogsDialog = ref(false)
const editingSchedule = ref<Schedule | null>(null)
const scheduleLogs = ref<any[]>([])
const scheduleFormRef = ref<FormInstance>()

const scheduleForm = reactive({
  name: '',
  deviceId: null as number | null,
  action: 'on',
  cronExpr: '0 8 * * *',
  description: '',
  isEnabled: true
})

onMounted(async () => {
  await loadSchedules()
  devices.value = await listDevices({ familyId: familyStore.currentFamilyId || undefined })
})

watch(() => familyStore.currentFamilyId, async () => {
  loadSchedules()
  devices.value = await listDevices({ familyId: familyStore.currentFamilyId || undefined })
})

async function loadSchedules() {
  loading.value = true
  try {
    schedules.value = await listSchedules({ familyId: familyStore.currentFamilyId || undefined })
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function editSchedule(schedule: Schedule) {
  editingSchedule.value = schedule
  Object.assign(scheduleForm, {
    name: schedule.name,
    deviceId: schedule.deviceId,
    action: schedule.action,
    cronExpr: schedule.cronExpr,
    description: schedule.description,
    isEnabled: schedule.isEnabled
  })
  showCreateDialog.value = true
}

async function saveSchedule() {
  saving.value = true
  try {
    if (!familyStore.currentFamilyId) {
      ElMessage.warning('请先选择或创建家庭')
      return
    }
    const data = { ...scheduleForm, familyId: familyStore.currentFamilyId }
    if (editingSchedule.value) {
      await updateSchedule(editingSchedule.value.id, data)
      ElMessage.success('更新成功')
    } else {
      await createSchedule(data)
      ElMessage.success('创建成功')
    }
    showCreateDialog.value = false
    resetForm()
    loadSchedules()
  } catch (e) {
    console.error(e)
  } finally {
    saving.value = false
  }
}

async function handleDeleteSchedule(schedule: Schedule) {
  try {
    await ElMessageBox.confirm(`确定要删除任务"${schedule.name}"吗？`, '提示', { type: 'warning' })
    await deleteSchedule(schedule.id)
    ElMessage.success('删除成功')
    loadSchedules()
  } catch (e: any) {
    if (e !== 'cancel') console.error(e)
  }
}

async function toggleSchedule(schedule: Schedule) {
  try {
    await updateSchedule(schedule.id, { isEnabled: schedule.isEnabled })
    ElMessage.success(schedule.isEnabled ? '已启用' : '已禁用')
  } catch (e) {
    console.error(e)
  }
}

async function viewLogs(schedule: Schedule) {
  try {
    scheduleLogs.value = await listScheduleLogs(schedule.id)
    showLogsDialog.value = true
  } catch (e) {
    console.error(e)
  }
}

function resetForm() {
  editingSchedule.value = null
  Object.assign(scheduleForm, {
    name: '',
    deviceId: null,
    action: 'on',
    cronExpr: '0 8 * * *',
    description: '',
    isEnabled: true
  })
}

function getActionLabel(action: string) {
  return scheduleActionOptions.find(o => o.value === action)?.label || action
}
</script>
