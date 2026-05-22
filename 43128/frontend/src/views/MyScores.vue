<template>
  <div class="page">
    <div class="page-title">我的成绩</div>
    <el-table :data="list" stripe>
      <el-table-column prop="score" label="成绩" />
      <el-table-column prop="rank" label="排名" />
      <el-table-column prop="points" label="积分" />
      <el-table-column prop="time_used" label="用时" />
      <el-table-column prop="remarks" label="备注" />
      <el-table-column prop="created_at" label="录入时间" />
      <el-table-column label="操作">
        <template #default="{ row }">
          <el-button link type="primary" @click="genCert(row)">生成证书</el-button>
        </template>
      </el-table-column>
    </el-table>
    <el-pagination style="margin-top:16px" layout="prev,pager,next" :total="total" :page-size="pageSize" v-model:current-page="page" @current-change="fetch" />
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { ElMessage } from 'element-plus'
import { scoreApi, certApi } from '@/api'

const list = ref<any[]>([])
const total = ref(0)
const page = ref(1)
const pageSize = ref(10)

async function fetch() {
  const res = await scoreApi.my({ page: page.value, page_size: pageSize.value })
  const data = res.data as any
  list.value = data.list || []
  total.value = data.total || 0
}

async function genCert(row: any) {
  await certApi.generate(row.id)
  ElMessage.success('已提交生成请求')
}

onMounted(fetch)
</script>
