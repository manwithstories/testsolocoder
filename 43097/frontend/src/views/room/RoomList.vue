<template>
  <div class="room-container">
    <el-card shadow="hover">
      <div class="header-bar">
        <div class="filter-bar">
          <el-select
            v-model="filterStatus"
            placeholder="选择状态"
            clearable
            style="width: 150px"
            @change="fetchList"
          >
            <el-option label="空闲" :value="RoomStatus.AVAILABLE" />
            <el-option label="已入住" :value="RoomStatus.OCCUPIED" />
            <el-option label="已预订" :value="RoomStatus.RESERVED" />
            <el-option label="维修中" :value="RoomStatus.MAINTENANCE" />
            <el-option label="清洁中" :value="RoomStatus.CLEANING" />
          </el-select>
          <el-select
            v-model="filterRoomType"
            placeholder="选择房型"
            clearable
            style="width: 180px"
            @change="fetchList"
          >
            <el-option
              v-for="type in roomTypeOptions"
              :key="type.id"
              :label="type.name"
              :value="type.id"
            />
          </el-select>
          <el-input
            v-model="searchKeyword"
            placeholder="搜索房号"
            clearable
            :prefix-icon="Search"
            style="width: 200px"
            @keyup.enter="handleSearch"
          />
          <el-button type="primary" :icon="Search" @click="handleSearch">搜索</el-button>
          <el-button :icon="Refresh" @click="resetFilter">重置</el-button>
        </div>
        <div class="action-bar">
          <el-upload
            :show-file-list="false"
            accept=".xlsx,.xls"
            :before-upload="handleImport"
          >
            <el-button :icon="Upload">批量导入</el-button>
          </el-upload>
          <el-button type="primary" :icon="Plus" @click="handleAdd">添加房间</el-button>
        </div>
      </div>

      <el-table :data="tableData" v-loading="loading" border stripe>
        <el-table-column prop="id" label="ID" width="80" />
        <el-table-column prop="roomNo" label="房号" width="100" />
        <el-table-column prop="floor" label="楼层" width="80">
          <template #default="{ row }">
            {{ row.floor }}楼
          </template>
        </el-table-column>
        <el-table-column label="房型" min-width="120">
          <template #default="{ row }">
            {{ row.roomType?.name || '-' }}
          </template>
        </el-table-column>
        <el-table-column label="价格" width="120">
          <template #default="{ row }">
            <span class="price-text">¥{{ row.roomType?.price?.toFixed(2) || '-' }}/晚</span>
          </template>
        </el-table-column>
        <el-table-column label="状态" width="100">
          <template #default="{ row }">
            <el-tag :type="getStatusType(row.status)">
              {{ getStatusText(row.status) }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column label="设施" min-width="200">
          <template #default="{ row }">
            <el-tag
              v-for="(facility, index) in row.roomType?.facilities || []"
              :key="index"
              size="small"
              style="margin-right: 4px; margin-bottom: 4px"
            >
              {{ facility }}
            </el-tag>
          </template>
        </el-table-column>
        <el-table-column prop="description" label="备注" min-width="150" show-overflow-tooltip />
        <el-table-column label="操作" width="200" fixed="right">
          <template #default="{ row }">
            <el-button type="primary" link :icon="Edit" @click="handleEdit(row)">编辑</el-button>
            <el-button type="danger" link :icon="Delete" @click="handleDelete(row)">删除</el-button>
          </template>
        </el-table-column>
      </el-table>

      <div class="pagination-bar">
        <el-pagination
          v-model:current-page="pagination.page"
          v-model:page-size="pagination.pageSize"
          :page-sizes="[10, 20, 50, 100]"
          :total="pagination.total"
          layout="total, sizes, prev, pager, next, jumper"
          @size-change="fetchList"
          @current-change="fetchList"
        />
      </div>
    </el-card>

    <el-dialog
      v-model="dialogVisible"
      :title="isEdit ? '编辑房间' : '添加房间'"
      width="500px"
      destroy-on-close
    >
      <el-form
        ref="formRef"
        :model="formData"
        :rules="formRules"
        label-width="100px"
        label-position="right"
      >
        <el-form-item label="房号" prop="roomNo">
          <el-input v-model="formData.roomNo" placeholder="请输入房号" />
        </el-form-item>
        <el-form-item label="楼层" prop="floor">
          <el-input-number
            v-model="formData.floor"
            :min="1"
            :max="99"
            style="width: 100%"
            placeholder="请输入楼层"
          />
        </el-form-item>
        <el-form-item label="房型" prop="roomTypeId">
          <el-select
            v-model="formData.roomTypeId"
            placeholder="请选择房型"
            style="width: 100%"
          >
            <el-option
              v-for="type in roomTypeOptions"
              :key="type.id"
              :label="type.name"
              :value="type.id"
            />
          </el-select>
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-select v-model="formData.status" placeholder="请选择状态" style="width: 100%">
            <el-option label="空闲" :value="RoomStatus.AVAILABLE" />
            <el-option label="已入住" :value="RoomStatus.OCCUPIED" />
            <el-option label="已预订" :value="RoomStatus.RESERVED" />
            <el-option label="维修中" :value="RoomStatus.MAINTENANCE" />
            <el-option label="清洁中" :value="RoomStatus.CLEANING" />
          </el-select>
        </el-form-item>
        <el-form-item label="备注" prop="description">
          <el-input
            v-model="formData.description"
            type="textarea"
            :rows="3"
            placeholder="请输入备注"
          />
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitLoading" @click="handleSubmit">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { Search, Refresh, Plus, Edit, Delete, Upload } from '@element-plus/icons-vue'
import * as XLSX from 'xlsx'
import {
  getRoomList,
  createRoom,
  updateRoom,
  deleteRoom,
  getAllRoomTypes
} from '@/api/room'
import { RoomStatus, type Room, type RoomType, type PageParams } from '@/types'

const loading = ref(false)
const submitLoading = ref(false)
const searchKeyword = ref('')
const filterStatus = ref<RoomStatus | ''>('')
const filterRoomType = ref<number | ''>('')
const tableData = ref<Room[]>([])
const roomTypeOptions = ref<RoomType[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const formRef = ref<FormInstance>()

const pagination = reactive({
  page: 1,
  pageSize: 10,
  total: 0
})

const formData = reactive<Partial<Room>>({
  roomNo: '',
  floor: 1,
  roomTypeId: undefined,
  status: RoomStatus.AVAILABLE,
  description: ''
})

const formRules: FormRules = {
  roomNo: [{ required: true, message: '请输入房号', trigger: 'blur' }],
  floor: [{ required: true, message: '请输入楼层', trigger: 'blur' }],
  roomTypeId: [{ required: true, message: '请选择房型', trigger: 'change' }],
  status: [{ required: true, message: '请选择状态', trigger: 'change' }]
}

const getStatusType = (status: RoomStatus) => {
  const typeMap: Record<RoomStatus, string> = {
    [RoomStatus.AVAILABLE]: 'success',
    [RoomStatus.OCCUPIED]: 'danger',
    [RoomStatus.RESERVED]: 'warning',
    [RoomStatus.MAINTENANCE]: 'info',
    [RoomStatus.CLEANING]: 'primary'
  }
  return typeMap[status] || 'info'
}

const getStatusText = (status: RoomStatus) => {
  const textMap: Record<RoomStatus, string> = {
    [RoomStatus.AVAILABLE]: '空闲',
    [RoomStatus.OCCUPIED]: '已入住',
    [RoomStatus.RESERVED]: '已预订',
    [RoomStatus.MAINTENANCE]: '维修中',
    [RoomStatus.CLEANING]: '清洁中'
  }
  return textMap[status] || '未知'
}

const fetchRoomTypes = async () => {
  try {
    const res = await getAllRoomTypes()
    roomTypeOptions.value = res
  } catch (error) {
    console.error('Failed to fetch room types:', error)
  }
}

const fetchList = async () => {
  loading.value = true
  try {
    const params: PageParams & { status?: RoomStatus; roomTypeId?: number; keyword?: string } = {
      page: pagination.page,
      pageSize: pagination.pageSize
    }
    if (filterStatus.value) {
      params.status = filterStatus.value
    }
    if (filterRoomType.value) {
      params.roomTypeId = filterRoomType.value
    }
    if (searchKeyword.value) {
      params.keyword = searchKeyword.value
    }
    const res = await getRoomList(params)
    tableData.value = res.list
    pagination.total = res.total
  } catch (error) {
    console.error('Failed to fetch room list:', error)
  } finally {
    loading.value = false
  }
}

const handleSearch = () => {
  pagination.page = 1
  fetchList()
}

const resetFilter = () => {
  searchKeyword.value = ''
  filterStatus.value = ''
  filterRoomType.value = ''
  pagination.page = 1
  fetchList()
}

const handleAdd = () => {
  isEdit.value = false
  Object.assign(formData, {
    roomNo: '',
    floor: 1,
    roomTypeId: undefined,
    status: RoomStatus.AVAILABLE,
    description: ''
  })
  dialogVisible.value = true
}

const handleEdit = (row: Room) => {
  isEdit.value = true
  Object.assign(formData, {
    id: row.id,
    roomNo: row.roomNo,
    floor: row.floor,
    roomTypeId: row.roomTypeId,
    status: row.status,
    description: row.description
  })
  dialogVisible.value = true
}

const handleDelete = (row: Room) => {
  ElMessageBox.confirm(`确定要删除房间"${row.roomNo}"吗？`, '删除确认', {
    confirmButtonText: '确定',
    cancelButtonText: '取消',
    type: 'warning'
  })
    .then(async () => {
      try {
        await deleteRoom(row.id)
        ElMessage.success('删除成功')
        fetchList()
      } catch (error) {
        console.error('Failed to delete room:', error)
      }
    })
    .catch(() => {})
}

const handleSubmit = async () => {
  if (!formRef.value) return
  await formRef.value.validate(async (valid) => {
    if (valid) {
      submitLoading.value = true
      try {
        if (isEdit.value && formData.id) {
          await updateRoom(formData.id, formData)
          ElMessage.success('编辑成功')
        } else {
          await createRoom(formData as Omit<Room, 'id' | 'createdAt' | 'updatedAt'>)
          ElMessage.success('添加成功')
        }
        dialogVisible.value = false
        fetchList()
      } catch (error) {
        console.error('Failed to submit room:', error)
      } finally {
        submitLoading.value = false
      }
    }
  })
}

const handleImport = async (file: File) => {
  try {
    const data = await file.arrayBuffer()
    const workbook = XLSX.read(data)
    const worksheet = workbook.Sheets[workbook.SheetNames[0]]
    const jsonData = XLSX.utils.sheet_to_json(worksheet)

    if (!jsonData.length) {
      ElMessage.warning('Excel文件中没有数据')
      return false
    }

    const importData = jsonData.map((item: any) => ({
      roomNo: String(item['房号'] || item['roomNo'] || ''),
      floor: Number(item['楼层'] || item['floor'] || 1),
      roomTypeId: Number(item['房型ID'] || item['roomTypeId'] || 1),
      status: (item['状态'] || item['status'] || RoomStatus.AVAILABLE) as RoomStatus,
      description: String(item['备注'] || item['description'] || '')
    }))

    ElMessageBox.confirm(
      `确定要导入 ${importData.length} 条房间数据吗？`,
      '导入确认',
      {
        confirmButtonText: '确定',
        cancelButtonText: '取消',
        type: 'info'
      }
    )
      .then(async () => {
        try {
          const successCount = 0
          for (const item of importData) {
            try {
              await createRoom(item as Omit<Room, 'id' | 'createdAt' | 'updatedAt'>)
            } catch (e) {
              console.error('Failed to import room:', item, e)
            }
          }
          ElMessage.success(`导入完成`)
          fetchList()
        } catch (error) {
          console.error('Failed to import rooms:', error)
        }
      })
      .catch(() => {})

    return false
  } catch (error) {
    console.error('Failed to parse Excel:', error)
    ElMessage.error('解析Excel文件失败')
    return false
  }
}

onMounted(() => {
  fetchRoomTypes()
  fetchList()
})
</script>

<style scoped lang="scss">
.room-container {
  .header-bar {
    display: flex;
    justify-content: space-between;
    align-items: center;
    margin-bottom: 20px;
    flex-wrap: wrap;
    gap: 12px;

    .filter-bar {
      display: flex;
      gap: 12px;
      flex-wrap: wrap;
    }

    .action-bar {
      display: flex;
      gap: 12px;
    }
  }

  .pagination-bar {
    display: flex;
    justify-content: flex-end;
    margin-top: 20px;
  }

  .price-text {
    color: #f56c6c;
    font-weight: 600;
  }
}
</style>
