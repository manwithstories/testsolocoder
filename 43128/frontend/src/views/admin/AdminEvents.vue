<template>
  <div class="page">
    <div class="page-title">赛事管理</div>
    <div class="card" style="display:flex;justify-content:space-between;align-items:center">
      <el-button type="primary" @click="$router.push('/admin/event/create')">+ 创建赛事</el-button>
      <el-input v-model="keyword" placeholder="搜索赛事名称" style="width:240px" clearable @clear="fetch" @keyup.enter="fetch" />
    </div>
    <el-table :data="list" stripe>
      <el-table-column prop="id" label="ID" width="70" />
      <el-table-column prop="name" label="赛事名称" min-width="180">
        <template #default="{ row }">
          <router-link :to="`/events/${row.id}`">{{ row.name }}</router-link>
        </template>
      </el-table-column>
      <el-table-column prop="location" label="地点" />
      <el-table-column label="时间" min-width="200">
        <template #default="{ row }">
          {{ fmt(row.start_date) }} ~ {{ fmt(row.end_date) }}
        </template>
      </el-table-column>
      <el-table-column label="报名截止">
        <template #default="{ row }">{{ fmt(row.registration_deadline) }}</template>
      </el-table-column>
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <span :class="['tag', row.is_published ? 'tag-published' : 'tag-draft']">
            {{ row.is_published ? '已发布' : '草稿' }}
          </span>
        </template>
      </el-table-column>
      <el-table-column label="操作" width="260">
        <template #default="{ row }">
          <el-button link type="primary" @click="$router.push(`/admin/event/${row.id}/edit`)">编辑</el-button>
          <el-button
            v-if="!row.is_published"
            link type="success"
            @click="publish(row)"
          >发布</el-button>
          <el-button
            v-else
            link type="warning"
            @click="unpublish(row)"
          >取消发布</el-button>
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
import { ElMessage, ElMessageBox } from 'element-plus'
import { eventApi } from '@/api'

const list = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)
const keyword = ref('')

function fmt(d: string) { return d ? new Date(d).toLocaleDateString() : '-' }

async function fetch() {
  const res = await eventApi.listAll({ page: page.value, page_size: pageSize.value })
  const data = res.data as any
  list.value = data.list || []
  total.value = data.total || 0
}

async function publish(row: any) {
  await ElMessageBox.confirm(`确定发布赛事【${row.name}】？`, '提示')
  await eventApi.publish(row.id)
  ElMessage.success('已发布')
  fetch()
}

async function unpublish(row: any) {
  await ElMessageBox.confirm(`确定取消发布赛事【${row.name}】？`, '提示')
  await eventApi.unpublish(row.id)
  ElMessage.success('已取消发布')
  fetch()
}

onMounted(fetch)
</script>
