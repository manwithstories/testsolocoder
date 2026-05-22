<template>
  <div class="guest-list">
    <div class="page-header">
      <el-input
        v-model="searchQuery"
        placeholder="搜索嘉宾..."
        :prefix-icon="Search"
        clearable
        style="width: 200px"
        @clear="fetchGuests"
        @keyup.enter="fetchGuests"
      />
      <el-select v-model="groupFilter" placeholder="分组" clearable style="width: 120px" @change="fetchGuests">
        <el-option label="新郎方" value="groom" />
        <el-option label="新娘方" value="bride" />
        <el-option label="共同" value="both" />
      </el-select>
      <el-select v-model="rsvpFilter" placeholder="RSVP" clearable style="width: 120px" @change="fetchGuests">
        <el-option label="待回复" value="pending" />
        <el-option label="已接受" value="accepted" />
        <el-option label="已拒绝" value="declined" />
      </el-select>
      <el-button type="primary" :icon="Plus" @click="showCreateDialog = true">
        添加嘉宾
      </el-button>
      <el-button :icon="Upload" @click="showImportDialog = true">
        导入
      </el-button>
      <el-button type="success" :icon="Download" @click="exportGuests">
        导出
      </el-button>
    </div>

    <el-tabs v-model="activeTab">
      <el-tab-pane label="嘉宾列表" name="list">
        <el-table :data="guests" v-loading="loading" stripe>
          <el-table-column label="姓名" min-width="140">
            <template #default="{ row }">
              <span :style="{ fontWeight: row.is_vip ? 'bold' : 'normal' }">
                {{ row.full_name }}
                <el-tag v-if="row.is_vip" type="warning" size="small" style="margin-left: 4px">VIP</el-tag>
              </span>
            </template>
          </el-table-column>
          <el-table-column prop="email" label="邮箱" width="180" />
          <el-table-column prop="phone" label="电话" width="130" />
          <el-table-column prop="group" label="分组" width="100">
            <template #default="{ row }">
              <el-tag size="small">{{ groupText(row.group) }}</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="RSVP" width="100">
            <template #default="{ row }">
              <el-tag :type="rsvpType(row.rsvp_status)" size="small">
                {{ rsvpText(row.rsvp_status) }}
              </el-tag>
            </template>
          </el-table-column>
          <el-table-column label="座位" width="100">
            <template #default="{ row }">
              <span v-if="row.table_id && row.seat_number">
                {{ getTableName(row.table_id) }}-{{ row.seat_number }}号
              </span>
              <el-tag v-else type="info" size="small">未安排</el-tag>
            </template>
          </el-table-column>
          <el-table-column label="操作" width="150" fixed="right">
            <template #default="{ row }">
              <el-button type="primary" link @click="editGuest(row)">编辑</el-button>
              <el-button type="danger" link @click="deleteGuest(row)">删除</el-button>
            </template>
          </el-table-column>
        </el-table>
      </el-tab-pane>

      <el-tab-pane label="座位表" name="tables">
        <div class="tables-container">
          <div class="tables-sidebar">
            <div class="sidebar-header">
              <span>餐桌列表</span>
              <el-button type="primary" size="small" :icon="Plus" @click="showTableDialog = true">添加</el-button>
            </div>
            <div
              v-for="table in tables"
              :key="table.id"
              class="table-item"
              :class="{ active: selectedTable?.id === table.id }"
              @click="selectTable(table)"
            >
              <div class="table-name">{{ table.table_name }}</div>
              <div class="table-info">
                <el-progress
                  :percentage="Math.round((table.table_guests?.length || 0) / table.capacity * 100)"
                  :stroke-width="6"
                  :show-text="false"
                  style="width: 80px"
                />
                <span>{{ table.table_guests?.length || 0 }}/{{ table.capacity }}</span>
              </div>
            </div>
            <el-empty v-if="tables.length === 0" description="暂无餐桌" :image-size="80" />
          </div>
          <div class="table-canvas" v-if="selectedTable">
            <div class="table-header">
              <h4>{{ selectedTable.table_name }}</h4>
              <div class="table-actions">
                <el-button type="danger" size="small" @click="deleteTable(selectedTable)">删除餐桌</el-button>
              </div>
            </div>
            <div class="seats-wrapper">
              <div
                v-for="seat in selectedTable.capacity"
                :key="seat"
                class="seat-item"
                :class="{ occupied: getSeatGuest(seat), selected: selectedSeat === seat }"
                @click="handleSeatClick(seat)"
                draggable="true"
                @dragover.prevent
                @drop="handleDrop($event, seat)"
              >
                <div class="seat-number">座位 {{ seat }}</div>
                <div v-if="getSeatGuest(seat)" class="seat-guest">
                  <el-icon><UserFilled /></el-icon>
                  <span>{{ getSeatGuest(seat)?.full_name }}</span>
                  <el-icon class="remove-icon" @click.stop="removeSeatAssignment(seat)"><Close /></el-icon>
                </div>
                <div v-else class="seat-empty">
                  <el-icon><Plus /></el-icon>
                  <span>点击安排</span>
                </div>
              </div>
            </div>
            <div class="drag-hint">
              <el-icon><InfoFilled /></el-icon>
              <span>提示：点击座位选择嘉宾，或拖拽嘉宾到座位上</span>
            </div>
          </div>
          <el-empty v-else description="请选择一个餐桌" style="flex: 1" />
        </div>
      </el-tab-pane>
    </el-tabs>

    <el-dialog v-model="showCreateDialog" :title="editingGuest ? '编辑嘉宾' : '添加嘉宾'" width="500px">
      <el-form ref="guestForm" :model="guestForm" :rules="guestRules" label-width="100px">
        <el-row :gutter="20">
          <el-col :span="12">
            <el-form-item label="姓" prop="last_name">
              <el-input v-model="guestForm.last_name" />
            </el-form-item>
          </el-col>
          <el-col :span="12">
            <el-form-item label="名" prop="first_name">
              <el-input v-model="guestForm.first_name" />
            </el-form-item>
          </el-col>
        </el-row>
        <el-form-item label="邮箱">
          <el-input v-model="guestForm.email" />
        </el-form-item>
        <el-form-item label="电话">
          <el-input v-model="guestForm.phone" />
        </el-form-item>
        <el-form-item label="分组">
          <el-select v-model="guestForm.group" style="width: 100%">
            <el-option label="新郎方" value="groom" />
            <el-option label="新娘方" value="bride" />
            <el-option label="共同" value="both" />
          </el-select>
        </el-form-item>
        <el-form-item label="RSVP">
          <el-select v-model="guestForm.rsvp_status" style="width: 100%">
            <el-option label="待回复" value="pending" />
            <el-option label="已接受" value="accepted" />
            <el-option label="已拒绝" value="declined" />
          </el-select>
        </el-form-item>
        <el-form-item label="VIP">
          <el-switch v-model="guestForm.is_vip" />
        </el-form-item>
        <el-form-item label="备注">
          <el-input v-model="guestForm.notes" type="textarea" :rows="2" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showCreateDialog = false">取消</el-button>
        <el-button type="primary" :loading="saving" @click="saveGuest">保存</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showImportDialog" title="导入嘉宾" width="400px">
      <el-upload
        drag
        action=""
        :auto-upload="false"
        :on-change="handleFileChange"
        accept=".csv,.xlsx,.xls"
      >
        <el-icon class="el-icon--upload"><UploadFilled /></el-icon>
        <div class="el-upload__text">将文件拖到此处，或<em>点击上传</em></div>
        <template #tip>
          <div class="el-upload__tip">支持 CSV、Excel 格式文件</div>
        </template>
      </el-upload>
      <template #footer>
        <el-button @click="showImportDialog = false">取消</el-button>
        <el-button type="primary" :loading="importing" :disabled="!importFile" @click="importGuests">
          导入
        </el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showTableDialog" title="添加餐桌" width="400px">
      <el-form :model="tableForm" label-width="80px">
        <el-form-item label="桌名">
          <el-input v-model="tableForm.table_name" placeholder="如：主桌、1号桌" />
        </el-form-item>
        <el-form-item label="桌号">
          <el-input-number v-model="tableForm.table_number" :min="1" />
        </el-form-item>
        <el-form-item label="容量">
          <el-input-number v-model="tableForm.capacity" :min="1" :max="30" />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="showTableDialog = false">取消</el-button>
        <el-button type="primary" @click="createTable">创建</el-button>
      </template>
    </el-dialog>

    <el-dialog v-model="showAssignDialog" title="安排嘉宾入座" width="400px" :close-on-click-modal="false">
      <div style="margin-bottom: 16px;">
        餐桌：<strong>{{ selectedTable?.table_name }}</strong>，座位号：<strong>{{ selectedSeat }}</strong>
      </div>
      <el-select
        v-model="selectedGuestId"
        placeholder="请选择嘉宾"
        filterable
        style="width: 100%"
        size="large"
      >
        <el-option
          v-for="guest in unassignedGuests"
          :key="guest.id"
          :label="guest.full_name + (guest.group ? ` (${groupText(guest.group)})` : '')"
          :value="guest.id"
        />
      </el-select>
      <template #footer>
        <el-button @click="showAssignDialog = false">取消</el-button>
        <el-button type="primary" :loading="assigning" @click="confirmAssign">确认安排</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { guestApi } from '@/api/guest'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Plus, Search, Upload, Download, UploadFilled, UserFilled, Close, InfoFilled } from '@element-plus/icons-vue'
import type { Guest, GuestTable } from '@/types'

const props = defineProps<{
  weddingId: number
}>()

const loading = ref(false)
const saving = ref(false)
const importing = ref(false)
const assigning = ref(false)
const guests = ref<Guest[]>([])
const tables = ref<any[]>([])
const searchQuery = ref('')
const groupFilter = ref('')
const rsvpFilter = ref('')
const activeTab = ref('list')
const showCreateDialog = ref(false)
const showImportDialog = ref(false)
const showTableDialog = ref(false)
const showAssignDialog = ref(false)
const editingGuest = ref<Guest | null>(null)
const importFile = ref<File | null>(null)
const selectedTable = ref<any>(null)
const selectedSeat = ref<number | null>(null)
const selectedGuestId = ref<number | null>(null)

const guestForm = reactive({
  first_name: '',
  last_name: '',
  email: '',
  phone: '',
  group: '',
  rsvp_status: 'pending',
  is_vip: false,
  notes: ''
})

const tableForm = reactive({
  table_name: '',
  table_number: 1,
  capacity: 10
})

const guestRules: FormRules = {
  first_name: [{ required: true, message: '请输入名', trigger: 'blur' }],
  last_name: [{ required: true, message: '请输入姓', trigger: 'blur' }]
}

const guestFormRef = ref<FormInstance>()

const unassignedGuests = computed(() => {
  return guests.value.filter(g => !g.table_id)
})

function groupText(group?: string) {
  const texts: Record<string, string> = { groom: '新郎方', bride: '新娘方', both: '共同' }
  return texts[group || ''] || '-'
}

function rsvpType(status: string) {
  const types: Record<string, string> = { pending: 'warning', accepted: 'success', declined: 'danger' }
  return types[status] || 'info'
}

function rsvpText(status: string) {
  const texts: Record<string, string> = { pending: '待回复', accepted: '已接受', declined: '已拒绝' }
  return texts[status] || status
}

function getTableName(tableId?: number) {
  const table = tables.value.find(t => t.id === tableId)
  return table?.table_name || ''
}

async function fetchGuests() {
  loading.value = true
  try {
    const res = await guestApi.getList(props.weddingId, {
      search: searchQuery.value,
      group: groupFilter.value,
      rsvp_status: rsvpFilter.value,
      page: 1,
      page_size: 500
    })
    guests.value = res.data.list
  } catch (error) {
    console.error('Failed to fetch guests:', error)
  } finally {
    loading.value = false
  }
}

async function fetchTables() {
  try {
    const res = await guestApi.getTables(props.weddingId)
    tables.value = res.data
  } catch (error) {
    console.error('Failed to fetch tables:', error)
  }
}

function editGuest(guest: Guest) {
  editingGuest.value = guest
  Object.assign(guestForm, {
    first_name: guest.first_name,
    last_name: guest.last_name,
    email: guest.email || '',
    phone: guest.phone || '',
    group: guest.group || '',
    rsvp_status: guest.rsvp_status,
    is_vip: guest.is_vip,
    notes: guest.notes || ''
  })
  showCreateDialog.value = true
}

async function saveGuest() {
  if (!guestFormRef.value) return
  
  await guestFormRef.value.validate(async (valid) => {
    if (valid) {
      saving.value = true
      try {
        const data = {
          ...guestForm,
          full_name: `${guestForm.last_name}${guestForm.first_name}`
        }
        if (editingGuest.value) {
          await guestApi.update(props.weddingId, editingGuest.value.id, data)
          ElMessage.success('更新成功')
        } else {
          await guestApi.create(props.weddingId, data)
          ElMessage.success('创建成功')
        }
        showCreateDialog.value = false
        fetchGuests()
        fetchTables()
      } catch (error: any) {
        ElMessage.error(error.message || '保存失败')
      } finally {
        saving.value = false
      }
    }
  })
}

async function deleteGuest(guest: Guest) {
  try {
    await ElMessageBox.confirm(`确定要删除嘉宾"${guest.full_name}"吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    
    await guestApi.delete(props.weddingId, guest.id)
    ElMessage.success('删除成功')
    fetchGuests()
    fetchTables()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete guest:', error)
    }
  }
}

function handleFileChange(file: any) {
  importFile.value = file.raw
}

async function importGuests() {
  if (!importFile.value) return
  
  importing.value = true
  try {
    await guestApi.import(props.weddingId, importFile.value)
    ElMessage.success('导入成功')
    showImportDialog.value = false
    fetchGuests()
  } catch (error: any) {
    ElMessage.error(error.message || '导入失败')
  } finally {
    importing.value = false
  }
}

async function exportGuests() {
  try {
    const res = await guestApi.export(props.weddingId) as any
    const blob = new Blob([res], {
      type: 'application/vnd.openxmlformats-officedocument.spreadsheetml.sheet'
    })
    const url = URL.createObjectURL(blob)
    const link = document.createElement('a')
    link.href = url
    link.download = 'guests.xlsx'
    link.click()
    URL.revokeObjectURL(url)
    ElMessage.success('导出成功')
  } catch (error) {
    ElMessage.error('导出失败')
  }
}

function selectTable(table: any) {
  selectedTable.value = table
  selectedSeat.value = null
}

function getSeatGuest(seatNumber: number) {
  if (!selectedTable.value) return null
  return selectedTable.value.table_guests?.find((g: Guest) => g.seat_number === seatNumber)
}

function handleSeatClick(seatNumber: number) {
  const existingGuest = getSeatGuest(seatNumber)
  if (existingGuest) {
    ElMessageBox.confirm(
      `座位 ${seatNumber} 当前安排的是"${existingGuest.full_name}"，要移除吗？`,
      '移除座位安排',
      {
        confirmButtonText: '移除',
        cancelButtonText: '取消',
        type: 'warning'
      }
    ).then(async () => {
      await guestApi.update(props.weddingId, existingGuest.id, {
        table_id: undefined,
        seat_number: 0
      })
      ElMessage.success('已移除座位安排')
      fetchGuests()
      fetchTables()
    }).catch(() => {})
    return
  }
  
  selectedSeat.value = seatNumber
  selectedGuestId.value = null
  showAssignDialog.value = true
}

function handleDrop(event: DragEvent, seatNumber: number) {
  event.preventDefault()
  const guestId = Number(event.dataTransfer?.getData('guestId'))
  if (guestId) {
    assignGuestToSeat(guestId, seatNumber)
  }
}

async function confirmAssign() {
  if (!selectedGuestId.value || !selectedSeat.value) {
    ElMessage.warning('请选择嘉宾')
    return
  }
  assigning.value = true
  try {
    await assignGuestToSeat(selectedGuestId.value, selectedSeat.value)
    showAssignDialog.value = false
    ElMessage.success('安排成功')
  } catch (error: any) {
    ElMessage.error(error.message || '安排失败')
  } finally {
    assigning.value = false
  }
}

async function assignGuestToSeat(guestId: number, seatNumber: number) {
  await guestApi.assignSeat(props.weddingId, {
    guest_id: guestId,
    table_id: selectedTable.value.id,
    seat_number: seatNumber
  })
  await fetchGuests()
  await fetchTables()
}

async function removeSeatAssignment(seatNumber: number) {
  const guest = getSeatGuest(seatNumber)
  if (!guest) return
  
  try {
    await ElMessageBox.confirm(`确定要移除"${guest.full_name}"的座位安排吗？`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await guestApi.update(props.weddingId, guest.id, {
      table_id: undefined,
      seat_number: 0
    })
    ElMessage.success('已移除')
    fetchGuests()
    fetchTables()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to remove assignment:', error)
    }
  }
}

async function createTable() {
  try {
    await guestApi.createTable(props.weddingId, tableForm)
    ElMessage.success('创建成功')
    showTableDialog.value = false
    tableForm.table_name = ''
    tableForm.table_number = 1
    tableForm.capacity = 10
    fetchTables()
  } catch (error: any) {
    ElMessage.error(error.message || '创建失败')
  }
}

async function deleteTable(table: any) {
  try {
    await ElMessageBox.confirm(`确定要删除餐桌"${table.table_name}"吗？该餐桌上的嘉宾安排将被移除！`, '提示', {
      confirmButtonText: '确定',
      cancelButtonText: '取消',
      type: 'warning'
    })
    await guestApi.deleteTable(props.weddingId, table.id)
    ElMessage.success('删除成功')
    if (selectedTable.value?.id === table.id) {
      selectedTable.value = null
    }
    fetchTables()
    fetchGuests()
  } catch (error) {
    if (error !== 'cancel') {
      console.error('Failed to delete table:', error)
    }
  }
}

onMounted(() => {
  fetchGuests()
  fetchTables()
})
</script>

<style scoped>
.guest-list {
  padding: 0;
}

.page-header {
  display: flex;
  gap: 16px;
  margin-bottom: 20px;
}

.tables-container {
  display: flex;
  gap: 20px;
  min-height: 500px;
  background: #fff;
  border-radius: 8px;
  padding: 20px;
}

.tables-sidebar {
  width: 220px;
  border-right: 1px solid #ebeef5;
  padding-right: 16px;
  display: flex;
  flex-direction: column;
}

.sidebar-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 16px;
  font-weight: 600;
}

.table-item {
  padding: 12px;
  margin-bottom: 8px;
  border: 1px solid #e4e7ed;
  border-radius: 6px;
  cursor: pointer;
  transition: all 0.3s;
}

.table-item:hover {
  border-color: #409EFF;
  background: #f5f7fa;
}

.table-item.active {
  border-color: #409EFF;
  background: #ecf5ff;
}

.table-name {
  font-weight: 600;
  margin-bottom: 6px;
}

.table-info {
  display: flex;
  align-items: center;
  gap: 8px;
  font-size: 12px;
  color: #909399;
}

.table-canvas {
  flex: 1;
  padding-left: 20px;
  display: flex;
  flex-direction: column;
}

.table-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 24px;
  padding-bottom: 12px;
  border-bottom: 1px solid #ebeef5;
}

.table-header h4 {
  margin: 0;
}

.seats-wrapper {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(160px, 1fr));
  gap: 16px;
  flex: 1;
}

.seat-item {
  border: 2px dashed #dcdfe6;
  border-radius: 8px;
  padding: 16px;
  min-height: 90px;
  cursor: pointer;
  transition: all 0.3s;
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  background: #fafbfc;
}

.seat-item:hover {
  border-color: #409EFF;
  background: #f0f7ff;
  transform: translateY(-2px);
}

.seat-item.selected {
  border-color: #E6A23C;
  background: #fdf6ec;
}

.seat-item.occupied {
  border-style: solid;
  border-color: #67C23A;
  background: #f0f9eb;
}

.seat-number {
  font-size: 12px;
  color: #909399;
  margin-bottom: 8px;
}

.seat-guest {
  display: flex;
  align-items: center;
  gap: 4px;
  font-size: 13px;
  font-weight: 600;
  color: #303133;
  position: relative;
  width: 100%;
  justify-content: center;
}

.remove-icon {
  color: #F56C6C;
  cursor: pointer;
  margin-left: 4px;
}

.remove-icon:hover {
  color: #c45656;
}

.seat-empty {
  display: flex;
  flex-direction: column;
  align-items: center;
  gap: 4px;
  color: #c0c4cc;
  font-size: 12px;
}

.drag-hint {
  margin-top: 20px;
  padding: 12px 16px;
  background: #fdf6ec;
  border-radius: 6px;
  color: #E6A23C;
  font-size: 13px;
  display: flex;
  align-items: center;
  gap: 6px;
}
</style>
