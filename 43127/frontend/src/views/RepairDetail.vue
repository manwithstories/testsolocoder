<template>
  <div class="repair-detail">
    <div class="page-header">
      <h2 class="page-title">工单详情</h2>
      <el-button @click="goBack">返回</el-button>
    </div>

    <el-row :gutter="20">
      <el-col :span="16">
        <div class="card">
          <div class="detail-header">
            <h1>{{ order?.title }}</h1>
            <el-tag :type="getStatusType(order?.status)" size="large">
              {{ getStatusText(order?.status) }}
            </el-tag>
          </div>

          <el-descriptions :column="2" border>
            <el-descriptions-item label="分类">{{ order?.category }}</el-descriptions-item>
            <el-descriptions-item label="优先级">
              <el-tag :type="getPriorityType(order?.priority)">
                {{ getPriorityText(order?.priority) }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="房源">{{ order?.property?.title }}</el-descriptions-item>
            <el-descriptions-item label="租户">{{ order?.tenant?.name }}</el-descriptions-item>
            <el-descriptions-item label="处理人">{{ order?.handler?.realName || '-' }}</el-descriptions-item>
            <el-descriptions-item label="创建时间">
              {{ formatDate(order?.createdAt) }}
            </el-descriptions-item>
            <el-descriptions-item label="完成时间" v-if="order?.completedAt">
              {{ formatDate(order?.completedAt) }}
            </el-descriptions-item>
          </el-descriptions>

          <div class="description">
            <h3>问题描述</h3>
            <p>{{ order?.description }}</p>
          </div>

          <div class="process-note" v-if="order?.processNote">
            <h3>处理备注</h3>
            <p>{{ order?.processNote }}</p>
          </div>

          <div class="images" v-if="order?.images">
            <h3>现场照片</h3>
            <div class="image-list">
              <el-image
                v-for="(img, index) in imageList"
                :key="index"
                :src="img"
                :preview-src-list="imageList"
                fit="cover"
                style="width: 150px; height: 150px; margin-right: 10px; border-radius: 4px"
              />
            </div>
          </div>
        </div>
      </el-col>

      <el-col :span="8">
        <div class="card">
          <h3>工单进度</h3>
          <el-steps :active="order?.status || 0" finish-status="success" align-center>
            <el-step title="待处理" />
            <el-step title="处理中" />
            <el-step title="已完成" />
            <el-step title="已关闭" />
          </el-steps>
        </div>
      </el-col>
    </el-row>
  </div>
</template>

<script setup lang="ts">
import { ref, computed, onMounted } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import type { RepairOrder } from '@/types'
import { getRepair } from '@/api/business'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()

const order = ref<RepairOrder | null>(null)

const imageList = computed(() => {
  if (!order.value?.images) return []
  return order.value.images.split(',').filter(Boolean)
})

onMounted(async () => {
  try {
    const res = await getRepair(Number(route.params.id))
    order.value = res.data
  } catch (error) {
    console.error('Failed to load repair:', error)
  }
})

function formatDate(date?: string) {
  return date ? dayjs(date).format('YYYY-MM-DD HH:mm') : '-'
}

function getStatusType(status?: number) {
  switch (status) {
    case 1: return 'warning'
    case 2: return 'primary'
    case 3: return 'success'
    default: return 'info'
  }
}

function getStatusText(status?: number) {
  switch (status) {
    case 1: return '待处理'
    case 2: return '处理中'
    case 3: return '已完成'
    default: return '已关闭'
  }
}

function getPriorityType(priority?: number) {
  switch (priority) {
    case 1: return 'danger'
    case 2: return 'warning'
    default: return 'info'
  }
}

function getPriorityText(priority?: number) {
  switch (priority) {
    case 1: return '紧急'
    case 2: return '高'
    default: return '普通'
  }
}

function goBack() {
  router.back()
}
</script>

<style scoped>
.repair-detail {
  padding: 0;
}

.detail-header {
  display: flex;
  justify-content: space-between;
  align-items: center;
  margin-bottom: 20px;
}

.detail-header h1 {
  font-size: 20px;
  margin: 0;
}

.description,
.process-note,
.images {
  margin-top: 20px;
}

.description h3,
.process-note h3,
.images h3 {
  font-size: 16px;
  margin-bottom: 10px;
}

.description p,
.process-note p {
  color: #606266;
  line-height: 1.6;
}

.image-list {
  display: flex;
  flex-wrap: wrap;
}
</style>
