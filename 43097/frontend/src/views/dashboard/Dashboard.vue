<template>
  <div class="dashboard-container">
    <el-row :gutter="24" class="stats-row">
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon icon-total">
              <el-icon :size="32"><OfficeBuilding /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ stats.availableRooms + stats.occupiedRooms }}</p>
              <p class="stat-label">总房间数</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon icon-checkin">
              <el-icon :size="32"><Calendar /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ stats.todayCheckIns }}</p>
              <p class="stat-label">今日入住</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon icon-checkout">
              <el-icon :size="32"><Sunny /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ stats.todayCheckOuts }}</p>
              <p class="stat-label">今日退房</p>
            </div>
          </div>
        </el-card>
      </el-col>
      <el-col :xs="12" :sm="6">
        <el-card class="stat-card" shadow="hover">
          <div class="stat-content">
            <div class="stat-icon icon-occupied">
              <el-icon :size="32"><HomeFilled /></el-icon>
            </div>
            <div class="stat-info">
              <p class="stat-value">{{ stats.occupiedRooms }}</p>
              <p class="stat-label">当前在住</p>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="24" class="room-status-row">
      <el-col :xs="24" :lg="16">
        <el-card class="room-status-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span class="card-title">房间状态看板</span>
            </div>
          </template>
          <div class="room-status-grid">
            <div
              v-for="room in mockRooms"
              :key="room.id"
              class="room-item"
              :class="`status-${room.status}`"
            >
              <span class="room-number">{{ room.roomNumber }}</span>
              <span class="room-floor">{{ room.floor }}楼</span>
            </div>
          </div>
          <div class="status-legend">
            <div class="legend-item">
              <span class="legend-color status-available"></span>
              <span>空闲</span>
            </div>
            <div class="legend-item">
              <span class="legend-color status-occupied"></span>
              <span>已入住</span>
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
        </el-card>
      </el-col>
      <el-col :xs="24" :lg="8">
        <el-card class="quick-stats-card" shadow="hover">
          <template #header>
            <div class="card-header">
              <span class="card-title">快捷统计</span>
            </div>
          </template>
          <div class="quick-stats-list">
            <div class="quick-stat-item">
              <span class="stat-label">今日预订</span>
              <span class="stat-value text-primary">{{ stats.todayBookings }}</span>
            </div>
            <div class="quick-stat-item">
              <span class="stat-label">待处理预订</span>
              <span class="stat-value text-warning">{{ stats.pendingBookings }}</span>
            </div>
            <div class="quick-stat-item">
              <span class="stat-label">今日营收</span>
              <span class="stat-value text-success">¥{{ stats.todayRevenue.toFixed(2) }}</span>
            </div>
            <div class="quick-stat-item">
              <span class="stat-label">本月营收</span>
              <span class="stat-value text-info">¥{{ stats.monthRevenue.toFixed(2) }}</span>
            </div>
            <div class="quick-stat-item">
              <span class="stat-label">入住率</span>
              <span class="stat-value text-danger">
                {{ totalRooms > 0 ? ((stats.occupiedRooms / totalRooms) * 100).toFixed(1) : 0 }}%
              </span>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { OfficeBuilding, Calendar, Sunny, HomeFilled } from '@element-plus/icons-vue'
import { getDashboardStats, getRoomStatusOverview } from '@/api/dashboard'
import { RoomStatus, type DashboardStats } from '@/types'

const stats = ref<DashboardStats>({
  todayCheckIns: 0,
  todayCheckOuts: 0,
  occupiedRooms: 0,
  availableRooms: 0,
  todayRevenue: 0,
  monthRevenue: 0,
  todayBookings: 0,
  pendingBookings: 0
})

const mockRooms = ref([
  { id: 1, roomNumber: '101', floor: 1, status: RoomStatus.AVAILABLE },
  { id: 2, roomNumber: '102', floor: 1, status: RoomStatus.OCCUPIED },
  { id: 3, roomNumber: '103', floor: 1, status: RoomStatus.RESERVED },
  { id: 4, roomNumber: '104', floor: 1, status: RoomStatus.AVAILABLE },
  { id: 5, roomNumber: '105', floor: 1, status: RoomStatus.CLEANING },
  { id: 6, roomNumber: '201', floor: 2, status: RoomStatus.OCCUPIED },
  { id: 7, roomNumber: '202', floor: 2, status: RoomStatus.AVAILABLE },
  { id: 8, roomNumber: '203', floor: 2, status: RoomStatus.MAINTENANCE },
  { id: 9, roomNumber: '204', floor: 2, status: RoomStatus.OCCUPIED },
  { id: 10, roomNumber: '205', floor: 2, status: RoomStatus.RESERVED },
  { id: 11, roomNumber: '301', floor: 3, status: RoomStatus.AVAILABLE },
  { id: 12, roomNumber: '302', floor: 3, status: RoomStatus.OCCUPIED },
  { id: 13, roomNumber: '303', floor: 3, status: RoomStatus.AVAILABLE },
  { id: 14, roomNumber: '304', floor: 3, status: RoomStatus.CLEANING },
  { id: 15, roomNumber: '305', floor: 3, status: RoomStatus.OCCUPIED },
  { id: 16, roomNumber: '401', floor: 4, status: RoomStatus.RESERVED },
  { id: 17, roomNumber: '402', floor: 4, status: RoomStatus.AVAILABLE },
  { id: 18, roomNumber: '403', floor: 4, status: RoomStatus.OCCUPIED },
  { id: 19, roomNumber: '404', floor: 4, status: RoomStatus.AVAILABLE },
  { id: 20, roomNumber: '405', floor: 4, status: RoomStatus.MAINTENANCE }
])

const totalRooms = computed(() => stats.value.availableRooms + stats.value.occupiedRooms)

const fetchData = async () => {
  try {
    const [statsData, roomStatusData] = await Promise.all([
      getDashboardStats(),
      getRoomStatusOverview()
    ])
    stats.value = statsData
  } catch (error) {
    console.error('Failed to fetch dashboard data:', error)
  }
}

onMounted(() => {
  fetchData()
})
</script>

<style scoped lang="scss">
.dashboard-container {
  .stats-row {
    margin-bottom: 24px;
  }

  .stat-card {
    border-radius: 8px;
    margin-bottom: 16px;

    .stat-content {
      display: flex;
      align-items: center;
      gap: 16px;

      .stat-icon {
        width: 64px;
        height: 64px;
        border-radius: 12px;
        display: flex;
        align-items: center;
        justify-content: center;
        color: #fff;

        &.icon-total {
          background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
        }

        &.icon-checkin {
          background: linear-gradient(135deg, #11998e 0%, #38ef7d 100%);
        }

        &.icon-checkout {
          background: linear-gradient(135deg, #f093fb 0%, #f5576c 100%);
        }

        &.icon-occupied {
          background: linear-gradient(135deg, #fa709a 0%, #fee140 100%);
        }
      }

      .stat-info {
        .stat-value {
          font-size: 28px;
          font-weight: 600;
          color: #303133;
          margin: 0;
          line-height: 1.2;
        }

        .stat-label {
          font-size: 14px;
          color: #909399;
          margin: 4px 0 0 0;
        }
      }
    }
  }

  .room-status-row {
    .card-header {
      display: flex;
      justify-content: space-between;
      align-items: center;

      .card-title {
        font-size: 16px;
        font-weight: 600;
        color: #303133;
      }
    }
  }

  .room-status-card {
    margin-bottom: 24px;

    .room-status-grid {
      display: grid;
      grid-template-columns: repeat(auto-fill, minmax(80px, 1fr));
      gap: 12px;
      margin-bottom: 24px;

      .room-item {
        padding: 12px 8px;
        border-radius: 8px;
        text-align: center;
        cursor: pointer;
        transition: all 0.3s;

        &:hover {
          transform: translateY(-2px);
          box-shadow: 0 4px 12px rgba(0, 0, 0, 0.15);
        }

        .room-number {
          display: block;
          font-size: 16px;
          font-weight: 600;
          color: #fff;
        }

        .room-floor {
          display: block;
          font-size: 12px;
          color: rgba(255, 255, 255, 0.8);
          margin-top: 4px;
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

        &.status-cleaning {
          background: linear-gradient(135deg, #409eff 0%, #66b1ff 100%);
        }

        &.status-maintenance {
          background: linear-gradient(135deg, #909399 0%, #a6a9ad 100%);
        }
      }
    }

    .status-legend {
      display: flex;
      flex-wrap: wrap;
      gap: 16px;
      padding-top: 16px;
      border-top: 1px solid #ebeef5;

      .legend-item {
        display: flex;
        align-items: center;
        gap: 8px;
        font-size: 13px;
        color: #606266;

        .legend-color {
          width: 16px;
          height: 16px;
          border-radius: 4px;

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
  }

  .quick-stats-card {
    margin-bottom: 24px;

    .quick-stats-list {
      .quick-stat-item {
        display: flex;
        justify-content: space-between;
        align-items: center;
        padding: 16px 0;
        border-bottom: 1px solid #f0f2f5;

        &:last-child {
          border-bottom: none;
        }

        .stat-label {
          font-size: 14px;
          color: #606266;
        }

        .stat-value {
          font-size: 16px;
          font-weight: 600;

          &.text-primary {
            color: #409eff;
          }

          &.text-success {
            color: #67c23a;
          }

          &.text-warning {
            color: #e6a23c;
          }

          &.text-danger {
            color: #f56c6c;
          }

          &.text-info {
            color: #909399;
          }
        }
      }
    }
  }
}

@media (max-width: 768px) {
  .dashboard-container {
    .stat-card .stat-content {
      gap: 12px;

      .stat-icon {
        width: 48px;
        height: 48px;

        :deep(.el-icon) {
          font-size: 24px;
        }
      }

      .stat-info .stat-value {
        font-size: 24px;
      }
    }

    .room-status-grid {
      grid-template-columns: repeat(auto-fill, minmax(70px, 1fr)) !important;
      gap: 8px !important;
    }
  }
}
</style>
