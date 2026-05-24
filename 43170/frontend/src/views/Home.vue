<template>
  <div class="home-page">
    <div class="banner">
      <div class="banner-content">
        <h1>专业摄影器材租赁平台</h1>
        <p>为摄影爱好者和专业摄影师提供高品质设备租赁服务</p>
        <el-button type="primary" size="large" @click="goToEquipments">立即浏览设备</el-button>
      </div>
    </div>

    <div class="page-container">
      <h2 class="section-title">热门设备</h2>
      <el-row :gutter="20">
        <el-col
          v-for="equipment in equipments"
          :key="equipment.id"
          :xs="24"
          :sm="12"
          :md="8"
          :lg="6"
        >
          <div class="card equipment-card" @click="goToDetail(equipment.id)">
            <el-carousel
              v-if="equipment.images && equipment.images.length > 0"
              height="200px"
              :indicator-position="'none'"
            >
              <el-carousel-item v-for="img in equipment.images" :key="img.id">
                <img
                  :src="`/uploads/${img.imageUrl}`"
                  :alt="equipment.name"
                  class="equipment-image"
                />
              </el-carousel-item>
            </el-carousel>
            <div v-else class="equipment-image-placeholder">
              <el-icon :size="48"><Camera /></el-icon>
            </div>
            <div class="equipment-info">
              <div class="equipment-name">{{ equipment.name }}</div>
              <div class="equipment-brand">{{ equipment.brand }} {{ equipment.model }}</div>
              <div class="equipment-rating">
                <el-rate :model-value="equipment.rating" disabled size="small" />
                <span>({{ equipment.reviewCount }})</span>
              </div>
              <div class="equipment-price">¥{{ equipment.dailyRent }}/天</div>
              <div class="equipment-deposit">押金: ¥{{ equipment.deposit }}</div>
            </div>
          </div>
        </el-col>
      </el-row>

      <div v-if="equipments.length === 0" class="empty-state">
        <el-icon><Box /></el-icon>
        <p>暂无设备</p>
      </div>
    </div>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import { searchApi } from '@/api/equipment'
import type { Equipment } from '@/types'

const router = useRouter()
const equipments = ref<Equipment[]>([])

onMounted(async () => {
  try {
    const response = await searchApi.searchEquipments({
      pageSize: 8,
      sortBy: 'rating',
      sortOrder: 'desc'
    })
    equipments.value = response.data
  } catch (error) {
    console.error('Failed to load equipments:', error)
  }
})

function goToEquipments() {
  router.push('/equipments')
}

function goToDetail(id: number) {
  router.push(`/equipments/${id}`)
}
</script>

<style scoped>
.banner {
  background: linear-gradient(135deg, #667eea 0%, #764ba2 100%);
  color: #fff;
  padding: 60px 20px;
  text-align: center;
}

.banner-content h1 {
  font-size: 36px;
  margin-bottom: 16px;
}

.banner-content p {
  font-size: 18px;
  margin-bottom: 24px;
  opacity: 0.9;
}

.equipment-image {
  width: 100%;
  height: 200px;
  object-fit: cover;
}

.equipment-image-placeholder {
  width: 100%;
  height: 200px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
  color: #c0c4cc;
}

.equipment-brand {
  font-size: 14px;
  color: #909399;
  margin-bottom: 8px;
}

.equipment-deposit {
  font-size: 12px;
  color: #909399;
  margin-top: 4px;
}
</style>
