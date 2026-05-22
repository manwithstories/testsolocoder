<template>
  <div class="page">
    <div class="page-title">赛事列表</div>
    <el-table :data="list" stripe>
      <el-table-column prop="name" label="赛事名称" min-width="180">
        <template #default="{ row }">
          <router-link :to="`/events/${row.id}`" style="color:#409eff;font-weight:500">{{ row.name }}</router-link>
        </template>
      </el-table-column>
      <el-table-column prop="location" label="地点" />
      <el-table-column label="时间" min-width="200">
        <template #default="{ row }">
          {{ formatDate(row.start_date) }} ~ {{ formatDate(row.end_date) }}
        </template>
      </el-table-column>
      <el-table-column label="报名截止">
        <template #default="{ row }">{{ formatDate(row.registration_deadline) }}</template>
      </el-table-column>
      <el-table-column prop="organizer" label="主办方" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <span :class="['tag', row.is_published ? 'tag-published' : 'tag-draft']">
            {{ row.is_published ? '已发布' : '草稿' }}
          </span>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination
      style="margin-top:16px"
      layout="total, prev, pager, next"
      :total="total"
      :page-size="pageSize"
      v-model:current-page="page"
      @current-change="fetch"
    />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { eventApi } from '@/api'

const list = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

function formatDate(d: string) {
  return d ? new Date(d).toLocaleDateString() : '-'
}

async function fetch() {
  const res = await eventApi.list({ page: page.value, page_size: pageSize.value })
  const data = res.data as any
  list.value = data.list || []
  total.value = data.total || 0
}

onMounted(fetch)
</script>
