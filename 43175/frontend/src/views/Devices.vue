<template>
  <div class="page-container">
    <div class="page-header">
      <h2>设备管理</h2>
      <el-button type="primary" :icon="Plus" @click="showCreateDialog = true">
        添加设备
      </el-button>
    </div>

    <div class="search-bar">
      <el-input v-model="searchForm.location" placeholder="搜索位置" style="width: 200px" clearable @input="loadDevices" />
      <el-select v-model="searchForm.deviceType" placeholder="设备类型" style="width: 160px" clearable @change="loadDevices">
        <el-option v-for="opt in deviceTypeOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
      </el-select>
      <el-select v-model="searchForm.status" placeholder="设备状态" style="width: 140px" clearable @change="loadDevices">
        <el-option label="在线" value="online" />
        <el-option label="离线" value="offline" />
      </el-select>
    </div>

    <div class="table-container">
      <el-table :data="devices" v-loading="loading" stripe>
        <el-table-column prop="name" label="设备名称" min-width="120">
          <template #default="{ row }">
            <el-link type="primary" @click="$router.push(`/devices/${row.id}`)">{{ row.name }}</el-link>
          </template>
        </el-table-column>
        <el-table-column prop="deviceType" label="类型" width="100">
          <template #default="{ row }">
            {{ getDeviceTypeLabel(row.deviceType) }}
          </template>
        </el-table-column>
        <el-table-column prop="location" label="位置" width="120" />
        <el-table-column prop="power" label="功率(W)" width="100" />
        <el-table-column prop="protocol" label="协议" width="100">
          <template #default="{ row }">
            {{ getProtocolLabel(row.protocol) }}
          </template>
        </el-table-column>
        <el-table-column prop="status" label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">{{ getStatusLabel(row.status) }}</el-tag>
          </template>
        </el-table-column>
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button size="small" :type="row.status === 'on' || row.status === 'online' ? 'warning' : 'success'"
              @click="toggleStatus(row)">
              {{ row.status === 'on' || row.status === 'online' ? '关闭' : '开启' }}
            </el-button>
            <el-button size="small" type="primary" link @click="editDevice(row)">编辑</el-button>
            <el-button size="small" type="danger" link @click="handleDeleteDevice(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>
      <el-empty v-if="devices.length === 0 && !loading" description="暂无设备" />
    </div>

    <el-dialog v-model="showCreateDialog" :title="editingDevice ? '编辑设备' : '添加设备'" width="500px">
      <el-form :model="deviceForm" :rules="deviceRules" ref="deviceFormRef" label-width="80px">
        <el-form-item label="设备名称" prop="name">
          <el-input v-model="deviceForm.name" />
        </el-form-item>
        <el-form-item label="设备类型" prop="deviceType">
          <el-select v-model="deviceForm.deviceType" style="width: 100%">
            <el-option v-for="opt in deviceTypeOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="厂商" prop="vendor">
          <el-input v-model="deviceForm.vendor" />
        </el-form-item>
        <el-form-item label="安装位置" prop="location">
          <el-input v-model="deviceForm.location" placeholder="如：客厅、卧室" />
        </el-form-item>
        <el-form-item label="功率(W)" prop="power">
          <el-input-number v-model="deviceForm.power" :min="0" :max="10000" />
        </el-form-item>
        <el-form-item label="通信协议" prop="protocol">
          <el-select v-model="deviceForm.protocol" style="width: 100%">
            <el-option v-for="opt in protocolOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveDevice">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  listDevices, createDevice, updateDevice, deleteDevice, updateDeviceStatus,
  deviceTypeOptions, protocolOptions, type Device
} from '@/api/device'
import { useFamilyStore } from '@/stores/family'

const familyStore = useFamilyStore()
const loading = ref(false)
const saving = ref(false)
const devices = ref<Device[]>([])
const showCreateDialog = ref(false)
const editingDevice = ref<Device | null>(null)
const deviceFormRef = ref<FormInstance>()

const searchForm = reactive({
  location: '',
  deviceType: '',
  status: ''
})

const deviceForm = reactive({
  name: '',
  deviceType: '',
  vendor: '',
  location: '',
  power: 0,
  protocol: ''
})

const deviceRules: FormRules = {
  name: [{ required: true, message: '请输入设备名称', trigger: 'blur' }],
  deviceType: [{ required: true, message: '请选择设备类型', trigger: 'change' }],
  power: [{ required: true, message: '请输入功率', trigger: 'blur' }],
  protocol: [{ required: true, message: '请选择通信协议', trigger: 'change' }]
}

onMounted(() => {
  loadDevices()
})

watch(() => familyStore.currentFamilyId, () => {
  loadDevices()
})

async function loadDevices() {
  loading.value = true
  try {
    const res = await listDevices({
      familyId: familyStore.currentFamilyId || undefined,
      deviceType: searchForm.deviceType || undefined,
      status: searchForm.status || undefined,
      location: searchForm.location || undefined
    })
    devices.value = res
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

function editDevice(row: Device) {
  editingDevice.value = row
  Object.assign(deviceForm, {
    name: row.name,
    deviceType: row.deviceType,
    vendor: row.vendor,
    location: row.location,
    power: row.power,
    protocol: row.protocol
  })
  showCreateDialog.value = true
}

async function saveDevice() {
  if (!deviceFormRef.value) return
  await deviceFormRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true
      try {
        if (editingDevice.value) {
          await updateDevice(editingDevice.value.id, deviceForm)
          ElMessage.success('更新成功')
        } else {
          if (!familyStore.currentFamilyId) {
            ElMessage.warning('请先选择或创建家庭')
            return
          }
          await createDevice({ ...deviceForm, familyId: familyStore.currentFamilyId })
          ElMessage.success('添加成功')
        }
        showCreateDialog.value = false
        resetForm()
        loadDevices()
      } catch (e) {
        console.error(e)
      } finally {
        saving.value = false
      }
    }
  })
}

async function handleDeleteDevice(row: Device) {
  try {
    await ElMessageBox.confirm(`确定要删除设备"${row.name}"吗？`, '提示', { type: 'warning' })
    await deleteDevice(row.id)
    ElMessage.success('删除成功')
    loadDevices()
  } catch (e: any) {
    if (e !== 'cancel') console.error(e)
  }
}

async function toggleStatus(row: Device) {
  const newStatus = row.status === 'on' || row.status === 'online' ? 'off' : 'on'
  try {
    await updateDeviceStatus(row.id, newStatus)
    ElMessage.success('状态已更新')
    loadDevices()
  } catch (e) {
    console.error(e)
  }
}

function resetForm() {
  editingDevice.value = null
  Object.assign(deviceForm, {
    name: '',
    deviceType: '',
    vendor: '',
    location: '',
    power: 0,
    protocol: ''
  })
}

function getDeviceTypeLabel(type: string) {
  return deviceTypeOptions.find(o => o.value === type)?.label || type
}

function getProtocolLabel(proto: string) {
  return protocolOptions.find(o => o.value === proto)?.label || proto
}

function getStatusLabel(status: string) {
  const map: Record<string, string> = { online: '在线', offline: '离线', on: '开启', off: '关闭' }
  return map[status] || status
}

function getStatusType(status: string) {
  const map: Record<string, string> = { online: 'success', offline: 'info', on: 'success', off: 'warning' }
  return map[status] || 'info'
}
</script>
