<template>
  <div class="page">
    <div class="page-title">{{ event?.name || '加载中...' }}</div>
    <div v-if="event" class="card">
      <p><b>地点：</b>{{ event.location }}</p>
      <p><b>时间：</b>{{ formatDate(event.start_date) }} - {{ formatDate(event.end_date) }}</p>
      <p><b>报名截止：</b>{{ formatDate(event.registration_deadline) }}</p>
      <p><b>主办方：</b>{{ event.organizer }}</p>
      <p v-if="event.description"><b>简介：</b>{{ event.description }}</p>
    </div>
    <div v-if="event" class="card">
      <h3>赛事项目</h3>
      <el-table :data="event.items || []">
        <el-table-column prop="name" label="项目名称" />
        <el-table-column prop="category" label="类别" />
        <el-table-column label="年龄">
          <template #default="{ row }">
            {{ row.min_age || 0 }} - {{ row.max_age || '不限' }}
          </template>
        </el-table-column>
        <el-table-column label="名额">
          <template #default="{ row }">
            {{ row.quota }} / 候补 {{ row.waitlist_quota || 0 }}
          </template>
        </el-table-column>
        <el-table-column prop="fee" label="费用" />
        <el-table-column label="操作">
          <template #default="{ row }">
            <el-button size="small" type="primary" @click="register(row)">报名</el-button>
          </template>
        </el-table-column>
      </el-table>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import { eventApi, regApi } from '@/api'

const route = useRoute()
const event = ref<any>(null)

function formatDate(d: string) { return d ? new Date(d).toLocaleString() : '' }

async function register(item: any) {
  try {
    await ElMessageBox.confirm(`确定报名【${item.name}】吗？`, '提示')
    await regApi.create({ event_item_id: item.id, reg_type: 'individual' })
    ElMessage.success('报名成功')
  } catch (_) { /* ignore cancel */ }
}

onMounted(async () => {
  const id = Number(route.params.id)
  const res = await eventApi.get(id)
  event.value = res.data
})
</script>
