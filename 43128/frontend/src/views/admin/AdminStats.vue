<template>
  <div class="page">
    <div class="page-title">数据统计与导出</div>
    <div class="card">
      <el-button type="primary" @click="load">刷新统计</el-button>
      <el-button :href="statsApi.export" link type="primary">导出报表</el-button>
    </div>
    <div class="card">
      <pre>{{ JSON.stringify(data, null, 2) }}</pre>
    </div>
  </div>
</template>

<script setup lang="ts">
import { onMounted, ref } from 'vue'
import { statsApi } from '@/api'

const data = ref<any>({})

async function load() {
  const res = await statsApi.overview({})
  data.value = res.data
}

onMounted(load)
</script>
