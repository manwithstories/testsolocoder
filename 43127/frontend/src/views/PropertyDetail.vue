<template>
  <div class="property-detail">
    <div class="page-header">
      <h2 class="page-title">房源详情</h2>
      <el-button @click="goBack">返回</el-button>
    </div>

    <el-row :gutter="20">
      <el-col :span="16">
        <div class="card">
          <el-carousel height="400px" v-if="property?.images?.length">
            <el-carousel-item v-for="(img, index) in property.images" :key="index">
              <img :src="img.url" style="width: 100%; height: 100%; object-fit: cover" />
            </el-carousel-item>
          </el-carousel>
          <div v-else class="no-image">暂无图片</div>

          <div class="property-info">
            <h1>{{ property?.title }}</h1>
            <div class="info-row">
              <span class="price">¥{{ property?.rent }}/月</span>
              <span class="deposit">押金: ¥{{ property?.deposit }}</span>
            </div>
            <div class="meta">
              <span>{{ property?.layout }}</span>
              <span>{{ property?.area }}㎡</span>
              <span>{{ property?.region }}</span>
              <span>{{ property?.community }}</span>
            </div>
            <div class="description">
              <h3>房源描述</h3>
              <p>{{ property?.description }}</p>
            </div>
            <div class="facilities" v-if="property?.facilities?.length">
              <h3>配套设施</h3>
              <div class="facility-list">
                <span v-for="f in property.facilities" :key="f.id" class="facility-tag">
                  {{ f.icon }} {{ f.name }}
                </span>
              </div>
            </div>
          </div>
        </div>
      </el-col>

      <el-col :span="8">
        <div class="card">
          <h3>房源信息</h3>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="小区">{{ property?.community }}</el-descriptions-item>
            <el-descriptions-item label="地址">{{ property?.address }}</el-descriptions-item>
            <el-descriptions-item label="楼层">{{ property?.floor }}</el-descriptions-item>
            <el-descriptions-item label="楼栋">{{ property?.building }}</el-descriptions-item>
            <el-descriptions-item label="房间号">{{ property?.roomNo }}</el-descriptions-item>
            <el-descriptions-item label="付款方式">{{ property?.paymentType }}</el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="getStatusType(property?.status)">
                {{ getStatusText(property?.status) }}
              </el-tag>
            </el-descriptions-item>
          </el-descriptions>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { Property } from '@/types'
import { getProperty } from '@/api/property'

const route = useRoute()
const router = useRouter()

const property = ref<Property | null>(null)

onMounted(async () => {
  try {
    const res = await getProperty(Number(route.params.id))
    property.value = res.data
  } catch (error) {
    console.error('Failed to load property:', error)
  }
})

function getStatusType(status?: number) {
  switch (status) {
    case 1: return 'success'
    case 2: return 'warning'
    default: return 'info'
  }
}

function getStatusText(status?: number) {
  switch (status) {
    case 1: return '上架'
    case 2: return '已租'
    default: return '下架'
  }
}

function goBack() {
  router.back()
}
</script>

<style scoped>
.property-detail {
  padding: 0;
}

.no-image {
  height: 400px;
  display: flex;
  align-items: center;
  justify-content: center;
  background: #f5f7fa;
  color: #909399;
}

.property-info {
  padding: 20px;
}

.property-info h1 {
  font-size: 24px;
  margin-bottom: 15px;
}

.info-row {
  display: flex;
  align-items: baseline;
  gap: 20px;
  margin-bottom: 15px;
}

.price {
  font-size: 28px;
  color: #f56c6c;
  font-weight: 600;
}

.deposit {
  color: #909399;
}

.meta {
  display: flex;
  gap: 20px;
  color: #606266;
  margin-bottom: 20px;
}

.description h3,
.facilities h3 {
  font-size: 16px;
  margin-bottom: 10px;
}

.description p {
  color: #606266;
  line-height: 1.6;
}

.facility-list {
  display: flex;
  flex-wrap: wrap;
  gap: 10px;
}

.facility-tag {
  padding: 4px 12px;
  background: #ecf5ff;
  border-radius: 4px;
  font-size: 13px;
}
</style>
