<template>
  <div class="rooms">
    <el-card class="filter-card">
      <el-form :inline="true" :model="filterForm">
        <el-form-item label="楼层">
          <el-select v-model="filterForm.floor" placeholder="全部楼层" clearable @change="loadRooms">
            <el-option v-for="floor in floors" :key="floor" :label="floor" :value="floor" />
          </el-select>
        </el-form-item>
        <el-form-item label="设备">
          <el-input v-model="filterForm.equipment" placeholder="搜索设备" clearable @keyup.enter="loadRooms" />
        </el-form-item>
        <el-form-item>
          <el-button type="primary" @click="loadRooms">
            <el-icon><Search /></el-icon>
            搜索
          </el-button>
        </el-form-item>
        <el-form-item v-if="isAdmin || isSpaceAdmin">
          <el-button type="success" @click="showCreateDialog">
            <el-icon><Plus /></el-icon>
            新建会议室
          </el-button>
        </el-form-item>
      </el-form>
    </el-card>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="8" v-for="room in rooms" :key="room.id">
        <el-card class="room-card" shadow="hover">
          <div class="room-cover">
            <el-image v-if="room.photos?.length" :src="room.photos[0].url" fit="cover" class="room-image" />
            <div v-else class="room-placeholder">
              <el-icon :size="48"><OfficeBuilding /></el-icon>
            </div>
            <div class="room-badge" :class="room.status === 1 ? 'active' : 'inactive'">
              {{ room.status === 1 ? '可用' : '停用' }}
            </div>
          </div>
          <div class="room-info">
            <h3 class="room-name">{{ room.name }}</h3>
            <div class="room-meta">
              <span><el-icon><Location /></el-icon> {{ room.floor }}层</span>
              <span><el-icon><User /></el-icon> {{ room.capacity }}人</span>
            </div>
            <div class="room-price">
              <span class="price">¥{{ room.price_per_hour }}</span>
              <span class="unit">/小时</span>
            </div>
            <div class="room-availability">
              <el-icon><Clock /></el-icon>
              <span>{{ room.available_start || '08:00' }} - {{ room.available_end || '22:00' }}</span>
            </div>
            <div class="room-equipment" v-if="room.equipment">
              <el-tag v-for="eq in parseEquipment(room.equipment)" :key="eq" size="small" type="info">{{ eq }}</el-tag>
            </div>
            <div class="room-actions" v-if="isAdmin || isSpaceAdmin">
              <el-button type="primary" size="small" @click="showEditDialog(room)">编辑</el-button>
              <el-button type="danger" size="small" @click="deleteRoom(room)">删除</el-button>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-empty v-if="rooms.length === 0 && !loading" description="暂无会议室" style="margin-top: 60px" />

    <el-dialog v-model="dialogVisible" :title="isEdit ? '编辑会议室' : '新建会议室'" width="600px">
      <el-form ref="roomFormRef" :model="roomForm" :rules="roomRules" label-width="100px">
        <el-form-item label="会议室名称" prop="name">
          <el-input v-model="roomForm.name" placeholder="请输入会议室名称" />
        </el-form-item>
        <el-form-item label="所在楼层" prop="floor">
          <el-input v-model="roomForm.floor" placeholder="请输入楼层" />
        </el-form-item>
        <el-form-item label="容纳人数" prop="capacity">
          <el-input-number v-model="roomForm.capacity" :min="1" :max="100" />
        </el-form-item>
        <el-form-item label="位置" prop="location">
          <el-input v-model="roomForm.location" placeholder="请输入位置描述" />
        </el-form-item>
        <el-form-item label="每小时价格" prop="price_per_hour">
          <el-input-number v-model="roomForm.price_per_hour" :min="0" :precision="2" />
        </el-form-item>
        <el-form-item label="可用开始时间">
          <el-time-picker
            v-model="roomForm.available_start"
            format="HH:mm"
            value-format="HH:mm"
            placeholder="选择开始时间"
          />
        </el-form-item>
        <el-form-item label="可用结束时间">
          <el-time-picker
            v-model="roomForm.available_end"
            format="HH:mm"
            value-format="HH:mm"
            placeholder="选择结束时间"
          />
        </el-form-item>
        <el-form-item label="设备清单">
          <el-select v-model="roomForm.equipment" multiple placeholder="选择设备" style="width: 100%">
            <el-option label="投影仪" value="投影仪" />
            <el-option label="电视" value="电视" />
            <el-option label="白板" value="白板" />
            <el-option label="视频会议" value="视频会议" />
            <el-option label="电话会议" value="电话会议" />
            <el-option label="WiFi" value="WiFi" />
            <el-option label="空调" value="空调" />
            <el-option label="饮水机" value="饮水机" />
          </el-select>
        </el-form-item>
        <el-form-item label="描述">
          <el-input v-model="roomForm.description" type="textarea" :rows="3" placeholder="请输入会议室描述" />
        </el-form-item>
        <el-form-item label="状态" v-if="isEdit">
          <el-radio-group v-model="roomForm.status">
            <el-radio :value="1">可用</el-radio>
            <el-radio :value="0">停用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
      <template #footer>
        <el-button @click="dialogVisible = false">取消</el-button>
        <el-button type="primary" :loading="submitting" @click="saveRoom">确定</el-button>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, onMounted } from 'vue'
import { ElMessage, ElMessageBox, type FormInstance, type FormRules } from 'element-plus'
import { api } from '@/api'
import { useUserStore } from '@/stores/user'

const userStore = useUserStore()
const isAdmin = ref(userStore.isAdmin)
const isSpaceAdmin = ref(userStore.isSpaceAdmin)

const loading = ref(false)
const rooms = ref<any[]>([])
const floors = ref<string[]>([])
const dialogVisible = ref(false)
const isEdit = ref(false)
const submitting = ref(false)
const roomFormRef = ref<FormInstance>()
const editRoomId = ref<number | null>(null)

const filterForm = reactive({
  floor: '',
  equipment: ''
})

const roomForm = reactive({
  name: '',
  floor: '',
  capacity: 1,
  location: '',
  price_per_hour: 0,
  available_start: '08:00',
  available_end: '22:00',
  equipment: [] as string[],
  description: '',
  status: 1
})

const roomRules: FormRules = {
  name: [{ required: true, message: '请输入会议室名称', trigger: 'blur' }],
  floor: [{ required: true, message: '请输入楼层', trigger: 'blur' }],
  capacity: [{ required: true, message: '请输入容纳人数', trigger: 'change' }],
  price_per_hour: [{ required: true, message: '请输入价格', trigger: 'change' }]
}

onMounted(() => {
  loadRooms()
  loadFloors()
})

async function loadRooms() {
  loading.value = true
  try {
    const res: any = await api.getRooms({
      page: 1,
      page_size: 100,
      floor: filterForm.floor || undefined,
      equipment: filterForm.equipment || undefined
    })
    rooms.value = res.data?.rooms || []
  } catch (e) {
    console.error(e)
  } finally {
    loading.value = false
  }
}

async function loadFloors() {
  try {
    const res: any = await api.getFloors()
    floors.value = res.data || []
  } catch (e) {
    console.error(e)
  }
}

function parseEquipment(equipment: string) {
  try {
    return JSON.parse(equipment)
  } catch {
    return []
  }
}

function showCreateDialog() {
  isEdit.value = false
  editRoomId.value = null
  Object.assign(roomForm, {
    name: '',
    floor: '',
    capacity: 1,
    location: '',
    price_per_hour: 0,
    available_start: '08:00',
    available_end: '22:00',
    equipment: [],
    description: '',
    status: 1
  })
  dialogVisible.value = true
}

function showEditDialog(room: any) {
  isEdit.value = true
  editRoomId.value = room.id
  Object.assign(roomForm, {
    name: room.name,
    floor: room.floor,
    capacity: room.capacity,
    location: room.location,
    price_per_hour: room.price_per_hour,
    available_start: room.available_start || '08:00',
    available_end: room.available_end || '22:00',
    equipment: parseEquipment(room.equipment),
    description: room.description,
    status: room.status
  })
  dialogVisible.value = true
}

async function saveRoom() {
  if (!roomFormRef.value) return
  await roomFormRef.value.validate(async (valid) => {
    if (valid) {
      submitting.value = true
      try {
        if (isEdit.value && editRoomId.value) {
          await api.updateRoom(editRoomId.value, roomForm)
          ElMessage.success('更新成功')
        } else {
          await api.createRoom(roomForm)
          ElMessage.success('创建成功')
        }
        dialogVisible.value = false
        loadRooms()
      } catch (e: any) {
        console.error(e)
      } finally {
        submitting.value = false
      }
    }
  })
}

async function deleteRoom(room: any) {
  try {
    await ElMessageBox.confirm(`确定删除会议室 "${room.name}" 吗？`, '提示', {
      type: 'warning'
    })
    await api.deleteRoom(room.id)
    ElMessage.success('删除成功')
    loadRooms()
  } catch (e: any) {
    if (e !== 'cancel') console.error(e)
  }
}
</script>

<style scoped>
.filter-card {
  border-radius: 8px;
}

.room-card {
  border-radius: 8px;
  margin-bottom: 20px;
}

.room-card :deep(.el-card__body) {
  padding: 0;
}

.room-cover {
  position: relative;
  height: 160px;
  overflow: hidden;
  border-radius: 8px 8px 0 0;
}

.room-image {
  width: 100%;
  height: 100%;
}

.room-placeholder {
  width: 100%;
  height: 100%;
  display: flex;
  align-items: center;
  justify-content: center;
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
}

.room-badge {
  position: absolute;
  top: 12px;
  right: 12px;
  padding: 4px 12px;
  border-radius: 12px;
  font-size: 12px;
  color: #fff;
}

.room-badge.active {
  background: #67C23A;
}

.room-badge.inactive {
  background: #909399;
}

.room-info {
  padding: 16px;
}

.room-name {
  margin: 0 0 8px;
  font-size: 16px;
  color: #303133;
}

.room-meta {
  display: flex;
  gap: 16px;
  color: #909399;
  font-size: 13px;
  margin-bottom: 8px;
}

.room-meta span {
  display: flex;
  align-items: center;
  gap: 4px;
}

.room-price {
  display: flex;
  align-items: baseline;
  gap: 4px;
  margin-bottom: 8px;
}

.room-price .price {
  font-size: 20px;
  font-weight: 600;
  color: #F56C6C;
}

.room-price .unit {
  font-size: 12px;
  color: #909399;
}

.room-availability {
  display: flex;
  align-items: center;
  gap: 4px;
  color: #606266;
  font-size: 13px;
  margin-bottom: 8px;
}

.room-equipment {
  display: flex;
  flex-wrap: wrap;
  gap: 4px;
  margin-bottom: 12px;
}

.room-actions {
  display: flex;
  gap: 8px;
  padding-top: 12px;
  border-top: 1px solid #ebeef5;
}
</style>
