<template>
  <div class="page-container">
    <div class="page-header">
      <el-button :icon="ArrowLeft" @click="goBack">返回</el-button>
      <h2>{{ wedding?.title || '婚礼详情' }}</h2>
    </div>

    <el-tabs v-model="activeTab" type="border-card">
      <el-tab-pane label="基本信息" name="info">
        <el-descriptions :column="2" border>
          <el-descriptions-item label="婚礼名称">{{ wedding?.title }}</el-descriptions-item>
          <el-descriptions-item label="婚礼日期">{{ formatDate(wedding?.wedding_date) }}</el-descriptions-item>
          <el-descriptions-item label="新郎">{{ wedding?.groom_name }}</el-descriptions-item>
          <el-descriptions-item label="新娘">{{ wedding?.bride_name }}</el-descriptions-item>
          <el-descriptions-item label="婚礼风格">{{ wedding?.style || '-' }}</el-descriptions-item>
          <el-descriptions-item label="主题色">{{ wedding?.theme_color || '-' }}</el-descriptions-item>
          <el-descriptions-item label="场地">{{ wedding?.venue || '-' }}</el-descriptions-item>
          <el-descriptions-item label="预算">¥{{ formatNumber(wedding?.budget) }}</el-descriptions-item>
          <el-descriptions-item label="预计人数">{{ wedding?.guest_count }}</el-descriptions-item>
          <el-descriptions-item label="状态">
            <el-tag :type="statusType(wedding?.status)">{{ statusText(wedding?.status) }}</el-tag>
          </el-descriptions-item>
          <el-descriptions-item label="场地地址" :span="2">{{ wedding?.venue_address || '-' }}</el-descriptions-item>
          <el-descriptions-item label="备注" :span="2">{{ wedding?.description || '-' }}</el-descriptions-item>
        </el-descriptions>
      </el-tab-pane>

      <el-tab-pane label="供应商" name="vendors">
        <VendorList :wedding-id="weddingId" />
      </el-tab-pane>

      <el-tab-pane label="嘉宾管理" name="guests">
        <GuestList :wedding-id="weddingId" />
      </el-tab-pane>

      <el-tab-pane label="预算管理" name="budget">
        <BudgetManagement :wedding-id="weddingId" />
      </el-tab-pane>

      <el-tab-pane label="任务清单" name="tasks">
        <TaskList :wedding-id="weddingId" />
      </el-tab-pane>

      <el-tab-pane label="合同文档" name="documents">
        <DocumentList :wedding-id="weddingId" />
      </el-tab-pane>
    </el-tabs>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted, computed } from 'vue'
import { useRoute, useRouter } from 'vue-router'
import { useWeddingStore } from '@/store/wedding'
import { weddingApi } from '@/api/wedding'
import { ArrowLeft } from '@element-plus/icons-vue'
import dayjs from 'dayjs'
import type { Wedding } from '@/types'
import VendorList from '@/views/VendorList.vue'
import GuestList from '@/views/GuestList.vue'
import BudgetManagement from '@/views/Budget.vue'
import TaskList from '@/views/TaskList.vue'
import DocumentList from '@/views/DocumentList.vue'

const route = useRoute()
const router = useRouter()
const weddingStore = useWeddingStore()

const weddingId = computed(() => Number(route.params.id))
const wedding = ref<Wedding | null>(null)
const activeTab = ref('info')

function formatDate(date?: string) {
  return date ? dayjs(date).format('YYYY-MM-DD') : '-'
}

function formatNumber(num?: number) {
  return num?.toLocaleString() || '0'
}

function statusType(status?: string) {
  const types: Record<string, string> = {
    planning: 'warning',
    confirmed: 'success',
    completed: 'primary',
    cancelled: 'danger'
  }
  return types[status || ''] || 'info'
}

function statusText(status?: string) {
  const texts: Record<string, string> = {
    planning: '筹备中',
    confirmed: '已确认',
    completed: '已完成',
    cancelled: '已取消'
  }
  return texts[status || ''] || status || '-'
}

async function fetchWedding() {
  try {
    if (weddingStore.currentWedding?.id === weddingId.value) {
      wedding.value = weddingStore.currentWedding
    } else {
      const res = await weddingApi.getById(weddingId.value)
      wedding.value = res.data
      weddingStore.setCurrentWedding(res.data)
    }
  } catch (error) {
    console.error('Failed to fetch wedding:', error)
  }
}

function goBack() {
  router.push('/weddings')
}

onMounted(fetchWedding)
</script>

<style scoped>
.page-container {
  padding: 0;
}

.page-header {
  display: flex;
  align-items: center;
  gap: 16px;
  margin-bottom: 20px;
}

.page-header h2 {
  margin: 0;
}
</style>
