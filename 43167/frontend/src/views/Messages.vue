<template>
  <el-card>
    <template #header>
      <div style="display:flex; justify-content:space-between; align-items:center">
        <span>消息中心</span>
        <el-button @click="markAll">全部已读</el-button>
      </div>
    </template>
    <el-table :data="list">
      <el-table-column prop="title" label="标题" />
      <el-table-column prop="content" label="内容" />
      <el-table-column prop="read" label="状态">
        <template #default="{ row }">
          <el-tag :type="row.read ? 'info' : 'warning'">{{ row.read ? '已读' : '未读' }}</el-tag>
        </template>
      </el-table-column>
      <el-table-column label="操作">
        <template #default="{ row }">
          <el-button v-if="!row.read" size="small" @click="read(row.id)">标为已读</el-button>
        </template>
      </el-table-column>
    </el-table>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import request from '@/utils/request'

const list = ref<any[]>([])

async function loadList() {
  const res: any = await request.get('/messages')
  list.value = res.list || res || []
}

onMounted(loadList)

async function read(id: number) {
  await request.post(`/messages/${id}/read`)
  loadList()
}
async function markAll() {
  await request.post('/messages/read-all')
  list.value = list.value.map((x) => ({ ...x, read: true }))
}
</script>
