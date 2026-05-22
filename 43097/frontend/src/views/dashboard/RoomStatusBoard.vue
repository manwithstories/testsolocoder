<template>
  <div class="page-container">
    <div class="page-header">
      <h2 class="page-title">房间状态看板</h2>
      <div class="page-actions">
        <el-button :icon="Refresh" @click="fetchData">刷新</el-button>
      </div>
    </div>

    <div class="filter-bar">
      <div class="filter-form">
        <div class="filter-item">
          <label class="filter-label">楼层筛选</label>
          <el-select v-model="selectedFloor" placeholder="全部楼层" style="width: 150px" clearable>
            <el-option
              v-for="floor in floors"
              :key="floor"
              :label="`${floor}楼`"
              :value="floor"
            />
          </el-select>
        </div>
        <div class="filter-item">
          <label class="filter-label">房间状态</label>
          <el-select v-model="selectedStatus" placeholder="全部状态" style="width: 150px" clearable>
            <el-option
              v-for="status in statusOptions"
              :key="status.value"
              :label="status.label"
              :value="status.value"
            />
          </el-select>
        </div>
        <div class="filter-item">
          <label class="filter-label">房型</label>
          <el-select v-model="selectedRoomType" placeholder="全部房型" style="width: 180px" clearable>
            <el-option
              v-for="type in roomTypes"
              :key="type.id"
              :label="type.name"
              :value="type.id"
            />
          </el-select>
        </div>
        <div class="filter-actions">
          <el-button type="primary" :icon="Search" @click="handleSearch">查询</el-button>
          <el-button :icon="Refresh" @click="handleReset">重置</el-button>
        </div>
      </div>
    </div>

    <div class="common-card mb-20">
      <div class="card-header">
        <h3 class="card-title">房间状态统计</h3>
      </div>
      <div class="card-body">
        <el-row :gutter="16">
          <el-col :xs="12" :sm="6">
            <div class="status-stat status-available" @click="filterByStatus(RoomStatus.AVAILABLE)">
              <div class="stat-number">{{ statusCounts.available }}</div>
              <div class="stat-label">空闲</div>
            </div>
          </el-col>
          <el-col :xs="12" :sm="6">
            <div class="status-stat status-occupied" @click="filterByStatus(RoomStatus.OCCUPIED)">
              <div class="stat-number">{{ statusCounts.occupied }}</div>
              <div class="stat-label">入住中</div>
            </div>
          </el-col>
          <el-col :xs="12" :sm="6">
            <div class="status-stat status-reserved" @click="filterByStatus(RoomStatus.RESERVED)">
              <div class="stat-number">{{ statusCounts.reserved }}</div>
              <div class="stat-label">已预订</div>
            </div>
          </el-col>
          <el-col :xs="12" :sm="6">
            <div class="status-stat status-maintenance" @click="filterByStatus(RoomStatus.MAINTENANCE)">
              <div class="stat-number">{{ statusCounts.maintenance }}</div>
              <div class="stat-label">维修中</div>
            </div>
          </el-col>
        </el-row>
      </div>
    </div>

    <div class="common-card">
      <div class="card-header">
        <h3 class="card-title">房间列表</h3>
        <div class="status-legend">
          <div class="legend-item">
            <span class="legend-color status-available"></span>
            <span>空闲</span>
          </div>
          <div class="legend-item">
            <span class="legend-color status-occupied"></span>
            <span>入住中</span>
          </div>
          <div class="legend-item">
            <span class="legend-color status-reserved"></span>
            <span>已预订</span>
          </div>
          <div class="legend-item">
            <span class="legend-color status-cleaning"></span>
            <span>清洁中</span>
          </div>
          <div class="legend-item">
            <span class="legend-color status-maintenance"></span>
            <span>维修中</span>
          </div>
        </div>
      </div>
      <div class="card-body">
        <div v-for="floor in filteredFloors" :key="floor" class="floor-section">
          <div class="floor-title">
            <el-icon><OfficeBuilding /></el-icon>
            <span>{{ floor }}楼</span>
            <span class="floor-count">
              {{ getFloorRooms(floor).length }} 间房
            </span>
          </div>
          <div class="room-grid">
            <div
              v-for="room in getFloorRooms(floor)"
              :key="room.id"
              class="room-card"
              :class="`status-${room.status}`"
              @click="handleRoomClick(room)"
            >
              <div class="room-header">
                <span class="room-number">{{ room.roomNumber }}</span>
                <span class="room-status-tag">{{ getStatusLabel(room.status) }}</span>
              </div>
              <div class="room-type">{{ room.roomType?.name || '标准间' }}</div>
              <div class="room-info" v-if="room.status === RoomStatus.OCCUPIED && room.currentGuest">
                <div class="guest-info">
                  <el-icon><User /></el-icon>
                  <span class="truncate">{{ room.currentGuest.name }}</span>
                </div>
                <div class="guest-info">
                  <el-icon><Phone /></el-icon>
                  <span>{{ room.currentGuest.phone }}</span>
                </div>
              </div>
              <div class="room-info" v-else-if="room.status === RoomStatus.RESERVED && room.booking">
                <div class="guest-info">
                  <el-icon><User /></el-icon>
                  <span class="truncate">{{ room.booking.guestName }}</span>
                </div>
                <div class="booking-info">
                  <el-icon><Calendar /></el-icon>
                  <span>{{ room.booking.checkInDate }}</span>
                </div>
              </div>
              <div class="room-footer">
                <span class="room-price">¥{{ room.roomType?.price || 199 }}/晚</span>
              </div>
            </div>
          </div>
        </div>

        <el-empty v-if="filteredRooms.length === 0" description="暂无符合条件的房间" />
      </div>
    </div>

    <el-dialog
      v-model="detailDialogVisible"
      :title="`房间详情 - ${currentRoom?.roomNumber}`"
      width="500px"
      :close-on-click-modal="false"
    >
      <div v-if="currentRoom" class="room-detail">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="房号">
            {{ currentRoom.roomNumber }}
          </el-descriptions-item>
          <el-descriptions-item label="楼层">
            {{ currentRoom.floor }}楼
          </el-descriptions-item>
          <el-descriptions-item label="房型">
            {{ currentRoom.roomType?.name || '标准间' }}
          </el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="getStatusTagType(currentRoom.status)">
              {{ getStatusLabel(currentRoom.status) }}
            </el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="价格">
            ¥{{ currentRoom.roomType?.price || 199 }}/晚
          </el-descriptions-item>
          <el-descriptions-item label="面积">
            {{ currentRoom.roomType?.area || 25 }}㎡
          </el-descriptions-item>
          <el-descriptions-item label="床型">
            {{ currentRoom.roomType?.bedType || '大床房' }}
          </el-descriptions-item>
          <el-descriptions-item label="可住人数">
            {{ currentRoom.roomType?.capacity || 2 }}人
          </el-descriptions-item>
          <el-descriptions-item label="设施" :span="2">
            <el-tag
              v-for="facility in currentRoom.roomType?.facilities || ['WiFi', '电视', '空调']"
              :key="facility"
              size="small"
              style="margin-right: 4px; margin-bottom: 4px;"
            >
              {{ facility }}
            </el-tag>
          </el-descriptions-item>
        </el-descriptions>

        <div v-if="currentRoom.status === RoomStatus.OCCUPIED && currentRoom.currentGuest" class="guest-detail">
          <h4>当前客人信息</h4>
          <el-descriptions :column="2" border size="small">
            <el-descriptions-item label="姓名">
              {{ currentRoom.currentGuest.name }}
            </el-descriptions-item>
            <el-descriptions-item label="电话">
              {{ currentRoom.currentGuest.phone }}
            </el-descriptions-item>
            <el-descriptions-item label="身份证号">
              {{ currentRoom.currentGuest.idCard }}
            </el-descriptions-item>
            <el-descriptions-item label="入住时间">
              {{ currentRoom.currentGuest.checkInTime }}
            </el-descriptions-item>
            <el-descriptions-item label="预计退房" :span="2">
              {{ currentRoom.currentGuest.expectedCheckOutTime }}
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </div>
      <template #footer>
        <div class="dialog-footer">
          <el-button @click="detailDialogVisible = false">关闭</el-button>
          <el-button
            v-if="currentRoom?.status === RoomStatus.OCCUPIED"
            type="warning"
            @click="handleCheckOut"
          >
            办理退房
          </el-button>
          <el-button
            v-if="currentRoom?.status === RoomStatus.AVAILABLE"
            type="primary"
            @click="handleCheckIn"
          >
            办理入住
          </el-button>
        </div>
      </template>
    </el-dialog>
  </div>
</template>

<script setup lang="ts">
import { ref, reactive, computed, onMounted } from 'vue'
import { ElMessage } from 'element-plus'
import {
  Refresh,
  Search,
  OfficeBuilding,
  User,
  Phone,
  Calendar
} from '@element-plus/icons-vue'
import { getRoomList, getAllRoomTypes } from '@/api/room'
import { RoomStatus, type Room, type RoomType } from '@/types'

interface RoomWithGuest extends Room {
  currentGuest?: {
    name: string
    phone: string
    idCard: string
    checkInTime: string
    expectedCheckOutTime: string
  }
  booking?: {
    guestName: string
    checkInDate: string
  }
}

const loading = ref(false)
const detailDialogVisible = ref(false)
const currentRoom = ref<RoomWithGuest | null>(null)

const selectedFloor = ref<number | null>(null)
const selectedStatus = ref<RoomStatus | null>(null)
const selectedRoomType = ref<number | null>(null)

const roomTypes = ref<RoomType[]>([])
const allRooms = ref<RoomWithGuest[]>([])

const statusOptions = [
  { label: '空闲', value: RoomStatus.AVAILABLE },
  { label: '入住中', value: RoomStatus.OCCUPIED },
  { label: '已预订', value: RoomStatus.RESERVED },
  { label: '清洁中', value: RoomStatus.CLEANING },
  { label: '维修中', value: RoomStatus.MAINTENANCE }
]

const guestNames = ['张三', '李四', '王五', '赵六', '钱七', '孙八', '周九', '吴十', '郑十一', '王十二']

const floors = computed(() => {
  const floorSet = new Set(allRooms.value.map(room => room.floor))
  return Array.from(floorSet).sort((a, b) => a - b)
})

const statusCounts = computed(() => {
  const counts = {
    available: 0,
    occupied: 0,
    reserved: 0,
    cleaning: 0,
    maintenance: 0
  }
  
  allRooms.value.forEach(room => {
    switch (room.status) {
      case RoomStatus.AVAILABLE:
        counts.available++
        break
      case RoomStatus.OCCUPIED:
        counts.occupied++
        break
      case RoomStatus.RESERVED:
        counts.reserved++
        break
      case RoomStatus.CLEANING:
        counts.cleaning++
        break
      case RoomStatus.MAINTENANCE:
        counts.maintenance++
        break
    }
  })
  
  return counts
})

const filteredRooms = computed(() => {
  return allRooms.value.filter(room => {
    if (selectedFloor.value !== null && room.floor !== selectedFloor.value) {
      return false
    }
    if (selectedStatus.value !== null && room.status !== selectedStatus.value) {
      return false
    }
    if (selectedRoomType.value !== null && room.roomTypeId !== selectedRoomType.value) {
      return false
    }
    return true
  })
})

const filteredFloors = computed(() => {
  const floorSet = new Set(filteredRooms.value.map(room => room.floor))
  return Array.from(floorSet).sort((a, b) => a - b)
})

const getFloorRooms = (floor: number) => {
  return filteredRooms.value.filter(room => room.floor === floor)
}

const getStatusLabel = (status: RoomStatus) => {
  const map: Record<RoomStatus, string> = {
    [RoomStatus.AVAILABLE]: '空闲',
    [RoomStatus.OCCUPIED]: '入住中',
    [RoomStatus.RESERVED]: '已预订',
    [RoomStatus.CLEANING]: '清洁中',
    [RoomStatus.MAINTENANCE]: '维修中'
  }
  return map[status] || status
}

const getStatusTagType = (status: RoomStatus) => {
  const map: Record<RoomStatus, string> = {
    [RoomStatus.AVAILABLE]: 'success',
    [RoomStatus.OCCUPIED]: 'danger',
    [RoomStatus.RESERVED]: 'warning',
    [RoomStatus.CLEANING]: 'primary',
    [RoomStatus.MAINTENANCE]: 'info'
  }
  return map[status] || 'info'
}

const filterByStatus = (status: RoomStatus) => {
  selectedStatus.value = selectedStatus.value === status ? null : status
  handleSearch()
}

const handleRoomClick = (room: RoomWithGuest) => {
  currentRoom.value = room
  detailDialogVisible.value = true
}

const handleCheckIn = () => {
  ElMessage.success(`正在为房间 ${currentRoom.value?.roomNumber} 办理入住`)
  detailDialogVisible.value = false
}

const handleCheckOut = () => {
  ElMessage.success(`正在为房间 ${currentRoom.value?.roomNumber} 办理退房`)
  detailDialogVisible.value = false
}

const handleSearch = () => {
}

const handleReset = () => {
  selectedFloor.value = null
  selectedStatus.value = null
  selectedRoomType.value = null
}

const generateMockData = (): RoomWithGuest[] => {
  const rooms: RoomWithGuest[] = []
  const roomTypeNames = ['标准大床房', '豪华双床房', '商务套房', '家庭房', '总统套房']
  const bedTypes = ['大床房', '双床房', '特大床']
  const statuses = [
    RoomStatus.AVAILABLE,
    RoomStatus.AVAILABLE,
    RoomStatus.OCCUPIED,
    RoomStatus.OCCUPIED,
    RoomStatus.OCCUPIED,
    RoomStatus.RESERVED,
    RoomStatus.CLEANING,
    RoomStatus.MAINTENANCE
  ]

  let id = 1
  for (let floor = 1; floor <= 6; floor++) {
    for (let num = 1; num <= 8; num++) {
      const roomNumber = `${floor}0${num}`
      const roomTypeId = Math.floor(Math.random() * 5) + 1
      const status = statuses[Math.floor(Math.random() * statuses.length)]
      
      const room: RoomWithGuest = {
        id,
        roomNumber,
        floor,
        roomTypeId,
        status,
        roomType: {
          id: roomTypeId,
          name: roomTypeNames[roomTypeId - 1] || '标准间',
          price: 199 + roomTypeId * 100,
          capacity: 2 + Math.floor(roomTypeId / 2),
          bedType: bedTypes[Math.floor(Math.random() * bedTypes.length)],
          area: 20 + roomTypeId * 5,
          facilities: ['WiFi', '电视', '空调', '独立卫浴'],
          status: true,
          createdAt: '',
          updatedAt: ''
        },
        createdAt: '',
        updatedAt: ''
      }

      if (status === RoomStatus.OCCUPIED) {
        const guestIndex = Math.floor(Math.random() * guestNames.length)
        room.currentGuest = {
          name: guestNames[guestIndex],
          phone: `138${Math.floor(Math.random() * 100000000)}`,
          idCard: `110101199${Math.floor(Math.random() * 10)}${String(Math.floor(Math.random() * 12)).padStart(2, '0')}${String(Math.floor(Math.random() * 28)).padStart(2, '0')}${Math.floor(Math.random() * 10000)}`,
          checkInTime: '2024-01-15 14:30',
          expectedCheckOutTime: '2024-01-18 12:00'
        }
      }

      if (status === RoomStatus.RESERVED) {
        const guestIndex = Math.floor(Math.random() * guestNames.length)
        room.booking = {
          guestName: guestNames[guestIndex],
          checkInDate: '2024-01-20'
        }
      }

      rooms.push(room)
      id++
    }
  }

  return rooms
}

const fetchData = async () => {
  loading.value = true
  try {
    const [roomTypeData, roomData] = await Promise.all([
      getAllRoomTypes(),
      getRoomList({ page: 1, pageSize: 100 })
    ]) as any

    if (roomTypeData) {
      roomTypes.value = roomTypeData
    } else {
      roomTypes.value = [
        { id: 1, name: '标准大床房', price: 199, capacity: 2, bedType: '大床房', area: 25, facilities: ['WiFi', '电视', '空调'], status: true, createdAt: '', updatedAt: '' },
        { id: 2, name: '豪华双床房', price: 299, capacity: 2, bedType: '双床房', area: 30, facilities: ['WiFi', '电视', '空调', '浴缸'], status: true, createdAt: '', updatedAt: '' },
        { id: 3, name: '商务套房', price: 499, capacity: 2, bedType: '特大床', area: 45, facilities: ['WiFi', '电视', '空调', '浴缸', '商务桌'], status: true, createdAt: '', updatedAt: '' },
        { id: 4, name: '家庭房', price: 399, capacity: 4, bedType: '双床房', area: 40, facilities: ['WiFi', '电视', '空调', '儿童用品'], status: true, createdAt: '', updatedAt: '' },
        { id: 5, name: '总统套房', price: 999, capacity: 2, bedType: '特大床', area: 80, facilities: ['WiFi', '电视', '空调', '浴缸', '桑拿', '管家服务'], status: true, createdAt: '', updatedAt: '' }
      ]
    }

    if (roomData && roomData.list) {
      allRooms.value = roomData.list
    } else {
      allRooms.value = generateMockData()
    }
  } catch (error) {
    console.error('Failed to fetch room data:', error)
    roomTypes.value = [
      { id: 1, name: '标准大床房', price: 199, capacity: 2, bedType: '大床房', area: 25, facilities: ['WiFi', '电视', '空调'], status: true, createdAt: '', updatedAt: '' },
      { id: 2, name: '豪华双床房', price: 299, capacity: 2, bedType: '双床房', area: 30, facilities: ['WiFi', '电视', '空调', '浴缸'], status: true, createdAt: '', updatedAt: '' },
      { id: 3, name: '商务套房', price: 499, capacity: 2, bedType: '特大床', area: 45, facilities: ['WiFi', '电视', '空调', '浴缸', '商务桌'], status: true, createdAt: '', updatedAt: '' },
      { id: 4, name: '家庭房', price: 399, capacity: 4, bedType: '双床房', area: 40, facilities: ['WiFi', '电视', '空调', '儿童用品'], status: true, createdAt: '', updatedAt: '' },
      { id: 5, name: '总统套房', price: 999, capacity: 2, bedType: '特大床', area: 80, facilities: ['WiFi', '电视', '空调', '浴缸', '桑拿', '管家服务'], status: true, createdAt: '', updatedAt: '' }
    ]
    allRooms.value = generateMockData()
  } finally {
    loading.value = false
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped lang="scss">
.status-stat {
  padding: 20px;
  border-radius: 8px;
  text-align: center;
  cursor: pointer;
  transition: all 0.3s;
  color: #fff;
  margin-bottom: 16px;

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.15);
  }

  .stat-number {
    font-size: 32px;
    font-weight: 700;
    line-height: 1.2;
  }

  .stat-label {
    font-size: 14px;
    margin-top: 4px;
    opacity: 0.9;
  }

  &.status-available {
    background: linear-gradient(135deg, #67c23a 0%, #85ce61 100%);
  }

  &.status-occupied {
    background: linear-gradient(135deg, #f56c6c 0%, #f78989 100%);
  }

  &.status-reserved {
    background: linear-gradient(135deg, #e6a23c 0%, #ebb563 100%);
  }

  &.status-maintenance {
    background: linear-gradient(135deg, #909399 0%, #a6a9ad 100%);
  }
}

.status-legend {
  display: flex;
  gap: 16px;
  flex-wrap: wrap;

  .legend-item {
    display: flex;
    align-items: center;
    gap: 6px;
    font-size: 13px;
    color: #606266;

    .legend-color {
      width: 14px;
      height: 14px;
      border-radius: 3px;

      &.status-available {
        background: #67c23a;
      }

      &.status-occupied {
        background: #f56c6c;
      }

      &.status-reserved {
        background: #e6a23c;
      }

      &.status-cleaning {
        background: #409eff;
      }

      &.status-maintenance {
        background: #909399;
      }
    }
  }
}

.floor-section {
  margin-bottom: 24px;

  &:last-child {
    margin-bottom: 0;
  }

  .floor-title {
    display: flex;
    align-items: center;
    gap: 8px;
    font-size: 16px;
    font-weight: 600;
    color: #303133;
    margin-bottom: 12px;
    padding-left: 8px;
    border-left: 4px solid #409eff;

    .floor-count {
      font-size: 13px;
      font-weight: normal;
      color: #909399;
      margin-left: 4px;
    }
  }
}

.room-grid {
  display: grid;
  grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
  gap: 12px;
}

.room-card {
  background: #fff;
  border-radius: 8px;
  padding: 16px;
  border: 2px solid transparent;
  cursor: pointer;
  transition: all 0.3s;
  position: relative;
  overflow: hidden;

  &::before {
    content: '';
    position: absolute;
    top: 0;
    left: 0;
    right: 0;
    height: 4px;
  }

  &:hover {
    transform: translateY(-2px);
    box-shadow: 0 6px 20px rgba(0, 0, 0, 0.12);
  }

  &.status-available {
    &::before {
      background: #67c23a;
    }
    border-color: rgba(103, 194, 58, 0.3);
  }

  &.status-occupied {
    &::before {
      background: #f56c6c;
    }
    border-color: rgba(245, 108, 108, 0.3);
  }

  &.status-reserved {
    &::before {
      background: #e6a23c;
    }
    border-color: rgba(230, 162, 60, 0.3);
  }

  &.status-cleaning {
    &::before {
      background: #409eff;
    }
    border-color: rgba(64, 158, 255, 0.3);
  }

  &.status-maintenance {
    &::before {
      background: #909399;
    }
    border-color: rgba(144, 147, 153, 0.3);
  }

  .room-header {
    display: flex;
    justify-content: space-between;
    align-items: flex-start;
    margin-bottom: 8px;

    .room-number {
      font-size: 20px;
      font-weight: 700;
      color: #303133;
    }

    .room-status-tag {
      font-size: 11px;
      padding: 2px 8px;
      border-radius: 10px;
      color: #fff;

      .status-available & {
        background: #67c23a;
      }

      .status-occupied & {
        background: #f56c6c;
      }

      .status-reserved & {
        background: #e6a23c;
      }

      .status-cleaning & {
        background: #409eff;
      }

      .status-maintenance & {
        background: #909399;
      }
    }
  }

  .room-type {
    font-size: 13px;
    color: #606266;
    margin-bottom: 8px;
  }

  .room-info {
    font-size: 12px;
    color: #909399;
    margin-bottom: 4px;

    .guest-info,
    .booking-info {
      display: flex;
      align-items: center;
      gap: 4px;
      margin-bottom: 2px;

      .el-icon {
        font-size: 12px;
      }
    }
  }

  .room-footer {
    margin-top: 8px;
    padding-top: 8px;
    border-top: 1px solid #f0f2f5;

    .room-price {
      font-size: 16px;
      font-weight: 600;
      color: #f56c6c;
    }
  }
}

.room-detail {
  .guest-detail {
    margin-top: 20px;

    h4 {
      font-size: 14px;
      font-weight: 600;
      color: #303133;
      margin-bottom: 12px;
    }
  }
}

@media (max-width: 768px) {
  .status-legend {
    justify-content: center;
    gap: 12px;
  }

  .room-grid {
    grid-template-columns: repeat(auto-fill, minmax(150px, 1fr));
    gap: 8px;
  }

  .room-card {
    padding: 12px;

    .room-header .room-number {
      font-size: 18px;
    }
  }

  .status-stat {
    padding: 16px;

    .stat-number {
      font-size: 24px;
    }
  }
}
</style>
