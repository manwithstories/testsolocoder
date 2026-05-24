<template>
  <el-card v-if="watch">
    <el-row :gutter="24">
      <el-col :span="12">
        <el-carousel height="400px">
          <el-carousel-item v-for="p in watch.photos || []" :key="p.id">
            <img :src="p.url" style="width:100%; height:100%; object-fit:contain" />
          </el-carousel-item>
        </el-carousel>
      </el-col>
      <el-col :span="12">
        <h2>{{ watch.brand }} {{ watch.model }}</h2>
        <p>参考号: {{ watch.reference_no }}</p>
        <p>年份: {{ watch.year }}</p>
        <p>机芯: {{ watch.movement }}</p>
        <p>表径: {{ watch.case_size_mm }}mm</p>
        <p>表壳: {{ watch.case_material }}</p>
        <p>表盘: {{ watch.dial_color }}</p>
        <p>表带: {{ watch.bracelet }}</p>
        <p>成色: {{ watch.condition }}</p>
        <p>价格: ¥{{ watch.price.toFixed(2) }}</p>
        <p>状态: {{ watch.status }} <el-tag v-if="watch.authed" type="success">已鉴定</el-tag></p>
        <p>描述: {{ watch.description }}</p>
        <el-button type="primary" @click="favorite">收藏</el-button>
        <el-button type="success" @click="bid">出价</el-button>
      </el-col>
    </el-row>
  </el-card>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRoute } from 'vue-router'
import { ElMessage, ElMessageBox } from 'element-plus'
import request from '@/utils/request'
import type { Watch } from '@/types'

const route = useRoute()
const watch = ref<Watch | null>(null)

onMounted(async () => {
  const id = route.params.id
  watch.value = await request.get(`/watches/${id}`)
})

async function favorite() {
  if (!watch.value) return
  await request.post('/favorites', { watch_id: watch.value.id })
  ElMessage.success('已收藏')
}

async function bid() {
  try {
    const { value } = await ElMessageBox.prompt('请输入出价金额', '出价', {
      inputPattern: /\d+(\.\d+)?/,
      inputErrorMessage: '金额格式不正确'
    })
    await request.post('/trades')
    ElMessage.info('请前往交易模块查看开放交易并出价')
  } catch {}
}
</script>
