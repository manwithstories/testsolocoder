<template>
  <div class="page">
    <div class="page-title">消息通知</div>
    <el-button style="margin-bottom:12px" @click="markAll">全部已读</el-button>
    <el-table :data="list" stripe>
      <el-table-column prop="title" label="标题" />
      <el-table-column prop="content" label="内容" />
      <el-table-column label="状态">
        <template #default="{ row }">
          <span v-if="!row.is_read" class="tag tag-waitlist">未读</span>
          <span v-else class="muted">已读</span>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="时间" />
      <el-table-column label="操作">
        <template #default="{ row }">
          <el-button v-if="!row.is_read" link type="primary" @click="mark(row.id)">标为已读</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination style="margin-top:16px" layout="prev,pager,next" :total="total" :page-size="pageSize" v-model:current-page="page" @current-change="fetch" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { msgApi } from '@/api'

const list = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

async function fetch() {
  const res = await msgApi.list({ page: page.value, page_size: pageSize.value })
  const data = res.data as any
  list.value = data.list || []
  total.value = data.total || 0
}

async function mark(id: number) {
  await msgApi.markRead(id)
  fetch()
}
async function markAll() {
  await msgApi.markAll()
  fetch()
}

onMounted(fetch)
</script>
