<template>
  <div class="home">
    <el-row :gutter="20">
      <el-col :span="6" v-for="(stat, i) in stats" :key="i">
        <el-card shadow="hover" class="stat-card">
          <div class="stat-content">
            <div class="stat-icon" :style="{ background: stat.color }">
              <el-icon :size="24"><component :is="stat.icon" /></el-icon>
            </div>
            <div class="stat-info">
              <div class="stat-value">{{ stat.value }}</div>
              <div class="stat-label">{{ stat.label }}</div>
            </div>
          </div>
        </el-card>
      </el-col>
    </el-row>

    <el-row :gutter="20" style="margin-top: 20px">
      <el-col :span="16">
        <el-card>
          <template #header>
            <span>推荐设备</span>
          </template>
          <el-row :gutter="20">
            <el-col :span="8" v-for="drone in drones" :key="drone.id">
              <el-card shadow="hover" class="drone-card" @click="goDetail(drone.id)">
                <div class="drone-image">
                  <el-image :src="drone.images?.split(',')[0]" fit="cover">
                    <template #error>
                      <div class="image-placeholder">
                        <el-icon :size="48"><Box /></el-icon>
                      </div>
                    </template>
                  </el-image>
                </div>
                <div class="drone-info">
                  <h4>{{ drone.name }}</h4>
                  <p>{{ drone.brand }} · {{ drone.model }}</p>
                  <div class="drone-meta">
                    <el-rate :model-value="drone.rating" disabled size="small" />
                    <span class="price">¥{{ drone.price_per_day }}/天</span>
                  </div>
                </div>
              </el-card>
            </el-col>
          </el-row>
        </el-card>
      </el-col>
      <el-col :span="8">
        <el-card>
          <template #header>
            <span>快捷操作</span>
          </template>
          <div class="quick-actions">
            <el-button type="primary" size="large" @click="$router.push('/drones')">
              浏览设备
            </el-button>
            <el-button type="success" size="large" @click="$router.push('/services')">
              航拍服务
            </el-button>
            <el-button v-if="role === 'client'" type="warning" size="large" @click="$router.push('/service/create')">
              发布需求
            </el-button>
            <el-button v-if="role === 'owner'" type="info" size="large" @click="$router.push('/drone/create')">
              添加设备
            </el-button>
          </div>
        </el-card>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRouter } from 'vue-router'
import { useUserStore } from '@/stores/user'
import request from '@/utils/request'

const router = useRouter()
const userStore = useUserStore()
const role = computed(() => userStore.role)

const drones = ref<Drone[]>([])

const stats = ref([
  { label: '在线设备', value: 0, icon: 'Box', color: '#409EFF' },
  { label: '进行中订单', value: 0, icon: 'List', color: '#67C23A' },
  { label: '开放服务', value: 0, icon: 'Service', color: '#E6A23C' },
  { label: '今日飞行', value: 0, icon: 'Location', color: '#F56C6C' }
])

onMounted(() => {
  fetchDrones()
})

async function fetchDrones() {
  try {
    const res: any = await request.get('/drones', { params: { status: 'online', page_size: 6 } })
    drones.value = res.data.list || []
  } catch (e) {
    console.error(e)
  }
}

function goDetail(id: number) {
  router.push(`/drone/${id}`)
}
</script>

<style scoped>
.home {
  padding: 0;
}
.stat-card {
  cursor: default;
}
.stat-content {
  display: flex;
  align-items: center;
  gap: 16px;
}
.stat-icon {
  width: 56px;
  height: 56px;
  border-radius: 12px;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #fff;
}
.stat-value {
  font-size: 24px;
  font-weight: bold;
  color: #303133;
}
.stat-label {
  font-size: 14px;
  color: #909399;
}
.drone-card {
  margin-bottom: 20px;
  cursor: pointer;
}
.drone-image {
  width: 100%;
  height: 160px;
  overflow: hidden;
  border-radius: 4px;
}
.image-placeholder {
  width: 100%;
  height: 100%;
  background: #f0f2f5;
  display: flex;
  align-items: center;
  justify-content: center;
  color: #c0c4cc;
}
.drone-info {
  padding-top: 12px;
}
.drone-info h4 {
  margin: 0 0 4px;
  color: #303133;
}
.drone-info p {
  margin: 0 0 8px;
  color: #909399;
  font-size: 13px;
}
.drone-meta {
  display: flex;
  justify-content: space-between;
  align-items: center;
}
.price {
  color: #f56c6c;
  font-weight: bold;
}
.quick-actions {
  display: flex;
  flex-direction: column;
  gap: 12px;
}
.quick-actions .el-button {
  width: 100%;
}
</style>
