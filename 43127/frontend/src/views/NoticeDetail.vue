<template>
  <div class="notice-detail">
    <div class="page-header">
      <h2 class="page-title">公告详情</h2>
      <el-button @click="goBack">返回</el-button>
    </div>

    <el-row :gutter="20">
      <el-col :span="16">
        <div class="card">
          <div class="detail-header">
            <div>
              <el-tag v-if="notice?.isTop" type="danger" style="margin-right: 8px;">置顶</el-tag>
              <el-tag v-if="notice?.type === 2" type="warning" style="margin-right: 8px;">紧急</el-tag>
              <h1>{{ notice?.title }}</h1>
            </div>
          </div>

          <div class="meta">
            <span>发布人: {{ notice?.publisher?.realName }}</span>
            <span>发布时间: {{ formatDateTime(notice?.createdAt) }}</span>
            <span>推送范围: {{ notice?.building || '全部' }}</span>
          </div>

          <div class="content">
            <p>{{ notice?.content }}</p>
          </div>
        </div>
      </el-col>

      <el-col :span="8">
        <div class="card">
          <h3>公告信息</h3>
          <el-descriptions :column="1" border>
            <el-descriptions-item label="类型">
              <el-tag :type="notice?.type === 2 ? 'warning' : 'info'">
                {{ notice?.type === 2 ? '紧急通知' : '普通公告' }}
              </el-tag>
            </el-descriptions-item>
            <el-descriptions-item label="状态">
              <el-tag :type="notice?.status === 1 ? 'success' : 'info'">
                {{ notice?.status === 1 ? '已发布' : '已下架' }}
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
import type { Notice } from '@/types'
import { getNotice } from '@/api/business'
import dayjs from 'dayjs'

const route = useRoute()
const router = useRouter()

const notice = ref<Notice | null>(null)

onMounted(async () => {
  try {
    const res = await getNotice(Number(route.params.id))
    notice.value = res.data
  } catch (error) {
    console.error('Failed to load notice:', error)
  }
})

function formatDateTime(date?: string) {
  return date ? dayjs(date).format('YYYY-MM-DD HH:mm') : '-'
}

function goBack() {
  router.back()
}
</script>

<style scoped>
.notice-detail {
  padding: 0;
}

.detail-header {
  margin-bottom: 20px;
  padding-bottom: 15px;
  border-bottom: 1px solid #ebeef5;
}

.detail-header h1 {
  font-size: 20px;
  margin: 10px 0 0 0;
}

.meta {
  display: flex;
  gap: 20px;
  color: #909399;
  font-size: 13px;
  margin-bottom: 20px;
}

.content p {
  color: #606266;
  line-height: 1.8;
  font-size: 14px;
}
</style>
