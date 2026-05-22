<template>
  <div class="page">
    <div class="page-title">我的证书</div>
    <div v-if="list.length === 0" class="card muted">暂无证书，请先参赛并取得成绩</div>
    <el-table v-else :data="list" stripe>
      <el-table-column prop="certificate_no" label="证书编号" />
      <el-table-column prop="certificate_name" label="证书名称" />
      <el-table-column prop="rank" label="名次" width="80" />
      <el-table-column prop="score" label="成绩" width="100" />
      <el-table-column label="状态" width="100">
        <template #default="{ row }">
          <span :class="['tag', row.status === 'generated' ? 'tag-confirmed' : 'tag-waitlist']">
            {{ statusText(row.status) }}
          </span>
        </template>
      </el-table-column>
      <el-table-column prop="created_at" label="生成时间" />
      <el-table-column label="操作" width="120">
        <template #default="{ row }">
          <el-button
            v-if="row.status === 'generated'"
            link type="primary"
            @click="download(row)"
          >下载</el-button>
          <el-button
            v-else-if="row.status === 'failed' && row.retry_count < 3"
            link type="warning"
            @click="retry(row)"
          >重试</el-button>
          <span v-else class="muted">{{ statusText(row.status) }}</span>
        </template>
      </el-table-column>
    </el-table>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { certApi } from '@/api'

const list = ref<any[]>([])

function statusText(s: string) {
  const map: Record<string, string> = {
    generated: '已生成', generating: '生成中', failed: '失败',
  }
  return map[s] || s
}

function download(row: any) {
  window.open(certApi.download(row.id), '_blank')
}

async function retry(row: any) {
  await certApi.generate(row.score_id || row.id)
  ElMessage.success('已提交重试请求')
  fetch()
}

async function fetch() {
  const res = await certApi.my()
  list.value = (res.data as any[]) || []
}

onMounted(fetch)
</script>
