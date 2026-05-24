<template>
  <div class="page-container">
    <div class="page-header">
      <h2>设备分组</h2>
      <el-button type="primary" :icon="Plus" @click="showCreateDialog = true">
        新建分组
      </el-button>
    </div>

    <div class="search-bar">
      <el-input v-model="searchForm.name" placeholder="搜索分组名称" style="width: 200px" clearable @input="loadGroups" />
      <el-select v-model="searchForm.type" placeholder="分组类型" style="width: 160px" clearable @change="loadGroups">
        <el-option v-for="opt in groupTypeOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
      </el-select>
    </div>

    <div class="card-grid">
      <div v-for="group in groups" :key="group.id" class="device-card">
        <div class="device-header">
          <div>
            <el-tag :type="getGroupTypeColor(group.type)" size="small">
              {{ getGroupTypeLabel(group.type) }}
            </el-tag>
            <span class="device-name" style="margin-left: 8px;">{{ group.name }}</span>
          </div>
          <el-dropdown @command="(cmd: string) => handleCommand(cmd, group)">
            <el-icon class="cursor-pointer"><MoreFilled /></el-icon>
            <template #dropdown>
              <el-dropdown-menu>
                <el-dropdown-item command="edit">编辑</el-dropdown-item>
                <el-dropdown-item command="delete">删除</el-dropdown-item>
              </el-dropdown-menu>
            </template>
          </el-dropdown>
        </div>
        <div class="device-info">
          <p v-if="group.description">{{ group.description }}</p>
          <p>设备数量: {{ group.devices?.length || 0 }}</p>
        </div>
        <div class="device-actions">
          <el-button size="small" type="success" @click="batchControl(group, 'on')">全部开启</el-button>
          <el-button size="small" type="warning" @click="batchControl(group, 'off')">全部关闭</el-button>
          <el-button size="small" @click="showDevices(group)">查看设备</el-button>
        </div>
      </div>
    </div>

    <el-empty v-if="groups.length === 0 && !loading" description="暂无分组" />

    <el-dialog v-model="showCreateDialog" :title="editingGroup ? '编辑分组' : '新建分组'" width="500px">
      <el-form :model="groupForm" :rules="groupRules" ref="groupFormRef" label-width="80px">
        <el-form-item label="分组名称" prop="name">
          <el-input v-model="groupForm.name" />
        </el-form-item>
        <el-form-item label="分组类型" prop="type">
          <el-select v-model="groupForm.type" style="width: 100%">
            <el-option v-for="opt in groupTypeOptions" :key="opt.value" :label="opt.label" :value="opt.value" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="groupForm.description" type="textarea" :rows="3" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveGroup">确定</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showDevicesDialog" title="分组设备" width="600px">
      <div v-if="currentGroup">
        <div class="mb-20 flex-between">
          <span>添加设备到 {{ currentGroup.name }}</span>
          <el-select v-model="selectedDeviceId" placeholder="选择设备" style="width: 200px" clearable>
            <el-option v-for="d in availableDevices" :key="d.id" :label="d.name" :value="d.id" />
          </el-select>
          <el-button type="primary" size="small" :disabled="!selectedDeviceId" @click="handleAddDeviceToGroup">添加</el-button>
        </div>
        <el-table :data="currentGroup.devices || []" style="width: 100%">
          <el-table-column prop="name" label="设备名称" />
          <el-table-column prop="deviceType" label="类型">
            <template #default="{ row }">{{ getDeviceTypeLabel(row.deviceType) }}</template>
          </el-table-column>
          <el-table-column prop="status" label="状态">
            <template #default="{ row }">
              <el-tag :type="getStatusType(row.status)">{{ getStatusLabel(row.status) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="100">
            <template #default="{ row }">
              <el-button type="danger" link size="small" @click="removeDevice(row)">移除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </div>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted, watch } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus } from '@element-plus/icons-vue'
import {
  listGroups, createGroup, updateGroup, deleteGroup,
  addDeviceToGroup, removeDeviceFromGroup, batchControlGroup,
  groupTypeOptions, type DeviceGroup
} from '@/api/group'
import { listDevices, deviceTypeOptions, type Device } from '@/api/device'
import { useFamilyStore } from '@/stores/family'

const familyStore = useFamilyStore()
const loading = ref(false)
const saving = ref(false)
const groups = ref<DeviceGroup[]>([])
const showCreateDialog = ref(false)
const showDevicesDialog = ref(false)
const editingGroup = ref<DeviceGroup | null>(null)
const currentGroup = ref<DeviceGroup | null>(null)
const selectedDeviceId = ref<number | null>(null)
const availableDevices = ref<Device[]>([])
const groupFormRef = ref<FormInstance>()

const searchForm = reactive({ name: '', type: '' })

const groupForm = reactive({
  name: '',
  type: 'room',
  description: ''
})

const groupRules: FormRules = {
  name: [{ required: true, message: '请输入分组名称', trigger: 'blur' }],
  type: [{ required: true, message: '请选择分组类型', trigger: 'change' }]
}

onMounted(() => {
  loadGroups()
  loadAvailableDevices()
})

watch(() => familyStore.currentFamilyId, () => {
  loadGroups()
  loadAvailableDevices()
})

async function loadAvailableDevices() {
  try {
    availableDevices.value = await listDevices({ familyId: familyStore.currentFamilyId || undefined })
  } catch (e) {
    console.error(e)
  }
}

async function loadGroups() {
  loading.value = true
  try {
    let res = await listGroups({ familyId: familyStore.currentFamilyId || undefined })
    if (searchForm.type) {
      res = res.filter(g => g.type === searchForm.type)
    }
    if (searchForm.name) {
      res = res.filter(g => g.name.includes(searchForm.name))
    }
    groups.value = res
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function handleCommand(cmd: string, group: DeviceGroup) {
  if (cmd === 'edit') {
    editingGroup.value = group
    Object.assign(groupForm, { name: group.name, type: group.type, description: group.description })
    showCreateDialog.value = true
  } else if (cmd === 'delete') {
    try {
      await ElMessageBox.confirm(`确定要删除分组"${group.name}"吗？`, '提示', { type: 'warning' })
      await deleteGroup(group.id)
      ElMessage.success('删除成功')
      loadGroups()
    } catch (e: any) {
      if (e !== 'cancel') console.error(e)
    }
  }
}

async function saveGroup() {
  if (!groupFormRef.value) return
  await groupFormRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true
      try {
        if (editingGroup.value) {
          await updateGroup(editingGroup.value.id, groupForm)
          ElMessage.success('更新成功')
        } else {
          if (!familyStore.currentFamilyId) {
            ElMessage.warning('请先选择或创建家庭')
            return
          }
          await createGroup({ ...groupForm, familyId: familyStore.currentFamilyId })
          ElMessage.success('创建成功')
        }
        showCreateDialog.value = false
        resetForm()
        loadGroups()
      } catch (e) {
        console.error(e)
      } finally {
        saving.value = false
      }
    }
  })
}

async function showDevices(group: DeviceGroup) {
  currentGroup.value = group
  const allDevices = await listDevices()
  const groupDeviceIds = (group.devices || []).map((d: any) => d.id)
  availableDevices.value = allDevices.filter((d: Device) => !groupDeviceIds.includes(d.id))
  showDevicesDialog.value = true
}

async function handleAddDeviceToGroup() {
  if (!currentGroup.value || !selectedDeviceId.value) return
  try {
    await addDeviceToGroup(currentGroup.value.id, selectedDeviceId.value)
    ElMessage.success('添加成功')
    const res = await getGroup(currentGroup.value.id)
    currentGroup.value = res
    const groupDeviceIds = (res.devices || []).map((d: any) => d.id)
    availableDevices.value = availableDevices.value.filter(d => !groupDeviceIds.includes(d.id))
    selectedDeviceId.value = null
  } catch (e) {
    console.error(e)
  }
}

async function removeDevice(device: any) {
  if (!currentGroup.value) return
  try {
    await removeDeviceFromGroup(currentGroup.value.id, device.id)
    ElMessage.success('移除成功')
    const res = await getGroup(currentGroup.value.id)
    currentGroup.value = res
  } catch (e) {
    console.error(e)
  }
}

async function batchControl(group: DeviceGroup, action: string) {
  try {
    await batchControlGroup(group.id, action)
    ElMessage.success('批量操作成功')
    loadGroups()
  } catch (e) {
    console.error(e)
  }
}

async function getGroup(id: number): Promise<DeviceGroup> {
  const groups = await listGroups()
  return groups.find(g => g.id === id) || {} as DeviceGroup
}

function resetForm() {
  editingGroup.value = null
  Object.assign(groupForm, { name: '', type: 'room', description: '' })
}

function getGroupTypeLabel(type: string) {
  return groupTypeOptions.find(o => o.value === type)?.label || type
}

function getGroupTypeColor(type: string) {
  const map: Record<string, string> = { room: '', floor: 'success', function: 'warning', custom: 'info' }
  return map[type] || ''
}

function getDeviceTypeLabel(type: string) {
  return deviceTypeOptions.find(o => o.value === type)?.label || type
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
