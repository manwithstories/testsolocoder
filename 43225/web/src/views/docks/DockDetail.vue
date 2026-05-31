<template>
  <div class="dock-detail" v-if="dock">
    <el-page-header @back="$router.back()" :content="dock.name" class="page-header" />

    <el-row :gutter="16">
      <el-col :xs="24" :md="16">
        <div class="card-container">
          <img :src="dock.image_url || '/placeholder-dock.jpg'" class="dock-image" />
        </div>

        <div class="card-container" style="margin-top: 16px">
          <h3>码头信息</h3>
          <el-descriptions :column="2" border>
            <el-descriptions-item label="地址">{{ dock.address }}</el-descriptions-item>
            <el-descriptions-item label="城市">{{ dock.city || '-' }}</el-descriptions-item>
            <el-descriptions-item label="国家">{{ dock.country || '-' }}</el-descriptions-item>
            <el-descriptions-item label="营业时间">
              {{ dock.open_time || '-' }} - {{ dock.close_time || '-' }}
            </el-descriptions-item>
            <el-descriptions-item label="设施">
              <el-tag v-for="amenity in dock.amenities?.split(',') || []" :key="amenity" size="small">
                {{ amenity.trim() }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="评分">
              <el-rate :model-value="dock.average_rating" disabled />
              <span>({{ dock.review_count }}条评价)</span>
            </el-descriptions-item>
          </el-descriptions>

          <h4 style="margin-top: 20px">描述</h4>
          <p>{{ dock.description || '暂无描述' }}</p>
        </div>

        <div class="card-container" style="margin-top: 16px">
          <div class="section-header">
            <h3>泊位列表</h3>
            <el-button type="primary" size="small" @click="$router.push('/berth-reservations')">
              预约泊位
            </el-button>
          </div>
          <el-table :data="dock.berths || []" style="width: 100%">
            <el-table-column prop="number" label="泊位号" width="100" />
            <el-table-column prop="berth_type" label="类型" width="100">
              <template #default="{ row }">
                {{ getBerthTypeText(row.berth_type) }}
              </template>
            </el-table-column>
            <el-table-column label="尺寸">
              <template #default="{ row }">
                {{ row.max_length }}m × {{ row.max_width }}m
              </template>
            </el-table-column>
            <el-table-column label="设施">
              <template #default="{ row }">
                <el-tag v-if="row.has_water" size="small" type="success">水</el-tag>
                <el-tag v-if="row.has_electric" size="small" type="primary">电</el-tag>
                <el-tag v-if="row.has_internet" size="small" type="warning">网络</el-tag>
              </template>
            </el-table-column>
            <el-table-column label="价格">
              <template #default="{ row }">
                ¥{{ row.daily_rate }}/天 或 ¥{{ row.hourly_rate }}/小时
              </template>
            </el-table-column>
            <el-table-column label="状态" width="100">
              <template #default="{ row }">
                <el-tag :type="row.status === 'available' ? 'success' : 'warning'" size="small">
                  {{ getStatusText(row.status) }}
                </el-tag>
              </template>
            </el-table-column>
          </el-table>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { getDockApi } from '@/api/berth'
import type { Dock } from '@/types/berth'

const route = useRoute()
const dock = ref<Dock | null>(null)

const getBerthTypeText = (type: string) => {
  const map: Record<string, string> = {
    small: '小型',
    medium: '中型',
    large: '大型'
  }
  return map[type] || type
}

const getStatusText = (status: string) => {
  const map: Record<string, string> = {
    available: '可用',
    occupied: '已占用',
    reserved: '已预约',
    maintenance: '维护中'
  }
  return map[status] || status
}

onMounted(async () => {
  const id = route.params.id as string
  try {
    const res: any = await getDockApi(id)
    dock.value = res.data
  } catch (error) {
    console.error('Failed to fetch dock:', error)
  }
})
</script>

<style lang="scss" scoped>
.dock-detail {
  .page-header {
    margin-bottom: 16px;
    padding: 12px 20px;
    background: #fff;
    border-radius: 8px;
  }

  .card-container {
    background: #fff;
    border-radius: 8px;
    padding: 20px;

    .dock-image {
      width: 100%;
      max-height: 300px;
      object-fit: cover;
      border-radius: 8px;
    }

    .section-header {
      display: flex;
      justify-content: space-between;
      align-items: center;
      margin-bottom: 16px;
    }
  }
}
</style>
