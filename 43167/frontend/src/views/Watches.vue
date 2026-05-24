<template>
  <div>
    <el-card>
      <template #header>
        <div class="header">
          <h3>手表市场</h3>
          <el-input v-model="keyword" placeholder="搜索品牌/型号" style="width: 220px" clearable />
          <el-input v-model="brand" placeholder="品牌" style="width: 140px" clearable />
          <el-button type="primary" @click="load">搜索</el-button>
        </div>
      </template>
      <el-row :gutter="16">
        <el-col v-for="w in list" :key="w.id" :span="6" style="margin-bottom: 16px">
          <el-card shadow="hover" @click="goDetail(w.id)">
            <img v-if="w.photos?.length" :src="w.photos[0].url" style="width:100%; height:180px; object-fit:cover" />
            <h4>{{ w.brand }} {{ w.model }}</h4>
            <div>¥{{ w.price.toFixed(2) }}</div>
            <div>年份 {{ w.year }} · {{ w.case_size_mm }}mm</div>
          </el-card>
        </el-col>
      </el-row>
    </el-card>
  </div>
</template>

<script setup lang="ts">
import { ref, onMounted } from 'vue'
import { useRouter } from 'vue-router'
import request from '@/utils/request'
import type { Watch } from '@/types'

const router = useRouter()
const list = ref<Watch[]>([])
const keyword = ref('')
const brand = ref('')

async function load() {
  const res: any = await request.get('/watches', { params: { keyword: keyword.value, brand: brand.value } })
  list.value = res.list || []
}

function goDetail(id: number) {
  router.push(`/watches/${id}`)
}

onMounted(load)
</script>

<style scoped>
.header { display: flex; align-items: center; gap: 8px; }
.header h3 { margin: 0; margin-right: auto; }
</style>
